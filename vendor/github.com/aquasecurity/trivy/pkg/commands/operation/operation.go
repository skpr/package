package operation

import (
	"context"
	"crypto/tls"
	"os"
	"strings"

	"github.com/samber/lo"

	"github.com/aquasecurity/trivy/pkg/flag"

	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"golang.org/x/xerrors"

	"github.com/aquasecurity/trivy-db/pkg/metadata"
	"github.com/aquasecurity/trivy/pkg/db"
	"github.com/aquasecurity/trivy/pkg/fanal/cache"
	"github.com/aquasecurity/trivy/pkg/log"
	"github.com/aquasecurity/trivy/pkg/utils"
)

// SuperSet binds cache dependencies
var SuperSet = wire.NewSet(
	cache.NewFSCache,
	wire.Bind(new(cache.LocalArtifactCache), new(cache.FSCache)),
	NewCache,
)

// Cache implements the local cache
type Cache struct {
	cache.Cache
}

// NewCache is the factory method for Cache
func NewCache(c flag.CacheOptions) (Cache, error) {
	if strings.HasPrefix(c.CacheBackend, "redis://") {
		log.Logger.Infof("Redis cache: %s", c.CacheBackendMasked())
		options, err := redis.ParseURL(c.CacheBackend)
		if err != nil {
			return Cache{}, err
		}

		if !lo.IsEmpty(c.RedisOptions) {
			caCert, cert, err := utils.GetTLSConfig(c.RedisCACert, c.RedisCert, c.RedisKey)
			if err != nil {
				return Cache{}, err
			}

			options.TLSConfig = &tls.Config{
				RootCAs:      caCert,
				Certificates: []tls.Certificate{cert},
				MinVersion:   tls.VersionTLS12,
			}
		}

		redisCache := cache.NewRedisCache(options, c.CacheTTL)
		return Cache{Cache: redisCache}, nil
	}

	if c.CacheTTL != 0 {
		log.Logger.Warn("'--cache-ttl' is only available with Redis cache backend")
	}

	// standalone mode
	fsCache, err := cache.NewFSCache(utils.CacheDir())
	if err != nil {
		return Cache{}, xerrors.Errorf("unable to initialize fs cache: %w", err)
	}
	return Cache{Cache: fsCache}, nil
}

// Reset resets the cache
func (c Cache) Reset() (err error) {
	if err := c.ClearDB(); err != nil {
		return xerrors.Errorf("failed to clear the database: %w", err)
	}
	if err := c.ClearArtifacts(); err != nil {
		return xerrors.Errorf("failed to clear the artifact cache: %w", err)
	}
	return nil
}

// ClearDB clears the DB cache
func (c Cache) ClearDB() (err error) {
	log.Logger.Info("Removing DB file...")
	if err = os.RemoveAll(utils.CacheDir()); err != nil {
		return xerrors.Errorf("failed to remove the directory (%s) : %w", utils.CacheDir(), err)
	}
	return nil
}

// ClearArtifacts clears the artifact cache
func (c Cache) ClearArtifacts() error {
	log.Logger.Info("Removing artifact caches...")
	if err := c.Clear(); err != nil {
		return xerrors.Errorf("failed to remove the cache: %w", err)
	}
	return nil
}

// DownloadDB downloads the DB
func DownloadDB(appVersion, cacheDir, dbRepository string, quiet, insecure, skipUpdate bool) error {
	client := db.NewClient(cacheDir, quiet, insecure, db.WithDBRepository(dbRepository))
	ctx := context.Background()
	needsUpdate, err := client.NeedsUpdate(appVersion, skipUpdate)
	if err != nil {
		return xerrors.Errorf("database error: %w", err)
	}

	if needsUpdate {
		log.Logger.Info("Need to update DB")
		log.Logger.Infof("DB Repository: %s", dbRepository)
		log.Logger.Info("Downloading DB...")
		if err = client.Download(ctx, cacheDir); err != nil {
			return xerrors.Errorf("failed to download vulnerability DB: %w", err)
		}
	}

	// for debug
	if err = showDBInfo(cacheDir); err != nil {
		return xerrors.Errorf("failed to show database info: %w", err)
	}
	return nil
}

func showDBInfo(cacheDir string) error {
	m := metadata.NewClient(cacheDir)
	meta, err := m.Get()
	if err != nil {
		return xerrors.Errorf("something wrong with DB: %w", err)
	}
	log.Logger.Debugf("DB Schema: %d, UpdatedAt: %s, NextUpdate: %s, DownloadedAt: %s",
		meta.Version, meta.UpdatedAt, meta.NextUpdate, meta.DownloadedAt)
	return nil
}

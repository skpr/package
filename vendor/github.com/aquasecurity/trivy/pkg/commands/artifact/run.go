package artifact

import (
	"context"
	"errors"
	"os"

	"github.com/hashicorp/go-multierror"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
	"golang.org/x/xerrors"

	"github.com/aquasecurity/go-version/pkg/semver"
	"github.com/aquasecurity/trivy-db/pkg/db"
	tcache "github.com/aquasecurity/trivy/pkg/cache"
	"github.com/aquasecurity/trivy/pkg/commands/operation"
	"github.com/aquasecurity/trivy/pkg/fanal/analyzer"
	"github.com/aquasecurity/trivy/pkg/fanal/analyzer/config"
	"github.com/aquasecurity/trivy/pkg/fanal/analyzer/secret"
	"github.com/aquasecurity/trivy/pkg/fanal/artifact"
	"github.com/aquasecurity/trivy/pkg/fanal/cache"
	"github.com/aquasecurity/trivy/pkg/flag"
	"github.com/aquasecurity/trivy/pkg/log"
	"github.com/aquasecurity/trivy/pkg/module"
	pkgReport "github.com/aquasecurity/trivy/pkg/report"
	"github.com/aquasecurity/trivy/pkg/result"
	"github.com/aquasecurity/trivy/pkg/rpc/client"
	"github.com/aquasecurity/trivy/pkg/scanner"
	"github.com/aquasecurity/trivy/pkg/types"
	"github.com/aquasecurity/trivy/pkg/utils"
)

// TargetKind represents what kind of artifact Trivy scans
type TargetKind string

const (
	TargetContainerImage TargetKind = "image"
	TargetFilesystem     TargetKind = "fs"
	TargetRootfs         TargetKind = "rootfs"
	TargetRepository     TargetKind = "repo"
	TargetImageArchive   TargetKind = "archive"
	TargetSBOM           TargetKind = "sbom"

	devVersion = "dev"
)

var (
	defaultPolicyNamespaces = []string{"appshield", "defsec", "builtin"}
	SkipScan                = errors.New("skip subsequent processes")
)

// InitializeScanner defines the initialize function signature of scanner
type InitializeScanner func(context.Context, ScannerConfig) (scanner.Scanner, func(), error)

type ScannerConfig struct {
	// e.g. image name and file path
	Target string

	// Cache
	ArtifactCache      cache.ArtifactCache
	LocalArtifactCache cache.LocalArtifactCache

	// Client/Server options
	RemoteOption client.ScannerOption

	// Artifact options
	ArtifactOption artifact.Option
}

type Runner interface {
	// ScanImage scans an image
	ScanImage(ctx context.Context, opts flag.Options) (types.Report, error)
	// ScanFilesystem scans a filesystem
	ScanFilesystem(ctx context.Context, opts flag.Options) (types.Report, error)
	// ScanRootfs scans rootfs
	ScanRootfs(ctx context.Context, opts flag.Options) (types.Report, error)
	// ScanRepository scans repository
	ScanRepository(ctx context.Context, opts flag.Options) (types.Report, error)
	// ScanSBOM scans SBOM
	ScanSBOM(ctx context.Context, opts flag.Options) (types.Report, error)
	// Filter filter a report
	Filter(ctx context.Context, opts flag.Options, report types.Report) (types.Report, error)
	// Report a writes a report
	Report(opts flag.Options, report types.Report) error
	// Close closes runner
	Close(ctx context.Context) error
}

type runner struct {
	cache  cache.Cache
	dbOpen bool

	// WASM modules
	module *module.Manager
}

type runnerOption func(*runner)

// WithCacheClient takes a custom cache implementation
// It is useful when Trivy is imported as a library.
func WithCacheClient(c cache.Cache) runnerOption {
	return func(r *runner) {
		r.cache = c
	}
}

// NewRunner initializes Runner that provides scanning functionalities.
// It is possible to return SkipScan and it must be handled by caller.
func NewRunner(ctx context.Context, cliOptions flag.Options, opts ...runnerOption) (Runner, error) {
	r := &runner{}
	for _, opt := range opts {
		opt(r)
	}

	if err := r.initCache(cliOptions); err != nil {
		return nil, xerrors.Errorf("cache error: %w", err)
	}

	// Update the vulnerability database if needed.
	if err := r.initDB(cliOptions); err != nil {
		return nil, xerrors.Errorf("DB error: %w", err)
	}

	// Initialize WASM modules
	m, err := module.NewManager(ctx)
	if err != nil {
		return nil, xerrors.Errorf("WASM module error: %w", err)
	}
	m.Register()
	r.module = m

	return r, nil
}

// Close closes everything
func (r *runner) Close(ctx context.Context) error {
	var errs error
	if err := r.cache.Close(); err != nil {
		errs = multierror.Append(errs, err)
	}

	if r.dbOpen {
		if err := db.Close(); err != nil {
			errs = multierror.Append(errs, err)
		}
	}

	if err := r.module.Close(ctx); err != nil {
		errs = multierror.Append(errs, err)
	}
	return errs
}

func (r *runner) ScanImage(ctx context.Context, opts flag.Options) (types.Report, error) {
	// Disable the lock file scanning
	opts.DisabledAnalyzers = analyzer.TypeLockfiles

	var s InitializeScanner
	switch {
	case opts.Input != "" && opts.ServerAddr == "":
		// Scan image tarball in standalone mode
		s = archiveStandaloneScanner
	case opts.Input != "" && opts.ServerAddr != "":
		// Scan image tarball in client/server mode
		s = archiveRemoteScanner
	case opts.Input == "" && opts.ServerAddr == "":
		// Scan container image in standalone mode
		s = imageStandaloneScanner
	case opts.Input == "" && opts.ServerAddr != "":
		// Scan container image in client/server mode
		s = imageRemoteScanner
	}

	return r.scanArtifact(ctx, opts, s)
}

func (r *runner) ScanFilesystem(ctx context.Context, opts flag.Options) (types.Report, error) {
	// Disable the individual package scanning
	opts.DisabledAnalyzers = append(opts.DisabledAnalyzers, analyzer.TypeIndividualPkgs...)

	return r.scanFS(ctx, opts)
}

func (r *runner) ScanRootfs(ctx context.Context, opts flag.Options) (types.Report, error) {
	// Disable the lock file scanning
	opts.DisabledAnalyzers = append(opts.DisabledAnalyzers, analyzer.TypeLockfiles...)

	return r.scanFS(ctx, opts)
}

func (r *runner) scanFS(ctx context.Context, opts flag.Options) (types.Report, error) {
	var s InitializeScanner
	if opts.ServerAddr == "" {
		// Scan filesystem in standalone mode
		s = filesystemStandaloneScanner
	} else {
		// Scan filesystem in client/server mode
		s = filesystemRemoteScanner
	}

	return r.scanArtifact(ctx, opts, s)
}

func (r *runner) ScanRepository(ctx context.Context, opts flag.Options) (types.Report, error) {
	// Do not scan OS packages
	opts.VulnType = []string{types.VulnTypeLibrary}

	// Disable the OS analyzers and individual package analyzers
	opts.DisabledAnalyzers = append(analyzer.TypeIndividualPkgs, analyzer.TypeOSes...)

	return r.scanArtifact(ctx, opts, repositoryStandaloneScanner)
}

func (r *runner) ScanSBOM(ctx context.Context, opts flag.Options) (types.Report, error) {
	var s InitializeScanner
	if opts.ServerAddr == "" {
		// Scan cycloneDX in standalone mode
		s = sbomStandaloneScanner
	} else {
		// Scan cycloneDX in client/server mode
		s = sbomRemoteScanner
	}

	return r.scanArtifact(ctx, opts, s)
}

func (r *runner) scanArtifact(ctx context.Context, opts flag.Options, initializeScanner InitializeScanner) (types.Report, error) {
	report, err := scan(ctx, opts, initializeScanner, r.cache)
	if err != nil {
		return types.Report{}, xerrors.Errorf("scan error: %w", err)
	}

	return report, nil
}

func (r *runner) Filter(ctx context.Context, opts flag.Options, report types.Report) (types.Report, error) {
	results := report.Results

	// Filter results
	for i := range results {
		err := result.Filter(ctx, &results[i], opts.Severities, opts.IgnoreUnfixed, opts.IncludeNonFailures,
			opts.IgnoreFile, opts.IgnorePolicy, opts.IgnoredLicenses)
		if err != nil {
			return types.Report{}, xerrors.Errorf("unable to filter vulnerabilities: %w", err)
		}
	}
	return report, nil
}

func (r *runner) Report(opts flag.Options, report types.Report) error {
	if err := pkgReport.Write(report, pkgReport.Option{
		AppVersion:         opts.AppVersion,
		Format:             opts.Format,
		Output:             opts.Output,
		Tree:               opts.DependencyTree,
		Severities:         opts.Severities,
		OutputTemplate:     opts.Template,
		IncludeNonFailures: opts.IncludeNonFailures,
		Trace:              opts.Trace,
	}); err != nil {
		return xerrors.Errorf("unable to write results: %w", err)
	}

	return nil
}

func (r *runner) initDB(opts flag.Options) error {
	// When scanning config files or running as client mode, it doesn't need to download the vulnerability database.
	if opts.ServerAddr != "" || !slices.Contains(opts.SecurityChecks, types.SecurityCheckVulnerability) {
		return nil
	}

	// download the database file
	noProgress := opts.Quiet || opts.NoProgress
	if err := operation.DownloadDB(opts.AppVersion, opts.CacheDir, opts.DBRepository, noProgress, opts.Insecure, opts.SkipDBUpdate); err != nil {
		return err
	}

	if opts.DownloadDBOnly {
		return SkipScan
	}

	if err := db.Init(opts.CacheDir); err != nil {
		return xerrors.Errorf("error in vulnerability DB initialize: %w", err)
	}
	r.dbOpen = true

	return nil
}

func (r *runner) initCache(opts flag.Options) error {
	// Skip initializing cache when custom cache is passed
	if r.cache != nil {
		return nil
	}

	// client/server mode
	if opts.ServerAddr != "" {
		remoteCache := tcache.NewRemoteCache(opts.ServerAddr, opts.CustomHeaders, opts.Insecure)
		r.cache = tcache.NopCache(remoteCache)
		return nil
	}

	// standalone mode
	utils.SetCacheDir(opts.CacheDir)
	cacheClient, err := operation.NewCache(opts.CacheOptions)
	if err != nil {
		return xerrors.Errorf("unable to initialize the cache: %w", err)
	}
	log.Logger.Debugf("cache dir:  %s", utils.CacheDir())

	if opts.Reset {
		defer cacheClient.Close()
		if err = cacheClient.Reset(); err != nil {
			return xerrors.Errorf("cache reset error: %w", err)
		}
		return SkipScan
	}
	if opts.ClearCache {
		defer cacheClient.Close()
		if err = cacheClient.ClearArtifacts(); err != nil {
			return xerrors.Errorf("cache clear error: %w", err)
		}
		return SkipScan
	}

	r.cache = cacheClient
	return nil
}

// Run performs artifact scanning
//func Run(cliCtx *cli.Context, targetKind TargetKind) error {
//	opt, err := InitOption(cliCtx)
//	if err != nil {
//		return xerrors.Errorf("InitOption: %w", err)
//	}
//
//	return run(cliCtx.Context, opt, targetKind)
//}

// Run performs artifact scanning
func Run(ctx context.Context, opts flag.Options, targetKind TargetKind) (err error) {
	ctx, cancel := context.WithTimeout(ctx, opts.Timeout)
	defer cancel()

	defer func() {
		if xerrors.Is(err, context.DeadlineExceeded) {
			log.Logger.Warn("Increase --timeout value")
		}
	}()

	if opts.GenerateDefaultConfig {
		log.Logger.Info("Writing the default config to trivy-default.yaml...")
		return viper.SafeWriteConfigAs("trivy-default.yaml")
	}

	r, err := NewRunner(ctx, opts)
	if err != nil {
		if errors.Is(err, SkipScan) {
			return nil
		}
		return xerrors.Errorf("init error: %w", err)
	}
	defer r.Close(ctx)

	var report types.Report
	switch targetKind {
	case TargetContainerImage, TargetImageArchive:
		if report, err = r.ScanImage(ctx, opts); err != nil {
			return xerrors.Errorf("image scan error: %w", err)
		}
	case TargetFilesystem:
		if report, err = r.ScanFilesystem(ctx, opts); err != nil {
			return xerrors.Errorf("filesystem scan error: %w", err)
		}
	case TargetRootfs:
		if report, err = r.ScanRootfs(ctx, opts); err != nil {
			return xerrors.Errorf("rootfs scan error: %w", err)
		}
	case TargetRepository:
		if report, err = r.ScanRepository(ctx, opts); err != nil {
			return xerrors.Errorf("repository scan error: %w", err)
		}
	case TargetSBOM:
		if report, err = r.ScanSBOM(ctx, opts); err != nil {
			return xerrors.Errorf("sbom scan error: %w", err)
		}
	}

	report, err = r.Filter(ctx, opts, report)
	if err != nil {
		return xerrors.Errorf("filter error: %w", err)
	}

	if err = r.Report(opts, report); err != nil {
		return xerrors.Errorf("report error: %w", err)
	}

	Exit(opts, report.Results.Failed())

	return nil
}

func disabledAnalyzers(opts flag.Options) []analyzer.Type {
	// Specified analyzers to be disabled depending on scanning modes
	// e.g. The 'image' subcommand should disable the lock file scanning.
	analyzers := opts.DisabledAnalyzers

	// It doesn't analyze apk commands by default.
	if !opts.ScanRemovedPkgs {
		analyzers = append(analyzers, analyzer.TypeApkCommand)
	}

	// Do not analyze programming language packages when not running in 'library'
	if !slices.Contains(opts.VulnType, types.VulnTypeLibrary) {
		analyzers = append(analyzers, analyzer.TypeLanguages...)
	}

	// Do not perform secret scanning when it is not specified.
	if !slices.Contains(opts.SecurityChecks, types.SecurityCheckSecret) {
		analyzers = append(analyzers, analyzer.TypeSecret)
	}

	// Do not perform misconfiguration scanning when it is not specified.
	if !slices.Contains(opts.SecurityChecks, types.SecurityCheckConfig) &&
		!slices.Contains(opts.SecurityChecks, types.SecurityCheckRbac) {
		analyzers = append(analyzers, analyzer.TypeConfigFiles...)
	}

	// Scanning file headers and license files is expensive.
	// It is performed only when '--security-checks license' and '--license-full' are specified.
	if !slices.Contains(opts.SecurityChecks, types.SecurityCheckLicense) || !opts.LicenseFull {
		analyzers = append(analyzers, analyzer.TypeLicenseFile)
	}

	return analyzers
}

func initScannerConfig(opts flag.Options, cacheClient cache.Cache) (ScannerConfig, types.ScanOptions, error) {
	target := opts.Target
	if opts.Input != "" {
		target = opts.Input
	}

	scanOptions := types.ScanOptions{
		VulnType:            opts.VulnType,
		SecurityChecks:      opts.SecurityChecks,
		ScanRemovedPackages: opts.ScanRemovedPkgs, // this is valid only for 'image' subcommand
		ListAllPackages:     opts.ListAllPkgs,
		LicenseCategories:   opts.LicenseCategories,
	}

	if slices.Contains(opts.SecurityChecks, types.SecurityCheckVulnerability) {
		log.Logger.Info("Vulnerability scanning is enabled")
		log.Logger.Debugf("Vulnerability type:  %s", scanOptions.VulnType)
	}

	// ScannerOption is filled only when config scanning is enabled.
	var configScannerOptions config.ScannerOption
	if slices.Contains(opts.SecurityChecks, types.SecurityCheckConfig) {
		log.Logger.Info("Misconfiguration scanning is enabled")
		configScannerOptions = config.ScannerOption{
			Trace:            opts.Trace,
			Namespaces:       append(opts.PolicyNamespaces, defaultPolicyNamespaces...),
			PolicyPaths:      opts.PolicyPaths,
			DataPaths:        opts.DataPaths,
			FilePatterns:     opts.FilePatterns,
			HelmValues:       opts.HelmValues,
			HelmValueFiles:   opts.HelmValueFiles,
			HelmFileValues:   opts.HelmFileValues,
			HelmStringValues: opts.HelmStringValues,
			TerraformTFVars:  opts.TerraformTFVars,
		}
	}

	// Do not load config file for secret scanning
	if slices.Contains(opts.SecurityChecks, types.SecurityCheckSecret) {
		ver := canonicalVersion(opts.AppVersion)
		log.Logger.Info("Secret scanning is enabled")
		log.Logger.Info("If your scanning is slow, please try '--security-checks vuln' to disable secret scanning")
		log.Logger.Infof("Please see also https://aquasecurity.github.io/trivy/%s/docs/secret/scanning/#recommendation for faster secret detection", ver)
	} else {
		opts.SecretConfigPath = ""
	}

	if slices.Contains(opts.SecurityChecks, types.SecurityCheckLicense) {
		if opts.LicenseFull {
			log.Logger.Info("Full license scanning is enabled")
		} else {
			log.Logger.Info("License scanning is enabled")
		}
	}

	return ScannerConfig{
		Target:             target,
		ArtifactCache:      cacheClient,
		LocalArtifactCache: cacheClient,
		RemoteOption: client.ScannerOption{
			RemoteURL:     opts.ServerAddr,
			CustomHeaders: opts.CustomHeaders,
			Insecure:      opts.Insecure,
		},
		ArtifactOption: artifact.Option{
			DisabledAnalyzers: disabledAnalyzers(opts),
			SkipFiles:         opts.SkipFiles,
			SkipDirs:          opts.SkipDirs,
			InsecureSkipTLS:   opts.Insecure,
			Offline:           opts.OfflineScan,
			NoProgress:        opts.NoProgress || opts.Quiet,
			RepoBranch:        opts.RepoBranch,
			RepoCommit:        opts.RepoCommit,
			RepoTag:           opts.RepoTag,

			// For misconfiguration scanning
			MisconfScannerOption: configScannerOptions,

			// For secret scanning
			SecretScannerOption: secret.ScannerOption{
				ConfigPath: opts.SecretConfigPath,
			},
		},
	}, scanOptions, nil
}

func scan(ctx context.Context, opts flag.Options, initializeScanner InitializeScanner, cacheClient cache.Cache) (
	types.Report, error) {

	scannerConfig, scanOptions, err := initScannerConfig(opts, cacheClient)
	if err != nil {
		return types.Report{}, err
	}

	s, cleanup, err := initializeScanner(ctx, scannerConfig)
	if err != nil {
		return types.Report{}, xerrors.Errorf("unable to initialize a scanner: %w", err)
	}
	defer cleanup()

	report, err := s.ScanArtifact(ctx, scanOptions)
	if err != nil {
		return types.Report{}, xerrors.Errorf("scan failed: %w", err)
	}
	return report, nil
}

func Exit(opts flag.Options, failedResults bool) {
	if opts.ExitCode != 0 && failedResults {
		os.Exit(opts.ExitCode)
	}
}

func canonicalVersion(ver string) string {
	if ver == devVersion {
		return ver
	}
	v, err := semver.Parse(ver)
	if err != nil {
		return devVersion
	}
	// Replace pre-release with "dev"
	// e.g. v0.34.0-beta1+snapshot-1
	if v.IsPreRelease() || v.Metadata() != "" {
		return devVersion
	}

	// Add "v" prefix, "0.34.0" => "v0.34.0" for the url
	return "v" + ver
}

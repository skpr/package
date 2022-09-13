package artifact

import (
	"context"

	"golang.org/x/xerrors"

	"github.com/aquasecurity/trivy/pkg/scanner"
	"github.com/aquasecurity/trivy/pkg/types"
)

// imageStandaloneScanner initializes a container image scanner in standalone mode
// $ trivy image alpine:3.15
func imageStandaloneScanner(ctx context.Context, conf ScannerConfig) (scanner.Scanner, func(), error) {
	dockerOpt, err := types.GetDockerOption(conf.ArtifactOption.InsecureSkipTLS)
	if err != nil {
		return scanner.Scanner{}, nil, err
	}
	s, cleanup, err := initializeDockerScanner(ctx, conf.Target, conf.ArtifactCache, conf.LocalArtifactCache,
		dockerOpt, conf.ArtifactOption)
	if err != nil {
		return scanner.Scanner{}, func() {}, xerrors.Errorf("unable to initialize a docker scanner: %w", err)
	}
	return s, cleanup, nil
}

// archiveStandaloneScanner initializes an image archive scanner in standalone mode
// $ trivy image --input alpine.tar
func archiveStandaloneScanner(ctx context.Context, conf ScannerConfig) (scanner.Scanner, func(), error) {
	s, err := initializeArchiveScanner(ctx, conf.Target, conf.ArtifactCache, conf.LocalArtifactCache, conf.ArtifactOption)
	if err != nil {
		return scanner.Scanner{}, func() {}, xerrors.Errorf("unable to initialize the archive scanner: %w", err)
	}
	return s, func() {}, nil
}

// imageRemoteScanner initializes a container image scanner in client/server mode
// $ trivy image --server localhost:4954 alpine:3.15
func imageRemoteScanner(ctx context.Context, conf ScannerConfig) (
	scanner.Scanner, func(), error) {
	// Scan an image in Docker Engine, Docker Registry, etc.
	dockerOpt, err := types.GetDockerOption(conf.ArtifactOption.InsecureSkipTLS)
	if err != nil {
		return scanner.Scanner{}, nil, err
	}

	s, cleanup, err := initializeRemoteDockerScanner(ctx, conf.Target, conf.ArtifactCache, conf.RemoteOption,
		dockerOpt, conf.ArtifactOption)
	if err != nil {
		return scanner.Scanner{}, nil, xerrors.Errorf("unable to initialize the docker scanner: %w", err)
	}
	return s, cleanup, nil
}

// archiveRemoteScanner initializes an image archive scanner in client/server mode
// $ trivy image --server localhost:4954 --input alpine.tar
func archiveRemoteScanner(ctx context.Context, conf ScannerConfig) (scanner.Scanner, func(), error) {
	// Scan tar file
	s, err := initializeRemoteArchiveScanner(ctx, conf.Target, conf.ArtifactCache, conf.RemoteOption, conf.ArtifactOption)
	if err != nil {
		return scanner.Scanner{}, nil, xerrors.Errorf("unable to initialize the archive scanner: %w", err)
	}
	return s, func() {}, nil
}

// filesystemStandaloneScanner initializes a filesystem scanner in standalone mode
func filesystemStandaloneScanner(ctx context.Context, conf ScannerConfig) (scanner.Scanner, func(), error) {
	s, cleanup, err := initializeFilesystemScanner(ctx, conf.Target, conf.ArtifactCache, conf.LocalArtifactCache, conf.ArtifactOption)
	if err != nil {
		return scanner.Scanner{}, func() {}, xerrors.Errorf("unable to initialize a filesystem scanner: %w", err)
	}
	return s, cleanup, nil
}

// filesystemRemoteScanner initializes a filesystem scanner in client/server mode
func filesystemRemoteScanner(ctx context.Context, conf ScannerConfig) (scanner.Scanner, func(), error) {
	s, cleanup, err := initializeRemoteFilesystemScanner(ctx, conf.Target, conf.ArtifactCache, conf.RemoteOption, conf.ArtifactOption)
	if err != nil {
		return scanner.Scanner{}, func() {}, xerrors.Errorf("unable to initialize a filesystem scanner: %w", err)
	}
	return s, cleanup, nil
}

// filesystemStandaloneScanner initializes a repository scanner in standalone mode
func repositoryStandaloneScanner(ctx context.Context, conf ScannerConfig) (scanner.Scanner, func(), error) {
	s, cleanup, err := initializeRepositoryScanner(ctx, conf.Target, conf.ArtifactCache, conf.LocalArtifactCache, conf.ArtifactOption)
	if err != nil {
		return scanner.Scanner{}, func() {}, xerrors.Errorf("unable to initialize a filesystem scanner: %w", err)
	}
	return s, cleanup, nil
}

// sbomStandaloneScanner initializes a SBOM scanner in standalone mode
func sbomStandaloneScanner(ctx context.Context, conf ScannerConfig) (scanner.Scanner, func(), error) {
	s, cleanup, err := initializeSBOMScanner(ctx, conf.Target, conf.ArtifactCache, conf.LocalArtifactCache, conf.ArtifactOption)
	if err != nil {
		return scanner.Scanner{}, func() {}, xerrors.Errorf("unable to initialize a cycloneDX scanner: %w", err)
	}
	return s, cleanup, nil
}

// sbomRemoteScanner initializes a SBOM scanner in client/server mode
func sbomRemoteScanner(ctx context.Context, conf ScannerConfig) (scanner.Scanner, func(), error) {
	s, cleanup, err := initializeRemoteSBOMScanner(ctx, conf.Target, conf.ArtifactCache, conf.RemoteOption, conf.ArtifactOption)
	if err != nil {
		return scanner.Scanner{}, func() {}, xerrors.Errorf("unable to initialize a cycloneDX scanner: %w", err)
	}
	return s, cleanup, nil
}

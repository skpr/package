package report

import (
	"io"
	"strings"
	"sync"

	"github.com/aquasecurity/trivy/pkg/report/predicate"
	"github.com/aquasecurity/trivy/pkg/report/table"

	"golang.org/x/xerrors"

	dbTypes "github.com/aquasecurity/trivy-db/pkg/types"
	"github.com/aquasecurity/trivy/pkg/log"
	"github.com/aquasecurity/trivy/pkg/report/cyclonedx"
	"github.com/aquasecurity/trivy/pkg/report/github"
	"github.com/aquasecurity/trivy/pkg/report/spdx"
	"github.com/aquasecurity/trivy/pkg/types"
)

const (
	SchemaVersion = 2

	FormatTable      = "table"
	FormatJSON       = "json"
	FormatTemplate   = "template"
	FormatSarif      = "sarif"
	FormatCycloneDX  = "cyclonedx"
	FormatSPDX       = "spdx"
	FormatSPDXJSON   = "spdx-json"
	FormatGitHub     = "github"
	FormatCosignVuln = "cosign-vuln"
)

var (
	SupportedSBOMFormats = []string{FormatCycloneDX, FormatSPDX, FormatSPDXJSON, FormatGitHub}
)

type Option struct {
	AppVersion string

	Format         string
	Output         io.Writer
	Tree           bool
	Severities     []dbTypes.Severity
	OutputTemplate string

	// For misconfigurations
	IncludeNonFailures bool
	Trace              bool

	// For licenses
	LicenseRiskThreshold int
	IgnoredLicenses      []string
}

// Write writes the result to output, format as passed in argument
func Write(report types.Report, option Option) error {
	var writer Writer
	switch option.Format {
	case FormatTable:
		writer = &table.Writer{
			Output:               option.Output,
			Severities:           option.Severities,
			Tree:                 option.Tree,
			ShowMessageOnce:      &sync.Once{},
			IncludeNonFailures:   option.IncludeNonFailures,
			Trace:                option.Trace,
			LicenseRiskThreshold: option.LicenseRiskThreshold,
			IgnoredLicenses:      option.IgnoredLicenses,
		}
	case FormatJSON:
		writer = &JSONWriter{Output: option.Output}
	case FormatGitHub:
		writer = &github.Writer{Output: option.Output, Version: option.AppVersion}
	case FormatCycloneDX:
		// TODO: support xml format option with cyclonedx writer
		writer = cyclonedx.NewWriter(option.Output, option.AppVersion)
	case FormatSPDX, FormatSPDXJSON:
		writer = spdx.NewWriter(option.Output, option.AppVersion, option.Format)
	case FormatTemplate:
		// We keep `sarif.tpl` template working for backward compatibility for a while.
		if strings.HasPrefix(option.OutputTemplate, "@") && strings.HasSuffix(option.OutputTemplate, "sarif.tpl") {
			log.Logger.Warn("Using `--template sarif.tpl` is deprecated. Please migrate to `--format sarif`. See https://github.com/aquasecurity/trivy/discussions/1571")
			writer = SarifWriter{Output: option.Output, Version: option.AppVersion}
			break
		}
		var err error
		if writer, err = NewTemplateWriter(option.Output, option.OutputTemplate); err != nil {
			return xerrors.Errorf("failed to initialize template writer: %w", err)
		}
	case FormatSarif:
		writer = SarifWriter{Output: option.Output, Version: option.AppVersion}
	case FormatCosignVuln:
		writer = predicate.NewVulnWriter(option.Output, option.AppVersion)
	default:
		return xerrors.Errorf("unknown format: %v", option.Format)
	}

	if err := writer.Write(report); err != nil {
		return xerrors.Errorf("failed to write results: %w", err)
	}
	return nil
}

// Writer defines the result write operation
type Writer interface {
	Write(types.Report) error
}

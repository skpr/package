package types

const (
	ArtifactJSONSchemaVersion = 1
	BlobJSONSchemaVersion     = 2
)

const (
	// Programming language dependencies
	Bundler    = "bundler"
	GemSpec    = "gemspec"
	Cargo      = "cargo"
	Composer   = "composer"
	Npm        = "npm"
	NuGet      = "nuget"
	DotNetCore = "dotnet-core"
	Pip        = "pip"
	Pipenv     = "pipenv"
	Poetry     = "poetry"
	PythonPkg  = "python-pkg"
	NodePkg    = "node-pkg"
	Yarn       = "yarn"
	Pnpm       = "pnpm"
	Jar        = "jar"
	Pom        = "pom"
	GoBinary   = "gobinary"
	GoModule   = "gomod"
	JavaScript = "javascript"
	RustBinary = "rustbinary"

	// Config files
	YAML           = "yaml"
	JSON           = "json"
	Dockerfile     = "dockerfile"
	Terraform      = "terraform"
	CloudFormation = "cloudformation"
	Kubernetes     = "kubernetes"
	Ansible        = "ansible"
	Helm           = "helm"
	Rbac           = "rbac"
	Cloud          = "cloud"

	// Licensing
	License = "license"

	// Language-specific file names
	NuGetPkgsLock   = "packages.lock.json"
	NuGetPkgsConfig = "packages.config"

	GoMod = "go.mod"
	GoSum = "go.sum"

	MavenPom = "pom.xml"

	NpmPkgLock = "package-lock.json"
	YarnLock   = "yarn.lock"
	PnpmLock   = "pnpm-lock.yaml"

	ComposerLock = "composer.lock"

	PipRequirements = "requirements.txt"
	PipfileLock     = "Pipfile.lock"
	PoetryLock      = "poetry.lock"

	GemfileLock = "Gemfile.lock"

	CargoLock = "Cargo.lock"
)

package main

import (
	"os"

	"github.com/alecthomas/kingpin"
	docker "github.com/fsouza/go-dockerclient"

	"github.com/skpr/package/pkg/builder"
)

var (
	cliDockerUser = kingpin.Flag("docker-username", "Username for Docker authentication").Envar("DOCKER_USERNAME").String()
	cliDockerPass = kingpin.Flag("docker-password", "Password for Docker authentication").Envar("DOCKER_PASSWORD").String()
	cliRegistry   = kingpin.Flag("verbose", "Verbose mode.").String()
	cliContext    = kingpin.Flag("context", "Path to use as a context for building images.").Default(".").String()
	cliNoPush     = kingpin.Flag("no-push", "Don't push images to the registry after being built. Used for local debugging.").Bool()
	cliDirectory  = kingpin.Flag("directory", "The location of the package directory").Default(".skpr/package").String()
	cliDebug      = kingpin.Flag("debug", "Show debug information").Bool()
	cliVersion    = kingpin.Arg("version", "Version of the application which is being packaged").Required().String()
)

func main() {
	kingpin.Parse()

	params := builder.Params{
		Directory: *cliDirectory,
		Debug:     *cliDebug,
		Writer:    os.Stdout,
		Registry:  *cliRegistry,
		Version:   *cliVersion,
		Context:   *cliContext,
		NoPush:    *cliNoPush,
		Auth: docker.AuthConfiguration{
			Username: *cliDockerUser,
			Password: *cliDockerPass,
		},
	}

	_, err := builder.BuildAndPush(params)
	if err != nil {
		panic(err)
	}
}

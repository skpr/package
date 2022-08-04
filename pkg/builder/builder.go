package builder

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/segmentio/textio"
	"golang.org/x/sync/errgroup"

	"github.com/skpr/package/pkg/color"
	"github.com/skpr/package/pkg/utils/aws/ecr"
	"github.com/skpr/package/pkg/utils/finder"
	"github.com/skpr/package/pkg/utils/image"
)

// DockerClientInterface provides an interface that allows us to test the builder.
type DockerClientInterface interface {
	BuildImage(options docker.BuildImageOptions) error
	PushImage(options docker.PushImageOptions, auth docker.AuthConfiguration) error
}

// Builder is the docker image builder.
type Builder struct {
	dockerClient DockerClientInterface
}

// Params used for building the applications.
type Params struct {
	Directory string
	Debug     bool
	Writer    io.Writer
	Registry  string
	Version   string
	Context   string
	NoPush    bool
	Auth      docker.AuthConfiguration
}

const (
	// ImageNameCompile is used for compiling the application.
	ImageNameCompile = "compile"

	// BuildArgCompileImage is used for referencing the compile image.
	BuildArgCompileImage = "COMPILE_IMAGE"
)

// NewBuilder creates a new Builder.
func NewBuilder(dockerClient DockerClientInterface) *Builder {
	return &Builder{
		dockerClient: dockerClient,
	}
}

// BuildOutput provided to tasks which trigger a build.
type BuildOutput struct {
	Images map[string]string `json:"image" yaml:"image"`
}

// BuildAndPush a packaged set of images.
func BuildAndPush(params Params) (BuildOutput, error) {
	var output BuildOutput

	// @todo, Consider abstracting this if another registry + credentials pair is required.
	if ecr.IsRegistry(params.Registry) {
		auth, err := ecr.UpgradeAuth(params.Registry, params.Auth)
		if err != nil {
			return output, fmt.Errorf("failed to upgrade AWS ECR authentication: %w", err)
		}

		params.Auth = auth
	}

	dockerfiles, err := finder.FindDockerfiles(params.Directory)
	if err != nil {
		return output, fmt.Errorf("failed to find dockerfiles: %w", err)
	}

	if params.Debug {
		fmt.Println("Found the following dockerfiles:")
		for key, path := range dockerfiles {
			fmt.Printf("%-10s %q\n", key, path)
		}
	}

	// Print deprecation notice.
	for key, path := range dockerfiles {
		if strings.HasSuffix(path, ".dockerfile") {
			fmt.Printf("[DEPRECATED] Dockerfile location %q is deprecated. Use \"%s/%s/Dockerfile\" instead.\n", path, filepath.Dir(path), key)
		}
	}

	dockerclient, err := docker.NewClientFromEnv()
	if err != nil {
		return output, fmt.Errorf("failed to setup Docker client: %w", err)
	}

	builder := NewBuilder(dockerclient)

	output, err = builder.Build(dockerfiles, params)
	if err != nil {
		return output, err
	}

	return output, nil
}

// Build the images.
func (b *Builder) Build(dockerfiles finder.Dockerfiles, params Params) (BuildOutput, error) {
	resp := BuildOutput{
		Images: make(map[string]string),
	}

	compileDockerfile, ok := dockerfiles[ImageNameCompile]
	if !ok {
		return resp, fmt.Errorf("%q is a required dockerfile", ImageNameCompile)
	}

	// We build the compile image first, as it is the base image for other images.
	compileBuild := docker.BuildImageOptions{
		Name:         image.Name(params.Registry, params.Version, ImageNameCompile),
		Dockerfile:   compileDockerfile,
		ContextDir:   params.Context,
		OutputStream: prefix(params.Writer, ImageNameCompile),
	}

	// We need to build the 'compile' image first.
	fmt.Fprintf(params.Writer, "Building image: %s\n", compileBuild.Name)
	start := time.Now()
	err := b.dockerClient.BuildImage(compileBuild)
	if err != nil {
		return resp, err
	}
	fmt.Fprintf(params.Writer, "Built compile image in %s\n", time.Since(start).Round(time.Second))

	// Remove compile from list of dockerfiles.
	delete(dockerfiles, ImageNameCompile)

	args := []docker.BuildArg{
		{
			Name:  BuildArgCompileImage,
			Value: image.Name(params.Registry, params.Version, ImageNameCompile),
		},
	}
	var builds []docker.BuildImageOptions
	for imageName, dockerfile := range dockerfiles {
		builds = append(builds, docker.BuildImageOptions{
			Name:         image.Name(params.Registry, params.Version, imageName),
			Dockerfile:   dockerfile,
			ContextDir:   params.Context,
			OutputStream: prefix(params.Writer, imageName),
			BuildArgs:    args,
		})
	}

	bg, ctx := errgroup.WithContext(context.Background())

	for _, build := range builds {
		// https://golang.org/doc/faq#closures_and_goroutines
		build := build

		// Allows us to cancel build executions.
		build.Context = ctx

		fmt.Fprintf(params.Writer, "Building image: %s\n", build.Name)

		bg.Go(func() error {
			start = time.Now()
			err := b.dockerClient.BuildImage(build)
			if err != nil {
				return err
			}
			fmt.Fprintf(params.Writer, "Built %s image in %s\n", build.Name, time.Since(start).Round(time.Second))
			return nil
		})
	}
	err = bg.Wait()
	if err != nil {
		return resp, err
	}

	if params.NoPush {
		return resp, nil
	}

	var pushes []docker.PushImageOptions
	for imageName := range dockerfiles {
		// Compile image is only for building, so we don't push.
		if imageName == ImageNameCompile {
			continue
		}

		tag := image.Tag(params.Version, imageName)

		resp.Images[imageName] = fmt.Sprintf("%s:%s", params.Registry, tag)

		pushes = append(pushes, docker.PushImageOptions{
			Name: params.Registry,
			Tag:  tag,
		})
	}

	pg, ctx := errgroup.WithContext(context.Background())

	for _, push := range pushes {
		// https://golang.org/doc/faq#closures_and_goroutines
		push := push

		// Allows us to cancel push executions.
		push.Context = ctx

		fmt.Fprintf(params.Writer, "Pushing image: %s:%s\n", push.Name, push.Tag)

		pg.Go(func() error {
			start = time.Now()
			err = b.dockerClient.PushImage(push, params.Auth)
			if err != nil {
				return err
			}
			fmt.Fprintf(params.Writer, "Pushed %s:%s image in %s\n", push.Name, push.Tag, time.Since(start).Round(time.Second))
			return nil
		})
	}
	err = pg.Wait()
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Helper function to prefix all output for a stream.
func prefix(w io.Writer, name string) io.Writer {
	return textio.NewPrefixWriter(w, fmt.Sprintf("%s\t", color.Wrap(strings.ToUpper(name))))
}

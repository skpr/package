package builder

import (
	"bytes"
	"testing"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/stretchr/testify/assert"

	"github.com/skpr/package/pkg/builder/mock"
	"github.com/skpr/package/pkg/utils/finder"
)

func TestBuild(t *testing.T) {

	dockerClient := &mock.DockerClient{}
	dockerClient.BuildWg.Add(4)
	dockerClient.PushWg.Add(3)

	dockerFiles := make(finder.Dockerfiles)
	dockerFiles["compile"] = ".skpr/package/compile/Dockerfile"
	dockerFiles["cli"] = ".skpr/package/cli/Dockerfile"
	dockerFiles["app"] = ".skpr/package/app/Dockerfile"
	dockerFiles["web"] = ".skpr/package/web/Dockerfile"

	var b bytes.Buffer

	params := Params{
		Writer:   &b,
		Registry: "foo",
		Version:  "222",
		Context:  "bar",
		NoPush:   false,
		Auth:     docker.AuthConfiguration{},
	}

	builder := NewBuilder(dockerClient)
	_, err := builder.Build(dockerFiles, params)
	assert.NoError(t, err)

	assert.Equal(t, 4, dockerClient.BuildCount())
	assert.Equal(t, 3, dockerClient.PushCount())

}

package task

import (
	"testing"

	types "github.com/more-than-code/deploybot-service-builder/deploybot-types"
	"github.com/more-than-code/deploybot-service-builder/util"
)

func TestBuildImage(t *testing.T) {
	r := NewRunner()
	err := r.DoTask(types.Task{Type: types.TaskBuild, Config: types.BuildConfig{
		ImageName:  "binartist/geoy-graph",
		ImageTag:   "latest",
		RepoUrl:    "https://github.com/joe-and-his-friends/geoy-services.git",
		RepoName:   "geoy-services",
		Dockerfile: "graph/app.dockerfile",
	}}, nil)

	if err != nil {
		t.Error(err)
	}
}

func TestPushImage(t *testing.T) {
	h := util.NewContainerHelper("unix:///var/run/docker.sock")

	err := h.PushImage("binartist/geoy-graph")
	if err != nil {
		t.Error(err)
	}
}

func TestRunContainer(t *testing.T) {
	env := []string{}

	r := NewRunner()
	err := r.DoTask(types.Task{Type: types.TaskDeploy, Config: types.DeployConfig{ExposedPort: "9000", HostPort: "9000", Env: env, ImageName: "binartist/mo-service-graph", ImageTag: "latest", ServiceName: "graph", AutoRemove: false}}, nil)

	if err != nil {
		t.Error(err)
	}

}

package task

import (
	"encoding/json"
	"os"

	dTypes "github.com/docker/docker/api/types"
	"github.com/kelseyhightower/envconfig"
	types "github.com/more-than-code/deploybot-service-builder/deploybot-types"
	"github.com/more-than-code/deploybot-service-builder/util"
)

type RunnerConfig struct {
	ProjectsPath string `envconfig:"PROJECTS_PATH"`
	DockerHost   string `envconfig:"DOCKER_HOST"`
}

type Runner struct {
	cfg     RunnerConfig
	cHelper *util.ContainerHelper
}

func NewRunner() *Runner {
	var cfg RunnerConfig
	err := envconfig.Process("", &cfg)
	if err != nil {
		panic(err)
	}

	return &Runner{cfg: cfg, cHelper: util.NewContainerHelper(cfg.DockerHost)}
}

func (r *Runner) DoTask(t types.Task, arguments []string) error {

	var c types.BuildConfig

	bs, err := json.Marshal(t.Config)

	if err != nil {
		return err
	}

	err = json.Unmarshal(bs, &c)

	if err != nil {
		return err
	}

	if c.RepoBranch == "" {
		c.RepoBranch = "main"
	}

	// !!! Never omit the trailing slash, otherwise util.TarFiles will fail
	path := r.cfg.ProjectsPath + "/" + c.RepoName + "_" + c.RepoBranch + "/"

	os.RemoveAll(path)
	err = util.CloneRepo(path, c.RepoUrl, c.RepoBranch)

	if err != nil {
		return err
	}

	files, err := util.TarFiles(path)

	if err != nil {
		return err
	}

	imageNameTag := c.ImageName + ":" + c.ImageTag

	err = r.cHelper.BuildImage(files, &dTypes.ImageBuildOptions{Dockerfile: c.Dockerfile, Tags: []string{imageNameTag}, BuildArgs: c.Args, Version: dTypes.BuilderBuildKit})

	if err != nil {
		return err
	}

	r.cHelper.PushImage(imageNameTag)

	return nil
}

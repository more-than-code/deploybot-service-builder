package util

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/kelseyhightower/envconfig"
)

type ContainerHelperConfig struct {
	DhUsername string `envconfig:"DH_USERNAME"`
	DhPassword string `envconfig:"DH_PASSWORD"`
}

type ContainerHelper struct {
	cli *client.Client
	cfg ContainerHelperConfig
}

func NewContainerHelper(host string) *ContainerHelper {
	var cfg ContainerHelperConfig
	err := envconfig.Process("", &cfg)
	if err != nil {
		panic(err)
	}

	cli, err := client.NewClientWithOpts(client.WithHost(host), client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	return &ContainerHelper{cli: cli, cfg: cfg}
}

func (h *ContainerHelper) BuildImage(buildContext io.Reader, buidOptions *types.ImageBuildOptions) error {
	buildResponse, err := h.cli.ImageBuild(context.Background(), buildContext, *buidOptions)

	if err != nil {
		return err
	}

	res, err := io.ReadAll(buildResponse.Body)
	if err != nil {
		return err
	}

	log.Println(string(res))

	buildResponse.Body.Close()

	return nil
}

func (h *ContainerHelper) PushImage(name string) error {
	authConfig := types.AuthConfig{
		Username: h.cfg.DhUsername,
		Password: h.cfg.DhPassword,
	}
	encodedJSON, _ := json.Marshal(authConfig)
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)

	res, err := h.cli.ImagePush(context.Background(), name, types.ImagePushOptions{RegistryAuth: authStr})

	if err != nil {
		return err
	}

	defer res.Close()
	io.Copy(os.Stdout, res)
	return nil
}

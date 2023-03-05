package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"github.com/more-than-code/deploybot-service-builder/task"
)

type Config struct {
	ServerPort int `envconfig:"SERVER_PORT"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		panic(err)
	}

	g := gin.Default()

	t := task.NewScheduler()
	g.POST("/ghWebhook", t.GhWebhookHandler())
	g.POST("/streamWebhook", t.StreamWebhookHandler())
	g.GET("/healthCheck", t.HealthCheckHandler())

	g.Run(fmt.Sprintf(":%d", cfg.ServerPort))
}

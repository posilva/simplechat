// Package main is the entry point for simple chat app
package main

import (
	"embed"
	"html/template"

	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/posilva/simplechat/internal/adapters/input/handler"
	"github.com/posilva/simplechat/internal/adapters/output/moderator"
	"github.com/posilva/simplechat/internal/adapters/output/notifier"
	"github.com/posilva/simplechat/internal/adapters/output/registry"
	"github.com/posilva/simplechat/internal/adapters/output/repository"
	"github.com/posilva/simplechat/internal/core/services"
	"github.com/posilva/simplechat/internal/testutil"
)

//go:embed templates/*
var f embed.FS

func main() {

	r := gin.Default()

	templ := template.Must(template.New("").ParseFS(f, "templates/*.html"))
	r.SetHTMLTemplate(templ)

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Simple Chat",
		})
	})

	repo, err := repository.NewDynamoDBRepository(repository.DefaultiLocalAWSClientConfig(), testutil.DynamoDBLocalTableName)
	if err != nil {
		panic(fmt.Errorf("failed to create repositort %v", err))
	}
	reg := registry.NewInMemoryRegistry()
	notif, err := notifier.NewRabbitMQNotifierWithLocal(testutil.RabbitMQLocalURL, reg)
	if err != nil {
		panic(fmt.Errorf("failed to create notifier: %v", err))
	}
	mod := moderator.NewIgnoreModerator()

	wsHandler := handler.NewWebSockerHandler(*services.NewChatService(repo, notif, mod))
	r.GET("/ws", wsHandler.Handle)

	err = r.Run(":8081")
	if err != nil {
		panic(fmt.Errorf("failed to start the server %v", err))
	}
}

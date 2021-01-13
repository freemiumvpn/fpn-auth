package client

import (
	"context"
	"net/http"

	"github.com/freemiumvpn/fpn-auth/internal/client/handlers"
	"github.com/freemiumvpn/fpn-auth/internal/client/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type (
	// Client Auth
	Client struct {
		server *http.Server
	}
)

// New configs
func New(ctx context.Context, address string) *Client {
	gin.SetMode(gin.ReleaseMode)
	handler := gin.New()
	middlewares.Bind(handler)
	handlers.Bind(handler)

	server := &http.Server{
		Addr:    address,
		Handler: handler,
	}

	client := &Client{}
	client.server = server
	return client
}

// Server instance
func (client Client) Serve() error {
	logrus.WithFields(logrus.Fields{
		"address": client.server.Addr,
	}).Info("[ Client ] Listening")

	return client.server.ListenAndServe()
}

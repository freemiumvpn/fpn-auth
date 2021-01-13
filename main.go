package main

import (
	"os"

	"github.com/freemiumvpn/fpn-auth/internal/client"
	"github.com/freemiumvpn/fpn-auth/shared/logger"
	"github.com/freemiumvpn/fpn-auth/shared/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
)

const (
	appName    = "fpn-auth"
	appVersion = "0.1.0"
)

var (
	config struct {
		LogLevel          string
		LogOutput         string
		AddressPrometheus string
		AddressClient     string
	}

	appFlags = []cli.Flag{
		&cli.StringFlag{
			Name:        "log-level",
			Usage:       "Log Level",
			EnvVars:     []string{"LOG_LEVEL"},
			Value:       "info",
			Destination: &config.LogLevel,
		},
		&cli.StringFlag{
			Name:        "log-output",
			Usage:       "Log Output (text, json)",
			EnvVars:     []string{"LOG_OUTPUT"},
			Value:       "json",
			Destination: &config.LogOutput,
		},
		&cli.StringFlag{
			Name:        "prometheus-address",
			Usage:       "Prometheus Address exposes '/__/metrics' ",
			EnvVars:     []string{"PROMETHEUS_ADDRESS"},
			Value:       ":8081",
			Destination: &config.AddressPrometheus,
		},
		&cli.StringFlag{
			Name:        "client-address",
			Usage:       "Client Address exposes the vpn wrapper",
			EnvVars:     []string{"CLIENT_ADDRESS"},
			Value:       ":8989",
			Destination: &config.AddressClient,
		},
	}
)

func appAction(cliCtx *cli.Context) error {
	if err := logger.New(config.LogLevel, config.LogOutput); err != nil {
		return err
	}

	ctx := cliCtx.Context
	errorGroup, ctx := errgroup.WithContext(ctx)

	metrics := prometheus.New(ctx, config.AddressPrometheus)
	authClient := client.New(ctx, config.AddressClient)

	errorGroup.Go(func() error {
		return metrics.Serve(ctx)
	})

	errorGroup.Go(func() error {
		return authClient.Serve()
	})

	return errorGroup.Wait()
}

func main() {
	app := cli.NewApp()
	app.Name = appName
	app.Flags = appFlags
	app.Version = appVersion
	app.Action = appAction

	if err := app.Run(os.Args); err != nil {
		logrus.WithError(err).Fatal("Failed to run service")
	}
}

package main

import (
	"context"

	"github.com/omegaatt36/logaggd/app"
	"github.com/omegaatt36/logaggd/app/logaggd"

	"github.com/urfave/cli"
)

// Main starts process in cli.
func Main(ctx context.Context, c *cli.Context) {
	server := logaggd.Server{}
	server.Start(ctx, c.String("listen-addr"))
}

func main() {
	app := app.App{
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "listen-addr",
				Value: ":7800",
			},
		},
		Main: Main,
	}

	app.Run()
}

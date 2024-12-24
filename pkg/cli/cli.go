package cli

import (
	"context"

	"github.com/carlosonunez/flight-summarizer/pkg/server"
	"github.com/urfave/cli/v3"
)

// ServerCLI creates the top-level CLI structure for the summarizer server.
func ServerCLI() *cli.Command {
	return &cli.Command{
		Name:  "summarizer",
		Usage: "Flight summarizer!",
		Commands: []*cli.Command{
			{
				Name:  "serve",
				Usage: "Start the Flight Summarizer server",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "port",
						Aliases: []string{"p"},
						Usage:   "Defines the server port at `PORT`",
						Value:   8080,
						Sources: cli.EnvVars("SERVER_PORT"),
					},
					&cli.StringFlag{
						Name:    "host",
						Aliases: []string{"a"},
						Usage:   "Defines the IP `ADDRESS` to bind against",
						Value:   "127.0.0.1",
						Sources: cli.EnvVars("SERVER_HOST"),
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return server.Start(&server.ServerOptions{
						Port:          int(cmd.Int("port")),
						ListenAddress: cmd.String("host"),
					})
				},
			},
		},
	}
}

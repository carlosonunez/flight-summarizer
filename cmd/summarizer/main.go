package main

import (
	"context"
	"log"
	"os"

	"github.com/carlosonunez/flight-summarizer/pkg/cli"
)

func main() {
	if err := cli.ServerCLI().Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

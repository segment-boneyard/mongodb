package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/segmentio/kit"
	"github.com/segmentio/kit/config"
	"github.com/segmentio/kit/schema"
)

const (
	Version = "0.0.1"
)

func main() {
	app := cli.NewApp()
	app.Name = "kit"
	app.Version = Version
	app.Flags = []cli.Flag{
		cli.HelpFlag,
		cli.BoolFlag{
			Name:  "debug, D",
			Usage: "enable debug logging level",
		},
	}
	app.Commands = []cli.Command{
		GetCreateCommand(),
	}
	app.Before = func(c *cli.Context) error {
		done := make(chan bool, 1)
		config.SetProviders([]config.ProviderType{})
		go kit.Run(schema.Service{Name: "kit-cli", Version: Version, Handler: func() {}}, done)
		<-done
		return nil
	}
	app.Run(os.Args)
}

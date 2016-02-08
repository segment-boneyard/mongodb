package main

import (
	"os"
	"path"

	"github.com/codegangsta/cli"
	"github.com/segmentio/kit/log"
)

func GetCreateCommand() cli.Command {
	return cli.Command{
		Name:  "create",
		Usage: "scaffolds a new kit project",
		Before: func(c *cli.Context) error {
			if !c.IsSet("name") {
				return log.Error("Missing -name flag")
			}
			return nil
		},
		Action: func(c *cli.Context) {
			var err error
			directoryPath := path.Clean(c.String("dir"))
			serviceName := c.String("name")

			// If directory is empty, get the current working directory
			// and point the directoryPath
			if directoryPath == "" {
				directoryPath, err = os.Getwd()
				log.Debugf("No directory was set, creating new folder in %s", directoryPath)
				if err != nil {
					log.Fatalf("Error reading current working directory: %s", err.Error())
					return
				}
				directoryPath = path.Join(directoryPath, serviceName)
				log.Debugf("Set path to %s", directoryPath)
			}

			// Create directoryPath
			if err := os.MkdirAll(directoryPath, 0777); err != nil {
				// If the directory flag was not given, and we tried to create
				// a new directory, but that already exists, abort.
				if !c.IsSet("dir") && err == os.ErrExist {
					log.Errorf("The directory %s already exists, aborting", directoryPath)
					return
				} else {
					log.Errorf("Error creating the directory: %s", err.Error())
					return
				}
			}

			// Change the current working directory
			log.Debug("Changing current working directory to %s", directoryPath)
			os.Chdir(directoryPath)
		},
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "name",
				Usage: "service name",
			},
			cli.StringFlag{
				Name:  "dir, d",
				Usage: "directory where to create the service",
			},
			cli.StringFlag{
				Name:  "version",
				Usage: "service version",
			},
		},
	}
}

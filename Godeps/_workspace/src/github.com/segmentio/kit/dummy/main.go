package main

import (
	"time"

	"github.com/segmentio/kit"
	"github.com/segmentio/kit/config"
	"github.com/segmentio/kit/log"
	"github.com/segmentio/kit/schema"
)

func configAbcGetter() string {
	return config.Get("abc").(string)
}

func handler() {
	for {
		time.Sleep(1 * time.Second)
		log.Debug("I am a debug")
	}
}

func main() {
	kit.Run(schema.Service{
		Name:    "example",
		Version: "1.0.0",
		Handler: handler,
	})
}

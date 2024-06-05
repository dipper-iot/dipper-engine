package main

import (
	"github.com/dipper-iot/dipper-engine/engine"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	c := engine.New()
	if err := c.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

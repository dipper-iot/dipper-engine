package main

import (
	"github.com/dipper-iot/dipper-engine/engine"
	"log"
	"os"
)

func main() {
	c := engine.New()
	if err := c.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

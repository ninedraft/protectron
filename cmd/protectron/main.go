package main

import (
	"log"

	"github.com/ninedraft/protectron/pkg/protectron"
)

func main() {
	if err := protectron.Cli().Execute(); err != nil {
		log.Fatal(err)
	}
}

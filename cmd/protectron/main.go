package main

import (
	"log"

	"github.com/ninedraft/protectron/pkg/cli"
)

func main() {
	if err := cli.Command().Execute(); err != nil {
		log.Fatal(err)
	}
}

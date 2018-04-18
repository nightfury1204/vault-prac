package main

import (
	"os"

	"github.com/nightfury1204/vault-prac/cert-issuer/commands"
)

func main() {
	if err := commands.NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}

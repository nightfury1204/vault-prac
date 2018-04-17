package main

import (
	"os"

	"github.com/nightfury1204/vault-prac/vault-bootstrapper/commands"
)

func main() {
	if err := commands.NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}

// Package main starts the authz-checker
package main

import (
	"github.com/Stoakes/authz-checker/cmd"
	"github.com/Stoakes/go-pkg/log"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Bg().Fatal(err.Error())
	}
}

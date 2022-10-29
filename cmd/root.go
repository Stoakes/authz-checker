// Package cmd gathers all authz-checker commands
package cmd

import (
	"github.com/spf13/cobra"
)

// RootCmd is the entrypoint for every vault-exporter command
var RootCmd = &cobra.Command{
	Use:               "authz-checker",
	Short:             "Are you allowed to access this ? Let me tell you.",
	SilenceUsage:      true,
	DisableAutoGenTag: true,
	Long:              `.`,
}

func init() {
	RootCmd.AddCommand(startCmd)
}

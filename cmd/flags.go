package cmd

import (
	"fmt"
	"os"

	tocoldp "github.com/sfborg/to-coldp/pkg"
	"github.com/sfborg/to-coldp/pkg/config"
	"github.com/spf13/cobra"
)

type flagFunc func(cmd *cobra.Command)

func versionFlag(cmd *cobra.Command) {
	b, _ := cmd.Flags().GetBool("version")
	if b {
		version := tocoldp.GetVersion()
		fmt.Printf(
			"\nVersion: %s\nBuild:   %s\n\n",
			version.Version,
			version.Build,
		)
		os.Exit(0)
	}
}

func nameUsageFlag(cmd *cobra.Command) {
	b, _ := cmd.Flags().GetBool("version")
	if b {
		opts = append(opts, config.OptWithNameUsage(b))
	}
}

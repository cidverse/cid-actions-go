package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func getRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   `cid-central`,
		Short: `cid is a cli to run pipeline actions locally and as part of your ci/cd process`,
		Long:  `cid is a cli to run pipeline actions locally and as part of your ci/cd process`,
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
			os.Exit(0)
		},
	}

	rootCmd.AddCommand(runCmd())
	rootCmd.AddCommand(metadataCmd())

	return rootCmd
}

// Execute executes the root command.
func Execute() error {
	rootCmd := getRootCmd()
	return rootCmd.Execute()
}

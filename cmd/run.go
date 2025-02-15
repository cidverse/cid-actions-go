package cmd

import (
	"fmt"
	"os"

	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/spf13/cobra"
)

func runCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: `executes the specified action`,
		Run: func(cmd *cobra.Command, args []string) {
			// sdk
			sdk, sdkErr := cidsdk.NewSDK()
			if sdkErr != nil {
				fmt.Println("Fatal: Failed to initialize SDK: " + sdkErr.Error())
				os.Exit(1)
			}

			// actions
			var actions = getActions(sdk)

			// execute
			action := actions[args[0]]
			if action == nil {
				fmt.Printf("Fatal: action %s is not known!", args[0])
				os.Exit(1)
			}

			err := action.Execute()
			if err != nil {
				fmt.Printf("Fatal: action encountered an error, %s", err.Error())
				os.Exit(1)
			}
		},
	}

	return cmd
}

package cmd

import (
	"fmt"
	"github.com/cidverse/cid-actions-go/actions/api"
	"github.com/cidverse/cid-actions-go/actions/gobuild"
	"github.com/cidverse/cid-actions-go/actions/golintgolangci"
	"github.com/cidverse/cid-actions-go/actions/gotest"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
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
		actions := map[string]api.Action{
			// golang
			"go-build":         gobuild.Action{Sdk: *sdk},
			"go-test":          gotest.Action{Sdk: *sdk},
			"go-lint-golangci": golintgolangci.Action{Sdk: *sdk},
		}

		// execute
		action := actions[args[0]]
		if action != nil {
			err := action.Execute()
			if err != nil {
				fmt.Println("Fatal: Actions returned error status" + err.Error())
				os.Exit(1)
			}
		}
	},
}

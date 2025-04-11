package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/spf13/cobra"
)

func metadataCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cid-metadata",
		Short: `returns metadata about all available actions`,
		Run: func(cmd *cobra.Command, args []string) {
			// sdk
			sdk, sdkErr := cidsdk.NewSDK()
			if sdkErr != nil {
				fmt.Println("Fatal: Failed to initialize SDK: " + sdkErr.Error())
				os.Exit(1)
			}

			// actions
			var actions = getActions(sdk)
			var metadata []cidsdk.ActionMetadata

			for _, action := range actions {
				metadata = append(metadata, action.Metadata())
			}

			// output as json
			output, err := json.MarshalIndent(metadata, "", "  ")
			if err != nil {
				fmt.Println("Fatal: Failed to marshal metadata: " + err.Error())
				os.Exit(1)
			}

			fmt.Println(string(output))
		},
	}

	return cmd
}

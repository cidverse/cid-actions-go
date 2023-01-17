package cmd

import (
	"fmt"
	"os"

	"github.com/cidverse/cid-actions-go/actions/api"
	"github.com/cidverse/cid-actions-go/actions/container"
	"github.com/cidverse/cid-actions-go/actions/fossa"
	"github.com/cidverse/cid-actions-go/actions/ggshield"
	"github.com/cidverse/cid-actions-go/actions/gitleaks"
	"github.com/cidverse/cid-actions-go/actions/golang"
	"github.com/cidverse/cid-actions-go/actions/gosec"
	"github.com/cidverse/cid-actions-go/actions/helm"
	"github.com/cidverse/cid-actions-go/actions/hugo"
	"github.com/cidverse/cid-actions-go/actions/java"
	"github.com/cidverse/cid-actions-go/actions/mkdocs"
	"github.com/cidverse/cid-actions-go/actions/node"
	"github.com/cidverse/cid-actions-go/actions/python"
	"github.com/cidverse/cid-actions-go/actions/sonarqube"
	"github.com/cidverse/cid-actions-go/actions/syft"
	"github.com/cidverse/cid-actions-go/actions/techdocs"
	"github.com/cidverse/cid-actions-go/actions/upx"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/spf13/cobra"
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
			// changeloggenerate
			// container
			"buildah-build":   container.BuildahBuildAction{Sdk: *sdk},
			"buildah-publish": container.BuildahPublishAction{Sdk: *sdk},
			// fossa
			"fossa-scan": fossa.SourceScanAction{Sdk: *sdk},
			// ggshield
			"ggshield-scan": ggshield.ScanAction{Sdk: *sdk},
			// gitleaks
			"gitleaks-scan": gitleaks.ScanAction{Sdk: *sdk},
			// golang
			"go-build": golang.BuildAction{Sdk: *sdk},
			"go-test":  golang.TestAction{Sdk: *sdk},
			"go-lint":  golang.LintAction{Sdk: *sdk},
			// gosec
			"gosec-scan": gosec.ScanAction{Sdk: *sdk},
			// helm
			"helm-build":         helm.BuildAction{Sdk: *sdk},
			"helm-lint":          helm.LintAction{Sdk: *sdk},
			"helm-publish-nexus": helm.PublishNexusAction{Sdk: *sdk},
			// java
			"java-build":   java.BuildAction{Sdk: *sdk},
			"java-test":    java.TestAction{Sdk: *sdk},
			"java-publish": java.PublishAction{Sdk: *sdk},
			// node
			"node-build": node.BuildAction{Sdk: *sdk},
			"node-test":  node.TestAction{Sdk: *sdk},
			// python
			"python-build": python.BuildAction{Sdk: *sdk},
			"python-test":  python.TestAction{Sdk: *sdk},
			"python-lint":  python.LintAction{Sdk: *sdk},
			// sonarqube
			"sonarqube-scan": sonarqube.ScanAction{Sdk: *sdk},
			// syft
			"syft-sbom-build":   syft.BuildAction{Sdk: *sdk},
			"grype-sbom-report": syft.ReportAction{Sdk: *sdk},
			// mkdocs
			"mkdocs-start": mkdocs.StartAction{Sdk: *sdk},
			"mkdocs-build": mkdocs.BuildAction{Sdk: *sdk},
			// hugo
			"hugo-start": hugo.StartAction{Sdk: *sdk},
			"hugo-build": hugo.BuildAction{Sdk: *sdk},
			// techdocs
			"techdocs-start":   techdocs.StartAction{Sdk: *sdk},
			"techdocs-build":   techdocs.BuildAction{Sdk: *sdk},
			"techdocs-publish": techdocs.PublishAction{Sdk: *sdk},
			// trivy
			// TODO: "trivy-scan": trivy.ScanAction{Sdk: *sdk},
			// upx-optimize
			"opx-optimize": upx.OptimizeAction{Sdk: *sdk},
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

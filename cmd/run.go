package cmd

import (
	"fmt"
	"os"

	"github.com/cidverse/cid-actions-go/actions/api"
	"github.com/cidverse/cid-actions-go/actions/applicationinspector"
	"github.com/cidverse/cid-actions-go/actions/changeloggenerate"
	"github.com/cidverse/cid-actions-go/actions/container"
	"github.com/cidverse/cid-actions-go/actions/cosign"
	"github.com/cidverse/cid-actions-go/actions/fossa"
	"github.com/cidverse/cid-actions-go/actions/ggshield"
	"github.com/cidverse/cid-actions-go/actions/github"
	"github.com/cidverse/cid-actions-go/actions/gitlab"
	"github.com/cidverse/cid-actions-go/actions/gitleaks"
	"github.com/cidverse/cid-actions-go/actions/golang"
	"github.com/cidverse/cid-actions-go/actions/gosec"
	"github.com/cidverse/cid-actions-go/actions/helm"
	"github.com/cidverse/cid-actions-go/actions/hugo"
	"github.com/cidverse/cid-actions-go/actions/java"
	"github.com/cidverse/cid-actions-go/actions/mkdocs"
	"github.com/cidverse/cid-actions-go/actions/node"
	"github.com/cidverse/cid-actions-go/actions/ossf"
	"github.com/cidverse/cid-actions-go/actions/python"
	"github.com/cidverse/cid-actions-go/actions/qodana"
	"github.com/cidverse/cid-actions-go/actions/semgrep"
	"github.com/cidverse/cid-actions-go/actions/sonarqubescan"
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
			"changelog-generate": changeloggenerate.GenerateAction{Sdk: *sdk},
			// container
			"buildah-build":   container.BuildahBuildAction{Sdk: *sdk},
			"buildah-publish": container.BuildahPublishAction{Sdk: *sdk},
			// cosign
			"cosign-container-sign":   cosign.SignAction{Sdk: *sdk},
			"cosign-container-attach": cosign.AttachAction{Sdk: *sdk},
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
			"sonarqube-scan": sonarqubescan.ScanAction{Sdk: *sdk},
			// syft
			"syft-container-sbom-generate": syft.ContainerGenerateAction{Sdk: *sdk},
			"syft-artifact-sbom-generate":  syft.ArtifactGenerateAction{Sdk: *sdk},
			"grype-container-sbom-report":  syft.ReportAction{Sdk: *sdk},
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
			// semgrep
			"semgrep-scan": semgrep.ScanAction{Sdk: *sdk},
			// qodana
			"qodana-scan": qodana.ScanAction{Sdk: *sdk},
			// github
			"github-sarif-upload":    github.SarifUploadAction{Sdk: *sdk},
			"github-release-publish": github.PublishAction{Sdk: *sdk},
			// gitlab
			"gitlab-release-publish": gitlab.PublishAction{Sdk: *sdk},
			// ossf
			"ossf-scorecard-scan": ossf.ScorecardScanAction{Sdk: *sdk},
			// applicationinspector
			"applicationinspector-scan": applicationinspector.ScanAction{Sdk: *sdk},
		}

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

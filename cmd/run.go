package cmd

import (
	"fmt"
	"os"

	"github.com/cidverse/cid-actions-go/actions/api"
	"github.com/cidverse/cid-actions-go/actions/applicationinspector"
	"github.com/cidverse/cid-actions-go/actions/changeloggenerate"
	"github.com/cidverse/cid-actions-go/actions/container/containerbuild"
	"github.com/cidverse/cid-actions-go/actions/container/containerpublish"
	"github.com/cidverse/cid-actions-go/actions/cosign/cosignattach"
	"github.com/cidverse/cid-actions-go/actions/cosign/cosignsign"
	"github.com/cidverse/cid-actions-go/actions/fossa/fossasourcescan"
	"github.com/cidverse/cid-actions-go/actions/ggshield"
	"github.com/cidverse/cid-actions-go/actions/github/githubpublishrelease"
	"github.com/cidverse/cid-actions-go/actions/github/githubpublishsarif"
	"github.com/cidverse/cid-actions-go/actions/gitlab/gitlabreleasepublish"
	"github.com/cidverse/cid-actions-go/actions/gitleaks/gitleaksscan"
	"github.com/cidverse/cid-actions-go/actions/golang/golangbuild"
	"github.com/cidverse/cid-actions-go/actions/golang/golanglint"
	"github.com/cidverse/cid-actions-go/actions/golang/golangtest"
	"github.com/cidverse/cid-actions-go/actions/gosec/gosecscan"
	"github.com/cidverse/cid-actions-go/actions/helm/helmbuild"
	"github.com/cidverse/cid-actions-go/actions/helm/helmlint"
	"github.com/cidverse/cid-actions-go/actions/helm/helmpublishnexus"
	"github.com/cidverse/cid-actions-go/actions/helm/helmpublishregistry"
	"github.com/cidverse/cid-actions-go/actions/hugo/hugobuild"
	"github.com/cidverse/cid-actions-go/actions/hugo/hugostart"
	"github.com/cidverse/cid-actions-go/actions/java/javabuild"
	"github.com/cidverse/cid-actions-go/actions/java/javagradlewrapperscan"
	"github.com/cidverse/cid-actions-go/actions/java/javapublish"
	"github.com/cidverse/cid-actions-go/actions/java/javatest"
	"github.com/cidverse/cid-actions-go/actions/mkdocs/mkdocsbuild"
	"github.com/cidverse/cid-actions-go/actions/mkdocs/mkdocsstart"
	"github.com/cidverse/cid-actions-go/actions/node/nodebuild"
	"github.com/cidverse/cid-actions-go/actions/node/nodetest"
	"github.com/cidverse/cid-actions-go/actions/ossf/scorecardscan"
	"github.com/cidverse/cid-actions-go/actions/python/pythonbuild"
	"github.com/cidverse/cid-actions-go/actions/python/pythonlint"
	"github.com/cidverse/cid-actions-go/actions/python/pythontest"
	"github.com/cidverse/cid-actions-go/actions/qodana/qodanascan"
	"github.com/cidverse/cid-actions-go/actions/semgrep/semgrepscan"
	"github.com/cidverse/cid-actions-go/actions/sonarqube/sonarqubescan"
	"github.com/cidverse/cid-actions-go/actions/syft/grypesbomreport"
	"github.com/cidverse/cid-actions-go/actions/syft/syftartifactsbomgenerate"
	"github.com/cidverse/cid-actions-go/actions/syft/syftcontainersbombuild"
	"github.com/cidverse/cid-actions-go/actions/techdocs/techdocsbuild"
	"github.com/cidverse/cid-actions-go/actions/techdocs/techdocspublish"
	"github.com/cidverse/cid-actions-go/actions/techdocs/techdocsstart"
	"github.com/cidverse/cid-actions-go/actions/upx/upxoptimize"
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
			"changelog-generate": changeloggenerate.Action{Sdk: *sdk},
			// container
			"buildah-build":   containerbuild.Action{Sdk: *sdk},
			"buildah-publish": containerpublish.Action{Sdk: *sdk},
			// cosign
			"cosign-container-sign":   cosignsign.SignAction{Sdk: *sdk},
			"cosign-container-attach": cosignattach.AttachAction{Sdk: *sdk},
			// fossa
			"fossa-scan": fossasourcescan.Action{Sdk: *sdk},
			// ggshield
			"ggshield-scan": ggshield.Action{Sdk: *sdk},
			// gitleaks
			"gitleaks-scan": gitleaksscan.Action{Sdk: *sdk},
			// golang
			"go-build": golangbuild.BuildAction{Sdk: *sdk},
			"go-test":  golangtest.TestAction{Sdk: *sdk},
			"go-lint":  golanglint.Action{Sdk: *sdk},
			// gosec
			"gosec-scan": gosecscan.ScanAction{Sdk: *sdk},
			// helm
			"helm-build":            helmbuild.BuildAction{Sdk: *sdk},
			"helm-lint":             helmlint.LintAction{Sdk: *sdk},
			"helm-publish-nexus":    helmpublishnexus.PublishNexusAction{Sdk: *sdk},
			"helm-publish-registry": helmpublishregistry.PublishRegistryAction{Sdk: *sdk},
			// java
			"java-build":               javabuild.Action{Sdk: *sdk},
			"java-test":                javatest.Action{Sdk: *sdk},
			"java-publish":             javapublish.Action{Sdk: *sdk},
			"java-gradle-wrapper-scan": javagradlewrapperscan.Action{Sdk: *sdk},
			// node
			"node-build": nodebuild.Action{Sdk: *sdk},
			"node-test":  nodetest.Action{Sdk: *sdk},
			// python
			"python-build": pythonbuild.BuildAction{Sdk: *sdk},
			"python-test":  pythontest.TestAction{Sdk: *sdk},
			"python-lint":  pythonlint.LintAction{Sdk: *sdk},
			// sonarqube
			"sonarqube-scan": sonarqubescan.Action{Sdk: *sdk},
			// syft
			"syft-container-sbom-generate": syftcontainersbombuild.Action{Sdk: *sdk},
			"syft-artifact-sbom-generate":  syftartifactsbomgenerate.Action{Sdk: *sdk},
			"grype-container-sbom-report":  grypesbomreport.Action{Sdk: *sdk},
			// mkdocs
			"mkdocs-start": mkdocsstart.StartAction{Sdk: *sdk},
			"mkdocs-build": mkdocsbuild.BuildAction{Sdk: *sdk},
			// hugo
			"hugo-start": hugostart.Action{Sdk: *sdk},
			"hugo-build": hugobuild.Action{Sdk: *sdk},
			// techdocs
			"techdocs-start":   techdocsstart.Action{Sdk: *sdk},
			"techdocs-build":   techdocsbuild.Action{Sdk: *sdk},
			"techdocs-publish": techdocspublish.Action{Sdk: *sdk},
			// trivy
			// TODO: "trivy-scan": trivy.Action{Sdk: *sdk},
			// upx-optimize
			"opx-optimize": upxoptimize.OptimizeAction{Sdk: *sdk},
			// semgrep
			"semgrep-scan": semgrepscan.ScanAction{Sdk: *sdk},
			// qodana
			"qodana-scan": qodanascan.ScanAction{Sdk: *sdk},
			// github
			"github-sarif-upload":    githubpublishsarif.Action{Sdk: *sdk},
			"github-release-publish": githubpublishrelease.Action{Sdk: *sdk},
			// gitlab
			"gitlab-release-publish": gitlabreleasepublish.PublishAction{Sdk: *sdk},
			// ossf
			"ossf-scorecard-scan": scorecardscan.Action{Sdk: *sdk},
			// applicationinspector
			"applicationinspector-scan": applicationinspector.Action{Sdk: *sdk},
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

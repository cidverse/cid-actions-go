package cmd

import (
	"github.com/cidverse/cid-actions-go/actions/ansible/ansibledeploy"
	"github.com/cidverse/cid-actions-go/actions/ansible/ansiblelint"
	"github.com/cidverse/cid-actions-go/actions/applicationinspector"
	"github.com/cidverse/cid-actions-go/actions/changeloggenerate"
	"github.com/cidverse/cid-actions-go/actions/container/containerbuild"
	"github.com/cidverse/cid-actions-go/actions/container/containerpublish"
	"github.com/cidverse/cid-actions-go/actions/cosign/cosignattach"
	"github.com/cidverse/cid-actions-go/actions/cosign/cosignsign"
	"github.com/cidverse/cid-actions-go/actions/donet/dotnetbuild"
	"github.com/cidverse/cid-actions-go/actions/donet/dotnettest"
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
	"github.com/cidverse/cid-actions-go/actions/gradle/gradlebuild"
	"github.com/cidverse/cid-actions-go/actions/gradle/gradlepublish"
	"github.com/cidverse/cid-actions-go/actions/gradle/gradletest"
	"github.com/cidverse/cid-actions-go/actions/helm/helmbuild"
	"github.com/cidverse/cid-actions-go/actions/helm/helmdeploy"
	"github.com/cidverse/cid-actions-go/actions/helm/helmlint"
	"github.com/cidverse/cid-actions-go/actions/helm/helmpublishnexus"
	"github.com/cidverse/cid-actions-go/actions/helm/helmpublishregistry"
	"github.com/cidverse/cid-actions-go/actions/helmfile/helmfiledeploy"
	"github.com/cidverse/cid-actions-go/actions/helmfile/helmfilelint"
	"github.com/cidverse/cid-actions-go/actions/hugo/hugobuild"
	"github.com/cidverse/cid-actions-go/actions/hugo/hugostart"
	"github.com/cidverse/cid-actions-go/actions/maven/mavenbuild"
	"github.com/cidverse/cid-actions-go/actions/maven/mavenpublish"
	"github.com/cidverse/cid-actions-go/actions/maven/maventest"
	"github.com/cidverse/cid-actions-go/actions/mkdocs/mkdocsbuild"
	"github.com/cidverse/cid-actions-go/actions/mkdocs/mkdocsstart"
	"github.com/cidverse/cid-actions-go/actions/node/nodebuild"
	"github.com/cidverse/cid-actions-go/actions/node/nodetest"
	"github.com/cidverse/cid-actions-go/actions/ossf/scorecardscan"
	"github.com/cidverse/cid-actions-go/actions/python/pythonbuild"
	"github.com/cidverse/cid-actions-go/actions/python/pythonlint"
	"github.com/cidverse/cid-actions-go/actions/python/pythontest"
	"github.com/cidverse/cid-actions-go/actions/qodana/qodanascan"
	"github.com/cidverse/cid-actions-go/actions/renovate/renovatelint"
	"github.com/cidverse/cid-actions-go/actions/rust/rustbuild"
	"github.com/cidverse/cid-actions-go/actions/rust/rusttest"
	"github.com/cidverse/cid-actions-go/actions/semgrep/semgrepscan"
	"github.com/cidverse/cid-actions-go/actions/sonarqube/sonarqubescan"
	"github.com/cidverse/cid-actions-go/actions/syft/grypesbomreport"
	"github.com/cidverse/cid-actions-go/actions/syft/syftartifactsbomgenerate"
	"github.com/cidverse/cid-actions-go/actions/syft/syftcontainersbombuild"
	"github.com/cidverse/cid-actions-go/actions/techdocs/techdocsbuild"
	"github.com/cidverse/cid-actions-go/actions/techdocs/techdocspublish"
	"github.com/cidverse/cid-actions-go/actions/techdocs/techdocsstart"
	"github.com/cidverse/cid-actions-go/actions/upx/upxoptimize"
	"github.com/cidverse/cid-actions-go/actions/zizmor/zizmorscan"
	cidsdk "github.com/cidverse/cid-sdk-go"
)

func getActions(sdk *cidsdk.SDK) map[string]cidsdk.Action {
	actions := []cidsdk.Action{
		// ansible
		ansiblelint.Action{Sdk: *sdk},
		ansibledeploy.Action{Sdk: *sdk},
		// changeloggenerate
		changeloggenerate.Action{Sdk: *sdk},
		// container
		containerbuild.Action{Sdk: *sdk},
		containerpublish.Action{Sdk: *sdk},
		// cosign
		cosignsign.Action{Sdk: *sdk},
		cosignattach.Action{Sdk: *sdk},
		// dotnet
		dotnetbuild.Action{Sdk: *sdk},
		dotnettest.Action{Sdk: *sdk},
		// fossa
		fossasourcescan.Action{Sdk: *sdk},
		// ggshield
		ggshield.Action{Sdk: *sdk},
		// gitleaks
		gitleaksscan.Action{Sdk: *sdk},
		// golang
		golangbuild.Action{Sdk: *sdk},
		golangtest.Action{Sdk: *sdk},
		golanglint.Action{Sdk: *sdk},
		// gosec
		gosecscan.Action{Sdk: *sdk},
		// helm
		helmbuild.Action{Sdk: *sdk},
		helmlint.Action{Sdk: *sdk},
		helmpublishnexus.Action{Sdk: *sdk},
		helmpublishregistry.Action{Sdk: *sdk},
		helmdeploy.Action{Sdk: *sdk},
		// helmfile
		helmfilelint.Action{Sdk: *sdk},
		helmfiledeploy.Action{Sdk: *sdk},
		// gradle
		gradlebuild.Action{Sdk: *sdk},
		gradletest.Action{Sdk: *sdk},
		gradlepublish.Action{Sdk: *sdk},
		// maven
		mavenbuild.Action{Sdk: *sdk},
		maventest.Action{Sdk: *sdk},
		mavenpublish.Action{Sdk: *sdk},
		// node
		nodebuild.Action{Sdk: *sdk},
		nodetest.Action{Sdk: *sdk},
		// python
		pythonbuild.Action{Sdk: *sdk},
		pythontest.Action{Sdk: *sdk},
		pythonlint.Action{Sdk: *sdk},
		// renovate
		renovatelint.Action{Sdk: *sdk},
		// rust
		rustbuild.Action{Sdk: *sdk},
		rusttest.Action{Sdk: *sdk},
		// sonarqube
		sonarqubescan.Action{Sdk: *sdk},
		// syft
		syftcontainersbombuild.Action{Sdk: *sdk},
		syftartifactsbomgenerate.Action{Sdk: *sdk},
		grypesbomreport.Action{Sdk: *sdk},
		// mkdocs
		mkdocsstart.Action{Sdk: *sdk},
		mkdocsbuild.Action{Sdk: *sdk},
		// hugo
		hugostart.Action{Sdk: *sdk},
		hugobuild.Action{Sdk: *sdk},
		// techdocs
		techdocsstart.Action{Sdk: *sdk},
		techdocsbuild.Action{Sdk: *sdk},
		techdocspublish.Action{Sdk: *sdk},
		// trivy
		// TODO: trivy.Action{Sdk: *sdk},
		// upx-optimize
		upxoptimize.Action{Sdk: *sdk},
		// semgrep
		semgrepscan.Action{Sdk: *sdk},
		// qodana
		qodanascan.Action{Sdk: *sdk},
		// github
		githubpublishsarif.Action{Sdk: *sdk},
		githubpublishrelease.Action{Sdk: *sdk},
		// gitlab
		gitlabreleasepublish.Action{Sdk: *sdk},
		// ossf
		scorecardscan.Action{Sdk: *sdk},
		// applicationinspector
		applicationinspector.Action{Sdk: *sdk},
		// spectral
		// TODO: spectral.LintAction{Sdk: *sdk},
		// zizmor
		zizmorscan.Action{Sdk: *sdk},
	}

	actionMap := make(map[string]cidsdk.Action, len(actions))
	for _, action := range actions {
		actionMap[action.Metadata().Name] = action
	}

	return actionMap
}

package sonarqube

import (
	"path"
	"strings"

	cidsdk "github.com/cidverse/cid-sdk-go"
)

type ScanAction struct {
	Sdk cidsdk.SDKClient
}

type ScanConfig struct {
	SonarHostURL       string `json:"sonar_host_url"  env:"SONAR_HOST_URL"`
	SonarOrganization  string `json:"sonar_organization"  env:"SONAR_ORGANIZATION"`
	SonarProjectKey    string `json:"sonar_project_key"  env:"SONAR_PROJECTKEY"`
	SonarDefaultBranch string `json:"sonar_default_branch"  env:"SONAR_DEFAULT_BRANCH"`
	SonarToken         string `json:"sonar_token"  env:"SONAR_TOKEN"`
}

func (a ScanAction) Execute() (err error) {
	cfg := ScanConfig{}
	ctx, err := a.Sdk.ProjectAction(&cfg)
	if err != nil {
		return err
	}

	// default to cloud host
	if cfg.SonarHostURL == "" {
		cfg.SonarHostURL = "https://sonarcloud.io"
	}
	if cfg.SonarProjectKey == "" {
		cfg.SonarProjectKey = ctx.Env["NCI_PROJECT_ID"]
	}
	if cfg.SonarDefaultBranch == "" {
		cfg.SonarDefaultBranch = "develop"
	}

	// ensure that the default branch is configured correctly
	prepareProject(cfg.SonarHostURL, cfg.SonarToken, cfg.SonarOrganization, cfg.SonarProjectKey, ctx.Env["NCI_PROJECT_NAME"], ctx.Env["NCI_PROJECT_DESCRIPTION"], cfg.SonarDefaultBranch)

	// run scan
	var scanArgs []string
	scanArgs = append(scanArgs, `-D sonar.host.url=`+cfg.SonarHostURL)
	scanArgs = append(scanArgs, `-D sonar.login=`+cfg.SonarToken)
	if cfg.SonarOrganization != "" {
		scanArgs = append(scanArgs, `-D sonar.organization=`+cfg.SonarOrganization)
	}
	scanArgs = append(scanArgs, `-D sonar.projectKey=`+cfg.SonarProjectKey)
	scanArgs = append(scanArgs, `-D sonar.projectName=`+ctx.Env["NCI_PROJECT_NAME"])
	scanArgs = append(scanArgs, `-D sonar.branch.name=`+ctx.Env["NCI_COMMIT_REF_SLUG"])
	scanArgs = append(scanArgs, `-D sonar.sources=.`)
	scanArgs = append(scanArgs, `-D sonar.tests=.`)

	// module specific parameters
	var sourceInclusion []string
	var sourceExclusions []string
	var testInclusion []string
	var testExclusions []string
	for _, module := range ctx.Modules {
		if module.BuildSystem == string(cidsdk.BuildSystemGradle) || module.BuildSystem == string(cidsdk.BuildSystemMaven) {
			sourceInclusion = append(sourceInclusion, "**/src/main/java/**", "**/src/main/kotlin/**")
			testInclusion = append(testInclusion, "**/src/test/java/**", "**/src/test/kotlin/**")
			scanArgs = append(scanArgs, `-D sonar.coverage.jacoco.xmlReportPaths=`+path.Join(ctx.Config.ArtifactDir, "**", "test", "jacoco.xml"))
			scanArgs = append(scanArgs, `-D sonar.java.binaries=.`)
			scanArgs = append(scanArgs, `-D sonar.java.test.binaries=.`)
		} else if module.BuildSystem == string(cidsdk.BuildSystemGoMod) {
			sourceExclusions = append(sourceExclusions, "**/*_test.go", "**/vendor/**", "**/testdata/*")
			testInclusion = append(testInclusion, "**/*_test.go")
			testExclusions = append(testExclusions, "**/vendor/**")
			scanArgs = append(scanArgs, `-D sonar.go.coverage.reportPaths=`+path.Join(ctx.Config.ArtifactDir, module.Slug, "go-test", "coverage.out"))
			scanArgs = append(scanArgs, `-D sonar.go.tests.reportPaths=`+path.Join(ctx.Config.ArtifactDir, module.Slug, "go-test", "coverage.json"))
		}
	}
	scanArgs = append(scanArgs, `-D sonar.inclusions=`+strings.Join(sourceInclusion, ","))
	scanArgs = append(scanArgs, `-D sonar.exclusions=`+strings.Join(sourceExclusions, ","))
	scanArgs = append(scanArgs, `-D sonar.test.inclusions=`+strings.Join(testInclusion, ","))
	scanArgs = append(scanArgs, `-D sonar.test.exclusions=`+strings.Join(testExclusions, ","))

	_, err = a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
		Command: `sonar-scanner ` + strings.Join(scanArgs, " "),
		WorkDir: ctx.ProjectDir,
	})
	if err != nil {
		return err
	}

	return nil
}

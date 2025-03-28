package sonarqubescan

import (
	"fmt"
	"os"
	"strings"

	"github.com/cidverse/cid-actions-go/actions/sonarqube/sonarqubecommon"
	"github.com/cidverse/cid-actions-go/util"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/gosimple/slug"
)

type Action struct {
	Sdk cidsdk.SDKClient
}

type Config struct {
	SonarHostURL       string `json:"sonar_host_url"  env:"SONAR_HOST_URL"`
	SonarOrganization  string `json:"sonar_organization"  env:"SONAR_ORGANIZATION"`
	SonarProjectKey    string `json:"sonar_project_key"  env:"SONAR_PROJECTKEY"`
	SonarDefaultBranch string `json:"sonar_default_branch"  env:"SONAR_DEFAULT_BRANCH"`
	SonarToken         string `json:"sonar_token"  env:"SONAR_TOKEN"`
}

func (a Action) Metadata() cidsdk.ActionMetadata {
	return cidsdk.ActionMetadata{
		Name:        "sonarqube-scan",
		Description: "Scans the repository for security issues using SonarQube.",
		Category:    "sast",
		Scope:       cidsdk.ActionScopeProject,
		Rules: []cidsdk.ActionRule{
			{
				Type:       "cel",
				Expression: `ENV["SONAR_TOKEN"] != "" && NCI_COMMIT_REF_TYPE == "branch"`,
			},
		},
		Access: cidsdk.ActionAccess{
			Environment: []cidsdk.ActionAccessEnv{
				{
					Name:        "SONAR_HOST_URL",
					Description: `The SonarQube host URL.`,
				},
				{
					Name:        "SONAR_ORGANIZATION",
					Description: `The SonarQube organization.`,
				},
				{
					Name:        "SONAR_PROJECTKEY",
					Description: `The SonarQube project key.`,
				},
				{
					Name:        "SONAR_DEFAULT_BRANCH",
					Description: `The SonarQube default branch.`,
				},
				{
					Name:        "SONAR_TOKEN",
					Description: `The SonarQube authentication token.`,
					Required:    true,
					Secret:      true,
				},
				{
					Name:        "NCI_.*",
					Description: `The project properties sonar needs to identify the repository, commit, merge request, etc.`,
					Pattern:     true,
				},
			},
			Executables: []cidsdk.ActionAccessExecutable{
				{
					Name: "cargo",
				},
			},
			Network: []cidsdk.ActionAccessNetwork{
				{
					Host: "sonarcloud.io:443",
				},
				{
					Host: "api.sonarcloud.io:443",
				},
				{
					Host: "scanner.sonarcloud.io:443",
				},
			},
		},
	}
}

func (a Action) Execute() (err error) {
	cfg := Config{}
	ctx, err := a.Sdk.ProjectAction(&cfg)
	if err != nil {
		return err
	}

	// default to cloud host
	if cfg.SonarHostURL == "" {
		cfg.SonarHostURL = "https://sonarcloud.io"
	}
	if cfg.SonarProjectKey == "" {
		cfg.SonarProjectKey = slug.Make(ctx.Env["NCI_REPOSITORY_HOST_SERVER"]) + "-" + ctx.Env["NCI_PROJECT_ID"]
	}
	if cfg.SonarDefaultBranch == "" {
		cfg.SonarDefaultBranch = util.FirstNonEmpty([]string{ctx.Env["NCI_PROJECT_DEFAULT_BRANCH"], "main"})
	}

	// ensure that the default branch is configured correctly
	_ = a.Sdk.Log(cidsdk.LogMessageRequest{Level: "info", Message: "creating project and setting default branch if missing", Context: map[string]interface{}{"default-branch": cfg.SonarDefaultBranch, "host": cfg.SonarHostURL, "project-key": cfg.SonarProjectKey, "organization": cfg.SonarOrganization}})
	err = sonarqubecommon.PrepareProject(cfg.SonarHostURL, cfg.SonarToken, cfg.SonarOrganization, cfg.SonarProjectKey, ctx.Env["NCI_PROJECT_NAME"], ctx.Env["NCI_PROJECT_DESCRIPTION"], cfg.SonarDefaultBranch)
	if err != nil {
		return fmt.Errorf("failed to prepare sonarqube project: %w", err)
	}

	// run scan
	scanArgs := []string{
		`-D sonar.host.url=` + cfg.SonarHostURL,
		`-D sonar.token=` + cfg.SonarToken,
		`-D sonar.projectKey=` + cfg.SonarProjectKey,
		`-D sonar.projectName=` + ctx.Env["NCI_PROJECT_NAME"],
		`-D sonar.sources=.`,
	}
	if cfg.SonarOrganization != "" {
		scanArgs = append(scanArgs, `-D sonar.organization=`+cfg.SonarOrganization)
	}

	// set version
	if ctx.Env["NCI_COMMIT_REF_TYPE"] == "tag" {
		scanArgs = append(scanArgs, `-D sonar.projectVersion=`+ctx.Env["NCI_COMMIT_REF_NAME"])
	}

	// publish sarif reports to sonarqube
	_ = a.Sdk.Log(cidsdk.LogMessageRequest{Level: "debug", Message: fmt.Sprintf("query artifacts with %s", "type == \"report\"")})
	artifacts, err := a.Sdk.ArtifactList(cidsdk.ArtifactListRequest{Query: `artifact_type == "report"`})
	if err != nil {
		return fmt.Errorf("failed to list report artifacts: %w", err)
	}
	files := make(map[string][]string, 0)
	_ = a.Sdk.Log(cidsdk.LogMessageRequest{Level: "debug", Message: fmt.Sprintf("found %d reports", len(*artifacts))})
	for _, artifact := range *artifacts {
		if artifact.Format == "sarif" {
			targetFile := cidsdk.JoinPath(ctx.Config.TempDir, fmt.Sprintf("%s-%s", artifact.Module, artifact.Name))
			var dlErr = a.Sdk.ArtifactDownload(cidsdk.ArtifactDownloadRequest{
				ID:         artifact.ID,
				TargetFile: targetFile,
			})
			if dlErr != nil {
				_ = a.Sdk.Log(cidsdk.LogMessageRequest{Level: "warn", Message: "failed to retrieve sarif report", Context: map[string]interface{}{"artifact": fmt.Sprintf("%s-%s", artifact.Module, artifact.Name)}})
				continue
			}

			files["sarif"] = append(files["sarif"], targetFile)
		} else if artifact.Format == "go-coverage" {
			targetFile := cidsdk.JoinPath(ctx.Config.TempDir, fmt.Sprintf("%s-%s", artifact.Module, artifact.Name))
			var dlErr = a.Sdk.ArtifactDownload(cidsdk.ArtifactDownloadRequest{
				ID:         artifact.ID,
				TargetFile: targetFile,
			})
			if dlErr != nil {
				_ = a.Sdk.Log(cidsdk.LogMessageRequest{Level: "warn", Message: "failed to retrieve sarif report", Context: map[string]interface{}{"artifact": fmt.Sprintf("%s-%s", artifact.Module, artifact.Name)}})
				continue
			}

			if artifact.FormatVersion == "out" {
				files["go-coverage-out"] = append(files["go-coverage-out"], targetFile)
			} else if artifact.FormatVersion == "json" {
				files["go-coverage-json"] = append(files["go-coverage-json"], targetFile)
			}
		} else if artifact.Format == "jacoco" {
			targetFile := cidsdk.JoinPath(ctx.Config.TempDir, fmt.Sprintf("%s-%s", artifact.Module, artifact.Name))
			var dlErr = a.Sdk.ArtifactDownload(cidsdk.ArtifactDownloadRequest{
				ID:         artifact.ID,
				TargetFile: targetFile,
			})
			if dlErr != nil {
				_ = a.Sdk.Log(cidsdk.LogMessageRequest{Level: "warn", Message: "failed to retrieve jacoco report", Context: map[string]interface{}{"artifact": fmt.Sprintf("%s-%s", artifact.Module, artifact.Name)}})
				continue
			}

			files["java-jacoco"] = append(files["java-jacoco"], targetFile)
		}
	}
	if len(files["sarif"]) > 0 {
		scanArgs = append(scanArgs, `-D sonar.sarifReportPaths=`+strings.Join(files["sarif"], ","))
	}
	if len(files["go-coverage-out"]) > 0 {
		scanArgs = append(scanArgs, `-D sonar.go.coverage.reportPaths=`+strings.Join(files["go-coverage-out"], ","))
	}
	if len(files["go-coverage-json"]) > 0 {
		scanArgs = append(scanArgs, `-D sonar.go.tests.reportPaths=`+strings.Join(files["go-coverage-json"], ","))
	}
	if len(files["java-jacoco"]) > 0 {
		scanArgs = append(scanArgs, `-D sonar.coverage.jacoco.xmlReportPaths=`+strings.Join(files["java-jacoco"], ","))
	}

	// module specific parameters
	var sourceInclusion []string
	var sourceExclusions = []string{"**/.git/**"}
	var testInclusion []string
	var testExclusions []string
	for _, module := range ctx.Modules {
		if module.BuildSystem == string(cidsdk.BuildSystemGradle) || module.BuildSystem == string(cidsdk.BuildSystemMaven) {
			sourceInclusion = append(sourceInclusion, "**/src/main/java/**", "**/src/main/kotlin/**")
			testInclusion = append(testInclusion, "**/src/test/java/**", "**/src/test/kotlin/**")
			scanArgs = append(scanArgs, `-D sonar.java.binaries=.`)
			scanArgs = append(scanArgs, `-D sonar.java.test.binaries=.`)
		} else if module.BuildSystem == string(cidsdk.BuildSystemGoMod) {
			sourceExclusions = append(sourceExclusions, "**/*_test.go", "**/vendor/**", "**/mocks/**", "**/testdata/*")
			testInclusion = append(testInclusion, "**/*_test.go")
			testExclusions = append(testExclusions, "**/vendor/**")
		}
	}
	if len(sourceInclusion) > 0 {
		scanArgs = append(scanArgs, `-D sonar.inclusions=`+strings.Join(sourceInclusion, ","))
	}
	if len(sourceExclusions) > 0 {
		scanArgs = append(scanArgs, `-D sonar.exclusions=`+strings.Join(sourceExclusions, ","))
	}
	if len(testInclusion) > 0 {
		scanArgs = append(scanArgs, `-D sonar.test.inclusions=`+strings.Join(testInclusion, ","))
	}
	if len(testExclusions) > 0 {
		scanArgs = append(scanArgs, `-D sonar.test.exclusions=`+strings.Join(testExclusions, ","))
	}

	// merge request
	if ctx.Env["NCI_PIPELINE_TRIGGER"] == "merge_request" {
		scanArgs = append(scanArgs, `-D sonar.pullrequest.key=`+ctx.Env["NCI_MERGE_REQUEST_ID"])

		if _, ok := ctx.Env["NCI_MERGE_REQUEST_SOURCE_BRANCH_NAME"]; ok {
			scanArgs = append(scanArgs, `-D sonar.pullrequest.branch=`+ctx.Env["NCI_MERGE_REQUEST_SOURCE_BRANCH_NAME"])
		}
		if _, ok := ctx.Env["NCI_MERGE_REQUEST_TARGET_BRANCH_NAME"]; ok {
			scanArgs = append(scanArgs, `-D sonar.pullrequest.base=`+ctx.Env["NCI_MERGE_REQUEST_TARGET_BRANCH_NAME"])
		}
	} else {
		scanArgs = append(scanArgs, fmt.Sprintf(`-D sonar.branch.name=%q`, ctx.Env["NCI_COMMIT_REF_NAME"]))
	}

	// execute
	scanResult, err := a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
		Command: `sonar-scanner -X ` + strings.Join(scanArgs, " "),
		WorkDir: ctx.ProjectDir,
		Env: map[string]string{
			"SONAR_SCANNER_OPTS": strings.Join([]string{os.Getenv("CID_PROXY_JVM"), os.Getenv("SONAR_SCANNER_OPTS")}, " "),
		},
	})
	if err != nil {
		return err
	} else if scanResult.Code != 0 {
		return fmt.Errorf("sonar scan failed, exit code %d", scanResult.Code)
	}

	return nil
}

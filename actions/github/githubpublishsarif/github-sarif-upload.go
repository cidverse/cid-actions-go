package githubpublishsarif

import (
	"context"
	"fmt"
	"strings"

	"github.com/cidverse/cid-actions-go/pkg/githubapi"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/cidverse/cidverseutils/compress"
	"github.com/go-playground/validator/v10"
	"github.com/google/go-github/v73/github"
	"golang.org/x/oauth2"
)

type Action struct {
	Sdk cidsdk.SDKClient
}

type Config struct {
	GitHubToken string `json:"github_token"  env:"GITHUB_TOKEN"`
}

func (a Action) Metadata() cidsdk.ActionMetadata {
	return cidsdk.ActionMetadata{
		Name:        "github-sarif-upload",
		Description: "Uploads all SARIF reports to GitHub CodeScanning. Supports merge requests.",
		Category:    "sast",
		Scope:       cidsdk.ActionScopeProject,
		Rules: []cidsdk.ActionRule{
			{
				Type:       "cel",
				Expression: `hasPrefix(ENV["NCI_REPOSITORY_REMOTE"], "https://github.com/")`,
			},
		},
		Access: cidsdk.ActionAccess{
			Environment: []cidsdk.ActionAccessEnv{
				{
					Name:        "GITHUB_TOKEN",
					Description: "The GitHub token to use for uploading the SARIF file. This token is available automatically in the GitHub Actions environment.",
					Required:    true,
					Secret:      true,
				},
			},
			Executables: []cidsdk.ActionAccessExecutable{},
		},
		Input: cidsdk.ActionInput{
			Artifacts: []cidsdk.ActionArtifactType{
				{
					Type:   "report",
					Format: "sarif",
				},
			},
		},
	}
}

func (a Action) GetConfig(d *cidsdk.ProjectActionData) (Config, error) {
	cfg := Config{}
	cidsdk.PopulateFromEnv(&cfg, d.Env)

	// validate
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

func (a Action) Execute() (err error) {
	// query action data
	d, err := a.Sdk.ProjectActionDataV1()
	if err != nil {
		return err
	}

	// parse config
	cfg, err := a.GetConfig(d)
	if err != nil {
		return err
	}

	// properties
	organization := githubapi.GetGithubOrganizationFromRemote(d.Env["NCI_REPOSITORY_REMOTE"])
	repository := githubapi.GetGithubRepositoryFromRemote(d.Env["NCI_REPOSITORY_REMOTE"])

	// GitHub Client
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.GitHubToken},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	client := github.NewClient(tc)

	// iterate over all sarif reports
	artifacts, err := a.Sdk.ArtifactList(cidsdk.ArtifactListRequest{Query: `artifact_type == "report" && format == "sarif"`})
	if err != nil {
		return err
	}
	for _, report := range *artifacts {
		// get report content
		sarif, reportErr := a.Sdk.ArtifactDownloadByteArray(cidsdk.ArtifactDownloadByteArrayRequest{
			ID: report.ID,
		})
		if reportErr != nil {
			return fmt.Errorf("failed to load report %s: %w", report.Name, reportErr)
		}

		// encoding
		sarifEncoded, reportErr := compress.GZIPBase64EncodeBytes(sarif)
		if reportErr != nil {
			return fmt.Errorf("failed to encode sarif report (gzip/base64): %w", err)
		}

		// git reference (sarif upload with pull request ref will result in pull request comments)
		ref := d.Env["NCI_COMMIT_REF_VCS"]
		if d.Env["NCI_PIPELINE_TRIGGER"] == "merge_request" && d.Env["NCI_MERGE_REQUEST_ID"] != "" {
			ref = fmt.Sprintf("refs/pull/%s/merge", d.Env["NCI_MERGE_REQUEST_ID"])
		}

		// upload
		_ = a.Sdk.Log(cidsdk.LogMessageRequest{Level: "info", Message: "uploading sarif report to github code scanning api", Context: map[string]interface{}{"report": report.Name, "ref": ref, "commit_hash": d.Env["NCI_COMMIT_HASH"]}})
		sarifAnalysis := &github.SarifAnalysis{CommitSHA: github.Ptr(d.Env["NCI_COMMIT_HASH"]), Ref: github.Ptr(ref), Sarif: github.Ptr(sarifEncoded), CheckoutURI: github.Ptr(d.Config.ProjectDir)}
		sarifId, _, reportErr := client.CodeScanning.UploadSarif(context.Background(), organization, repository, sarifAnalysis)

		if reportErr != nil {
			// "job scheduled on GitHub side" is not an error, job just isn't completed yet
			if strings.Contains(reportErr.Error(), "job scheduled on GitHub side") {
				_ = a.Sdk.Log(cidsdk.LogMessageRequest{Level: "info", Message: "sarif upload successful", Context: map[string]interface{}{"report": report.Name, "state": "github_job_pending"}})
			} else {
				return fmt.Errorf("failed to upload sarif to github code-scanning api: %s", reportErr.Error())
			}
		} else if sarifId != nil {
			_ = a.Sdk.Log(cidsdk.LogMessageRequest{Level: "info", Message: "sarif upload successful", Context: map[string]interface{}{"report": report.Name, "state": "ok", "id": *sarifId.ID, "url": *sarifId.URL}})
		}
	}

	return nil
}

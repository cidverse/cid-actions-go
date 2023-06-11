package githubpublishsarif

import (
	"context"
	"fmt"
	"strings"

	"github.com/cidverse/cid-actions-go/pkg/githubapi"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/cidverse/cidverseutils/pkg/encoding"
	"github.com/google/go-github/v53/github"
	"golang.org/x/oauth2"
)

type Action struct {
	Sdk cidsdk.SDKClient
}

type Config struct {
	GitHubToken string `json:"github_token"  env:"GITHUB_TOKEN"`
}

func (a Action) Execute() (err error) {
	cfg := Config{}
	ctx, err := a.Sdk.ProjectAction(&cfg)
	if err != nil {
		return err
	}

	// properties
	organization := githubapi.GetGithubOrganizationFromRemote(ctx.Env["NCI_REPOSITORY_REMOTE"])
	repository := githubapi.GetGithubRepositoryFromRemote(ctx.Env["NCI_REPOSITORY_REMOTE"])

	// GitHub Client
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.GitHubToken},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	client := github.NewClient(tc)

	// iterate over all sarif reports
	artifacts, err := a.Sdk.ArtifactList(cidsdk.ArtifactListRequest{Query: `artifact_type == "report" && format == "sarif" && format_version == "2.1.0"`})
	if err != nil {
		return err
	}
	for _, report := range *artifacts {
		// get report content
		sarif, err := a.Sdk.ArtifactDownloadByteArray(cidsdk.ArtifactDownloadByteArrayRequest{
			ID: report.ID,
		})
		if err != nil {
			return fmt.Errorf("failed to load report %s", report.Name)
		}

		// encoding
		sarifEncoded, err := encoding.GZIPBase64EncodeBytes(sarif)
		if err != nil {
			return fmt.Errorf("failed to encode sarif report (gzip/base64)")
		}

		// git reference (sarif upload with pull request ref will result in pull request comments)
		ref := ctx.Env["NCI_COMMIT_REF_VCS"]
		if ctx.Env["NCI_PIPELINE_TRIGGER"] == "merge_request" && ctx.Env["NCI_MERGE_REQUEST_ID"] != "" {
			ref = fmt.Sprintf("refs/pull/%s/merge", ctx.Env["NCI_MERGE_REQUEST_ID"])
		}

		// upload
		_ = a.Sdk.Log(cidsdk.LogMessageRequest{Level: "info", Message: "uploading sarif report to github code scanning api", Context: map[string]interface{}{"report": report.Name, "ref": ref, "commit_hash": ctx.Env["NCI_COMMIT_HASH"]}})
		sarifAnalysis := &github.SarifAnalysis{CommitSHA: github.String(ctx.Env["NCI_COMMIT_HASH"]), Ref: github.String(ref), Sarif: github.String(sarifEncoded), CheckoutURI: github.String(ctx.Config.ProjectDir)}
		sarifId, _, err := client.CodeScanning.UploadSarif(context.Background(), organization, repository, sarifAnalysis)
		if err != nil {
			// "job scheduled on GitHub side" is not a error, job just isn't completed yet
			if strings.Contains(err.Error(), "job scheduled on GitHub side") {
				_ = a.Sdk.Log(cidsdk.LogMessageRequest{Level: "info", Message: "sarif upload successful", Context: map[string]interface{}{"report": report.Name, "state": "github_job_pending"}})
			} else {
				return fmt.Errorf("failed to upload sarif to github code-scanning api: %s", err.Error())
			}
		} else {
			_ = a.Sdk.Log(cidsdk.LogMessageRequest{Level: "info", Message: "sarif upload successful", Context: map[string]interface{}{"report": report.Name, "state": "ok", "id": *sarifId.ID, "url": *sarifId.URL}})
		}
	}

	return nil
}

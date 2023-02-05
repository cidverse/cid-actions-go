package github

import (
	"context"
	"fmt"

	"github.com/cidverse/cid-actions-go/pkg/githubapi"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/cidverse/cidverseutils/pkg/encoding"
	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
)

type SarifUploadAction struct {
	Sdk cidsdk.SDKClient
}

type SarifUploadConfig struct {
	GitHubToken string `json:"github_token"  env:"GITHUB_TOKEN"`
}

func (a SarifUploadAction) Execute() (err error) {
	cfg := SarifUploadConfig{}
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
	artifacts, err := a.Sdk.ArtifactList(cidsdk.ArtifactListRequest{ArtifactType: "report", Format: "sarif", FormatVersion: "2.1.0"})
	if err != nil {
		return err
	}
	for _, report := range *artifacts {
		// get report content
		sarif, err := a.Sdk.ArtifactDownloadByteArray(cidsdk.ArtifactDownloadByteArrayRequest{
			Module: report.Module,
			Type:   string(report.Type),
			Name:   report.Name,
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
		if ctx.Env["NCI_PIPELINE_TRIGGER"] == "pull_request" && len(ctx.Env["NCI_PIPELINE_PULL_REQUEST_ID"]) > 0 {
			ref = fmt.Sprintf("refs/pull/%s/merge", ctx.Env["NCI_PIPELINE_PULL_REQUEST_ID"])
		}

		// upload
		sarifAnalysis := &github.SarifAnalysis{CommitSHA: github.String(ctx.Env["NCI_COMMIT_SHA"]), Ref: github.String(ref), Sarif: github.String(sarifEncoded), CheckoutURI: github.String(ctx.Config.ProjectDir), ToolName: github.String("cid")}
		_, _, err = client.CodeScanning.UploadSarif(context.Background(), organization, repository, sarifAnalysis)
		if err != nil {
			return fmt.Errorf("failed to upload sarif to github code-scanning api: %s", err.Error())
		}
	}

	return nil
}

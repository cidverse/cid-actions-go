package cosignattach

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/cidverse/cid-actions-go/actions/cosign/cosigncommon"
	"github.com/cidverse/cid-actions-go/util/container"
	cidsdk "github.com/cidverse/cid-sdk-go"
)

type Action struct {
	Sdk cidsdk.SDKClient
}

type AttachConfig struct {
	CosignMode                   string `json:"cosign_mode" env:"COSIGN_MODE"`
	CosignKey                    string `json:"cosign_key" env:"COSIGN_KEY"`
	CosignPassword               string `json:"cosign_password" env:"COSIGN_PASSWORD"`
	CosignTransparencyLogDisable bool   `json:"cosign_tlog_disable" env:"COSIGN_TLOG_DISABLE"`
}

func (a Action) Metadata() cidsdk.ActionMetadata {
	return cidsdk.ActionMetadata{
		Name:        "cosign-container-sbom-attach",
		Description: `Cosign allows to attach SBOMs to a container image.`,
		Category:    "security",
		Scope:       cidsdk.ActionScopeProject,
		Rules: []cidsdk.ActionRule{
			{
				Type:       "cel",
				Expression: `MODULE_BUILD_SYSTEM == "container" && getMapValue(ENV, "COSIGN_KEY") != "" && getMapValue(ENV, "COSIGN_PASSWORD") != ""`,
			},
		},
		Access: cidsdk.ActionAccess{
			Environment: []cidsdk.ActionAccessEnv{
				{
					Name:        "COSIGN_MODE",
					Description: "The cosign mode, either 'KEYLESS' or 'PRIVATEKEY'.",
				},
				{
					Name:        "COSIGN_KEY",
					Description: "The cosign key, base64 encoded.",
				},
				{
					Name:        "COSIGN_PASSWORD",
					Description: "The password for the cosign key.",
				},
				{
					Name:        "COSIGN_TLOG_DISABLE",
					Description: "Disable using the public rekor transparency log.",
				},
			},
			Executables: []cidsdk.ActionAccessExecutable{
				{
					Name: "cosign",
				},
			},
		},
	}
}

func (a Action) Execute() (err error) {
	cfg := AttachConfig{}
	ctx, err := a.Sdk.ModuleAction(&cfg)
	if err != nil {
		return err
	}

	// private key for signing
	certFile := cidsdk.JoinPath(ctx.Config.TempDir, "private.key")
	if cfg.CosignMode == cosigncommon.CosignModePrivateKey {
		data, err := base64.StdEncoding.DecodeString(cfg.CosignKey)
		if err != nil {
			return fmt.Errorf("failed to decode private key file: %s", err.Error())
		}
		err = a.Sdk.FileWrite(certFile, data)
		if err != nil {
			return fmt.Errorf("failed to write private key file: %s", err.Error())
		}
	}

	// target image reference
	imageRef, err := a.Sdk.ArtifactDownloadByteArray(cidsdk.ArtifactDownloadByteArrayRequest{
		ID: fmt.Sprintf("%s|oci-image|image.txt", ctx.Module.Slug),
	})
	if err != nil {
		return fmt.Errorf("failed to parse image reference from %s", err.Error())
	}

	// digest
	digest, err := a.Sdk.ArtifactDownloadByteArray(cidsdk.ArtifactDownloadByteArrayRequest{
		ID: fmt.Sprintf("%s|oci-image|digest.txt", ctx.Module.Slug),
	})
	if err != nil {
		return fmt.Errorf("failed to read image digest: %s", err.Error())
	}

	// get manifests
	artifacts, err := a.Sdk.ArtifactList(cidsdk.ArtifactListRequest{Query: fmt.Sprintf(`module == "%s" && artifact_type == "oci-image" && format == "manifest"`, ctx.Module.Slug)})
	if err != nil {
		return fmt.Errorf("failed to query manifest artifact list: %s", err.Error())
	}
	if len(*artifacts) > 0 {
		_ = a.Sdk.Log(cidsdk.LogMessageRequest{Level: "warn", Message: "attachments for manifests are not supported yet"})
		return nil
	}

	// query reports
	reportList, err := a.Sdk.ArtifactList(cidsdk.ArtifactListRequest{Query: fmt.Sprintf(`module == "%s" && artifact_type == "report"`, ctx.Module.Slug)})
	if err != nil {
		return fmt.Errorf("failed to query sbom artifact list: %s", err.Error())
	}

	// attach reports
	for _, report := range *reportList {
		if !shouldAttachReport(report.Format, report.FormatVersion) {
			continue
		}

		// log
		_ = a.Sdk.Log(cidsdk.LogMessageRequest{
			Level:   "info",
			Message: "attaching file to container image",
			Context: map[string]interface{}{
				"file":            report.Name,
				"format":          report.Format,
				"format_version":  report.FormatVersion,
				"container-image": imageRef,
			},
		})

		// download file
		targetFile := cidsdk.JoinPath(ctx.Config.TempDir, report.Name)
		err := a.Sdk.ArtifactDownload(cidsdk.ArtifactDownloadRequest{
			ID:         report.ID,
			TargetFile: targetFile,
		})
		if err != nil {
			return err
		}

		// attach
		opts := []string{
			"--type",
			formatVersionToAttestationType(report.Format, report.FormatVersion),
			"--predicate",
			targetFile,
		}
		if cfg.CosignTransparencyLogDisable {
			opts = append(opts, "--no-tlog-upload=true")
		}
		if cfg.CosignMode == "KEYLESS" {
			_, err = a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
				Command: fmt.Sprintf(`cosign attest %s %s@%s`, strings.Join(opts, " "), container.GetImageReferenceWithoutTag(string(imageRef)), digest),
				WorkDir: ctx.Module.ModuleDir,
				Env: map[string]string{
					"COSIGN_EXPERIMENTAL": "1",
				},
			})
			if err != nil {
				return err
			}
		} else if cfg.CosignMode == cosigncommon.CosignModePrivateKey {
			_, err = a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
				Command: fmt.Sprintf(`cosign attest --key "%s" %s %s@%s`, certFile, strings.Join(opts, " "), container.GetImageReferenceWithoutTag(string(imageRef)), digest),
				WorkDir: ctx.Module.ModuleDir,
				Env: map[string]string{
					"COSIGN_PASSWORD": cfg.CosignPassword,
				},
			})
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("COSIGN_MODE [%s] is not supported, choose either PRIVATEKEY or KEYLESS", cfg.CosignMode)
		}
	}

	return nil
}

func shouldAttachReport(format string, formatVersion string) bool {
	if format == "container-sbom" {
		if formatVersion == "spdx-json" || formatVersion == "syft-json" || formatVersion == "cyclonedx-json" {
			return true
		}
	} else if format == "container-slsaprovenance" {
		return true
	}

	return false
}

func formatVersionToAttestationType(format string, formatVersion string) string {
	if format == "container-slsaprovenance" {
		return formatVersion
	} else if format == "container-sbom" {
		if formatVersion == "spdx-json" {
			return "spdxjson"
		} else if formatVersion == "syft-json" {
			return "https://syft.dev/bom"
		}
	}

	return "custom"
}

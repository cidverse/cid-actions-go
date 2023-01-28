package cosign

import (
	"encoding/base64"
	"fmt"
	"path"
	"strings"

	"github.com/cidverse/cid-actions-go/util"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/containers/image/v5/manifest"
)

type AttachAction struct {
	Sdk cidsdk.SDKClient
}

type AttachConfig struct {
	CosignMode                   string `json:"cosign_mode" env:"COSIGN_MODE"`
	CosignKey                    string `json:"cosign_key" env:"COSIGN_KEY"`
	CosignPassword               string `json:"cosign_password" env:"COSIGN_PASSWORD"`
	CosignTransparencyLogDisable bool   `json:"cosign_tlog_disable" env:"COSIGN_TLOG_DISABLE"`
}

func (a AttachAction) Execute() (err error) {
	cfg := AttachConfig{}
	ctx, err := a.Sdk.ModuleAction(&cfg)
	if err != nil {
		return err
	}

	// target image reference
	imageRef, err := a.Sdk.ArtifactDownloadByteArray(cidsdk.ArtifactDownloadByteArrayRequest{
		Module: ctx.Module.Slug,
		Type:   "oci-image",
		Name:   "image.txt",
	})
	if err != nil {
		return fmt.Errorf("failed to parse image reference from %s", err.Error())
	}

	// get manifests
	artifacts, err := a.Sdk.ArtifactList(cidsdk.ArtifactListRequest{
		Module:       ctx.Module.Slug,
		ArtifactType: "oci-image",
		Format:       "manifest",
	})
	if err != nil {
		return fmt.Errorf("failed to query manifest artifact list: %s", err.Error())
	}

	// process manifests
	var manifests []util.ManifestLayer
	for _, artifact := range *artifacts {
		if artifact.FormatVersion == "v2s2" {
			// fetch
			manifestJSON, err := a.Sdk.ArtifactDownloadByteArray(cidsdk.ArtifactDownloadByteArrayRequest{
				Module: ctx.Module.Slug,
				Type:   "oci-image",
				Name:   artifact.Name,
			})
			if err != nil {
				return fmt.Errorf("failed to read manifest.json: %s", err.Error())
			}

			// parse
			mf, err := manifest.ListFromBlob(manifestJSON, "application/vnd.docker.distribution.manifest.list.v2+json")
			if err != nil {
				return fmt.Errorf("failed to parse manifest.json: %s", err.Error())
			}
			manifests = util.GetManifestLayersV2S2(mf.(*manifest.Schema2List))
		} else {
			return fmt.Errorf("manifest format [%s] is not supported", artifact.FormatVersion)
		}
	}

	// query sbom
	sbomList, err := a.Sdk.ArtifactList(cidsdk.ArtifactListRequest{
		Module:       ctx.Module.Slug,
		ArtifactType: "report",
		Format:       "container-sbom",
	})
	if err != nil {
		return fmt.Errorf("failed to query sbom artifact list: %s", err.Error())
	}

	// private key for signing
	certFile := path.Join(ctx.Config.TempDir, "private.key")
	if cfg.CosignMode == cosignModePrivateKey {
		data, err := base64.StdEncoding.DecodeString(cfg.CosignKey)
		if err != nil {
			return fmt.Errorf("failed to decode private key file: %s", err.Error())
		}
		err = a.Sdk.FileWrite(certFile, data)
		if err != nil {
			return fmt.Errorf("failed to write private key file: %s", err.Error())
		}
	}

	// attach SBOMs individually if available
	for _, mfImage := range manifests {
		for _, sbom := range *sbomList {
			if strings.HasPrefix(sbom.Name, mfImage.PlatformFilename+".") {
				// skip unsupported formats
				if sbom.FormatVersion != "spdx-json" && sbom.FormatVersion != "syft-json" && sbom.FormatVersion != "cyclonedx-json" {
					continue
				}

				_ = a.Sdk.Log(cidsdk.LogMessageRequest{
					Level:   "info",
					Message: "attaching sbom to container image",
					Context: map[string]interface{}{
						"sbom":            sbom.Name,
						"format":          sbom.FormatVersion,
						"container-image": imageRef,
					},
				})

				// download sbom
				targetFile := path.Join(ctx.Config.TempDir, sbom.Name)
				err := a.Sdk.ArtifactDownload(cidsdk.ArtifactDownloadRequest{
					Module:     sbom.Module,
					Type:       string(sbom.Type),
					Name:       sbom.Name,
					TargetFile: targetFile,
				})
				if err != nil {
					return err
				}

				// attach sbom
				opts := []string{
					"--type",
					formatVersionToAttestationType(sbom.FormatVersion),
					"--predicate",
					targetFile,
				}
				if cfg.CosignTransparencyLogDisable {
					opts = append(opts, "--no-tlog-upload=true")
				}
				if cfg.CosignMode == "KEYLESS" {
					_, err = a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
						Command: fmt.Sprintf(`cosign attest %s %s@%s`, strings.Join(opts, " "), getContainerImageReferenceWithoutTag(string(imageRef)), mfImage.Digest),
						WorkDir: ctx.Module.ModuleDir,
						Env: map[string]string{
							"COSIGN_EXPERIMENTAL": "1",
						},
					})
					if err != nil {
						return err
					}
				} else if cfg.CosignMode == cosignModePrivateKey {
					_, err = a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
						Command: fmt.Sprintf(`cosign attest --key "%s" %s %s@%s`, certFile, strings.Join(opts, " "), getContainerImageReferenceWithoutTag(string(imageRef)), mfImage.Digest),
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
		}
	}

	return nil
}

func formatVersionToAttestationType(input string) string {
	if input == "spdx-json" {
		return "spdxjson"
	} else if input == "syft-json" {
		return "https://syft.dev/bom"
	}

	return "custom"
}

func getContainerImageReferenceWithoutTag(input string) string {
	parts := strings.SplitN(input, ":", 2)
	return parts[0]
}

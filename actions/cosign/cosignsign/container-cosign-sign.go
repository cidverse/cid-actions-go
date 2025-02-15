package cosignsign

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

type SignConfig struct {
	CosignMode                   string `json:"cosign_mode" env:"COSIGN_MODE"`
	CosignKey                    string `json:"cosign_key" env:"COSIGN_KEY"`
	CosignPassword               string `json:"cosign_password" env:"COSIGN_PASSWORD"`
	CosignTransparencyLogDisable bool   `json:"cosign_tlog_disable" env:"COSIGN_TLOG_DISABLE"`
}

func (a Action) Metadata() cidsdk.ActionMetadata {
	return cidsdk.ActionMetadata{
		Name: "cosign-container-sign",
		Description: `Cosign is a tool that enables container image signing, verification, and storage in an OCI registry.
        It supports various signing methods, including keyless signing using a public good certificate authority and transparency log, hardware and KMS signing, and bring-your-own PKI.
        With Cosign, you can sign your container images to ensure their authenticity and integrity, and store them in a trusted registry.
        The signed images can then be verified by others to ensure that they have not been tampered with.`,
		Category: "security",
		Scope:    cidsdk.ActionScopeProject,
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
	cfg := SignConfig{}
	ctx, err := a.Sdk.ModuleAction(&cfg)
	if err != nil {
		return err
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

	opts := []string{
		"--recursive",
		fmt.Sprintf(`-a "repo=%s"`, ctx.Env["NCI_REPOSITORY_REMOTE"]),
		fmt.Sprintf(`-a "commit_hash=%s"`, ctx.Env["NCI_COMMIT_HASH"]),
	}
	if cfg.CosignTransparencyLogDisable == true {
		opts = append(opts, "--no-tlog-upload=true")
	}

	if cfg.CosignMode == cosigncommon.CosignModeKeyless {
		// sign container
		signResult, err := a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
			Command: fmt.Sprintf(`cosign sign %s %s@%s`, strings.Join(opts, " "), container.GetImageReferenceWithoutTag(string(imageRef)), digest),
			WorkDir: ctx.ProjectDir,
			Env: map[string]string{
				"COSIGN_EXPERIMENTAL": "1",
			},
		})
		if err != nil {
			return err
		} else if signResult.Code != 0 {
			return fmt.Errorf("cosign sign failed, exit code %d", signResult.Code)
		}
	} else if cfg.CosignMode == cosigncommon.CosignModePrivateKey {
		// private key
		certFile := cidsdk.JoinPath(ctx.Config.TempDir, "private.key")
		data, err := base64.StdEncoding.DecodeString(cfg.CosignKey)
		if err != nil {
			return fmt.Errorf("failed to decode private key file: %s", err.Error())
		}
		err = a.Sdk.FileWrite(certFile, data)
		if err != nil {
			return fmt.Errorf("failed to write private key file: %s", err.Error())
		}

		// sign container
		signResult, err := a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
			Command: fmt.Sprintf(`cosign sign --key "%s" %s %s@%s`, certFile, strings.Join(opts, " "), container.GetImageReferenceWithoutTag(string(imageRef)), digest),
			WorkDir: ctx.ProjectDir,
			Env: map[string]string{
				"COSIGN_PASSWORD": cfg.CosignPassword,
			},
		})
		if err != nil {
			return err
		} else if signResult.Code != 0 {
			return fmt.Errorf("cosign sign failed, exit code %d", signResult.Code)
		}
	} else {
		return fmt.Errorf("COSIGN_MODE [%s] is not supported, choose either PRIVATEKEY or KEYLESS", cfg.CosignMode)
	}

	return nil
}

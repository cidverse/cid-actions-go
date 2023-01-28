package cosign

import (
	"encoding/base64"
	"fmt"
	"path"
	"strings"

	cidsdk "github.com/cidverse/cid-sdk-go"
)

type SignAction struct {
	Sdk cidsdk.SDKClient
}

type SignConfig struct {
	CosignMode                   string `json:"cosign_mode" env:"COSIGN_MODE"`
	CosignKey                    string `json:"cosign_key" env:"COSIGN_KEY"`
	CosignPassword               string `json:"cosign_password" env:"COSIGN_PASSWORD"`
	CosignTransparencyLogDisable bool   `json:"cosign_tlog_disable" env:"COSIGN_TLOG_DISABLE"`
}

func (a SignAction) Execute() (err error) {
	cfg := SignConfig{}
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

	// fetch
	digest, err := a.Sdk.ArtifactDownloadByteArray(cidsdk.ArtifactDownloadByteArrayRequest{
		Module: ctx.Module.Slug,
		Type:   "oci-image",
		Name:   "digest.txt",
	})
	if err != nil {
		return fmt.Errorf("failed to read image digest: %s", err.Error())
	}

	opts := []string{
		"--recursive",
		fmt.Sprintf(`-a "repo=%s"`, ctx.Env["NCI_REPOSITORY_REMOTE"]),
		fmt.Sprintf(`-a "commit_hash=%s"`, ctx.Env["NCI_COMMIT_SHA"]),
	}
	if cfg.CosignTransparencyLogDisable == true {
		opts = append(opts, "--no-tlog-upload=true")
	}

	if cfg.CosignMode == cosignModeKeyless {
		// sign container
		_, err = a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
			Command: fmt.Sprintf(`cosign sign %s %s@%s`, strings.Join(opts, " "), getContainerImageReferenceWithoutTag(string(imageRef)), digest),
			WorkDir: ctx.ProjectDir,
			Env: map[string]string{
				"COSIGN_EXPERIMENTAL": "1",
			},
		})
		if err != nil {
			return err
		}
	} else if cfg.CosignMode == cosignModePrivateKey {
		// private key
		certFile := path.Join(ctx.Config.TempDir, "private.key")
		data, err := base64.StdEncoding.DecodeString(cfg.CosignKey)
		if err != nil {
			return fmt.Errorf("failed to decode private key file: %s", err.Error())
		}
		err = a.Sdk.FileWrite(certFile, data)
		if err != nil {
			return fmt.Errorf("failed to write private key file: %s", err.Error())
		}

		// sign container
		_, err = a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
			Command: fmt.Sprintf(`cosign sign --key "%s" %s %s@%s`, certFile, strings.Join(opts, " "), getContainerImageReferenceWithoutTag(string(imageRef)), digest),
			WorkDir: ctx.ProjectDir,
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

	return nil
}

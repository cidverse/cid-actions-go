package container

import (
	"encoding/base64"
	"fmt"
	"path"

	cidsdk "github.com/cidverse/cid-sdk-go"
)

type SignAction struct {
	Sdk cidsdk.SDKClient
}

type SignConfig struct {
	CosignMode     string `json:"cosign_mode" env:"COSIGN_MODE"`
	CosignKey      string `json:"cosign_key" env:"COSIGN_KEY"`
	CosignPassword string `json:"cosign_password" env:"COSIGN_PASSWORD"`
}

func (a SignAction) Execute() (err error) {
	cfg := SignConfig{}
	ctx, err := a.Sdk.ModuleAction(&cfg)
	if err != nil {
		return err
	}

	// target image reference
	imageRefFile := path.Join(ctx.Config.ArtifactDir, ctx.Module.Slug, "oci-image", "image.txt")
	imageRef, err := a.Sdk.FileRead(imageRefFile)
	if err != nil {
		return fmt.Errorf("failed to parse image reference from %s", err.Error())
	}

	if cfg.CosignMode == "KEYLESS" {
		// sign container
		_, err = a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
			Command: fmt.Sprintf(`cosign sign -a "repo=%s" -a "commit_hash=%s" %s`, ctx.Env["NCI_REPOSITORY_REMOTE"], ctx.Env["NCI_COMMIT_SHA"], imageRef),
			WorkDir: ctx.ProjectDir,
			Env: map[string]string{
				"COSIGN_EXPERIMENTAL": "1",
			},
		})
		if err != nil {
			return err
		}
	} else if cfg.CosignMode == "PRIVATEKEY" {
		// files
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
			Command: fmt.Sprintf(`cosign sign --key "%s" -a "repo=%s" -a "commit_hash=%s" %s`, certFile, ctx.Env["NCI_REPOSITORY_REMOTE"], ctx.Env["NCI_COMMIT_SHA"], imageRef),
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

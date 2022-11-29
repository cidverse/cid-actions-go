package techdocs

import (
	"fmt"
	"path"
	"strings"

	cidsdk "github.com/cidverse/cid-sdk-go"
)

type PublishAction struct {
	Sdk cidsdk.SDKClient
}

type PublishConfig struct {
	Entity           string `json:"entity"  env:"TECHDOCS_ENTITY"`
	Target           string `json:"target" env:"TECHDOCS_PUBLISH_TARGET"`
	S3Endpoint       string `json:"s3_endpoint"  env:"TECHDOCS_S3_ENDPOINT"`
	S3Bucket         string `json:"s3_bucket"  env:"TECHDOCS_S3_BUCKET"`
	S3AccessKey      string `json:"s3_accesskey"  env:"TECHDOCS_S3_ACCESS_KEY"`
	S3SecretKey      string `json:"s3_accesssecret"  env:"TECHDOCS_S3_SECRET_KEY"`
	S3ForcePathStyle bool   `json:"s3_force_path_style"  env:"TECHDOCS_S3_FORCE_PATH_STYLE"`
	Insecure         bool   `json:"insecure"  env:"TECHDOCS_INSECURE"`
}

func (a PublishAction) Execute() (err error) {
	cfg := PublishConfig{}
	ctx, err := a.Sdk.ModuleAction(&cfg)
	if err != nil {
		return err
	}

	if ctx.Module.BuildSystem != string(cidsdk.BuildSystemMkdocs) || ctx.Module.BuildSystemSyntax != string(cidsdk.MkdocsTechdocs) {
		return fmt.Errorf("not supported: %s/%s", ctx.Module.BuildSystem, ctx.Module.BuildSystemSyntax)
	}
	outputDir := path.Join(ctx.Config.ArtifactDir, ctx.Module.Slug, "html")

	publishEnv := make(map[string]string)
	var publishArgs []string
	publishArgs = append(publishArgs, fmt.Sprintf(`--entity %s`, cfg.Entity))
	publishArgs = append(publishArgs, fmt.Sprintf(`--directory %s`, outputDir))
	if cfg.Insecure {
		publishEnv["NODE_TLS_REJECT_UNAUTHORIZED"] = "0"
	}
	if cfg.Target == "awsS3" {
		// auth
		publishEnv["AWS_ACCESS_KEY_ID"] = cfg.S3AccessKey
		publishEnv["AWS_SECRET_ACCESS_KEY"] = cfg.S3SecretKey

		// args
		publishArgs = append(publishArgs, fmt.Sprintf(`--publisher-type %s`, cfg.Target))
		publishArgs = append(publishArgs, fmt.Sprintf(`--awsEndpoint %s`, cfg.S3Endpoint))
		publishArgs = append(publishArgs, fmt.Sprintf(`--storage-name %s`, cfg.S3Bucket))
		if cfg.S3ForcePathStyle {
			publishArgs = append(publishArgs, `--awsS3ForcePathStyle`)
		}
	}

	_, err = a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
		Command: `techdocs-cli publish ` + strings.Join(publishArgs, " "),
		WorkDir: ctx.ProjectDir,
		Env:     publishEnv,
	})
	if err != nil {
		return err
	}

	return nil
}

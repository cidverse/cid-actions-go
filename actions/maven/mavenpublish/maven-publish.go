package mavenpublish

import (
	"fmt"
	"strings"

	"github.com/cidverse/cid-actions-go/actions/gradle/gradlecommon"
	"github.com/cidverse/cid-actions-go/actions/maven/mavencommon"
	"github.com/cidverse/cid-actions-go/util"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/go-playground/validator/v10"
)

type Action struct {
	Sdk cidsdk.SDKClient
}

type Config struct {
	WrapperVerification     bool   `json:"wrapper_verification" env:"WRAPPER_VERIFICATION"`
	MavenVersion            string `json:"maven_version"        env:"MAVEN_VERSION"`
	MavenRepositoryUrl      string `json:"maven_repo_url"       env:"MAVEN_REPO_URL"`
	MavenRepositoryUsername string `json:"maven_repo_username"  env:"MAVEN_REPO_USERNAME"`
	MavenRepositoryPassword string `json:"maven_repo_password"  env:"MAVEN_REPO_PASSWORD"`
	GPGSignPrivateKey       string `json:"gpg_sign_private_key" env:"MAVEN_GPG_SIGN_PRIVATEKEY"`
	GPGSignPassword         string `json:"gpg_sign_password"    env:"MAVEN_GPG_SIGN_PASSWORD"`
	GPGSignKeyId            string `json:"gpg_sign_key_id"      env:"MAVEN_GPG_SIGN_KEYID"`
}

func (a Action) Metadata() cidsdk.ActionMetadata {
	return cidsdk.ActionMetadata{
		Name: "maven-publish",
		Description: `This action publishes maven artifacts.

        **Publication**

        Username and password are not required when publishing to GitHub Packages. (https://maven.pkg.github.com/<your_username>/<your_repository>)

        **Signing**

        To sign the artifacts, MAVEN_GPG_SIGN_PRIVATEKEY and MAVEN_GPG_SIGN_PASSWORD must be set.
        You can generate a private key using the gpg command line tools:

        gpg --full-generate-key
        // When asked for the key type, select RSA (sign only).
        // For the key size, select 4096.
        // For the expiration date, select 0 to make the key never expire.
        // The Key-ID will appear in the output, it will look something like this A1B2C3D4E5F6G7H8.
        // Use a strong password to protect the key, store it in the MAVEN_GPG_SIGN_PASSWORD secret.
        gpg --armor --export-secret-keys A1B2C3D4E5F6G7H8 | base64 -w0
        // Now store the secret key in the MAVEN_GPG_SIGN_PRIVATEKEY
        `,
		Category: "publish",
		Scope:    cidsdk.ActionScopeModule,
		Rules: []cidsdk.ActionRule{
			{
				Type:       "cel",
				Expression: `MODULE_BUILD_SYSTEM == "maven" && getMapValue(ENV, "MAVEN_REPO_URL") != ""`,
			},
		},
		Access: cidsdk.ActionAccess{
			Environment: []cidsdk.ActionAccessEnv{
				{
					Name:        "MAVEN_VERSION",
					Description: "Overwrites the version of the maven artifact to publish, defaults to the git tag or branch name.",
				},
				{
					Name:        "MAVEN_REPO_URL",
					Description: "The URL of the maven repository to publish to.",
					Required:    true,
				},
				{
					Name:        "MAVEN_REPO_USERNAME",
					Description: "The username to use for authentication with the maven repository.",
				},
				{
					Name:        "MAVEN_REPO_PASSWORD",
					Description: "The password to use for authentication with the maven repository.",
				},
				{
					Name:        "MAVEN_GPG_SIGN_PRIVATEKEY",
					Description: "The ASCII-armored private key (base64 encoded).",
				},
				{
					Name:        "MAVEN_GPG_SIGN_PASSWORD",
					Description: "The password for the private key.",
				},
				{
					Name:        "MAVEN_GPG_SIGN_KEYID",
					Description: "The GPG key ID, only required when using sub keys.",
				},
				{
					Name:        "GITHUB_ACTOR",
					Description: "The GitHub actor to use for pushing the artifacts to a maven repository on GitHub Packages (https://maven.pkg.github.com/).",
				},
				{
					Name:        "GITHUB_TOKEN",
					Description: "The GitHub token to use for pushing the artifacts to a maven repository on GitHub Packages (https://maven.pkg.github.com/).",
					Secret:      true,
				},
			},
			Executables: []cidsdk.ActionAccessExecutable{
				{
					Name: "java",
				},
			},
			Network: util.MergeActionAccessNetwork(gradlecommon.NetworkJvm, gradlecommon.NetworkGradle, gradlecommon.NetworkPublish),
		},
	}
}

func (a Action) GetConfig(d *cidsdk.ModuleActionData) (Config, error) {
	cfg := Config{}
	cidsdk.PopulateFromEnv(&cfg, d.Env)

	// version
	if cfg.MavenVersion == "" {
		cfg.MavenVersion = gradlecommon.GetVersion(d.Env["NCI_COMMIT_REF_TYPE"], d.Env["NCI_COMMIT_REF_RELEASE"], d.Env["NCI_COMMIT_HASH_SHORT"])
	}

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
	d, err := a.Sdk.ModuleActionDataV1()
	if err != nil {
		return err
	}

	// parse config
	cfg, err := a.GetConfig(d)
	if err != nil {
		return err
	}

	// github packages
	if strings.HasPrefix(cfg.MavenRepositoryUrl, "https://maven.pkg.github.com/") {
		if cfg.MavenRepositoryUsername == "" {
			cfg.MavenRepositoryUsername = d.Env["GITHUB_ACTOR"]
		}
		if cfg.MavenRepositoryPassword == "" {
			cfg.MavenRepositoryPassword = d.Env["GITHUB_TOKEN"]
		}
	}

	// wrapper
	mavenWrapper := cidsdk.JoinPath(d.Module.ModuleDir, "mvnw")
	if !a.Sdk.FileExists(mavenWrapper) {
		return fmt.Errorf("maven wrapper not found at %s", mavenWrapper)
	}

	// version
	cmdResult, err := a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
		Command: mavencommon.MavenWrapperCommand(fmt.Sprintf("versions:set -DnewVersion=%q", cfg.MavenVersion)),
		WorkDir: d.Module.ModuleDir,
	})
	if err != nil {
		return err
	} else if cmdResult.Code != 0 {
		return fmt.Errorf("maven build failed, exit code %d", cmdResult.Code)
	}

	// publish
	cmdResult, err = a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
		Command: mavencommon.MavenWrapperCommand(`deploy --batch-mode`),
		WorkDir: d.Module.ModuleDir,
		Env: map[string]string{
			"MAVEN_REPO_URL":      cfg.MavenRepositoryUrl,
			"MAVEN_REPO_USERNAME": cfg.MavenRepositoryUsername,
			"MAVEN_REPO_PASSWORD": cfg.MavenRepositoryPassword,
		},
	})
	if err != nil {
		return err
	} else if cmdResult.Code != 0 {
		return fmt.Errorf("maven build failed, exit code %d", cmdResult.Code)
	}

	return nil
}

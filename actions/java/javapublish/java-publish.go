package javapublish

import (
	"fmt"
	"strings"

	"github.com/cidverse/cid-actions-go/actions/java/javacommon"
	cidsdk "github.com/cidverse/cid-sdk-go"
)

type Action struct {
	Sdk cidsdk.SDKClient
}

type Config struct {
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
		Name: "java-publish",
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
				Expression: `MODULE_BUILD_SYSTEM == "gradle" && getMapValue(ENV, "MAVEN_REPO_URL") != ""`,
			},
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
				},
			},
			Executables: []cidsdk.ActionAccessExecutable{
				{
					Name: "java",
				},
			},
		},
	}
}

func (a Action) Execute() (err error) {
	cfg := Config{}
	ctx, err := a.Sdk.ModuleAction(&cfg)
	if err != nil {
		return err
	}

	// github packages
	if strings.HasPrefix(cfg.MavenRepositoryUrl, "https://maven.pkg.github.com/") {
		if cfg.MavenRepositoryUsername == "" {
			cfg.MavenRepositoryUsername = ctx.Env["GITHUB_ACTOR"]
		}
		if cfg.MavenRepositoryPassword == "" {
			cfg.MavenRepositoryPassword = ctx.Env["GITHUB_TOKEN"]
		}
	}

	// version
	if cfg.MavenVersion == "" {
		cfg.MavenVersion = javacommon.GetVersion(ctx.Env["NCI_COMMIT_REF_TYPE"], ctx.Env["NCI_COMMIT_REF_RELEASE"], ctx.Env["NCI_COMMIT_HASH_SHORT"])
	}

	// publish
	if ctx.Module.BuildSystem == string(cidsdk.BuildSystemGradle) {
		// verify gradle wrapper
		err = javacommon.VerifyGradleWrapper(ctx.Module.ModuleDir)
		if err != nil {
			return err
		}

		gradleWrapper := cidsdk.JoinPath(ctx.Module.ModuleDir, "gradlew")
		if !a.Sdk.FileExists(gradleWrapper) {
			return fmt.Errorf("gradle wrapper not found at %s", gradleWrapper)
		}

		// TODO: run "gradle tasks --all" and check if the "publish" task is available?
		publishEnv := make(map[string]string)
		if cfg.GPGSignKeyId != "" {
			publishEnv["ORG_GRADLE_PROJECT_signingKeyId"] = cfg.GPGSignKeyId
		}
		if cfg.GPGSignPrivateKey != "" {
			publishEnv["ORG_GRADLE_PROJECT_signingKey"] = cfg.GPGSignPrivateKey
		}
		if cfg.GPGSignPassword != "" {
			publishEnv["ORG_GRADLE_PROJECT_signingPassword"] = cfg.GPGSignPassword
		}
		publishEnv["MAVEN_REPO_URL"] = cfg.MavenRepositoryUrl
		publishEnv["MAVEN_REPO_USERNAME"] = cfg.MavenRepositoryUsername
		publishEnv["MAVEN_REPO_PASSWORD"] = cfg.MavenRepositoryPassword

		// args
		publishArgs := []string{
			fmt.Sprintf(`-Pversion=%q`, cfg.MavenVersion),
			`publish`,
			`--no-daemon`,
			`--warning-mode=all`,
			`--console=plain`,
			`--stacktrace`,
		}
		publishResult, err := a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
			Command: fmt.Sprintf("java-exec %s %s", gradleWrapper, strings.Join(publishArgs, " ")),
			WorkDir: ctx.Module.ModuleDir,
			Env:     publishEnv,
		})
		if err != nil {
			return err
		} else if publishResult.Code != 0 {
			return fmt.Errorf("gradle publish failed, exit code %d", publishResult.Code)
		}
	} else if ctx.Module.BuildSystem == string(cidsdk.BuildSystemMaven) {
		mavenWrapper := cidsdk.JoinPath(ctx.Module.ModuleDir, "mvnw")
		if !a.Sdk.FileExists(mavenWrapper) {
			return fmt.Errorf("maven wrapper not found at %s", mavenWrapper)
		}

		buildArgs := []string{
			`deploy`,
			`--batch-mode`,
		}
		buildResult, err := a.Sdk.ExecuteCommand(cidsdk.ExecuteCommandRequest{
			Command: fmt.Sprintf("java-exec %s %s", mavenWrapper, strings.Join(buildArgs, " ")),
			WorkDir: ctx.Module.ModuleDir,
		})
		if err != nil {
			return err
		} else if buildResult.Code != 0 {
			return fmt.Errorf("maven publish failed, exit code %d", buildResult.Code)
		}
	}

	return nil
}

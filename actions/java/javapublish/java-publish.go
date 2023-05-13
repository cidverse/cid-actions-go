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
		cfg.MavenVersion = javacommon.GetVersion(ctx.Env["NCI_COMMIT_REF_TYPE"], ctx.Env["NCI_COMMIT_REF_RELEASE"], ctx.Env["NCI_COMMIT_SHA_SHORT"])
	}

	// publish
	if ctx.Module.BuildSystem == string(cidsdk.BuildSystemGradle) {
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

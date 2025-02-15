package javagradlewrapperscan

import (
	"fmt"

	cidsdk "github.com/cidverse/cid-sdk-go"
)

type Action struct {
	Sdk cidsdk.SDKClient
}

type Config struct {
}

func (a Action) Metadata() cidsdk.ActionMetadata {
	return cidsdk.ActionMetadata{
		Name: "java-gradle-wrapper-scan",
		Description: `Verifies the integrity of the gradle-wrapper.

        - queries information from https://services.gradle.org/versions/all
        - verifies the distribution hash in the gradle-wrapper.properties file against checksumUrl
        - verifies the gradle-wrapper.jar hash against wrapperChecksumUrl`,
		Category: "sast",
		Scope:    cidsdk.ActionScopeModule,
		Rules: []cidsdk.ActionRule{
			{
				Type:       "cel",
				Expression: `MODULE_BUILD_SYSTEM == "gradle"`,
			},
		},
		Access: cidsdk.ActionAccess{
			Environment: []cidsdk.ActionAccessEnv{},
			Executables: []cidsdk.ActionAccessExecutable{},
		},
	}
}

func (a Action) Execute() error {
	cfg := Config{}
	ctx, err := a.Sdk.ModuleAction(&cfg)
	if err != nil {
		return err
	}

	// check for the gradle wrapper
	propertiesFile := cidsdk.JoinPath(ctx.Module.ModuleDir, "gradle", "wrapper", "gradle-wrapper.properties")
	if a.Sdk.FileExists(propertiesFile) {
		props, err := ParseGradleWrapperProperties(propertiesFile)
		if err != nil {
			return fmt.Errorf("failed to parse gradle-wrapper.properties file: %w", err)
		}

		// find release
		version := ParseVersionInDistributionURL(props["distributionUrl"])
		if version == "" {
			return fmt.Errorf("failed to parse gradle version from distributionUrl: %s", props["distributionUrl"])
		}
		release, err := FindGradleRelease(version, true)
		if err != nil {
			return fmt.Errorf("failed to find gradle release for version %s: %w", version, err)
		}

		// distribution checksum
		if release.Checksum != props["distributionSha256Sum"] {
			return fmt.Errorf("distributionSha256Sum does not match expected value: %s != %s", release.Checksum, props["distributionSha256Sum"])
		}

		// verify checksums
		wrapperHash, err := hashFileSHA256(cidsdk.JoinPath(ctx.Module.ModuleDir, "gradle", "wrapper", "gradle-wrapper.jar"))
		if err != nil {
			return fmt.Errorf("failed to hash gradle/wrapper/gradle-wrapper.jar: %w", err)
		}
		if wrapperHash != release.WrapperChecksum {
			return fmt.Errorf("gradle/wrapper/gradle-wrapper.jar checksum does not match expected value: %s != %s", wrapperHash, release.WrapperChecksum)
		}
	}

	return nil
}

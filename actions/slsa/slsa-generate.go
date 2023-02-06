package slsa

import (
	"fmt"

	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/common"
	slsa "github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/v0.2"
)

type GenerateAction struct {
	Sdk cidsdk.SDKClient
}

type GenerateConfig struct {
}

func (a GenerateAction) Execute() (err error) {
	cfg := GenerateConfig{}
	ctx, err := a.Sdk.ModuleAction(&cfg)
	if err != nil {
		return err
	}

	prov := slsa.ProvenancePredicate{}

	// builder
	prov.BuildType = fmt.Sprintf("https://github.com/cidverse/cid@%s", "0.0.0")
	prov.Builder = common.ProvenanceBuilder{
		ID: fmt.Sprintf("https://github.com/cidverse/cid@%s", "0.0.0"),
	}

	// build config (stages and steps)
	prov.BuildConfig = []string{"test"}

	// invocation
	prov.Invocation = slsa.ProvenanceInvocation{
		ConfigSource: slsa.ConfigSource{
			URI: fmt.Sprintf("%s+%s@%s", ctx.Env["NCI_REPOSITORY_KIND"], ctx.Env["NCI_REPOSITORY_REMOTE"], ctx.Env["NCI_COMMIT_REF_NAME"]),
			Digest: common.DigestSet{
				"sha1": ctx.Env["NCI_COMMIT_SHA"],
			},
			// TODO: current workflow + current stage
			EntryPoint: "workflow=main,stage=build",
		},

		// Non user-controllable environment vars needed to reproduce the build.
		Environment: map[string]interface{}{
			"id":      ctx.Env["NCI_WORKER_ID"],
			"name":    ctx.Env["NCI_WORKER_NAME"],
			"arch":    ctx.Env["NCI_WORKER_ARCH"],
			"version": ctx.Env["NCI_WORKER_VERSION"],
		},

		// Parameters coming from the trigger event.
		Parameters: ctx.Env,
	}

	// metadata
	prov.Metadata = &slsa.ProvenanceMetadata{
		BuildInvocationID: "build-id",
		BuildStartedOn:    nil,
		BuildFinishedOn:   nil,
		Completeness: slsa.ProvenanceComplete{
			Parameters:  true,
			Environment: true,
			Materials:   true,
		},
		Reproducible: false,
	}

	// materials
	prov.Materials = []common.ProvenanceMaterial{
		{
			URI: "git+https://github.com/curl/curl-docker@master",
			Digest: common.DigestSet{
				"sha1": "d6525c840a62b398424a78d792f457477135d0cf",
			},
		},
		{
			URI: "github_hosted_vm:ubuntu-18.04:20210123.1",
		},
	}

	// store slsa provenance

	return nil
}

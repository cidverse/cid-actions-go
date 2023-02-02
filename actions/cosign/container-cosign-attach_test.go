package cosign

import (
	"testing"

	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/cidverse/cid-sdk-go/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCosignAttachManifest(t *testing.T) {
	sdk := mocks.NewSDKClient(t)
	sdk.On("ModuleAction", mock.Anything).Return(DockerfileTestData(false), nil).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*AttachConfig)
		arg.CosignMode = "KEYLESS"
	})
	sdk.On("ArtifactDownloadByteArray", cidsdk.ArtifactDownloadByteArrayRequest{
		Module: "my-project",
		Type:   "oci-image",
		Name:   "image.txt",
	}).Return([]byte(`docker.io/hello-world`), nil)
	sdk.On("ArtifactDownloadByteArray", cidsdk.ArtifactDownloadByteArrayRequest{
		Module: "my-project",
		Type:   "oci-image",
		Name:   "digest.txt",
	}).Return([]byte(`sha256:c38b49430bfe198766f03d135e58af0803588f89a26759d0c90d00f3a2aafde0`), nil)
	sdk.On("ArtifactList", cidsdk.ArtifactListRequest{Module: "my-project", ArtifactType: "oci-image", Format: "manifest"}).Return(&[]cidsdk.ActionArtifact{
		{
			BuildID:       "0",
			JobID:         "0",
			Module:        "my-module",
			Name:          "manifest.json",
			Type:          "oci-image",
			Format:        "manifest",
			FormatVersion: "oci",
		},
	}, nil)

	action := AttachAction{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}

func TestCosignAttach(t *testing.T) {
	sdk := mocks.NewSDKClient(t)
	sdk.On("ModuleAction", mock.Anything).Return(DockerfileTestData(false), nil).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*AttachConfig)
		arg.CosignMode = "KEYLESS"
	})
	sdk.On("ArtifactDownloadByteArray", cidsdk.ArtifactDownloadByteArrayRequest{
		Module: "my-project",
		Type:   "oci-image",
		Name:   "image.txt",
	}).Return([]byte(`docker.io/hello-world`), nil)
	sdk.On("ArtifactDownloadByteArray", cidsdk.ArtifactDownloadByteArrayRequest{
		Module: "my-project",
		Type:   "oci-image",
		Name:   "digest.txt",
	}).Return([]byte(`sha256:c38b49430bfe198766f03d135e58af0803588f89a26759d0c90d00f3a2aafde0`), nil)
	sdk.On("ArtifactList", cidsdk.ArtifactListRequest{Module: "my-project", ArtifactType: "oci-image", Format: "manifest"}).Return(&[]cidsdk.ActionArtifact{}, nil)
	sdk.On("ArtifactList", cidsdk.ArtifactListRequest{Module: "my-project", ArtifactType: "report"}).Return(&[]cidsdk.ActionArtifact{
		{
			BuildID:       "0",
			JobID:         "0",
			Module:        "my-module",
			Name:          "linux_amd64.syft.json",
			Type:          "report",
			Format:        "container-sbom",
			FormatVersion: "syft-json",
		},
		{
			BuildID:       "0",
			JobID:         "0",
			Module:        "my-module",
			Name:          "linux_amd64.spdx.json",
			Type:          "report",
			Format:        "container-sbom",
			FormatVersion: "spdx-json",
		},
		{
			BuildID:       "0",
			JobID:         "0",
			Module:        "my-module",
			Name:          "slsaprovenance.json",
			Type:          "report",
			Format:        "container-slsaprovenance",
			FormatVersion: "https://slsa.dev/provenance/v0.2",
		},
	}, nil)

	// syft
	sdk.On("ArtifactDownload", cidsdk.ArtifactDownloadRequest{
		Module:     "my-module",
		Type:       "report",
		Name:       "linux_amd64.syft.json",
		TargetFile: "/my-project/.tmp/linux_amd64.syft.json",
	}).Return(nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: `cosign attest --type https://syft.dev/bom --predicate /my-project/.tmp/linux_amd64.syft.json docker.io/hello-world@sha256:c38b49430bfe198766f03d135e58af0803588f89a26759d0c90d00f3a2aafde0`,
		WorkDir: "/my-project",
		Env: map[string]string{
			"COSIGN_EXPERIMENTAL": "1",
		},
	}).Return(nil, nil)

	// spdx
	sdk.On("ArtifactDownload", cidsdk.ArtifactDownloadRequest{
		Module:     "my-module",
		Type:       "report",
		Name:       "linux_amd64.spdx.json",
		TargetFile: "/my-project/.tmp/linux_amd64.spdx.json",
	}).Return(nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: `cosign attest --type spdxjson --predicate /my-project/.tmp/linux_amd64.spdx.json docker.io/hello-world@sha256:c38b49430bfe198766f03d135e58af0803588f89a26759d0c90d00f3a2aafde0`,
		WorkDir: "/my-project",
		Env: map[string]string{
			"COSIGN_EXPERIMENTAL": "1",
		},
	}).Return(nil, nil)

	// slsaprovenance
	sdk.On("ArtifactDownload", cidsdk.ArtifactDownloadRequest{
		Module:     "my-module",
		Type:       "report",
		Name:       "slsaprovenance.json",
		TargetFile: "/my-project/.tmp/slsaprovenance.json",
	}).Return(nil)
	sdk.On("ExecuteCommand", cidsdk.ExecuteCommandRequest{
		Command: `cosign attest --type https://slsa.dev/provenance/v0.2 --predicate /my-project/.tmp/slsaprovenance.json docker.io/hello-world@sha256:c38b49430bfe198766f03d135e58af0803588f89a26759d0c90d00f3a2aafde0`,
		WorkDir: "/my-project",
		Env: map[string]string{
			"COSIGN_EXPERIMENTAL": "1",
		},
	}).Return(nil, nil)

	action := AttachAction{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}

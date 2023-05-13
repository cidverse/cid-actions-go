package changeloggenerate

import (
	"testing"
	"time"

	"github.com/cidverse/cid-actions-go/actions/api"
	"github.com/cidverse/cid-actions-go/pkg/changelog"
	"github.com/cidverse/cid-actions-go/pkg/core/test"
	cidsdk "github.com/cidverse/cid-sdk-go"
	"github.com/cidverse/cid-sdk-go/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var Templates = []string{
	"github.changelog",
}

var CommitPattern = []string{
	"(?P<type>[A-Za-z]+)((?:\\((?P<scope>[^()\\r\\n]*)\\)|\\()?(?P<breaking>!)?)(:\\s?(?P<subject>.*))?",
}

var TitleMaps = map[string]string{
	"build":    "Build System",
	"ci":       "CI",
	"docs":     "Documentation",
	"feat":     "Features",
	"fix":      "Bug Fixes",
	"perf":     "Performance",
	"refactor": "Refactor",
	"style":    "Style",
	"test":     "Test",
	"chore":    "Internal",
}

var NoteKeywords = []changelog.NoteKeyword{
	{
		Keyword: "NOTE",
		Title:   "Notes",
	},
	{
		Keyword: "BREAKING CHANGE",
		Title:   "Breaking Changes",
	},
}

func TestChangelogGenerateWithPreviousRelease(t *testing.T) {
	sdk := test.Setup(t)
	sdk.On("ProjectAction", mock.Anything).Return(api.GetProjectActionData(false), nil).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*changelog.Config)
		arg.Templates = Templates
		arg.CommitPattern = CommitPattern
		arg.TitleMaps = TitleMaps
		arg.NoteKeywords = NoteKeywords
	})
	sdk.On("VCSReleases", cidsdk.VCSReleasesRequest{}).Return(&[]cidsdk.VCSRelease{
		{
			Version: "1.2.0",
			Ref:     cidsdk.VCSTag{RefType: "tag", Value: "v1.2.0"},
		},
		{
			Version: "1.1.0",
			Ref:     cidsdk.VCSTag{RefType: "tag", Value: "v1.1.0"},
		},
		{
			Version: "1.0.0",
			Ref:     cidsdk.VCSTag{RefType: "tag", Value: "v1.0.0"},
		},
	}, nil)
	sdk.On("VCSCommits", cidsdk.VCSCommitsRequest{
		FromHash: "hash/abcdef123456",
		ToHash:   "tag/v1.1.0",
		Limit:    1000,
	}).Return(&[]cidsdk.VCSCommit{
		{
			HashShort:   "123456a",
			Hash:        "f7331a7bc3a0531cf8aa4c982d7fefefffcbe8bc",
			Message:     "feat: add cool new feature",
			Description: "",
			Author:      cidsdk.VCSAuthor{Name: "A Person", Email: "email@example.com"},
			Committer:   cidsdk.VCSAuthor{Name: "A Person", Email: "email@example.com"},
			Tags:        nil,
			AuthoredAt:  time.Now(),
			CommittedAt: time.Now(),
			Changes:     nil,
			Context:     nil,
		},
	}, nil)
	sdk.On("ArtifactUpload", cidsdk.ArtifactUploadRequest{
		File:    "github.changelog",
		Content: "## Features\n- add cool new feature\n\n",
		Type:    "changelog",
	}).Return(nil)

	action := Action{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}

func TestChangelogGenerateFirstRelease(t *testing.T) {
	sdk := mocks.NewSDKClient(t)
	sdk.On("ProjectAction", mock.Anything).Return(api.GetProjectActionData(false), nil).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*changelog.Config)
		arg.Templates = Templates
		arg.CommitPattern = CommitPattern
		arg.TitleMaps = TitleMaps
		arg.NoteKeywords = NoteKeywords
	})
	sdk.On("VCSReleases", cidsdk.VCSReleasesRequest{}).Return(&[]cidsdk.VCSRelease{}, nil)
	sdk.On("VCSCommits", cidsdk.VCSCommitsRequest{
		FromHash: "hash/abcdef123456",
		ToHash:   "",
		Limit:    1000,
	}).Return(&[]cidsdk.VCSCommit{
		{
			HashShort:   "123456a",
			Hash:        "f7331a7bc3a0531cf8aa4c982d7fefefffcbe8bc",
			Message:     "feat: add cool new feature",
			Description: "",
			Author:      cidsdk.VCSAuthor{Name: "A Person", Email: "email@example.com"},
			Committer:   cidsdk.VCSAuthor{Name: "A Person", Email: "email@example.com"},
			Tags:        nil,
			AuthoredAt:  time.Now(),
			CommittedAt: time.Now(),
			Changes:     nil,
			Context:     nil,
		},
	}, nil)
	sdk.On("ArtifactUpload", cidsdk.ArtifactUploadRequest{
		File:    "github.changelog",
		Content: "## Features\n- add cool new feature\n\n",
		Type:    "changelog",
	}).Return(nil)

	action := Action{Sdk: sdk}
	err := action.Execute()
	assert.NoError(t, err)
}

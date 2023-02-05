package githubapi

import (
	"testing"
)

func TestGetGithubOrganizationFromRemote(t *testing.T) {
	tests := []struct {
		gitRemoteURL string
		want         string
	}{
		{"https://github.com/cidverse/normalizeci.git", "cidverse"},
		{"https://github.com/org-name/repo-name.git", "org-name"},
		{"https://github.com/username/repo-name.git", "username"},
	}
	for _, test := range tests {
		if got := GetGithubOrganizationFromRemote(test.gitRemoteURL); got != test.want {
			t.Errorf("GetGithubOrganization(%q) = %q, want %q", test.gitRemoteURL, got, test.want)
		}
	}
}

func TestGetGithubRepositoryFromRemote(t *testing.T) {
	tests := []struct {
		gitRemoteURL string
		want         string
	}{
		{"https://github.com/cidverse/normalizeci.git", "normalizeci"},
		{"https://github.com/org-name/repo-name.git", "repo-name"},
		{"https://github.com/username/repo-name.git", "repo-name"},
	}
	for _, test := range tests {
		if got := GetGithubRepositoryFromRemote(test.gitRemoteURL); got != test.want {
			t.Errorf("GetGithubRepository(%q) = %q, want %q", test.gitRemoteURL, got, test.want)
		}
	}
}

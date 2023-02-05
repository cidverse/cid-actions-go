package githubapi

import (
	"strings"
)

func GetGithubOrganizationFromRemote(gitRemoteURL string) string {
	parts := strings.Split(gitRemoteURL, "/")
	return parts[len(parts)-2]
}

func GetGithubRepositoryFromRemote(gitRemoteURL string) string {
	parts := strings.Split(gitRemoteURL, "/")
	repoNameWithExtension := parts[len(parts)-1]
	repoName := strings.TrimSuffix(repoNameWithExtension, ".git")
	return repoName
}

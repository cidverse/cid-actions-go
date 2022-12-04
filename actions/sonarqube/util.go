package sonarqube

import (
	"github.com/rs/zerolog/log"
)

func prepareProject(server string, accessToken string, organization string, projectKey string, projectName string, projectDescription string, mainBranch string) {
	// query branches
	branchList, branchListErr := getDefaultBranch(server, accessToken, projectKey)
	if branchListErr != nil {
		// no access or project doesn't exist - create
		createErr := createProject(server, accessToken, organization, projectKey, projectName)
		if createErr != nil {
			log.Error().Err(createErr).Msg("failed to create sonarqube project")
		}

		// rename main branch
		renameErr := renameMainBranch(server, accessToken, projectKey, mainBranch)
		if renameErr != nil {
			log.Error().Err(renameErr).Msg("failed to rename sonarqube main branch")
		}

		return
	}

	currentMainBranch := ""
	for _, branch := range branchList.Branches {
		if branch.IsMain {
			currentMainBranch = branch.Name
		}
	}

	// rename main branch if needed
	if mainBranch != currentMainBranch {
		log.Info().Str("current-main-branch", currentMainBranch).Str("new-main-branch", mainBranch).Msg("changing sonarqube main branch")

		// delete possible conflicts
		deleteErr := deleteBranch(server, accessToken, projectKey, mainBranch)
		if deleteErr != nil {
			log.Debug().Err(deleteErr).Str("branch", mainBranch).Msg("failed to delete sonarqube branch")
		}

		// rename main branch
		renameErr := renameMainBranch(server, accessToken, projectKey, mainBranch)
		if renameErr != nil {
			log.Error().Err(renameErr).Msg("failed to rename sonarqube main branch")
		}
	}
}

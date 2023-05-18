package sonarqubecommon

import (
	"fmt"
)

func PrepareProject(server string, accessToken string, organization string, projectKey string, projectName string, projectDescription string, mainBranch string) error {
	// query branches
	branchList, branchListErr := GetDefaultBranch(server, accessToken, projectKey)
	if branchListErr != nil {
		// no access or project doesn't exist - create
		createErr := CreateProject(server, accessToken, organization, projectKey, projectName)
		if createErr != nil {
			return fmt.Errorf("failed to create sonarqube project: %s", createErr.Error())
		}

		// rename main branch
		renameErr := RenameMainBranch(server, accessToken, projectKey, mainBranch)
		if renameErr != nil {
			return fmt.Errorf("failed to rename sonarqube main branch: %s", renameErr.Error())
		}

		return nil
	}

	currentMainBranch := ""
	for _, branch := range branchList.Branches {
		if branch.IsMain {
			currentMainBranch = branch.Name
		}
	}

	// rename main branch if needed
	if mainBranch != currentMainBranch {
		// delete possible conflicts
		deleteErr := DeleteBranch(server, accessToken, projectKey, mainBranch)
		if deleteErr != nil {
			return fmt.Errorf("failed to delete branch %s: %s", mainBranch, deleteErr.Error())
		}

		// rename main branch
		renameErr := RenameMainBranch(server, accessToken, projectKey, mainBranch)
		if renameErr != nil {
			return fmt.Errorf("failed to rename main branch: %s", renameErr.Error())
		}
	}

	return nil
}

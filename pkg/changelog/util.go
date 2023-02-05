package changelog

import (
	"embed"
	"errors"
	"path/filepath"

	"github.com/cidverse/cidverseutils/pkg/filesystem"
)

// GetFileContent returns the file content from either the directory or the embedded filesystem in that order
func GetFileContent(folder string, fs embed.FS, file string) (string, error) {
	if filesystem.FileExists(filepath.Join(folder, file)) {
		content, err := filesystem.GetFileContent(filepath.Join(folder, file))
		if err != nil {
			return "", err
		}

		return content, nil
	}

	// look in internal fs
	content, contentErr := GetFileContentFromEmbedFS(fs, "templates/"+file)
	if contentErr == nil {
		return content, nil
	}

	return "", errors.New("can't find template file " + file)
}

func GetFileContentFromEmbedFS(fs embed.FS, file string) (string, error) {
	fileBytes, fileErr := fs.ReadFile(file)

	if fileErr != nil {
		return "", fileErr
	}

	return string(fileBytes), nil
}

func AddLinks(input string) string {
	return input
}

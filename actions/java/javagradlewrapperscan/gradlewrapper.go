package javagradlewrapperscan

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

func ParseGradleWrapperProperties(filePath string) (map[string]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	props := make(map[string]string)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 && line[0] != '#' {
			fields := strings.SplitN(line, "=", 2)
			if len(fields) == 2 {
				props[strings.TrimSpace(fields[0])] = strings.TrimSpace(fields[1])
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return props, nil
}

func ParseVersionInDistributionURL(url string) string {
	re := regexp.MustCompile(`gradle-(\d+\.\d+\.\d+)[^/]*$`)
	matches := re.FindStringSubmatch(url)
	if len(matches) < 2 {
		return ""
	}
	return matches[1]
}

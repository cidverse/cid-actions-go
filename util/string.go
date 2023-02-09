package util

func FirstNonEmpty(strings []string) string {
	for _, str := range strings {
		if str != "" {
			return str
		}
	}
	return ""
}

package docker

import "strings"

func trimContainerName(contNames []string) string {
	return strings.TrimPrefix(contNames[0], "/")
}

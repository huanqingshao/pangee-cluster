package common

import (
	"fmt"
)

func GetRepoUrl(source string) (string, error) {
	url := ""

	switch source {
	case "github":
		url = "https://raw.githubusercontent.com/Horan-Z/test/refs/heads"
	case "gitee":
		url = "https://gitee.com/Horan-Z/test/raw"
	default:
		return url, fmt.Errorf("unsupported source type: %s", source)
	}

	return url, nil
}

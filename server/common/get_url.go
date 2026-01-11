package common

import (
	"fmt"
)

func GetRepoUrl(source string) (string, error) {
	url := ""

	// TODO: 修改地址
	switch source {
	case "github":
		url = "https://raw.githubusercontent.com/Horan-Z/test"
	case "gitee":
		url = "https://gitee.com/Horan-Z/test/raw"
	default:
		return url, fmt.Errorf("unsupported source type: %s", source)
	}

	return url, nil
}

func GetDownloadUrl(source string) (string, error) {
	url := ""

	// TODO: 修改地址
	switch source {
	case "github":
		url = "https://github.com/Horan-Z/test/releases/download"
	case "gitee":
		url = "https://gitee.com/Horan-Z/test/releases/download"
	default:
		return url, fmt.Errorf("unsupported source type: %s", source)
	}

	return url, nil
}

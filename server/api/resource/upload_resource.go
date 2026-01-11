package resource

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/opencmit/pangee-cluster/api/command"
	"github.com/opencmit/pangee-cluster/common"
	"github.com/opencmit/pangee-cluster/constants"
)

func UploadResource(c *gin.Context) {
	// 限制上传大小50MB
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 50<<20)

	enableProxyOnDownload := c.PostForm("enableProxyOnDownload")
	httpProxy := c.PostForm("httpProxy")
	file, err := c.FormFile("file")
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to get file", err)
		return
	}

	if !strings.HasSuffix(strings.ToLower(file.Filename), ".zip") {
		common.HandleError(c, http.StatusInternalServerError, "please upload a zip file", err)
		return
	}

	// 创建临时目录保存上传的ZIP文件
	tempDir, err := os.MkdirTemp("", "resource-upload")
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to create temp directory", err)
		return
	}
	defer os.RemoveAll(tempDir)

	tempFilePath := filepath.Join(tempDir, "resource-pack.zip")
	err = c.SaveUploadedFile(file, tempFilePath)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to save file", err)
		return
	}

	// 解压ZIP文件
	unzipDir := filepath.Join(tempDir, strings.TrimSuffix(file.Filename, ".zip"))
	err = os.MkdirAll(unzipDir, 0755)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to create unzip directory", err)
		return
	}

	// 使用命令行解压
	unzipCmd := exec.Command("unzip", tempFilePath, "-d", unzipDir)
	if err := unzipCmd.Run(); err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to unzip file", err)
		return
	}
	defer os.RemoveAll(unzipDir)

	packageYamlPath := filepath.Join(unzipDir, "package.yaml")
	info, err := os.Stat(packageYamlPath)
	if err != nil || info.IsDir() {
		common.HandleError(c, http.StatusInternalServerError, "package.yaml not found in zip root directory", err)
		return
	}

	// 使用yq获取版本号
	yqCmd := exec.Command("yq", "eval", ".metadata.version", packageYamlPath)
	output, err := yqCmd.Output()
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to get version from package.yaml", err)
		return
	}

	version := strings.TrimSpace(string(output))
	if version == "" {
		common.HandleError(c, http.StatusInternalServerError, "version not found in package.yaml", nil)
		return
	}

	uploadDir := filepath.Join(constants.GET_DATA_RESOURCE_DIR(), version)
	os.MkdirAll(uploadDir, 0755)

	mvCmd := exec.Command("mv", tempFilePath, uploadDir)
	mvOutput, err := mvCmd.CombinedOutput()
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError,
			fmt.Sprintf("failed to move zip file: %v, output: %s", err, string(mvOutput)),
			err)
		return
	}

	postExec := func(status command.ExecuteExitStatus) (string, error) {
		success := status.Success
		var message string
		if success {
			message = "\033[32m[ " + "Kuboardspray resource package has been loaded successfully." + " ]\033[0m \n"
			message += "\033[32m[ " + "Kuboardspray 资源包已成功上传，请回到资源包页面查看。" + " ]\033[0m \n"
		} else {
			message = "\033[31m\033[01m\033[05m[" + "Failed to download pangeecluster resource package, please reivew the log and fix the issues." + "]\033[0m \n"
			message += "\033[31m\033[01m\033[05m[" + "下载资源包失败，请回顾日志，找到错误信息，并解决问题后，再次尝试。" + "]\033[0m \n"
		}
		return "\n" + message, nil
	}

	cmd := command.Execute{
		OwnerType: "resource",
		OwnerName: version,
		Cmd:       "./pull-resource-package.sh",
		Args: func(execute_dir string) []string {
			return []string{
				uploadDir,
				"null",
				"null",
				enableProxyOnDownload,
				httpProxy,
			}
		},
		Type:     "download",
		PostExec: postExec,
	}

	if err := cmd.Exec(); err != nil {
		common.HandleError(c, http.StatusInternalServerError, "Faild to InstallCluster. ", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data": gin.H{
			"filename": file.Filename,
			"size":     file.Size,
			"pid":      cmd.R_Pid,
			"version":  version, // 返回获取到的版本号
		},
	})
}

func moveDirContents(src, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := os.MkdirAll(dstPath, 0755); err != nil {
				return err
			}
			if err := moveDirContents(srcPath, dstPath); err != nil {
				return err
			}
			os.RemoveAll(srcPath)
		} else {
			if err := os.Rename(srcPath, dstPath); err != nil {
				return err
			}
		}
	}
	return nil
}

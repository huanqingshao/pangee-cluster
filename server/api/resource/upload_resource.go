package resource

import (
	"fmt"
	"io"
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

	enableProxy := c.PostForm("enableProxy")
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

	version := strings.TrimSuffix(file.Filename, ".zip")
	uploadDir := constants.GET_DATA_DIR() + "/temp"
	unzipDir := filepath.Join(constants.GET_DATA_RESOURCE_DIR(), version, "content")
	// filepath.Dir(unzipDir) 是 unzipDir 的上级目录
	parentDir := filepath.Dir(unzipDir)
	os.MkdirAll(uploadDir, 0755)
	os.MkdirAll(unzipDir, 0755)
	tempFilePath := filepath.Join(uploadDir, file.Filename)
	err = c.SaveUploadedFile(file, tempFilePath)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to save file", err)
		return
	}

	err = unzip(tempFilePath, unzipDir)
	if err != nil {
		os.Remove(tempFilePath)
		common.HandleError(c, http.StatusInternalServerError, "failed to unzip file", err)
		return
	}

	packageYamlSource := filepath.Join(unzipDir, "package.yaml")
	packageYamlDest := filepath.Join(parentDir, "package.yaml")

	// 复制文件
	if err := copyFile(packageYamlSource, packageYamlDest); err != nil {
		// 出错时清理所有生成的文件和目录
		os.RemoveAll(parentDir)
		common.HandleError(c, http.StatusInternalServerError, "failed to copy package.yaml", err)
		return
	}

	os.Remove(tempFilePath)

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
				parentDir,
				version,
				"null",
				enableProxy,
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
			"unzipTo":  unzipDir,
			"pid":      cmd.R_Pid,
		},
	})
}

// 解压zip文件到指定目录
func unzip(src, dest string) error {
	cmd := exec.Command("rm", "-rf", dest)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("unzip failed: %v, output: %s", err, string(output))
	}
	cmd = exec.Command("unzip", src, "-d", dest)
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("unzip failed: %v, output: %s", err, string(output))
	}
	return nil
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

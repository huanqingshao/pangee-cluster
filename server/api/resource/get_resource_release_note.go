package resource

import (
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opencmit/pangee-cluster/common"
	"github.com/opencmit/pangee-cluster/constants"
)

func GetLocalResourceReleaseNote(c *gin.Context) {
	var req GetResourceRequest
	c.ShouldBindUri(&req)

	releaseNoteContent, err := os.ReadFile(GET_RESOURCE_RELEASE_NOTE_PATH(req.Name))
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "cannot open file: "+GET_RESOURCE_RELEASE_NOTE_PATH(req.Name), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data": gin.H{
			"release_note": string(releaseNoteContent),
		},
	})
}

func GET_RESOURCE_RELEASE_NOTE_PATH(name string) string {
	return constants.GET_DATA_RESOURCE_DIR() + "/" + name + "/content/release.md"
}

func GetRemoteResourceReleaseNote(c *gin.Context) {
	var req GetResourceRequest
	c.ShouldBindUri(&req)

	source := c.DefaultQuery("source", "github")
	requestPath, err := common.GetRepoUrl(source)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "cannot get url from source: "+source, err)
		return
	}

	requestPath += "/" + req.Name + "/release.md"

	client := &http.Client{}
	request, err := http.NewRequest("GET", requestPath, nil)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "cannot create request", err)
		return
	}

	// 添加查询参数防止缓存
	q := request.URL.Query()
	q.Add("nocache", strconv.FormatInt(time.Now().UnixNano(), 10))
	request.URL.RawQuery = q.Encode()

	resp, err := client.Do(request)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "cannot get release.md", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应体内容
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to read response body", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data": gin.H{
			"release_note": string(bodyBytes),
		},
	})
}

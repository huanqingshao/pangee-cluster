package resource

import (
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opencmit/pangee-cluster/api/command"
	"github.com/opencmit/pangee-cluster/common"
	"github.com/opencmit/pangee-cluster/constants"
	"gopkg.in/yaml.v3"
)

type GetResourceRequest struct {
	Name string `uri:"name"`
}

func GetLocalResource(c *gin.Context) {
	var req GetResourceRequest
	c.ShouldBindUri(&req)

	packageContent, err := os.ReadFile(GET_RESOURCE_YAML_PATH(req.Name))
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "cannot open file: "+GET_RESOURCE_YAML_PATH(req.Name), err)
		return
	}

	res := gin.H{}
	yaml.Unmarshal(packageContent, res)

	history, err := command.ReadTaskHistory("resource", req.Name)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "cannot read task history", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data": gin.H{
			"package": res,
			"history": history,
		},
	})
}

func GET_RESOURCE_YAML_PATH(name string) string {
	return constants.GET_DATA_RESOURCE_DIR() + "/" + name + "/package.yaml"
}

func GetRemoteResource(c *gin.Context) {

	var req GetResourceRequest
	c.ShouldBindUri(&req)

	source := c.DefaultQuery("source", "github")
	requestPath, err := common.GetRepoUrl(source)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "cannot get url from source: "+source, err)
		return
	}

	requestPath += "/" + req.Name + "/package.yaml"

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
		common.HandleError(c, http.StatusInternalServerError, "cannot get package.yaml", err)
		return
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to read response body", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data": gin.H{
			"package": string(bodyBytes),
		},
	})
}

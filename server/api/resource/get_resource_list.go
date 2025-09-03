package resource

import (
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opencmit/pangee-cluster/common"
	"gopkg.in/yaml.v3"
)

func GetPackageList(c *gin.Context) {
	// 从请求参数获取baseURL，默认为GitHub
	source := c.DefaultQuery("source", "github")

	yamlData, err := fetchYamlData(source)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取YAML数据失败"})
		return
	}

	var result struct {
		Items []interface{} `yaml:"items"`
	}

	err = yaml.Unmarshal(yamlData, &result)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "解析YAML数据失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": result.Items})
}

func fetchYamlData(source string) ([]byte, error) {

	requestPath, err := common.GetRepoUrl(source)
	if err != nil {
		return nil, err
	}

	requestPath += "/list/package-list.yaml"

	client := &http.Client{}
	request, err := http.NewRequest("GET", requestPath, nil)
	if err != nil {
		return nil, err
	}

	// 添加查询参数防止缓存
	q := request.URL.Query()
	q.Add("nocache", strconv.FormatInt(time.Now().UnixNano(), 10))
	request.URL.RawQuery = q.Encode()

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

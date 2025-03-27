package filebrowser

import (
	"net/http"

	"os"
	"strings"

	"github.com/eip-work/kuboard-spray/common"
	"github.com/eip-work/kuboard-spray/constants"
	"github.com/gin-gonic/gin"
)

type FileRequest struct {
	FilePath string `form:"filePath" binding:"required"`
}

func Get(c *gin.Context) {
	var req FileRequest
	c.ShouldBindQuery(&req)

	path := constants.GET_DATA_DIR() + strings.Replace(req.FilePath, "data://", "", -1)

	content, err := os.ReadFile(path)
	if err != nil {
		common.HandleError(c, http.StatusNotFound, "cannot open file: "+path, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data": gin.H{
			"content": string(content),
		},
	})
}

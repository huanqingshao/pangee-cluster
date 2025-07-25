package resource

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/opencmit/pangee-cluster/common"
	"github.com/opencmit/pangee-cluster/constants"
)

func GetResourceReleaseNote(c *gin.Context) {
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

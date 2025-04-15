package operation_v2

import (
	"net/http"

	"github.com/eip-work/kuboard-spray/api/command"
	"github.com/eip-work/kuboard-spray/common"
	"github.com/gin-gonic/gin"
)

func CheckStepStatus(c *gin.Context) {
	var request command.CheckStepStatusRequest
	c.ShouldBindUri(&request)
	c.ShouldBindJSON(&request)

	stepStatus, err := command.CheckStepStatusExec(request)

	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, err.Error(), err)
	}

	c.JSON(http.StatusOK, common.KuboardSprayResponse{
		Code:    http.StatusOK,
		Message: "success",
		Data:    stepStatus,
	})

}

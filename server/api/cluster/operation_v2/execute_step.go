package operation_v2

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opencmit/pangee-cluster/api/cluster/cluster_common"
	"github.com/opencmit/pangee-cluster/api/command"
	"github.com/opencmit/pangee-cluster/common"
)

type ExecuteStepRequest struct {
	Cluster      string `uri:"cluster" binding:"required"`
	Operation    string `uri:"operation" binding:"required"`
	Step         string `uri:"step" binding:"required"`
	Fork         int    `json:"fork"`
	Verbose      string `json:"verbose"`
	ExcludeNodes string `json:"nodes_to_exclude"`
}

func ExecuteStep(c *gin.Context) {
	var req ExecuteStepRequest
	c.ShouldBindUri(&req)
	c.ShouldBindJSON(&req)

	inventory, _, err := updateResourcePackageVarsToInventory(req)

	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to process inventory", err)
		return
	}

	postExec := func(status command.ExecuteExitStatus) (string, error) {
		if status.Success {
			return "\n执行成功: " + req.Cluster + "/" + req.Operation + "/" + req.Step, nil
		}
		return "\n执行失败: " + req.Cluster + "/" + req.Operation + "/" + req.Step, nil
	}

	playbook := "operations/" + req.Operation + "/" + req.Step + "/playbook.yaml"

	pid := req.Operation + "/" + req.Step + "/" + time.Now().Format("2006-01-02_15-04-05.999")

	cmd := command.Execute{
		OwnerType: "cluster",
		OwnerName: req.Cluster,
		Cmd:       "ansible-playbook",
		Args: func(execute_dir string) []string {
			result := []string{"-i", execute_dir + "/inventory.yaml", playbook}
			result = appendCommonParams(result, req, false)
			return result
		},
		Dir:  cluster_common.ResourcePackageDirForInventory(inventory),
		Type: req.Operation,
		PreExec: func(execute_dir string) error {
			return common.SaveYamlFile(execute_dir+"/inventory.yaml", inventory)
		},
		PostExec: postExec,
		Pid:      pid,
	}

	if err := cmd.Exec(); err != nil {
		common.HandleError(c, http.StatusInternalServerError, "Faild to InstallCluster. ", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data": gin.H{
			"pid": cmd.R_Pid,
		},
	})
}

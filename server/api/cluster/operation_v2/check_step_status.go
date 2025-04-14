package operation_v2

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/eip-work/kuboard-spray/api/ansible_rpc"
	"github.com/eip-work/kuboard-spray/api/cluster/cluster_common"
	"github.com/eip-work/kuboard-spray/api/command"
	"github.com/eip-work/kuboard-spray/common"
	"github.com/eip-work/kuboard-spray/constants"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type CheckStepStatusResponse map[string](map[string]ansible_rpc.AnsibleResultNode)

type CheckStepStatusRequest struct {
	Cluster   string `uri:"cluster" binding:"required"`
	Operation string `uri:"operation" binding:"required"`
	Step      string `uri:"step" binding:"required"`
}

func CheckStepStatus(c *gin.Context) {
	var request CheckStepStatusRequest
	c.ShouldBindUri(&request)
	c.ShouldBindJSON(&request)

	startTime := time.Now()

	cluster, err := cluster_common.ClusterMetadataByName(request.Cluster)

	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to get cluster metadata", err)
		return
	}

	playbook := "operations/" + request.Operation + "/" + request.Step + "/status.yaml"

	if !common.PathExists(cluster.ResourcePackageDir + "/" + playbook) {
		c.JSON(http.StatusNotFound, common.KuboardSprayResponse{
			Code:    http.StatusNotFound,
			Message: playbook + " is not found.",
		})
		return
	}

	cmd := command.Run{
		Cmd: "ansible-playbook",
		Args: []string{
			cluster.ResourcePackageDir + "/" + playbook,
			"-i", cluster.InventoryPath,
			"--forks", "200",
			"--timeout", "3",
			"-e", "kuboardspray_ssh_args='-o ConnectionAttempts=1 -o UserKnownHostsFile=/dev/null -F /dev/null'",
		},
		Env:     []string{"ANSIBLE_CONFIG=./ansible.cfg", "ANSIBLE_CACHE_PLUGIN_CONNECTION=" + constants.GET_DATA_CLUSTER_DIR() + "/" + request.Cluster + "/fact"},
		Timeout: 60,
		Dir:     "./ansible-rpc",
	}

	stdout, stderr, err := cmd.Run()
	duration := time.Now().UnixNano() - startTime.UnixNano()
	logrus.Trace("duration: ", duration/1000000)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to run", err)
		return
	}

	result := &ansible_rpc.AnsibleResult{}
	if err := json.Unmarshal(stdout, result); err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to Unmarshal result: ["+string(stdout)+"]", err)
		logrus.Trace("stdout: ", string(stdout), "\nstderr: ", string(stderr))
		logrus.Trace("duration: ", duration/1000000)
		return
	}

	stepStatus := CheckStepStatusResponse{}

	for _, task := range result.Plays[0].Tasks {
		for nodeName, node := range task.Hosts {
			if stepStatus[nodeName] == nil {
				stepStatus[nodeName] = make(map[string]ansible_rpc.AnsibleResultNode)
			}
			stepStatus[nodeName][task.Task.Name] = node
		}
	}

	c.JSON(http.StatusOK, common.KuboardSprayResponse{
		Code:    http.StatusOK,
		Message: "success",
		Data:    stepStatus,
	})

}

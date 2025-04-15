package command

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/eip-work/kuboard-spray/api/cluster/cluster_common"
	"github.com/eip-work/kuboard-spray/common"
	"github.com/eip-work/kuboard-spray/constants"
	"github.com/sirupsen/logrus"
)

// 只适用于 Cluster 的 Operation/Step

type CheckStepStatusResponse map[string](map[string]AnsibleResultNode)

type CheckStepStatusRequest struct {
	Cluster   string `uri:"cluster" binding:"required"`
	Operation string `uri:"operation" binding:"required"`
	Step      string `uri:"step" binding:"required"`
}

func CheckStepStatusExec(request CheckStepStatusRequest) (CheckStepStatusResponse, error) {

	startTime := time.Now()

	cluster, err := cluster_common.ClusterMetadataByName(request.Cluster)

	if err != nil {
		// common.HandleError(c, http.StatusInternalServerError, "failed to get cluster metadata", err)
		return nil, errors.New("failed to get cluster metadata")
	}

	playbook := "operations/" + request.Operation + "/" + request.Step + "/status.yaml"

	if !common.PathExists(cluster.ResourcePackageDir + "/" + playbook) {
		// c.JSON(http.StatusNotFound, common.KuboardSprayResponse{
		// 	Code:    http.StatusNotFound,
		// 	Message: playbook + " is not found.",
		// })
		return nil, errors.New(playbook + " is not found.")
	}

	cmd := Run{
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
		// common.HandleError(c, http.StatusInternalServerError, "failed to run", err)
		return nil, errors.New(err.Error())
	}

	result := &AnsibleResult{}
	if err := json.Unmarshal(stdout, result); err != nil {
		// common.HandleError(c, http.StatusInternalServerError, "failed to Unmarshal result: ["+string(stdout)+"]", err)
		logrus.Trace("stdout: ", string(stdout), "\nstderr: ", string(stderr))
		logrus.Trace("duration: ", duration/1000000)
		return nil, errors.New("failed to Unmarshal result: [" + string(stdout) + "]")
	}

	stepStatus := CheckStepStatusResponse{}

	for _, task := range result.Plays[0].Tasks {
		for nodeName, node := range task.Hosts {
			if stepStatus[nodeName] == nil {
				stepStatus[nodeName] = make(map[string]AnsibleResultNode)
			}
			stepStatus[nodeName][task.Task.Name] = node
		}
	}
	return stepStatus, nil
}

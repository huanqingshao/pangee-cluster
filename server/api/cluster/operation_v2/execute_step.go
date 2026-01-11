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

// 对于 inventory 中以 path 指定的 hosts 内所有属性 pangeecluster_resource_package 为 originalAction 的 host，
// 重置其 pangeecluster_resource_package 字段值为 none
func resetNodeAction(inventory map[string]interface{}, path string, originalAction string) {
	hosts := common.MapGet(inventory, path).(map[string]interface{})
	for _, hostItf := range hosts {
		if host, ok := hostItf.(map[string]interface{}); ok {
			if common.MapGet(host, "pangeecluster_node_action") == originalAction {
				common.MapSet(host, "pangeecluster_node_action", "none")
			}
		}
	}
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
			// FIXME 待验证

			tempClusterMetadata, err := cluster_common.ClusterMetadataByName(req.Cluster)
			if err != nil {
				return "解析数据出错", err
			}

			// logrus.Info("将要获取 data.operations")
			operations := common.MapGet(tempClusterMetadata.ResourcePackage, "data.operations").([]interface{})
			for _, operationItf := range operations { // 遍历所有的 operation
				// logrus.Info("遍历 operation", index)
				operation := operationItf.(map[string]interface{})
				nodeAction := common.MapGetString(operation, "pangeecluster_node_action")
				if common.MapGetString(operation, "name") == req.Operation && nodeAction != "" { // 如果当前 operation 的 pangeecluster_node_action 不为空
					// logrus.Info("匹配到 operation", nodeAction)
					steps := common.MapGet(operation, "steps").([]interface{})
					for index, stepItf := range steps { // 遍历所有的 steps
						// logrus.Info("匹配到 step", nodeAction)
						step := stepItf.(map[string]interface{})
						if common.MapGetString(step, "name") == req.Step && index == len(steps)-1 { // 如果当前 step 为最后一个，重置 pangeecluster_resource_package 字段
							// logrus.Info("是最后一个 step")
							resetNodeAction(inventory, "all.hosts", nodeAction)
							resetNodeAction(inventory, "all.children.target.children.etcd.hosts", nodeAction)
							resetNodeAction(inventory, "all.children.target.children.harbor.hosts", nodeAction)
							resetNodeAction(inventory, "all.children.target.children.k8s_cluster.children.kube_control_plane.hosts", nodeAction)
							resetNodeAction(inventory, "all.children.target.children.k8s_cluster.children.kube_node.hosts", nodeAction)

							// logrus.Info("保存到文件", tempClusterMetadata.InventoryPath)
							// 保存 inventory
							common.SaveYamlFile(tempClusterMetadata.InventoryPath, inventory)
						}
					}
					break
				}
			}

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

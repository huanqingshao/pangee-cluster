package operation_v2

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/eip-work/kuboard-spray/api/cluster/cluster_common"
	"github.com/eip-work/kuboard-spray/common"
	"github.com/eip-work/kuboard-spray/constants"
)

func updateResourcePackageVarsToInventory(req ExecuteStepRequest) (map[string]interface{}, map[string]interface{}, error) {

	clusterMetadata, err := cluster_common.ClusterMetadataByName(req.Cluster)
	if err != nil {
		return nil, nil, err
	}
	inventory := clusterMetadata.Inventory
	resourcePackage := clusterMetadata.ResourcePackage

	common.MapSet(inventory, "all.vars.kuboardspray_no_log", !(req.Verbose == "v" || req.Verbose == "vvv"))

	versionJsonPath := constants.GetKuboardSprayWebDir() + "/version.json"
	versionJson, err := os.ReadFile(versionJsonPath)
	if err != nil {
		return nil, nil, err
	}
	version := map[string]string{}
	json.Unmarshal(versionJson, &version)

	v := version["version"]
	v = strings.TrimSuffix(v, "-amd64")
	v = strings.TrimSuffix(v, "-arm64")

	common.MapSet(inventory, "all.vars.kuboardspray_version", v)

	common.PopulateKuboardSprayVars(inventory, "cluster", req.Cluster)

	common.MapSet(inventory, "all.vars.ansible_ssh_pipelining", false)

	return inventory, resourcePackage, nil
}

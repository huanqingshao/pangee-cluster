package fact

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/eip-work/kuboard-spray/api/command"
	"github.com/eip-work/kuboard-spray/constants"
)

func nodefact_cached(req GetNodeFactRequest) (*command.AnsibleResultNode, error) {

	factDir := constants.GET_DATA_DIR() + "/" + req.NodeOwnerType + "/" + req.NodeOwner + "/fact"
	factPath := factDir + "/" + req.Node + "_" + req.Ip + "_" + req.Port

	fact_bytes, err := ioutil.ReadFile(factPath)
	if err != nil {
		return nil, err
	}

	fact := &command.AnsibleResultNode{}
	if err := json.Unmarshal(fact_bytes, fact); err != nil {
		return nil, errors.New("Failed to Unmarshal result " + err.Error())
	}
	return fact, nil

}

package delivery

import (
	"encoding/json"
	"log"
	"net/http"
)

func (router *Router) GetAllNodesHandler(w http.ResponseWriter, r *http.Request) {
	nodesInfo, err := router.nodeManager.ListNodesInfo()
	if err != nil {
		RespondErr(w, http.StatusInternalServerError, err)
		return
	}

	resps := make([]NodeInfoResponse, 0, len(nodesInfo))
	for _, nodeInfo := range nodesInfo {
		resps = append(resps, NodeInfoResponse{
			NodeName:    nodeInfo.NodeName,
			XrayRunning: nodeInfo.IsXrayRunning,
		})
	}

	RespondSuccess(w, http.StatusOK, resps)
}

func (router *Router) AddNode(w http.ResponseWriter, r *http.Request) {

	var req AddNodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondErr(w, http.StatusInternalServerError, err)
		return
	}

	res, err := router.nodeManager.AddNode(r.Context(), req.NodeAddress)
	if err != nil {
		log.Printf("failed to add node, err: %s", err)
		RespondErr(w, http.StatusInternalServerError, err)
		return
	}

	RespondSuccess(w, http.StatusOK, &NodeInfoResponse{
		NodeName:    res.NodeName,
		XrayRunning: res.IsXrayRunning,
	})
}

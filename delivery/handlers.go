package delivery

import (
	"encoding/json"
	"errors"
	"gofronet-foundation/gofro-control/nodes"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Router struct {
	nodeManager *nodes.NodeManager
}

func NewRouter(nodeManager *nodes.NodeManager) *Router {
	return &Router{
		nodeManager: nodeManager,
	}
}

func (router *Router) Register(r chi.Router) {
	r.Route("/nodes", func(r chi.Router) {
		r.Get("/", router.GetAllNodesHandler)
		r.Post("/", router.AddNode)

		r.Route("/{node_name}", func(r chi.Router) {
			r.Get("/config", router.GetCurrentConfig)
			r.Put("/config", router.UpdateConfig)

			r.Post("/start", router.Start)
			r.Post("/stop", router.Stop)
			r.Post("/restart", router.Restart)
		})
	})
}

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
		RespondErr(w, http.StatusInternalServerError, err)
		return
	}

	RespondSuccess(w, http.StatusOK, &NodeInfoResponse{
		NodeName:    res.NodeName,
		XrayRunning: res.IsXrayRunning,
	})
}

func (router *Router) UpdateConfig(w http.ResponseWriter, r *http.Request) {

	nodeName := chi.URLParam(r, "node_name")
	if nodeName == "" {
		RespondErr(w, http.StatusBadRequest, errors.New("node_name required"))
		return
	}

	var req UpdateConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondErr(w, http.StatusInternalServerError, err)
		return
	}

	node, err := router.nodeManager.GetNodeByName(r.Context(), nodeName)
	if err != nil {
		RespondErr(w, http.StatusInternalServerError, err)
		return
	}

	if err := node.UpdateConfig(r.Context(), req.NewConfig); err != nil {
		RespondErr(w, http.StatusInternalServerError, err)
		return
	}

	RespondSuccess(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (router *Router) GetCurrentConfig(w http.ResponseWriter, r *http.Request) {

	nodeName := chi.URLParam(r, "node_name")
	if nodeName == "" {
		RespondErr(w, http.StatusBadRequest, errors.New("node_name required"))
		return
	}

	node, err := router.nodeManager.GetNodeByName(r.Context(), nodeName)
	if err != nil {
		RespondErr(w, http.StatusInternalServerError, err)
		return
	}

	currentConfig, err := node.GetCurrentConfig(r.Context())
	if err != nil {
		RespondErr(w, http.StatusInternalServerError, err)
		return
	}

	RespondSuccess(w, http.StatusOK, &GetCurrentConfigResponse{
		NodeName:      nodeName,
		CurrentConfig: currentConfig,
	})
}

func (router *Router) Start(w http.ResponseWriter, r *http.Request) {
	nodeName := chi.URLParam(r, "node_name")
	if nodeName == "" {
		RespondErr(w, http.StatusBadRequest, errors.New("node_name required"))
		return
	}

	node, err := router.nodeManager.GetNodeByName(r.Context(), nodeName)
	if err != nil {
		RespondErr(w, http.StatusInternalServerError, err)
		return
	}

	if err := node.Start(r.Context()); err != nil {
		RespondErr(w, http.StatusInternalServerError, err)
		return
	}

	RespondSuccess(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (router *Router) Stop(w http.ResponseWriter, r *http.Request) {
	nodeName := chi.URLParam(r, "node_name")
	if nodeName == "" {
		RespondErr(w, http.StatusBadRequest, errors.New("node_name required"))
		return
	}

	node, err := router.nodeManager.GetNodeByName(r.Context(), nodeName)
	if err != nil {
		RespondErr(w, http.StatusInternalServerError, err)
		return
	}

	if err := node.Stop(r.Context()); err != nil {
		RespondErr(w, http.StatusInternalServerError, err)
		return
	}

	RespondSuccess(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (router *Router) Restart(w http.ResponseWriter, r *http.Request) {
	nodeName := chi.URLParam(r, "node_name")
	if nodeName == "" {
		RespondErr(w, http.StatusBadRequest, errors.New("node_name required"))
		return
	}

	node, err := router.nodeManager.GetNodeByName(r.Context(), nodeName)
	if err != nil {
		RespondErr(w, http.StatusInternalServerError, err)
		return
	}

	if err := node.Restart(r.Context()); err != nil {
		RespondErr(w, http.StatusInternalServerError, err)
		return
	}

	RespondSuccess(w, http.StatusOK, map[string]string{"status": "ok"})
}

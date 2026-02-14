package delivery

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

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

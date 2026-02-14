package delivery

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

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

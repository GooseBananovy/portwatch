package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/goosebananovy/portwatch/internal/model/metrics/docker"
	"github.com/goosebananovy/portwatch/internal/model/metrics/network"
	"github.com/goosebananovy/portwatch/internal/model/metrics/procs"
	"github.com/goosebananovy/portwatch/internal/model/metrics/sys"
	"github.com/goosebananovy/portwatch/internal/provider"
)

type MetricsHandler struct {
	provider provider.Provider
}

func NewMetricsHandler(p provider.Provider) *MetricsHandler {
	return &MetricsHandler{
		provider: p,
	}
}

func (mh *MetricsHandler) All(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")

	var response struct {
		SystemStats sys.Stats     `json:"system_stats"`
		Network     network.Stats `json:"network"`
		Docker      docker.Cons   `json:"docker"`
		Procs       procs.Procs   `json:"procs"`
	}

	var err error

	response.SystemStats, err = mh.provider.Sys(ctx)
	if err != nil {
		log.Printf("failed to get system stats: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	response.Network, err = mh.provider.Network(ctx)
	if err != nil {
		log.Printf("failed to get network stats: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	response.Docker, err = mh.provider.Docker(ctx)
	if err != nil {
		log.Printf("failed to get docker stats: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	response.Procs, err = mh.provider.Procs(ctx)
	if err != nil {
		log.Printf("failed to get procs stats: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(response)
}

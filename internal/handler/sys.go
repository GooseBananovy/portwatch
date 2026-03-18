package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/goosebananovy/portwatch/internal/model/metrics/sys"
)

func (mh *MetricsHandler) Uptime(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")

	var response struct {
		UptimeSec sys.UptimeSec `json:"uptime"`
	}

	var err error
	response.UptimeSec, err = mh.provider.Uptime(ctx)
	if err != nil {
		log.Printf("failed to get uptime: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(response)
}

func (mh *MetricsHandler) Cpu(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")

	var response sys.Cpu
	var err error

	response, err = mh.provider.Cpu(ctx)
	if err != nil {
		log.Printf("failed to get cpu stats: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(response)
}

func (mh *MetricsHandler) Ram(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")

	var response sys.Ram
	var err error

	response, err = mh.provider.Ram(ctx)
	if err != nil {
		log.Printf("failed to get ram stats: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(response)
}

func (mh *MetricsHandler) Disk(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")

	var response sys.Disk
	var err error

	response, err = mh.provider.Disk(ctx)
	if err != nil {
		log.Printf("failed to get disk stats: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(response)
}

func (mh *MetricsHandler) Sys(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")

	var response sys.Stats
	var err error

	response, err = mh.provider.Sys(ctx)
	if err != nil {
		log.Printf("failed to get system stats: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(response)
}

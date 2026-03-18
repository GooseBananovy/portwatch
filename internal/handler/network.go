package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/goosebananovy/portwatch/internal/model/metrics/network"
)

func (mh *MetricsHandler) Network(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")

	var response network.Stats
	var err error

	response, err = mh.provider.Network(ctx)
	if err != nil {
		log.Printf("failed to get network stats: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(response)
}

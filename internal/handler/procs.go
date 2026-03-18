package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/goosebananovy/portwatch/internal/model/metrics/procs"
)

func (mh *MetricsHandler) Procs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")

	var response procs.Procs
	var err error

	response, err = mh.provider.Procs(ctx)
	if err != nil {
		log.Printf("failed to get procs stats: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(response)
}

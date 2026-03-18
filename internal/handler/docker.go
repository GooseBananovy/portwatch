package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/goosebananovy/portwatch/internal/model/metrics/docker"
)

func (mh *MetricsHandler) Docker(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")

	var response docker.Cons
	var err error

	response, err = mh.provider.Docker(ctx)
	if err != nil {
		log.Printf("failed to get docker stats: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(response)
}

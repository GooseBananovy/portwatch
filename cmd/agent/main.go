package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/goosebananovy/portwatch/internal/handler"
	"github.com/goosebananovy/portwatch/internal/provider/linux"
)

func main() {
	lp := linux.NewLinuxProvider()
	mh := handler.NewMetricsHandler(lp)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /portwatch/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		type healthResponse struct {
			Status string `json:"status"`
		}

		response := healthResponse{
			Status: "ok",
		}
		json.NewEncoder(w).Encode(response)
	})

	mux.HandleFunc("GET /portwatch/uptime", mh.Uptime)
	mux.HandleFunc("GET /portwatch/cpu", mh.Cpu)
	mux.HandleFunc("GET /portwatch/ram", mh.Ram)
	mux.HandleFunc("GET /portwatch/disk", mh.Disk)
	mux.HandleFunc("GET /portwatch/sys", mh.Sys)

	mux.HandleFunc("GET /portwatch/procs", mh.Procs)

	mux.HandleFunc("GET /portwatch/network", mh.Network)

	mux.HandleFunc("GET /portwatch/docker", mh.Docker)

	mux.HandleFunc("GET /portwatch/all", mh.All)

	if err := http.ListenAndServe(":9090", mux); err != nil {
		log.Fatal(err)
	}
}

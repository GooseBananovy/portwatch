package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	host := os.Getenv("PORTWATCH_HOST")

	if len(os.Args) < 2 {
		log.Fatal("Bad usage")
	}

	command := os.Args[1]
	var response *http.Response
	var err error

	switch command {
	case "uptime":
		response, err = http.Get(host + "/portwatch/uptime")
	case "health":
		response, err = http.Get(host + "/portwatch/health")
	case "cpu":
		response, err = http.Get(host + "/portwatch/cpu")
	case "ram":
		response, err = http.Get(host + "/portwatch/ram")
	case "disk":
		response, err = http.Get(host + "/portwatch/disk")
	case "sys":
		response, err = http.Get(host + "/portwatch/sys")
	case "procs":
		response, err = http.Get(host + "/portwatch/procs")
	case "network":
		response, err = http.Get(host + "/portwatch/network")
	case "docker":
		response, err = http.Get(host + "/portwatch/docker")
	case "all":
		response, err = http.Get(host + "/portwatch/all")
	default:
		log.Fatal("Unknown command: " + command)
	}

	if err != nil {
		log.Fatalf("Failed to make request: %v", err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Fatalf("Request failed with status: %s", response.Status)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	fmt.Println(string(body))
}

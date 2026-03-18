package linux

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/goosebananovy/portwatch/internal/model/metrics/docker"
)

type consInfoScreen struct {
	CpuStats struct {
		SystemCpuUsage uint64 `json:"system_cpu_usage"`
		CpuUsage       struct {
			TotalUsage uint64 `json:"total_usage"`
		} `json:"cpu_usage"`
	} `json:"cpu_stats"`

	MemoryStats struct {
		Usage uint64 `json:"usage"`
		Limit uint64 `json:"limit"`
	} `json:"memory_stats"`

	PreCpuStats struct {
		SystemCpuUsage uint64 `json:"system_cpu_usage"`
		CpuUsage       struct {
			TotalUsage uint64 `json:"total_usage"`
		} `json:"cpu_usage"`
	} `json:"precpu_stats"`
}

func (lp *LinuxProvider) Docker(ctx context.Context) (docker.Cons, error) {
	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
				return (&net.Dialer{}).DialContext(ctx, "unix", "/var/run/docker.sock")
			},
		},
	}

	resp, err := client.Get("http://localhost/containers/json")
	if err != nil {
		return nil, fmt.Errorf("failed to get containers json list: %w", err)
	}

	var consList []struct {
		Id     string   `json:"Id"`
		Names  []string `json:"Names"`
		Status string   `json:"Status"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&consList); err != nil {
		return nil, fmt.Errorf("failed to parse list containers json: %w", err)
	}
	resp.Body.Close()

	var result docker.Cons

	for i := range consList {
		id := consList[i].Id

		if resp, err = client.Get("http://localhost/containers/" + id + "/stats?stream=false"); err != nil {
			return nil, fmt.Errorf("failed to get info about container "+id+": %w", err)
		}

		var screen consInfoScreen

		if err = json.NewDecoder(resp.Body).Decode(&screen); err != nil {
			return nil, fmt.Errorf("failed to parse json info about container "+id+": %w", err)
		}
		resp.Body.Close()

		deltaConCpu := screen.CpuStats.CpuUsage.TotalUsage - screen.PreCpuStats.CpuUsage.TotalUsage
		deltaSysCpu := screen.CpuStats.SystemCpuUsage - screen.PreCpuStats.SystemCpuUsage

		result = append(result, docker.Container{
			Name:     consList[i].Names[0],
			Status:   consList[i].Status,
			Cpu:      float64(deltaConCpu) / float64(deltaSysCpu) * 100,
			RamUsed:  screen.MemoryStats.Usage,
			RamTotal: screen.MemoryStats.Limit,
		})
	}

	return result, nil
}

package provider

import (
	"context"

	"github.com/goosebananovy/portwatch/internal/model/metrics/docker"
	"github.com/goosebananovy/portwatch/internal/model/metrics/network"
	"github.com/goosebananovy/portwatch/internal/model/metrics/procs"
	"github.com/goosebananovy/portwatch/internal/model/metrics/sys"
)

type Provider interface {
	Docker(ctx context.Context) (docker.Cons, error)
	Network(ctx context.Context) (network.Stats, error)
	Procs(ctx context.Context) (procs.Procs, error)
	Cpu(ctx context.Context) (sys.Cpu, error)
	Ram(ctx context.Context) (sys.Ram, error)
	Disk(ctx context.Context) (sys.Disk, error)
	Uptime(ctx context.Context) (sys.UptimeSec, error)
	Sys(ctx context.Context) (sys.Stats, error)
}

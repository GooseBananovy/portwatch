package linux

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/sys/unix"

	"github.com/goosebananovy/portwatch/internal/model/metrics/sys"
)

func (lp *LinuxProvider) Uptime(ctx context.Context) (sys.UptimeSec, error) {
	pathUptime := lp.procBase + "/uptime"
	rawData, err := os.ReadFile(pathUptime)
	if err != nil {
		return 0, fmt.Errorf("failed to read %s file: %w", pathUptime, err)
	}

	uptime, err := strconv.ParseFloat(strings.Fields(string(rawData))[0], 64)
	if err != nil {
		return 0, fmt.Errorf("failed to convert %s content to float: %w", pathUptime, err)
	}

	if uptime <= 0 {
		return 0, errors.New("got non-positive uptime duration")
	}

	return sys.UptimeSec(uptime), nil
}

func (lp *LinuxProvider) Ram(ctx context.Context) (sys.Ram, error) {
	pathMeminfo := lp.procBase + "/meminfo"
	rawData, err := os.ReadFile(pathMeminfo)
	if err != nil {
		return sys.Ram{}, fmt.Errorf("failed to read %s file: %w", pathMeminfo, err)
	}

	var memTotal uint64
	var memAvailable uint64

	for line := range strings.SplitSeq((string(rawData)), "\n") {
		if strings.HasPrefix(line, "MemTotal:") {
			memTotal, err = strconv.ParseUint(strings.Fields(line)[1], 10, 64)

			if err != nil {
				return sys.Ram{}, fmt.Errorf("failed to convert %s content to uint: %w", pathMeminfo, err)
			}

		} else if strings.HasPrefix(line, "MemAvailable:") {
			memAvailable, err = strconv.ParseUint(strings.Fields(line)[1], 10, 64)

			if err != nil {
				return sys.Ram{}, fmt.Errorf("failed to convert %s content to uint: %w", pathMeminfo, err)
			}

		}
	}

	if memAvailable > memTotal {
		return sys.Ram{}, errors.New("available memory greater then total")
	}

	return sys.Ram{
		TotalBytes:     memTotal * 1024,
		AvailableBytes: memAvailable * 1024,
		UsedBytes:      (memTotal - memAvailable) * 1024,
	}, nil
}

func (lp *LinuxProvider) Cpu(ctx context.Context) (sys.Cpu, error) {
	pathStat := lp.procBase + "/stat"
	rawData, err := os.ReadFile(pathStat)
	if err != nil {
		return sys.Cpu{}, fmt.Errorf("failed to read %s file: %w", pathStat, err)
	}

	var count uint
	var timeFirstScreen []uint64
	var idleFirstScreen []uint64
	var totalTimeFirstScreen uint64
	var totalIdleFirstScreen uint64

	for line := range strings.SplitSeq(string(rawData), "\n") {
		if strings.HasPrefix(line, "cpu") {
			total, idle, err := parseCpuLine(line)
			if err != nil {
				return sys.Cpu{}, err
			}

			if strings.Fields(line)[0] == "cpu" {
				totalTimeFirstScreen = total
				totalIdleFirstScreen = idle
			} else {
				timeFirstScreen = append(timeFirstScreen, total)
				idleFirstScreen = append(idleFirstScreen, idle)
				count++
			}
		}
	}

	time.Sleep(1 * time.Second)

	rawData, err = os.ReadFile(pathStat)
	if err != nil {
		return sys.Cpu{}, fmt.Errorf("failed to read %s file: %w", pathStat, err)
	}

	var loads []float64
	var totalLoad float64
	var cpuIdx int

	for line := range strings.SplitSeq(string(rawData), "\n") {
		if strings.HasPrefix(line, "cpu") {
			total, idle, err := parseCpuLine(line)
			if err != nil {
				return sys.Cpu{}, err
			}

			if strings.Fields(line)[0] == "cpu" {
				totalLoad = (float64(total-totalTimeFirstScreen) - float64(idle-totalIdleFirstScreen)) / float64(total-totalTimeFirstScreen) * 100
			} else {
				loads = append(loads, (float64(total-timeFirstScreen[cpuIdx])-float64(idle-idleFirstScreen[cpuIdx]))/float64(total-timeFirstScreen[cpuIdx])*100)
				cpuIdx++
			}
		}
	}

	return sys.Cpu{
		Count:     count,
		Loads:     loads,
		TotalLoad: totalLoad,
	}, nil
}

func parseCpuLine(line string) (uint64, uint64, error) {
	var idle uint64
	var total uint64
	for i, strVal := range strings.Fields(line)[1:] {
		val, err := strconv.ParseUint(strVal, 10, 64)
		if err != nil {
			return 0, 0, fmt.Errorf("failed to parse string to integer: %w", err)
		}

		total += val
		if i == 3 {
			idle += val
		}
	}
	return total, idle, nil
}

func (lp *LinuxProvider) Disk(ctx context.Context) (sys.Disk, error) {
	pathMounts := lp.netProcBase + "/mounts"
	rawData, err := os.ReadFile(pathMounts)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s file: %w", pathMounts, err)
	}

	var paths []string
	for line := range strings.SplitSeq(string(rawData), "\n") {
		if strings.HasPrefix(line, "/dev/") {
			paths = append(paths, strings.Fields(line)[1])
		}
	}

	var result sys.Disk

	for _, path := range paths {
		var stat unix.Statfs_t
		if err = unix.Statfs(path, &stat); err != nil {
			return nil, fmt.Errorf("failed to get stats about partition: %w", err)
		}

		result = append(result, sys.Partition{
			Path:       path,
			TotalBytes: stat.Blocks * uint64(stat.Bsize),
			FreeBytes:  stat.Bfree * uint64(stat.Bsize),
			UsedBytes:  (stat.Blocks - stat.Bfree) * uint64(stat.Bsize),
		})
	}

	return result, nil
}

func (lp *LinuxProvider) Sys(ctx context.Context) (sys.Stats, error) {
	var result sys.Stats

	var err error
	result.UptimeSec, err = lp.Uptime(ctx)
	if err != nil {
		return sys.Stats{}, err
	}

	result.Cpu, err = lp.Cpu(ctx)
	if err != nil {
		return sys.Stats{}, err
	}

	result.Ram, err = lp.Ram(ctx)
	if err != nil {
		return sys.Stats{}, err
	}

	result.Disk, err = lp.Disk(ctx)
	if err != nil {
		return sys.Stats{}, err
	}

	return result, nil
}

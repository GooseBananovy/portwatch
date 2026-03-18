package linux

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/goosebananovy/portwatch/internal/model/metrics/procs"
)

func (lp *LinuxProvider) Procs(ctx context.Context) (procs.Procs, error) {
	entries, err := os.ReadDir(lp.procBase)
	if err != nil {
		return nil, fmt.Errorf("failed to get files entries in %s dir: %w", lp.procBase, err)
	}

	var result procs.Procs
	procWorked := make(map[uint64]uint64)

	cpuWorkedFirstScreen, err := cpuWorked(lp.procBase)
	if err != nil {
		return nil, err
	}

	ram, err := lp.Ram(ctx)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if _, err = strconv.Atoi(entry.Name()); entry.IsDir() && err == nil {
			var proc procs.Proc

			fileStatus := lp.procBase + "/" + entry.Name() + "/status"
			fileStat := lp.procBase + "/" + entry.Name() + "/stat"

			rawInfo, err := os.ReadFile(fileStatus)
			if err != nil {
				return nil, fmt.Errorf("failed to read "+fileStatus+"file: %w", err)
			}

			for line := range strings.SplitSeq(string(rawInfo), "\n") {
				if strings.HasPrefix(line, "Name") {
					proc.Name = strings.Fields(line)[1]
				} else if strings.HasPrefix(line, "Pid") {
					proc.Pid, err = strconv.ParseUint(strings.Fields(line)[1], 10, 64) //TODO: Чтобы убрать все проверки на negative бесполезные, надо везде сделать так просто
					if err != nil {
						return nil, fmt.Errorf("failed to convert pid to uint: %w", err)
					}
				} else if strings.HasPrefix(line, "VmRSS") {
					ramUsed, err := strconv.ParseUint(strings.Fields(line)[1], 10, 64)
					if err != nil {
						return nil, fmt.Errorf("failed to convert bytes quantity to uint: %w", err)
					}

					proc.Ram = float64(ramUsed*1024) / float64(ram.TotalBytes) * 100
				}
			}

			result = append(result, proc)

			utime, stime, err := parseProcStatFile(fileStat)
			if err != nil {
				return nil, err
			}

			procWorked[proc.Pid] = uint64(utime + stime)

		}
	}

	time.Sleep(1 * time.Second)

	cpuWorkedSecondScreen, err := cpuWorked(lp.procBase)
	if err != nil {
		return nil, err
	}

	for i := range result {
		pid := result[i].Pid
		fileStat := lp.procBase + "/" + strconv.FormatUint(pid, 10) + "/stat"

		utime, stime, err := parseProcStatFile(fileStat)
		if err != nil {
			continue
		}

		result[i].Cpu = float64((utime+stime)-procWorked[pid]) / float64(cpuWorkedSecondScreen-cpuWorkedFirstScreen) * 100
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Cpu > result[j].Cpu
	})

	if len(result) > 30 {
		result = result[:30]
	}

	return result, nil
}

func parseProcStatFile(path string) (uint64, uint64, error) {
	rawInfo, err := os.ReadFile(path)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to read "+path+"file: %w", err)
	}

	line := string(rawInfo)
	lineArgs := strings.Fields(line[strings.LastIndex(line, ")")+1:])

	utime, err := strconv.ParseInt(lineArgs[11], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to convert utime to int: %w", err)
	}

	stime, err := strconv.ParseInt(lineArgs[12], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to convert stime to int: %w", err)
	}

	if utime < 0 || stime < 0 {
		return 0, 0, errors.New("got negative time")
	}

	return uint64(utime), uint64(stime), nil
}

func cpuWorked(procBase string) (uint64, error) {
	pathStat := procBase + "/stat"
	rawData, err := os.ReadFile(pathStat)
	if err != nil {
		return 0, fmt.Errorf("failed to read %s file: %w", pathStat, err)
	}

	line := strings.Split(string(rawData), "\n")[0]
	if len(line) == 0 || strings.Fields(line)[0] != "cpu" {
		return 0, fmt.Errorf("invalid content in %s file: %w", pathStat, err)
	}

	var total uint64
	for i, strVal := range strings.Fields(line)[1:] {
		if i != 3 {
			val, err := strconv.ParseInt(strVal, 10, 64)
			if err != nil {
				return 0, fmt.Errorf("failed to parse string to integer: %w", err)
			}

			if val < 0 {
				return 0, errors.New("got negative time")
			}

			total += uint64(val)
		}
	}

	return total, nil
}

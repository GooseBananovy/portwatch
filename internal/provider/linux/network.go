package linux

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/goosebananovy/portwatch/internal/model/metrics/network"
)

func (lp *LinuxProvider) Network(ctx context.Context) (network.Stats, error) {
	pathNetDev := lp.procBase + "/1/net/dev"
	rawData, err := os.ReadFile(pathNetDev)
	if err != nil {
		return network.Stats{}, fmt.Errorf("failed to read %s file: %w", pathNetDev, err)
	}

	var ingoingBFirstScreen uint64
	var outgoingBFirstScreen uint64

	for line := range strings.SplitSeq(string(rawData), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "wlan") || strings.HasPrefix(line, "eth") || strings.HasPrefix(line, "ens") || strings.HasPrefix(line, "enp") {

			ingoingB, outgoingB, err := parseNetDevLine(line)
			if err != nil {
				return network.Stats{}, err
			}

			ingoingBFirstScreen += uint64(ingoingB)
			outgoingBFirstScreen += uint64(outgoingB)
		}
	}

	var traffic network.Traffic

	time.Sleep(1 * time.Second)

	rawData, err = os.ReadFile(pathNetDev)
	if err != nil {
		return network.Stats{}, fmt.Errorf("failed to read %s file: %w", pathNetDev, err)
	}

	for line := range strings.SplitSeq(string(rawData), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "wlan") || strings.HasPrefix(line, "eth") || strings.HasPrefix(line, "ens") || strings.HasPrefix(line, "enp") {
			ingoingB, outgoingB, err := parseNetDevLine(line)
			if err != nil {
				return network.Stats{}, err
			}

			traffic.IngoingBPS += ingoingB
			traffic.OutgoingBPS += outgoingB
		}
	}

	traffic.IngoingBPS -= ingoingBFirstScreen
	traffic.OutgoingBPS -= outgoingBFirstScreen

	pathNetTCP := lp.procBase + "/1/net/tcp"
	rawData, err = os.ReadFile(pathNetTCP)
	if err != nil {
		return network.Stats{}, fmt.Errorf("failed to read %s file: %w", pathNetTCP, err)
	}

	var tcpCount network.TCPCount

	for i, line := range strings.Split(string(rawData), "\n") {
		if i > 1 && len(strings.Fields(line)) > 3 && strings.Fields(line)[3] == "01" {
			tcpCount++
		}
	}

	return network.Stats{
		Traffic:  traffic,
		TcpCount: tcpCount,
	}, nil
}

func parseNetDevLine(line string) (uint64, uint64, error) {
	ingoingB, err := strconv.ParseInt(strings.Fields(line)[1], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to convert net dev content to int: %w", err)
	}

	if ingoingB < 0 {
		return 0, 0, errors.New("got negative recieve bytes quantity")
	}

	outgoingB, err := strconv.ParseInt(strings.Fields(line)[9], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to convert net dev content to int: %w", err)
	}

	if outgoingB < 0 {
		return 0, 0, errors.New("got negative transmit bytes quantity")
	}
	return uint64(ingoingB), uint64(outgoingB), nil
}

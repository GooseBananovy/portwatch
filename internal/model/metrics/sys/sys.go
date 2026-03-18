package sys

type Cpu struct {
	Count     uint      `json:"count_cores"`
	Loads     []float64 `json:"loads"`
	TotalLoad float64   `json:"total_load"`
}

type Ram struct {
	TotalBytes     uint64 `json:"total"`
	UsedBytes      uint64 `json:"used"`
	AvailableBytes uint64 `json:"available"`
}

type Disk []Partition

type Partition struct {
	Path       string `json:"path"`
	TotalBytes uint64 `json:"total"`
	UsedBytes  uint64 `json:"used"`
	FreeBytes  uint64 `json:"free"`
}

type UptimeSec uint64

type Stats struct {
	Cpu       Cpu       `json:"cpu"`
	Ram       Ram       `json:"ram"`
	Disk      Disk      `json:"disk"`
	UptimeSec UptimeSec `json:"uptime"`
}

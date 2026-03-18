package sys

type Cpu struct {
	Count     uint
	Loads     []float64
	TotalLoad float64
}

type Ram struct {
	TotalBytes     uint64
	UsedBytes      uint64
	AvailableBytes uint64
}

type Disk []Partition

type Partition struct {
	Path       string
	TotalBytes uint64
	UsedBytes  uint64
	FreeBytes  uint64
}

type UptimeSec uint64

type Stats struct {
	Cpu       Cpu
	Ram       Ram
	Disk      Disk
	UptimeSec UptimeSec
}

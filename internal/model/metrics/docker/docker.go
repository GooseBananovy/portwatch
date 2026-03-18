package docker

type Container struct {
	Name     string  `json:"name"`
	Status   string  `json:"status"`
	Cpu      float64 `json:"cpu"`
	RamUsed  uint64  `json:"ram_used"`
	RamTotal uint64  `json:"ram_total"`
}

type Cons []Container

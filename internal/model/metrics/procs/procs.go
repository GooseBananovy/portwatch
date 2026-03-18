package procs

type Proc struct {
	Name string  `json:"name"`
	Pid  uint64  `json:"pid"`
	Cpu  float64 `json:"cpu"`
	Ram  float64 `json:"ram"`
}

type Procs []Proc

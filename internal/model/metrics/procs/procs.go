package procs

type Proc struct {
	Name string
	Pid  uint64
	Cpu  float64
	Ram  float64
}

type Procs []Proc

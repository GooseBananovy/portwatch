package docker

type Container struct {
	Name     string
	Status   string
	Cpu      float64
	RamUsed  uint64
	RamTotal uint64
}

type Cons []Container

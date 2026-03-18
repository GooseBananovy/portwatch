package network

type Traffic struct {
	IngoingBPS  uint64
	OutgoingBPS uint64
}

type TCPCount uint

type Stats struct {
	Traffic  Traffic
	TcpCount TCPCount
}

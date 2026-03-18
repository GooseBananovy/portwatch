package network

type Traffic struct {
	IngoingBPS  uint64 `json:"ingoing_bps"`
	OutgoingBPS uint64 `json:"outgoing_bps"`
}

type TCPCount uint

type Stats struct {
	Traffic  Traffic  `json:"traffic"`
	TcpCount TCPCount `json:"tcp_cons_count"`
}

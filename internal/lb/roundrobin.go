package lb

import (
	"sync/atomic"
)

type RoundRobin struct {
	addrs []string
	idx   uint64
}

func NewRoundRobin(addrs []string) *RoundRobin {
	return &RoundRobin{addrs: addrs}
}

func (r *RoundRobin) Next() string {
	n := uint64(len(r.addrs))
	if n == 0 {
		return ""
	}
	i := atomic.AddUint64(&r.idx, 1)
	return r.addrs[(i-1)%n]
}

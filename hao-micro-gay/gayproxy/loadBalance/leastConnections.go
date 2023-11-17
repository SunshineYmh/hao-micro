package loadbalance

import (
	"sync"
)

type LeastConnectionBalancer struct {
	servers []string
	conns   []int32
	lock    sync.Mutex
}

func NewLeastConnectionBalancer() *LeastConnectionBalancer {
	return &LeastConnectionBalancer{}
}

func (lb *LeastConnectionBalancer) Add(server string) {
	lb.lock.Lock()
	defer lb.lock.Unlock()
	lb.servers = append(lb.servers, server)
	lb.conns = append(lb.conns, 0)
}

func (lb *LeastConnectionBalancer) Next() string {
	var (
		minLoc  int
		minConn = int32(^uint32(0) >> 1)
	)

	lb.lock.Lock()
	defer lb.lock.Unlock()

	for i := range lb.conns {
		if lb.conns[i] < minConn {
			minConn = lb.conns[i]
			minLoc = i
		}
	}

	lb.conns[minLoc] += 1
	return lb.servers[minLoc]
}

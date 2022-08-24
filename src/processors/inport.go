// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package processors

import (
	"sync"
)

type InPort struct {
	mu          sync.Mutex
	ch          chan interface{}
	name        string
	// 需要记录自己有多少个下游
	edges       []*OutPort
	closed      bool
	closeCounts int
}

func NewInPort(name string) *InPort {
	return &InPort{
		name: name,
		ch:   make(chan interface{}),
	}
}

func (pt *InPort) Name() string {
	return pt.name
}

func (pt *InPort) AddEdge(rpt *OutPort) {
	pt.mu.Lock()
	defer pt.mu.Unlock()
	pt.edges = append(pt.edges, rpt)
}

func (pt *InPort) From(rpt *OutPort) {
	rpt.AddEdge(pt)
	pt.AddEdge(rpt)
}

func (pt *InPort) Send(v interface{}) {
	pt.mu.Lock()
	defer pt.mu.Unlock()

	if pt.closed {
		return
	}
	pt.ch <- v
}

func (pt *InPort) Recv() <-chan interface{} {
	return pt.ch
}

// 由发送端关闭 channel。
func (pt *InPort) Close() {
	pt.mu.Lock()
	defer pt.mu.Unlock()
	if !pt.closed {
		// 全部 OutPort 调了 Close，才会真正关闭 channal。
		pt.closeCounts++
		if pt.closeCounts >= len(pt.edges) {
			close(pt.ch)
			pt.closed = true
		}
	}
}

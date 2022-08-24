// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package processors

import (
	"context"
	"sync/atomic"
	"time"
)

type (
	NextFunc     func(interface{})
	DoneFunc     func()
	EventHandler interface{}

	BaseProcessor struct {
		in          *InPort
		out         *OutPort
		name        string
		ctx         context.Context
		duration    time.Duration
		pauseChan   chan struct{}
		resumeChan  chan struct{}
		nextHandler NextFunc
		doneHandler DoneFunc
	}
)

func NewBaseProcessor(name string) BaseProcessor {
	return BaseProcessor{
		in:         NewInPort(name),
		out:        NewOutPort(name),
		ctx:        context.Background(),
		name:       name,
		pauseChan:  make(chan struct{}),
		resumeChan: make(chan struct{}),
	}
}

func (p *BaseProcessor) Name() string {
	return p.name
}

func (p *BaseProcessor) Duration() time.Duration {
	return time.Duration(atomic.LoadInt64((*int64)(&p.duration)))
}

func (p *BaseProcessor) In() *InPort {
	return p.in
}

func (p *BaseProcessor) Out() *OutPort {
	return p.out
}

func (p *BaseProcessor) To(receivers ...IProcessor) {
	for _, receiver := range receivers {
		p.out.To(receiver.In())
	}
}

func (p *BaseProcessor) From(senders ...IProcessor) {
	for _, sender := range senders {
		source := sender.Out()
		p.in.From(source)
	}
}

func (p *BaseProcessor) Execute() {
	// Nothing.
}

func (p *BaseProcessor) Pause() {
	p.pauseChan <- struct{}{}
}

func (p *BaseProcessor) Resume() {
	p.resumeChan <- struct{}{}
}

func (p *BaseProcessor) SetContext(ctx context.Context) {
	p.ctx = ctx
}

func (p *BaseProcessor) Subscribe(eventHandlers ...EventHandler) {
	in := p.In()
	out := p.Out()
	ctx := p.ctx

	for _, handler := range eventHandlers {
		switch handler := handler.(type) {
		case func():
			p.doneHandler = handler
		case func(interface{}):
			p.nextHandler = handler
		}
	}

	defer func() {
		out.Close()
		close(p.pauseChan)
		close(p.resumeChan)
	}()

	for {
	Loop:
		select {
		case <-p.pauseChan:
			for range p.resumeChan {
				goto Loop
			}
			return

		case <-ctx.Done():
			if p.nextHandler != nil {
				p.nextHandler(ctx.Err())
			}
			return
		case x, ok := <-in.Recv():
			if !ok { // 如果接收端没有新数据，就执行自己的逻辑，否则执行 nextHandler 处理接收到的数据。
				if p.doneHandler != nil {
					start := time.Now()
					p.doneHandler()
					atomic.AddInt64((*int64)(&p.duration), int64(time.Since(start)))
				}
				return
			}
			if p.nextHandler != nil {
				start := time.Now()
				p.nextHandler(x)
				atomic.AddInt64((*int64)(&p.duration), int64(time.Since(start)))
			}
		}
	}
}

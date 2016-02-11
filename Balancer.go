package queue

import (
	"fmt"
	"log"
	"sync"
)

func NewBalancer(name string, size, queueSize int) *Balancer { // {{{
	return &Balancer{
		name:      name,
		size:      size,
		queueSize: queueSize,
		shutdown:  make(chan bool),
		pools:     &Pools{},
	}
} // }}}

type Balancer struct {
	sync.Mutex

	name          string
	size          int
	queueSize     int
	shutdown      chan bool
	workerFactory WorkerFactoryInterface
	pools         *Pools
}

func (this *Balancer) Name() string { // {{{
	return this.name
} // }}}

func (this *Balancer) Size() int { // {{{
	return this.size
} // }}}

func (this *Balancer) Info() string { // {{{
	return fmt.Sprintf("%sBalancer:%d#%d", this.Name(), this.size, this.queueSize)
} // }}}

func (this *Balancer) SetWorkerFactory(f WorkerFactoryInterface) { // {{{
	this.workerFactory = f
} // }}}

func (this *Balancer) WorkerFactory() WorkerFactoryInterface { // {{{
	if this.workerFactory == nil {
		this.workerFactory = NewWorkerFactory(this.name)
	}

	return this.workerFactory
} // }}}

func (this *Balancer) Dispatch(job JobInterface) { // {{{
	if this.pools.Len() > 0 {
		this.pools.Dispatch(job)
	} else {
		go log.Println("[Balancer:Error] Pools is empty")
	}
} // }}}

func (this *Balancer) Run() { // {{{
	for i := 1; i <= this.size; i++ {
		pool := NewPool(this.name, i, this.queueSize)
		pool.SetWorkerFactory(this.WorkerFactory().New())
		this.pools.Push(pool)
		pool.Start()
	}
} // }}}

func (this *Balancer) Close() { // {{{
	for this.pools.Len() > 0 {
		pool := this.pools.Pop()
		pool.Close()
	}
} // }}}

func (this *Balancer) Start() { // {{{
	go this.Run()
} // }}}

func (this *Balancer) Stop() { // {{{
	go this.Close()
} // }}}

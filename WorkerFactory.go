package queue

import (
	"sync"
)

func NewWorkerFactory(name string) *WorkerFactory { // {{{
	factory := &WorkerFactory{}
	factory.Initialize(name)

	return factory
} // }}}

type WorkerFactory struct {
	WorkerFactoryInterface

	sync.Mutex

	id   int
	name string
}

func (this *WorkerFactory) Initialize(name string) { // {{{
	this.id = 0
	this.name = name
} // }}}

func (this *WorkerFactory) Name() string { // {{{
	return this.name
} // }}}

func (this *WorkerFactory) NextId() int { // {{{
	this.Lock()
	defer this.Unlock()
	this.id++

	return this.id
} // }}}

func (this *WorkerFactory) Create(poolId int, workerQueue WorkerQueue) WorkerInterface { // {{{
	return NewWorker(this.Name(), poolId, this.NextId(), workerQueue)
} // }}}

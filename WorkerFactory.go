package queue

import (
	"fmt"
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

func (this *WorkerFactory) New() WorkerFactoryInterface { // {{{
	return NewWorkerFactory(this.Name())
} // }}}

func (this *WorkerFactory) Name() string { // {{{
	return this.name
} // }}}

func (this *WorkerFactory) Info() string { // {{{
	return fmt.Sprintf("%sWorkerFactory#%d", this.name, this.id)
} // }}}

func (this *WorkerFactory) NextId() int { // {{{
	this.Lock()
	defer this.Unlock()
	this.id++

	return this.id
} // }}}

func (this *WorkerFactory) Create(poolId int, workerQueue WorkerQueue) WorkerInterface { // {{{
	return NewWorker(this.name, poolId, this.NextId(), workerQueue)
} // }}}

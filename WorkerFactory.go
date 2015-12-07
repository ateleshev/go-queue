package queue

import (
	"log"
	"sync"
)

func NewWorkerFactory(name string, workerQueue WorkerQueue, logger *log.Logger) *WorkerFactory { // {{{
	factory := &WorkerFactory{}
	factory.Initialize(name, workerQueue, logger)

	return factory
} // }}}

type WorkerFactory struct {
	WorkerFactoryInterface

	sync.Mutex

	id          int
	name        string
	workerQueue WorkerQueue
	logger      *log.Logger
}

func (this *WorkerFactory) Initialize(name string, workerQueue WorkerQueue, logger *log.Logger) { // {{{
	this.id = 0
	this.name = name
	this.workerQueue = workerQueue
	this.logger = logger
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

func (this *WorkerFactory) WorkerQueue() WorkerQueue { // {{{
	return this.workerQueue
} // }}}

func (this *WorkerFactory) Logger() *log.Logger { // {{{
	return this.logger
} // }}}

func (this *WorkerFactory) Create() WorkerInterface { // {{{
	return NewWorker(this.Name(), this.NextId(), this.WorkerQueue(), this.Logger())
} // }}}

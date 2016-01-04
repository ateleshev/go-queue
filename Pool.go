package queue

import (
	"fmt"
	//	"log"
)

func NewPool(name string, size, queueSize int) *Pool { // {{{
	return &Pool{
		name:      name,
		size:      size,
		queueSize: queueSize,
		shutdown:  make(chan bool),
		workers:   make(map[int]WorkerMap, size),
		queue:     make(PoolQueue, size),
	}
} // }}}

type Pool struct {
	name      string
	size      int
	queueSize int
	shutdown  chan bool
	factory   WorkerFactoryInterface
	workers   map[int]WorkerMap
	queue     PoolQueue
}

func (this *Pool) Name() string { // {{{
	return this.name
} // }}}

func (this *Pool) Size() int { // {{{
	return this.size
} // }}}

func (this *Pool) Info() string { // {{{
	return fmt.Sprintf("%s:%d#%d", this.name, this.size, this.queueSize)
} // }}}

func (this *Pool) SetFactory(factory WorkerFactoryInterface) { // {{{
	this.factory = factory
} // }}}

func (this *Pool) Factory() WorkerFactoryInterface { // {{{
	if this.factory == nil {
		this.factory = NewWorkerFactory(this.Name())
	}

	return this.factory
} // }}}

func (this *Pool) Perform(job JobInterface) { // {{{
	workerQueue := <-this.queue //  |<-- 1. Select pool
	worker := <-workerQueue     //  |    2. Get worker from pool
	worker <- job               //  |    3. Start executing of job
	this.queue <- workerQueue   //  |--> 4. Wait next job
} // }}}

func (this *Pool) Run() { // {{{
	for i := 0; i < this.Size(); i++ {
		workerMap := make(WorkerMap, this.queueSize)
		workerQueue := make(WorkerQueue, this.queueSize)
		for n := 0; n < this.queueSize; n++ {
			worker := this.Factory().Create(i+1, workerQueue)
			worker.Start()
			workerMap[worker.Id()] = worker
		}

		this.workers[i] = workerMap
		this.queue <- workerQueue
	}
} // }}}

// Close channels
func (this *Pool) close() { // {{{
	close(this.queue)
} // }}}

func (this *Pool) Close() { // {{{
	defer this.close()

	// Stop workers
	for i := 0; i < this.Size(); i++ {
		for id, worker := range this.workers[i] {
			worker.Stop()
			delete(this.workers[i], id)
		}
	}
} // }}}

func (this *Pool) Start() { // {{{
	go this.Run()
} // }}}

func (this *Pool) Stop() { // {{{
	go this.Close()
} // }}}

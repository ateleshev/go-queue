package queue

import (
	"fmt"
)

func NewWorker(name string, poolId int, id int, workerQueue WorkerQueue) *Worker { // {{{
	worker := &Worker{}
	worker.Initialize(name, poolId, id, worker, workerQueue)

	return worker
} // }}}

type Worker struct {
	id          int
	poolId      int
	name        string
	instance    interface{}
	shutdown    chan bool
	jobQueue    JobQueue
	workerQueue WorkerQueue
	pending     uint
}

func (this *Worker) Initialize(name string, poolId int, id int, instance interface{}, workerQueue WorkerQueue) { // {{{
	this.id = id
	this.poolId = poolId
	this.name = name
	this.instance = instance
	this.shutdown = make(chan bool)
	this.jobQueue = make(JobQueue)
	this.workerQueue = workerQueue
} // }}}

func (this *Worker) Id() int { // {{{
	return this.id
} // }}}

func (this *Worker) PoolId() int { // {{{
	return this.poolId
} // }}}

func (this *Worker) Name() string { // {{{
	return this.name
} // }}}

func (this *Worker) Info() string { // {{{
	return fmt.Sprintf("%sWorker:%d#%d", this.name, this.poolId, this.id)
} // }}}

func (this *Worker) Run() { // {{{
	// log.Printf("[Worker:%s#%d] Run\n", this.Name(), this.Id())
	for {
		// Add ourselves into the worker queue.
		this.workerQueue <- this.jobQueue

		select {
		case job := <-this.jobQueue:
			this.pending++
			// log.Printf("[Worker:%s#%d] Execute Job\n", this.Name(), this.Id())
			job.Execute(this.instance)
			job.Done()
			this.pending--
			break
		case <-this.shutdown:
			defer this.close()
			// log.Printf("[Worker:%s#%d] Shutdown\n", this.Name(), this.Id())
			return
		}
	}
} // }}}

func (this *Worker) close() { // {{{
	close(this.jobQueue)
} // }}}

func (this *Worker) Close() { // {{{
	this.shutdown <- true
} // }}}

func (this *Worker) Start() { // {{{
	go this.Run()
} // }}}

func (this *Worker) Stop() { // {{{
	go this.Close()
} // }}}

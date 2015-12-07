package queue

import (
	"log"
)

func NewWorker(name string, id int, workerQueue WorkerQueue, logger *log.Logger) *Worker { // {{{
	worker := &Worker{}
	worker.Initialize(name, id, worker, workerQueue, logger)

	return worker
} // }}}

type Worker struct {
	id          int
	name        string
	child       interface{}
	shutdown    chan bool
	jobQueue    JobQueue
	workerQueue WorkerQueue
	logger      *log.Logger
}

func (this *Worker) Initialize(name string, id int, child interface{}, workerQueue WorkerQueue, logger *log.Logger) { // {{{
	this.id = id
	this.name = name
	this.child = child
	this.shutdown = make(chan bool)
	this.jobQueue = make(JobQueue)
	this.workerQueue = workerQueue
	this.logger = logger
} // }}}

func (this *Worker) Id() int { // {{{
	return this.id
} // }}}

func (this *Worker) Name() string { // {{{
	return this.name
} // }}}

func (this *Worker) Run() { // {{{
	// log.Printf("[Worker:%s#%d] Run\n", this.Name(), this.Id())
	for {
		// Add ourselves into the worker queue.
		this.workerQueue <- this.jobQueue

		select {
		case job := <-this.jobQueue:
			// log.Printf("[Worker:%s#%d] Execute Job\n", this.Name(), this.Id())
			job.Execute(this.child)
			job.Done()
			break
		case <-this.shutdown:
			// log.Printf("[Worker:%s#%d] Shutdown\n", this.Name(), this.Id())
			return
		}
	}
} // }}}

func (this *Worker) close() { // {{{
	close(this.jobQueue)
} // }}}

func (this *Worker) Close() { // {{{
	defer this.close()
	this.shutdown <- true
} // }}}

func (this *Worker) Start() { // {{{
	go this.Run()
} // }}}

func (this *Worker) Stop() { // {{{
	go this.Close()
} // }}}

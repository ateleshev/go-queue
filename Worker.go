package queue

//import (
//	"log"
//)

func NewWorker(name string, id int, workerQueue WorkerQueue) *Worker {
	return &Worker{
		// == Protected ==
		id:       id,
		name:     name,
		shutdown: make(chan bool),
		// == Public ==
		JobQueue:    make(JobQueue),
		WorkerQueue: workerQueue,
	}
}

type Worker struct {
	WorkerInterface

	id          int
	name        string
	shutdown    chan bool
	JobQueue    JobQueue
	WorkerQueue WorkerQueue
}

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
		this.WorkerQueue <- this.JobQueue

		select {
		case job := <-this.JobQueue:
			// log.Printf("[Worker:%s#%d] Execute Job\n", this.Name(), this.Id())
			job.Execute()
			job.Done()
			break
		case <-this.shutdown:
			// log.Printf("[Worker:%s#%d] Shutdown\n", this.Name(), this.Id())
			return
		}
	}
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

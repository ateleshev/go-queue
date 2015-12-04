package queue

import (
	"fmt"
	"log"
	"os"
)

var (
	ServerLoggerPrefixFormat = "[QS:%s] "
	ServerLoggerFlags        = log.Ldate | log.Ltime // log.Ldate | log.Ltime | log.Lshortfile // log.Ldate | log.Ltime | log.Llongfile
)

func NewServer(name string, depthJobsQueue int, depthWorkersQueue int) *Server { // {{{
	return &Server{
		name:     name,
		workers:  make(map[int]*Worker, depthWorkersQueue),
		shutdown: make(chan bool),

		// == Depth ==
		DepthJobsQueue:    depthJobsQueue,
		DepthWorkersQueue: depthWorkersQueue,

		// == Queue ==
		JobQueue:    make(JobQueue, depthJobsQueue),
		WorkerQueue: make(WorkerQueue, depthWorkersQueue),

		Logger:    log.New(os.Stderr, fmt.Sprintf(ServerLoggerPrefixFormat, name), ServerLoggerFlags),
		Logging:   true,
		Debugging: true,
	}
} // }}}

type Server struct {
	name string

	workers  map[int]*Worker
	shutdown chan bool

	// [Depth]
	DepthJobsQueue    int // depth of jobs queue
	DepthWorkersQueue int // depth of workers queue
	// [Queue]
	JobQueue    JobQueue
	WorkerQueue WorkerQueue

	Logger    *log.Logger
	Logging   bool
	Debugging bool
}

func (this *Server) log(format string, v ...interface{}) { // {{{
	if this.Logging {
		this.Logger.Printf(format, v...)
	}
} // }}}

func (this *Server) perform(job JobInterface) { // {{{
	worker := <-this.WorkerQueue
	worker <- job

	go this.log("Execute job\n")
} // }}}

func (this *Server) serve() { // {{{
	for {
		select {
		case job := <-this.JobQueue:
			go this.perform(job)
		case <-this.shutdown:
			defer this.Close()
			go this.log("Shutdown\n")
			return
		}
	}
} // }}}

func (this *Server) Name() string { // {{{
	return this.name
} // }}}

func (this *Server) Dispatch(job JobInterface) { // {{{
	this.JobQueue <- job
} // }}}

func (this *Server) Run() { // {{{
	for i := 0; i < this.DepthWorkersQueue; i++ {
		worker := NewWorker(this.Name(), i+1, this.WorkerQueue, this.Logger)
		worker.Start()
		this.workers[worker.Id()] = worker
	}
	go this.log("Run %d workers\n", this.DepthWorkersQueue)
	go this.serve()
} // }}}

func (this *Server) Close() { // {{{
	for id, worker := range this.workers {
		worker.Stop()
		delete(this.workers, id)
	}
} // }}}

func (this *Server) Start() { // {{{
	go this.Run()
} // }}}

func (this *Server) Stop() { // {{{
	go this.Close()
} // }}}

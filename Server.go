package queue

import (
	"log"
)

// poolSize = 30 (Max duration)
// queueSize = 1000 (Max hotels per city)

func NewServer(name string, poolSize, queueSize int) *Server { // {{{
	size := poolSize * queueSize
	return &Server{
		size:     size,
		pool:     NewPool(name, poolSize, queueSize),
		jobQueue: make(JobQueue, size),
		shutdown: make(chan bool),

		Logging:   true,
		Debugging: true,
	}
} // }}}

type Server struct {
	size     int // depth of JobQueue
	pool     *Pool
	jobQueue JobQueue
	shutdown chan bool

	Logger    *log.Logger
	Logging   bool
	Debugging bool
}

func (this *Server) log(format string, v ...interface{}) { // {{{
	if this.Logging && this.Logger != nil {
		this.Logger.Printf(format, v...)
	}
} // }}}

func (this *Server) Name() string { // {{{
	return this.pool.Name()
} // }}}

func (this *Server) Info() string { // {{{
	return this.pool.Info()
} // }}}

func (this *Server) SetWorkerFactory(workerFactory WorkerFactoryInterface) { // {{{
	this.pool.SetFactory(workerFactory)
} // }}}

func (this *Server) Dispatch(job JobInterface) { // {{{
	this.jobQueue <- job
} // }}}

func (this *Server) Serve() { // {{{
	for {
		select {
		case job := <-this.jobQueue:
			go this.pool.Perform(job)
			//go this.log("[Server] Perform")
		case <-this.shutdown:
			defer this.close()
			//go this.log("[Server] Shutdown")
			return
		}
	}
} // }}}

func (this *Server) Run() { // {{{
	this.pool.Run()

	go this.Serve()
} // }}}

// Close channels
func (this *Server) close() { // {{{
	defer close(this.jobQueue)
	this.pool.Close()
} // }}}

func (this *Server) Close() { // {{{
	this.shutdown <- true
} // }}}

func (this *Server) Start() { // {{{
	go this.Run()
} // }}}

func (this *Server) Stop() { // {{{
	go this.Close()
} // }}}

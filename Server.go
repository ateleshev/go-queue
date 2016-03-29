package queue

import (
	"fmt"
	"log"
)

// poolSize = 30 (Max duration)
// queueSize = 1000 (Max hotels per city)

func NewServer(name string, poolSize, queueSize int) *Server { // {{{
	size := poolSize * queueSize
	return &Server{
		name:     name,
		size:     size,
		balancer: NewBalancer(name, poolSize, queueSize),
		jobQueue: make(JobQueue, size),
		shutdown: make(chan bool),

		Logging:   true,
		Debugging: true,
	}
} // }}}

type Server struct {
	name     string
	size     int // depth of JobQueue
	balancer *Balancer
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
	return this.name
} // }}}

func (this *Server) Info() string { // {{{
	return fmt.Sprintf("%sServer(%s)", this.Name(), this.balancer.Info())
} // }}}

func (this *Server) SetWorkerFactory(f WorkerFactoryInterface) { // {{{
	this.balancer.SetWorkerFactory(f)
} // }}}

func (this *Server) Dispatch(job JobInterface) { // {{{
	this.jobQueue <- job
} // }}}

func (this *Server) Serve() { // {{{
	for {
		select {
		case job := <-this.jobQueue:
			this.balancer.Dispatch(job)
			//go this.log("[Server] Perform")
		case <-this.shutdown:
			defer this.close()
			//go this.log("[Server] Shutdown")
			return
		}
	}
} // }}}

func (this *Server) Run() { // {{{
	this.balancer.Run()

	go this.Serve()
} // }}}

// Close channels
func (this *Server) close() { // {{{
	defer close(this.jobQueue)
	this.balancer.Close()
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

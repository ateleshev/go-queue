package queue

import (
	//	"log"
	"time"
)

type Job struct {
	JobInterface

	// == Protected ==
	done chan bool
}

func (this *Job) Init() { // {{{
	// log.Printf("[Job] Init\n")
	this.done = make(chan bool)
} // }}}

//func (this *Job) Execute() { // {{{
//	log.Printf("[Job] Execute\n")
//} // }}}

func (this *Job) Wait() { // {{{
	select {
	case <-this.done:
		// log.Println("[Job] Finished")
		break
	}
} // }}}

func (this *Job) WaitDeadline(duration time.Duration) { // {{{
	select {
	case <-this.done:
		// log.Println("[Job] Finished")
		break
	case <-time.After(duration):
		// log.Println("[Job] Timed out")
		break
	}
} // }}}

func (this *Job) Done() { // {{{
	// log.Println("[Job] Done")
	this.done <- true
} // }}}

func (this *Job) Completed() <-chan bool { // {{{
	// log.Println("[Job] Completed")
	return this.done
} // }}}

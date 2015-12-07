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

func (this *Job) Initialize() { // {{{
	// log.Printf("[Job] Initialize\n")
	this.done = make(chan bool)
} // }}}

//func (this *Job) Execute(w interface{}) { // {{{
//  worker := w.(*Worker)
//	log.Printf("[Job:%s#%d] Execute\n", worker.Name(), worker.Id())
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

package queue_test

import (
	"testing"

	"github.com/ateleshev/go-queue"
)

type TestJob struct {
	queue.Job

	Message  string
	WorkerId int
}

func (this *TestJob) Execute(w interface{}) { // {{{
	worker := w.(*queue.Worker)
	this.WorkerId = worker.Id()
} // }}}

func TestExecuteJob(t *testing.T) { // {{{
	workerQueue := make(queue.WorkerQueue)
	defer close(workerQueue)

	w := queue.NewWorker("TestJobWorker", 1, 1, workerQueue)
	t.Logf("Created new Worker: '%s'", w.Info())

	j := &TestJob{Message: "Get worker identifier"}
	j.Initialize()
	t.Logf("Created new Job: '%s'", j.Message)

	j.Execute(w)
	if w.Id() != j.WorkerId {
		t.Errorf("[Error:TestJob] Job does not executed: Worker.Id=%d | Job.WorkerId=%d\n", w.Id(), j.WorkerId)
	} else {
		t.Logf("Job executed")
	}
} // }}}

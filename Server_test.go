package queue_test

import (
	"sync"
	"testing"
	"time"

	"github.com/ArtemTeleshev/go-queue"
)

type TestServerJob struct {
	queue.Job

	Message    string
	WorkerInfo string
}

func (this *TestServerJob) Execute(w interface{}) { // {{{
	worker := w.(*queue.Worker)

	// Sleep 30ms
	time.Sleep(30 * time.Millisecond)

	// Perform work
	this.WorkerInfo = worker.Info()
} // }}}

func finalizeTestServerJob(j *TestServerJob, i int, wg *sync.WaitGroup, t *testing.T) { // {{{
	defer wg.Done()
	startedAt := time.Now()
	j.Wait()
	finishedAt := time.Now()
	t.Logf("[TestServerJob:%d] Executed: '%s' -> [%.6fs]", i, j.WorkerInfo, finishedAt.Sub(startedAt).Seconds())
} // }}}

func TestServer(t *testing.T) { // {{{
	s := queue.NewServer("TestWorker", PoolSize, QueueSize)
	t.Logf("[Server] Created new: '%s'", s.Name())

	s.Start()
	defer s.Close()
	t.Logf("[Server] Started: '%s'", s.Info())

	var j *TestServerJob
	var wg sync.WaitGroup
	startedAt := time.Now()

	for i := 0; i < JobsQuantity; i++ {
		j = &TestServerJob{Message: "Get worker info from server"}
		j.Initialize()
		t.Logf("[TestServerJob:%d] Created: '%s'", i, j.Message)

		wg.Add(1)
		s.Dispatch(j)
		go finalizeTestServerJob(j, i, &wg, t)
	}

	wg.Wait()
	finishedAt := time.Now()
	t.Logf("[Server] Executed %d of tasks -> [%.6fs]", JobsQuantity, finishedAt.Sub(startedAt).Seconds())
} // }}}

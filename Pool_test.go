package queue_test

import (
	"sync"
	"testing"
	"time"

	"github.com/ArtemTeleshev/go-queue"
)

type TestPoolJob struct {
	queue.Job

	Message    string
	WorkerInfo string
}

func (this *TestPoolJob) Execute(w interface{}) { // {{{
	worker := w.(*queue.Worker)

	// Sleep 25ms
	time.Sleep(25 * time.Millisecond)

	// Perform work
	this.WorkerInfo = worker.Info()
} // }}}

func finalizeTestPoolJob(j *TestPoolJob, i int, wg *sync.WaitGroup, t *testing.T) { // {{{
	defer wg.Done()
	startedAt := time.Now()
	j.Wait()
	finishedAt := time.Now()
	t.Logf("[TestPoolJob:%d] Executed: '%s' -> [%.6fs]", i, j.WorkerInfo, finishedAt.Sub(startedAt).Seconds())
} // }}}

func TestPool(t *testing.T) { // {{{
	p := queue.NewPool("TestWorker", PoolSize, QueueSize)
	t.Logf("[Pool] Created new: '%s'", p.Name())

	p.Start()
	defer p.Close()
	t.Logf("[Pool] Started: '%s'", p.Info())

	var j *TestPoolJob
	var wg sync.WaitGroup
	startedAt := time.Now()

	for i := 0; i < JobsQuantity; i++ {
		j = &TestPoolJob{Message: "Get worker info from pool"}
		j.Initialize()
		t.Logf("[TestPoolJob:%d] Created: '%s'", i, j.Message)

		wg.Add(1)
		p.Perform(j)
		go finalizeTestPoolJob(j, i, &wg, t)
	}

	wg.Wait()
	finishedAt := time.Now()
	t.Logf("[Pool] Executed %d of tasks -> [%.6fs]", JobsQuantity, finishedAt.Sub(startedAt).Seconds())
} // }}}

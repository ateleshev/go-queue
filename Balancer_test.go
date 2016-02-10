package queue_test

import (
	"sync"
	"testing"
	"time"

	"github.com/ArtemTeleshev/go-queue"
)

type TestBalancerJob struct {
	queue.Job

	Message    string
	WorkerInfo string
}

func (this *TestBalancerJob) Execute(w interface{}) { // {{{
	worker := w.(*queue.Worker)

	// Sleep 25ms
	time.Sleep(25 * time.Millisecond)

	// Perform work
	this.WorkerInfo = worker.Info()
} // }}}

func finalizeTestBalancerJob(j *TestBalancerJob, i int, wg *sync.WaitGroup, t *testing.T) { // {{{
	defer wg.Done()
	startedAt := time.Now()
	j.Wait()
	finishedAt := time.Now()
	t.Logf("[TestBalancerJob:%d] Executed: '%s' -> [%.6fs]", i, j.WorkerInfo, finishedAt.Sub(startedAt).Seconds())
} // }}}

func TestBalancer(t *testing.T) { // {{{
	b := queue.NewBalancer("TestCacheAPISearch", PoolSize, QueueSize)
	t.Logf("[Balancer] Created new: '%s'", b.Name())

	b.Run()
	defer b.Close()
	t.Logf("[Balancer] Started: '%s'", b.Info())

	var j *TestBalancerJob
	var wg sync.WaitGroup
	startedAt := time.Now()

	for i := 0; i < JobsQuantity; i++ {
		j = &TestBalancerJob{Message: "Get worker info from Balancer"}
		j.Initialize()
		t.Logf("[TestBalancerJob:%d] Created: '%s'", i, j.Message)

		wg.Add(1)
		b.Dispatch(j)
		go finalizeTestBalancerJob(j, i, &wg, t)
	}

	wg.Wait()
	finishedAt := time.Now()
	t.Logf("[Balancer] Executed %d of tasks -> [%.6fs]", JobsQuantity, finishedAt.Sub(startedAt).Seconds())
} // }}}

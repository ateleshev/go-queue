package queue_test

import (
	"testing"

	"github.com/ArtemTeleshev/go-queue"
)

func TestWorkerFactory(t *testing.T) { // {{{
	workerQueue := make(queue.WorkerQueue)
	defer close(workerQueue)

	f := queue.NewWorkerFactory("TestWorker")
	t.Logf("Created new WorkerFactory: '%s'", f.Name())

	w1 := f.Create(1, workerQueue)
	t.Logf("Created new Worker: '%s'", w1.Info())

	w2 := f.Create(1, workerQueue)
	t.Logf("Created new Worker: '%s'", w2.Info())

	if w1.Id() >= w2.Id() {
		t.Errorf("[Error:TestWorkerFactory] Factory cannot increment identifier of worker: %d, %d \n", w1.Id(), w2.Id())
	} else {
		t.Logf("TestWorkerFactory works correct: '%s', '%s'", w1.Info(), w2.Info())
	}
} // }}}

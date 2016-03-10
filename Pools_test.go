package queue_test

import (
	"testing"
	"time"

	"github.com/ateleshev/go-queue"
)

func TestPools(t *testing.T) { // {{{
	startedAt := time.Now()

	pools := &queue.Pools{}
	p1 := queue.NewPool("TestPool", 1, QueueSize)
	t.Logf("[Pool] Created new: '%s'", p1.Name())
	pools.Push(p1)

	p2 := queue.NewPool("TestPool", 2, QueueSize)
	t.Logf("[Pool] Created new: '%s'", p2.Name())
	pools.Push(p2)

	for i := 0; i < 20; i++ {
		t.Logf("[Pool] Index: '%d'", pools.Index())
	}

	finishedAt := time.Now()
	t.Logf("[Pool] Executed [%.6fs]", finishedAt.Sub(startedAt).Seconds())
} // }}}

package queue

import (
	"log"
)

type WorkerInterface interface {
	Initialize(string, int, interface{}, WorkerQueue, *log.Logger)

	Id() int
	Name() string

	// In current goroutine
	Run()
	Close()

	// In new goroutine (background)
	Start()
	Stop()
}

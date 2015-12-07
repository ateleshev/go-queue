package queue

import (
	"log"
)

type WorkerFactoryInterface interface {
	Initialize(string, WorkerQueue, *log.Logger)

	Name() string
	NextId() int
	WorkerQueue() WorkerQueue
	Logger() *log.Logger

	Create() WorkerInterface
}

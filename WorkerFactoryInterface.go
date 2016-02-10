package queue

type WorkerFactoryInterface interface {
	Initialize(string) // ({Name})

	New() WorkerFactoryInterface
	Name() string
	NextId() int

	Create(int, WorkerQueue) WorkerInterface // ({PoolId}, {WorkerQueue})
}

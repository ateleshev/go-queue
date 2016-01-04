package queue

type WorkerFactoryInterface interface {
	Initialize(string) // ({Name})

	Name() string
	NextId() int

	Create(int, WorkerQueue) WorkerInterface // ({PoolId}, {WorkerQueue})
}

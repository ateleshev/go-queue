package queue

type WorkerInterface interface {
	Initialize(string, int, int, interface{}, WorkerQueue) // ({Name}, {PoolId}, {Id}, {WorkerInstance}, {WorkerQueue})

	Id() int
	PoolId() int
	Name() string
	Info() string

	// In current goroutine
	Run()
	Close()

	// In new goroutine (background)
	Start()
	Stop()
}

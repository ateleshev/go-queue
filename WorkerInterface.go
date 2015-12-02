package queue

type WorkerInterface interface {
	Id() int
	Name() string
	Run()
	Close()
	Start()
	Stop()
}

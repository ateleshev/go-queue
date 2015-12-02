package queue

type JobInterface interface {
	Init()
	Execute()
	Wait()
	Done()
	Completed() <-chan bool
}

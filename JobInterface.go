package queue

type JobInterface interface {
	Initialize()
	Execute()
	Wait()
	Done()
	Completed() <-chan bool
}

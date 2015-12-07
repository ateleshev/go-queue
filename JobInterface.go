package queue

type JobInterface interface {
	Initialize()
	Execute(interface{})
	Wait()
	Done()
	Completed() <-chan bool
}

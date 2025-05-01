package main

type Mess interface {
	Submit(task *task)
	ReadTask() *task
}

type mess struct {
	Task_ chan *task
}

func (tar *mess) Submit(task *task) {
	tar.Task_ <- task
}

func (tar *mess) ReadTask() *task {
	return <-tar.Task_
}

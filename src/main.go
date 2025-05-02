package main

func main() {
	var worker workers
	var mess chan *mess
	worker.Init(cimess{}, mess)
	var listener listener
	listener.Init()

}

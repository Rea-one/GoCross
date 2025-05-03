package main

type Server interface {
	Init()
	Start()
}

type server struct {
	listener Listener
	workers  Workers
	ipasser_ chan task
	opasser_ chan task
	pgconfig *cimess
}

func (tar *server) Init() {
	tar.listener.Init(tar.ipasser_, tar.opasser_)
	tar.workers.Init(tar.pgconfig, tar.ipasser_, tar.opasser_)
}

func (tar *server) Start() {
	tar.workers.Start()
	tar.listener.Start()
}

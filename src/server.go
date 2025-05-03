package main

import "time"

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
	tar.pgconfig = &cimess{
		database_:  "postgres",
		dominator_: "postgres",
		host_:      "localhost",
		password_:  "123456",
	}
	tar.listener.Init(tar.ipasser_, tar.opasser_)
	tar.workers.Init(tar.pgconfig, tar.ipasser_, tar.opasser_)
}

func (tar *server) Start() {
	go tar.workers.Start()
	go tar.listener.Start()
	for {
		time.Sleep(time.Second)
	}
}

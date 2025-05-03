package main

import (
	"time"
)

type Server interface {
	Init()
	Start()
}

type server struct {
	listener listener
	manager  manager
	iomap_   *iomap
	signal_  chan string
	pgconfig *cimess
}

func (tar *server) Init() {
	tar.pgconfig = &cimess{
		database_:  "postgres",
		dominator_: "postgres",
		host_:      "localhost",
		password_:  "123456",
	}
	tar.listener.Init(tar.signal_, tar.iomap_)
	tar.manager.Init(tar.signal_, tar.iomap_, tar.pgconfig)
}

func (tar *server) Start() {
	go tar.manager.Start()
	go tar.listener.Start()
	for {
		time.Sleep(time.Second)
	}
}

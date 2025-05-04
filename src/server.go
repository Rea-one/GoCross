package gocross

import (
	"time"
)

type ActServer interface {
	Init()
	Start()
}

type Server struct {
	listener listener
	manager  manager
	iomap_   *iomap
	signal_  chan string
	pgconfig *cimess
	lsconfig *string
}

func (tar *Server) Init() {
	tar.pgconfig = &cimess{
		database_:  "postgres",
		dominator_: "postgres",
		host_:      "localhost",
		password_:  "12",
	}
	config := "127.0.0.1:25054"
	tar.lsconfig = &config
	tar.iomap_ = new(iomap)
	tar.iomap_.Init()
	tar.listener.Init(tar.signal_, tar.iomap_, tar.lsconfig)
	tar.manager.Init(tar.signal_, tar.iomap_, tar.pgconfig)
}

func (tar *Server) Start() {
	go tar.manager.Start()
	go tar.listener.Start()
	for {
		time.Sleep(time.Second)
	}
}

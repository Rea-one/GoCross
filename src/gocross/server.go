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
	checker  *Checker
	signal_  chan string
	config   *cimess
}

func (tar *Server) Init() {
	tar.config = &cimess{
		database_:  "GoCross",
		dominator_: "postgres",
		password_:  "123456",
		pg_host_:   "127.0.0.1:5432",
		mn_host_:   "127.0.0.1:25059",
		host_:      "127.0.0.1:25054",
	}
	tar.signal_ = make(chan string, 4)
	tar.checker = new(Checker)
	tar.checker.Init()
	tar.listener.Init(tar.signal_, tar.checker.iom_, tar.config)
	tar.manager.Init(tar.signal_, tar.checker, tar.config)
}

func (tar *Server) Start() {
	go tar.manager.Start()
	go tar.listener.Start()
	for {
		time.Sleep(time.Second)
	}
}

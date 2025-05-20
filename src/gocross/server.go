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
		// 这里使用的是PostgreSQL数据库
		database_: "GoCross",
		// 数据库管理员
		dominator_: "postgres",
		// 数据库管理员密码
		password_: "123456",
		// 数据库主机地址
		pg_host_: "127.0.0.1:5432",
		// minio地址 目前其实没有使用
		mn_host_: "127.0.0.1:25059",
		// 监听地址
		host_: "0.0.0.0:25054",
	}
	// 交互通知
	tar.signal_ = make(chan string, 4)
	// ID映射表 用户ID -> IP
	tar.checker = new(Checker)
	tar.checker.Init()
	tar.listener.Init(tar.signal_, tar.checker.iom_, tar.config)
	tar.manager.Init(tar.signal_, tar.checker, tar.config)
}

func (tar *Server) Start() {
	// manager和listener的行为不同，
	// listener会根据监听数量即时创建receiver
	// worker会根据配置文件创建worker，当接收到来自listener的信号时，worker会调度worker完成任务
	go tar.manager.Start()
	go tar.listener.Start()
	for {
		time.Sleep(time.Second)
	}
}

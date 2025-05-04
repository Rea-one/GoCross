package gocross

import (
	"net"
	"time"
)

type Receiver interface {
	Init(int, net.Conn, chan task, chan task)
	Start()
	Stop()
	GetIP() string
}

type receiver struct {
	id_      int
	conn_    net.Conn
	stop_    bool
	release_ chan int
	ipasser_ chan task
	opasser_ chan task
}

func (tar *receiver) Init(id int, conn net.Conn,
	ip chan task, op chan task) {
	tar.conn_ = conn
	tar.ipasser_ = ip
	tar.opasser_ = op
}

func (tar *receiver) Start() {
	go tar.write()
	tar.read()
	if tar.stop_ {
		tar.release_ <- tar.id_
	}
}

func (tar *receiver) read() {
	buf := make([]byte, 1024)
	for {
		if tar.stop_ {
			break
		}
		n, err := tar.conn_.Read(buf)
		if err != nil {
			tar.conn_.Close()
			break
		}
		data := make([]byte, n)
		id := tar.GetIP()

		new_task := task{
			ID:       id,
			Query:    string(data),
			Result:   "",
			Deadline: time.Now().Add(time.Second * 10),
		}
		tar.ipasser_ <- new_task
	}
}

func (tar *receiver) write() {
	for {
		if tar.stop_ {
			break
		}
		task := <-tar.opasser_
		tar.conn_.Write([]byte(task.Result))
	}
}

func (tar *receiver) Stop() {
	tar.stop_ = true
}

func (tar *receiver) GetIP() string {
	return tar.conn_.RemoteAddr().String()
}

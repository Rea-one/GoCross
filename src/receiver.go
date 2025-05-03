package main

import (
	"fmt"
	"net"
	"time"
)

type Receiver interface {
	Init(net.Conn, chan task, chan task)
	Start()
}

type receiver struct {
	conn_    net.Conn
	record   string
	ipasser_ chan task
	opasser_ chan task
}

func (tar *receiver) Init(conn net.Conn, ip chan task, op chan task) {
	tar.conn_ = conn
	tar.ipasser_ = ip
	tar.opasser_ = op
}

func (tar *receiver) Start() {
	buf := make([]byte, 1024)
	for {
		n, err := tar.conn_.Read(buf)
		if err != nil {
			tar.conn_.Close()
			break
		}
		data := make([]byte, n)
		id := tar.conn_.RemoteAddr().String()

		new_task := task{
			ID:       id,
			Query:    string(data),
			Result:   "",
			Deadline: time.Now().Add(time.Second * 10),
		}

		select {
		case tar.ipasser_ <- new_task:
		default:
			// 可选：处理通道满的情况
			fmt.Println("通信器已满，请稍后再试")
		}
	}
}

func (tar *receiver) read() {
	buf := make([]byte, 1024)
	for {
		n, err := tar.conn_.Read(buf)
		if err != nil {
			tar.conn_.Close()
			break
		}
		data := make([]byte, n)
		id := tar.conn_.RemoteAddr().String()

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
		task := <-tar.opasser_
		tar.conn_.Write([]byte(task.Result))
	}
}

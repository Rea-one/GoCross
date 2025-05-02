package main

import (
	"fmt"
	"net"
)

type Receiver interface {
	Init(net.Conn)
	Start()
}

type receiver struct {
	conn_  net.Conn
	record string
	mess_  chan *task
}

func (tar *receiver) Init(conn net.Conn) {
	tar.conn_ = conn
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
		copy(data, buf[:n])
		select {
		case tar.mess_ <- &task{data: data}:
		default:
			// 可选：处理通道满的情况
			fmt.Println("message channel full")
		}
	}
}

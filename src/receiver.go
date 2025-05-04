package gocross

import (
	"net"
	"time"
)

type Receiver interface {
	Init(int, net.Conn, chan int, chan task, chan task)
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
	release chan int, ip chan task, op chan task) {
	tar.conn_ = conn
	tar.id_ = id
	tar.stop_ = false
	tar.release_ = release
	tar.ipasser_ = ip
	tar.opasser_ = op
}

func (tar *receiver) Start() {
	go tar.write()
	go tar.read()
}

func (tar *receiver) read() {
	buf := make([]byte, 1024)
	for !tar.stop_ {
		n, err := tar.conn_.Read(buf)
		if err != nil {
			tar.conn_.Close()
			break
		}
		data := make([]byte, n)
		copy(data, buf[:n])
		id := tar.GetIP()

		new_task := task{
			ID:       id,
			Query:    "",
			Result:   "",
			Deadline: time.Now().Add(time.Second * 10),
		}
		mess := string(data)
		if mess == "nomore" {
			new_task.Result = mess
		} else {
			new_task.Query = mess
		}
		tar.ipasser_ <- new_task
	}
}

func (tar *receiver) write() {
	for !tar.stop_ {
		select {
		case task := <-tar.opasser_:
			if task.GetResult() == "nomore" {
				tar.conn_.Close()
				tar.release_ <- tar.id_
				break
			}
			tar.conn_.Write([]byte(task.GetResult()))
		default:
			time.Sleep(time.Millisecond * 300)
		}
	}
}

func (tar *receiver) Stop() {
	tar.stop_ = true
}

func (tar *receiver) GetIP() string {
	return tar.conn_.RemoteAddr().String()
}

package main

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Worker interface {
	Start()
	Stop()
	Wait()
	Init(*pgxpool.Conn)
}

type worker struct {
	stop_     chan bool
	messages_ chan []string
	link_     *pgxpool.Conn
}

func (tar *worker) Start() {
	go func() {
		select {
		case <-tar.stop_:
			return
		default:
			select {
			case message := <-tar.messages_:
				tar.link_.QueryRow(context.Background(), message[0], message[1])
			default:
				time.Sleep(time.Second)
			}
		}
	}()
}

func (tar *worker) Stop() {
	tar.stop_ <- true
}

func (tar *worker) Wait() {

}

func (tar *worker) Init(link *pgxpool.Conn) {
	tar.stop_ = make(chan bool)
	tar.messages_ = make(chan []string)
	tar.link_ = link
}

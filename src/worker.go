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
	Init(*pgxpool.Conn, *messque)
}

type worker struct {
	stop_     chan bool
	messages_ *messque
	link_     *pgxpool.Conn
}

func (tar *worker) Start() {
	go func() {
		select {
		case <-tar.stop_:
			return
		default:
			for {
				if tar.messages_ != nil {
					message := tar.messages_.Read()
					if message != nil {
						tar.link_.Query(context.Background(), *message)
					}
					if tar.messages_.Inish() {
						tar.messages_.OClose()
						tar.messages_ = nil
					}
				} else {
					time.Sleep(time.Millisecond * 100)
				}

			}
		}
	}()
}

func (tar *worker) Stop() {
	tar.stop_ <- true
}

func (tar *worker) Wait() {

}

func (tar *worker) Init(link *pgxpool.Conn, messages *messque) {
	tar.stop_ = make(chan bool)
	tar.messages_ = messages
	tar.link_ = link
}

package main

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Worker interface {
	Start()
	Stop()
	Wait()
	Init(*pgxpool.Conn, *messque)
}

type worker struct {
	stop_   chan bool
	task_   <-chan *task
	output_ chan<- *task
	link_   *pgxpool.Conn
}

func (tar *worker) Start() {
	go func() {
		for {
			select {
			case tar.stop_ <- true:
				return
			default:
				select {
				case task := <-tar.task_:
					rows, err := tar.link_.Query(context.Background(), task.Query)
					if err != nil {
						task.Result = err.Error()
					} else {
						task.Result = processRows(rows)
					}
				case <-time.After(time.Second * 5):
					continue
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

func (tar *worker) Init(link *pgxpool.Conn, t chan *task) {
	tar.stop_ = make(chan bool)
	tar.link_ = link
	tar.task_ = t
}

func processRows(rows pgx.Rows) string {
	var result string
	for rows.Next() {
		var row string
		rows.Scan(&row)
		result += row
	}
	return result
}

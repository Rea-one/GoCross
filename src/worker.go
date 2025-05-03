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
	Action(*task)
	Change(chan task, chan task)
	Init(int, *pgxpool.Conn, chan task, chan task, chan int)
}

type worker struct {
	id_      int
	stop_    bool
	input_   <-chan task
	output_  chan<- task
	link_    *pgxpool.Conn
	release_ chan<- int
}

func (tar *worker) Start() {
	go func() {
		for {
			if !tar.stop_ {
				select {
				case task := <-tar.input_:
					if task.Result == "close" {
						tar.Stop()
						break
					}
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
	tar.stop_ = true
	defer tar.link_.Release()
}

func (tar *worker) Wait() {

}

func (tar *worker) Init(id int, link *pgxpool.Conn,
	it chan task, ot chan task, rsl chan int) {
	tar.id_ = id
	tar.stop_ = false
	tar.link_ = link
	tar.input_ = it
	tar.output_ = ot
	tar.release_ = rsl
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

func (tar *worker) Change(it chan task, ot chan task) {
	tar.input_ = it
	tar.output_ = ot
}

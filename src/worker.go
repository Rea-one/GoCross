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
	Init(int, *pgxpool.Conn, chan task, chan task, chan int)
}

type worker struct {
	id_      int
	stop_    chan bool
	input_   <-chan task
	output_  chan<- task
	link_    *pgxpool.Conn
	pick_up_ chan<- int
}

func (tar *worker) Start() {
	go func() {
		for {
			select {
			case tar.stop_ <- true:
				return
			default:
				select {
				case task := <-tar.input_:
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

func (tar *worker) Init(id int, link *pgxpool.Conn,
	it chan task, ot chan task, pick chan int) {
	tar.id_ = id
	tar.stop_ = make(chan bool)
	tar.link_ = link
	tar.input_ = it
	tar.output_ = ot
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

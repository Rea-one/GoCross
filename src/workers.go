package main

import (
	"context"
	"runtime"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Workers interface {
	Start()
	Stop()
	Wait()
	Init(*cimess, chan task, chan task)
}

type workers struct {
	rec_       *pgxpool.Config
	conn_pool_ *pgxpool.Pool
	group_     *mList[worker]
	pick_up_   chan int
	omess_     chan task
	imess_     chan task
	ipasser_   chan task
	opasser_   chan task
}

func (tar *workers) Start() {
	var id int = 0
	w := tar.group_.Till()
	for w != nil {
		conn, _ := tar.conn_pool_.Acquire(context.Background())
		w.Get().Init(id, conn, tar.imess_, tar.omess_, tar.pick_up_)
		w.Get().Start()
		w = w.Till()
		id++
	}

	for {
		select {
		case task := <-tar.ipasser_:

		default:
			time.Sleep(time.Millisecond * 100)
		}
		select {
		case task := <-tar.omess_:

		default:
			time.Sleep(time.Millisecond * 100)
		}
	}
}

func (tar *workers) Stop() {
	for _, w := range tar.group_ {
		w.Stop()
	}
}

func (tar *workers) Wait() {
	for _, w := range tar.group_ {
		w.Wait()
	}
}

func (tar *workers) Init(mess *cimess, ip chan task, op chan task) {
	config, _ := pgxpool.ParseConfig(mess.String())

	cpuNum := int32(runtime.NumCPU())

	config.MaxConns = cpuNum * 4 // 动态调整
	config.MinConns = cpuNum * 2

	tar.imess_ = make(chan task)
	tar.omess_ = make(chan task)

	tar.ipasser_ = ip
	tar.opasser_ = op
}

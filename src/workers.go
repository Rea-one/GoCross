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
	group_     []*worker
	pick_up_   chan int
	omess_     chan task
	imess_     chan task
	ipasser_   chan task
	opasser_   chan task
}

func (tar *workers) Start() {
	for id := range tar.group_ {
		conn, _ := tar.conn_pool_.Acquire(context.Background())
		tar.group_[id].Init(id, conn, tar.imess_, tar.omess_, tar.pick_up_)
		tar.group_[id].Start()
	}

	go func() {
		for {
			select {
			case task := <-tar.omess_:
				tar.opasser_ <- task
			default:
				time.Sleep(time.Millisecond * 100)
			}
		}
	}()
	for {
		select {
		case task := <-tar.ipasser_:
			tar.imess_ <- task
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

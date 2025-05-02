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
	Init(cimess, chan *mess)
}

type workers struct {
	rec_       *pgxpool.Config
	conn_pool_ *pgxpool.Pool
	group_     [4]*worker
	pmess_     []chan *task
	mess_      chan *task
}

func (tar *workers) Start() {
	for _, w := range tar.group_ {
		conn, _ := tar.conn_pool_.Acquire(context.Background())
		cur_task := make(chan *task)
		tar.pmess_ = append(tar.pmess_, cur_task)
		w.Init(conn, cur_task)
		w.Start()
	}
	var order int = 0
	for {
		select {
		case tasks := <-tar.mess_:
			if order < len(tar.pmess_) {
				tar.pmess_[order] <- tasks
				order = (order + 1) % len(tar.pmess_)
			}
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

func (tar *workers) Init(mess cimess, m chan chan *mess) {
	config, _ := pgxpool.ParseConfig(mess.String())

	cpuNum := int32(runtime.NumCPU())

	config.MaxConns = cpuNum * 4 // 动态调整
	config.MinConns = cpuNum * 2

	tar.mess_ = m
}

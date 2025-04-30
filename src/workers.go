package main

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Workers interface {
	Start()
	Stop()
	Wait()
	Init(cimess)
}

type workers struct {
	rec_       *pgxpool.Config
	conn_pool_ *pgxpool.Pool
	group_     []worker
	mess_      *demess
}

func (tar *workers) Start() {
	var conn *pgxpool.Conn
	for _, w := range tar.group_ {
		conn, _ = tar.conn_pool_.Acquire(context.Background())
		w.Init(conn)
		go w.Start()
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

func (tar *workers) Init(mess cimess) {
	var err error
	tar.rec_, err = pgxpool.ParseConfig(mess.String())
	if err != nil {
		os.Exit(1)
	}

	tar.rec_.MaxConns = 4
	tar.rec_.MinConns = 3

	tar.conn_pool_, err = pgxpool.NewWithConfig(context.Background(), tar.rec_)
	if err != nil {
		os.Exit(1)
	}
}

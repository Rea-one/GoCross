package main

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
)

type Conn interface {
	Init(cimess)
}

type CIMess interface {
	String()
}

type PGconn struct {
	conn *pgx.Conn
}

type cimess struct {
	dominator_ string
	password_  string
	host_      string
	database_  string
}

func (tar *PGconn) Init(mess cimess) {
	var err error
	var target string = "postgres://" + mess.dominator_ + ":" +
		mess.password_ + "@" + mess.host_ + "/" + mess.database_
	tar.conn, err = pgx.Connect(context.Background(), target)

	if err != nil {
		os.Exit(1)
	}
}

func (tar *cimess) String() string {
	return tar.dominator_ + ":" + tar.password_ + "@" + tar.host_ + "/" + tar.database_
}

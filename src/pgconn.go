package main

import (
	"github.com/jackc/pgx"
)

type PGconn struct {
	conn pgx.Conn
}

type CImess struct {
	dominator string
	password  string
	host      string
	database  string
}

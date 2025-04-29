package main

type worker struct {
	stop_     chan bool
	messages_ chan []string
	link_     chan *PGconn
}

package main

type Workers interface {
	Start()
	Stop()
	Wait()
	Init(CImess)
}

type workers struct {
	rec   CImess
	conn  PGconn
	group []worker
}

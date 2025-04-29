package main

type Workers interface {
	Start()
	Stop()
	Wait()
	Init(host string, dom string, database string, password string)
}

type workers struct {
	host     string
	dom      string
	database string
	password string
}

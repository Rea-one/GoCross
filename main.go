package main

import (
	"GoCross/src/gocross"
)

func main() {
	var server gocross.Server
	server.Init()
	server.Start()
}

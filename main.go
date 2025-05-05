package main

import (
	"now/gocross"
)

func main() {
	var server gocross.Server
	server.Init()
	server.Start()
}

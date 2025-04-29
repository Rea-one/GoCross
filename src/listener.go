package main

import "google.golang.org/grpc"

type Listener struct {
	grpc.Server
}

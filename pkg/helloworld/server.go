package helloworld

import (
	"context"
	"log"
)

type HelloWorldServer struct {
	UnimplementedGreeterServer
}

func (h *HelloWorldServer) SayHello(ctx context.Context, req *HelloRequest) (*HelloReply, error) {
	log.Printf("Received: %v", req.GetName())

	return &HelloReply{
		Message: "Hello " + req.GetName(),
	}, nil
}

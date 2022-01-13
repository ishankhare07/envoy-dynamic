package helloworld

import (
	"context"
	"log"

	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
)

type HelloWorldServer struct {
	snapshot cache.SnapshotCache
	UnimplementedGreeterServer
}

func (h *HelloWorldServer) SayHello(ctx context.Context, req *HelloRequest) (*HelloReply, error) {
	log.Printf("Received: %v", req.GetName())

	h.snapshot.ClearSnapshot(req.GetName())

	return &HelloReply{
		Message: "Hello " + req.GetName(),
	}, nil
}

func NewHelloWorldServer(initialSnapshot cache.SnapshotCache) *HelloWorldServer {
	return &HelloWorldServer{
		snapshot: initialSnapshot,
	}
}

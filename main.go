package main

import (
	"context"

	aggregateGRPCServer "github.com/ishankhare07/envoy-dynamic/pkg/server"
)

func main() {
	ctx := context.Background()
	aggregateGRPCServer.RunServer(ctx, 18000)
}

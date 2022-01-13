package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"

	clusterservice "github.com/envoyproxy/go-control-plane/envoy/service/cluster/v3"
	discoverygrpc "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	endpointservice "github.com/envoyproxy/go-control-plane/envoy/service/endpoint/v3"
	listenerservice "github.com/envoyproxy/go-control-plane/envoy/service/listener/v3"
	routeservice "github.com/envoyproxy/go-control-plane/envoy/service/route/v3"
	runtimeservice "github.com/envoyproxy/go-control-plane/envoy/service/runtime/v3"
	secretservice "github.com/envoyproxy/go-control-plane/envoy/service/secret/v3"

	helloworldservice "github.com/ishankhare07/envoy-dynamic/pkg/helloworld"

	generator "github.com/ishankhare07/envoy-dynamic/pkg/snapshot"

	"github.com/ishankhare07/envoy-dynamic/pkg/logger"

	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"github.com/envoyproxy/go-control-plane/pkg/test/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var log *logger.Logger
var snapshotCache cache.SnapshotCache

func init() {
	log = &logger.Logger{Debug: true}
}

func registerServer(grpcServer *grpc.Server, server server.Server) {
	// register services
	discoverygrpc.RegisterAggregatedDiscoveryServiceServer(grpcServer, server)
	endpointservice.RegisterEndpointDiscoveryServiceServer(grpcServer, server)
	clusterservice.RegisterClusterDiscoveryServiceServer(grpcServer, server)
	routeservice.RegisterRouteDiscoveryServiceServer(grpcServer, server)
	listenerservice.RegisterListenerDiscoveryServiceServer(grpcServer, server)
	secretservice.RegisterSecretDiscoveryServiceServer(grpcServer, server)
	runtimeservice.RegisterRuntimeDiscoveryServiceServer(grpcServer, server)
	// helloworldservice.RegisterGreeterServer(grpcServer, &helloworldservice.HelloWorldServer{})
}

func RunServer(ctx context.Context, port uint) {
	grpcServer := grpc.NewServer()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Errorf("%v", err)
	}

	snapshotCache = cache.NewSnapshotCache(false, cache.IDHash{}, log)

	snapshot := generator.GenerateSnapshot("listener_0", "local_route", "example_cluster", "www.envoyproxy.io", 10000, 80)
	if err := snapshot.Consistent(); err != nil {
		log.Errorf("snapshot inconsistency: %+v\n%+v", snapshot, err)
		os.Exit(1)
	}

	t, _ := json.MarshalIndent(snapshot, "", "  ")
	log.Debugf("will serve snapshot %s", t)

	// Add the snapshot to the cache
	if err := snapshotCache.SetSnapshot("test-id", snapshot); err != nil {
		log.Errorf("snapshot error %q for %+v", err, snapshot)
		os.Exit(1)
	}

	cb := &test.Callbacks{Debug: log.Debug}
	srv := server.NewServer(ctx, snapshotCache, cb)

	registerServer(grpcServer, srv)
	helloworldservice.RegisterGreeterServer(grpcServer, helloworldservice.NewHelloWorldServer(snapshotCache))
	reflection.Register(grpcServer)

	log.Infof("management server listening on port %d\n", port)
	if err = grpcServer.Serve(lis); err != nil {
		log.Errorf("%v", err)
		panic(err)
	}
}

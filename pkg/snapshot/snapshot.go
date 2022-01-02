package snapshot

import (
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	hcm "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	"github.com/golang/protobuf/ptypes"
)

func makeHTTPConnManager(routeName string) *hcm.HttpConnectionManager {
	return &hcm.HttpConnectionManager{
		CodecType:  hcm.HttpConnectionManager_AUTO,
		StatPrefix: "http",
		// AccessLog: []*accesslog.AccessLog{
		// 	{
		// 		Name:       "envoy.access_loggers.stdout",
		// 		ConfigType: &accesslog.AccessLog_TypedConfig{},
		// 	},
		// },
		RouteSpecifier: &hcm.HttpConnectionManager_Rds{
			Rds: &hcm.Rds{
				ConfigSource:    makeConfigSource(),
				RouteConfigName: routeName,
			},
		},
		HttpFilters: []*hcm.HttpFilter{
			{
				Name: wellknown.Router,
			},
		},
	}
}

func makeConfigSource() *core.ConfigSource {
	source := &core.ConfigSource{}
	source.ResourceApiVersion = resource.DefaultAPIVersion
	source.ConfigSourceSpecifier = &core.ConfigSource_ApiConfigSource{
		ApiConfigSource: &core.ApiConfigSource{
			TransportApiVersion:       resource.DefaultAPIVersion,
			ApiType:                   core.ApiConfigSource_GRPC,
			SetNodeOnFirstMessageOnly: true,
			GrpcServices: []*core.GrpcService{
				{
					TargetSpecifier: &core.GrpcService_EnvoyGrpc_{
						EnvoyGrpc: &core.GrpcService_EnvoyGrpc{
							ClusterName: "xds_cluster",
						},
					},
				},
			},
		},
	}

	return source
}

func makeListener(listenerName, routeName string, listenerPort uint32) *listener.Listener {
	manager := makeHTTPConnManager(routeName)
	typedConfig, err := ptypes.MarshalAny(manager)
	if err != nil {
		panic(err)
	}

	return &listener.Listener{
		Name: listenerName,
		Address: &core.Address{
			Address: &core.Address_SocketAddress{
				SocketAddress: &core.SocketAddress{
					Protocol: core.SocketAddress_TCP,
					Address:  "0.0.0.0",
					PortSpecifier: &core.SocketAddress_PortValue{
						PortValue: listenerPort,
					},
				},
			},
		},
		FilterChains: []*listener.FilterChain{
			{
				Filters: []*listener.Filter{
					{
						Name: wellknown.HTTPConnectionManager,
						ConfigType: &listener.Filter_TypedConfig{
							TypedConfig: typedConfig,
						},
					},
				},
			},
		},
	}
}

func GenerateSnapshot(listenerName, routeName, clusterName, upstreamHost string, listenerPort, upstreamPort uint32) cache.Snapshot {
	return cache.NewSnapshot("1",
		// endpoints
		[]types.Resource{},
		// clusters
		[]types.Resource{makeCluster(clusterName, upstreamHost, upstreamPort)},
		// routes
		[]types.Resource{makeRoute(routeName, clusterName, upstreamHost)},
		// listeners
		[]types.Resource{makeListener(listenerName, routeName, listenerPort)},
		// runtimes
		[]types.Resource{},
		// secrets
		[]types.Resource{},
	)
}

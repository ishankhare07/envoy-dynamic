package snapshot

import (
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
)

func makeRoute(routeName, clusterName, upstreamHost string) *route.RouteConfiguration {
	return &route.RouteConfiguration{
		Name: routeName,
		VirtualHosts: []*route.VirtualHost{
			{
				Name:    "local_service",
				Domains: []string{"*"},
				Routes: []*route.Route{
					{
						Match: &route.RouteMatch{
							PathSpecifier: &route.RouteMatch_Prefix{
								Prefix: "/",
							},
						},
						Action: &route.Route_Route{
							Route: &route.RouteAction{
								ClusterSpecifier: &route.RouteAction_Cluster{
									Cluster: clusterName,
								},
								HostRewriteSpecifier: &route.RouteAction_HostRewriteLiteral{
									HostRewriteLiteral: upstreamHost,
								},
							},
						},
					},
				},
			},
		},
	}
}

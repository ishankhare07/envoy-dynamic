package snapshot

import (
	"time"

	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	"github.com/golang/protobuf/ptypes"
)

func makeCluster(clusterName, upstreamHost string, upstreamPort uint32) *cluster.Cluster {
	return &cluster.Cluster{
		Name:           clusterName,
		ConnectTimeout: ptypes.DurationProto(5 * time.Second),
		ClusterDiscoveryType: &cluster.Cluster_Type{
			Type: cluster.Cluster_LOGICAL_DNS,
		},
		LbPolicy:        cluster.Cluster_ROUND_ROBIN,
		LoadAssignment:  makeEndpoint(clusterName, upstreamHost, upstreamPort),
		DnsLookupFamily: cluster.Cluster_V4_ONLY,
	}
}

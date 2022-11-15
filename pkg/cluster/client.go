package cluster

import (
	"context"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vnks"
)

type Client interface {
	ClustersUuidGet(ctx context.Context, uuid *string) (*vnks.ClusterRes, error)
	ClustersUuidKubeconfigGet(ctx context.Context, uuid *string) (*vnks.KubeconfigRes, error)
}

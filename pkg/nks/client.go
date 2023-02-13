package nks

import (
	"context"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vnks"
	"github.com/pkg/errors"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func (m Manager) GetCluster() (*vnks.Cluster, error) {
	ctx := context.Background()
	cluster, err := m.clusterClient.ClustersUuidGet(ctx, &m.clusterUuid)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get kubernetes cluster from api")
	}
	return cluster.Cluster, nil
}

func (m Manager) GetKubeconfig() (*clientcmdapi.Config, error) {
	ctx := context.Background()
	cluster, err := m.GetCluster()

	if *cluster.Status == "CREATING" {
		return nil, errors.New("kubernetes cluster is not running")
	}

	resp, err := m.clusterClient.ClustersUuidKubeconfigGet(ctx, &m.clusterUuid)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get kubeconfig from api")
	}

	cfg, err := clientcmd.Load([]byte(*resp.Kubeconfig))
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert string to kubeconfig")
	}

	return cfg, nil
}

package nks

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vnks"
	"github.com/NaverCloudPlatform/ncp-iam-authenticator/pkg/cluster"
)

type Manager struct {
	clusterClient cluster.Client
	clusterUuid   string
	region        string
}

func NewManager(clusterUuid string, region string, apiKey *ncloud.APIKey) *Manager {
	return &Manager{
		vnks.NewAPIClient(vnks.NewConfiguration(region, apiKey)).V2Api, clusterUuid, region,
	}
}

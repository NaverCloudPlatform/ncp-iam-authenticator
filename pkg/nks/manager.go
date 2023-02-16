package nks

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vnks"
)

type Manager struct {
	clusterClient Client
	clusterUuid   string
	region        string
}

func NewManager(clusterUuid string, region string, apiKey *ncloud.APIKey) *Manager {
	return &Manager{
		vnks.NewAPIClient(vnks.NewConfiguration(region, apiKey)).V2Api, clusterUuid, region,
	}
}

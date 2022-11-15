package mocks

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vnks"
	"github.com/NaverCloudPlatform/ncp-iam-authenticator/pkg/utils"
)

var CreatingClusterRes = vnks.ClusterRes{
	Cluster: &vnks.Cluster{
		Uuid:              utils.ToPointer("85782929-3DC0-4DAC-9DF9-C4CC5AA8TEST"),
		AcgName:           utils.ToPointer("nks-5526-87d5f"),
		Name:              utils.ToPointer("test-mock-data"),
		Capacity:          utils.ToPointer("vCPU 2EA, Memory 8GB"),
		ClusterType:       utils.ToPointer("SVR.VNKS.STAND.C002.M008.NET.SSD.B050.G002"),
		NodeCount:         utils.ToPointer(int32(1)),
		NodeMaxCount:      utils.ToPointer(int32(1)),
		CpuCount:          utils.ToPointer(int32(2)),
		MemorySize:        utils.ToPointer(int32(8)),
		CreatedAt:         utils.ToPointer("2022-10-31T05:19:27.000Z"),
		Endpoint:          utils.ToPointer("https://85782929-3DC0-4DAC-9DF9-C4CC5AA8TEST.kr.vnks.ntruss.com"),
		K8sVersion:        utils.ToPointer("1.23.9-nks.1"),
		RegionCode:        utils.ToPointer("KR"),
		Status:            utils.ToPointer("CREATING"),
		KubeNetworkPlugin: utils.ToPointer("cilium"),
		SubnetLbName:      utils.ToPointer("test-lb | KR-2 | 192.168.0.0/24 | Private"),
		SubnetLbNo:        utils.ToPointer(int32(1)),
		LbPublicSubnetNo:  nil,
		SubnetName:        utils.ToPointer("test"),
		SubnetNoList:      []*int32{utils.ToPointer(int32(1))},
		UpdatedAt:         utils.ToPointer("2022-10-31T05:19:27.000Z"),
		VpcName:           utils.ToPointer("test | 192.168.0.0/16"),
		VpcNo:             utils.ToPointer(int32(1)),
		ZoneCode:          utils.ToPointer("KR-2"),
		ZoneNo:            utils.ToPointer(int32(3)),
		LoginKeyName:      utils.ToPointer("test"),
		Log:               &vnks.ClusterLogInput{Audit: utils.ToPointer(false)},
		NodePool: []*vnks.NodePoolRes{
			{
				InstanceNo:     utils.ToPointer(int32(1)),
				K8sVersion:     utils.ToPointer("1.23.9"),
				Name:           utils.ToPointer("test"),
				NodeCount:      utils.ToPointer(int32(1)),
				SubnetNoList:   []*int32{utils.ToPointer(int32(1))},
				SubnetNameList: []*string{utils.ToPointer("test | KR-2 | 192.168.0.0/24 | Private")},
				ProductCode:    utils.ToPointer("SVR.VSVR.STAND.C002.M008.NET.SSD.B050.G002"),
				Status:         utils.ToPointer("CREATING"),
				Autoscale: &vnks.AutoscaleOption{
					Enabled: utils.ToPointer(false),
					Max:     utils.ToPointer(int32(0)),
					Min:     utils.ToPointer(int32(0)),
				},
			},
		},
	},
}
var RunningClusteRes = vnks.ClusterRes{
	Cluster: &vnks.Cluster{
		Uuid:              utils.ToPointer("85782929-3DC0-4DAC-9DF9-C4CC5AA8TEST"),
		AcgName:           utils.ToPointer("nks-5526-87d5f"),
		Name:              utils.ToPointer("test-mock-data"),
		Capacity:          utils.ToPointer("vCPU 2EA, Memory 8GB"),
		ClusterType:       utils.ToPointer("SVR.VNKS.STAND.C002.M008.NET.SSD.B050.G002"),
		NodeCount:         utils.ToPointer(int32(1)),
		NodeMaxCount:      utils.ToPointer(int32(1)),
		CpuCount:          utils.ToPointer(int32(2)),
		MemorySize:        utils.ToPointer(int32(8)),
		CreatedAt:         utils.ToPointer("2022-10-31T05:19:27.000Z"),
		Endpoint:          utils.ToPointer("https://85782929-3DC0-4DAC-9DF9-C4CC5AA8TEST.kr.vnks.ntruss.com"),
		K8sVersion:        utils.ToPointer("1.23.9-nks.1"),
		RegionCode:        utils.ToPointer("KR"),
		Status:            utils.ToPointer("RUNNING"),
		KubeNetworkPlugin: utils.ToPointer("cilium"),
		SubnetLbName:      utils.ToPointer("test-lb | KR-2 | 192.168.0.0/24 | Private"),
		SubnetLbNo:        utils.ToPointer(int32(1)),
		LbPublicSubnetNo:  nil,
		SubnetName:        utils.ToPointer("test"),
		SubnetNoList:      []*int32{utils.ToPointer(int32(1))},
		UpdatedAt:         utils.ToPointer("2022-10-31T05:19:27.000Z"),
		VpcName:           utils.ToPointer("test | 192.168.0.0/16"),
		VpcNo:             utils.ToPointer(int32(1)),
		ZoneCode:          utils.ToPointer("KR-2"),
		ZoneNo:            utils.ToPointer(int32(3)),
		LoginKeyName:      utils.ToPointer("test"),
		Log:               &vnks.ClusterLogInput{Audit: utils.ToPointer(false)},
		NodePool: []*vnks.NodePoolRes{
			{
				InstanceNo:     utils.ToPointer(int32(1)),
				K8sVersion:     utils.ToPointer("1.23.9"),
				Name:           utils.ToPointer("test"),
				NodeCount:      utils.ToPointer(int32(1)),
				SubnetNoList:   []*int32{utils.ToPointer(int32(1))},
				SubnetNameList: []*string{utils.ToPointer("test | KR-2 | 192.168.0.0/24 | Private")},
				ProductCode:    utils.ToPointer("SVR.VSVR.STAND.C002.M008.NET.SSD.B050.G002"),
				Status:         utils.ToPointer("RUN"),
				Autoscale: &vnks.AutoscaleOption{
					Enabled: utils.ToPointer(false),
					Max:     utils.ToPointer(int32(0)),
					Min:     utils.ToPointer(int32(0)),
				},
			},
		},
	},
}

var KubeconfigRes = vnks.KubeconfigRes{Kubeconfig: utils.ToPointer("apiVersion: v1\nclusters:\n  - cluster:\n      certificate-authority-data: mockcadata=\n      server: https://85782929-3DC0-4DAC-9DF9-C4CC5AA8TEST.kr.vnks.ntruss.com\n    name: kubernetes\ncontexts:\n  - context:\n      cluster: kubernetes\n      user: kubernetes-admin\n    name: kubernetes-admin@kubernetes\ncurrent-context: kubernetes-admin@kubernetes\nkind: Config\npreferences: {}\n")}

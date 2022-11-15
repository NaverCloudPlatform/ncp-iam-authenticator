package kubeconfig

import (
	"context"
	"errors"
	"github.com/NaverCloudPlatform/ncp-iam-authenticator/pkg/cluster"
	"github.com/NaverCloudPlatform/ncp-iam-authenticator/pkg/cluster/mocks"
	"github.com/NaverCloudPlatform/ncp-iam-authenticator/pkg/utils"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
)

func TestManager_ApplyIamToKubeconfig(t *testing.T) {
	type fields struct {
		clusterUuid   string
		clusterClient cluster.Client
		region        string
	}
	type args struct {
		config               *KubectlConfig
		profile              string
		configFile           string
		setDefaultConfigFile bool
	}
	type want struct {
		config *KubectlConfig
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			"no profile, no credentialConfig",
			fields{
				clusterUuid:   "27E9D8AD-4076-4CE6-B6C7-1C483CCBF18A",
				clusterClient: nil,
				region:        "KRS",
			},
			args{
				config: &KubectlConfig{
					ApiVersion:     "v1",
					CurrentContext: "kubernetes-admin@kubernetes",
					Clusters: []*KubectlClusterWithName{
						{
							Name: "kubernetes",
							Cluster: KubectlCluster{
								Server:                   "kubernetes-server-url",
								CertificateAuthorityData: "kubernetes-ca-data",
							},
						},
					},
					Contexts: []*KubectlContextWithName{
						{
							Name: "kubernetes-admin@kubernetes",
							Context: KubectlContext{
								Cluster: "kubernetes",
								User:    "kubernetes-admin",
							},
						},
					},
					Kind: "Config",
				},
				profile:              "",
				configFile:           "",
				setDefaultConfigFile: true,
			},
			want{config: &KubectlConfig{
				Kind:           "Config",
				ApiVersion:     "v1",
				CurrentContext: "nks-user@kubernetes",
				Clusters: []*KubectlClusterWithName{
					{
						Name: "kubernetes",
						Cluster: KubectlCluster{
							Server:                   "kubernetes-server-url",
							CertificateAuthorityData: "kubernetes-ca-data",
						},
					},
				},
				Contexts: []*KubectlContextWithName{
					{
						Name: "nks-user@kubernetes",
						Context: KubectlContext{
							Cluster: "kubernetes",
							User:    "nks-user",
						},
					},
				},
				Users: []*KubectlUserWithName{
					{
						Name: "nks-user",
						User: KubectlUser{
							Exec: ExecConfig{
								Command: "ncp-iam-authenticator",
								Args: []string{
									"token",
									"--clusterUuid",
									"27E9D8AD-4076-4CE6-B6C7-1C483CCBF18A",
									"--region",
									"KRS",
								},
								Env:        nil,
								APIVersion: "client.authentication.k8s.io/v1beta1",
							},
						},
					},
				},
			}},
		},
		{
			"profile, no credentialConfig",
			fields{
				clusterUuid:   "27E9D8AD-4076-4CE6-B6C7-1C483CCBF18A",
				clusterClient: nil,
				region:        "KRS",
			},
			args{
				config: &KubectlConfig{
					ApiVersion:     "v1",
					CurrentContext: "kubernetes-admin@kubernetes",
					Clusters: []*KubectlClusterWithName{
						{
							Name: "kubernetes",
							Cluster: KubectlCluster{
								Server:                   "kubernetes-server-url",
								CertificateAuthorityData: "kubernetes-ca-data",
							},
						},
					},
					Contexts: []*KubectlContextWithName{
						{
							Name: "kubernetes-admin@kubernetes",
							Context: KubectlContext{
								Cluster: "kubernetes",
								User:    "kubernetes-admin",
							},
						},
					},
					Kind: "Config",
				},
				profile:              "tester-krs",
				configFile:           "",
				setDefaultConfigFile: true,
			},
			want{config: &KubectlConfig{
				Kind:           "Config",
				ApiVersion:     "v1",
				CurrentContext: "nks-user@kubernetes",
				Clusters: []*KubectlClusterWithName{
					{
						Name: "kubernetes",
						Cluster: KubectlCluster{
							Server:                   "kubernetes-server-url",
							CertificateAuthorityData: "kubernetes-ca-data",
						},
					},
				},
				Contexts: []*KubectlContextWithName{
					{
						Name: "nks-user@kubernetes",
						Context: KubectlContext{
							Cluster: "kubernetes",
							User:    "nks-user",
						},
					},
				},
				Users: []*KubectlUserWithName{
					{
						Name: "nks-user",
						User: KubectlUser{
							Exec: ExecConfig{
								Command: "ncp-iam-authenticator",
								Args: []string{
									"token",
									"--clusterUuid",
									"27E9D8AD-4076-4CE6-B6C7-1C483CCBF18A",
									"--region",
									"KRS",
									"--profile",
									"tester-krs",
								},
								Env:        nil,
								APIVersion: "client.authentication.k8s.io/v1beta1",
							},
						},
					},
				},
			}},
		},
		{
			"no profile, no credentialConfig",
			fields{
				clusterUuid:   "27E9D8AD-4076-4CE6-B6C7-1C483CCBF18A",
				clusterClient: nil,
				region:        "KRS",
			},
			args{
				config: &KubectlConfig{
					ApiVersion:     "v1",
					CurrentContext: "kubernetes-admin@kubernetes",
					Clusters: []*KubectlClusterWithName{
						{
							Name: "kubernetes",
							Cluster: KubectlCluster{
								Server:                   "kubernetes-server-url",
								CertificateAuthorityData: "kubernetes-ca-data",
							},
						},
					},
					Contexts: []*KubectlContextWithName{
						{
							Name: "kubernetes-admin@kubernetes",
							Context: KubectlContext{
								Cluster: "kubernetes",
								User:    "kubernetes-admin",
							},
						},
					},
					Kind: "Config",
				},
				profile:              "",
				configFile:           "~/test/.configure.cp",
				setDefaultConfigFile: false,
			},
			want{config: &KubectlConfig{
				Kind:           "Config",
				ApiVersion:     "v1",
				CurrentContext: "nks-user@kubernetes",
				Clusters: []*KubectlClusterWithName{
					{
						Name: "kubernetes",
						Cluster: KubectlCluster{
							Server:                   "kubernetes-server-url",
							CertificateAuthorityData: "kubernetes-ca-data",
						},
					},
				},
				Contexts: []*KubectlContextWithName{
					{
						Name: "nks-user@kubernetes",
						Context: KubectlContext{
							Cluster: "kubernetes",
							User:    "nks-user",
						},
					},
				},
				Users: []*KubectlUserWithName{
					{
						Name: "nks-user",
						User: KubectlUser{
							Exec: ExecConfig{
								Command: "ncp-iam-authenticator",
								Args: []string{
									"token",
									"--clusterUuid",
									"27E9D8AD-4076-4CE6-B6C7-1C483CCBF18A",
									"--region",
									"KRS",
									"--credentialConfig",
									"~/test/.configure.cp",
								},
								Env:        nil,
								APIVersion: "client.authentication.k8s.io/v1beta1",
							},
						},
					},
				},
			}},
		},
		{
			"profile, credentialConfig",
			fields{
				clusterUuid:   "27E9D8AD-4076-4CE6-B6C7-1C483CCBF18A",
				clusterClient: nil,
				region:        "KRS",
			},
			args{
				config: &KubectlConfig{
					ApiVersion:     "v1",
					CurrentContext: "kubernetes-admin@kubernetes",
					Clusters: []*KubectlClusterWithName{
						{
							Name: "kubernetes",
							Cluster: KubectlCluster{
								Server:                   "kubernetes-server-url",
								CertificateAuthorityData: "kubernetes-ca-data",
							},
						},
					},
					Contexts: []*KubectlContextWithName{
						{
							Name: "kubernetes-admin@kubernetes",
							Context: KubectlContext{
								Cluster: "kubernetes",
								User:    "kubernetes-admin",
							},
						},
					},
					Kind: "Config",
				},
				profile:              "tester-krs",
				configFile:           "~/test/.configure.cp",
				setDefaultConfigFile: false,
			},
			want{config: &KubectlConfig{
				Kind:           "Config",
				ApiVersion:     "v1",
				CurrentContext: "nks-user@kubernetes",
				Clusters: []*KubectlClusterWithName{
					{
						Name: "kubernetes",
						Cluster: KubectlCluster{
							Server:                   "kubernetes-server-url",
							CertificateAuthorityData: "kubernetes-ca-data",
						},
					},
				},
				Contexts: []*KubectlContextWithName{
					{
						Name: "nks-user@kubernetes",
						Context: KubectlContext{
							Cluster: "kubernetes",
							User:    "nks-user",
						},
					},
				},
				Users: []*KubectlUserWithName{
					{
						Name: "nks-user",
						User: KubectlUser{
							Exec: ExecConfig{
								Command: "ncp-iam-authenticator",
								Args: []string{
									"token",
									"--clusterUuid",
									"27E9D8AD-4076-4CE6-B6C7-1C483CCBF18A",
									"--region",
									"KRS",
									"--credentialConfig",
									"~/test/.configure.cp",
									"--profile",
									"tester-krs",
								},
								Env:        nil,
								APIVersion: "client.authentication.k8s.io/v1beta1",
							},
						},
					},
				},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Manager{
				clusterUuid:   tt.fields.clusterUuid,
				clusterClient: tt.fields.clusterClient,
				region:        tt.fields.region,
			}
			m.ApplyIamToKubeconfig(tt.args.config, tt.args.profile, tt.args.configFile, tt.args.setDefaultConfigFile)
			if !reflect.DeepEqual(tt.args.config, tt.want.config) {
				t.Error(spew.Sprintf("GetKubeconfig() got = %v\n, want %v", tt.args.config, tt.want.config))
			}
		})
	}
}

func TestManager_GetKubeconfig(t *testing.T) {
	type fields struct {
		clusterUuid   string
		clusterClient cluster.Client
		region        string
	}
	//clusterGetRunning := "cluster-get-running"
	//kubeconfigGetError := "kubeconfig-get-error"
	//kubeconfigGetSuccess := "kubeconfig-get-success"

	ctx := context.Background()
	client := mocks.NewClient(t)
	client.On("ClustersUuidGet", ctx, utils.ToPointer("cluster-get-error")).Return(nil, errors.New("cluster get failed"))
	client.On("ClustersUuidGet", ctx, utils.ToPointer("cluster-get-creating")).Return(&mocks.CreatingClusterRes, nil)
	client.On("ClustersUuidGet", ctx, mock.Anything).Return(&mocks.RunningClusteRes, nil)
	client.On("ClustersUuidKubeconfigGet", ctx, utils.ToPointer("kubeconfig-get-error")).Return(nil, errors.New("kubeconfig get error"))
	client.On("ClustersUuidKubeconfigGet", ctx, mock.Anything).Return(&mocks.KubeconfigRes, nil)

	tests := []struct {
		name    string
		fields  fields
		want    *KubectlConfig
		wantErr bool
	}{
		{
			"cluster get failed",
			fields{
				clusterUuid:   "cluster-get-error",
				clusterClient: client,
				region:        "KR",
			},
			nil,
			true,
		},
		{
			"cluster status is creating",
			fields{
				clusterUuid:   "cluster-get-creating",
				clusterClient: client,
				region:        "KR",
			},
			nil,
			true,
		},
		{
			"cluster status is running && kubeconfig get error",
			fields{
				clusterUuid:   "kubeconfig-get-error",
				clusterClient: client,
				region:        "KR",
			},
			nil,
			true,
		},
		{
			"cluster status is running && kubeconfig get",
			fields{
				clusterUuid:   "kubeconfig-get-success",
				clusterClient: client,
				region:        "KR",
			},
			&KubectlConfig{
				Kind:           "Config",
				ApiVersion:     "v1",
				CurrentContext: "kubernetes-admin@kubernetes",
				Clusters: []*KubectlClusterWithName{
					{
						Name: "kubernetes",
						Cluster: KubectlCluster{
							Server:                   "https://85782929-3DC0-4DAC-9DF9-C4CC5AA8TEST.kr.vnks.ntruss.com",
							CertificateAuthorityData: "mockcadata=",
						},
					},
				},
				Contexts: []*KubectlContextWithName{
					{
						Name: "kubernetes-admin@kubernetes",
						Context: KubectlContext{
							Cluster: "kubernetes",
							User:    "kubernetes-admin",
						},
					},
				},
				Users: nil,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Manager{
				clusterUuid:   tt.fields.clusterUuid,
				clusterClient: tt.fields.clusterClient,
				region:        tt.fields.region,
			}
			got, err := m.GetKubeconfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetKubeconfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Error(spew.Sprintf("GetKubeconfig() got = %v, want %v", got, tt.want))
			}
		})
	}
}

func TestNewManager(t *testing.T) {
	type args struct {
		clusterUuid   string
		clusterClient cluster.Client
		region        string
	}
	tests := []struct {
		name string
		args args
		want *Manager
	}{
		{
			"new manager", args{
				clusterUuid:   "KR",
				clusterClient: nil,
				region:        "KRS",
			}, &Manager{
				clusterUuid:   "KR",
				clusterClient: nil,
				region:        "KRS",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewManager(tt.args.clusterUuid, tt.args.clusterClient, tt.args.region); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewManager() = %v, want %v", got, tt.want)
			}
		})
	}
}

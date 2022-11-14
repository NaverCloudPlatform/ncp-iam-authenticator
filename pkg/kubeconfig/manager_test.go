package kubeconfig

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/davecgh/go-spew/spew"
	"reflect"
	"testing"
)

func TestManager_ApplyIamToKubeconfig(t *testing.T) {
	type fields struct {
		clusterUuid  string
		ncloudConfig *ncloud.Configuration
		region       string
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
			name: "no profile, no credentialConfig",
			fields: fields{
				clusterUuid:  "27E9D8AD-4076-4CE6-B6C7-1C483CCBF18A",
				ncloudConfig: nil,
				region:       "KRS",
			},
			args: args{
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
			want: want{config: &KubectlConfig{
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
			name: "profile, no credentialConfig",
			fields: fields{
				clusterUuid:  "27E9D8AD-4076-4CE6-B6C7-1C483CCBF18A",
				ncloudConfig: nil,
				region:       "KRS",
			},
			args: args{
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
			want: want{config: &KubectlConfig{
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
			name: "no profile, no credentialConfig",
			fields: fields{
				clusterUuid:  "27E9D8AD-4076-4CE6-B6C7-1C483CCBF18A",
				ncloudConfig: nil,
				region:       "KRS",
			},
			args: args{
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
			want: want{config: &KubectlConfig{
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
			name: "profile, credentialConfig",
			fields: fields{
				clusterUuid:  "27E9D8AD-4076-4CE6-B6C7-1C483CCBF18A",
				ncloudConfig: nil,
				region:       "KRS",
			},
			args: args{
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
			want: want{config: &KubectlConfig{
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
				clusterUuid:  tt.fields.clusterUuid,
				ncloudConfig: tt.fields.ncloudConfig,
				region:       tt.fields.region,
			}
			m.ApplyIamToKubeconfig(tt.args.config, tt.args.profile, tt.args.configFile, tt.args.setDefaultConfigFile)
			if !reflect.DeepEqual(tt.args.config, tt.want.config) {
				t.Error(spew.Sprintf("GetKubeconfig() got = %v\n, want %v", tt.args.config, tt.want.config))
			}
		})
	}
}

//func TestManager_GetKubeconfig(t *testing.T) {
//	type fields struct {
//		clusterUuid  string
//		ncloudConfig *ncloud.Configuration
//		region       string
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		want    *KubectlConfig
//		wantErr bool
//	}{
//		{},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			m := Manager{
//				clusterUuid:  tt.fields.clusterUuid,
//				ncloudConfig: tt.fields.ncloudConfig,
//				region:       tt.fields.region,
//			}
//			got, err := m.GetKubeconfig()
//			if (err != nil) != tt.wantErr {
//				t.Errorf("GetKubeconfig() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("GetKubeconfig() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

func TestNewManager(t *testing.T) {
	type args struct {
		clusterUuid  string
		ncloudConfig *ncloud.Configuration
		region       string
	}
	tests := []struct {
		name string
		args args
		want *Manager
	}{
		{
			name: "create new manager",
			args: args{
				clusterUuid: "test-1230918209380",
				ncloudConfig: &ncloud.Configuration{
					BasePath:      "test-1230918209380",
					Host:          "test-1230918209380",
					Scheme:        "test-1230918209380",
					DefaultHeader: nil,
					UserAgent:     "test-1230918209380",
					HTTPClient:    nil,
					APIKey:        nil,
					Credentials:   nil,
				},
				region: "test-1230918209380",
			},
			want: &Manager{
				clusterUuid: "test-1230918209380",
				ncloudConfig: &ncloud.Configuration{
					BasePath:      "test-1230918209380",
					Host:          "test-1230918209380",
					Scheme:        "test-1230918209380",
					DefaultHeader: nil,
					UserAgent:     "test-1230918209380",
					HTTPClient:    nil,
					APIKey:        nil,
					Credentials:   nil,
				},
				region: "test-1230918209380",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewManager(tt.args.clusterUuid, tt.args.ncloudConfig, tt.args.region); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewManager() = %v, want %v", got, tt.want)
			}
		})
	}
}

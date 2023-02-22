package utils

import (
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"reflect"
	"testing"
)

func TestIsEmptyString(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"empty string",
			args{str: " \t\n"},
			true,
		},
		{
			"not empty string",
			args{str: "kk"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEmptyString(tt.args.str); got != tt.want {
				t.Errorf("IsEmptyString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateKubeconfigDupliacted(t *testing.T) {
	type args struct {
		clusterName string
		userName    string
		contextName string
		config      *clientcmdapi.Config
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"empty config, not duplicated",
			args{
				clusterName: "cluster",
				userName:    "user",
				contextName: "context",
				config: &clientcmdapi.Config{
					Kind:           "",
					APIVersion:     "",
					Preferences:    clientcmdapi.Preferences{},
					Clusters:       map[string]*clientcmdapi.Cluster{},
					AuthInfos:      map[string]*clientcmdapi.AuthInfo{},
					Contexts:       map[string]*clientcmdapi.Context{},
					CurrentContext: "",
					Extensions:     nil,
				},
			},
			false,
		}, {
			"not duplicated",
			args{
				clusterName: "cluster",
				userName:    "user",
				contextName: "context",
				config: &clientcmdapi.Config{
					Kind:        "",
					APIVersion:  "",
					Preferences: clientcmdapi.Preferences{},
					Clusters: map[string]*clientcmdapi.Cluster{
						"cluster1": nil,
					},
					AuthInfos: map[string]*clientcmdapi.AuthInfo{
						"user1": nil,
					},
					Contexts: map[string]*clientcmdapi.Context{
						"context1": nil,
					},
					CurrentContext: "",
					Extensions:     nil,
				},
			},
			false,
		}, {
			"cluster duplicated",
			args{
				clusterName: "cluster",
				userName:    "user",
				contextName: "context",
				config: &clientcmdapi.Config{
					Kind:        "",
					APIVersion:  "",
					Preferences: clientcmdapi.Preferences{},
					Clusters: map[string]*clientcmdapi.Cluster{
						"cluster": nil,
					},
					AuthInfos: map[string]*clientcmdapi.AuthInfo{
						"user1": nil,
					},
					Contexts: map[string]*clientcmdapi.Context{
						"context1": nil,
					},
					CurrentContext: "",
					Extensions:     nil,
				},
			},
			true,
		}, {
			"user duplicated",
			args{
				clusterName: "cluster",
				userName:    "user",
				contextName: "context",
				config: &clientcmdapi.Config{
					Kind:        "",
					APIVersion:  "",
					Preferences: clientcmdapi.Preferences{},
					Clusters: map[string]*clientcmdapi.Cluster{
						"cluster1": nil,
					},
					AuthInfos: map[string]*clientcmdapi.AuthInfo{
						"user": nil,
					},
					Contexts: map[string]*clientcmdapi.Context{
						"context1": nil,
					},
					CurrentContext: "",
					Extensions:     nil,
				},
			},
			true,
		}, {
			"context duplicated",
			args{
				clusterName: "cluster",
				userName:    "user",
				contextName: "context",
				config: &clientcmdapi.Config{
					Kind:        "",
					APIVersion:  "",
					Preferences: clientcmdapi.Preferences{},
					Clusters: map[string]*clientcmdapi.Cluster{
						"cluster1": nil,
					},
					AuthInfos: map[string]*clientcmdapi.AuthInfo{
						"user1": nil,
					},
					Contexts: map[string]*clientcmdapi.Context{
						"context": nil,
					},
					CurrentContext: "",
					Extensions:     nil,
				},
			},
			true,
		}, {
			"all duplicated",
			args{
				clusterName: "cluster",
				userName:    "user",
				contextName: "context",
				config: &clientcmdapi.Config{
					Kind:        "",
					APIVersion:  "",
					Preferences: clientcmdapi.Preferences{},
					Clusters: map[string]*clientcmdapi.Cluster{
						"cluster": nil,
					},
					AuthInfos: map[string]*clientcmdapi.AuthInfo{
						"user": nil,
					},
					Contexts: map[string]*clientcmdapi.Context{
						"context": nil,
					},
					CurrentContext: "",
					Extensions:     nil,
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateKubeconfigDupliacted(tt.args.clusterName, tt.args.userName, tt.args.contextName, tt.args.config); (err != nil) != tt.wantErr {
				t.Errorf("ValidateKubeconfigDupliacted() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPrettyJsonBytes(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			"compact to pretty",
			args{b: []byte("{\"apiVersion\":\"v1\",\"clusters\":[{\"cluster\":{\"certificate-authority-data\":\"asd=\",\"server\":\"https://3f03d958-4c4f-4af3-8bb5-bfca6df1d200.kr.vnks.ntruss.com\"},\"name\":\"nks_kr_test_3f03d958-4c4f-4af3-8bb5-bfca6df1d200\"}],\"contexts\":[{\"context\":{\"cluster\":\"nks_kr_test_3f03d958-4c4f-4af3-8bb5-bfca6df1d200\",\"user\":\"nks_kr_test_3f03d958-4c4f-4af3-8bb5-bfca6df1d200\"},\"name\":\"nks_kr_test_3f03d958-4c4f-4af3-8bb5-bfca6df1d200\"}],\"current-context\":\"nks_kr_test_3f03d958-4c4f-4af3-8bb5-bfca6df1d200\",\"kind\":\"Config\",\"preferences\":{},\"users\":[{\"name\":\"nks_kr_test_3f03d958-4c4f-4af3-8bb5-bfca6df1d200\",\"user\":{\"exec\":{\"apiVersion\":\"client.authentication.k8s.io/v1beta1\",\"args\":[\"token\",\"--clusterUuid\",\"3f03d958-4c4f-4af3-8bb5-bfca6df1d200\",\"--region\",\"KR\",\"--credentialConfig\",\"/Users/user/config\",\"--profile\",\"pub\"],\"command\":\"ncp-iam-authenticator\",\"env\":null,\"provideClusterInfo\":false}}}]}")},
			[]byte("{\n    \"apiVersion\": \"v1\",\n    \"clusters\": [\n        {\n            \"cluster\": {\n                \"certificate-authority-data\": \"asd=\",\n                \"server\": \"https://3f03d958-4c4f-4af3-8bb5-bfca6df1d200.kr.vnks.ntruss.com\"\n            },\n            \"name\": \"nks_kr_test_3f03d958-4c4f-4af3-8bb5-bfca6df1d200\"\n        }\n    ],\n    \"contexts\": [\n        {\n            \"context\": {\n                \"cluster\": \"nks_kr_test_3f03d958-4c4f-4af3-8bb5-bfca6df1d200\",\n                \"user\": \"nks_kr_test_3f03d958-4c4f-4af3-8bb5-bfca6df1d200\"\n            },\n            \"name\": \"nks_kr_test_3f03d958-4c4f-4af3-8bb5-bfca6df1d200\"\n        }\n    ],\n    \"current-context\": \"nks_kr_test_3f03d958-4c4f-4af3-8bb5-bfca6df1d200\",\n    \"kind\": \"Config\",\n    \"preferences\": {},\n    \"users\": [\n        {\n            \"name\": \"nks_kr_test_3f03d958-4c4f-4af3-8bb5-bfca6df1d200\",\n            \"user\": {\n                \"exec\": {\n                    \"apiVersion\": \"client.authentication.k8s.io/v1beta1\",\n                    \"args\": [\n                        \"token\",\n                        \"--clusterUuid\",\n                        \"3f03d958-4c4f-4af3-8bb5-bfca6df1d200\",\n                        \"--region\",\n                        \"KR\",\n                        \"--credentialConfig\",\n                        \"/Users/user/config\",\n                        \"--profile\",\n                        \"pub\"\n                    ],\n                    \"command\": \"ncp-iam-authenticator\",\n                    \"env\": null,\n                    \"provideClusterInfo\": false\n                }\n            }\n        }\n    ]\n}"),
		},
		{
			"not json",
			args{b: []byte("{apiVersion: v1, kind: Config, clusters: [{cluster: {proxy-url: 'http://proxy.example.org:3128', server: 'https://k8s.example.org/k8s/clusters/c-xxyyzz'}, name: development}], users: [{name: developer}], contexts: [{context: null, name: development}]}")},
			[]byte(""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PrettyJsonBytes(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PrettyJsonBytes() = %s, want %s", got, tt.want)
			}
		})
	}
}

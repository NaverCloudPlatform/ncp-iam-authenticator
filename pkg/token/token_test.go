package token

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud/credentials"
	"reflect"
	"testing"
)

func TestNewGenerator(t *testing.T) {
	tests := []struct {
		name    string
		want    Generator
		wantErr bool
	}{
		{
			"create new generator", generator{}, false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewGenerator()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGenerator() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGenerator() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generator_FormatJSON(t *testing.T) {
	type args struct {
		token Token
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := generator{}
			got, err := g.FormatJSON(tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("FormatJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FormatJSON() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generator_Get(t *testing.T) {
	type args struct {
		credential *credentials.Credentials
		clusterId  string
		region     string
	}
	tests := []struct {
		name    string
		args    args
		want    *Token
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := generator{}
			got, err := g.Get(tt.args.credential, tt.args.clusterId, tt.args.region)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getPathWithParams(t *testing.T) {
	type args struct {
		clusterId string
		region    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getPathWithParams(tt.args.clusterId, tt.args.region); got != tt.want {
				t.Errorf("getPathWithParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getStageFromRegion(t *testing.T) {
	type args struct {
		region string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"empty string", args{region: ""}, "v1",
		},
		{
			"FKR", args{region: "FKR"}, "v1",
		},
		{
			"KR", args{region: "KR"}, "v1",
		},
		{
			"SGN", args{region: "SGN"}, "sgn-v1",
		},
		{
			"KRS", args{region: "KRS"}, "krs-v1",
		},
		{
			"JPN", args{region: "JPN"}, "jpn-v1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getStageFromRegion(tt.args.region); got != tt.want {
				t.Errorf("getStageFromRegion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_makeSignature(t *testing.T) {
	type args struct {
		method    string
		uri       string
		accessKey string
		secretKey string
		timestamp string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"make signature successfully", args{
				method:    "GET",
				uri:       "/cluster/kubeconfig",
				accessKey: "access",
				secretKey: "secret",
				timestamp: "1668407407855",
			}, "LTJli9+OKT2KvUXxiKslMfu5FIOmDN83avehOvgUFp0=",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := makeSignature(tt.args.method, tt.args.uri, tt.args.accessKey, tt.args.secretKey, tt.args.timestamp); got != tt.want {
				t.Errorf("makeSignature() = %v, want %v", got, tt.want)
			}
		})
	}
}

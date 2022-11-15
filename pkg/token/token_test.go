package token

import (
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
		{
			"generate formatJSON",
			args{token: Token{Token: "test"}},
			"{\"kind\":\"ExecCredential\",\"apiVersion\":\"client.authentication.k8s.io/v1beta1\",\"spec\":{},\"status\":{\"token\":\"test\"}}",
			false,
		},
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
		{
			"get path with params",
			args{
				clusterId: "80CDD145-453B-473F-8078-D84789A5DAD3",
				region:    "KRS",
			},
			"/iam/" + getStageFromRegion("KRS") + "/user?clusterUuid=80CDD145-453B-473F-8078-D84789A5DAD3",
		},
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
			"make signature successfully",
			args{
				method:    "GET",
				uri:       "/cluster/kubeconfig",
				accessKey: "access",
				secretKey: "secret",
				timestamp: "1668407407855",
			},
			"LTJli9+OKT2KvUXxiKslMfu5FIOmDN83avehOvgUFp0=",
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

func Test_makeToken(t *testing.T) {
	type args struct {
		timestamp string
		accessKey string
		secretKey string
		clusterId string
		region    string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"make token",
			args{
				timestamp: "1668407407855",
				accessKey: "access",
				secretKey: "secret",
				clusterId: "D11328F1-ECA9-4F1B-BA22-921F61D9C5FF",
				region:    "KRS",
			},
			"k8s-ncp-v1.eyJ0aW1lc3RhbXAiOiIxNjY4NDA3NDA3ODU1IiwiYWNjZXNzS2V5IjoiYWNjZXNzIiwic2lnbmF0dXJlIjoiZW1lYWhmVzRDbXpzeElRNWlwMGNkOXlxRVJYSitaSnNSZlFwR05tS2RYaz0iLCJwYXRoIjoiL2lhbS9rcnMtdjEvdXNlcj9jbHVzdGVyVXVpZD1EMTEzMjhGMS1FQ0E5LTRGMUItQkEyMi05MjFGNjFEOUM1RkYifQ==",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := makeToken(tt.args.timestamp, tt.args.accessKey, tt.args.secretKey, tt.args.clusterId, tt.args.region)
			if (err != nil) != tt.wantErr {
				t.Errorf("makeToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("makeToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}

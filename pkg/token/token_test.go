package token

import (
	"encoding/base64"
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud/credentials"
	"github.com/NaverCloudPlatform/ncp-iam-authenticator/pkg/constants"
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
			args{token: Token{
				Token:      "test",
				Expiration: time.Date(2026, 9, 1, 13, 4, 0, 0, time.UTC),
			}},
			"{\"kind\":\"ExecCredential\",\"apiVersion\":\"client.authentication.k8s.io/v1beta1\",\"spec\":{\"interactive\":false},\"status\":{\"expirationTimestamp\":\"2026-09-01T13:04:00Z\",\"token\":\"test\"}}",
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
		{
			"PCS01", args{region: "PCS01"}, "v1",
		},
		{
			"FCS01", args{region: "FCS01"}, "v1",
		},
		{
			"GCS01", args{region: "GCS01"}, "v1",
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

func decodeClaim(t *testing.T, token string) Claim {
	t.Helper()

	if !strings.HasPrefix(token, constants.TokenPrefix) {
		t.Fatalf("token = %v, want prefix %v", token, constants.TokenPrefix)
	}

	raw, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(token, constants.TokenPrefix))
	if err != nil {
		t.Fatalf("failed to decode token: %v", err)
	}

	var claim Claim
	if err := json.Unmarshal(raw, &claim); err != nil {
		t.Fatalf("failed to unmarshal claim: %v", err)
	}

	return claim
}

func Test_generator_Get(t *testing.T) {
	type args struct {
		accessKey string
		secretKey string
		clusterId string
		region    string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"get token",
			args{
				accessKey: "access",
				secretKey: "secret",
				clusterId: "D11328F1-ECA9-4F1B-BA22-921F61D9C5FF",
				region:    "KRS",
			},
			false,
		},
		{
			"get token with empty region",
			args{
				accessKey: "access",
				secretKey: "secret",
				clusterId: "D11328F1-ECA9-4F1B-BA22-921F61D9C5FF",
				region:    "",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := generator{}
			credential := credentials.NewValueProviderCreds(tt.args.accessKey, tt.args.secretKey)

			before := time.Now()
			got, err := g.Get(credential, tt.args.clusterId, tt.args.region)
			after := time.Now()

			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			if got.Expiration.Before(before.Add(tokenValidity)) || got.Expiration.After(after.Add(tokenValidity)) {
				t.Errorf("Get() Expiration = %v, want within [%v, %v]",
					got.Expiration, before.Add(tokenValidity), after.Add(tokenValidity))
			}

			claim := decodeClaim(t, got.Token)
			gotTimestamp, err := strconv.ParseInt(claim.TimeStamp, 10, 64)
			if err != nil {
				t.Fatalf("failed to parse claim timestamp %q: %v", claim.TimeStamp, err)
			}
			if wantTimestamp := makeTimestamp(got.Expiration.Add(-tokenValidity)); gotTimestamp != wantTimestamp {
				t.Errorf("Get() claim timestamp = %v, want %v (Expiration - tokenValidity)",
					gotTimestamp, wantTimestamp)
			}

			wantToken, err := makeToken(claim.TimeStamp, tt.args.accessKey, tt.args.secretKey, tt.args.clusterId, tt.args.region)
			if err != nil {
				t.Fatalf("makeToken() error = %v", err)
			}
			if got.Token != wantToken {
				t.Errorf("Get() Token = %v, want %v", got.Token, wantToken)
			}
		})
	}
}

package credentials

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"testing"
)

func TestConfig_Valid(t *testing.T) {
	type fields struct {
		APIKey *ncloud.APIKey
		ApiUrl string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			"access key empty",
			fields{
				APIKey: &ncloud.APIKey{
					AccessKey: " \n\t",
					SecretKey: "secret",
				},
				ApiUrl: "www.ncloud.com",
			},
			false,
		},
		{
			"secret key empty",
			fields{
				APIKey: &ncloud.APIKey{
					AccessKey: "",
					SecretKey: " \n\t",
				},
				ApiUrl: "www.ncloud.com",
			},
			false,
		},
		{
			"api url empty",
			fields{
				APIKey: &ncloud.APIKey{
					AccessKey: " ",
					SecretKey: "secret",
				},
				ApiUrl: " \n\t",
			},
			false,
		},
		{
			"valid",
			fields{
				APIKey: &ncloud.APIKey{
					AccessKey: "access",
					SecretKey: "secret",
				},
				ApiUrl: "www.ncloud.com",
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				APIKey: tt.fields.APIKey,
				ApiUrl: tt.fields.ApiUrl,
			}
			if got := c.Valid(); got != tt.want {
				t.Errorf("Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}

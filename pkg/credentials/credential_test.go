package credentials

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"gopkg.in/ini.v1"
	"reflect"
	"testing"
)

func TestConfig_NewCredentialFromCommandLine(t *testing.T) {
	type fields struct {
		APIKey *ncloud.APIKey
		ApiUrl string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				APIKey: tt.fields.APIKey,
				ApiUrl: tt.fields.ApiUrl,
			}
			if err := c.NewCredentialFromCommandLine(); (err != nil) != tt.wantErr {
				t.Errorf("NewCredentialFromCommandLine() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

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
		// TODO: Add test cases.
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

func TestConfig_WriteCredentialToFile(t *testing.T) {
	type fields struct {
		APIKey *ncloud.APIKey
		ApiUrl string
	}
	type args struct {
		configPath string
		profile    string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				APIKey: tt.fields.APIKey,
				ApiUrl: tt.fields.ApiUrl,
			}
			c.WriteCredentialToFile(tt.args.configPath, tt.args.profile)
		})
	}
}

func TestNewCredentialConfig(t *testing.T) {
	type args struct {
		configPath string
		profile    string
	}
	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCredentialConfig(tt.args.configPath, tt.args.profile)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCredentialConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCredentialConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewCredentialFromEnv(t *testing.T) {
	tests := []struct {
		name string
		want *Config
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCredentialFromEnv(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCredentialFromEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewCredentialFromFile(t *testing.T) {
	type args struct {
		configPath string
		profile    string
	}
	tests := []struct {
		name string
		args args
		want *Config
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCredentialFromFile(tt.args.configPath, tt.args.profile); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCredentialFromFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getNcloudConfigFile(t *testing.T) {
	type args struct {
		configPath string
	}
	tests := []struct {
		name    string
		args    args
		want    *ini.File
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getNcloudConfigFile(tt.args.configPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("getNcloudConfigFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getNcloudConfigFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getNcloudConfigFileSection(t *testing.T) {
	type args struct {
		file    *ini.File
		profile string
	}
	tests := []struct {
		name    string
		args    args
		want    *ini.Section
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getNcloudConfigFileSection(tt.args.file, tt.args.profile)
			if (err != nil) != tt.wantErr {
				t.Errorf("getNcloudConfigFileSection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getNcloudConfigFileSection() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getNewSection(t *testing.T) {
	type args struct {
		file    *ini.File
		profile string
	}
	tests := []struct {
		name    string
		args    args
		want    *ini.Section
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getNewSection(tt.args.file, tt.args.profile)
			if (err != nil) != tt.wantErr {
				t.Errorf("getNewSection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getNewSection() got = %v, want %v", got, tt.want)
			}
		})
	}
}

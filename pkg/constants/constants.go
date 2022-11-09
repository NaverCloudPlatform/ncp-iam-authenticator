package constants

const (
	Version = "1.0.6"

	TokenPrefix = "k8s-ncp-v1."

	ProfileEnv         = "NCLOUD_PROFILE"
	AccessKeyIdEnv     = "NCLOUD_ACCESS_KEY_ID"
	AccessKeyEnv       = "NCLOUD_ACCESS_KEY"
	SecretAccessKeyEnv = "NCLOUD_SECRET_ACCESS_KEY"
	SecretKeyEnv       = "NCLOUD_SECRET_KEY"
	ApiGwUrlEnv        = "NCLOUD_API_GW"

	AccessKeyIdFileKey     = "ncloud_access_key_id"
	SecretAccessKeyFileKey = "ncloud_secret_access_key"
	ApiUrlFileKey          = "ncloud_api_url"
)

const (
	DefaultConfigPath = ".ncloud"
	DefaultConfigFile = "configure"
)

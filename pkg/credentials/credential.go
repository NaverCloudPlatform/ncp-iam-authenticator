package credentials

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncp-iam-authenticator/pkg/constants"
	"github.com/NaverCloudPlatform/ncp-iam-authenticator/pkg/utils"
	"github.com/deiwin/interact"
	"github.com/pkg/errors"
	"gopkg.in/ini.v1"
	"log"
	"os"
	"strings"
)

var (
	checkNotEmpty = func(input string) error {
		// note that the inputs provided to these checks are already trimmed
		if input == "" {
			return errors.New("Input should not be empty!")
		}
		return nil
	}
)

type Config struct {
	*ncloud.APIKey
	ApiUrl string
}

func NewCredentialConfig(configPath string, profile string) (*Config, error) {
	cfg := NewCredentialFromEnv()

	if cfg != nil && cfg.Valid() {
		return cfg, nil
	}

	cfg = NewCredentialFromFile(configPath, profile)

	if cfg != nil && cfg.Valid() {
		os.Setenv(constants.ApiGwUrlEnv, cfg.ApiUrl)
		return cfg, nil
	}

	if err := cfg.NewCredentialFromCommandLine(); err != nil {
		return nil, errors.Wrap(err, "failed to set credential from cli")
	}

	if cfg != nil && cfg.Valid() {
		os.Setenv(constants.ApiGwUrlEnv, cfg.ApiUrl)
		cfg.WriteCredentialToFile(configPath, profile)
		return cfg, nil
	}

	return nil, errors.New("failed to get credential config")
}

func (c *Config) Valid() bool {
	return !utils.IsEmptyString(c.AccessKey) && !utils.IsEmptyString(c.SecretKey) && !utils.IsEmptyString(c.ApiUrl)
}

func (c *Config) WriteCredentialToFile(configPath string, profile string) {
	file, err := getNcloudConfigFile(configPath)
	fileSection, err := getNewSection(file, profile)

	if err != nil {
		log.Fatal(err)
		return
	}

	_, err = fileSection.NewKey(constants.AccessKeyIdFileKey, c.AccessKey)
	_, err = fileSection.NewKey(constants.SecretAccessKeyFileKey, c.SecretKey)
	_, err = fileSection.NewKey(constants.ApiUrlFileKey, c.ApiUrl)

	err = file.SaveTo(configPath)

	if err != nil {
		log.Fatal(err)
		return
	}
}

func (c *Config) NewCredentialFromCommandLine() error {
	actor := interact.NewActor(os.Stdin, os.Stderr)

	msg := "Ncloud Access Key Id []"
	accessKey, err := actor.PromptAndRetry(msg, checkNotEmpty)

	if err != nil {
		return err
	}

	msg = "Ncloud Secret Access Key []"
	secretKey, err := actor.PromptAndRetry(msg, checkNotEmpty)

	if err != nil {
		return err
	}

	msg = "Ncloud API URL []"
	apiUrl, err := actor.PromptAndRetry(msg, checkNotEmpty)

	if err != nil {
		return err
	}

	c.AccessKey = accessKey
	c.SecretKey = secretKey
	c.ApiUrl = apiUrl
	return nil
}

func NewCredentialFromEnv() *Config {

	id := os.Getenv(strings.ToUpper(constants.AccessKeyIdEnv))
	if id == "" {
		id = os.Getenv(strings.ToUpper(constants.AccessKeyEnv))
	}

	secret := os.Getenv(strings.ToUpper(constants.SecretAccessKeyEnv))
	if secret == "" {
		secret = os.Getenv(strings.ToUpper(constants.SecretKeyEnv))
	}

	url := os.Getenv(strings.ToUpper(constants.ApiGwUrlEnv))

	if utils.IsEmptyString(id) || utils.IsEmptyString(secret) || utils.IsEmptyString(url) {
		return nil
	}

	return &Config{
		&ncloud.APIKey{
			AccessKey: id,
			SecretKey: secret,
		},
		url,
	}
}

func NewCredentialFromFile(configPath string, profile string) *Config {
	config := &Config{}

	apiKey := &ncloud.APIKey{
		AccessKey: "",
		SecretKey: "",
	}

	file, err := getNcloudConfigFile(configPath)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	fileSection, err := getNcloudConfigFileSection(file, profile)

	if err != nil {
		log.Fatal(err)
		return nil
	}

	apiKey.AccessKey = fileSection.Key(constants.AccessKeyIdFileKey).String()
	apiKey.SecretKey = fileSection.Key(constants.SecretAccessKeyFileKey).String()
	apiUrl := fileSection.Key(constants.ApiUrlFileKey).String()

	config.APIKey = apiKey
	config.ApiUrl = apiUrl

	return config
}

func getNcloudConfigFileSection(file *ini.File, profile string) (*ini.Section, error) {

	section := ini.DefaultSection

	if !utils.IsEmptyString(profile) {
		section = profile
	}

	return file.Section(section), nil
}

func getNcloudConfigFile(configPath string) (*ini.File, error) {
	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		if _, osError := os.Create(configPath); osError != nil {
			return nil, errors.Wrap(osError, "cannot create config file")
		}
	}

	file, err := ini.Load(configPath)
	if err != nil {
		return nil, errors.Wrap(err, "cannot load ini file")
	}
	return file, nil
}

func getNewSection(file *ini.File, profile string) (*ini.Section, error) {
	section := ini.DefaultSection
	if !utils.IsEmptyString(profile) {
		return file.NewSection(profile)
	}

	return file.NewSection(section)
}

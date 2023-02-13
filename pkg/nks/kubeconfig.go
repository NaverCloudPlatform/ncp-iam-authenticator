package nks

import (
	"fmt"
	"github.com/NaverCloudPlatform/ncp-iam-authenticator/pkg/constants"
	"github.com/NaverCloudPlatform/ncp-iam-authenticator/pkg/utils"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"os"
	"path/filepath"
	"strings"
)

type KubeconfigParam struct {
	ClusterName    string
	UserName       string
	ContextName    string
	CurrentContext bool
	Profile        string
	ConfigFile     string
}

func (m Manager) GetIamKubeconfig(param *KubeconfigParam) (*clientcmdapi.Config, error) {
	cfg, err := m.GetKubeconfig()
	if err != nil {
		return nil, errors.Wrap(err, "get kubeconfig failed")
	}

	cluster, exist := cfg.Clusters[constants.NksKubeconfigClusterName]
	if !exist {
		return nil, fmt.Errorf("kubeconfig don't get cluster %s", constants.NksKubeconfigClusterName)
	}

	cfg.Clusters[param.ClusterName] = cluster
	delete(cfg.Clusters, constants.NksKubeconfigClusterName)

	cfg.AuthInfos[param.UserName] = m.makeIamUser(param.Profile, param.ConfigFile)

	delete(cfg.Contexts, constants.NksKubeconfigContextName)
	cfg.Contexts[param.ContextName] = &clientcmdapi.Context{
		Cluster:  param.ClusterName,
		AuthInfo: param.UserName,
	}

	cfg.CurrentContext = param.ContextName

	return cfg, nil
}

func (m Manager) UpdateIamKubeconfig(param *KubeconfigParam, config *clientcmdapi.Config, overwrite bool) error {
	if !overwrite {
		var duplicateName []string
		if _, exist := config.Clusters[param.ClusterName]; exist {
			duplicateName = append(duplicateName, "cluster name: "+param.ClusterName)
		}
		if _, exist := config.Clusters[param.UserName]; exist {
			duplicateName = append(duplicateName, "user name: "+param.UserName)
		}
		if _, exist := config.Clusters[param.ContextName]; exist {
			duplicateName = append(duplicateName, "context name: "+param.ContextName)
		}
		if len(duplicateName) != 0 {
			return fmt.Errorf("some names are duplicated: %s. if you want to overwrite it, please use --overwrite flag", strings.Join(duplicateName, ", "))
		}
	}

	orgCfg, err := m.GetKubeconfig()
	if err != nil {
		return errors.Wrap(err, "get kubeconfig failed")
	}

	cluster, exist := orgCfg.Clusters[constants.NksKubeconfigClusterName]
	if !exist {
		return fmt.Errorf("kubeconfig don't get cluster %s", constants.NksKubeconfigClusterName)
	}
	config.Clusters[param.ClusterName] = cluster
	config.AuthInfos[param.UserName] = m.makeIamUser(param.ConfigFile, param.Profile)
	config.Contexts[param.ContextName] = &clientcmdapi.Context{
		Cluster:  param.ClusterName,
		AuthInfo: param.UserName,
	}

	if param.CurrentContext {
		log.Debug().Msg("set current-context")
		config.CurrentContext = param.ContextName
	}
	return nil
}

func (m Manager) makeIamUser(configFile, profile string) *clientcmdapi.AuthInfo {
	args := []string{
		"token",
		"--clusterUuid",
		m.clusterUuid,
		"--region",
		m.region,
	}

	home, _ := os.UserHomeDir()
	if configFile != filepath.Join(home, constants.NcloudConfigPath, constants.NcloudConfigFile) {
		args = append(args, "--credentialConfig")
		args = append(args, configFile)
	}
	if !utils.IsEmptyString(profile) {
		args = append(args, "--profile")
		args = append(args, profile)
	}

	return &clientcmdapi.AuthInfo{
		Exec: &clientcmdapi.ExecConfig{
			Command:    "ncp-iam-authenticator",
			Args:       args,
			APIVersion: "client.authentication.k8s.io/v1beta1",
		},
	}
}

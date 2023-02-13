package nks

import (
	"fmt"
	"github.com/NaverCloudPlatform/ncp-iam-authenticator/pkg/constants"
	"github.com/NaverCloudPlatform/ncp-iam-authenticator/pkg/utils"
	"github.com/pkg/errors"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"os"
	"path/filepath"
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

	args := []string{
		"token",
		"--clusterUuid",
		m.clusterUuid,
		"--region",
		m.region,
	}

	home, _ := os.UserHomeDir()
	if param.ConfigFile != filepath.Join(home, constants.NcloudConfigPath, constants.NcloudConfigFile) {
		args = append(args, "--credentialConfig")
		args = append(args, param.ConfigFile)
	}
	if !utils.IsEmptyString(param.Profile) {
		args = append(args, "--profile")
		args = append(args, param.Profile)
	}

	cfg.AuthInfos[param.UserName] = &clientcmdapi.AuthInfo{
		Exec: &clientcmdapi.ExecConfig{
			Command:    "ncp-iam-authenticator",
			Args:       args,
			APIVersion: "client.authentication.k8s.io/v1beta1",
		},
	}

	delete(cfg.Contexts, constants.NksKubeconfigContextName)
	cfg.Contexts[param.ContextName] = &clientcmdapi.Context{
		Cluster:  param.ClusterName,
		AuthInfo: param.UserName,
	}

	cfg.CurrentContext = param.ContextName

	return cfg, nil
}

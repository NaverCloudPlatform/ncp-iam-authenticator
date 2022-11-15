package kubeconfig

import (
	"context"
	"fmt"
	"github.com/NaverCloudPlatform/ncp-iam-authenticator/pkg/cluster"
	"github.com/NaverCloudPlatform/ncp-iam-authenticator/pkg/utils"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Manager struct {
	clusterUuid   string
	clusterClient cluster.Client
	region        string
}

func NewManager(clusterUuid string, clusterClient cluster.Client, region string) *Manager {
	return &Manager{
		clusterUuid, clusterClient, region,
	}
}

func (m Manager) GetKubeconfig() (*KubectlConfig, error) {
	ctx := context.Background()

	cluster, err := m.clusterClient.ClustersUuidGet(ctx, &m.clusterUuid)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get kubernetes cluster from api")
	}
	if *cluster.Cluster.Status == "CREATING" {
		return nil, errors.New("kubernetes cluster is not running")
	}

	receivedConfig, err := m.clusterClient.ClustersUuidKubeconfigGet(ctx, &m.clusterUuid)

	if err != nil {
		return nil, errors.Wrap(err, "failed to get kubeconfig from api")
	}

	kubectlConfig := KubectlConfig{}
	if err := yaml.Unmarshal([]byte(*receivedConfig.Kubeconfig), &kubectlConfig); err != nil {
		return nil, err
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to convert string to rest config")
	}
	return &kubectlConfig, nil
}

func (m Manager) ApplyIamToKubeconfig(config *KubectlConfig, profile string, configFile string, setDefaultConfigFile bool) {
	userName := "nks-user"
	currentContext := fmt.Sprintf("%s@%s", userName, config.Clusters[0].Name)
	config.CurrentContext = currentContext

	config.Contexts = []*KubectlContextWithName{
		{
			Name: currentContext,
			Context: KubectlContext{
				Cluster: config.Clusters[0].Name,
				User:    userName,
			},
		},
	}

	config.Users = []*KubectlUserWithName{
		{
			Name: userName,
			User: KubectlUser{
				Exec: ExecConfig{
					APIVersion: "client.authentication.k8s.io/v1beta1",
					Command:    "ncp-iam-authenticator",
					Args: []string{
						"token",
						"--clusterUuid",
						m.clusterUuid,
						"--region",
						m.region,
					},
				},
			},
		},
	}

	if !setDefaultConfigFile {
		config.Users[0].User.Exec.Args = append(config.Users[0].User.Exec.Args, "--credentialConfig")
		config.Users[0].User.Exec.Args = append(config.Users[0].User.Exec.Args, configFile)
	}
	if !utils.IsEmptyString(profile) {
		config.Users[0].User.Exec.Args = append(config.Users[0].User.Exec.Args, "--profile")
		config.Users[0].User.Exec.Args = append(config.Users[0].User.Exec.Args, profile)
	}
}

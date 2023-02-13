package cmd

import (
	"fmt"
	"github.com/NaverCloudPlatform/ncp-iam-authenticator/pkg/credentials"
	"github.com/NaverCloudPlatform/ncp-iam-authenticator/pkg/nks"
	"github.com/NaverCloudPlatform/ncp-iam-authenticator/pkg/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
	"sigs.k8s.io/yaml"
	"strings"
)

type createKubeconfigOptions struct {
	format      string
	output      string
	clusterUuid string
	region      string
	clusterName string
	userName    string
	contextName string
}

func (o *createKubeconfigOptions) SetDefault(clusterName string) {
	o.region = strings.ToUpper(o.region)
	defaultName := fmt.Sprintf("nks_%s_%s_%s", strings.ToLower(o.region), clusterName, o.clusterUuid)
	if utils.IsEmptyString(o.output) {
		o.output = fmt.Sprintf("kubeconfig-%s.yaml", o.clusterUuid)
	}

	var isClusterNameFlagEmpty, IsUserNameFlagEmpty bool
	if isClusterNameFlagEmpty = utils.IsEmptyString(o.clusterName); isClusterNameFlagEmpty {
		o.clusterName = defaultName
	}
	if IsUserNameFlagEmpty = utils.IsEmptyString(o.userName); IsUserNameFlagEmpty {
		o.userName = defaultName
	}
	if isClusterNameFlagEmpty && IsUserNameFlagEmpty {
		o.contextName = o.clusterName
	} else {
		o.contextName = fmt.Sprintf("%s@%s", o.userName, o.clusterName)
	}
}

var (
	nksManager *nks.Manager
)

func NewCmdCreateKubeconfig(rootOptions *rootOptions) *cobra.Command {
	options := &createKubeconfigOptions{}

	cmd := &cobra.Command{
		Use:   "create-kubeconfig",
		Short: "Get Kubeconfig to access kubernetes",
		Long:  ``,
		PreRun: func(cmd *cobra.Command, args []string) {
			credentialConfig, err := credentials.NewCredentialConfig(rootOptions.configFile, rootOptions.profile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to get credential config: %v", err)
				os.Exit(1)
			}

			log.Debug().
				Str("access_key", credentialConfig.APIKey.AccessKey).
				Str("secret_key", credentialConfig.APIKey.SecretKey).
				Str("api_gw_url", credentialConfig.ApiUrl).Msg("")

			nksManager = nks.NewManager(options.clusterUuid, options.region, credentialConfig.APIKey)

			cluster, err := nksManager.GetCluster()
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to get cluster: %v", err)
				os.Exit(1)
			}

			options.SetDefault(*cluster.Name)
			log.Debug().Str("options", fmt.Sprintf("%+v", options)).Msg("")
		},
		Run: func(cmd *cobra.Command, args []string) {
			kubeconfig, err := nksManager.GetIamKubeconfig(&nks.KubeconfigParam{
				ClusterName: options.clusterName,
				UserName:    options.userName,
				ContextName: options.contextName,
				Profile:     rootOptions.profile,
				ConfigFile:  rootOptions.configFile},
			)

			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to get iam kubeconfig: %v", err)
				os.Exit(1)
			}

			if options.format == "json" {
				yamlBytes, err := clientcmd.Write(*kubeconfig)
				if err != nil {
					fmt.Fprintf(os.Stderr, "failed to write kubeconfig string: %v", err)
					os.Exit(1)
				}
				jsonBytes, err := yaml.YAMLToJSON(yamlBytes)
				if err != nil {
					fmt.Fprintf(os.Stderr, "failed to convert yaml to json: %v", err)
					os.Exit(1)
				}

				dir := filepath.Dir(options.output)
				if _, err := os.Stat(dir); os.IsNotExist(err) {
					if err = os.MkdirAll(dir, 0755); err != nil {
						fmt.Fprintf(os.Stderr, "failed to make dir: %v", err)
						os.Exit(1)
					}
				}
				if err := os.WriteFile(options.output, jsonBytes, 0600); err != nil {
					fmt.Fprintf(os.Stderr, "failed to write kubeconfig file: %v", err)
				}
			} else {
				if err := clientcmd.WriteToFile(*kubeconfig, options.output); err != nil {
					fmt.Fprintf(os.Stderr, "failed to write kubeconfig file: %v", err)
					os.Exit(1)
				}
			}

			fmt.Fprintf(os.Stdout, "kubeconfig created successfully")
		},
	}

	cmd.PersistentFlags().StringVar(&options.clusterUuid, "clusterUuid", "", "clusterUuid")
	cmd.PersistentFlags().StringVar(&options.region, "region", "", "cluster region")
	cmd.PersistentFlags().StringVar(&options.format, "format", "yaml", "format")
	cmd.PersistentFlags().StringVarP(&options.output, "output", "o", "", "kubeconfig output path")
	cmd.PersistentFlags().StringVar(&options.clusterName, "clusterName", "", "kubeconfig output cluster name")
	cmd.PersistentFlags().StringVar(&options.userName, "userName", "", "kubeconfig output user name")

	if err := cmd.MarkPersistentFlagRequired("clusterUuid"); err != nil {
		fmt.Fprintf(os.Stderr, "failed to run create-kubeconfig: %v", err)
		os.Exit(1)
	}
	if err := cmd.MarkPersistentFlagRequired("region"); err != nil {
		fmt.Fprintf(os.Stderr, "failed to run create-kubeconfig: %v", err)
		os.Exit(1)
	}

	return cmd
}

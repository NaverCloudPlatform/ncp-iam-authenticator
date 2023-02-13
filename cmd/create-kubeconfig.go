package cmd

import (
	"fmt"
	"github.com/NaverCloudPlatform/ncp-iam-authenticator/pkg/credentials"
	"github.com/NaverCloudPlatform/ncp-iam-authenticator/pkg/nks"
	"github.com/NaverCloudPlatform/ncp-iam-authenticator/pkg/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
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

func NewCmdCreateKubeconfig(rootOptions *rootOptions) *cobra.Command {
	options := &createKubeconfigOptions{}

	cmd := &cobra.Command{
		Use:   "create-kubeconfig",
		Short: "Get Kubeconfig to access kubernetes",
		Long:  ``,
		PreRun: func(cmd *cobra.Command, args []string) {
			credentialConfig, err := credentials.NewCredentialConfig(rootOptions.configFile, rootOptions.profile)
			if err != nil {
				log.Error().Err(err).Msg("failed to get credential config")
				fmt.Fprintf(os.Stdout, "run create-kubeconfig failed. please check your credentialConfig and profile.")
				os.Exit(1)
			}

			log.Debug().
				Str("access_key", credentialConfig.APIKey.AccessKey).
				Str("secret_key", credentialConfig.APIKey.SecretKey).
				Str("api_gw_url", credentialConfig.ApiUrl).Msg("credential config")

			nksManager = nks.NewManager(options.clusterUuid, options.region, credentialConfig.APIKey)

			cluster, err := nksManager.GetCluster()
			if err != nil {
				log.Fatal().Err(err).Msg("failed to get cluster")
			}

			options.SetDefault(*cluster.Name)
			log.Debug().Str("options", fmt.Sprintf("%+v", options)).Msg("create-kubeconfig options")
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
				log.Fatal().Err(err).Msg("failed to get iam kubeconfig")
			}

			if err := utils.WriteKubeconfigToFile(kubeconfig, options.format, options.output); err != nil {
				log.Fatal().Err(err).Msg("failed to write kubeconfig to file")
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
		log.Error().Err(err).Msg("failed to get clusterUuid")
		fmt.Fprintf(os.Stdout, "failed to run create-kubeconfig. please check your clusterUuid")
		os.Exit(1)
	}
	if err := cmd.MarkPersistentFlagRequired("region"); err != nil {
		log.Error().Err(err).Msg("failed to get clusterUuid")
		fmt.Fprintf(os.Stdout, "failed to run create-kubeconfig. please check your clusterUuid")
		os.Exit(1)
	}

	return cmd
}

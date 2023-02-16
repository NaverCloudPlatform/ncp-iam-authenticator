package cmd

import (
	"fmt"
	"github.com/NaverCloudPlatform/ncp-iam-authenticator/pkg/constants"
	"github.com/NaverCloudPlatform/ncp-iam-authenticator/pkg/credentials"
	"github.com/NaverCloudPlatform/ncp-iam-authenticator/pkg/nks"
	"github.com/NaverCloudPlatform/ncp-iam-authenticator/pkg/utils"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"os"
	"path/filepath"
	"strings"
)

type updateKubeconfigOptions struct {
	format         string
	clusterUuid    string
	region         string
	clusterName    string
	userName       string
	contextName    string
	kubeconfig     string
	overwrite      bool
	currentContext bool
}

func (o *updateKubeconfigOptions) setDefault(clusterName string) error {
	o.region = strings.ToUpper(o.region)
	defaultName := fmt.Sprintf("nks_%s_%s_%s", strings.ToLower(o.region), clusterName, o.clusterUuid)

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

	if utils.IsEmptyString(o.kubeconfig) {
		configs := os.Getenv(constants.KubeconfigEnv)
		if utils.IsEmptyString(configs) {
			home, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			o.kubeconfig = filepath.Join(home, constants.KubeconfigRootDir, constants.KubeconfigFileName)
		} else {
			o.kubeconfig = filepath.SplitList(configs)[0]
		}
	}
	return nil
}

func (o *updateKubeconfigOptions) checkRequired() error {
	var errorList []string
	if o.region == "" {
		errorList = append(errorList, "--region")
	}
	if o.clusterUuid == "" {
		errorList = append(errorList, "--clusterUuid")
	}
	if len(errorList) != 0 {
		return fmt.Errorf("required flag %s not set", strings.Join(errorList, ", "))
	}
	return nil
}

func NewCmdUpdateKubeconfig(rootOptions *rootOptions) *cobra.Command {
	options := &updateKubeconfigOptions{}

	cmd := &cobra.Command{
		Use:   "update-kubeconfig",
		Short: "update Kubeconfig to access kubernetes",
		Long:  ``,
		PreRun: func(cmd *cobra.Command, args []string) {
			log.Debug().Str("options", fmt.Sprintf("%+v", options)).Msg("init update-kubeconfig options")
			if err := options.checkRequired(); err != nil {
				log.Error().Err(err).Msg("required flags not set")
				fmt.Fprintln(os.Stdout, "run update-kubeconfig failed. please check your required flags.")
				os.Exit(1)
			}

			credentialConfig, err := credentials.NewCredentialConfig(rootOptions.configFile, rootOptions.profile)
			if err != nil {
				log.Error().Err(err).Msg("failed to get credential config")
				fmt.Fprintln(os.Stdout, "run update-kubeconfig failed. please check your credentialConfig and profile.")
				os.Exit(1)
			}

			log.Debug().
				Str("access_key", credentialConfig.APIKey.AccessKey).
				Str("secret_key", credentialConfig.APIKey.SecretKey).
				Str("api_gw_url", credentialConfig.ApiUrl).Msg("credential config")

			nksManager = nks.NewManager(options.clusterUuid, strings.ToUpper(options.region), credentialConfig.APIKey)

			cluster, err := nksManager.GetCluster()
			if err != nil {
				log.Error().Err(err).Msg("failed to get cluster")
				fmt.Fprintln(os.Stdout, "run update-kubeconfig failed. please check your credentialConfig or clusterUuid.")
				os.Exit(1)
			}
			if *cluster.Status == "CREATING" {
				log.Error().Str("clusterStatus", *cluster.Status).Msg("cluster status is CREATING")
				fmt.Fprintln(os.Stdout, "run update-kubeconfig failed. please try again after cluster creation is complete.")
				os.Exit(1)
			}

			if err := options.setDefault(*cluster.Name); err != nil {
				log.Error().Err(err).Msg("failed to set options")
				fmt.Fprintln(os.Stdout, "run update-kubeconfig failed. please check your kubeconfig env or kubeconfig flag.")
			}
			log.Debug().Str("options", fmt.Sprintf("%+v", options)).Msg("update-kubeconfig options")
		},
		Run: func(cmd *cobra.Command, args []string) {
			var kubeconfig *clientcmdapi.Config

			if _, err := os.Stat(options.kubeconfig); errors.Is(err, os.ErrNotExist) {
				kubeconfig = clientcmdapi.NewConfig()
			} else {
				kubeconfig, err = clientcmd.LoadFromFile(options.kubeconfig)
				if err != nil {
					log.Error().Err(err).Msg("failed to load kubeconfig from file")
					fmt.Fprintln(os.Stdout, "run update-kubeconfig failed. please check your kubeconfig file or kubeconfig flag.")
					os.Exit(1)
				}
			}

			if !options.overwrite {
				if err := utils.ValidateKubeconfigDupliacted(options.clusterName, options.userName, options.contextName, kubeconfig); err != nil {
					log.Error().Err(err).Msg("duplicated name")
					fmt.Fprintln(os.Stdout, "run update-kubeconfig failed. please check your kubeconfig's clusterName, userName, contextName. if you want to overwrite it, please use --overwrite flag")
					os.Exit(1)
				}
			}
			if err := nksManager.UpdateIamKubeconfig(&nks.KubeconfigParam{
				ClusterName:    options.clusterName,
				UserName:       options.userName,
				ContextName:    options.contextName,
				Profile:        rootOptions.profile,
				ConfigFile:     rootOptions.configFile,
				CurrentContext: options.currentContext,
			}, kubeconfig, options.overwrite); err != nil {
				log.Fatal().Err(err).Msg("failed to update iam kubeconfig")
			}

			if err := utils.WriteKubeconfigToFile(kubeconfig, options.format, options.kubeconfig); err != nil {
				log.Fatal().Err(err).Msg("failed to write kubeconfig to file")
			}

			fmt.Fprintln(os.Stdout, "kubeconfig updated successfully.")
		},
	}

	cmd.PersistentFlags().StringVar(&options.clusterUuid, "clusterUuid", "", "clusterUuid")
	cmd.PersistentFlags().StringVar(&options.region, "region", "", "cluster region")
	cmd.PersistentFlags().StringVar(&options.format, "format", "yaml", "format")
	cmd.PersistentFlags().StringVar(&options.clusterName, "clusterName", "", "kubeconfig output cluster name")
	cmd.PersistentFlags().StringVar(&options.userName, "userName", "", "kubeconfig output user name")
	cmd.PersistentFlags().StringVar(&options.kubeconfig, "kubeconfig", "", "kubeconfig file path")
	cmd.PersistentFlags().BoolVar(&options.overwrite, "overwrite", false, "if the cluster name, user name, or context name is duplicated, overwrite them")
	cmd.PersistentFlags().BoolVar(&options.currentContext, "currentContext", true, "set current-context")

	return cmd
}

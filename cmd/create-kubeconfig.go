package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vnks"
	"github.com/NaverCloudPlatform/ncp-iam-authenticator/pkg/credentials"
	"github.com/NaverCloudPlatform/ncp-iam-authenticator/pkg/kubeconfig"
	"github.com/NaverCloudPlatform/ncp-iam-authenticator/pkg/utils"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
)

type createKubeconfigOptions struct {
	printDebugLog bool
	format        string
	output        string
	clusterUuid   string
	region        string
}

func NewCmdCreateKubeconfig(defaultOptions *defaultOptions) *cobra.Command {
	options := &createKubeconfigOptions{}

	cmd := &cobra.Command{
		Use:   "create-kubeconfig",
		Short: "Get Kubeconfig to access kubernetes",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			if options.printDebugLog {
				utils.PrintLog(os.Stdout, []string{
					fmt.Sprintf("PROFILE: %s", defaultOptions.profile),
					fmt.Sprintf("CONFIG FILE: %s", defaultOptions.configFile),
				})
			}

			options.region = strings.ToUpper(options.region)
			options.output = getOutputFileName(options.output, options.clusterUuid)

			credentialConfig, err := credentials.NewCredentialConfig(defaultOptions.configFile, defaultOptions.profile)
			if options.printDebugLog {
				utils.PrintLog(os.Stdout, []string{
					fmt.Sprintf("ACCESS_KEY: %s", credentialConfig.APIKey.AccessKey),
					fmt.Sprintf("API_GW_URL: %s", credentialConfig.ApiUrl),
				})
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "could not get credential config: %v", err)
				os.Exit(1)
			}

			ncloudConfig := vnks.NewConfiguration(options.region, credentialConfig.APIKey)

			kubeconfigManager := kubeconfig.NewManager(options.clusterUuid, ncloudConfig, options.region)

			kubeConfig, err := kubeconfigManager.GetKubeconfig()
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to get kubeconfig: %v", err)
				os.Exit(1)
			}

			kubeconfigManager.ApplyIamToKubeconfig(kubeConfig, defaultOptions.profile, defaultOptions.configFile, defaultOptions.setDefaultConfigFile)

			var kubeconfigBytes []byte

			if options.format == "yaml" {
				kubeconfigBytes, err = yaml.Marshal(kubeConfig)

				if err != nil {
					fmt.Fprintf(os.Stderr, "failed to marshal kubeconfig yaml: %v", err)
					os.Exit(1)
				}
			} else if options.format == "json" {
				kubeconfigBytes, err = json.Marshal(kubeConfig)

				if err != nil {
					fmt.Fprintf(os.Stderr, "failed to marshal kubeconfig yaml: %v", err)
					os.Exit(1)
				}
			}

			if err = os.WriteFile(options.output, kubeconfigBytes, 0644); err != nil {
				fmt.Fprintf(os.Stderr, "failed to write kubeconfig file: %v", err)
				os.Exit(1)
			}

			fmt.Fprintf(os.Stdout, "kubeconfig created successfully")
			return
		},
	}

	cmd.PersistentFlags().StringVar(&options.clusterUuid, "clusterUuid", "", "clusterUuid")
	cmd.PersistentFlags().StringVar(&options.region, "region", "", "cluster region")
	cmd.PersistentFlags().StringVar(&options.format, "format", "yaml", "format")
	cmd.PersistentFlags().StringVarP(&options.output, "output", "o", "", "kubeconfig output path")
	cmd.PersistentFlags().BoolVar(&options.printDebugLog, "debug", false, "debug option")

	cmd.MarkPersistentFlagRequired("clusterUuid")
	cmd.MarkPersistentFlagRequired("region")

	return cmd
}

func getOutputFileName(output, clusterUuid string) string {
	if utils.IsEmptyString(output) {
		return "kubeconfig-" + clusterUuid + ".yaml"
	}

	return output
}

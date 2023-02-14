package cmd

import (
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vnks"
	"github.com/NaverCloudPlatform/ncp-iam-authenticator/pkg/credentials"
	"github.com/NaverCloudPlatform/ncp-iam-authenticator/pkg/token"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
)

type tokenOptions struct {
	clusterUuid string
	region      string
}

func NewTokenCmd(defaultOptions *rootOptions) *cobra.Command {
	options := tokenOptions{}

	cmd := &cobra.Command{
		Use:   "token",
		Short: "Authenticate using SubAccount and get token for Kubernetes",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			gen, err := token.NewGenerator()
			if err != nil {
				log.Fatal().Err(err).Msg("new token generator failed")
			}

			credentialConfig, err := credentials.NewCredentialConfig(defaultOptions.configFile, defaultOptions.profile)
			if err != nil {
				log.Fatal().Err(err).Msg("failed to get credential config")
			}

			ncloudConfig := vnks.NewConfiguration(options.region, credentialConfig.APIKey)

			tok, err := gen.Get(ncloudConfig.GetCredentials(), options.clusterUuid, options.region)
			if err != nil {
				log.Fatal().Err(err).Msg("failed to gen token")
			}

			genToken, err := gen.FormatJSON(*tok)
			if err != nil {
				log.Fatal().Err(err).Msg("failed to token format json")
			}

			fmt.Fprint(os.Stdout, genToken)
		},
	}

	cmd.PersistentFlags().StringVar(&options.clusterUuid, "clusterUuid", "", "clusterUuid")
	cmd.PersistentFlags().StringVar(&options.region, "region", "", "cluster region")

	if err := cmd.MarkPersistentFlagRequired("clusterUuid"); err != nil {
		log.Error().Err(err).Msg("failed to get clusterUuid")
		fmt.Fprintln(os.Stdout, "failed to run update-kubeconfig. please check your clusterUuid")
		os.Exit(1)
	}
	if err := cmd.MarkPersistentFlagRequired("region"); err != nil {
		log.Error().Err(err).Msg("failed to get region")
		fmt.Fprintln(os.Stdout, "failed to run update-kubeconfig. please check your region")
		os.Exit(1)
	}

	return cmd
}

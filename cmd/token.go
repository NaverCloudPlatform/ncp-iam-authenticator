package cmd

import (
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vnks"
	"github.com/NaverCloudPlatform/ncp-iam-authenticator/pkg/credentials"
	"github.com/NaverCloudPlatform/ncp-iam-authenticator/pkg/token"
	"github.com/spf13/cobra"
	"os"
)

type tokenOptions struct {
	clusterUuid string
	region      string
}

func NewTokenCmd(defaultOptions *defaultOptions) *cobra.Command {
	options := tokenOptions{}

	cmd := &cobra.Command{
		Use:   "token",
		Short: "Authenticate using SubAccount and get token for Kubernetes",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			gen, err := token.NewGenerator()

			if err != nil {
				fmt.Fprintf(os.Stderr, "could not get token: %v", err)
				os.Exit(1)
			}

			credentialConfig, err := credentials.NewCredentialConfig(defaultOptions.configFile, defaultOptions.profile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "could not get credential config: %v", err)
				os.Exit(1)
			}

			ncloudConfig := vnks.NewConfiguration(options.region, credentialConfig.APIKey)

			tok, err := gen.Get(ncloudConfig.GetCredentials(), options.clusterUuid, options.region)
			if err != nil {
				fmt.Fprintf(os.Stderr, "could not get token: %v", err)
				os.Exit(1)
			}

			genToken, err := gen.FormatJSON(*tok)

			if err != nil {
				fmt.Fprintf(os.Stderr, "could not get token: %v\n", err)
				os.Exit(1)
			}

			fmt.Fprint(os.Stdout, genToken)
		},
	}

	cmd.PersistentFlags().StringVar(&options.clusterUuid, "clusterUuid", "", "clusterUuid")
	cmd.PersistentFlags().StringVar(&options.region, "region", "", "cluster region")

	if err := cmd.MarkPersistentFlagRequired("clusterUuid"); err != nil {
		fmt.Fprintf(os.Stderr, "failed to run token: %v", err)
		os.Exit(1)
	}
	if err := cmd.MarkPersistentFlagRequired("region"); err != nil {
		fmt.Fprintf(os.Stderr, "failed to run token: %v", err)
		os.Exit(1)
	}

	return cmd
}

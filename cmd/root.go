package cmd

import (
	"fmt"
	"github.com/NaverCloudPlatform/ncp-iam-authenticator/pkg/constants"
	"github.com/NaverCloudPlatform/ncp-iam-authenticator/pkg/utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

type rootOptions struct {
	profile       string
	configFile    string
	printDebugLog bool
}

func Execute() {
	if err := NewDefaultCmd().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to execute ncp-iam-authenticator: %v", err)
		os.Exit(1)
	}
}

func NewDefaultCmd() *cobra.Command {
	options := &rootOptions{}
	cmd := &cobra.Command{
		Use:   "ncp-iam-authenticator",
		Short: "ncloud kubernetes service iam authenticator",
		Long:  `cli written to authenticate with iam in ncloud kubernetes service`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			zerolog.SetGlobalLevel(zerolog.ErrorLevel)
			if options.printDebugLog {
				zerolog.SetGlobalLevel(zerolog.DebugLevel)
			}

			log.Debug().Str("profile", options.profile).Str("credentialConfig", options.configFile).Msg("")

			if utils.IsEmptyString(options.configFile) {
				home, err := os.UserHomeDir()
				if err != nil {
					fmt.Fprintf(os.Stdout, "failed to get homde dir: %v", err)
					os.Exit(1)
				}
				options.configFile = filepath.Join(home, constants.NcloudConfigPath, constants.NcloudConfigFile)
			}
		},
	}

	cmd.PersistentFlags().StringVar(&options.profile, "profile", "", "profile")
	cmd.PersistentFlags().StringVar(&options.configFile, "credentialConfig", os.Getenv(constants.ProfileEnv), "credential config path (default : ~/.ncloud/configure)")
	cmd.PersistentFlags().BoolVar(&options.printDebugLog, "debug", false, "debug option")

	cmd.AddCommand(NewVersionCmd())
	cmd.AddCommand(NewCmdCreateKubeconfig(options))
	cmd.AddCommand(NewTokenCmd(options))

	cmd.CompletionOptions.DisableDefaultCmd = true

	return cmd
}

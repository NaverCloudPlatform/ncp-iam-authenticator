package cmd

import (
	"fmt"
	"github.com/NaverCloudPlatform/ncp-iam-authenticator/pkg/constants"
	"github.com/NaverCloudPlatform/ncp-iam-authenticator/pkg/utils"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

type defaultOptions struct {
	profile              string
	configFile           string
	setDefaultConfigFile bool
}

func Execute() {
	if err := NewDefaultCmd().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to execute ncp-iam-authenticator: %v", err)
		os.Exit(1)
	}
}

func NewDefaultCmd() *cobra.Command {
	options := &defaultOptions{}

	cmd := &cobra.Command{
		Use:   "ncp-iam-authenticator",
		Short: "ncloud kubernetes service iam authenticator",
		Long:  `cli written to authenticate with iam in ncloud kubernetes service`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if utils.IsEmptyString(options.configFile) {
				home, err := os.UserHomeDir()
				cobra.CheckErr(err)
				options.setDefaultConfigFile = true
				options.configFile = filepath.Join(home, constants.DefaultConfigPath, constants.DefaultConfigFile)
			} else {
				options.setDefaultConfigFile = false
			}
		},
	}

	cmd.PersistentFlags().StringVar(&options.profile, "profile", "", "profile")
	cmd.PersistentFlags().StringVar(&options.configFile, "credentialConfig", os.Getenv(constants.ProfileEnv), "credential config path (default : ~/.ncloud/configure)")

	cmd.AddCommand(NewVersionCmd())
	cmd.AddCommand(NewCmdCreateKubeconfig(options))
	cmd.AddCommand(NewTokenCmd(options))

	cmd.CompletionOptions.DisableDefaultCmd = true

	return cmd
}

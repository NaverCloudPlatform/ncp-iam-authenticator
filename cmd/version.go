package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var version = ""

func NewVersionCmd() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Show the version info of the ncp-iam-authenticator",
		Long:  `Show the version info of the ncp-iam-authenticator`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version)
		},
	}

	return versionCmd
}

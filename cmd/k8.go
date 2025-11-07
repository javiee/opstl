package cmd

import (
	"github.com/spf13/cobra"
)

var namespace string

var k8Cmd = &cobra.Command{
	Use:   "kubernetes",
	Short: "Get all namespaces in the cluster",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rootCmd.AddCommand(k8Cmd)
	k8Cmd.Aliases = []string{"k8", "kubectl"}
	k8Cmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Get all namespaces in the cluster")
}

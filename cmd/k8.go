package cmd

import (
	"github.com/spf13/cobra"
)

var namespace string

var k8Cmd = &cobra.Command{
	Use:   "kubernetes",
	Short: "Kubernetes operations for your team's namespaces",
	Long: `Kubernetes operations for managing pods and events in your team's namespaces.

Available Subcommands:
  get-pods      List all pods owned by the configured team
  get-events    List all warning events in owned namespaces
  watch-pods    Watch and stream logs from pods in real-time

Examples:
  opstl kubernetes get-pods
  opstl k8 get-pods --status Running
  opstl kubectl get-events
  opstl k8 watch-pods -l app=myapp -r "myapp-.*"`,
}

func init() {
	rootCmd.AddCommand(k8Cmd)
	k8Cmd.Aliases = []string{"k8", "kubectl"}
	k8Cmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Get all namespaces in the cluster")
}

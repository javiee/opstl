package cmd

import (
	"github.com/spf13/cobra"
)

var (
	labelSelector string
	prefix        string
)

var watchPodsCmd = &cobra.Command{
	Use:   "watch-pods",
	Short: "Watch pods across owned namespaces by common labels",
	Run: func(cmd *cobra.Command, args []string) {
		watchPods() // Implementation will go here
	},
}

func init() {
	k8Cmd.AddCommand(watchPodsCmd)
	watchPodsCmd.Flags().StringVarP(&labelSelector, "label-selector", "l", "", "Label selector to filter pods")
	watchPodsCmd.Flags().StringVarP(&prefix, "prefix", "p", "", "Prefix pod names to watch")
}

func watchPods() {
	// Implementation for watching pods will go here
}

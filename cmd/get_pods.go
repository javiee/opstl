package cmd

import (
	"context"
	"fmt"
	"opstl/internal/utils"
	"strings"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	statusFilter string
)

var getPodsCmd = &cobra.Command{
	Use:   "get-pods",
	Short: "List all pods owned by the configured team",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Fetching events:")
		getPods()
	},
}

func init() {
	k8Cmd.AddCommand(getPodsCmd)
	getPodsCmd.Flags().StringVarP(&statusFilter, "status", "s", "", "Get pods filtered by status")
}

func getPods() {

	clientset, err := NewClientSet()
	if err != nil {
		fmt.Println("❌", err)
		return
	}

	namespaces, err := GetNamespaces(clientset)

	if err != nil {
		fmt.Println("❌", err)
		return
	}

	fmt.Printf("%-50s %-50s %-10s %-10s\n", "POD NAME", "NAMESPACE", "STATUS", "AGE")
	for _, ns := range namespaces {
		pods, err := clientset.CoreV1().Pods(ns).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			fmt.Printf("Error fetching pods in namespace %s: %v\n", ns, err)
			continue
		}
		fmt.Println(strings.Repeat("-", 65))
		for _, pod := range pods.Items {
			if statusFilter != "" && string(pod.Status.Phase) != statusFilter {
				continue
			}
			fmt.Printf("%-50s %-50s %-10s %-10s\n", pod.Name, pod.Namespace, pod.Status.Phase, utils.GetAge(pod.CreationTimestamp.Time))
		}

	}

}

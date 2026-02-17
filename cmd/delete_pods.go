package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var targetPattern string

var deletePodsCmd = &cobra.Command{
	Use:   "delete-pods",
	Short: "Delete pods matching a regex pattern",
	Long: `Delete pods across your team's namespaces that match a regex pattern.

Lists matching pods and asks for confirmation before deleting.

Flags:
  -t, --target    Regex pattern to match pod names (required)

Examples:
  opstl k8 delete-pods -t "redstone-.*"
  opstl k8 delete-pods --target "my-app-[a-z]+-.*"`,
	Run: func(cmd *cobra.Command, args []string) {
		deletePods()
	},
}

func init() {
	k8Cmd.AddCommand(deletePodsCmd)
	deletePodsCmd.Flags().StringVarP(&targetPattern, "target", "t", "", "Regex pattern to match pod names for deletion")
}

func deletePods() {
	if targetPattern == "" {
		fmt.Println("❌ Please provide a target regex pattern with -t/--target.")
		return
	}

	targetRegex, err := regexp.Compile(targetPattern)
	if err != nil {
		fmt.Printf("❌ Invalid regex pattern: %v\n", err)
		return
	}

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

	type podTarget struct {
		name      string
		namespace string
	}

	var matched []podTarget

	for _, ns := range namespaces {
		pods, err := clientset.CoreV1().Pods(ns).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			fmt.Printf("Error fetching pods in namespace %s: %v\n", ns, err)
			continue
		}
		for _, pod := range pods.Items {
			if targetRegex.MatchString(pod.Name) {
				matched = append(matched, podTarget{name: pod.Name, namespace: pod.Namespace})
			}
		}
	}

	if len(matched) == 0 {
		fmt.Println("No pods matched the pattern:", targetPattern)
		return
	}

	fmt.Printf("The following %d pod(s) match the pattern %q:\n\n", len(matched), targetPattern)
	fmt.Printf("%-50s %-50s\n", "POD NAME", "NAMESPACE")
	fmt.Println(strings.Repeat("-", 100))
	for _, p := range matched {
		fmt.Printf("%-50s %-50s\n", p.name, p.namespace)
	}

	fmt.Print("\nAre you sure you want to delete these pods? (y/n): ")
	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(strings.ToLower(answer))

	if answer != "y" && answer != "yes" {
		fmt.Println("Aborted.")
		return
	}

	fmt.Println()
	for _, p := range matched {
		err := clientset.CoreV1().Pods(p.namespace).Delete(context.TODO(), p.name, metav1.DeleteOptions{})
		if err != nil {
			fmt.Printf("❌ Failed to delete %s/%s: %v\n", p.namespace, p.name, err)
		} else {
			fmt.Printf("✅ Deleted %s/%s\n", p.namespace, p.name)
		}
	}
}

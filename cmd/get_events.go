package cmd

import (
	"context"
	"fmt"
	"opstl/internal/utils"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var getEventsCmd = &cobra.Command{
	Use:   "get-events",
	Short: "List all events in a namespace",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Fetching events from ownded namespaces:")
		getEvents()
	},
}

func init() {
	k8Cmd.AddCommand(getEventsCmd)
	// getEventsCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "Specify namespace")
}

func getEvents() {

	fmt.Println("This function will list all events in the specified namespace.")

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

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Name", "Type", "Reason", "AGE", "Message"})

	for _, ns := range namespaces {

		events, err := clientset.CoreV1().Events(ns).List(context.TODO(), metav1.ListOptions{FieldSelector: "type=Warning"})

		if err != nil {
			fmt.Printf("Error fetching events in namespace %s: %v\n", ns, err)
			continue
		}
		for _, event := range events.Items {
			t.AppendRow(table.Row{event.Name, event.Type, event.Reason, utils.GetAge(event.CreationTimestamp.Time), utils.Truncate(event.Message, 50)})
			t.AppendSeparator()
		}
	}
	t.Render()

}

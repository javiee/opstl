package cmd

import (
	"context"
	"fmt"

	"path/filepath"

	"github.com/spf13/viper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func NewClientSet() (*kubernetes.Clientset, error) {

	// config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	// if err != nil {
	// 	return nil, fmt.Errorf("Error loading kubeconfig: %v\n", err)
	// }
	// var clientset *kubernetes.Clientset

	// clientset, err = kubernetes.NewForConfig(config)
	// if err != nil {
	// 	return nil, fmt.Errorf("Error creating Kubernetes client: %v\n", err)
	// }
	// return clientset, nil

	// Find the kubeconfig path (default: ~/.kube/config)
	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")

	fmt.Printf("Using kubeconfig: %s\n", kubeconfig)

	// Load the kubeconfig (like kubectl does)
	configLoadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}

	// Create the loader that applies context and user overrides
	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		configLoadingRules,
		configOverrides,
	)

	// Actually build the rest.Config
	config, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("error loading kubeconfig: %v", err)
	}

	// Build the clientset (shared for all resources)
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("error creating clientset: %v", err)
	}

	return clientset, nil
}

func containsNS(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func GetNamespaces(clientset *kubernetes.Clientset) ([]string, error) {

	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{
		LabelSelector: viper.GetString("teamLabelSelector"),
	})
	if err != nil {
		return nil, fmt.Errorf("Error fetching namespaces: %v\n", err)
	}

	var namespaceNames []string

	filteredNamespaces := viper.GetStringSlice("filteredNamespaces")

	for _, ns := range namespaces.Items {
		if !containsNS(filteredNamespaces, ns.Name) {
			namespaceNames = append(namespaceNames, ns.Name)
		}
	}
	return namespaceNames, nil
}

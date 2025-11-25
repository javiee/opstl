package cmd

import (
	"bufio"
	"context"
	"fmt"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

var (
	labelSelector string
	nameRegex     string
	prefix        string
	podList       []v1.Pod
)

var watchPodsCmd = &cobra.Command{
	Use:   "watch-pods",
	Short: "Watch pods across owned namespaces by common labels",
	Run: func(cmd *cobra.Command, args []string) {
		main() // Implementation will go here
	},
}

func init() {
	k8Cmd.AddCommand(watchPodsCmd)
	watchPodsCmd.Flags().StringVarP(&labelSelector, "label-selector", "l", "", "Label selector to filter pods")
	watchPodsCmd.Flags().StringVarP(&nameRegex, "name-regex", "r", "", "Prefix pod names to watch")
}

func main() {

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

	WatchPodsInNamespaces(clientset, namespaces)

}

func WatchPodsInNamespaces(clientset *kubernetes.Clientset, namespaces []string) {

	for _, ns := range namespaces {
		fmt.Println("Starting pod informer for namespace:", ns)
		startPodInformer(clientset, ns)
	}

	// 🔥 BLOCK FOREVER
	select {}
}

func startPodInformer(clientset *kubernetes.Clientset, namespace string) {

	// fmt.Println("Starting pod informer for namespace:", namespace)
	options := []informers.SharedInformerOption{
		informers.WithNamespace(namespace),
	}

	if labelSelector != "" {
		options = append(options, informers.WithTweakListOptions(
			func(opts *metav1.ListOptions) {
				opts.LabelSelector = labelSelector
			},
		))
	}

	factory := informers.NewSharedInformerFactoryWithOptions(
		clientset,
		0,
		options...,
	)
	podInformer := factory.Core().V1().Pods().Informer()

	podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{

		AddFunc: func(obj interface{}) {
			pod := obj.(*v1.Pod)
			if pod.Status.Phase == v1.PodRunning {
				streamLogs(clientset, pod)
			}
		},

		UpdateFunc: func(oldObj, newObj interface{}) {
			oldPod := oldObj.(*v1.Pod)
			newPod := newObj.(*v1.Pod)

			// Detect Pending → Running transition
			if oldPod.Status.Phase != v1.PodRunning &&
				newPod.Status.Phase == v1.PodRunning {

				streamLogs(clientset, newPod)
			}
		},
	})
	stop := make(chan struct{})
	go factory.Start(stop)
}

func streamLogs(clientset *kubernetes.Clientset, pod *v1.Pod) {
	// Implementation for streaming logs from the pod
	fmt.Printf("▶️ Streaming logs for %s/%s\n", pod.Namespace, pod.Name)

	req := clientset.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &v1.PodLogOptions{
		Container:  pod.Spec.Containers[0].Name,
		Follow:     true,
		Timestamps: false,
	})

	stream, err := req.Stream(context.TODO())
	if err != nil {
		fmt.Printf("❌ Cannot stream logs: %v\n", err)
		return
	}
	defer stream.Close()

	scanner := bufio.NewScanner(stream)
	for scanner.Scan() {
		fmt.Printf("%s/%s:%s\n", pod.Namespace, pod.Name, scanner.Text())
	}

}

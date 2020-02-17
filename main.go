package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/client-go/util/retry"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	podName = flag.String("pod", "", "The name of the target pod")
	flag.Parse()

	if *path == "" {
		usage()
		os.Exit(2)
	}

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

        podClient := clientset.CoreV1().Pods(pod.Namespace)

	ephemeralContainer := &apiv1.EphemeralContainer{
		ObjectMeta: metav1.ObjectMeta{
			Name: "let-me-in",
		},
		"ephemeralContainers":[
			{
				"name":"debugger",
				"image":"busybox",
				"command":[
					"sh"
				],
				"resources":{
         			},
				"terminationMessagePolicy":"File",
				"imagePullPolicy":"IfNotPresent",
				"stdin":true,
				"tty":true
			}
		]
	}

	// Inject ephemeral container
	fmt.Println("Creating ephemeral container...")
	result, err := podClient.UpdateEphemeralContainers(context.TODO(), podName, ephemeralContainer, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())

}

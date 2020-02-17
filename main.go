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

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

        ephemeralContainerClient := clientset.ApiV1().EphemeralContainer(apiv1.NamespaceDefault)

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

	// Create Deployment
	fmt.Println("Creating ephemeral container...")
	result, err := ephemeralContainerClient.Create(context.TODO(), ephemeralContainer, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())

}

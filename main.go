package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	var podname *string = flag.String("pod", "", "The name of the target pod")
	var namespace *string = flag.String("namespace", "", "The name of the target namespace")
	flag.Parse()

	if *podname == "" {
		os.Exit(2)
	}

	if *namespace == "" {
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

        podClient := clientset.CoreV1().Pods(*namespace)
 
	ec := &apiv1.EphemeralContainers{
			ObjectMeta: v1.ObjectMeta{
				Name:                       "example-ec",
				Labels:                     {},
				Annotations:                {},
			},
			EphemeralContainers: {
				{
					EphemeralContainerCommon: v1.EphemeralContainerCommon{
						Name:                     "debugger",
						Image:                    "busybox",
						Command:                  {"sh"},
						Args:                     nil,
						WorkingDir:               "",
						Ports:                    nil,
						EnvFrom:                  nil,
						Env:                      nil,
						Resources:                v1.ResourceRequirements{},
						VolumeMounts:             nil,
						VolumeDevices:            nil,
						TerminationMessagePath:   "",
						TerminationMessagePolicy: "File",
						ImagePullPolicy:          "IfNotPresent",
						Stdin:                    true,
						StdinOnce:                false,
						TTY:                      true,
					},
					TargetContainerName: "*podname",
				},
			},
		}
	}

	// Inject ephemeral container
	fmt.Println("Creating ephemeral container...")
	result, err := podClient.UpdateEphemeralContainers(context.TODO(), *podname, ec, metav1.UpdateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())

}

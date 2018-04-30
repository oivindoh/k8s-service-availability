/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Note: the example only works with the code within the same release/branch.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var kubeconfig = flag.String("kubeconfig", "", "kubeconfig file path")

func main() {

	namespaceFlag := flag.String("namespace", "default", "namespace containing the service")
	servicenameFlag := flag.String("service", "default", "name of the service")
	flag.Parse()

	config, configerror := getConfig(*kubeconfig)
	if configerror != nil {
		panic(configerror.Error())
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("Checking for available addresses of %s in namespace %s\n", *servicenameFlag, *namespaceFlag)

	var availableAddresses int
	serviceQuery := fmt.Sprintf("metadata.name=%s", *servicenameFlag)

	for {
		endpoints, endpointserror := client.CoreV1().Endpoints(*namespaceFlag).List(metav1.ListOptions{FieldSelector: serviceQuery})
		if endpointserror != nil {
			panic(endpointserror.Error())
		}
		fmt.Printf("Service %s has %d endpoints\n", *servicenameFlag, len(endpoints.Items))
		for _, endpoint := range endpoints.Items {
			subsets := endpoint.Subsets
			fmt.Println(" Subset:")
			for _, subset := range subsets {
				fmt.Printf("  Available addresses: %d\n", len(subset.Addresses))
				fmt.Printf("  Not available addresses: %d\n", len(subset.NotReadyAddresses))
				availableAddresses += len(subset.Addresses)
			}
		}
		if availableAddresses > 0 {
			os.Exit(0)
		}
		time.Sleep(5 * time.Second)
	}
}

func getConfig(kubeconfig string) (*rest.Config, error) {
	isKubernetes := os.Getenv("KUBERNETES_SERVICE_HOST")
	if isKubernetes != "" {
		// We're in kube; we have a service account
		return rest.InClusterConfig()
	}
	if kubeconfig == "" {
		home := os.Getenv("HOME")
		kubeconfig = filepath.Join(home, ".kube", "config")
	}
	return clientcmd.BuildConfigFromFlags("", kubeconfig)
}

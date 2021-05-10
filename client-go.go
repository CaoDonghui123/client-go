package main

import (
	"context"
	"flag"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

func main() {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	// uses the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	//查询pod信息
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	for index, pod := range pods.Items {
		fmt.Printf("pod: %v  %v node:%v\n", index, pod.Name, pod.Spec.NodeName)
	}

	//查询api
	api, rs, err := discovery.ServerGroupsAndResources(clientset)
	for index, value := range rs {
		fmt.Printf("rs:%v  %v\n", index, value.GroupVersion)
	}
	for index, value := range api {
		fmt.Printf("rs:%v  %v\n", index, value.Name)
	}
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

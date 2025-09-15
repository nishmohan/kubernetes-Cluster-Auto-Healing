package main

import (
    "context"
    "fmt"
    "os"
    "path/filepath"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
)

func main() {
    // Build kubeconfig path (~/.kube/config)
    homeDir, err := os.UserHomeDir()
    if err != nil {
        panic(err)
    }
    kubeconfig := filepath.Join(homeDir, ".kube", "config")

    // Load config from kubeconfig file
    config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
    if err != nil {
        panic(err)
    }

    // Create Kubernetes clientset
    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        panic(err)
    }

    // Get pods from all namespaces
    pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
    if err != nil {
        panic(err)
    }

    fmt.Printf("Found %d pods:\n", len(pods.Items))
    for _, pod := range pods.Items {
        fmt.Printf("Namespace: %-15s Pod: %-35s Status: %s\n",
            pod.Namespace, pod.Name, pod.Status.Phase)
    }
}

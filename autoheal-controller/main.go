// main.go
package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

type actionLog struct {
	TS        string `json:"ts"`
	Event     string `json:"event"`
	Namespace string `json:"namespace"`
	Pod       string `json:"pod"`
	Action    string `json:"action"`
	Reason    string `json:"reason"`
	Outcome   string `json:"outcome"`
	Error     string `json:"error,omitempty"`
}

func jlog(a actionLog) {
	a.TS = time.Now().UTC().Format(time.RFC3339)
	b, _ := json.Marshal(a)
	log.Println(string(b))
}

func hasControllerOwner(p *corev1.Pod) bool {
	for _, o := range p.OwnerReferences {
		if o.Controller != nil && *o.Controller {
			return true
		}
	}
	return false
}

func isCrashLoop(p *corev1.Pod) bool {
	for _, cs := range p.Status.ContainerStatuses {
		if cs.State.Waiting != nil && cs.State.Waiting.Reason == "CrashLoopBackOff" && cs.RestartCount >= 3 {
			return true
		}
	}
	return false
}

func isEvicted(p *corev1.Pod) bool {
	return p.Status.Reason == "Evicted"
}

func deletePod(cs *kubernetes.Clientset, ns, name, reason string, force bool) {
	var grace *int64
	if force {
		z := int64(0)
		grace = &z
	}
	err := cs.CoreV1().Pods(ns).Delete(context.TODO(), name, metav1.DeleteOptions{GracePeriodSeconds: grace})
	if err != nil {
		jlog(actionLog{Event: "action", Namespace: ns, Pod: name, Action: "delete", Reason: reason, Outcome: "failure", Error: err.Error()})
		return
	}
	jlog(actionLog{Event: "action", Namespace: ns, Pod: name, Action: "delete", Reason: reason, Outcome: "success"})
}

func handle(cs *kubernetes.Clientset, obj interface{}) {
	pod := obj.(*corev1.Pod)
	ns := pod.Namespace
	name := pod.Name

	// Safety: act only on controller-managed pods
	if !hasControllerOwner(pod) {
		return
	}

	// CrashLoopBackOff → delete to reschedule
	if isCrashLoop(pod) {
		deletePod(cs, ns, name, "CrashLoopBackOff", false)
		return
	}

	// Evicted → cleanup
	if isEvicted(pod) {
		// small grace to avoid flapping
		if time.Since(pod.CreationTimestamp.Time) < 60*time.Second {
			return
		}
		deletePod(cs, ns, name, "Evicted", false)
		return
	}
}

func kubeConfig() (*rest.Config, error) {
	// 1) try in-cluster (normal for a Deployment)
	if cfg, err := rest.InClusterConfig(); err == nil {
		return cfg, nil
	}
	// 2) fall back to KUBECONFIG or default path (useful for local dev)
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		home := os.Getenv("HOME")
		if home != "" {
			kubeconfig = home + "/.kube/config"
		}
	}
	return clientcmd.BuildConfigFromFlags("", kubeconfig)
}

func main() {
	cfg, err := kubeConfig()
	if err != nil {
		log.Fatalf("kube config: %v", err)
	}
	cs, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		log.Fatalf("client: %v", err)
	}

	lw := cache.NewListWatchFromClient(
		cs.CoreV1().RESTClient(),
		"pods",
		corev1.NamespaceAll,
		fields.Everything(),
	)

	_, controller := cache.NewInformer(
		lw,
		&corev1.Pod{},
		0, // no periodic resync
		cache.ResourceEventHandlerFuncs{
			AddFunc:    func(obj interface{}) { handle(cs, obj) },
			UpdateFunc: func(_, newObj interface{}) { handle(cs, newObj) },
		},
	)

	stop := make(chan struct{})
	defer close(stop)
	controller.Run(stop)
}
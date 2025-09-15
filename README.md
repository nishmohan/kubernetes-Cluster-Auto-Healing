# Kubernetes Cluster Health Checker and Auto-Healing
Managing large-scale Kubernetes clusters manually is resource-intensive and error-prone. Frequent issues like pod failures, unresponsive nodes, and resource spikes can cause service degradation or outages. This project addresses these challenges by building a self-healing system with real-time monitoring, automated recovery, and dashboard visualization, reducing manual effort and improving system resilience.

 ## Project Goal
 1. Develop an automated health monitoring system for Kubernetes clusters, focusing on key metrics like node health, pod statuses, and resource utilization.
 2. Implement self-healing actions that restart failed pods, reschedule workloads, and, if necessary, trigger scaling events to balance workloads.
 3. Provide real-time alerts and notifications to inform the team of critical issues that may require manual intervention.
 4. Create a web dashboard to display real-time health status, historical data, and auto-healing logs for transparency and traceability.

 ## Technology Stack
  Language:	Go (Kubernetes client)
  Monitoring:	Prometheus, Node Exporter, Kube-State-Metrics
  Visualization:	Grafana
  Alerting:	Prometheus Alertmanager + Slack API
  Kubernetes API:	client-go (Golang client)
  Containerization:	Docker
  Deployment:	Helm or YAML manifests
  Optional Config Mgmt:	ConfigMap, Secrets, RBAC


 ## Sprint-by-Sprint Breakdown
  ### Sprint 1 – Setup
   This sprint sets the foundation for your Kubernetes Cluster Health Checker and Auto-Healing system.
  #### Pre-requisites setup
   - Setup Docker Desktop
   - Kubernetes cluster (minikube, kind, EKS and others)
   - A valid kubeconfig to authenticate
   - I installed Go to local system
   - Install Helm
  #### Kubernetes API Setup using Go
  I wrote go code as below structure:
  k8s-health-checker/
     ├── go.mod
     ├── go.sum
     └── main.go   
  Then installed required module using commands:
  ```
go get k8s.io/client-go@latest
go get k8s.io/apimachinery@latest
```
  

 

  



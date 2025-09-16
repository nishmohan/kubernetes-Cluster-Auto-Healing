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
  ```
  k8s-health-checker/
     ├── go.mod
     ├── go.sum
     └── main.go
```
  Then installed required module using commands:
  ```
go get k8s.io/client-go@latest
go get k8s.io/apimachinery@latest
```
I configured docker-desktop context as a cluster on local using following command:
```
kubectl config use-context docker-desktop
```
To install all Dependencies for Go used following command.
```
go mod tidy
go mod verify
```
Deploy Prometheus and Grafana into your Kubernetes cluster
```
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm install prometheus prometheus-community/kube-prometheus-stack --namespace monitoring --create-namespace

#Check the installation
kubectl get pods -n monitoring
```
### Access Prometheus and Grafana UIs
```
kubectl port-forward svc/prometheus-kube-prometheus-prometheus -n monitoring 9090
Visit Prometheus: http://localhost:9090

#Then in another terminal:
kubectl port-forward svc/prometheus-grafana -n monitoring 3000:80
Visit Grafana: http://localhost:3000

```
<img width="1914" height="564" alt="image" src="https://github.com/user-attachments/assets/4a13a0cc-bbee-4b08-8980-1c025229d53c" />
<img width="1761" height="809" alt="image" src="https://github.com/user-attachments/assets/73af5f77-8ef5-481a-9e55-3610ad056b21" />

 ### Sprint 2 - Health Monitoring
 I setup kube-state-metrics using commands:
 ```
helm upgrade --install kube-state-metrics prometheus-community/kube-state-metrics -n monitoring --create-namespace
kubectl -n monitoring get pods -l app.kubernetes.io/name=kube-state-metrics

#checked metrics and give result successfully
kubectl -n monitoring get svc kube-state-metrics

#setup a variable
$REL = kubectl -n monitoring get servicemonitor -o jsonpath='{.items[0].metadata.labels.release}'
```

I created kube-health-basic-alerts.yaml
After that I ran following commands to setup Alertmanager:
```
kubectl apply -f kube-health-basic-alerts.yaml
kubectl -n monitoring get prometheusrule kube-health-basic-alerts
kubectl -n monitoring port-forward svc/prometheus-kube-prometheus-alertmanager 9093:9093

```
Here is the screen shot for Alertmanager
<img width="1703" height="936" alt="image" src="https://github.com/user-attachments/assets/f053bf39-a8bc-44f9-bfcf-39c599188272" />

Then create pods and check it:
```
kubectl apply -f bad-pod.yaml
kubectl get pod crashy -w
```

 ### Sprint 3 - Pod-Level Auto-Healing
 Automate pod recovery processes and ensure that self-healing actions are logged for transparency.
 To achieve this I created autoheal-rbac.yaml and controller directory then installed dependencies:
 ```
 kubectl apply -f autoheal-rbac.yaml
 go mod init example.com/autoheal
 go get k8s.io/client-go@v0.30.0 k8s.io/api@v0.30.0 k8s.io/apimachinery@v0.30.0

 #Verify the go
  go mod tidy
  go mod verify

 #deploy the controller using commands
 kubectl apply -f autoheal-rbac.yaml
 kubectl apply -f autoheal-deploy.yaml

```
After set this autuheal process then checked it's runing or not. Also print the log using following commands:
```
kubectl -n autoheal get deploy autoheal
kubectl -n autoheal get pods -o wide
kubectl -n autoheal logs deploy/autoheal --tail=50

```
More details check below screen shots:

<img width="1104" height="65" alt="image" src="https://github.com/user-attachments/assets/7881d26a-a011-4a75-a5a0-eab556624414" />
<img width="1223" height="62" alt="image" src="https://github.com/user-attachments/assets/37042b41-a13e-46e7-a05d-d26a64a591df" />
<img width="1543" height="296" alt="image" src="https://github.com/user-attachments/assets/8faf5b0a-25ad-4e60-b533-f181f90ef9d1" />

Testing purpose I created crashy-deploy.yaml then ran commands and achive following steps:
  - Test self-healing actions in a staging environment.
  - Log each action for auditing and reporting.
```
kubectl -n default get pods -l app=crashy -w
```
Open new terminal and run command:
- Add automated resource cleanup for pods in "CrashLoopBackOff" or "Evicted" states.
```
kubectl -n autoheal logs deploy/autoheal -f | Select-String -Pattern CrashLoopBackOff
```
<img width="1548" height="345" alt="image" src="https://github.com/user-attachments/assets/10b2b263-dba5-4236-99c1-023bec46b3e7" />

After everything done I clean up using commands:
```
kubectl -n default delete deployment crashy
Remove-Item .\crashy-deploy.yaml

```
Also checked if RBAC is effective:
```
kubectl auth can-i --as=system:serviceaccount:autoheal:autoheal-sa list pods --all-namespaces
kubectl auth can-i --as=system:serviceaccount:autoheal:autoheal-sa delete pods --all-namespaces
kubectl auth can-i --as=system:serviceaccount:autoheal:autoheal-sa create events --all-namespaces
```


### Sprint 4 - Pod-Level Auto-Healing













  

 

  



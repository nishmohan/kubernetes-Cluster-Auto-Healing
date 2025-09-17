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
Ensure the system can scale and balance resources automatically, optimizing cluster performance and cost.

 #### Node Autoscaling
 - Use Cluster Autoscaler for automatic node scaling.
 - Configured it with AWS to add/remove nodes based on workload.
 - Set appropriate thresholds to avoid over-scaling or resource wastage.
Installed metrics-server:
```
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
kubectl get apiservices v1beta1.metrics.k8s.io
kubectl top nodes

```
I need to create multi-node local cluster for thi I have to install "kind" So I installed that and then ran commands:

```
kind create cluster --name autoscale-lab --config kind-3n.yaml
kubectl config use-context kind-autoscale-lab
kubectl get nodes -o wide
```
I have created EKS(AWS)- Node Autoscaler Cluster with an autoscalable node group:
```
eksctl create cluster --name autoscale-lab-mlal --region us-west-2 `
  --nodegroup-name ng1 --nodes 1 --nodes-min 1 --nodes-max 4 --managed
#then get the node
kubectl get nodes -o wide
```
You can refer the screen shots for cluster:
<img width="1914" height="894" alt="Screenshot 2025-09-17 094115" src="https://github.com/user-attachments/assets/7e51fdfd-e968-46e6-82c0-5a0e895c3ff9" />
<img width="1885" height="838" alt="Screenshot 2025-09-17 094650" src="https://github.com/user-attachments/assets/4cb9c2a8-3ad9-4820-ba7c-5badef13494f" />
<img width="1866" height="836" alt="Screenshot 2025-09-17 094731" src="https://github.com/user-attachments/assets/5cdfe37d-b6c2-4670-9093-c453ba088e74" />
<img width="1885" height="879" alt="Screenshot 2025-09-17 094752" src="https://github.com/user-attachments/assets/4909cbf2-6509-416f-8589-78328e0935c5" />
<img width="1892" height="901" alt="Screenshot 2025-09-17 094808" src="https://github.com/user-attachments/assets/4a0a6766-b115-4bb0-a0b6-a5c6867748cb" />
<img width="1863" height="863" alt="Screenshot 2025-09-17 094831" src="https://github.com/user-attachments/assets/0e7a145c-fb61-40f3-af7e-3eec859e9efb" />
<img width="1470" height="759" alt="Screenshot 2025-09-17 094854" src="https://github.com/user-attachments/assets/44f85a63-d23a-4d08-99e7-005ae68878e1" />
<img width="1516" height="736" alt="Screenshot 2025-09-17 095441" src="https://github.com/user-attachments/assets/9ced25d0-5c11-4c71-86c7-c3ae43c5c8b2" />
<img width="1531" height="760" alt="Screenshot 2025-09-17 095523" src="https://github.com/user-attachments/assets/ab1495fd-b444-49bd-b8e2-09b149e067e5" />
<img width="1910" height="789" alt="Screenshot 2025-09-17 094326" src="https://github.com/user-attachments/assets/52897113-db8c-49fd-8453-53d78e8b935b" />
<img width="1888" height="840" alt="Screenshot 2025-09-17 094602" src="https://github.com/user-attachments/assets/178c9d3e-be28-4de0-a728-e3315932ea31" />
<img width="1862" height="848" alt="Screenshot 2025-09-17 094626" src="https://github.com/user-attachments/assets/79b0fab4-cc65-4665-b809-c4e2fdb5354b" />

It's not allowing the IAM role to full access for cluster So I used EKS Pod Identity for Cluster Autoscaler and this avoids the blocked IAM action and is fully supported.

```
aws eks create-addon --cluster-name autoscale-lab-mlal --region us-west-2 --addon-name eks-pod-identity-agent
aws eks describe-addon --cluster-name autoscale-lab-mlal --region us-west-2 --addon-name eks-pod-identity-agent

#Create an IAM role for Cluster Autoscaler (Pod Identity trust)
aws iam create-role --role-name EKSCA-PodIdentity-Role --assume-role-policy-document file://ca-trust.json

```

Attach a policy that lets CA scale your ASG
```
$policyArn = aws iam create-policy --policy-name EKSCA-Policy  --policy-document file://ca-policy.json --query Policy.Arn --output text
aws iam attach-role-policy --role-name EKSCA-PodIdentity-Role --policy-arn $policyArn

```
Create the CA service account & Pod Identity association and that API is purpose-built to connect a K8s service account to an IAM role via Pod Identity.

```
# Make sure the SA exists first
kubectl -n kube-system create serviceaccount cluster-autoscaler --dry-run=client -o yaml | kubectl apply -f -

# Link the SA to your IAM role (Pod Identity association)
$roleArn = (aws iam get-role --role-name EKSCA-PodIdentity-Role --query Role.Arn --output text)
aws eks create-pod-identity-association `
  --cluster-name autoscale-lab-mlal `
  --namespace kube-system `
  --service-account cluster-autoscaler `
  --role-arn $roleArn `
  --region us-west-2
```
After that I installed the  Cluster Autoscaler and it will be targeted nodegroup directly.
```
helm upgrade --install cluster-autoscaler autoscaler/cluster-autoscaler -n kube-system -f ca-values.yaml
kubectl -n kube-system rollout status deploy/cluster-autoscaler

```
After that I setup HPA demo app for testing. 


### Sprint 5 - Alerts & Notifications
Implement a comprehensive alerting and notification system to keep the team informed of critical events and actions.

A. Integrate Slack or Teams API for real-time notifications on critical issues or auto-healing actions.
 - I did setup slack webhook: Create a Slack Incoming Webhook (browser-only) using URL: https://api.slack.com/apps
 - Click Create New App → From scratch and
    - App name: k8s-alerts
    - Pick your workspace.
  - In the left sidebar, go to Features → Incoming Webhooks and turn it ON.
  - Click Add New Webhook to Workspace, choose the channel (e.g. #k8s-alerts), then Allow.
  - Copy the generated Webhook URL (looks like https://hooks.slack.com/services/...).

    For reference check below screenshots:
    
![Untitled](https://github.com/user-attachments/assets/1ee3162e-0d5f-4521-8017-1c91cda360a2)
![Untitled-1](https://github.com/user-attachments/assets/56d9ed04-189d-483d-b5d6-5d94d9078f88)
![WhatsApp Image 2025-09-17 at 6 46 35 PM](https://github.com/user-attachments/assets/6fed3776-7dfc-4ab7-95df-57652940b57b)

Then configured the Slack Webhook URL to Docker-Dekstop:
```
$NS = "monitoring"
kubectl create namespace $NS --dry-run=client -o yaml | kubectl apply -f -

# Replace the URL with your actual Slack webhook:
$SLACK = "https://hooks.slack.com/services/ssss/ssss/zzzz"
kubectl -n $NS create secret generic am-slack-webhook `
  --from-literal=slack_api_url="$SLACK" `
  --dry-run=client -o yaml | kubectl apply -f -
```
B. Configure Alertmanager with customizable alerting rules for different severity levels.
 - Install kube-prometheus-stack (Prometheus + Alertmanager + Grafana) with Slack routing

   ```
   helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
   helm repo update

   #create kps-alerts-values.yaml
   helm upgrade --install kps prometheus-community/kube-prometheus-stack  -n $NS --create-namespace -f kps-alerts-values.yaml

   # verify pods
   kubectl -n $NS get pods

  ```
- Add useful alert rules (PrometheusRule) and create a file sprint5-rules.yaml then run following commands:
```
kubectl apply -f sprint5-rules.yaml
kubectl -n $NS get prometheusrules

```
C. Test alerts and notifications to ensure that DevOps teams receive timely updates.
If Crash any pod then it will send alert to slack channel and I used Direct message then it will send direct message.
Forcefully can do using below commands:
```
kubectl create ns s5-test
kubectl -n s5-test run badimage --image=doesnotexist:latest
# wait ~2–5 min, you should get PodCrashLooping (and maybe HighPodRestarts)
kubectl -n s5-test get pods -w
# cleanup
kubectl delete ns s5-test

```
Direct message will receive on slack:
![WhatsApp Image 2025-09-17 at 6 46 57 PM](https://github.com/user-attachments/assets/56e56240-6560-4c9c-81d9-deb18c98d001)


Here are the screen shots for Docker Desktop:
<img width="1871" height="878" alt="image" src="https://github.com/user-attachments/assets/e3e91526-7a64-4924-b6ed-ff435e46c305" />
<img width="1511" height="511" alt="image" src="https://github.com/user-attachments/assets/34c37def-750a-4e61-bc9f-7fdef70bd40e" />
<img width="1521" height="494" alt="image" src="https://github.com/user-attachments/assets/c1307247-207b-4a6e-ba39-744773f6858c" />


























  

 

  



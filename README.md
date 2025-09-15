# Kubernetes Cluster Health Checker and Auto-Healing

 ## Project Goal
 1. Develop an automated health monitoring system for Kubernetes clusters, focusing on key metrics like node health, pod statuses, and resource utilization.
 2. Implement self-healing actions that restart failed pods, reschedule workloads, and, if necessary, trigger scaling events to balance workloads.
 3. Provide real-time alerts and notifications to inform the team of critical issues that may require manual intervention.
 4. Create a web dashboard to display real-time health status, historical data, and auto-healing logs for transparency and traceability.

 ## Tool Used
- Programming Languages: Go (for performance and Kubernetes API compatibility)
- Kubernetes API: For interacting with and managing Kubernetes resources
- Prometheus: For monitoring and collecting metrics from the Kubernetes cluster
- Grafana: For visualizing metrics and health status on a real-time dashboard
- Alertmanager (Prometheus): For sending alerts to Slack or other communication tools
- Slack API: For notifications to DevOps teams
- Docker: To containerize the application and deploy it as a microservice



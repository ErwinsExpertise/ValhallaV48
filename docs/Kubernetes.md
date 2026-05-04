# Kubernetes Setup Guide

This guide covers deploying Valhalla to a Kubernetes cluster using Helm.

## Why Kubernetes?

Kubernetes provides:
- **High availability** - Automatic restarts and health checks
- **Scalability** - Easy to scale channels horizontally
- **Production-ready** - Battle-tested orchestration platform
- **Infrastructure as code** - Declarative configuration

## Prerequisites

- **Kubernetes cluster** - minikube, kind, K3s, or cloud provider (GKE, EKS, AKS)
- **kubectl** - Configured to connect to your cluster
- **Helm 3** - Package manager for Kubernetes
- **NX data** - See [Installation Guide](Installation.md) for v48 conversion
- **Container registry** - Or ability to load images directly (minikube, kind)

## Quick Start

### Step 1: Prepare Your Cluster

#### Option A: Local Development with minikube

```bash
# Install minikube
# See https://minikube.sigs.k8s.io/docs/start/

# Start cluster
minikube start --memory=4096 --cpus=2

# Verify
kubectl get nodes
```

#### Option B: Local Development with kind

```bash
# Install kind
# See https://kind.sigs.k8s.io/docs/user/quick-start/

# Create cluster
kind create cluster --name valhalla

# Verify
kubectl get nodes
```

#### Option C: Cloud Provider

Use your cloud provider's tools:
- **GKE**: `gcloud container clusters create valhalla`
- **EKS**: Use eksctl or AWS console
- **AKS**: `az aks create --name valhalla`

### Step 2: Prepare NX Data

For Kubernetes, you need to make your NX data available to pods. You have two options:

#### Option A: Prebake into Docker image
This is easiest option, simply build your Docker images with the NX data inside it. 

#### Option B: Persistent Volume (Recommended)

Create a PersistentVolume and copy your NX data to it. Example for local development:

```yaml
# data-pv.yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: data-nx-pv
spec:
  capacity:
    storage: 1Gi
  accessModes:
    - ReadOnlyMany
  hostPath:
    path: /data/valhalla
    type: DirectoryOrCreate
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: data-nx-pvc
  namespace: valhalla
spec:
  accessModes:
    - ReadOnlyMany
  resources:
    requests:
      storage: 1Gi
```

Apply and copy file:
```bash
kubectl apply -f data-pv.yaml

# For minikube
minikube ssh
sudo mkdir -p /data/valhalla
# Then copy your NX data to /data/valhalla/ using your preferred method

# For kind - mount when creating cluster
kind create cluster --config=kind-config.yaml
```

### Step 3: Build and Load the Image

```bash
# Build image
docker build -t valhalla:latest -f Dockerfile .

# Load into cluster
# For kind:
kind load docker-image valhalla:latest

# For minikube:
minikube image load valhalla:latest

# For cloud providers:
# Tag and push to your container registry
docker tag valhalla:latest gcr.io/your-project/valhalla:latest
docker push gcr.io/your-project/valhalla:latest
```

### Step 4: Deploy with Helm

```bash
# Create namespace
kubectl create namespace valhalla

# Install chart
helm install valhalla ./helm -n valhalla

# Watch pods start
kubectl get pods -n valhalla -w
```

### Step 5: Expose Services

By default, all services use ClusterIP (internal only). To access from outside:

#### Option A: Port Forwarding (Development)

```bash
# Forward login server
kubectl port-forward -n valhalla svc/login-server 8484:8484

# In another terminal, forward channels
kubectl port-forward -n valhalla svc/channel-server-1 8685:8685
kubectl port-forward -n valhalla svc/channel-server-2 8686:8686
```

#### Option B: Gateway API (Recommended for Production)

See [Exposing Services with Gateway API](#exposing-services-with-gateway-api) below.

## Helm Chart Configuration

### values.yaml

The Helm chart can be customized via `helm/values.yaml`:

```yaml
# Image configuration
image:
  repository: valhalla
  tag: latest
  pullPolicy: IfNotPresent

# Replica counts
replicaCount:
  login: 1
  world: 1
  cashshop: 1
  channels: 2

# Database configuration
database:
  address: "db"
  port: "3306"
  user: "root"
  password: "password"
  database: "maplestory"

mysql:
  host: mysql.mysql.svc.cluster.local
  port: 3306
  user: root
  password: password
  database: maplestory
  tls:
    mode: "required"
    serverName: "mydb.example.com"
    caFile: "/etc/valhalla/mysql-tls/ca.pem"
    existingSecret: "valhalla-mysql-tls"

nx:
  path: "/app/nx"
  existingClaim: "valhalla-nx"
  readOnly: true

# World settings
world:
  message: "Welcome to Valhalla!"
  ribbon: 2
  expRate: 1.0
  dropRate: 1.0
  mesosRate: 1.0

# Channel settings
channel:
  maxPop: 250

# Gateway settings
gateway:
  enabled: true
  createGateway: false
  gatewayClassName: valhalla-traefik
  publicAddress: ""
  autoDiscoverAddress: true

runtimeConfig:
  image: alpine/k8s:1.31.12
```

### Installing with Custom Values

```bash
# Create custom values file
cat > my-values.yaml <<EOF
world:
  expRate: 2.0
  dropRate: 1.5
  mesosRate: 1.5
  
channel:
  maxPop: 500
EOF

# Install with custom values
helm install valhalla ./helm -n valhalla -f my-values.yaml
```

For local non-TLS MySQL, leave `mysql.tls.mode` empty.

If you want to mount NX/WZ assets from a PVC instead of baking them into the image, set:

```yaml
nx:
  path: "/app/nx"
  existingClaim: "valhalla-nx"
  readOnly: true
```

That PVC is mounted into the login, channel, and cashshop pods, and the generated TOML gets:

```toml
[nx]
path = "/app/nx"
```

If you leave `nx.existingClaim` and `nx.path` empty, Helm does nothing and the server falls back to its normal built-in NX path resolution.

For TLS-required MySQL using system roots only:

```yaml
mysql:
  tls:
    mode: "required"
    serverName: "mydb.example.com"
```

For TLS with a custom CA or client certificate, create a Kubernetes secret containing the PEM files, mount it with `mysql.tls.existingSecret`, and point `caFile`, `certFile`, and `keyFile` at the mounted paths, for example:

```yaml
mysql:
  tls:
    mode: "required"
    serverName: "mydb.example.com"
    caFile: "/etc/valhalla/mysql-tls/ca.pem"
    certFile: "/etc/valhalla/mysql-tls/client-cert.pem"
    keyFile: "/etc/valhalla/mysql-tls/client-key.pem"
    existingSecret: "valhalla-mysql-tls"
```

### Upgrading Configuration

```bash
# Edit values.yaml or create new values file
vim helm/values.yaml

# Upgrade deployment
helm upgrade valhalla ./helm -n valhalla

# Rollback if needed
helm rollback valhalla -n valhalla
```

## Service Discovery

In Kubernetes, services use DNS names instead of IP addresses:

| Service | Docker Compose | Kubernetes |
|---------|---------------|------------|
| Login Server | `login_server` | `login-server` |
| World Server | `world_server` | `world-server` |
| Database | `db` | `db` |
| Channel 1 | `channel_server_1` | `channel-server-1` |

The Helm chart automatically adjusts configurations to use hyphens for K8s service names.

## Exposing Services with Gateway API

Gateway API can expose the TCP services MapleStory needs while keeping the Valhalla services as `ClusterIP`.

### Step 1: Use the bundled Gateway controller or bring your own

By default, this chart now installs a bundled Traefik deployment and ships the Gateway API experimental CRDs needed for `Gateway`, `GatewayClass`, and `TCPRoute`.

On a fresh cluster, `helm install` will install those CRDs automatically from `helm/crds` before creating the rest of the resources.

CRD caveat: Helm installs CRDs on first install, but does not fully manage CRD upgrades/removals. If the chart later bumps Gateway API CRD versions, treat that as an explicit cluster upgrade step.

If you prefer another controller such as kgateway, disable the bundled Traefik install and point the chart at your existing Gateway/GatewayClass and gateway service.

The runtime config init container no longer uses Bitnami. If you need to override the helper image used for LoadBalancer address discovery, set `runtimeConfig.image`.

### Step 2: Enable Gateway API in Helm values

```yaml
gateway:
  enabled: true
  createGateway: false
  name: valhalla-gateway
  gatewayClassName: valhalla-traefik
  publicAddress: ""
  autoDiscoverAddress: true
  listeners:
    login: true
    cashshop: true
    channels: true
```

If you already have a shared Gateway, keep `createGateway: false` and point `name`, `namespace`, and `gatewayClassName` at the existing Gateway.

For an external controller such as kgateway:

```yaml
controller:
  enabled: false

gateway:
  enabled: true
  createGateway: true
  name: valhalla-gateway
  gatewayClassName: kgateway
  autoDiscoverAddress: true
  discovery:
    serviceName: kgateway-proxy
    namespace: kgateway-system
```

**Important**: Each channel still needs its own TCP port/listener. The chart creates one listener and `TCPRoute` per channel based on `channel.replicas`.

### Step 3: Update Valhalla Configuration

If you already have a fixed reserved/public IPv4, you can set it directly:

```yaml
gateway:
  publicAddress: "<gateway-ipv4>"
```

If `gateway.publicAddress` is empty and `gateway.autoDiscoverAddress=true`, the channel and cashshop pods will wait for the configured gateway service's LoadBalancer address, resolve a hostname to IPv4 if needed, and render that into their startup config automatically.

`channel.clientConnectionAddress` and `cashshop.clientConnectionAddress` can still be set individually, but if left empty they follow this order:

1. `gateway.publicAddress`
2. Auto-discovered LoadBalancer service address

This is needed because the server ultimately needs a literal IPv4 address in its startup config.

Upgrade Helm deployment:
```bash
helm upgrade valhalla ./helm -n valhalla -f helm/values.yaml
```

### Step 4: Get External Address

Inspect the Gateway or gateway service status to find the published external address:

```bash
kubectl get gateway -n valhalla
kubectl describe gateway valhalla-gateway -n valhalla

# Bundled Traefik default:
kubectl get svc -n valhalla valhalla-traefik
```

### Step 5: Update MapleStory Client

Configure your client to connect to the Gateway's external IPv4 address.

## Scaling Channels

### Add More Channels

Edit `helm/values.yaml`:
```yaml
replicaCount:
  channels: 5  # Increase from 2 to 5
```

Increase the Helm channel replica count:
```yaml
channel:
  replicas: 5
```

Upgrade Valhalla:
```bash
helm upgrade valhalla ./helm -n valhalla
```

The chart will automatically add the matching Gateway listeners and `TCPRoute` resources for the extra channels.

If you use the bundled controller, extra channel ports are generated automatically from `channel.replicas` for the Traefik service, Traefik entrypoints, Gateway listeners, and `TCPRoute`s.

## Database

### Using External MySQL

For production, use a managed database service:

```yaml
# values.yaml
database:
  address: "mysql.example.com"
  port: "3306"
  user: "valhalla"
  password: "securePassword"
  database: "maplestory"
```

### Using In-Cluster MySQL

The Helm chart can deploy MySQL within the cluster (not recommended for production):

```yaml
mysql:
  enabled: true
  persistence:
    enabled: true
    size: 10Gi
```

## Monitoring

### Prometheus

Install Prometheus to scrape Valhalla metrics:

```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm install prometheus prometheus-community/kube-prometheus-stack -n monitoring --create-namespace
```

Configure ServiceMonitor:
```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: valhalla
  namespace: valhalla
spec:
  selector:
    matchLabels:
      app: valhalla
  endpoints:
    - port: metrics
      path: /metrics
```

### Grafana

Access Grafana (installed with Prometheus):
```bash
kubectl port-forward -n monitoring svc/prometheus-grafana 3000:80
```

Default credentials: admin/prom-operator

## Managing the Deployment

### View Pods

```bash
kubectl get pods -n valhalla
```

### View Logs

```bash
# Specific pod
kubectl logs -n valhalla login-server-xyz123

# Follow logs
kubectl logs -n valhalla -f login-server-xyz123

# All pods with label
kubectl logs -n valhalla -l app=channel-server
```

### Execute Commands in Pod

```bash
kubectl exec -it -n valhalla channel-server-1-xyz123 -- sh
```

### Restart Deployment

```bash
kubectl rollout restart deployment/login-server -n valhalla
```

### Scale Manually

```bash
kubectl scale deployment/channel-server --replicas=3 -n valhalla
```

## Troubleshooting

### Pods Not Starting

**Check pod status**:
```bash
kubectl describe pod -n valhalla <pod-name>
```

**Common issues**:
- Image pull error: Check image name and pull policy
- Missing NX data: Verify ConfigMap or PV is correctly mounted
- Database connection: Check database service and credentials

### CrashLoopBackOff

**Check logs**:
```bash
kubectl logs -n valhalla <pod-name> --previous
```

**Common causes**:
- Missing environment variables
- Database not ready
- Incorrect configuration

### Service Not Reachable

**Check service**:
```bash
kubectl get svc -n valhalla
kubectl describe svc -n valhalla login-server
```

**Test connectivity**:
```bash
# From inside cluster
kubectl run -it --rm debug --image=alpine --restart=Never -n valhalla -- sh
apk add curl netcat-openbsd
nc -zv login-server 8484
```

### ConfigMap/Secret Changes Not Reflected

Pods don't automatically restart when ConfigMaps/Secrets change:

```bash
# Force restart
kubectl rollout restart deployment/login-server -n valhalla
```

## Security Best Practices

1. **Use Secrets for sensitive data**:
   ```bash
   kubectl create secret generic db-credentials \
     --from-literal=password=securePassword \
     -n valhalla
   ```

2. **Set resource limits**:
   ```yaml
   resources:
     limits:
       memory: "512Mi"
       cpu: "500m"
     requests:
       memory: "256Mi"
       cpu: "250m"
   ```

3. **Use RBAC** for access control
4. **Enable Network Policies** to restrict traffic
5. **Run as non-root user** where possible
6. **Keep images updated** regularly

## Backup and Recovery

### Backup Database

```bash
# If using in-cluster MySQL
kubectl exec -n valhalla db-0 -- mysqldump -u root -ppassword maplestory > backup.sql

# If using PVC
kubectl exec -n valhalla db-0 -- mysqldump -u root -ppassword maplestory | gzip > backup.sql.gz
```

### Restore Database

```bash
kubectl exec -i -n valhalla db-0 -- mysql -u root -ppassword maplestory < backup.sql
```

## Production Checklist

- [ ] Use managed database service
- [ ] Set up SSL/TLS certificates
- [ ] Configure resource requests and limits
- [ ] Set up monitoring and alerting
- [ ] Configure automatic backups
- [ ] Use Secrets for sensitive data
- [ ] Set up logging aggregation
- [ ] Configure pod disruption budgets
- [ ] Test disaster recovery procedures
- [ ] Set up autoscaling (HPA) if needed
- [ ] Configure network policies
- [ ] Use separate namespaces for different environments

## Next Steps

- Configure server settings: [Configuration.md](Configuration.md)
- Learn about Docker deployment: [Docker.md](Docker.md)
- Understand local development: [Local.md](Local.md)
- Build from source: [Building.md](Building.md)

## Useful Commands Reference

```bash
# Deploy
helm install valhalla ./helm -n valhalla

# Upgrade
helm upgrade valhalla ./helm -n valhalla

# Rollback
helm rollback valhalla -n valhalla

# Uninstall
helm uninstall valhalla -n valhalla

# View values
helm get values valhalla -n valhalla

# Check status
helm status valhalla -n valhalla

# View pods
kubectl get pods -n valhalla

# View services
kubectl get svc -n valhalla

# View logs
kubectl logs -f -n valhalla <pod-name>

# Port forward
kubectl port-forward -n valhalla svc/login-server 8484:8484

# Execute in pod
kubectl exec -it -n valhalla <pod-name> -- sh
```

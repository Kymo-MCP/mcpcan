# Helm Chart for MCPCan

[![Version](https://img.shields.io/badge/version-v1.0.0--dev-blue.svg)](https://github.com/Kymo-MCP/mcpcan)
[![License](https://img.shields.io/badge/license-Apache%202.0-green.svg)](https://github.com/Kymo-MCP/mcpcan/blob/main/LICENSE)
[![Kubernetes](https://img.shields.io/badge/kubernetes-1.20%2B-blue.svg)](https://kubernetes.io/)
[![Helm](https://img.shields.io/badge/helm-3.0%2B-blue.svg)](https://helm.sh/)

MCPCan (Microservices Container Platform) is a comprehensive microservices platform that provides code package management, authorization services, API gateway, and web interface. This Helm chart enables easy deployment and management of the entire MCPCan ecosystem on Kubernetes.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Components](#components)
- [Usage](#usage)
- [Upgrading](#upgrading)
- [Uninstalling](#uninstalling)
- [Configuration Examples](#configuration-examples)
- [Troubleshooting](#troubleshooting)
- [Contributing](#contributing)

## Prerequisites

- Kubernetes 1.20+
- Helm 3.0+
- NGINX Ingress Controller (if ingress is enabled)
- Persistent storage (for data persistence)
- At least 2GB RAM and 2 CPU cores available in the cluster

## Installation

### Quick Start

```bash
# Add the MCPCan Helm repository
helm repo add mcpcan https://kymo-mcp.github.io/mcpcan/deploy/

# Update repository
helm repo update

# Install MCPCan with default configuration
helm install mcpcan mcpcan/mcpcan

# Or install from local chart
helm install mcpcan ./helm/
```

### Custom Installation

```bash
# Install with custom values
helm install mcpcan ./helm/ -f custom-values.yaml

# Install in specific namespace
helm install mcpcan ./helm/ --namespace mcp-system --create-namespace

# Install with inline value overrides
helm install mcpcan ./helm/ \
  --set global.domain=mcp.example.com \
  --set global.publicIP=192.168.1.100 \
  --set infrastructure.mysql.auth.rootPassword=secure-password
```

## Configuration

The following table lists the configurable parameters and their default values:

### Global Configuration

| Parameter                | Description                           | Default                               |
| ------------------------ | ------------------------------------- | ------------------------------------- |
| `global.domain`          | Domain name for the application       | `""`                                  |
| `global.publicIP`        | Public IP address for external access | `"192.168.1.100"`                     |
| `global.version`         | Application version tag               | `v2.1.4`            |
| `global.registry`        | Container registry URL                | `ccr.ccs.tencentyun.com/itqm-private` |
| `global.imagePullPolicy` | Image pull policy                     | `Always`                              |
| `global.appSecret`       | Application secret key                | `dev-app-secret`                      |

### Infrastructure Configuration

#### MySQL Configuration

| Parameter                                | Description                   | Default             |
| ---------------------------------------- | ----------------------------- | ------------------- |
| `infrastructure.mysql.enabled`           | Enable MySQL deployment       | `true`              |
| `infrastructure.mysql.auth.rootPassword` | MySQL root password           | `dev-root-password` |
| `infrastructure.mysql.auth.database`     | Default database name         | `mcp_dev`           |
| `infrastructure.mysql.auth.username`     | Application database username | `mcp_user`          |
| `infrastructure.mysql.auth.password`     | Application database password | `dev-password`      |
| `infrastructure.mysql.image.repository`  | MySQL image repository        | `mysql`             |
| `infrastructure.mysql.image.tag`         | MySQL version tag             | `"8.0.30"`          |
| `infrastructure.mysql.service.port`      | MySQL service port            | `3306`              |
| `infrastructure.mysql.service.nodePort`  | MySQL NodePort                | `31306`             |

#### Redis Configuration

| Parameter                               | Description             | Default              |
| --------------------------------------- | ----------------------- | -------------------- |
| `infrastructure.redis.enabled`          | Enable Redis deployment | `true`               |
| `infrastructure.redis.auth.password`    | Redis password          | `dev-redis-password` |
| `infrastructure.redis.auth.db`          | Default Redis database  | `0`                  |
| `infrastructure.redis.image.repository` | Redis image repository  | `redis`              |
| `infrastructure.redis.image.tag`        | Redis version tag       | `"6-alpine"`         |
| `infrastructure.redis.service.port`     | Redis service port      | `6379`               |
| `infrastructure.redis.service.nodePort` | Redis NodePort          | `31379`              |

### Services Configuration

#### Web Service

| Parameter                      | Description        | Default |
| ------------------------------ | ------------------ | ------- |
| `services.web.enabled`         | Enable web service | `true`  |
| `services.web.replicas`        | Number of replicas | `1`     |
| `services.web.service.port`    | Service port       | `3000`  |
| `services.web.ingress.enabled` | Enable ingress     | `true`  |
| `services.web.ingress.path`    | Ingress path       | `/`     |

#### Market Service

| Parameter                         | Description           | Default            |
| --------------------------------- | --------------------- | ------------------ |
| `services.market.enabled`         | Enable market service | `true`             |
| `services.market.replicas`        | Number of replicas    | `1`                |
| `services.market.service.port`    | Service port          | `8080`             |
| `services.market.ingress.enabled` | Enable ingress        | `true`             |
| `services.market.ingress.path`    | Ingress path          | `/api/market/(.*)` |

#### Authorization Service

| Parameter                        | Description                  | Default           |
| -------------------------------- | ---------------------------- | ----------------- |
| `services.authz.enabled`         | Enable authorization service | `true`            |
| `services.authz.replicas`        | Number of replicas           | `1`               |
| `services.authz.service.port`    | Service port                 | `8081`            |
| `services.authz.ingress.enabled` | Enable ingress               | `true`            |
| `services.authz.ingress.path`    | Ingress path                 | `/api/authz/(.*)` |



### Ingress Configuration

| Parameter                | Description        | Default      |
| ------------------------ | ------------------ | ------------ |
| `ingress.enabled`        | Enable ingress     | `true`       |
| `ingress.className`      | Ingress class name | `nginx`      |
| `ingress.tls.enabled`    | Enable TLS         | `false`      |
| `ingress.tls.secretName` | TLS secret name    | `domain-tls` |

## Components

MCPCan consists of the following components:

### Infrastructure Services

- **MySQL 8.0**: Primary database for application data
- **Redis 6**: Caching and session storage

### Application Services

- **Web Frontend**: Vue.js-based user interface (Port: 3000)
- **Market Service**: Code package management and marketplace (Port: 8080)
- **Authorization Service**: User authentication and authorization (Port: 8081)

### Supporting Components

- **ConfigMaps**: Configuration management
- **Secrets**: Sensitive data management
- **Ingress**: External access routing
- **Persistent Volumes**: Data persistence

## Usage

### Accessing the Application

After installation, you can access MCPCan through:

1. **Domain-based access** (if domain is configured):

   ```
   https://your-domain.com
   ```

2. **IP-based access**:

   ```
   http://your-public-ip
   ```

3. **Port forwarding** (for testing):
   ```bash
   kubectl port-forward svc/mcp-web-svc 3000:3000
   ```

### Service Endpoints

- **Web Interface**: `/`
- **Market API**: `/api/market/*`
- **Authorization API**: `/api/authz/*`

### Database Access

For direct database access (development only):

```bash
# MySQL
kubectl port-forward svc/mysql-svc 3306:3306
mysql -h localhost -P 3306 -u mcp_user -p

# Redis
kubectl port-forward svc/redis-svc 6379:6379
redis-cli -h localhost -p 6379
```

## Upgrading

### Upgrade the Chart

```bash
# Upgrade to latest version
helm upgrade mcpcan ./helm/

# Upgrade with new values
helm upgrade mcpcan ./helm/ -f new-values.yaml

# Upgrade with specific version
helm upgrade mcpcan ./helm/ --version v2.1.4
```

### Rolling Updates

The chart supports rolling updates for all services. To update a specific service:

```bash
# Update web service image
helm upgrade mcpcan ./helm/ --set global.version=v2.1.4

# Update with zero downtime
helm upgrade mcpcan ./helm/ --wait --timeout=600s
```

## Uninstalling

```bash
# Uninstall the release
helm uninstall mcpcan

# Uninstall from specific namespace
helm uninstall mcpcan --namespace mcp-system

# Remove persistent data (optional)
kubectl delete pvc -l app.kubernetes.io/instance=mcpcan
```

## Configuration Examples

### Production Configuration

```yaml
# production-values.yaml
global:
  domain: "mcp.company.com"
  appSecret: "production-secret-key-change-me"
  imagePullPolicy: IfNotPresent

infrastructure:
  mysql:
    auth:
      rootPassword: "secure-mysql-root-password"
      password: "secure-mysql-password"
    resources:
      requests:
        memory: "512Mi"
        cpu: "200m"
      limits:
        memory: "1Gi"
        cpu: "1000m"

  redis:
    auth:
      password: "secure-redis-password"
    resources:
      requests:
        memory: "256Mi"
        cpu: "100m"
      limits:
        memory: "512Mi"
        cpu: "500m"

services:
  web:
    replicas: 2
    resources:
      requests:
        memory: "128Mi"
        cpu: "100m"
      limits:
        memory: "256Mi"
        cpu: "500m"

ingress:
  tls:
    enabled: true
    secretName: "mcp-tls-secret"
```

### Development Configuration

```yaml
# development-values.yaml
global:
  publicIP: "localhost"
  imagePullPolicy: Always

infrastructure:
  mysql:
    service:
      enabledNodePort: true
  redis:
    service:
      enabledNodePort: true

ingress:
  enabled: false
```

### High Availability Configuration

```yaml
# ha-values.yaml
services:
  web:
    replicas: 3
  market:
    replicas: 2
  authz:
    replicas: 2

infrastructure:
  mysql:
    resources:
      requests:
        memory: "1Gi"
        cpu: "500m"
      limits:
        memory: "2Gi"
        cpu: "1000m"
```

## Troubleshooting

### Common Issues

#### 1. Pods in Pending State

**Problem**: Pods remain in `Pending` state.

**Solution**:

```bash
# Check node resources
kubectl describe nodes

# Check pod events
kubectl describe pod <pod-name>

# Check storage class
kubectl get storageclass
```

#### 2. Database Connection Issues

**Problem**: Services cannot connect to MySQL/Redis.

**Solution**:

```bash
# Check service endpoints
kubectl get endpoints

# Test database connectivity
kubectl run mysql-test --image=mysql:8.0 --rm -it --restart=Never -- \
  mysql -h mysql-svc -u mcp_user -p

# Check service logs
kubectl logs -l app=mcp-market
```

#### 3. Ingress Not Working

**Problem**: Cannot access services through ingress.

**Solution**:

```bash
# Check ingress controller
kubectl get pods -n ingress-nginx

# Check ingress configuration
kubectl describe ingress mcpcan-ingress

# Check ingress logs
kubectl logs -n ingress-nginx -l app.kubernetes.io/name=ingress-nginx
```

#### 4. Image Pull Errors

**Problem**: Cannot pull container images.

**Solution**:

```bash
# Check image pull secrets
kubectl get secrets

# Test image pull manually
kubectl run test --image=<your-registry>/image:tag --rm -it --restart=Never

# Update registry credentials
kubectl create secret docker-registry regcred \
  --docker-server=<your-registry> \
  --docker-username=<username> \
  --docker-password=<password>
```

### Debugging Commands

```bash
# Check all resources
kubectl get all -l app.kubernetes.io/instance=mcpcan

# Check pod logs
kubectl logs -l app=mcp-web --tail=100

# Check resource usage
kubectl top pods

# Get detailed pod information
kubectl describe pod <pod-name>

# Check persistent volumes
kubectl get pv,pvc

# Check configuration
helm get values mcpcan
```

### Performance Tuning

#### Resource Optimization

```yaml
# Adjust resource limits based on usage
services:
  web:
    resources:
      requests:
        memory: "64Mi"
        cpu: "50m"
      limits:
        memory: "256Mi"
        cpu: "500m"
```

#### Database Optimization

```yaml
# MySQL configuration
infrastructure:
  mysql:
    resources:
      requests:
        memory: "512Mi"
        cpu: "250m"
      limits:
        memory: "1Gi"
        cpu: "1000m"
```

## Contributing

We welcome contributions! Please see our [Contributing Guide](../CONTRIBUTING.md) for details.

### Development Setup

```bash
# Clone the repository
git clone https://github.com/Kymo-MCP/mcpcan.git

# Navigate to helm directory
cd mcpcan/deploy/helm

# Validate the chart
helm lint .

# Test installation
helm install test-release . --dry-run --debug
```

### Reporting Issues

Please report issues on our [GitHub Issues](https://github.com/Kymo-MCP/mcpcan/issues) page.

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](../LICENSE) file for details.

## Support

- **Documentation**: [GitHub Wiki](https://github.com/Kymo-MCP/mcpcan/wiki)
- **Issues**: [GitHub Issues](https://github.com/Kymo-MCP/mcpcan/issues)
- **Email**: opensource@kymo.cn

---

**Note**: This chart is designed for Kubernetes environments. For Docker Compose deployment, please refer to the documentation in this same directory.

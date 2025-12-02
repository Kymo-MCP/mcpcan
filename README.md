
<div align="center">
  <img width="1872" height="932" alt="image" src="https://github.com/user-attachments/assets/a6d8796f-435d-48fd-9833-264a94a5ca1d" />
</div>

<div align="center">

# MCP CAN
The open source integration platform for MCP Server.</br>
MCPCAN uses containers for flexible deployment of MCP services, resolving potential system configuration conflicts. It supports multi-protocol compatibility and conversion, enabling seamless integration between different MCP service architectures. It also provides visual monitoring, security authentication, and one-stop deployment capabilities.</br>

  <img src="https://img.shields.io/badge/Vue.js-3.2.47-4FC08D?style=for-the-badge&logo=vue.js&logoColor=ffffff" alt="Vue.js"/>
  <img src="https://img.shields.io/badge/TypeScript-5.0-blue?style=for-the-badge&logo=typescript" alt="TypeScript"/>
  <img src="https://img.shields.io/badge/MySQL-8.0-blue?style=for-the-badge&logo=mysql" alt="MySQL"/>
  <img src="https://img.shields.io/badge/Kubernetes-1.28-326ce5?style=for-the-badge&logo=kubernetes" alt="Kubernetes"/>
  <img src="https://img.shields.io/badge/License-GPL--3.0-blue?style=for-the-badge" alt="GPL-3.0"/>
<p align="center">
    <a href="https://demo.mcpcan.com" target="_blank">
        <img alt="Static Badge" src="https://img.shields.io/badge/Product-F04438"></a>
    <a href="https://dify.ai/pricing" target="_blank">
        <img alt="Static Badge" src="https://img.shields.io/badge/free-pricing?logo=free&color=%20%23155EEF&label=pricing&labelColor=%20%23528bff"></a>
    <a href="https://discord.gg/EegGj7G7Bz" target="_blank">
        <img src="https://img.shields.io/discord/1428637640856571995?logo=discord&labelColor=%20%235462eb&logoColor=%20%23f5f5f5&color=%20%235462eb"
            alt="chat on Discord"></a>
    <a href="https://twitter.com/intent/follow?screen_name=MCPCAN" target="_blank">
        <img src="https://img.shields.io/twitter/follow/MCPCAN?logo=X&color=%20%23f5f5f5"
            alt="follow on X(Twitter)"></a>
</p>

<p align="center">
   <strong>English</strong> | <a href="./README_CN.md">中文版</a> <br>
</p>
</div>

## 🎉Introduction
MCPCan is an open-source platform focused on efficient management of MCP (Model Context Protocol) services,Providing DevOps and development teams with comprehensive MCP service lifecycle management capabilities through a modern web interface.
MCPCan supports multi-protocol compatibility and conversion, enabling seamless integration between different MCP service architectures while providing visual monitoring, security authentication, and one-stop deployment capabilities.<br/>

## ✨ Key Features

- **🎯 Unified Management**: Centralized management of all MCP service instances and configurations
- **🔄 Protocol Conversion**: Supports seamless conversion between various MCP protocols
- **📊 Real-time Monitoring**: Provides detailed service status and performance monitoring
- **🔐 Security & Authentication**: Built-in identity authentication and permission management system
- **🚀 One-stop Deployment**: Quick release, configuration, and distribution of MCP services
- **📈 Scalability**: Cloud-native architecture based on Kubernetes

## ✨Demon & MainSite
For the best demo experience, try  <a href="https://demo.mcpcan.com">DemoSite : demo.mcpcan.com</a> directly.<br>
[MP4]<br>
To view our official website address, simply click <a href="https://www.mcpcan.com">MainSite : www.mcpcan.com</a>.
## 👨‍🚀Quickstart

For detailed deployment instructions, please refer to our [Deployment Guide](https://kymo-mcp.github.io/mcpcan-deploy/).

### 1. Get the Deployment Repository

```bash
# GitHub (International)
git clone https://github.com/Kymo-MCP/mcpcan-deploy.git
cd mcpcan-deploy

# Gitee (Recommended for China)
git clone https://gitee.com/kymomcp/mcpcan-deploy.git
cd mcpcan-deploy
```

### 2. Installation

**Fast Install (Recommended)**

Suitable for clean Linux servers. Automatically installs k3s, ingress-nginx, Helm, and deploys the MCPCAN platform.

```bash
# Standard Install (International Mirrors)
./scripts/install-fast.sh

# Accelerated Install (China Mirrors)
./scripts/install-fast.sh --cn
```

Once successful, access `http://<Your Public IP>` to start using it.

**Custom Install (Helm)**

Suitable for scenarios requiring custom domains, HTTPS, or modifying default configurations.

```bash
# 1. Install dependencies (Skip if k3s/Helm are already installed)
./scripts/install-run-environment.sh       # International Mirrors
# ./scripts/install-run-environment.sh --cn  # China Mirrors

# 2. Copy and modify configuration
cp helm/values.yaml helm/values-custom.yaml
# Edit helm/values-custom.yaml to set global.domain and other parameters

# 3. Install the platform
helm install mcpcan ./helm -f helm/values-custom.yaml \
  --namespace mcpcan --create-namespace --timeout 600s --wait
```

## 🚀Components

MCPCan consists of multiple key components, which collectively form the functional framework of MCPCan, providing users with comprehensive MCP service management capabilities.

| Project                                | Status                                                      | Description                                |
| -------------------------------------- | ----------------------------------------------------------- | ------------------------------------------ |
| [MCPCan-Web](frontend/)                | ![Status](https://img.shields.io/badge/status-active-green) | MCPCan Web UI (Vue.js Frontend)            |
| [MCPCan-Backend](backend/)             | ![Status](https://img.shields.io/badge/status-active-green) | MCPCan Backend Services (Go Microservices) |
| [MCPCan-Gateway](backend/cmd/gateway/) | ![Status](https://img.shields.io/badge/status-active-green) | MCP Gateway Service                        |
| [MCPCan-Market](backend/cmd/market/)   | ![Status](https://img.shields.io/badge/status-active-green) | MCP Service Marketplace                    |
| [MCPCan-Authz](backend/cmd/authz/)     | ![Status](https://img.shields.io/badge/status-active-green) | Authentication and Authorization Service   |

## 🐧Technology Stack

### 🐧Frontend

- **Framework**: Vue.js 3.5+ (Composition API)
- **Language**: TypeScript
- **Styling**: UnoCSS, SCSS
- **UI Components**: Element Plus
- **State Management**: Pinia
- **Build Tool**: Vite

### 🐧Backend

- **Language**: Go 1.24.2+
- **Framework**: Gin, gRPC
- **Database**: MySQL, Redis
- **Container**: Docker, Kubernetes

## 🐧Third-party Projects

- [mcpcan-deploy](https://github.com/Kymo-MCP/mcpcan-deploy) - Official Helm charts source repository for MCPCan
- [MCPCan Helm Charts](https://kymo-mcp.github.io/mcpcan-deploy/) - Official Helm charts repository for MCPCan

## 💝Contributing💝

Welcome to submit PR to contribute. Please refer to [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

Before contributing, please:

1. Read our [Code of Conduct](CODE_OF_CONDUCT.md)
2. Check existing issues and pull requests
3. Follow our coding standards and commit message conventions

## Security

If you discover a security vulnerability, please refer to our [Security Policy](SECURITY.md) for responsible disclosure guidelines.

## License

Copyright (c) 2024-2025 MCPCan Team, All rights reserved.

Licensed under The GNU General Public License version 3 (GPLv3) (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at

https://www.gnu.org/licenses/gpl-3.0.html

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

## Community & Support

- 📖 [Documentation](https://kymo-mcp.github.io/mcpcan-deploy/)
- 💬 [Discord Community](https://discord.com/channels/1428637640856571995/1428637896532820038)
- 🐛 [Issue Tracker](https://github.com/Kymo-MCP/mcpcan/issues)
- 📧 [Mailing List](mailto:opensource@kymo.cn)

## Acknowledgements

- Thanks to the [MCP Protocol](https://modelcontextprotocol.io/) community
- Thanks to all contributors and supporters
- Special thanks to the open-source projects that make MCPCan possible

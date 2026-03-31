<div align="center">
  <img width="1872" height="932" alt="image" src="images/image.png" />

</div>


<div align="center">

# MCP CAN

The open source integration platform for MCP Server.</br>
MCPCAN uses containers for flexible deployment of MCP services, resolving potential system configuration conflicts. It supports multi-protocol compatibility and conversion, enabling seamless integration between different MCP service architectures. It also provides visual monitoring, security authentication, and one-stop deployment capabilities.</br>

  <img src="https://img.shields.io/badge/Vue.js-3.2.47-4FC08D?style=for-the-badge&logo=vue.js&logoColor=ffffff" alt="Vue.js"/>
  <img src="https://img.shields.io/badge/TypeScript-5.0-blue?style=for-the-badge&logo=typescript" alt="TypeScript"/>
  <img src="https://img.shields.io/badge/MySQL-8.0-blue?style=for-the-badge&logo=mysql" alt="MySQL"/>
  <img src="https://img.shields.io/badge/Kubernetes-1.28-326ce5?style=for-the-badge&logo=kubernetes" alt="Kubernetes"/>
  <img src="https://img.shields.io/badge/License-Sustainable%20Use-orange?style=for-the-badge" alt="Sustainable Use License"/>
</div>
<p align="center">
   <strong>English</strong> | <a href="./README_CN.md">中文版</a> <br>
   <a href="https://demo.mcpcan.com">DemoSite : demo.mcpcan.com（login: admin/admin123）</a> | <a href="https://www.mcpcan.com">MainSite : www.mcpcan.com</a><br>
   <a href="https://www.mcpcan.com/docs/en/guide/welcome
   "><u>Document</a></u>
</p>
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

MCPCan is an open-source platform focused on efficient management of MCP (Model Context Protocol) services, providing DevOps and development teams with comprehensive MCP service lifecycle management capabilities through a modern web interface.
MCPCan supports multi-protocol compatibility and conversion, enabling seamless integration between different MCP service architectures while providing visual monitoring, security authentication, and one-stop deployment capabilities.

## 💡 Introduction

MCPCan is an open-source platform focused on efficient management of MCP (Model Context Protocol) services, providing DevOps and development teams with comprehensive MCP service lifecycle management capabilities through a modern web interface.
MCPCan supports multi-protocol compatibility and conversion, enabling seamless integration between different MCP service architectures while providing visual monitoring, security authentication, and one-stop deployment capabilities.<br/>

## ✨ Core Features

- **🎯 Unified Management**: Centralized management of all MCP service instances and configurations
- **🔄 Protocol Conversion**: Supports seamless conversion between various MCP protocols
- **📊 Real-time Monitoring**: Provides detailed service status and performance monitoring data
- **🔐 Security Authentication**: Built-in identity authentication and permission management system
- **🚀 One-stop Deployment**: Quick release, configuration, and distribution of MCP services
- **📈 Scalability**: Cloud-native architecture based on Kubernetes

## ✨ Demo and Official Website

For the best demo experience, try directly <a href="https://demo.mcpcan.com">DemoSite : demo.mcpcan.com（login: admin/admin123）</a>.<br>
</video>
Watch our demo video on Bilibili: <a href="https://www.bilibili.com/video/BV1htBXBbECr?t=3.2">BV1htBXBbECr</a><br>
To view our official website address, simply click <a href="https://www.mcpcan.com">MainSite : www.mcpcan.com</a>.

## 👨‍🚀 Quick Start

For detailed deployment instructions, please refer to our [Deployment Guide](https://www.mcpcan.com/docs/en/guide/install).

### 1. Get Deployment Repository

```bash
# Recommended: GitHub
git clone https://github.com/Kymo-MCP/mcpcan.git
cd mcpcan/deploy/docker-compose/
```

### 2. Start Services

**Docker Compose Quick Start (Recommended)**

Suitable for local development, testing, and lightweight production deployments.

```bash
# 1. Initialize configuration
cp example.env .env
# (Optional) Modify .env file for settings like REGISTRY_PREFIX to switch between global/CN mirrors

# 2. Generate final configuration
chmod +x replace.sh
./replace.sh

# 3. Start services
docker compose up -d

# 4. Access Web UI
Login: admin/admin123
```

After successful installation, access `http://localhost` (or `http://<Your Public IP>`) to start using.

**Helm Installation**

Suitable for Kubernetes environment deployment, please refer to [Helm Deployment Guide](https://kymo-mcp.github.io/deploy/).

## 🚀 Components

MCPCan adopts a microservices architecture, consisting of the following core components:

| Component Name                         | Directory Path              | Functional Description               |
| -------------------------------------- | --------------------------- | ------------------------------------ |
| **MCPCan-Web**                         | `frontend/`                 | Management UI based on Vue 3         |
| **MCP-Market**                         | `backend/cmd/market/`       | Core business: Marketplace & Instance management |
| **MCP-Authz**                          | `backend/cmd/authz/`        | RBAC, User & Dept management         |

## 🐧 Technology Stack

### 🐧 Frontend

- **Framework**: Vue 3.5.x (Composition API)
- **Build Tool**: Vite 7.0
- **UI Components**: Element Plus
- **State Management**: Pinia 3.0
- **Styling**: UnoCSS, SCSS
- **Editor**: Monaco Editor, Markdown-it

### 🐧 Backend

- **Language**: Go 1.25.x
- **Web Framework**: Gin, gRPC (grpc-gateway)
- **Data Storage**: MySQL (GORM), Redis (Cache/Token)
- **Core Libraries**: LangChainGo, MCP-Go, Docker SDK, Client-go
- **Ops & Deployment**: Traefik (Edge Gateway), Docker Compose, Helm (Kubernetes)

## 🐧 Related Projects

- [MCPCan Helm Charts](https://kymo-mcp.github.io/deploy/) - Official Helm charts index for MCPCan
- [Deployment Guide](docs/en/guide/install.md) - Detailed deployment and configuration guide

## 💝 Contributing Guide

Welcome to submit PR to contribute! Please refer to [Contributing](CONTRIBUTING.md) for detailed guidelines.

Before contributing, please ensure:

1. Read our [Code of Conduct](CODE_OF_CONDUCT.md)
2. Check existing issues and pull requests (avoid duplicate work)
3. Follow our coding standards and commit message conventions

## ✅ Security

If you discover a security vulnerability, please refer to our [Security Policy](SECURITY.md) for responsible disclosure guidelines.

## 📄 License

Copyright (c) 2024-2025 MCPCan Team, All rights reserved.

Licensed under the **Sustainable Use License**. Portions of this software (specifically Enterprise features in `.ee` folders) are subject to the **Enterprise License** (LICENSE_EE.md). 

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

## 👥 Community & Support

- 📖 [Documentation](https://kymo-mcp.github.io/deploy/)
- 💬 [Discord Community](https://discord.com/channels/1428637640856571995/1428637896532820038)
- 🐛 [Issue Tracker](https://github.com/Kymo-MCP/mcpcan/issues)
- 📧 [Mailing List](mailto:opensource@kymo.cn)
- 🌐 WeChat<br>
  <img src="images/WeChat group QR code.jpg" alt="alt text" width="170">

## 💕 Acknowledgments

- Thanks to the [MCP Protocol](https://modelcontextprotocol.io/) community
- Thanks to all contributors and supporters
- Special thanks to the open-source projects that make MCPCan possible

## 🌟 Star History

[![Star History Chart](https://api.star-history.com/svg?repos=Kymo-MCP/mcpcan&type=date&legend=top-left)](https://www.star-history.com/#Kymo-MCP/mcpcan&type=date&legend=top-left)

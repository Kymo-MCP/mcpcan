<div align="center">
  <img width="1872" height="932" alt="image" src="images/image_cn.png" />
</div>
<div align="center">

# MCP CAN

开源的 MCP 服务器集成平台。</br>
MCPCAN 使用容器实现 MCP 服务的灵活部署，解决潜在的系统配置冲突。它支持多协议兼容与转换，实现不同 MCP 服务架构之间的无缝集成。它还提供可视化监控、安全认证和一站式部署功能。</br>

  <img src="https://img.shields.io/badge/Vue.js-3.2.47-4FC08D?style=for-the-badge&logo=vue.js&logoColor=ffffff" alt="Vue.js"/>
  <img src="https://img.shields.io/badge/TypeScript-5.0-blue?style=for-the-badge&logo=typescript" alt="TypeScript"/>
  <img src="https://img.shields.io/badge/MySQL-8.0-blue?style=for-the-badge&logo=mysql" alt="MySQL"/>
  <img src="https://img.shields.io/badge/Kubernetes-1.28-326ce5?style=for-the-badge&logo=kubernetes" alt="Kubernetes"/>
  <img src="https://img.shields.io/badge/License-Apache%202.0-blue?style=for-the-badge" alt="Apache 2.0"/>
</div>
<p align="center">
   <a href="./README.md">English</a> | <strong>中文版</strong> <br>
   <a href="https://demo.mcpcan.com">DemoSite : demo.mcpcan.com</a> | <a href="https://www.mcpcan.com">MainSite : www.mcpcan.com</a><br>
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

MCPCan 是一个专注于高效管理 MCP（模型上下文协议）服务的开源平台，通过现代化的 Web 界面，为 DevOps 和开发团队提供全面的 MCP 服务生命周期管理功能。
MCPCan 支持多协议兼容和转换，实现不同 MCP 服务架构之间的无缝集成，同时提供可视化监控、安全认证和一站式部署功能。

## 💡 介绍

MCPCan 是一个专注于高效管理 MCP（模型上下文协议）服务的开源平台，通过现代化的 Web 界面，为 DevOps 和开发团队提供全面的 MCP 服务生命周期管理功能。
MCPCan 支持多协议兼容和转换，实现不同 MCP 服务架构之间的无缝集成，同时提供可视化监控、安全认证和一站式部署功能。<br/>

## ✨ 核心功能

- **🎯 统一管理**：集中管理所有 MCP 服务实例及配置项
- **🔄 协议转换**：支持多种 MCP 协议间无缝互转
- **📊 实时监控**：提供详尽的服务状态与性能监控数据
- **🔐 安全认证**：内置身份认证与权限管理体系
- **🚀 一站式部署**：MCP 服务快速发布、配置与分发
- **📈 可扩展性**：基于 Kubernetes 的云原生架构

## ✨ 演示和官网

为了获得最佳演示体验，请尝试直接 <a href="https://demo.mcpcan.com">DemoSite : demo.mcpcan.com</a><br>

要查看我们的官方网站地址，只需点击 <a href="https://www.mcpcan.com">MainSite : www.mcpcan.com</a>。

## 👨‍🚀 快速开始

有关详细部署说明，请参阅我们的[部署指南](https://kymo-mcp.github.io/mcpcan-deploy/)。

### 1. 获取部署仓库

```bash
# GitHub (国际)
git clone https://github.com/Kymo-MCP/mcpcan-deploy.git
cd mcpcan-deploy

# Gitee (中国推荐)
git clone https://gitee.com/kymomcp/mcpcan-deploy.git
cd mcpcan-deploy
```

### 2. 安装

**快速安装（推荐）**

适用于干净的 Linux 服务器。自动安装 k3s、ingress-nginx、Helm，并部署 MCPCAN 平台。

```bash
# 标准安装（国际镜像）
./scripts/install-fast.sh

# 加速安装（中国镜像）
./scripts/install-fast.sh --cn
```

安装成功后，访问 `http://<Your Public IP>` 开始使用。

**自定义安装（Helm）**

适用于需要自定义域名、HTTPS 或修改默认配置的场景。

```bash
# 1. 安装依赖（如果k3s/Helm已安装则跳过）
./scripts/install-run-environment.sh       # 国际镜像
# ./scripts/install-run-environment.sh --cn  # 中国镜像

# 2. 复制并修改配置
cp helm/values.yaml helm/values-custom.yaml
# 编辑helm/values-custom.yaml设置global.domain等参数

# 3. 安装平台
helm install mcpcan ./helm -f helm/values-custom.yaml \
  --namespace mcpcan --create-namespace --timeout 600s --wait
```

## 🚀 组件

MCPCan 由多个关键组件构成，这些组件共同构成 MCPCan 的功能框架，为用户提供全面的 MCP 服务管理功能。

| Project                                | Status                                                      | Description                                |
| -------------------------------------- | ----------------------------------------------------------- | ------------------------------------------ |
| [MCPCan-Web](frontend/)                | ![Status](https://img.shields.io/badge/status-active-green) | MCPCan Web UI (Vue.js Frontend)            |
| [MCPCan-Backend](backend/)             | ![Status](https://img.shields.io/badge/status-active-green) | MCPCan Backend Services (Go Microservices) |
| [MCPCan-Gateway](backend/cmd/gateway/) | ![Status](https://img.shields.io/badge/status-active-green) | MCP Gateway Service                        |
| [MCPCan-Market](backend/cmd/market/)   | ![Status](https://img.shields.io/badge/status-active-green) | MCP Service Marketplace                    |
| [MCPCan-Authz](backend/cmd/authz/)     | ![Status](https://img.shields.io/badge/status-active-green) | Authentication and Authorization Service   |

## 🐧 技术栈

### 🐧 前端

- **框架**：Vue.js 3.5+（组合式 API）
- **语言**：TypeScript
- **样式方案**：UnoCSS、SCSS
- **UI 组件库**：Element Plus
- **状态管理**：Pinia
- **构建工具**：Vite

### 🐧 后端

- **语言**：Go 1.24.2+
- **框架**：Gin、gRPC
- **数据库**：MySQL、Redis
- **容器化工具**：Docker、Kubernetes

## 🐧 第三方项目

- [mcpcan-deploy](https://github.com/Kymo-MCP/mcpcan-deploy) - MCPCan 的官方 Helm Charts 源代码库
- [MCPCan Helm Charts](https://kymo-mcp.github.io/mcpcan-deploy/) - MCPCan 的官方 Helm 图表库

## 💝 贡献指南

欢迎提交 PR 参与贡献！请参考[贡献](CONTRIBUTING.md)查看详细指引。

贡献前，请确保：

1. 阅读我们的[行为准则](CODE_OF_CONDUCT.md)
2. 检查现有 issue 和拉取请求（避免重复工作）
3. 遵循我们的编码规范和提交信息约定

## ✅ 安全

若发现安全漏洞，请参考我们的[安全政策](SECURITY.md)，按照负责任的披露准则进行报告。

## 📄 许可证

版权所有 (c) 2024-2025 MCPCan 团队，保留所有权利。

本软件基于 Apache 许可证第 2.0 版（以下简称“许可证”）授权；除非遵守许可证规定，否则不得使用本文件。您可通过以下链接获取许可证副本：

http://www.apache.org/licenses/LICENSE-2.0

除非适用法律要求或书面同意，否则根据许可证分发的软件均按“原样”提供，不附带任何明示或暗示的担保或条件。请查看许可证以了解具体的权限和限制条款。

## 👥 社区与支持

- 📖 [文档](https://kymo-mcp.github.io/mcpcan-deploy/)
- 💬 [Discord 社区](https://discord.com/channels/1428637640856571995/1428637896532820038)
- 🐛 [问题追踪](https://github.com/Kymo-MCP/mcpcan/issues)
- 📧 [邮件列表](mailto:opensource@kymo.cn)
- 🌐 微信<br>
  <img src="images/WeChat group QR code.jpg" alt="alt text" width="170">

## 💕 致谢

- 感谢[MCP 协议](https://modelcontextprotocol.io/)社区
- 感谢所有贡献者和支持者
- 特别致谢使 MCPCan 项目成为可能的开源项目

## 🌟 Star 历史

[![Star History Chart](https://api.star-history.com/svg?repos=Kymo-MCP/mcpcan&type=date&legend=top-left)](https://www.star-history.com/#Kymo-MCP/mcpcan&type=date&legend=top-left)

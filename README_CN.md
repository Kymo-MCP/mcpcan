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
  <img src="https://img.shields.io/badge/License-Sustainable%20Use-orange?style=for-the-badge" alt="Sustainable Use License"/>
</div>
<p align="center">
   <a href="./README.md">English</a> | <strong>中文版</strong> <br>
   <a href="https://demo.mcpcan.com">DemoSite : demo.mcpcan.com（登录: admin/admin123）</a> | <a href="https://www.mcpcan.com">MainSite : www.mcpcan.com</a><br>
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

有关详细部署说明，请参阅我们的[部署指南](https://www.mcpcan.com/docs/zh/guide/install)。

### 1. 获取部署仓库

```bash
# 推荐使用 GitHub
git clone https://github.com/Kymo-MCP/mcpcan.git
cd mcpcan/deploy/docker-compose/
```

### 2. 启动服务

**Docker Compose 快速启动（推荐）**

适用于本地开发、测试以及轻量级的生产环境部署。

```bash
# 1. 初始化配置文件
cp example.env .env
# (可选) 修改 .env 文件中的配置，如 REGISTRY_PREFIX 切换国内/国际镜像源

# 2. 生成最终配置
chmod +x replace.sh
./replace.sh

# 3. 启动服务
docker compose up -d

# 4. 访问 Web UI
登录: admin/admin123
```

安装成功后，访问 `http://localhost` (或 `http://<Your Public IP>`) 开始使用。

**Helm 安装**

适用于 Kubernetes 环境部署，请参考 [Helm 部署指南](https://kymo-mcp.github.io/deploy/)。

## 🚀 组件

MCPCan 采用微服务架构，由以下核心组件构成：

| 组件名称        | 目录路径                   | 功能描述                                 |
| --------------- | -------------------------- | ---------------------------------------- |
| **MCPCan-Web**  | `frontend/`                | 基于 Vue 3 的管理后端界面                |
| **MCP-Market**  | `backend/cmd/market/`      | 核心业务模块：应用市场、实例生命周期管理 |
| **MCP-Authz**   | `backend/cmd/authz/`       | 权限系统、RBAC、用户与部门管理           |

## 🐧 技术栈

### 🐧 前端

- **框架**：Vue 3.5.x (Composition API)
- **构建工具**：Vite 7.0
- **UI 组件库**：Element Plus
- **状态管理**：Pinia 3.0
- **样式方案**：UnoCSS, SCSS
- **编辑器**：Monaco Editor (代码编辑), Markdown-it

### 🐧 后端

- **开发语言**：Go 1.25.x
- **Web 框架**：Gin, gRPC (grpc-gateway)
- **数据存储**：MySQL (GORM), Redis (缓存/令牌)
- **核心库**：LangChainGo, MCP-Go, Docker SDK, Client-go
- **部署运维**：Traefik (边缘网关), Docker Compose, Helm (Kubernetes)

## 🐧 相关项目

- [MCPCan Helm Charts](https://kymo-mcp.github.io/deploy/) - MCPCan 的官方 Helm 图表索引库
- [部署文档](docs/zh/guide/install.md) - 详细的部署与配置指南

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

- 📖 [文档](https://kymo-mcp.github.io/deploy/)
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

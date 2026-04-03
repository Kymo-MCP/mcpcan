# 版本发布说明 (Release Notes)
**统计范围**: `fix/v2.1.3` ... `HEAD`

> 统计规则：已排除 `Merge branch...` 与 `version bump` 类提交。

## 🐛 问题修复 (Bug Fixes)
- **修复 stdio 代理协议错误映射导致的 MCP 调用失败**: AI 会话 MCP 配置生成改为优先使用 `proxyProtocol`，并基于 URL 后缀 (`/mcp`/`/sse`) 兜底判定协议；同时后端在 `type` 与 URL 冲突时自动纠正 transport，避免 `type=sse + /mcp` 触发 400。(`048c259`)
- **增强 MCP 网关 URL 归一化与工具调用兼容性**: 优化 MCP 网关 URL 归一化流程，补充相关测试与 LLM 工具调用兼容适配，降低容器/网关场景下初始化失败概率。(`f93c3ac`)
- **加固 OpenAPI 导入链路并统一部署镜像配置**: 修复 OpenAPI 导入与实例创建路径中的校验与容错问题，补充 Kubernetes label 与文件校验测试，并同步 Helm 与 Docker Compose 的镜像配置。(`5b9daf5`)
- **优化代理实例处理与默认部署行为**: 完善代理模式实例参数处理、初始化数据逻辑与 token 相关流程，补强前后端代理表单与 JSON 边界处理。(`1ea2690`)

## 🔧 其他 (Chore)
- 无。

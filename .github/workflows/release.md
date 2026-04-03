# Release Notes
**Scope**: `fix/v2.1.3` ... `HEAD`

> Rules: commits matching `Merge branch...` and version-bump-only changes are excluded.

## 🐛 Bug Fixes
- **Fixed MCP failures caused by incorrect stdio proxy transport mapping**: MCP config generation for AI sessions now prioritizes `proxyProtocol` and falls back to URL suffix detection (`/mcp` vs `/sse`). Backend now auto-corrects transport when `type` conflicts with the gateway URL, preventing `type=sse + /mcp` from returning HTTP 400. (`048c259`)
- **Improved MCP gateway URL normalization and tool-call compatibility**: Refined gateway URL normalization and added compatibility adjustments for LLM tool calling, reducing MCP initialization failures in containerized and gateway deployments. (`f93c3ac`)
- **Hardened OpenAPI import flow and aligned deployment image configuration**: Fixed validation and error-handling paths in OpenAPI import and instance creation, added Kubernetes label/file validation tests, and synchronized image settings across Helm and Docker Compose. (`5b9daf5`)
- **Improved proxy-instance handling and deployment defaults**: Enhanced proxy-mode instance parameter handling, initialization behavior, token-related flows, and frontend/backend boundary handling for proxy forms and JSON parsing. (`1ea2690`)

## 🔧 Chore
- None.

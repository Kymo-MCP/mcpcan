# MCPCan Docker Compose 部署指南

基于 Docker Compose 的一键式私有化部署方案。默认开启 HTTP 访问，支持通过扩展包启用 HTTPS。

## 1. 快速开始

只需三步即可完成基础部署：

### 第一步：初始化环境

```bash
cp .example.env .env
```

_(可选)_ 编辑 `.env` 修改 `VERSION` 或端口。

### 第二步：生成配置

```bash
# 赋予执行权限并运行
chmod +x replace.sh
./replace.sh
```

### 第三步：起动服务

```bash
docker compose up -d
```

---

## 2. 访问与验证

服务启动后（约 1 分钟后），可通过以下方式访问：

- **Web 控制台**: `http://localhost` (默认 80 端口)
- **Traefik 仪表盘**: `http://localhost:8090` (用于查看路由状态)
- **状态检查**: 执行 `docker compose ps` 确保所有容器显示为 `Up (healthy)`。

---

## 3. 进阶配置

### 启用 HTTPS (TLS)

系统支持“自带证书 (BYOC)”模式。如需开启 443 端口及加密访问，请参考专用指南：
👉 **[HTTPS 部署扩展方案](./HTTPS_SETUP.md)**

### 自定义环境变量

所有核心参数均在 `.env` 中定义，修改后请重新运行 `./replace.sh` 并执行 `docker compose up -d`。

| 变量              | 说明                         |
| :---------------- | :--------------------------- |
| `REGISTRY_PREFIX` | 镜像仓库前缀 (默认为 77kymo) |
| `VERSION`         | 系统版本标签                 |
| `ADMIN_USERNAME`  | 初始管理员账号               |
| `ADMIN_PASSWORD`  | 初始管理员密码               |
| `INIT_ADMIN_PASSWORD_POLICY` | 管理员密码初始化策略（`create_only`/`force_reset`） |
| `INIT_MENU_SYNC_POLICY` | 菜单同步策略（`off`/`add_only`/`full_sync`） |

---

## 4. 常用维护命令

| 任务             | 命令                                          |
| :--------------- | :-------------------------------------------- |
| **查看日志**     | `docker compose logs -f [服务名]`             |
| **重启服务**     | `docker compose restart`                      |
| **停止并清理**   | `docker compose down`                         |
| **更新镜像**     | `docker compose pull && docker compose up -d` |
| **查看资源占用** | `docker stats`                                |

---

## 5. 常见问题 (FAQ)

**Q: 修改了 .env 没生效？**
A: 请确保运行了 `./replace.sh` 来重新生成 `config/` 目录下的真实配置文件。

**Q: 数据库连接失败？**
A: 首次启动 MySQL 初始化较慢。请通过 `docker compose logs -f mysql` 查看是否已就绪。

**Q: 80 端口被占用？**
A: 修改 `.env` 中的 `MCP_ENTRY_SERVICE_PORT`，重新生成配置并重启服务。

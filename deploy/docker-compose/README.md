# MCPCan Docker Compose Deployment Guide

A one-click private deployment solution based on Docker Compose. HTTP access is enabled by default, with support for HTTPS via an extension package.

## 1. Quick Start

Complete your basic deployment in just three steps:

### Step 1: Initialize Environment

```bash
cp .example.env .env
```

_(Optional)_ Edit `.env` to modify `VERSION` or ports.

### Step 2: Generate Configuration

```bash
# Grant execution permissions and run
chmod +x replace.sh
./replace.sh
```

### Step 3: Start Services

```bash
docker compose up -d
```

---

## 2. Access and Verification

Once the services have started (typically after ~1 minute), you can access them via:

- **Web Console**: `http://localhost` (Default Port 80)
- **Traefik Dashboard**: `http://localhost:8090` (To view routing status)
- **Status Check**: Run `docker compose ps` to ensure all containers are `Up (healthy)`.

---

## 3. Advanced Configuration

### Enable HTTPS (TLS)

The system supports a "Bring Your Own Certificate (BYOC)" mode. To enable port 443 and secure access, please refer to the dedicated guide:
👉 **[HTTPS/TLS Deployment Guide](./HTTPS_SETUP.md)**

### Custom Environment Variables

All core parameters are defined in `.env`. After modifying, please rerun `./replace.sh` and execute `docker compose up -d`.

| Variable          | Description                                |
| :---------------- | :----------------------------------------- |
| `REGISTRY_PREFIX` | Image registry prefix (Defaults to 77kymo) |
| `VERSION`         | System version tag                         |
| `ADMIN_USERNAME`  | Initial admin username                     |
| `ADMIN_PASSWORD`  | Initial admin password                     |
| `INIT_ADMIN_PASSWORD_POLICY` | Admin password init policy (`create_only`/`force_reset`) |
| `INIT_MENU_SYNC_POLICY` | Menu sync policy (`off`/`add_only`/`full_sync`) |

---

## 4. Common Maintenance Commands

| Task                 | Command                                       |
| :------------------- | :-------------------------------------------- |
| **View Logs**        | `docker compose logs -f [service_name]`       |
| **Restart Services** | `docker compose restart`                      |
| **Stop and Clean**   | `docker compose down`                         |
| **Update Images**    | `docker compose pull && docker compose up -d` |
| **Resource Usage**   | `docker stats`                                |

---

## 5. FAQ

**Q: Changes in .env are not taking effect?**
A: Ensure you have run `./replace.sh` to regenerate the actual configuration files in the `config/` directory.

**Q: Database connection failed?**
A: Initializing MySQL for the first time is slow. Check `docker compose logs -f mysql` to see if it's ready.

**Q: Port 80 is occupied?**
A: Modify `MCP_ENTRY_SERVICE_PORT` in `.env`, regenerate the configuration, and restart the services.

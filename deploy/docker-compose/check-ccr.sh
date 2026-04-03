#!/usr/bin/env bash
set -euo pipefail

EXPECTED_REGISTRY="ccr.ccs.tencentyun.com/itqm-private"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

PASS=true

ok() { printf '[OK] %s\n' "$1"; }
warn() { printf '[WARN] %s\n' "$1"; }
fail() { printf '[FAIL] %s\n' "$1"; PASS=false; }

echo "=== MCPCAN CCR 一键检查 ==="
echo "目录: $SCRIPT_DIR"
echo "期望仓库: $EXPECTED_REGISTRY"
echo

# 1) 检查 .env
if [[ -f .env ]]; then
  actual_registry="$(grep -E '^REGISTRY_PREFIX=' .env | tail -n1 | cut -d'=' -f2- || true)"
  if [[ "${actual_registry}" == "${EXPECTED_REGISTRY}" ]]; then
    ok ".env 中 REGISTRY_PREFIX 已是 CCR"
  else
    fail ".env 中 REGISTRY_PREFIX=${actual_registry}（应为 ${EXPECTED_REGISTRY}）"
  fi
else
  fail "未找到 .env 文件"
fi

echo
# 2) 检查核心容器镜像
if docker ps --format '{{.Names}}' | rg -q '^mcp-market-dev$'; then
  while IFS= read -r line; do
    name="$(awk '{print $1}' <<<"$line")"
    image="$(awk '{print $2}' <<<"$line")"
    case "${name}" in
      mcp-entry-dev)
        if [[ "${image}" == traefik:* ]]; then
          ok "${name} 使用 ${image}（预期）"
        else
          warn "${name} 使用 ${image}（通常应为 traefik:*）"
        fi
        ;;
      mcp-market-dev|mcp-authz-dev|mcp-web-dev)
        if [[ "${image}" == "${EXPECTED_REGISTRY}"/* ]]; then
          ok "${name} 使用 CCR 镜像：${image}"
        else
          fail "${name} 未使用 CCR 镜像：${image}"
        fi
        ;;
    esac
  done < <(docker ps --format '{{.Names}} {{.Image}}' | rg '^(mcp-entry-dev|mcp-market-dev|mcp-authz-dev|mcp-web-dev) ')
else
  warn "mcp-market-dev 未运行，跳过核心容器检查"
fi

echo
# 3) 检查 mcp-market 运行时镜像前缀环境变量
if docker ps --format '{{.Names}}' | rg -q '^mcp-market-dev$'; then
  market_mirror="$(docker exec mcp-market-dev sh -lc 'printenv REGISTORY_IMAGE_MIRROR || true' | tr -d '\r')"
  if [[ "${market_mirror}" == "${EXPECTED_REGISTRY}" ]]; then
    ok "mcp-market-dev 内 REGISTORY_IMAGE_MIRROR=${market_mirror}"
  elif [[ -z "${market_mirror}" ]]; then
    fail "mcp-market-dev 内未设置 REGISTORY_IMAGE_MIRROR（会回落到 77kymo）"
  else
    fail "mcp-market-dev 内 REGISTORY_IMAGE_MIRROR=${market_mirror}（应为 ${EXPECTED_REGISTRY}）"
  fi
else
  warn "mcp-market-dev 未运行，跳过环境变量检查"
fi

echo
# 4) 检查动态实例容器镜像来源（最关键）
instance_lines="$(docker ps --format '{{.Names}} {{.Image}}' | rg '^(mcp-instance-|mcp-ext-)' || true)"
if [[ -z "${instance_lines}" ]]; then
  warn "当前没有 mcp-instance/mcp-ext 动态容器"
else
  bad_count=0
  total_count=0
  while IFS= read -r line; do
    [[ -z "${line}" ]] && continue
    total_count=$((total_count + 1))
    name="$(awk '{print $1}' <<<"$line")"
    image="$(awk '{print $2}' <<<"$line")"
    if [[ "${image}" == "${EXPECTED_REGISTRY}"/* ]]; then
      ok "${name} 使用 CCR 镜像：${image}"
    else
      bad_count=$((bad_count + 1))
      fail "${name} 混用非 CCR 镜像：${image}"
    fi
  done <<<"${instance_lines}"

  if [[ $bad_count -eq 0 ]]; then
    ok "动态实例容器检查通过（${total_count}/${total_count}）"
  else
    fail "动态实例容器存在混用（${bad_count}/${total_count}）"
  fi
fi

echo
# 5) 给出结论
if [[ "$PASS" == true ]]; then
  echo "✅ 结论：当前 Docker dev 环境已统一为 CCR。"
  exit 0
else
  echo "❌ 结论：存在非 CCR 配置或容器。"
  echo
  echo "可用修复命令（在本目录执行）："
  cat <<'FIX'
perl -i -pe 's#^REGISTRY_PREFIX=.*#REGISTRY_PREFIX=ccr.ccs.tencentyun.com/itqm-private#' .env
./replace.sh
cat > docker-compose.dev.ccr.override.yml <<'YAML'
services:
  mcp-market-svc:
    environment:
      - REGISTORY_IMAGE_MIRROR=ccr.ccs.tencentyun.com/itqm-private
YAML
docker compose -f docker-compose.dev.yml -f docker-compose.dev.ccr.override.yml down
docker ps -aq --format '{{.ID}} {{.Names}}' | awk '/mcp-instance-|mcp-ext-/{print $1}' | xargs -r docker rm -f
docker compose -f docker-compose.dev.yml -f docker-compose.dev.ccr.override.yml pull
docker compose -f docker-compose.dev.yml -f docker-compose.dev.ccr.override.yml up -d --build --force-recreate
FIX
  exit 1
fi

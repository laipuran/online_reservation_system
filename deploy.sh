#!/usr/bin/env bash
set -euo pipefail

PI_USER="puranlai"
PI_HOST="192.168.31.120"
PI_DIR="~/Code/online_reservation_system"
SSH_TARGET="${PI_USER}@${PI_HOST}"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
BE_TAR="${SCRIPT_DIR}/ors-be.tar"
FE_TAR="${SCRIPT_DIR}/ors-fe.tar"

GREEN='\033[0;32m'; YELLOW='\033[1;33m'; RED='\033[0;31m'; NC='\033[0m'
info()  { echo -e "${GREEN}[INFO]${NC} $*"; }
warn()  { echo -e "${YELLOW}[WARN]${NC} $*"; }
error() { echo -e "${RED}[ERROR]${NC} $*"; }

usage() {
  cat <<EOF
用法: $(basename "$0") [选项]

一键部署在线预约系统到树莓派。

选项:
  --skip-build    跳过镜像构建，仅传输和部署
  --help          显示此帮助信息

环境变量:
  JWT_SECRET       JWT 签名密钥 (默认: change-me-in-production)
  ALLOWED_ORIGINS  允许的跨域来源 (默认: *)
EOF
  exit 0
}

SKIP_BUILD=false
while [[ $# -gt 0 ]]; do
  case "$1" in
    --skip-build) SKIP_BUILD=true; shift ;;
    --help) usage ;;
    *) error "未知选项: $1"; usage ;;
  esac
done

cleanup() { rm -f "$BE_TAR" "$FE_TAR"; }
trap cleanup EXIT

if ! docker buildx version &>/dev/null; then
  error "Docker buildx 不可用，请安装 Docker Desktop 或 docker-buildx 插件"
  exit 1
fi

if ! ssh -o ConnectTimeout=5 -o BatchMode=yes "$SSH_TARGET" exit 2>/dev/null; then
  warn "SSH 免密登录未配置，尝试使用密码登录..."
fi

info "============================================"
info "  在线预约系统 - 树莓派部署脚本"
info "  目标: ${SSH_TARGET}:${PI_DIR}"
info "============================================"

if [ "$SKIP_BUILD" = false ]; then
  info "Step 1/5: 交叉编译后端镜像 (linux/arm64) ..."
  docker buildx build \
    --platform linux/arm64 \
    -t ors-be:latest \
    -o type=docker,dest="$BE_TAR" \
    -f "${SCRIPT_DIR}/Dockerfile.be" \
    "$SCRIPT_DIR"
  info "  ✓ 后端镜像已导出 -> ors-be.tar"

  info "Step 2/5: 交叉编译前端镜像 (linux/arm64) ..."
  docker buildx build \
    --platform linux/arm64 \
    -t ors-fe:latest \
    -o type=docker,dest="$FE_TAR" \
    -f "${SCRIPT_DIR}/Dockerfile.fe" \
    "$SCRIPT_DIR"
  info "  ✓ 前端镜像已导出 -> ors-fe.tar"
else
  info "Step 1-2/5: 跳过镜像构建 (--skip-build) ..."
fi

info "Step 3/5: 创建远程目录并传输文件 ..."
ssh "$SSH_TARGET" "mkdir -p ${PI_DIR}/deploy/nginx ${PI_DIR}/ors-be/migrations ${PI_DIR}/ors-be/db"

scp "$BE_TAR"         "$SSH_TARGET":"${PI_DIR}/ors-be.tar"
scp "$FE_TAR"         "$SSH_TARGET":"${PI_DIR}/ors-fe.tar"
scp "${SCRIPT_DIR}/docker-compose.yml"               "$SSH_TARGET":"${PI_DIR}/docker-compose.yml"
scp "${SCRIPT_DIR}/deploy/nginx/default.conf"        "$SSH_TARGET":"${PI_DIR}/deploy/nginx/default.conf"
scp -r "${SCRIPT_DIR}/ors-be/migrations/"            "$SSH_TARGET":"${PI_DIR}/ors-be/migrations/"
scp "${SCRIPT_DIR}/ors-be/db/seed.sql"               "$SSH_TARGET":"${PI_DIR}/ors-be/db/seed.sql"

scp "${SCRIPT_DIR}/.env"               "$SSH_TARGET":"${PI_DIR}/.env" 2>/dev/null || true
info "  ✓ 文件传输完成"

info "Step 4/5: 在树莓派上导入镜像 ..."
ssh "$SSH_TARGET" bash -s << 'REMOTE'
  set -euo pipefail
  cd ~/Code/online_reservation_system

  if [ -f ors-be.tar ]; then
    echo "  → 导入后端镜像..."
    docker load -i ors-be.tar
    rm -f ors-be.tar
  fi

  if [ -f ors-fe.tar ]; then
    echo "  → 导入前端镜像..."
    docker load -i ors-fe.tar
    rm -f ors-fe.tar
  fi
REMOTE
info "  ✓ 镜像导入完成"

info "Step 5/5: 启动服务 ..."
ssh "$SSH_TARGET" bash -s << 'REMOTE'
  set -euo pipefail
  cd ~/Code/online_reservation_system
  docker compose up -d
  echo "  → 服务状态:"
  docker compose ps
REMOTE
info "  ✓ 服务已启动"

info "============================================"
info "  部署完成！"
info "  访问 http://${PI_HOST} 即可使用"
info "============================================"

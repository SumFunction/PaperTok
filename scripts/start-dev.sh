#!/bin/bash
# PaperTok 本地开发一键启动：同时启动后端 (Go) 和前端 (Vite)
# 用法: ./scripts/start-dev.sh  或  bash scripts/start-dev.sh

set -e
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

BACKEND_PID=""
FRONTEND_PID=""

# 设置开发环境的 JWT_SECRET（仅用于本地开发，生产环境请使用安全的随机字符串）
export JWT_SECRET="${JWT_SECRET:-papertok_dev_secret_key_2024_do_not_use_in_production}"

cleanup() {
  echo ""
  echo "正在停止服务..."
  [ -n "$BACKEND_PID" ]  && kill "$BACKEND_PID" 2>/dev/null || true
  [ -n "$FRONTEND_PID" ] && kill "$FRONTEND_PID" 2>/dev/null || true
  exit 0
}

trap cleanup SIGINT SIGTERM EXIT

echo "=========================================="
echo "  PaperTok 本地开发环境"
echo "=========================================="
echo "  后端: http://localhost:8080"
echo "  前端: http://localhost:5173"
echo "  按 Ctrl+C 停止全部"
echo "=========================================="
echo ""

# 启动后端（需在 backend 目录下以找到 config.yaml）
(cd "$ROOT/backend" && go run cmd/server/main.go) &
BACKEND_PID=$!

# 稍等后端监听端口
sleep 2

# 启动前端
(cd "$ROOT/frontend" && npm run dev) &
FRONTEND_PID=$!

wait

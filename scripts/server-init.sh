#!/bin/bash
# PaperTok 服务器初始化脚本
# 在腾讯云 CentOS 服务器上运行此脚本

set -e

echo "=========================================="
echo "  PaperTok 服务器初始化脚本"
echo "=========================================="

# 检查是否以 root 运行
if [ "$EUID" -ne 0 ]; then
  echo "请使用 root 权限运行此脚本"
  exit 1
fi

# 更新系统
echo "[1/6] 更新系统..."
yum update -y

# 安装必要工具
echo "[2/6] 安装必要工具..."
yum install -y yum-utils device-mapper-persistent-data lvm2 git curl wget

# 安装 Docker
echo "[3/6] 安装 Docker..."
if ! command -v docker &> /dev/null; then
  # 添加 Docker 仓库（使用阿里云镜像）
  yum-config-manager --add-repo https://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo
  
  # 安装 Docker
  yum install -y docker-ce docker-ce-cli containerd.io
  
  # 启动 Docker
  systemctl start docker
  systemctl enable docker
  
  # 配置 Docker 镜像加速（阿里云）
  mkdir -p /etc/docker
  cat > /etc/docker/daemon.json << 'EOF'
{
  "registry-mirrors": [
    "https://mirror.ccs.tencentyun.com",
    "https://registry.docker-cn.com"
  ]
}
EOF
  systemctl daemon-reload
  systemctl restart docker
  
  echo "Docker 安装完成！"
else
  echo "Docker 已安装，跳过"
fi

# 安装 Docker Compose
echo "[4/6] 安装 Docker Compose..."
if ! command -v docker-compose &> /dev/null; then
  # 使用国内镜像下载
  curl -L "https://ghproxy.com/https://github.com/docker/compose/releases/download/v2.24.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
  chmod +x /usr/local/bin/docker-compose
  ln -sf /usr/local/bin/docker-compose /usr/bin/docker-compose
  echo "Docker Compose 安装完成！"
else
  echo "Docker Compose 已安装，跳过"
fi

# 创建项目目录并克隆代码
echo "[5/6] 克隆项目代码..."
PROJECT_DIR="/opt/papertok"
if [ ! -d "$PROJECT_DIR" ]; then
  git clone https://github.com/SumFunction/PaperTok.git $PROJECT_DIR
  echo "代码克隆完成！"
else
  echo "项目目录已存在，更新代码..."
  cd $PROJECT_DIR
  git pull origin main
fi

# 配置防火墙
echo "[6/6] 配置防火墙..."
if command -v firewall-cmd &> /dev/null; then
  firewall-cmd --permanent --add-port=80/tcp
  firewall-cmd --permanent --add-port=443/tcp
  firewall-cmd --permanent --add-port=8080/tcp
  firewall-cmd --reload
  echo "防火墙配置完成！"
else
  echo "未检测到 firewalld，请手动配置安全组开放 80、443、8080 端口"
fi

echo ""
echo "=========================================="
echo "  初始化完成！"
echo "=========================================="
echo ""
echo "下一步操作："
echo "1. 进入项目目录：cd /opt/papertok"
echo "2. 启动服务：docker-compose up -d --build"
echo "3. 查看日志：docker-compose logs -f"
echo ""
echo "访问地址：http://$(hostname -I | awk '{print $1}')"
echo ""

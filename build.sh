#!/bin/bash
set -e

VERSION=${1:-"v1.0.0"}
DEV_MODE=${2:-"false"}

echo "🚀 开始构建 Palworld Server Tool..."

# 检查必需的工具
for tool in node pnpm go curl; do
    if ! command -v $tool &> /dev/null; then
        echo "❌ 未找到 $tool，请先安装"
        exit 1
    fi
done

# 构建前端
echo "📦 构建前端项目..."

# 构建 web
cd web
pnpm install
pnpm build
cd ..

# 构建 pal-conf
cd pal-conf
pnpm install
pnpm build
cd ..

# 准备资源文件
echo "📁 准备资源文件..."
mkdir -p assets
cp -r pal-conf/dist/assets/* assets/
cp pal-conf/dist/index.html pal-conf.html

# 下载 sav_cli
echo "⬇️ 下载 sav_cli..."
ARCH=$(uname -m)
OS=$(uname -s | tr '[:upper:]' '[:lower:]')

if [ "$OS" = "darwin" ]; then
    if [ "$ARCH" = "arm64" ]; then
        SAV_CLI_ARCH="aarch64"
    else
        SAV_CLI_ARCH="x86_64"
    fi
    SAV_CLI_URL="https://github.com/zaigie/palworld-server-tool/releases/download/v0.9.9/sav_cli_darwin_$SAV_CLI_ARCH"
elif [ "$OS" = "linux" ]; then
    if [ "$ARCH" = "aarch64" ] || [ "$ARCH" = "arm64" ]; then
        SAV_CLI_ARCH="aarch64"
    else
        SAV_CLI_ARCH="x86_64"
    fi
    SAV_CLI_URL="https://github.com/zaigie/palworld-server-tool/releases/download/v0.9.9/sav_cli_linux_$SAV_CLI_ARCH"
else
    echo "❌ 不支持的操作系统: $OS"
    exit 1
fi

mkdir -p temp
curl -L -o temp/sav_cli "$SAV_CLI_URL"
chmod +x temp/sav_cli
mv temp/sav_cli ./sav_cli

# 下载地图文件
echo "🗺️ 下载地图文件..."
if [ ! -d "map" ]; then
    curl -L -o map.zip "https://github.com/zaigie/palworld-server-tool/releases/download/v0.9.9/map.zip"
    unzip -o map.zip
    rm map.zip
fi

# 编译后端
echo "🔨 编译 Go 后端..."
if [ "$DEV_MODE" = "true" ]; then
    go build -o pst main.go
else
    go build -ldflags="-s -w -X 'main.version=$VERSION'" -o pst main.go
fi

# 设置权限
chmod +x pst

# 设置环境变量
export SAVE__DECODE_PATH="./sav_cli"

# 清理临时文件
rm -rf temp

echo "✅ 构建完成！"
echo "运行命令: ./pst"
echo "访问地址: http://localhost:8080"

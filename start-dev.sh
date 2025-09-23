#!/bin/bash

# Seven AI 开发环境启动脚本

echo "🚀 启动 Seven AI 开发环境..."

# 检查是否安装了必要的工具
check_dependencies() {
    echo "📋 检查依赖..."
    
    if ! command -v node &> /dev/null; then
        echo "❌ Node.js 未安装，请先安装 Node.js"
        exit 1
    fi
    
    if ! command -v go &> /dev/null; then
        echo "❌ Go 未安装，请先安装 Go"
        exit 1
    fi
    
    if ! command -v mysql &> /dev/null; then
        echo "❌ MySQL 未安装，请先安装 MySQL"
        exit 1
    fi
    
    echo "✅ 依赖检查完成"
}

# 启动前端
start_frontend() {
    echo "🎨 启动前端服务..."
    cd frontend
    npm install
    npm run dev &
    cd ..
    echo "✅ 前端服务启动完成 (http://localhost:3000)"
}

# 启动后端
start_backend() {
    echo "⚙️  启动后端服务..."
    cd backend
    go mod tidy
    go run main.go &
    cd ..
    echo "✅ 后端服务启动完成 (http://localhost:8080)"
}

# 主函数
main() {
    check_dependencies
    start_frontend
    sleep 3
    start_backend
    
    echo ""
    echo "🎉 Seven AI 开发环境启动完成！"
    echo "📱 前端地址: http://localhost:3000"
    echo "🔧 后端地址: http://localhost:8080"
    echo "📊 健康检查: http://localhost:8080/health"
    echo ""
    echo "按 Ctrl+C 停止所有服务"
    
    # 等待用户中断
    wait
}

# 捕获中断信号
trap 'echo ""; echo "🛑 正在停止服务..."; kill $(jobs -p) 2>/dev/null; exit 0' INT

main

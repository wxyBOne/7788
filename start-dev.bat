@echo off
chcp 65001 >nul

echo 🚀 启动 Seven AI 开发环境...

:: 检查依赖
echo 📋 检查依赖...
where node >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ Node.js 未安装，请先安装 Node.js
    pause
    exit /b 1
)

where go >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ Go 未安装，请先安装 Go
    pause
    exit /b 1
)

echo ✅ 依赖检查完成

:: 启动前端
echo 🎨 启动前端服务...
cd frontend
call npm install
start "前端服务" cmd /k "npm run dev"
cd ..

:: 等待前端启动
timeout /t 5 /nobreak >nul

:: 启动后端
echo ⚙️  启动后端服务...
cd backend
go mod tidy
start "后端服务" cmd /k "go run main.go"
cd ..

echo.
echo 🎉 Seven AI 开发环境启动完成！
echo 📱 前端地址: http://localhost:3000
echo 🔧 后端地址: http://localhost:8080
echo 📊 健康检查: http://localhost:8080/health
echo.
echo 按任意键退出...
pause >nul

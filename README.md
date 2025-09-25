# Seven AI 角色扮演与养成系统

## 环境配置

### 1. 环境变量配置

项目使用环境变量来管理敏感配置信息。请创建 `.env` 文件并配置以下变量：

```bash
# 服务器配置
PORT=8080
ENVIRONMENT=development

# 数据库配置
DATABASE_URL=root:your_password@tcp(localhost:3306)/seven_ai?charset=utf8mb4&parseTime=True&loc=Local

# AI服务配置
AI_API_KEY=your_ai_api_key_here
AI_BASE_URL=https://api.openai.com/v1

# 语音服务配置
ASR_API_KEY=your_asr_api_key_here
TTS_API_KEY=your_tts_api_key_here

# 视觉服务配置
VISION_API_KEY=your_vision_api_key_here

# JWT密钥
JWT_SECRET=your_jwt_secret_key_here
```

### 2. 安全注意事项

- **不要**将 `.env` 文件提交到版本控制系统
- **不要**在代码中硬编码敏感信息
- 使用强密码和复杂的JWT密钥
- 定期轮换API密钥

### 3. 数据库设置

1. 创建MySQL数据库：
```sql
CREATE DATABASE seven_ai CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

2. 运行数据库迁移：
```bash
mysql -u root -p seven_ai < database/schema.sql
```

### 4. 项目结构

```
├── frontend/          # Vue.js 前端
├── backend/           # Go 后端
├── database/          # 数据库脚本
├── .env              # 环境变量（不提交到git）
├── .gitignore        # Git忽略文件
└── README.md         # 项目说明
```

## 开发指南

### 前端开发
```bash
cd frontend
npm install
npm run dev
```

### 后端开发
```bash
cd backend
go mod tidy
go run main.go
```

## 部署说明

1. 确保所有环境变量都已正确配置
2. 数据库已创建并运行迁移
3. 构建前端：`npm run build`
4. 运行后端服务

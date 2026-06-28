# 项目搭建日志

## 环境检查

### 1. Go 环境

```bash
go version
go version go1.26.4 windows/amd64
```

### 2. Docker 环境

```bash
docker --version
Docker version 27.0.3, build 7d4bcd8

docker-compose --version
Docker Compose version v2.27.0
```

### 3. Node.js 环境

```bash
node --version
v20.x.x

pnpm --version
8.x.x
```

## 项目启动步骤

### 1. 克隆项目

```bash
git clone https://github.com/multica-ai/multica.git
cd multica
```

### 2. 配置环境变量

创建 `.env` 文件或使用 `docker-compose.selfhost.yml`。

### 3. Docker Compose 启动

```bash
docker-compose -f docker-compose.selfhost.yml up
```

### 4. 服务验证

- 前端：http://localhost:8081
- API：http://localhost:8080

## 遇到的问题及解决

### 问题 1：端口冲突

**现象**：端口 3000 被 Windows 系统占用

**解决方案**：修改 `docker-compose.selfhost.yml`，将前端端口改为 8081

```yaml
FRONTEND_ORIGIN: http://localhost:8081
MULTICA_APP_URL: http://localhost:8081
```

### 问题 2：邮件验证码无法接收

**现象**：登录时未收到邮件验证码

**解决方案**：配置开发环境固定验证码

```yaml
MULTICA_DEV_VERIFICATION_CODE: 123456
```

### 问题 3：Docker Daemon 未启动

**现象**：Docker 命令无法执行

**解决方案**：启动 Docker Desktop

## 运行状态

- PostgreSQL：正常运行（端口 5432）
- Redis：正常运行（端口 6379）
- 后端 API：正常运行（端口 8080）
- 前端 Web：正常运行（端口 8081）

## 测试登录

使用固定验证码 `123456` 完成登录验证。
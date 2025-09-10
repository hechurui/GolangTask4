## 项目结构

```
GolangTask4/  
├── internal/            # 私有应用代码 
│   ├── api/             # API处理器 
│   │   ├── handlers/    # 请求处理函数 
│   │   ├── middlewares/ # 中间件 
│   │   └── routes/      # 路由定义 
│   ├── domain/          # 领域模型 
├── pkg/                 # 可公共使用的包 
│   ├── auth/            # 认证相关工具 
│   ├── database/        # 数据库连接工具 
├── .env                 # 环境变量配置 
├── go.mod               # 依赖管理 
├── main.go              # 程序主入口 
└── go.sum  
```

## 依赖安装

```bash
# 初始化项目 
go mod init GolangTask4  
 
# 核心依赖安装 
go get -u github.com/gin-gonic/gin              # Gin框架 
go get -u gorm.io/gorm                          # GORM ORM 
go get -u gorm.io/driver/mysql                  # MySQL驱动 
go get -u github.com/golang-jwt/jwt/v5          # JWT认证 
go get -u golang.org/x/crypto                   # 密码加密 
go get -u github.com/joho/godotenv              # 环境变量管理 

# 安装依赖

go mod tidy
```

### 运行环境

- Go 1.24 或更高版本
- MySQL 5.7+ 或 MySQL 8.0+

### 启动项目

```bash
go run main.go
```



## API 接口文档
系统提供的 API 接口如下：

# 用户认证
POST /api/register - 用户注册
POST /api/login - 用户登录
# 文章管理
GET /api/posts - 获取所有文章
GET /api/posts/:id - 获取单篇文章
POST /api/posts - 创建文章（需要认证）
PUT /api/posts/:id - 更新文章（需要认证，仅作者）
DELETE /api/posts/:id - 删除文章（需要认证，仅作者）
# 评论功能
GET /api/posts/:id/comments - 获取文章的所有评论
POST /api/posts/:id/comments - 创建评论（需要认证）


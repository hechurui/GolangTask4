package handlers

import (
	"GolangTask4/internal/domain"
	"GolangTask4/pkg/auth"
	"GolangTask4/pkg/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 用户请求结构体
type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// LoginResponse 登录响应结构体
type LoginResponse struct {
	Token string      `json:"token"`
	User  domain.User `json:"user"`
}

// 注册
func Register(c *gin.Context) {
	var user UserRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// 先检查DB是否初始化
	if database.DB == nil {
		c.JSON(500, gin.H{"error": "database not initialized"})
		return
	}
	// 检查用户名是否已存在
	var existingUser domain.User
	if err := database.DB.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	// 检查邮箱是否已存在
	if err := database.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	// 密码加密
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": "could not hash password"})
		return
	}

	// 创建新用户
	saveUser := domain.User{
		Username: user.Username,
		Password: hashedPassword,
		Email:    user.Email,
	}

	// 保存用户
	if err := database.DB.Create(&saveUser).Error; err != nil {
		c.JSON(500, gin.H{"error": "could not create user"})
		return
	}

	// 生成JWT token
	tokenString, err := auth.GenerateToken(saveUser.ID, saveUser.Username)
	if err != nil {
		c.JSON(500, gin.H{"error": "could not generate token"})
		return
	}

	c.JSON(200, LoginResponse{
		Token: tokenString,
		User: domain.User{
			ID:       saveUser.ID,
			Username: saveUser.Username,
			Email:    saveUser.Email,
		},
	})
}

// 登录
func Login(c *gin.Context) {
	var credentials struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 1、通过username查询出用户信息
	var user domain.User
	if err := database.DB.Where("username = ?", credentials.Username).First(&user).Error; err != nil {
		c.JSON(401, gin.H{"error": "invalid credentials username"})
		return
	}

	// 校验密码
	if checkResult := auth.CheckPassword(credentials.Password, user.Password); !checkResult {
		c.JSON(401, gin.H{"error": "invalid credentials"})
		return
	}

	// 生成JWT token
	tokenString, err := auth.GenerateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(500, gin.H{"error": "could not generate token"})
		return
	}

	c.JSON(200, LoginResponse{
		Token: tokenString,
		User: domain.User{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		},
	})
}

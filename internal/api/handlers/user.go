package handlers

import (
	"GolangTask4/internal/domain"
	"GolangTask4/pkg/auth"
	"GolangTask4/pkg/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

// LoginResponse 登录响应结构体
type LoginResponse struct {
	Token string      `json:"token"`
	User  domain.User `json:"user"`
}

// 注册
func register(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
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
	user.Password = string(hashedPassword)

	// 保存用户
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "could not create user"})
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

// 登录
func login(c *gin.Context) {
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
		c.JSON(401, gin.H{"error": "invalid credentials"})
		return
	}

	// 校验密码
	if checkResult := auth.CheckPassword(user.Password, credentials.Password); !checkResult {
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

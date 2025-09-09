package auth

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

// Claims JWT载荷结构
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT令牌
func GenerateToken(userID uint, username string) (string, error) {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		return "", err
	}

	// 获取JWT密钥和过期时间
	secret := os.Getenv("JWT_SECRET")
	expHours, _ := strconv.Atoi(os.Getenv("JWT_EXPIRATION_HOURS"))
	if expHours == 0 {
		expHours = 72 // 默认72小时过期
	}

	// 设置过期时间
	expirationTime := time.Now().Add(time.Duration(expHours) * time.Hour)

	// 创建claims
	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "blog-backend",
		},
	}

	// 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名token
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseToken 解析JWT令牌
func ParseToken(tokenString string) (*Claims, error) {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, errors.New("JWT密钥未配置")
	}

	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("不支持的签名方法")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	// 验证claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("无效的token")
}

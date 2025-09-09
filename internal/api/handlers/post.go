package handlers

import (
	"GolangTask4/internal/domain"
	"GolangTask4/pkg/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 创建文章
func CreatePost(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(403, gin.H{"error": "无权限"})
		return
	}

	var post struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(400, gin.H{"error": "无效的请求数据"})
		return
	}

	newPost := domain.Post{
		Title:   post.Title,
		Content: post.Content,
		UserID:  userID.(uint),
	}

	if result := database.DB.Create(&newPost); result.Error != nil {
		c.JSON(500, gin.H{"error": "创建文章失败"})
		return
	}

	c.JSON(200, gin.H{"文章创建成功->id": newPost.ID})
}

// 查询所有文章
func GetPosts(c *gin.Context) {
	var posts []domain.Post
	if err := database.DB.Preload("User").Find(&posts).Error; err != nil {
		c.JSON(404, gin.H{"error": "Failed to fetch posts"})
		return
	}

	c.JSON(200, posts)
}

// 根据id，查询单篇文章
func GetPost(c *gin.Context) {
	var post domain.Post
	id := c.Param("id")

	if err := database.DB.Preload("User").Preload("Comments").Preload("Comments.User").First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	c.JSON(200, post)
}

// 更新文章
func UpdatePost(c *gin.Context) {
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(403, gin.H{"error": "无权限修改此文章"})
		return
	}

	// 文章ID
	id := c.Param("id")

	var post domain.Post
	if err := database.DB.First(&post, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "文章未找到"})
		return
	}

	var updateData struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(400, gin.H{"error": "无效的请求数据"})
		return
	}

	database.DB.Model(&post).Updates(domain.Post{
		Title:   updateData.Title,
		Content: updateData.Content,
	})

	c.JSON(200, gin.H{"message": "文章更新成功"})
}

// 删除文章
func DeletePost(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(403, gin.H{"error": "无权限"})
		return
	}

	id := c.Param("id")
	var post domain.Post
	if err := database.DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	// 检查是否为文章作者
	if post.UserID != userID.(uint) {
		c.JSON(403, gin.H{"error": "参数异常，删除无权限"})
		return
	}

	// 删除文章
	if err := database.DB.Delete(&post).Error; err != nil {
		c.JSON(500, gin.H{"error": "文章删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "文章删除成功"})
}

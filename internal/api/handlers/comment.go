package handlers

import (
	"GolangTask4/internal/domain"
	"GolangTask4/pkg/database"

	"github.com/gin-gonic/gin"
)

// 创建评论
func CreateComment(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(403, gin.H{"error": "无权限"})
		return
	}
	postID := c.Param("postId")

	var comment struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(400, gin.H{"error": "无效的请求数据"})
		return
	}

	// 校验文章是否存在
	var post domain.Post
	if err := database.DB.First(&post, postID).Error; err != nil {
		c.JSON(404, gin.H{"error": "文章不存在"})
		return
	}

	newComment := domain.Comment{
		Content: comment.Content,
		UserID:  userID.(uint),
		PostID:  post.ID,
	}

	if result := database.DB.Create(&newComment); result.Error != nil {
		c.JSON(500, gin.H{"error": "创建评论失败"})
		return
	}

	c.JSON(200, gin.H{"评论创建成功->id": newComment.ID})
}

// 获取文章的评论
func GetPostComments(c *gin.Context) {
	postID := c.Param("postId")
	var comments []domain.Comment

	if err := database.DB.Where("post_id = ?", postID).Preload("User").Find(&comments).Error; err != nil {
		c.JSON(500, gin.H{"error": "获取评论失败"})
		return
	}

	c.JSON(200, comments)
}

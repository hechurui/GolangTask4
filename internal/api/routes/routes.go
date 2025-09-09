package routes

import (
	"GolangTask4/internal/api/handlers"
	"GolangTask4/internal/api/middlewares"

	"github.com/gin-gonic/gin"
)

// 设置路由
func SetupRoutes(router *gin.Engine) {
	// 公开路由
	public := router.Group("/api")
	{

		// 认证路由
		public.POST("/register", handlers.Register)
		public.POST("/login", handlers.Login)

		// 文章路由（公开访问）
		public.GET("/posts", handlers.GetPosts)
		public.GET("/posts/:id", handlers.GetPost)

		// 评论路由（公开访问）
		public.GET("/posts/:postId/comments", handlers.GetPostComments)
	}

	// 需要认证的路由
	protected := router.Group("/api")
	protected.Use(middlewares.AuthMiddleware())
	{
		// 文章路由（需要认证）
		protected.POST("/posts", handlers.CreatePost)
		protected.PUT("/posts/:id", handlers.UpdatePost)
		protected.DELETE("/posts/:id", handlers.DeletePost)

		// 评论路由（需要认证）
		protected.POST("/posts/:postId/comments", handlers.CreateComment)
	}
}

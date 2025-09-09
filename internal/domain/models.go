package domain

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Username  string         `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Email     string         `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"size:100;not null" json:"-"` // 不在JSON中返回密码
	Posts     []Post         `gorm:"foreignKey:UserID" json:"posts,omitempty"`
	Comments  []Comment      `gorm:"foreignKey:UserID" json:"comments,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Post 文章模型
type Post struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Title     string         `gorm:"size:200;not null" json:"title"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	UserID    uint           `gorm:"not null" json:"user_id"`
	User      User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Comments  []Comment      `gorm:"foreignKey:PostID" json:"comments,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Comment 评论模型
type Comment struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	UserID    uint           `gorm:"not null" json:"user_id"`
	PostID    uint           `gorm:"not null" json:"post_id"`
	User      User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Post      Post           `gorm:"foreignKey:PostID" json:"post,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

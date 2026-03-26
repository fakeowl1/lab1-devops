package model

import "time"

type Note struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Title      string    `gorm:"type:varchar(128);not null" json:"title"`
	Content    string    `gorm:"not null" json:"content"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"createdAt"`
}

type CreateNote struct {
	Title string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

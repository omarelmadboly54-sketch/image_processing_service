package models

import "time"

type Image struct {
	ID           int64  	`gorm:"primarykey" json:"-"`
	FileUrl      string 	`gorm:"not null" json:"file_url"`
	OriginalName string 	`gorm:"not null" json:"original_name"`
	FileSize     int64    	`gorm:"not null" json:"file_size"`
	Format       string 	`gorm:"not null" json:"format"`
	Width        int    	`gorm:"not null" json:"width"`
	Height       int    	`gorm:"not null" json:"height"`
	UserID       int64  	`gorm:"not null;column:user_id" json:"-"`
	CreatedAt    time.Time	`gorm:"autoCreateTime" json:"-"`
	UpdatedAt	 time.Time	`gorm:"autoUpdateTime" json:"-"`
}

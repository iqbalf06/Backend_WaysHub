package models

import "time"

type Video struct {
	ID          int       `json:"id" gorm:"primary_key:auto_increment"`
	Title       string    `json:"title" gorm:"type: varchar(255)"`
	Thumbnail   string    `gorm:"type: varchar(255)" json:"thumbnail"`
	Description string    `gorm:"type: varchar(255)" json:"description"`
	Video       string    `gorm:"type: varchar(255)" json:"video"`
	ChannelID   int       `json:"channel_id" `
	Channel     Channel   `json:"channel" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt   time.Time `json:"createdAt"`
	ViewCount   int       `json:"viewcount" form:"viewcount" gorm:"type: int"`
	Comment     []Comment `json:"comment" `
}

type VideoResponse struct {
	ID          int    `json:"id"`
	Title       string `json:"title" gorm:"type: varchar(255)"`
	Thumbnail   string `gorm:"type: varchar(255)" json:"thumbnail"`
	Description string `json:"description"`
	ViewCount   int    `json:"viewcount" form:"viewcount" gorm:"type: int"`
	// Channel     []Channel `json:"-"`
	CreatedAt   time.Time `json:"-"`
	Channelname string    `json:"channelName"`
	Photo       string    `json:"photo"`
	Cover       string    `json:"cover"`
}

func (VideoResponse) TableName() string {
	return "videos"
}

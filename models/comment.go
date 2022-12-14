package models

import "time"

type Comment struct {
	ID        int       `json:"id" gorm:"primary_key:auto_increment"`
	Comment   string    `gorm:"type: varchar(255)" json:"comment"`
	CreatedAt time.Time `json:"created_at"`
	ChannelID int       `json:"channel_id"`
	Channel   Channel   `json:"channel" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	VideoID   int       `json:"video_id"`
	Video     string    `json:"video" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

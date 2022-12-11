package videodto

import "time"

type VideoRequest struct {
	ID          int       `json:"id" gorm:"primary_key:auto_increment"`
	Title       string    `json:"title" gorm:"type: varchar(255)"`
	Thumbnail   string    `gorm:"type: varchar(255)" json:"thumbnail"`
	Description string    `gorm:"type: varchar(255)" json:"description"`
	Video       string    `gorm:"type: varchar(255)" json:"video"`
	ChannelID   int       `json:"channel_id"`
	CreatedAt   time.Time `json:"-"`

	ViewCount int `json:"viewcount" form:"viewcount" gorm:"type: int"`
}

type EditVideoRequest struct {
	Title       string `json:"title" gorm:"type: varchar(255)"`
	Thumbnail   string `gorm:"type: varchar(255)" json:"thumbnail"`
	Description string `gorm:"type: varchar(255)" json:"description"`
}

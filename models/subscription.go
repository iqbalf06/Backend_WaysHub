package models

type Subscription struct {
	ID            int `json:"id" gorm:"primary_key:auto_increment"`
	ChannelId     int `json:"channel_id"`
}

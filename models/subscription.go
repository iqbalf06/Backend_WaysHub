package models

type Subscription struct {
	ID        int     `json:"id" gorm:"primary_key:auto_increment"`
	ChannelID int     `json:"channel_id"`
	Channel   Channel `json:"channel"`
}

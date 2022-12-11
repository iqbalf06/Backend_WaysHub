package repositories

import (
	"wayshub/models"

	"gorm.io/gorm"
)

type VideoRepository interface {
	FindVideos() ([]models.Video, error)
	GetVideo(ID int) (models.Video, error)
	GetChannelVideo(ChannelID int) ([]models.Video, error)
	AddVideo(video models.Video) (models.Video, error)
	EditVideo(video models.Video) (models.Video, error)
	DeleteVideo(video models.Video) (models.Video, error)
	ViewCount(id int) error
}

func RepositoryVideo(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindVideos() ([]models.Video, error) {
	var videos []models.Video
	err := r.db.Preload("Channel").Find(&videos).Error

	return videos, err
}

func (r *repository) GetVideo(ID int) (models.Video, error) {
	var video models.Video
	err := r.db.First(&video, ID).Error

	return video, err
}

func (r *repository) GetChannelVideo(ChannelID int) ([]models.Video, error) {
	var video []models.Video
	err := r.db.Debug().Preload("Channel").Where("channel_id=?", ChannelID).Find(&video).Error

	return video, err
}

func (r *repository) AddVideo(video models.Video) (models.Video, error) {
	err := r.db.Create(&video).Error

	return video, err
}

func (r *repository) EditVideo(video models.Video) (models.Video, error) {
	err := r.db.Save(&video).Error

	return video, err
}

func (r *repository) ViewCount(id int) error {
	var view models.Video
	r.db.First(&view, id)
	view.ViewCount += 1
	err := r.db.Save(&view).Error
	return err
}

func (r *repository) DeleteVideo(video models.Video) (models.Video, error) {
	err := r.db.Delete(&video).Error

	return video, err
}

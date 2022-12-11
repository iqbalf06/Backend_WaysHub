package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
	dto "wayshub/dto/result"
	videodto "wayshub/dto/video"
	"wayshub/models"
	"wayshub/repositories"

	"context"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

type handlerVideo struct {
	VideoRepository repositories.VideoRepository
}

func HandlerVideo(VideoRepository repositories.VideoRepository) *handlerVideo {
	return &handlerVideo{VideoRepository}
}

func (h *handlerVideo) AddVideo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userID := int(userInfo["id"].(float64))

	// Get dataFile from midleware and store to filethumbnail variable here ...
	dataContex := r.Context().Value("dataThumbnail")
	filethumbnail := dataContex.(string)

	videoContex := r.Context().Value("dataVideo")
	filevideo := videoContex.(string)

	var ctx = context.Background()
	var CLOUD_NAME = os.Getenv("CLOUD_NAME")
	var API_KEY = os.Getenv("API_KEY")
	var API_SECRET = os.Getenv("API_SECRET")

	// Add your Cloudinary credentials ...
	cld, _ := cloudinary.NewFromParams(CLOUD_NAME, API_KEY, API_SECRET)
	respthumbnail, err := cld.Upload.Upload(ctx, filethumbnail, uploader.UploadParams{Folder: "WaysHub/thumbnail"})

	if err != nil {
		fmt.Println(err.Error())
	}

	// Upload file to Cloudinary ...
	respvideo, err := cld.Upload.Upload(ctx, filevideo, uploader.UploadParams{Folder: "WaysHub/video"})

	if err != nil {
		fmt.Println(err.Error())
	}

	request := videodto.VideoRequest{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
	}
	fmt.Println(request)
	validation := validator.New()
	err = validation.Struct(request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	video := models.Video{
		Title:       request.Title,
		Thumbnail:   respthumbnail.SecureURL,
		Description: request.Description,
		Video:       respvideo.SecureURL,
		ChannelID:   userID,
		CreatedAt:   time.Now(),
	}

	video, err = h.VideoRepository.AddVideo(video)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	video, _ = h.VideoRepository.GetVideo(video.ID)

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: "success", Data: video}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerVideo) FindVideos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	videos, err := h.VideoRepository.FindVideos()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: "success", Data: videos}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerVideo) GetVideo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var video models.Video
	video, err := h.VideoRepository.GetVideo(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	video.Video = os.Getenv("PATH_FILE") + video.Video

	videos, err := h.VideoRepository.GetVideo(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: "success", Data: videos}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerVideo) GetChannelVideo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	// userID := int(userInfo["id"].(float64))

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var video []models.Video
	video, err := h.VideoRepository.GetChannelVideo(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: "success", Data: video}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerVideo) EditVideo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	dataContex := r.Context().Value("dataThumbnail")
	filethumbnail := dataContex.(string)

	videoContex := r.Context().Value("dataVideo")
	filevideo := videoContex.(string)

	request := videodto.EditVideoRequest{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
	}

	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}
	video, _ := h.VideoRepository.GetVideo(id)

	if request.Title != "" {
		video.Title = request.Title
	}

	if filethumbnail != "false" {
		video.Thumbnail = filethumbnail
	}

	if request.Description != "" {
		video.Description = request.Description
	}

	if filevideo != "false" {
		video.Video = filevideo
	}

	video, err = h.VideoRepository.EditVideo(video)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: "success", Data: video}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerVideo) UpdateView(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	video, _ := h.VideoRepository.GetVideo(id)

	video.ViewCount = video.ViewCount + 1

	fmt.Println(video.ViewCount)

	data, err := h.VideoRepository.EditVideo(video)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: "success", Data: data}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerVideo) DeleteVideo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	video, err := h.VideoRepository.GetVideo(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	_, err = h.VideoRepository.DeleteVideo(video)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	data := models.Video{
		ID: video.ID,
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: "success", Data: convertResponseVideo(data)}
	json.NewEncoder(w).Encode(response)
}

func convertResponseVideo(u models.Video) videodto.DeleteResponse {
	return videodto.DeleteResponse{
		ID: u.ID,
	}
}

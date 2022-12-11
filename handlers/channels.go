package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	authdto "wayshub/dto/auth"
	channelsdto "wayshub/dto/channels"
	dto "wayshub/dto/result"
	"wayshub/models"
	"wayshub/repositories"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

type handlerChannel struct {
	ChannelRepository repositories.ChannelRepository
}

func HandlerChannel(ChannelRepository repositories.ChannelRepository) *handlerChannel {
	return &handlerChannel{ChannelRepository}
}

func (h *handlerChannel) FindChannels(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userID := int(userInfo["id"].(float64))

	channels, err := h.ChannelRepository.FindChannels(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
	}

	for i, p := range channels {
		channels[i].Photo = os.Getenv("PATH_FILE") + p.Photo
	}

	for i, p := range channels {
		channels[i].Cover = os.Getenv("PATH_FILE") + p.Cover
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: "success", Data: channels}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerChannel) GetChannel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	channel, err := h.ChannelRepository.GetChannel(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: "success", Data: channel}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerChannel) GetMyChannel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userID := int(userInfo["id"].(float64))

	channel, err := h.ChannelRepository.GetMyChannel(userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: "success", Data: channel}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerChannel) EditChannel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userID := int(userInfo["id"].(float64))

	ContexPhoto := r.Context().Value("dataPhoto")
	filephoto := ContexPhoto.(string)
	ContexCover := r.Context().Value("dataCover")
	filecover := ContexCover.(string)

	var ctx = context.Background()
	var CLOUD_NAME = os.Getenv("CLOUD_NAME")
	var API_KEY = os.Getenv("API_KEY")
	var API_SECRET = os.Getenv("API_SECRET")

	// Add your Cloudinary credentials ...
	cld, _ := cloudinary.NewFromParams(CLOUD_NAME, API_KEY, API_SECRET)
	respcover, err := cld.Upload.Upload(ctx, filecover, uploader.UploadParams{Folder: "WaysHub/cover"})

	if err != nil {
		fmt.Println(err.Error())
	}

	// Upload file to Cloudinary ...
	respphoto, err := cld.Upload.Upload(ctx, filephoto, uploader.UploadParams{Folder: "WaysHub/photo"})

	if err != nil {
		fmt.Println(err.Error())
	}

	request := authdto.RegisterRequest{
		Channelname: r.FormValue("channelName"),
		Email:       r.FormValue("email"),
		Password:    r.FormValue("password"),
	}

	channel, err := h.ChannelRepository.GetChannel(int(userID))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	if request.Channelname != "" {
		channel.Channelname = request.Channelname
	}

	if request.Email != "" {
		channel.Email = request.Email
	}

	if request.Password != "" {
		channel.Password = request.Password
	}

	if filephoto != "false" {
		channel.Photo = respphoto.SecureURL
	}

	if filecover != "false" {
		channel.Cover = respcover.SecureURL
	}

	data, err := h.ChannelRepository.EditChannel(channel)
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

func (h *handlerChannel) DeleteChannel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userRole := userInfo["role"]
	userID := int(userInfo["id"].(float64))

	if userID != id && userRole != "admin" {
		w.WriteHeader(http.StatusUnauthorized)
		response := dto.ErrorResult{Code: http.StatusUnauthorized, Message: "you're not admin"}
		json.NewEncoder(w).Encode(response)
		return
	}

	channel, err := h.ChannelRepository.GetChannel(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	data, err := h.ChannelRepository.DeleteChannel(channel)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: "success", Data: convertResponse(data)}
	json.NewEncoder(w).Encode(response)
}

func convertResponse(u models.Channel) channelsdto.DeleteResponse {
	return channelsdto.DeleteResponse{
		ID: u.ID,
	}
}

package routes

import (
	"wayshub/handlers"
	"wayshub/pkg/middleware"
	"wayshub/pkg/mysql"
	"wayshub/repositories"

	"github.com/gorilla/mux"
)

func CommentRoutes(r *mux.Router) {
	commentRepository := repositories.RepositoryComment(mysql.DB)
	h := handlers.HandlerComment(commentRepository)

	r.HandleFunc("/video/comments/{id}", middleware.Auth(h.AddComment)).Methods("POST")
	r.HandleFunc("/video/comments", middleware.Auth(h.FindComments)).Methods("GET")
	r.HandleFunc("/video/comment/{id}", middleware.Auth(h.GetCommentbyVideo)).Methods("GET")
	r.HandleFunc("/video/{id}/comment/{id}", middleware.Auth(h.GetComment)).Methods("GET")
	r.HandleFunc("/video/{id}/comment/{id}", middleware.Auth(h.DeleteComment)).Methods("DELETE")
}

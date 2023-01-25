package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jscastaneda-esp/rest-ws-go/models"
	"github.com/jscastaneda-esp/rest-ws-go/repository"
	"github.com/jscastaneda-esp/rest-ws-go/server"
	"github.com/jscastaneda-esp/rest-ws-go/services"
	"github.com/segmentio/ksuid"
)

type UpsertPostRequest struct {
	PostContent string `json:"postContent"`
}

type CreatePostResponse struct {
	Id          string `json:"id"`
	PostContent string `json:"postContent"`
}

type GenericResponse struct {
	Message string `json:"message"`
}

func CreatePostHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, status, err := services.GetClaimsToken(r.Header.Get("Authorization"), s.Config().JWTSecret)
		if err != nil {
			http.Error(w, err.Error(), status)
			log.Println("GetUserData:", err)
			return
		}

		request := new(UpsertPostRequest)
		if err := json.NewDecoder(r.Body).Decode(request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Println("Json Decode:", err)
			return
		}

		id, err := ksuid.NewRandom()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println("NewRandom:", err)
			return
		}

		post := &models.Post{
			BaseModel: models.BaseModel{
				Id: id.String(),
			},
			PostContent: request.PostContent,
			UserId:      claims.UserId,
		}
		err = repository.InsertPost(r.Context(), post)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println("InsertPost:", err)
			return
		}

		postMessage := &models.WebSocketMessage{
			Type:    "Post_Created",
			Payload: post,
		}
		s.Hub().Broadcast(postMessage, nil)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(CreatePostResponse{
			Id:          post.Id,
			PostContent: post.PostContent,
		})
	}
}

func GetPostByIdHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		post, err := repository.GetPostById(r.Context(), params["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println("GetPostById:", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(post)
	}
}

func ListPostHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pageStr := r.URL.Query().Get("page")
		if pageStr == "" {
			pageStr = "1"
		}

		rowsFetchStr := r.URL.Query().Get("rowsFetch")
		if rowsFetchStr == "" {
			rowsFetchStr = s.Config().RowsDefault
		}

		page, err := strconv.ParseUint(pageStr, 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Println("ParseUint:", err)
			return
		}

		rowsFetch, err := strconv.ParseUint(rowsFetchStr, 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Println("ParseUint:", err)
			return
		}

		posts, err := repository.ListPosts(r.Context(), page, rowsFetch)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println("ListPosts:", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(posts)
	}
}

func UpdatePostHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, status, err := services.GetClaimsToken(r.Header.Get("Authorization"), s.Config().JWTSecret)
		if err != nil {
			http.Error(w, err.Error(), status)
			log.Println("GetUserData:", err, status)
		}

		params := mux.Vars(r)
		request := new(UpsertPostRequest)
		if err := json.NewDecoder(r.Body).Decode(request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Println("Json Decode:", err)
			return
		}

		post := &models.Post{
			BaseModel: models.BaseModel{
				Id: params["id"],
			},
			PostContent: request.PostContent,
			UserId:      claims.UserId,
		}
		err = repository.UpdatePost(r.Context(), post)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println("UpdatePost:", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(GenericResponse{
			Message: "Post updated",
		})
	}
}

func DeletePostHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, status, err := services.GetClaimsToken(r.Header.Get("Authorization"), s.Config().JWTSecret)
		if err != nil {
			http.Error(w, err.Error(), status)
			log.Println("GetUserData:", err, status)
			return
		}

		params := mux.Vars(r)
		err = repository.DeletePost(r.Context(), params["id"], claims.UserId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println("DeletePost:", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(GenericResponse{
			Message: "Post deleted",
		})
	}
}

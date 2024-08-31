package main

import (
	"database/sql"
	"errors"
	"github.com/go-chi/chi/v5"
	db "github.com/trenchesdeveloper/social-blue/internal/db/sqlc"
	"github.com/trenchesdeveloper/social-blue/internal/dto"
	"log"
	"net/http"
	"strconv"
)

func (s *server) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var input dto.Post
	if err := readJSON(w, r, &input); err != nil {
		s.badRequestError(w, r, err)
		return
	}

	post, err := s.store.CreatePost(r.Context(), db.CreatePostParams{
		Title:   input.Title,
		Content: input.Content,
		//TODO: get the user id from the request
		UserID: 1,
	})

	if err != nil {
		s.internalServerError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, post)

}

func (s *server) getPostHandler(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "postID")

	log.Println("postID: ", postID)

	id, err := strconv.ParseInt(postID, 10, 64)
	if err != nil {
		s.internalServerError(w, r, err)
		return
	}
	post, err := s.store.GetPostByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.notFoundError(w, r)
			return
		}
		s.internalServerError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, post)

}

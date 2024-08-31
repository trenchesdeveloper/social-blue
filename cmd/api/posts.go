package main

import (
	db "github.com/trenchesdeveloper/social-blue/internal/db/sqlc"
	"github.com/trenchesdeveloper/social-blue/internal/dto"
	"net/http"
)

func (s *server) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var input dto.Post
	if err := readJSON(w, r, &input); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	post, err := s.store.CreatePost(r.Context(), db.CreatePostParams{
		Title:   input.Title,
		Content: input.Content,
		//TODO: get the user id from the request
		UserID: 1,
	})

	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, post)

}

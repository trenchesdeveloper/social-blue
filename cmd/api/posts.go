package main

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-chi/chi/v5"
	db "github.com/trenchesdeveloper/social-blue/internal/db/sqlc"
	"github.com/trenchesdeveloper/social-blue/internal/dto"
	"log"
	"net/http"
	"strconv"
)

type postContextKey string

const postKey postContextKey = "post"

// CreatePost godoc
//
//	@Summary		Create a post
//	@Description	Create a post
//	@ID				create-post
//	@Accept			json
//	@Produce		json
//	@Param			input	body		dto.CreatPostDto	true	"Post data"
//	@Success		200		{object}	db.CreatePostParams
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
func (s *server) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var input dto.CreatPostDto
	if err := readJSON(w, r, &input); err != nil {
		s.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(input); err != nil {
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

	jsonResponse(w, http.StatusOK, post)

}

func (s *server) getPostHandler(w http.ResponseWriter, r *http.Request) {
	post, err := s.getPostFromCtx(r.Context())

	if err != nil {
		s.internalServerError(w, r, err)
		return
	}

	// fetch the comments
	comments, err := s.store.GetCommentsByPostID(r.Context(), post.ID)
	if err != nil {
		s.internalServerError(w, r, err)
		return
	}

	postWithComments := dto.GetPostWithCommentsDto{
		ID:        post.ID,
		Content:   post.Content,
		Title:     post.Title,
		UserID:    post.UserID,
		Tags:      post.Tags,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
		Comments:  comments,
	}

	jsonResponse(w, http.StatusOK, postWithComments)

}

func (s *server) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	post, err := s.getPostFromCtx(r.Context())

	if err != nil {
		s.internalServerError(w, r, err)
		return
	}

	var input dto.UpdatePostDto
	if err := readJSON(w, r, &input); err != nil {
		s.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(input); err != nil {
		s.badRequestError(w, r, err)
		return
	}

	log.Println("current post", post)

	UpdatedPost, err := s.store.UpdatePost(r.Context(), db.UpdatePostParams{
		ID:      post.ID,
		Column2: input.Content,
		Column3: input.Title,
		Column4: input.Tags,
		Version: post.Version,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.notFoundError(w, r)
			return
		}
		s.internalServerError(w, r, err)
		return
	}

	jsonResponse(w, http.StatusOK, UpdatedPost)
}

func (s *server) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "postID")

	id, err := strconv.ParseInt(postID, 10, 64)
	if err != nil {
		s.internalServerError(w, r, err)
		return
	}

	err = s.store.DeletePost(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.notFoundError(w, r)
			return
		}
		s.internalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *server) postsContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		postID := chi.URLParam(r, "postID")

		log.Println("postID: ", postID)

		id, err := strconv.ParseInt(postID, 10, 64)
		if err != nil {
			s.internalServerError(w, r, err)
			return
		}
		post, err := s.store.GetPostByID(r.Context(), id)
		log.Println("post from db", post)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				s.notFoundError(w, r)
				return
			}
			s.internalServerError(w, r, err)
			return
		}

		ctx = context.WithValue(ctx, postKey, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *server) getPostFromCtx(ctx context.Context) (db.GetPostByIDRow, error) {
	post, ok := ctx.Value(postKey).(db.GetPostByIDRow)

	if !ok {
		return db.GetPostByIDRow{}, errors.New("post not found")
	}
	log.Println("post from context", post)
	return post, nil
}

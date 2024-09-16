package main

import (
	"database/sql"
	db "github.com/trenchesdeveloper/social-blue/internal/db/sqlc"
	"github.com/trenchesdeveloper/social-blue/internal/dto"
	"net/http"
)

func (s *server) GetUserFeedsHandler(w http.ResponseWriter, r *http.Request) {
	// get the user from the context
	//_, err := s.getUserFromContext(r.Context())
	//if err != nil {
	//	s.internalServerError(w, r, err)
	//	return
	//}

	fq := PaginatedFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
	}

	fq, err := fq.Parse(r)
	if err != nil {
		s.badRequestError(w, r, err)
		return
	}

	// validate the query
	if err := Validate.Struct(fq); err != nil {
		s.badRequestError(w, r, err)
		return
	}

	// get the user feeds
	feeds, err := s.store.GetUserFeed(r.Context(), db.GetUserFeedParams{
		Limit:  int32(fq.Limit),
		Offset: int32(fq.Offset),
		UserID: 1,
		Column4: sql.NullString{
			String: fq.Search,
			Valid:  fq.Search != "",
		},
		Tags: fq.Tags,
	})
	if err != nil {
		s.internalServerError(w, r, err)
		return
	}

	postWithMedata := []dto.PostWithMetadata{}
	for _, feed := range feeds {
		postWithMedata = append(postWithMedata, dto.PostWithMetadata{
			GetPostWithCommentsDto: dto.GetPostWithCommentsDto{
				ID:        feed.ID,
				Content:   feed.Content,
				Title:     feed.Title,
				UserID:    feed.UserID,
				Tags:      feed.Tags,
				CreatedAt: feed.CreatedAt,
				UpdatedAt: feed.UpdatedAt,
				User: db.User{
					ID:       feed.UserID,
					Username: feed.Username.String,
				},
			},
		})
	}

	jsonRespose(w, http.StatusOK, postWithMedata)
}

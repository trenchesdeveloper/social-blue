// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"context"
)

type Querier interface {
	CreatePost(ctx context.Context, arg CreatePostParams) (CreatePostRow, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (CreateUserRow, error)
	DeletePost(ctx context.Context, id int64) error
	DeleteUser(ctx context.Context, id int64) error
	GetCommentsByPostID(ctx context.Context, postID int64) ([]GetCommentsByPostIDRow, error)
	GetPostByID(ctx context.Context, id int64) (GetPostByIDRow, error)
	GetUserByID(ctx context.Context, id int64) (GetUserByIDRow, error)
	GetUserByUsername(ctx context.Context, username string) (GetUserByUsernameRow, error)
	ListPosts(ctx context.Context) ([]ListPostsRow, error)
	ListUsers(ctx context.Context) ([]ListUsersRow, error)
	UpdatePost(ctx context.Context, arg UpdatePostParams) (UpdatePostRow, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (UpdateUserRow, error)
}

var _ Querier = (*Queries)(nil)

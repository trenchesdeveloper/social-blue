package dto

type UserResponseDto struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type RegisterUserDto struct {
	FirstName string `json:"first_name" validate:"required, max=100"`
	LastName  string `json:"last_name" validate:"required, max=100"`
	Username  string `json:"username" validate:"required, max=100"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6"`
}

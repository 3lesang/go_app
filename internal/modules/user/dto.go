package user

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Phone    string `json:"phone" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username"`
	Password string `json:"password" validate:"required"`
}

type DeleteUsersRequest struct {
	IDs []int64 `json:"ids"`
}

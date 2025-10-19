package auth

type SignInParams struct {
	Identify string `json:"identify"`
	Password string `json:"password"`
}

type SignUpParams struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"email"`
	Phone    string `json:"phone" validate:"required"`
	Username string `json:"username"`
	Password string `json:"password" validate:"required"`
}

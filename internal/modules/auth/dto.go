package auth

type LoginRequest struct {
	Identify string `json:"identify" example:"admin"`
	Password string `json:"password" example:"admin"`
}

type RegisterRequest struct {
	Name     string `json:"name" validate:"required" example:"Admin"`
	Email    string `json:"email" validate:"email" example:"admin@gmail.com"`
	Phone    string `json:"phone" validate:"required" example:"0984807356"`
	Username string `json:"username" example:"admin"`
	Password string `json:"password" validate:"required" example:"admin"`
}

package customer

type PaginatedResponse[T any] struct {
	Page       int   `json:"page" example:"1"`
	PageSize   int   `json:"page_size" example:"10"`
	TotalItems int64 `json:"total_items" example:"125"`
	TotalPages int   `json:"total_pages" example:"13"`
	Data       []T   `json:"data"`
}

type CreateCustomerRequest struct {
	Name     string `json:"name" validate:"required"`
	Phone    string `json:"phone" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type DeleteCustomersRequest struct {
	IDs []int64 `json:"ids"`
}

type CustomerLoginRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type UpdateMeRequest struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

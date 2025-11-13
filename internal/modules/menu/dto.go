package menu

type PaginatedResponse[T any] struct {
	Page       int   `json:"page" example:"1"`
	PageSize   int   `json:"page_size" example:"10"`
	TotalItems int64 `json:"total_items" example:"125"`
	TotalPages int   `json:"total_pages" example:"13"`
	Data       []T   `json:"data"`
}

type CreateMenuRequest struct {
	Name     string `json:"name" validate:"required"`
	Position string `json:"position" validate:"required"`
}

type UpdateMenuRequest struct {
	Name     string `json:"name" validate:"required"`
	Position string `json:"position" validate:"required"`
}

type DeleteMenusRequest struct {
	IDs []int64 `json:"ids"`
}

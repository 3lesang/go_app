package collection

type PaginatedResponse[T any] struct {
	Page       int   `json:"page" example:"1"`
	PageSize   int   `json:"page_size" example:"10"`
	TotalItems int64 `json:"total_items" example:"125"`
	TotalPages int   `json:"total_pages" example:"13"`
	Data       []T   `json:"data"`
}

type CollectionResponse struct {
	ID              int64  `json:"id"`
	Name            string `json:"name"`
	Slug            string `json:"slug"`
	File            string `json:"file"`
	MetaTitle       string `json:"meta_title"`
	MetaDescription string `json:"meta_description"`
	Layout          string `json:"layout"`
	Products        any    `json:"products"`
}

type CreateCollectionRequest struct {
	Name            string  `json:"name" validate:"required"`
	Slug            string  `json:"slug" validate:"required"`
	File            string  `json:"file"`
	MetaTitle       string  `json:"meta_title"`
	MetaDescription string  `json:"meta_decscription"`
	Layout          string  `json:"layout"`
	Conditions      string  `json:"conditions"`
	ProductIDs      []int64 `json:"product_ids"`
}

type UpdateCollectionRequest struct {
	Name            string  `json:"name" validate:"required"`
	Slug            string  `json:"slug" validate:"required"`
	File            string  `json:"file"`
	MetaTitle       string  `json:"meta_title"`
	MetaDescription string  `json:"meta_decscription"`
	Layout          string  `json:"layout"`
	Conditions      string  `json:"conditions"`
	ProductIDs      []int64 `json:"product_ids"`
}

type DeleteCollectionsRequest struct {
	IDs []int64 `json:"ids"`
}

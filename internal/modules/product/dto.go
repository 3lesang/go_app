package product

type CreateProductRequest struct {
	Name        string `json:"name" validate:"required"`
	Slug        string `json:"slug" validate:"required"`
	OriginPrice int64  `json:"origin_price" validate:"required"`
}

type UpdateProductRequest struct {
	Name        string `json:"name" validate:"required"`
	Slug        string `json:"slug" validate:"required"`
	OriginPrice int64  `json:"origin_price" validate:"required"`
}

type DeleteProductsRequest struct {
	IDs []int64 `json:"ids"`
}

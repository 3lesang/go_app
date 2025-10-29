package collection

type CreateCollectionRequest struct {
	Name string `json:"name" validate:"required"`
	Slug string `json:"slug" validate:"required"`
}

type UpdateCollectionRequest struct {
	Name string `json:"name" validate:"required"`
	Slug string `json:"slug" validate:"required"`
}

type DeleteCollectionsRequest struct {
	IDs []int64 `json:"ids"`
}

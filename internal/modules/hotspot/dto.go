package hotspot

type PaginatedResponse[T any] struct {
	Page       int   `json:"page" example:"1"`
	PageSize   int   `json:"page_size" example:"10"`
	TotalItems int64 `json:"total_items" example:"125"`
	TotalPages int   `json:"total_pages" example:"13"`
	Data       []T   `json:"data"`
}

type CreateSpot struct {
	X         float32 `json:"x"`
	Y         float32 `json:"y"`
	ProductID int64   `json:"product_id"`
}

type CreateHotspotRequest struct {
	File  string       `json:"file"`
	Spots []CreateSpot `json:"spots"`
}

type UpdateSpot struct {
	ID        int64   `json:"id"`
	X         float32 `json:"x"`
	Y         float32 `json:"y"`
	ProductID int64   `json:"product_id"`
}

type UpdateHotspotRequest struct {
	Spots []UpdateSpot `json:"spots"`
}

type DeleteHotspotsRequest struct {
	IDs []int64 `json:"ids"`
}

package review

type PaginatedResponse[T any] struct {
	Page       int   `json:"page" example:"1"`
	PageSize   int   `json:"page_size" example:"10"`
	TotalItems int64 `json:"total_items" example:"125"`
	TotalPages int   `json:"total_pages" example:"13"`
	Data       []T   `json:"data"`
}

type CreateReviewRequest struct {
	Rating     int    `json:"rating" validate:"gte=1,lte=5" example:"5"`
	Comment    string `json:"comment"`
	ProductID  int64  `json:"product_id" form:"product_id"`
	CustomerID int64  `json:"customer_id" form:"customer_id"`
}

type AverageRatingResponse struct {
	AverageRating float64  `json:"average_rating" example:"4.2"`
	TotalReviews  int64    `json:"total_reviews" example:"15"`
	TotalFiles    int64    `json:"total_files"`
	Files         []string `json:"files"`
}

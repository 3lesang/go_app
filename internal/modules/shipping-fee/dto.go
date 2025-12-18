package shippingfee

type PaginatedResponse[T any] struct {
	Page       int   `json:"page" example:"1"`
	PageSize   int   `json:"page_size" example:"10"`
	TotalItems int64 `json:"total_items" example:"125"`
	TotalPages int   `json:"total_pages" example:"13"`
	Data       []T   `json:"data"`
}

type ShippingFeeResponse struct {
	ID            int64 `json:"id"`
	MinWeight     int32 `json:"min_weight"`
	MaxWeight     int32 `json:"max_weight"`
	FeeAmount     int32 `json:"fee_amount"`
	MinOrderValue int32 `json:"min_order_value"`
	FreeShipping  bool  `json:"free_shipping"`
}

type CreateShippingFeeRequest struct {
	MinWeight     int32 `json:"min_weight"`
	MaxWeight     int32 `json:"max_weight"`
	FeeAmount     int32 `json:"fee_amount"`
	MinOrderValue int32 `json:"min_order_value"`
	FreeShipping  bool  `json:"free_shipping"`
}

type UpdateShippingFeeRequest struct {
	MinWeight     int32 `json:"min_weight"`
	MaxWeight     int32 `json:"max_weight"`
	FeeAmount     int32 `json:"fee_amount"`
	MinOrderValue int32 `json:"min_order_value"`
	FreeShipping  bool  `json:"free_shipping"`
}

type DeleteShippingFeesRequest struct {
	IDs []int64 `json:"ids"`
}

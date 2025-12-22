package order

type PaginatedResponse[T any] struct {
	Page       int   `json:"page" example:"1"`
	PageSize   int   `json:"page_size" example:"10"`
	TotalItems int64 `json:"total_items" example:"125"`
	TotalPages int   `json:"total_pages" example:"13"`
	Data       []T   `json:"data"`
}

type CreateOrderItems struct {
	Quantity  int32 `json:"quantity"`
	SalePrice int32 `json:"sale_price"`
	ProductID int64 `json:"product_id"`
	VariantID int64 `json:"variant_id"`
}

type CreateOrderAddress struct {
	FullName    string `json:"full_name"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	AddressLine string `json:"address_line"`
}

type CreateOrderRequest struct {
	TotalAmount       int32              `json:"total_amount"`
	DiscountAmount    int32              `json:"discount_amount"`
	ShippingFeeAmount int32              `json:"shipping_fee_amount"`
	Address           CreateOrderAddress `json:"address"`
	Items             []CreateOrderItems `json:"items"`
}

type UpdateOrderRequest struct {
	Status       string `json:"status"`
	CancelReason string `json:"cancel_reason"`
}

type DeleteOrdersRequest struct {
	IDs []int64 `json:"ids"`
}

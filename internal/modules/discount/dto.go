package discount

import "time"

type PaginatedResponse[T any] struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
	Data       []T `json:"data"`
}

type CreateDiscountRequest struct {
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	Code             string    `json:"code"`
	DiscountType     string    `json:"discount_type" validate:"required,oneof=code automatic"`
	Status           string    `json:"status" validate:"required,oneof=draft active scheduled expired"`
	UsageLimit       int32     `json:"usage_limit"`
	PerCustomerLimit int32     `json:"per_customer_limit"`
	StartsAt         time.Time `json:"starts_at" validate:"required"`
	EndsAt           time.Time `json:"ends_at"`
}

type CreateDiscountTargetRequest struct {
	TargetType string  `json:"target_type"`
	IDs        []int64 `json:"ids"`
}

type CreateDiscountEffectRequest struct {
	EffectType string `json:"effect_type"`
	Value      string `json:"value"`
	AppliesTo  string `json:"applies_to"`
}

type CreateDiscountConditionRequest struct {
	ConditionType string `json:"condition_type"`
	Operator      string `json:"operator"`
	Value         string `json:"value"`
}

type UpdateDiscountRequest struct {
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	Status           string    `json:"status" validate:"required,oneof=draft active scheduled expired"`
	EndsAt           time.Time `json:"ends_at"`
	UsageLimit       int32     `json:"usage_limit"`
	PerCustomerLimit int32     `json:"per_customer_limit"`
	StartsAt         time.Time `json:"starts_at"`
}

type UpdateDiscountTargetRequest struct {
	TargetType string `json:"target_type"`
}

type UpdateDiscountEffectRequest struct {
	EffectType string `json:"effect_type"`
	Value      string `json:"value"`
	AppliesTo  string `json:"applies_to"`
}

type UpdateDiscountConditionRequest struct {
	ConditionType string `json:"condition_type"`
	Value         string `json:"value"`
}

type OneDiscountResponse struct {
	ID               int64     `json:"id"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	Code             string    `json:"code"`
	DiscountType     string    `json:"discount_type"`
	Status           string    `json:"status"`
	UsageLimit       int       `json:"usage_limit"`
	PerCustomerLimit int       `json:"per_customer_limit"`
	StartsAt         time.Time `json:"start_at"`
	EndsAt           time.Time `json:"ends_at"`
}

type BulkDeleteDiscountsRequest struct {
	IDs []int64 `json:"ids" validate:"required,dive,gt=0"`
}

type UpsertCustomerUsageRequest struct {
	DiscountID int64 `json:"discount_id" validate:"required"`
	CustomerID int64 `json:"customer_id" validate:"required"`
}

package product

type PaginatedResponse[T any] struct {
	Page       int   `json:"page" example:"1"`
	PageSize   int   `json:"page_size" example:"10"`
	TotalItems int64 `json:"total_items" example:"125"`
	TotalPages int   `json:"total_pages" example:"13"`
	Data       []T   `json:"data"`
}

type Option struct {
	ID        int64         `json:"id"`
	Name      string        `json:"name"`
	No        int32         `json:"no"`
	ProductID int64         `json:"product_id"`
	Values    []OptionValue `json:"values"`
}

type OptionValue struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	No       int32  `json:"no"`
	OptionID int64  `json:"option_id"`
}

type VariantOption struct {
	OptionID   int64  `json:"option_id"`
	OptionName string `json:"option_name"`
	ValueID    int64  `json:"value_id"`
	Value      string `json:"value"`
}

type OneVariant struct {
	ID          int64           `json:"id"`
	OriginPrice int32           `json:"origin_price"`
	SalePrice   int32           `json:"sale_price"`
	Stock       int32           `json:"stock"`
	SKU         string          `json:"sku"`
	File        string          `json:"file"`
	Options     []VariantOption `json:"options"`
}

type OneProductResponse struct {
	ID              int64        `json:"id"`
	Name            string       `json:"name"`
	Slug            string       `json:"slug"`
	OriginPrice     int32        `json:"origin_price"`
	SalePrice       int32        `json:"sale_price"`
	Stock           int32        `json:"stock"`
	SKU             string       `json:"sku"`
	MetaTitle       string       `json:"meta_title"`
	MetaDescription string       `json:"meta_description"`
	IsActive        bool         `json:"is_active"`
	CategoryID      *int64       `json:"category_id"`
	Files           any          `json:"files"`
	Tags            any          `json:"tags"`
	Options         []Option     `json:"options"`
	Variants        []OneVariant `json:"variants"`
	Collections     any          `json:"collections"`
}

type CreateVariant struct {
	OriginPrice int32           `json:"origin_price"`
	SalePrice   int32           `json:"sale_price"`
	Stock       int32           `json:"stock"`
	Sku         string          `json:"sku"`
	File        string          `json:"file"`
	No          int32           `json:"no"`
	Options     []VariantOption `json:"options"`
}

type CreateOptionValue struct {
	Name string `json:"name"`
}

type CreateOptions struct {
	Name   string              `json:"name"`
	Values []CreateOptionValue `json:"values"`
}

type ProductFiles struct {
	No        int32  `json:"no"`
	IsPrimary bool   `json:"is_primary"`
	Name      string `json:"name"`
}

type CreateProductRequest struct {
	Name            string          `json:"name" validate:"required"`
	Slug            string          `json:"slug" validate:"required"`
	OriginPrice     int32           `json:"origin_price" validate:"gte=0"`
	SalePrice       int32           `json:"sale_price" validate:"gte=0"`
	Stock           int32           `json:"stock"`
	SKU             string          `json:"sku"`
	MetaTitle       string          `json:"meta_title"`
	MetaDescription string          `json:"meta_description"`
	CategoryID      int64           `json:"category_id"`
	Tags            []string        `json:"tags"`
	Files           []ProductFiles  `json:"files"`
	CollectionIDs   []int64         `json:"collection_ids"`
	Options         []CreateOptions `json:"options"`
	Variants        []CreateVariant `json:"variants"`
}

type UpdateOptionValue struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type UpdateOptions struct {
	ID     int64               `json:"id"`
	Name   string              `json:"name"`
	Values []UpdateOptionValue `json:"values"`
}

type UpdateVariants struct {
	ID          int64           `json:"id"`
	OriginPrice int32           `json:"origin_price"`
	SalePrice   int32           `json:"sale_price"`
	Stock       int32           `json:"stock"`
	Sku         string          `json:"sku"`
	File        string          `json:"file"`
	Options     []VariantOption `json:"options"`
}

type UpdateProductRequest struct {
	Name            string           `json:"name" validate:"required"`
	Slug            string           `json:"slug" validate:"required"`
	OriginPrice     int32            `json:"origin_price" validate:"gte=0"`
	SalePrice       int32            `json:"sale_price" validate:"gte=0"`
	Stock           int32            `json:"stock"`
	SKU             string           `json:"sku"`
	MetaTitle       string           `json:"meta_title"`
	MetaDescription string           `json:"meta_description"`
	CategoryID      int64            `json:"category_id"`
	Tags            []string         `json:"tags"`
	Files           []ProductFiles   `json:"files"`
	CollectionIDs   []int64          `json:"collection_ids"`
	Options         []UpdateOptions  `json:"options"`
	Variants        []UpdateVariants `json:"variants"`
}

type DeleteProductsRequest struct {
	IDs []int64 `json:"ids"`
}

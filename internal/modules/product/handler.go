package product

import (
	"app/internal/db"
	product_db "app/internal/db/product"
	"context"
	"math"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
)

// GetProductsHandler godoc
// @Summary      Get product list
// @Description  Returns a list of products
// @Tags         products
// @Security BearerAuth
// @Produce      json
// @Param        page      query     int     false  "Page number"  default(1)
// @Param        page_size query     int     false  "Page size"    default(10)
// @Success      200  {object}  PaginatedResponse[any]
// @Router       /products [get]
func GetProductsHandler(c *fiber.Ctx) error {

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	offset := (page - 1) * pageSize

	ctx := context.Background()

	results, err := db.ProductQueries.GetProducts(ctx, product_db.GetProductsParams{Limit: int32(pageSize), Offset: int32(offset)})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	total, err := db.ProductQueries.CountProducts(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return c.JSON(PaginatedResponse[product_db.GetProductsRow]{
		Page:       page,
		PageSize:   pageSize,
		TotalItems: total,
		TotalPages: totalPages,
		Data:       results,
	})
}

// GetProductHandler godoc
// @Summary      Get a product
// @Description  Returns a product by ID
// @Tags         products
// @Security BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  OneProductResponse
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /products/{id} [get]
func GetProductHandler(c *fiber.Ctx) error {
	param := c.Params("id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid param",
		})
	}
	ctx := context.Background()
	product, err := db.ProductQueries.GetProduct(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	files, _ := db.ProductQueries.GetFilesByProductID(ctx, product.ID)
	tags, _ := db.ProductQueries.GetTagsByProductID(ctx, product.ID)
	optionsDB, _ := db.ProductQueries.GetOptionsByProductID(ctx, product.ID)

	optionIDs := make([]int64, len(optionsDB))
	for i, o := range optionsDB {
		optionIDs[i] = o.ID
	}

	valuesDB, _ := db.ProductQueries.GetOptionValuesByOptionIDs(ctx, optionIDs)

	valuesMap := make(map[int64][]OptionValue)
	for _, v := range valuesDB {
		valuesMap[v.OptionID] = append(valuesMap[v.OptionID], OptionValue{
			ID:       v.ID,
			Name:     v.Name,
			No:       v.No.Int32,
			OptionID: v.OptionID,
		})
	}

	var optionsWithValues []Option
	for _, o := range optionsDB {
		option := Option{
			ID:     o.ID,
			Name:   o.Name,
			No:     o.No,
			Values: valuesMap[o.ID],
		}
		optionsWithValues = append(optionsWithValues, option)
	}

	variantsDB, _ := db.ProductQueries.GetVariantsByProductID(ctx, product.ID)

	variantIDs := make([]int64, len(variantsDB))
	for i, v := range variantsDB {
		variantIDs[i] = v.ID
	}

	optionRows, _ := db.ProductQueries.GetVariantOptionsByVariantIDs(ctx, variantIDs)

	optionsMap := make(map[int64][]VariantOption)
	for _, row := range optionRows {
		optionsMap[int64(row.VariantID)] = append(optionsMap[int64(row.VariantID)], VariantOption{
			OptionID:   row.OptionID,
			OptionName: row.OptionName,
			ValueID:    row.ValueID,
			Value:      row.ValueName,
		})
	}

	variants := make([]OneVariant, len(variantsDB))
	for i, v := range variantsDB {
		variants[i] = OneVariant{
			ID:          v.ID,
			SKU:         v.Sku,
			OriginPrice: v.OriginPrice,
			SalePrice:   v.SalePrice,
			Stock:       v.Stock,
			Options:     optionsMap[v.ID],
			File:        v.File.String,
		}
	}

	collections, _ := db.ProductQueries.GetCollectionsByProductID(ctx, id)

	return c.Status(fiber.StatusOK).JSON(OneProductResponse{
		ID:              product.ID,
		Name:            product.Name,
		Slug:            product.Slug,
		OriginPrice:     product.OriginPrice,
		SalePrice:       product.SalePrice,
		Stock:           product.Stock.Int32,
		SKU:             product.Sku.String,
		Weight:          product.Weight.Int32,
		Long:            product.Long.Int32,
		Wide:            product.Wide.Int32,
		High:            product.High.Int32,
		MetaTitle:       product.MetaTitle,
		MetaDescription: product.MetaDescription,
		CategoryID:      &product.CategoryID.Int64,
		Files:           files,
		Tags:            tags,
		Options:         optionsWithValues,
		Variants:        variants,
		Collections:     collections,
		IsActive:        product.IsActive,
	})
}

// GetProductBySlugHandler godoc
// @Summary      Get a product
// @Description  Returns a product by slug
// @Tags         products
// @Produce      json
// @Param        slug   path      string  true  "Product slug"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /products/slug/{slug} [get]
func GetProductBySlugHandler(c *fiber.Ctx) error {
	param := c.Params("slug")
	ctx := context.Background()
	result, err := db.ProductQueries.GetProductBySlug(ctx, param)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(result)
}

// GetProductByCategoryHandler godoc
// @Summary      Get a product
// @Description  Returns a product by category id
// @Tags         products
// @Produce      json
// @Param        id   path      string  true  "Category id"
// @Param        page      query     int     false  "Page number"  default(1)
// @Param        page_size query     int     false  "Page size"    default(10)
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /products/categories/{id} [get]
func GetProductByCategoryHandler(c *fiber.Ctx) error {
	param := c.Params("id")
	id, _ := strconv.ParseInt(param, 10, 64)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	offset := (page - 1) * pageSize

	ctx := context.Background()
	result, err := db.ProductQueries.GetProductsByCategory(ctx, product_db.GetProductsByCategoryParams{
		CategoryID: pgtype.Int8{Int64: id, Valid: true},
		Limit:      int32(pageSize),
		Offset:     int32(offset),
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	total, err := db.ProductQueries.CountProductsByCategory(ctx, pgtype.Int8{Int64: id, Valid: true})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return c.Status(fiber.StatusOK).JSON(PaginatedResponse[product_db.GetProductsByCategoryRow]{
		Page:       page,
		PageSize:   pageSize,
		TotalItems: total,
		TotalPages: totalPages,
		Data:       result,
	})
}

// CreateProductHandler godoc
// @Summary      Create a new product
// @Description  Creates a new product and returns the created product
// @Tags         products
// @Security BearerAuth
// @Accept       json
// @Produce      json
// @Param        payload  body	CreateProductRequest  true  "Create product data"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /products [post]
func CreateProductHandler(c *fiber.Ctx) error {
	var req CreateProductRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"message": err.Error(),
		})
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	ctx := context.Background()
	params := product_db.CreateProductParams{
		Name:        req.Name,
		Slug:        req.Slug,
		OriginPrice: req.OriginPrice,
		SalePrice:   req.SalePrice,
		Stock:       pgtype.Int4{Int32: req.Stock, Valid: true},
		Sku:         pgtype.Text{String: req.SKU, Valid: true},
		Weight: pgtype.Int4{
			Int32: req.Weight,
			Valid: true,
		},
		Long: pgtype.Int4{Int32: req.Long, Valid: true},
		Wide: pgtype.Int4{Int32: req.Wide, Valid: true},
		High: pgtype.Int4{Int32: req.High, Valid: true},
		CategoryID: pgtype.Int8{
			Int64: req.CategoryID,
			Valid: req.CategoryID > 0,
		},
		MetaTitle:       req.MetaTitle,
		MetaDescription: req.MetaDescription,
	}

	productID, err := db.ProductQueries.CreateProduct(ctx, params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if len(req.Files) > 0 {
		fileParams := product_db.BulkInsertProductFilesParams{}
		for _, f := range req.Files {
			fileParams.Names = append(fileParams.Names, f.Name)
			fileParams.IsPrimaries = append(fileParams.IsPrimaries, f.IsPrimary)
			fileParams.Nos = append(fileParams.Nos, f.No)
			fileParams.ProductIds = append(fileParams.ProductIds, productID)
		}
		if err := db.ProductQueries.BulkInsertProductFiles(ctx, fileParams); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	if len(req.Tags) > 0 {
		tagParams := product_db.BulkInsertProductTagsParams{}
		for _, t := range req.Tags {
			tagParams.Names = append(tagParams.Names, t)
			tagParams.ProductIds = append(tagParams.ProductIds, productID)
		}
		if err := db.ProductQueries.BulkInsertProductTags(ctx, tagParams); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	if len(req.Options) > 0 {
		optionParams := product_db.BulkInsertOptionsParams{}
		for i, o := range req.Options {
			optionParams.Names = append(optionParams.Names, o.Name)
			optionParams.Nos = append(optionParams.Nos, int32(i))
			optionParams.ProductIds = append(optionParams.ProductIds, productID)
		}
		options, _ := db.ProductQueries.BulkInsertOptions(ctx, optionParams)
		optionValueParams := product_db.BulkInsertOptionValuesParams{}

		optionMap := make(map[string]int64)
		for _, oDB := range options {
			optionMap[oDB.Name] = oDB.ID
		}

		for _, o := range req.Options {
			optionID := optionMap[o.Name]
			for i, v := range o.Values {
				optionValueParams.Names = append(optionValueParams.Names, v.Name)
				optionValueParams.Nos = append(optionValueParams.Nos, int32(i))
				optionValueParams.OptionIds = append(optionValueParams.OptionIds, optionID)
			}
		}

		optionValues, err := db.ProductQueries.BulkInsertOptionValues(ctx, optionValueParams)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		optionValueIDMap := make(map[int64]map[string]int64)
		for _, ov := range optionValues {
			if optionValueIDMap[ov.OptionID] == nil {
				optionValueIDMap[ov.OptionID] = map[string]int64{}
			}
			optionValueIDMap[ov.OptionID][ov.Name] = ov.ID
		}

		variantParams := product_db.BulkInsertVariantsParams{}

		for i, v := range req.Variants {
			variantParams.OriginPrices = append(variantParams.OriginPrices, v.OriginPrice)
			variantParams.SalePrices = append(variantParams.SalePrices, v.SalePrice)
			variantParams.Files = append(variantParams.Files, v.File)
			variantParams.Stocks = append(variantParams.Stocks, v.Stock)
			variantParams.Skus = append(variantParams.Skus, v.Sku)
			variantParams.Nos = append(variantParams.Nos, int32(i))
			variantParams.ProductIds = append(variantParams.ProductIds, productID)
		}

		variantRows, err := db.ProductQueries.BulkInsertVariants(ctx, variantParams)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		variantOptionParams := product_db.BulkInsertVariantOptionParams{}
		for i, v := range req.Variants {
			variantID := int64(variantRows[i])
			for _, opt := range v.Options {
				optionID := optionMap[opt.OptionName]
				valueID := optionValueIDMap[optionID][opt.Value]
				variantOptionParams.VariantIds = append(variantOptionParams.VariantIds, variantID)
				variantOptionParams.OptionIds = append(variantOptionParams.OptionIds, optionID)
				variantOptionParams.OptionValueIds = append(variantOptionParams.OptionValueIds, valueID)
			}
		}
		if err := db.ProductQueries.BulkInsertVariantOption(ctx, variantOptionParams); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	if len(req.CollectionIDs) > 0 {
		collectionParams := product_db.BulkInsertProductCollectionParams{}
		for _, collectionID := range req.CollectionIDs {
			collectionParams.CollectionIds = append(collectionParams.CollectionIds, collectionID)
			collectionParams.ProductIds = append(collectionParams.ProductIds, productID)
		}
		if err = db.ProductQueries.BulkInsertProductCollection(ctx, collectionParams); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	return c.SendStatus(fiber.StatusCreated)
}

// UpdateProductHandler godoc
// @Summary      Update a product
// @Description  Updates a product
// @Tags         products
// @Security BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Param        payload  body	UpdateProductRequest  true  "Update product data"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /products/{id} [put]
func UpdateProductHandler(c *fiber.Ctx) error {
	param := c.Params("id")
	productID, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid param",
		})
	}
	var req UpdateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	ctx := context.Background()
	params := product_db.UpdateProductParams{
		ID:          productID,
		Name:        req.Name,
		Slug:        req.Slug,
		OriginPrice: req.OriginPrice,
		SalePrice:   req.SalePrice,
		Stock:       pgtype.Int4{Int32: req.Stock, Valid: true},
		Sku:         pgtype.Text{String: req.SKU, Valid: true},
		Weight: pgtype.Int4{
			Int32: req.Weight,
			Valid: true,
		},
		Long:            pgtype.Int4{Int32: req.Long, Valid: true},
		Wide:            pgtype.Int4{Int32: req.Wide, Valid: true},
		High:            pgtype.Int4{Int32: req.High, Valid: true},
		CategoryID:      pgtype.Int8{Int64: req.CategoryID, Valid: req.CategoryID > 0},
		MetaTitle:       req.MetaTitle,
		MetaDescription: req.MetaDescription,
	}
	if err := db.ProductQueries.UpdateProduct(ctx, params); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := db.ProductQueries.DeleteProductFiles(ctx, productID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if len(req.Files) > 0 {
		fileParams := product_db.BulkInsertProductFilesParams{}
		for _, f := range req.Files {
			fileParams.Names = append(fileParams.Names, f.Name)
			fileParams.IsPrimaries = append(fileParams.IsPrimaries, f.IsPrimary)
			fileParams.Nos = append(fileParams.Nos, f.No)
			fileParams.ProductIds = append(fileParams.ProductIds, productID)
		}
		if err := db.ProductQueries.BulkInsertProductFiles(ctx, fileParams); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	if err := db.ProductQueries.DeleteProductTags(ctx, productID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if len(req.Tags) > 0 {
		tagParams := product_db.BulkInsertProductTagsParams{}
		for _, t := range req.Tags {
			tagParams.Names = append(tagParams.Names, t)
			tagParams.ProductIds = append(tagParams.ProductIds, productID)
		}
		if err := db.ProductQueries.BulkInsertProductTags(ctx, tagParams); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	if err := db.ProductQueries.DeleteCollectionsByProductID(ctx, productID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if len(req.CollectionIDs) > 0 {
		collectionParams := product_db.BulkInsertProductCollectionParams{}
		for _, collectionID := range req.CollectionIDs {
			collectionParams.CollectionIds = append(collectionParams.CollectionIds, collectionID)
			collectionParams.ProductIds = append(collectionParams.ProductIds, productID)
		}
		if err = db.ProductQueries.BulkInsertProductCollection(ctx, collectionParams); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	optIdsMap := map[string]int64{}
	optionValueIDMap := make(map[int64]map[string]int64)

	if len(req.Options) > 0 {
		updateOptionsParams := product_db.BulkUpdateOptionsParams{}
		updateOptionValueParams := product_db.BulkUpdateOptionValuesParams{}
		createOptionValueParams := product_db.BulkInsertOptionValuesParams{}

		productOptIds := []int64{}
		optionIDs := []int64{}
		valueIDs := []int64{}

		for _, o := range req.Options {
			if o.ID == 0 {
				continue
			}
			if len(o.Values) > 0 {
				for vIdx, v := range o.Values {
					if v.ID == 0 {
						createOptionValueParams.Names = append(createOptionValueParams.Names, v.Name)
						createOptionValueParams.OptionIds = append(createOptionValueParams.OptionIds, o.ID)
						createOptionValueParams.Nos = append(createOptionValueParams.Nos, int32(vIdx))
						continue
					}
					optionIDs = append(optionIDs, o.ID)
					valueIDs = append(valueIDs, v.ID)
					updateOptionValueParams.Ids = append(updateOptionValueParams.Ids, v.ID)
					updateOptionValueParams.Names = append(updateOptionValueParams.Names, v.Name)
				}
			}
			productOptIds = append(productOptIds, o.ID)
			updateOptionsParams.Ids = append(updateOptionsParams.Ids, o.ID)
			updateOptionsParams.Names = append(updateOptionsParams.Names, o.Name)
		}

		if len(optionIDs) > 0 {
			db.ProductQueries.DeleteOptionsNotInIDs(ctx, product_db.DeleteOptionsNotInIDsParams{ProductID: productID, Ids: productOptIds})
		}
		if len(valueIDs) > 0 {
			db.ProductQueries.DeleteOptionValuesNotInIDs(ctx, product_db.DeleteOptionValuesNotInIDsParams{
				OptionIds: optionIDs,
				ValueIds:  valueIDs,
			})
		}

		optionValuesDB, _ := db.ProductQueries.BulkInsertOptionValues(ctx, createOptionValueParams)

		for _, ov := range optionValuesDB {
			if optionValueIDMap[ov.OptionID] == nil {
				optionValueIDMap[ov.OptionID] = map[string]int64{}
			}
			optionValueIDMap[ov.OptionID][ov.Name] = ov.ID
		}

		db.ProductQueries.BulkUpdateOptionValues(ctx, updateOptionValueParams)
		db.ProductQueries.BulkUpdateOptions(ctx, updateOptionsParams)

		createOptionParams := product_db.BulkInsertOptionsParams{}
		for i, o := range req.Options {
			if o.ID != 0 {
				continue
			}
			createOptionParams.Names = append(createOptionParams.Names, o.Name)
			createOptionParams.ProductIds = append(createOptionParams.ProductIds, productID)
			createOptionParams.Nos = append(createOptionParams.Nos, int32(i))
		}
		optsDB, _ := db.ProductQueries.BulkInsertOptions(ctx, createOptionParams)

		for _, oDB := range optsDB {
			optIdsMap[oDB.Name] = oDB.ID
		}

		createOptionValueParams = product_db.BulkInsertOptionValuesParams{}

		for _, o := range req.Options {
			if o.ID != 0 {
				continue
			}
			if len(o.Values) == 0 {
				continue
			}
			for i, v := range o.Values {
				createOptionValueParams.Names = append(createOptionValueParams.Names, v.Name)
				createOptionValueParams.Nos = append(createOptionValueParams.Nos, int32(i))
				createOptionValueParams.OptionIds = append(createOptionValueParams.OptionIds, optIdsMap[o.Name])
			}
		}

		optionValuesDB, _ = db.ProductQueries.BulkInsertOptionValues(ctx, createOptionValueParams)

		for _, ov := range optionValuesDB {
			if optionValueIDMap[ov.OptionID] == nil {
				optionValueIDMap[ov.OptionID] = map[string]int64{}
			}
			optionValueIDMap[ov.OptionID][ov.Name] = ov.ID
		}
	}

	if len(req.Variants) > 0 {
		updateVariantParams := product_db.BulkUpdateVariantsParams{}
		createVariantParams := product_db.BulkInsertVariantsParams{}
		variantIDs := []int64{}
		createVariantOptionParams := product_db.BulkInsertVariantOptionParams{}
		vIdxs := []int{}
		for i, v := range req.Variants {
			if v.ID == 0 {
				vIdxs = append(vIdxs, i)
				createVariantParams.Files = append(createVariantParams.Files, v.File)
				createVariantParams.OriginPrices = append(createVariantParams.OriginPrices, v.OriginPrice)
				createVariantParams.SalePrices = append(createVariantParams.SalePrices, v.SalePrice)
				createVariantParams.Stocks = append(createVariantParams.Stocks, v.Stock)
				createVariantParams.Skus = append(createVariantParams.Skus, v.Sku)
				createVariantParams.Nos = append(createVariantParams.Nos, int32(i))
				createVariantParams.ProductIds = append(createVariantParams.ProductIds, productID)
				continue
			}
			variantIDs = append(variantIDs, v.ID)

			for _, vo := range v.Options {
				optionID := optIdsMap[vo.OptionName]
				valueID := optionValueIDMap[optionID][vo.Value]
				createVariantOptionParams.OptionIds = append(createVariantOptionParams.OptionIds, optionID)
				createVariantOptionParams.VariantIds = append(createVariantOptionParams.VariantIds, valueID)
				createVariantOptionParams.VariantIds = append(createVariantOptionParams.VariantIds, v.ID)
			}

			updateVariantParams.Ids = append(updateVariantParams.Ids, v.ID)
			updateVariantParams.OriginPrices = append(updateVariantParams.OriginPrices, v.OriginPrice)
			updateVariantParams.SalePrices = append(updateVariantParams.SalePrices, v.SalePrice)
			updateVariantParams.Files = append(updateVariantParams.Files, v.File)
			updateVariantParams.Stocks = append(updateVariantParams.Stocks, v.Stock)
			updateVariantParams.Skus = append(updateVariantParams.Skus, v.Sku)
		}
		db.ProductQueries.DeleteVariantsNotInIDsByProductID(ctx, product_db.DeleteVariantsNotInIDsByProductIDParams{
			Ids:       variantIDs,
			ProductID: productID,
		})
		db.ProductQueries.BulkUpdateVariants(ctx, updateVariantParams)
		db.ProductQueries.BulkInsertVariantOption(ctx, createVariantOptionParams)
		vdbIDs, _ := db.ProductQueries.BulkInsertVariants(ctx, createVariantParams)

		createVariantOptionParams = product_db.BulkInsertVariantOptionParams{}

		for vIdx, v := range req.Variants {
			variantID := v.ID
			if variantID == 0 {
				var tmpIdx int
				for i, v := range vIdxs {
					if v == vIdx {
						tmpIdx = i
						break
					}
				}
				variantID = int64(vdbIDs[tmpIdx])
			}

			for _, opt := range v.Options {
				if v.ID != 0 && opt.ValueID != 0 {
					continue
				}
				optionID := opt.OptionID
				valueID := opt.ValueID
				if optionID == 0 {
					optionID = optIdsMap[opt.OptionName]
				}
				if valueID == 0 {
					valueID = optionValueIDMap[optionID][opt.Value]
				}
				createVariantOptionParams.VariantIds = append(createVariantOptionParams.VariantIds, variantID)
				createVariantOptionParams.OptionIds = append(createVariantOptionParams.OptionIds, optionID)
				createVariantOptionParams.OptionValueIds = append(createVariantOptionParams.OptionValueIds, valueID)
			}
		}

		db.ProductQueries.BulkInsertVariantOption(ctx, createVariantOptionParams)
	}
	return c.SendStatus(fiber.StatusOK)
}

// DeleteProductsHandler godoc
// @Summary      Delete multiple products
// @Description  Deletes multiple products by their IDs
// @Tags         products
// @Security BearerAuth
// @Accept       json
// @Produce      json
// @Param        ids  body      DeleteProductsRequest  true  "List of product IDs"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /products [delete]
func DeleteProductsHandler(c *fiber.Ctx) error {
	var req DeleteProductsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	ctx := context.Background()
	if err := db.ProductQueries.BulkDeleteProducts(ctx, req.IDs); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusOK)
}

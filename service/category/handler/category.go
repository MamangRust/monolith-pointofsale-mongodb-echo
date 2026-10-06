package handler

import (
	pbcategories "github.com/MamangRust/monolith-point-of-sale-pb/categories"
	pbcommon "github.com/MamangRust/monolith-point-of-sale-pb/common"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"strconv"

	"github.com/MamangRust/monolith-point-of-sale-category/repository"
	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-shared/convert"
)

func mapPaginationMeta(meta *pbcommon.PaginationMeta) *pbcommon.PaginationMeta {
	if meta == nil {
		return nil
	}
	return &pbcommon.PaginationMeta{
		CurrentPage:  meta.CurrentPage,
		PageSize:     meta.PageSize,
		TotalPages:   meta.TotalPages,
		TotalRecords: meta.TotalRecords,
	}
}

func mapResponseCategory(category *salemongo.Category) *pbcategories.CategoryResponse {
	if category == nil {
		return nil
	}
	var description, slugCategory string
	if category.Description != nil {
		description = *category.Description
	}
	if category.SlugCategory != nil {
		slugCategory = *category.SlugCategory
	}
	var createdAtStr, updatedAtStr string
	if !category.CreatedAt.IsZero() {
		createdAtStr = convert.FormatTimePtr(category.CreatedAt)
	}
	if !category.UpdatedAt.IsZero() {
		updatedAtStr = convert.FormatTimePtr(category.UpdatedAt)
	}
	return &pbcategories.CategoryResponse{
		Id:           int32(category.CategoryID),
		Name:         category.Name,
		Description:  description,
		SlugCategory: slugCategory,
		CreatedAt:    createdAtStr,
		UpdatedAt:    updatedAtStr,
	}
}

func mapResponseGetCategory(category *repository.CategoryRow) *pbcategories.CategoryResponse {
	if category == nil {
		return nil
	}
	var description, slugCategory string
	if category.Description != nil {
		description = *category.Description
	}
	if category.SlugCategory != nil {
		slugCategory = *category.SlugCategory
	}
	var createdAtStr, updatedAtStr string
	if !category.CreatedAt.IsZero() {
		createdAtStr = convert.FormatTimePtr(category.CreatedAt)
	}
	if !category.UpdatedAt.IsZero() {
		updatedAtStr = convert.FormatTimePtr(category.UpdatedAt)
	}
	return &pbcategories.CategoryResponse{
		Id:           int32(category.CategoryID),
		Name:         category.Name,
		Description:  description,
		SlugCategory: slugCategory,
		CreatedAt:    createdAtStr,
		UpdatedAt:    updatedAtStr,
	}
}

func mapResponsesCategory(categories []*repository.CategoryRow) []*pbcategories.CategoryResponse {
	var mappedCategories []*pbcategories.CategoryResponse
	for _, category := range categories {
		mappedCategories = append(mappedCategories, mapResponseGetCategory(category))
	}
	return mappedCategories
}

// mapResponsesCategoryModels maps domain models (used by FindByIds).
func mapResponsesCategoryModels(categories []*salemongo.Category) []*pbcategories.CategoryResponse {
	var mappedCategories []*pbcategories.CategoryResponse
	for _, category := range categories {
		if category == nil {
			continue
		}
		mappedCategories = append(mappedCategories, mapResponseCategory(category))
	}
	return mappedCategories
}

func mapResponseCategoryDeleteAt(category *salemongo.Category) *pbcategories.CategoryResponseDeleteAt {
	if category == nil {
		return nil
	}
	var description, slugCategory string
	if category.Description != nil {
		description = *category.Description
	}
	if category.SlugCategory != nil {
		slugCategory = *category.SlugCategory
	}
	var createdAtStr, updatedAtStr string
	if !category.CreatedAt.IsZero() {
		createdAtStr = convert.FormatTimePtr(category.CreatedAt)
	}
	if !category.UpdatedAt.IsZero() {
		updatedAtStr = convert.FormatTimePtr(category.UpdatedAt)
	}
	var deletedAt *wrapperspb.StringValue
	if category.DeletedAt != nil {
		deletedAt = convert.TimeToWrappers(category.DeletedAt)
	}

	return &pbcategories.CategoryResponseDeleteAt{
		Id:           int32(category.CategoryID),
		Name:         category.Name,
		Description:  description,
		SlugCategory: slugCategory,
		CreatedAt:    createdAtStr,
		UpdatedAt:    updatedAtStr,
		DeletedAt:    deletedAt,
	}
}

func mapResponseGetCategoryActive(category *repository.CategoryActiveRow) *pbcategories.CategoryResponseDeleteAt {
	if category == nil {
		return nil
	}
	var description, slugCategory string
	if category.Description != nil {
		description = *category.Description
	}
	if category.SlugCategory != nil {
		slugCategory = *category.SlugCategory
	}
	var createdAtStr, updatedAtStr string
	if !category.CreatedAt.IsZero() {
		createdAtStr = convert.FormatTimePtr(category.CreatedAt)
	}
	if !category.UpdatedAt.IsZero() {
		updatedAtStr = convert.FormatTimePtr(category.UpdatedAt)
	}
	var deletedAt *wrapperspb.StringValue
	if category.DeletedAt != nil {
		deletedAt = wrapperspb.String(category.DeletedAt.Format("2006-01-02 15:04:05"))
	}

	return &pbcategories.CategoryResponseDeleteAt{
		Id:           int32(category.CategoryID),
		Name:         category.Name,
		Description:  description,
		SlugCategory: slugCategory,
		CreatedAt:    createdAtStr,
		UpdatedAt:    updatedAtStr,
		DeletedAt:    deletedAt,
	}
}

func mapResponsesCategoryActive(categories []*repository.CategoryActiveRow) []*pbcategories.CategoryResponseDeleteAt {
	var mappedCategories []*pbcategories.CategoryResponseDeleteAt
	for _, category := range categories {
		mappedCategories = append(mappedCategories, mapResponseGetCategoryActive(category))
	}
	return mappedCategories
}

func mapResponseGetCategoryTrashed(category *repository.CategoryTrashedRow) *pbcategories.CategoryResponseDeleteAt {
	if category == nil {
		return nil
	}
	var description, slugCategory string
	if category.Description != nil {
		description = *category.Description
	}
	if category.SlugCategory != nil {
		slugCategory = *category.SlugCategory
	}
	var createdAtStr, updatedAtStr string
	if !category.CreatedAt.IsZero() {
		createdAtStr = convert.FormatTimePtr(category.CreatedAt)
	}
	if !category.UpdatedAt.IsZero() {
		updatedAtStr = convert.FormatTimePtr(category.UpdatedAt)
	}
	var deletedAt *wrapperspb.StringValue
	if category.DeletedAt != nil {
		deletedAt = wrapperspb.String(category.DeletedAt.Format("2006-01-02 15:04:05"))
	}

	return &pbcategories.CategoryResponseDeleteAt{
		Id:           int32(category.CategoryID),
		Name:         category.Name,
		Description:  description,
		SlugCategory: slugCategory,
		CreatedAt:    createdAtStr,
		UpdatedAt:    updatedAtStr,
		DeletedAt:    deletedAt,
	}
}

func mapResponsesCategoryTrashed(categories []*repository.CategoryTrashedRow) []*pbcategories.CategoryResponseDeleteAt {
	var mappedCategories []*pbcategories.CategoryResponseDeleteAt
	for _, category := range categories {
		mappedCategories = append(mappedCategories, mapResponseGetCategoryTrashed(category))
	}
	return mappedCategories
}

func mapResponseCategoryMonthlyPrice(category *repository.CategoryMonthlyRow) *pbcategories.CategoryMonthPriceResponse {
	if category == nil {
		return nil
	}
	return &pbcategories.CategoryMonthPriceResponse{
		Month:        strconv.Itoa(int(category.Month)),
		CategoryId:   category.CategoryID,
		CategoryName: category.CategoryName,
		OrderCount:   int32(category.OrderCount),
		ItemsSold:    int32(category.ItemsSold),
		TotalRevenue: category.TotalRevenue,
	}
}

func mapResponsesCategoryMonthlyPrices(c []*repository.CategoryMonthlyRow) []*pbcategories.CategoryMonthPriceResponse {
	var categoryRecords []*pbcategories.CategoryMonthPriceResponse
	for _, category := range c {
		categoryRecords = append(categoryRecords, mapResponseCategoryMonthlyPrice(category))
	}
	return categoryRecords
}

func mapResponseCategoryMonthlyPriceById(category *repository.CategoryMonthlyByIdRow) *pbcategories.CategoryMonthPriceResponse {
	if category == nil {
		return nil
	}
	return &pbcategories.CategoryMonthPriceResponse{
		Month:        strconv.Itoa(int(category.Month)),
		CategoryId:   category.CategoryID,
		CategoryName: category.CategoryName,
		OrderCount:   int32(category.OrderCount),
		ItemsSold:    int32(category.ItemsSold),
		TotalRevenue: category.TotalRevenue,
	}
}

func mapResponsesCategoryMonthlyPricesById(c []*repository.CategoryMonthlyByIdRow) []*pbcategories.CategoryMonthPriceResponse {
	var categoryRecords []*pbcategories.CategoryMonthPriceResponse
	for _, category := range c {
		categoryRecords = append(categoryRecords, mapResponseCategoryMonthlyPriceById(category))
	}
	return categoryRecords
}

func mapResponseCategoryMonthlyPriceByMerchant(category *repository.CategoryMonthlyByMerchantRow) *pbcategories.CategoryMonthPriceResponse {
	if category == nil {
		return nil
	}
	return &pbcategories.CategoryMonthPriceResponse{
		Month:        strconv.Itoa(int(category.Month)),
		CategoryId:   category.CategoryID,
		CategoryName: category.CategoryName,
		OrderCount:   int32(category.OrderCount),
		ItemsSold:    int32(category.ItemsSold),
		TotalRevenue: category.TotalRevenue,
	}
}

func mapResponsesCategoryMonthlyPricesByMerchant(c []*repository.CategoryMonthlyByMerchantRow) []*pbcategories.CategoryMonthPriceResponse {
	var categoryRecords []*pbcategories.CategoryMonthPriceResponse
	for _, category := range c {
		categoryRecords = append(categoryRecords, mapResponseCategoryMonthlyPriceByMerchant(category))
	}
	return categoryRecords
}

func mapResponseCategoryYearlyPrice(category *repository.CategoryYearlyRow) *pbcategories.CategoryYearPriceResponse {
	if category == nil {
		return nil
	}
	return &pbcategories.CategoryYearPriceResponse{
		Year:               strconv.Itoa(int(category.Year)),
		CategoryId:         category.CategoryID,
		CategoryName:       category.CategoryName,
		OrderCount:         int32(category.OrderCount),
		ItemsSold:          int32(category.ItemsSold),
		TotalRevenue:       category.TotalRevenue,
		UniqueProductsSold: int32(category.UniqueProductsSold),
	}
}

func mapResponsesCategoryYearlyPrices(c []*repository.CategoryYearlyRow) []*pbcategories.CategoryYearPriceResponse {
	var categoryRecords []*pbcategories.CategoryYearPriceResponse
	for _, category := range c {
		categoryRecords = append(categoryRecords, mapResponseCategoryYearlyPrice(category))
	}
	return categoryRecords
}

func mapResponseCategoryYearlyPriceById(category *repository.CategoryYearlyByIdRow) *pbcategories.CategoryYearPriceResponse {
	if category == nil {
		return nil
	}
	return &pbcategories.CategoryYearPriceResponse{
		Year:               strconv.Itoa(int(category.Year)),
		CategoryId:         category.CategoryID,
		CategoryName:       category.CategoryName,
		OrderCount:         int32(category.OrderCount),
		ItemsSold:          int32(category.ItemsSold),
		TotalRevenue:       category.TotalRevenue,
		UniqueProductsSold: int32(category.UniqueProductsSold),
	}
}

func mapResponsesCategoryYearlyPricesById(c []*repository.CategoryYearlyByIdRow) []*pbcategories.CategoryYearPriceResponse {
	var categoryRecords []*pbcategories.CategoryYearPriceResponse
	for _, category := range c {
		categoryRecords = append(categoryRecords, mapResponseCategoryYearlyPriceById(category))
	}
	return categoryRecords
}

func mapResponseCategoryYearlyPriceByMerchant(category *repository.CategoryYearlyByMerchantRow) *pbcategories.CategoryYearPriceResponse {
	if category == nil {
		return nil
	}
	return &pbcategories.CategoryYearPriceResponse{
		Year:               strconv.Itoa(int(category.Year)),
		CategoryId:         category.CategoryID,
		CategoryName:       category.CategoryName,
		OrderCount:         int32(category.OrderCount),
		ItemsSold:          int32(category.ItemsSold),
		TotalRevenue:       category.TotalRevenue,
		UniqueProductsSold: int32(category.UniqueProductsSold),
	}
}

func mapResponsesCategoryYearlyPricesByMerchant(c []*repository.CategoryYearlyByMerchantRow) []*pbcategories.CategoryYearPriceResponse {
	var categoryRecords []*pbcategories.CategoryYearPriceResponse
	for _, category := range c {
		categoryRecords = append(categoryRecords, mapResponseCategoryYearlyPriceByMerchant(category))
	}
	return categoryRecords
}

func mapResponseCashierMonthlyTotalPrice(c *repository.CategoryMonthlyTotalPriceRow) *pbcategories.CategoriesMonthlyTotalPriceResponse {
	if c == nil {
		return nil
	}
	return &pbcategories.CategoriesMonthlyTotalPriceResponse{
		Year:         strconv.Itoa(int(c.Year)),
		Month:        strconv.Itoa(int(c.Month)),
		TotalRevenue: c.TotalPrice,
	}
}

func mapResponseCategoryMonthlyTotalPrices(c []*repository.CategoryMonthlyTotalPriceRow) []*pbcategories.CategoriesMonthlyTotalPriceResponse {
	var CategoryRecords []*pbcategories.CategoriesMonthlyTotalPriceResponse
	for _, Category := range c {
		CategoryRecords = append(CategoryRecords, mapResponseCashierMonthlyTotalPrice(Category))
	}
	return CategoryRecords
}

func mapResponseCashierMonthlyTotalPriceById(c *repository.CategoryMonthlyTotalPriceByIdRow) *pbcategories.CategoriesMonthlyTotalPriceResponse {
	if c == nil {
		return nil
	}
	return &pbcategories.CategoriesMonthlyTotalPriceResponse{
		Year:         strconv.Itoa(int(c.Year)),
		Month:        strconv.Itoa(int(c.Month)),
		TotalRevenue: c.TotalPrice,
	}
}

func mapResponseCategoryMonthlyTotalPricesById(c []*repository.CategoryMonthlyTotalPriceByIdRow) []*pbcategories.CategoriesMonthlyTotalPriceResponse {
	var CategoryRecords []*pbcategories.CategoriesMonthlyTotalPriceResponse
	for _, Category := range c {
		CategoryRecords = append(CategoryRecords, mapResponseCashierMonthlyTotalPriceById(Category))
	}
	return CategoryRecords
}

func mapResponseCashierMonthlyTotalPriceByMerchant(c *repository.CategoryMonthlyTotalPriceByMerchantRow) *pbcategories.CategoriesMonthlyTotalPriceResponse {
	if c == nil {
		return nil
	}
	return &pbcategories.CategoriesMonthlyTotalPriceResponse{
		Year:         strconv.Itoa(int(c.Year)),
		Month:        strconv.Itoa(int(c.Month)),
		TotalRevenue: c.TotalPrice,
	}
}

func mapResponseCategoryMonthlyTotalPricesByMerchant(c []*repository.CategoryMonthlyTotalPriceByMerchantRow) []*pbcategories.CategoriesMonthlyTotalPriceResponse {
	var CategoryRecords []*pbcategories.CategoriesMonthlyTotalPriceResponse
	for _, Category := range c {
		CategoryRecords = append(CategoryRecords, mapResponseCashierMonthlyTotalPriceByMerchant(Category))
	}
	return CategoryRecords
}

func mapResponseCategoryYearlyTotalSale(c *repository.CategoryYearlyTotalPriceRow) *pbcategories.CategoriesYearlyTotalPriceResponse {
	if c == nil {
		return nil
	}
	return &pbcategories.CategoriesYearlyTotalPriceResponse{
		Year:         strconv.Itoa(int(c.Year)),
		TotalRevenue: c.TotalPrice,
	}
}

func mapResponseCategoryYearlyTotalPrices(c []*repository.CategoryYearlyTotalPriceRow) []*pbcategories.CategoriesYearlyTotalPriceResponse {
	var CategoryRecords []*pbcategories.CategoriesYearlyTotalPriceResponse
	for _, Category := range c {
		CategoryRecords = append(CategoryRecords, mapResponseCategoryYearlyTotalSale(Category))
	}
	return CategoryRecords
}

func mapResponseCategoryYearlyTotalSaleById(c *repository.CategoryYearlyTotalPriceByIdRow) *pbcategories.CategoriesYearlyTotalPriceResponse {
	if c == nil {
		return nil
	}
	return &pbcategories.CategoriesYearlyTotalPriceResponse{
		Year:         strconv.Itoa(int(c.Year)),
		TotalRevenue: c.TotalPrice,
	}
}

func mapResponseCategoryYearlyTotalPricesById(c []*repository.CategoryYearlyTotalPriceByIdRow) []*pbcategories.CategoriesYearlyTotalPriceResponse {
	var CategoryRecords []*pbcategories.CategoriesYearlyTotalPriceResponse
	for _, Category := range c {
		CategoryRecords = append(CategoryRecords, mapResponseCategoryYearlyTotalSaleById(Category))
	}
	return CategoryRecords
}

func mapResponseCategoryYearlyTotalSaleByMerchant(c *repository.CategoryYearlyTotalPriceByMerchantRow) *pbcategories.CategoriesYearlyTotalPriceResponse {
	if c == nil {
		return nil
	}
	return &pbcategories.CategoriesYearlyTotalPriceResponse{
		Year:         strconv.Itoa(int(c.Year)),
		TotalRevenue: c.TotalPrice,
	}
}

func mapResponseCategoryYearlyTotalPricesByMerchant(c []*repository.CategoryYearlyTotalPriceByMerchantRow) []*pbcategories.CategoriesYearlyTotalPriceResponse {
	var CategoryRecords []*pbcategories.CategoriesYearlyTotalPriceResponse
	for _, Category := range c {
		CategoryRecords = append(CategoryRecords, mapResponseCategoryYearlyTotalSaleByMerchant(Category))
	}
	return CategoryRecords
}

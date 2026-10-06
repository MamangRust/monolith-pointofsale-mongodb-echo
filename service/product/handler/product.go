package handler

import (
	pbcommon "github.com/MamangRust/monolith-point-of-sale-pb/common"
	pbproducts "github.com/MamangRust/monolith-point-of-sale-pb/products"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"time"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-product/repository"
	"github.com/MamangRust/monolith-point-of-sale-shared/convert"
)

func parseStrPointer(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func parseInt32Pointer(i *int32) int32 {
	if i != nil {
		return *i
	}
	return 0
}

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

func formatTimePtr(t *time.Time) string {
	if t != nil {
		return t.Format("2006-01-02 15:04:05")
	}
	return ""
}

func formatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func formatOptionalTime(t *time.Time) *wrapperspb.StringValue {
	if t == nil {
		return nil
	}
	return wrapperspb.String(t.Format("2006-01-02 15:04:05"))
}

func mapResponseProduct(product *salemongo.Product) *pbproducts.ProductResponse {
	if product == nil {
		return nil
	}
	return &pbproducts.ProductResponse{
		Id:           int32(product.ProductID),
		MerchantId:   int32(product.MerchantID),
		CategoryId:   int32(product.CategoryID),
		Name:         product.Name,
		Description:  parseStrPointer(product.Description),
		Price:        int32(product.Price),
		CountInStock: int32(product.CountInStock),
		Brand:        parseStrPointer(product.Brand),
		Weight:       parseInt32Pointer(product.Weight),
		SlugProduct:  parseStrPointer(product.SlugProduct),
		ImageProduct: parseStrPointer(product.ImageProduct),
		Barcode:      parseStrPointer(product.Barcode),
		CreatedAt:    formatTimePtr(product.CreatedAt),
		UpdatedAt:    formatTimePtr(product.UpdatedAt),
	}
}

func mapResponsesProduct(products []*repository.ProductRow) []*pbproducts.ProductResponse {
	var mappedProducts []*pbproducts.ProductResponse
	for _, p := range products {
		if p == nil {
			continue
		}
		mappedProducts = append(mappedProducts, &pbproducts.ProductResponse{
			Id:           int32(p.ProductID),
			MerchantId:   int32(p.MerchantID),
			CategoryId:   int32(p.CategoryID),
			Name:         p.Name,
			Description:  parseStrPointer(p.Description),
			Price:        int32(p.Price),
			CountInStock: int32(p.CountInStock),
			Brand:        parseStrPointer(p.Brand),
			Weight:       parseInt32Pointer(p.Weight),
			SlugProduct:  parseStrPointer(p.SlugProduct),
			ImageProduct: parseStrPointer(p.ImageProduct),
			Barcode:      parseStrPointer(p.Barcode),
			CreatedAt:    formatTimePtr(p.CreatedAt),
			UpdatedAt:    formatTimePtr(p.UpdatedAt),
		})
	}
	return mappedProducts
}

func mapResponsesProductByMerchant(products []*repository.ProductMerchantRow, merchantId int32) []*pbproducts.ProductResponse {
	var mappedProducts []*pbproducts.ProductResponse
	for _, p := range products {
		if p == nil {
			continue
		}
		mappedProducts = append(mappedProducts, &pbproducts.ProductResponse{
			Id:           int32(p.ProductID),
			MerchantId:   merchantId,
			Name:         p.Name,
			Description:  parseStrPointer(p.Description),
			Price:        int32(p.Price),
			CountInStock: int32(p.CountInStock),
			Brand:        parseStrPointer(p.Brand),
			ImageProduct: parseStrPointer(p.ImageProduct),
			CreatedAt:    formatTimePtr(p.CreatedAt),
		})
	}
	return mappedProducts
}

func mapResponsesProductByCategory(products []*repository.ProductCategoryNameRow) []*pbproducts.ProductResponse {
	var mappedProducts []*pbproducts.ProductResponse
	for _, p := range products {
		if p == nil {
			continue
		}
		mappedProducts = append(mappedProducts, &pbproducts.ProductResponse{
			Id:           int32(p.ProductID),
			MerchantId:   int32(p.MerchantID),
			CategoryId:   int32(p.CategoryID),
			Name:         p.Name,
			Description:  parseStrPointer(p.Description),
			Price:        int32(p.Price),
			CountInStock: int32(p.CountInStock),
			Brand:        parseStrPointer(p.Brand),
			Weight:       parseInt32Pointer(p.Weight),
			SlugProduct:  parseStrPointer(p.SlugProduct),
			ImageProduct: parseStrPointer(p.ImageProduct),
			Barcode:      parseStrPointer(p.Barcode),
			CreatedAt:    formatTimePtr(p.CreatedAt),
			UpdatedAt:    formatTimePtr(p.UpdatedAt),
		})
	}
	return mappedProducts
}

func mapResponseProductDeleteAt(product *salemongo.Product) *pbproducts.ProductResponseDeleteAt {
	if product == nil {
		return nil
	}
	var deletedAt *wrapperspb.StringValue
	if product.DeletedAt != nil {
		deletedAt = convert.TimeToWrappers(product.DeletedAt)
	}

	return &pbproducts.ProductResponseDeleteAt{
		Id:           int32(product.ProductID),
		MerchantId:   int32(product.MerchantID),
		CategoryId:   int32(product.CategoryID),
		Name:         product.Name,
		Description:  parseStrPointer(product.Description),
		Price:        int32(product.Price),
		CountInStock: int32(product.CountInStock),
		Brand:        parseStrPointer(product.Brand),
		Weight:       parseInt32Pointer(product.Weight),
		SlugProduct:  parseStrPointer(product.SlugProduct),
		ImageProduct: parseStrPointer(product.ImageProduct),
		Barcode:      parseStrPointer(product.Barcode),
		CreatedAt:    formatTimePtr(product.CreatedAt),
		UpdatedAt:    formatTimePtr(product.UpdatedAt),
		DeletedAt:    deletedAt,
	}
}

func mapResponsesProductActive(products []*repository.ProductActiveRow) []*pbproducts.ProductResponseDeleteAt {
	var mappedProducts []*pbproducts.ProductResponseDeleteAt
	for _, p := range products {
		if p == nil {
			continue
		}
		mappedProducts = append(mappedProducts, &pbproducts.ProductResponseDeleteAt{
			Id:           int32(p.ProductID),
			MerchantId:   int32(p.MerchantID),
			CategoryId:   int32(p.CategoryID),
			Name:         p.Name,
			Description:  parseStrPointer(p.Description),
			Price:        int32(p.Price),
			CountInStock: int32(p.CountInStock),
			Brand:        parseStrPointer(p.Brand),
			Weight:       parseInt32Pointer(p.Weight),
			SlugProduct:  parseStrPointer(p.SlugProduct),
			ImageProduct: parseStrPointer(p.ImageProduct),
			Barcode:      parseStrPointer(p.Barcode),
			CreatedAt:    formatTimePtr(p.CreatedAt),
			UpdatedAt:    formatTimePtr(p.UpdatedAt),
			DeletedAt:    formatOptionalTime(p.DeletedAt),
		})
	}
	return mappedProducts
}

func mapResponsesProductTrashed(products []*repository.ProductTrashedRow) []*pbproducts.ProductResponseDeleteAt {
	var mappedProducts []*pbproducts.ProductResponseDeleteAt
	for _, p := range products {
		if p == nil {
			continue
		}
		mappedProducts = append(mappedProducts, &pbproducts.ProductResponseDeleteAt{
			Id:           int32(p.ProductID),
			MerchantId:   int32(p.MerchantID),
			CategoryId:   int32(p.CategoryID),
			Name:         p.Name,
			Description:  parseStrPointer(p.Description),
			Price:        int32(p.Price),
			CountInStock: int32(p.CountInStock),
			Brand:        parseStrPointer(p.Brand),
			Weight:       parseInt32Pointer(p.Weight),
			SlugProduct:  parseStrPointer(p.SlugProduct),
			ImageProduct: parseStrPointer(p.ImageProduct),
			Barcode:      parseStrPointer(p.Barcode),
			CreatedAt:    formatTimePtr(p.CreatedAt),
			UpdatedAt:    formatTimePtr(p.UpdatedAt),
			DeletedAt:    formatOptionalTime(p.DeletedAt),
		})
	}
	return mappedProducts
}

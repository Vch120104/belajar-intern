package masteritemservice

import (
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type DiscountPercentService interface {
	GetAllDiscountPercent(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int)
	// SaveDiscountPercent(req masteritempayloads.DiscountPercentResponse) (bool, error)
	// GetDiscountPercentById(Id int) (masteritempayloads.DiscountPercentResponse, error)
	// ChangeStatusDiscountPercent(Id int) (bool, error)
}

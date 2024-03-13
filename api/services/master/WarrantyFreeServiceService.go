package masterservice

import (
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type WarrantyFreeServiceService interface {
	GetAllWarrantyFreeService(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int)
	GetWarrantyFreeServiceById(Id int) map[string]interface{}
	SaveWarrantyFreeService(req masterpayloads.WarrantyFreeServiceRequest) bool
}

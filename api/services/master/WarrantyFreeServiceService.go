package masterservice

import (
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type WarrantyFreeServiceService interface {
	GetAllWarrantyFreeService(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetWarrantyFreeServiceById(Id int) (map[string]interface{}, *exceptions.BaseErrorResponse)
	SaveWarrantyFreeService(req masterpayloads.WarrantyFreeServiceRequest) (bool, *exceptions.BaseErrorResponse)
	ChangeStatusWarrantyFreeService(Id int) (bool, *exceptions.BaseErrorResponse)
}

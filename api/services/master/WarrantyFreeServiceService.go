package masterservice

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type WarrantyFreeServiceService interface {
	GetAllWarrantyFreeService(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse)
	GetWarrantyFreeServiceById(Id int) (map[string]interface{}, *exceptionsss_test.BaseErrorResponse)
	SaveWarrantyFreeService(req masterpayloads.WarrantyFreeServiceRequest) (bool, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusWarrantyFreeService(Id int) (bool, *exceptionsss_test.BaseErrorResponse)
}

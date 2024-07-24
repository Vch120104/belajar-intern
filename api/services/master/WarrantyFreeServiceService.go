package masterservice

import (
	masterentities "after-sales/api/entities/master"
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type WarrantyFreeServiceService interface {
	GetAllWarrantyFreeService(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetWarrantyFreeServiceById(Id int) (map[string]interface{}, *exceptions.BaseErrorResponse)
	SaveWarrantyFreeService(req masterpayloads.WarrantyFreeServiceRequest) (masterentities.WarrantyFreeService, *exceptions.BaseErrorResponse)
	ChangeStatusWarrantyFreeService(Id int) (masterpayloads.WarrantyFreeServicePatchResponse, *exceptions.BaseErrorResponse)
	UpdateWarrantyFreeService(req masterentities.WarrantyFreeService, id int)(masterentities.WarrantyFreeService,*exceptions.BaseErrorResponse)
}

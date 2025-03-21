package masterservice

import (
	masterentities "after-sales/api/entities/master"
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type PackageMasterService interface {
	GetAllPackageMaster([]utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetAllPackageMasterDetail(pagination.Pagination, int) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetByIdPackageMaster(int) (map[string]interface{}, *exceptions.BaseErrorResponse)
	GetByIdPackageMasterDetail(int) (map[string]interface{}, *exceptions.BaseErrorResponse)
	GetByCodePackageMaster(string) (masterentities.PackageMaster, *exceptions.BaseErrorResponse)
	PostPackageMaster(masterpayloads.PackageMasterResponse) (masterentities.PackageMaster, *exceptions.BaseErrorResponse)
	PostPackageMasterDetail(masterpayloads.PackageMasterDetail, int) (masterentities.PackageMasterDetail, *exceptions.BaseErrorResponse)
	ChangeStatusItemPackage(int) (masterentities.PackageMaster, *exceptions.BaseErrorResponse)
	ActivateMultiIdPackageMasterDetail(string) (bool, *exceptions.BaseErrorResponse)
	DeactivateMultiIdPackageMasterDetail(string) (bool, *exceptions.BaseErrorResponse)
	CopyToOtherModel(int, string, int) (int, *exceptions.BaseErrorResponse)
}

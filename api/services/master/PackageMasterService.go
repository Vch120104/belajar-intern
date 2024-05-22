package masterservice

import (
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type PackageMasterService interface {
	GetAllPackageMaster([]utils.FilterCondition, pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetAllPackageMasterDetail(pagination.Pagination, int) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetByIdPackageMaster(int) (map[string]interface{}, *exceptions.BaseErrorResponse)
	GetByIdPackageMasterDetail(int, int, int) (map[string]interface{}, *exceptions.BaseErrorResponse)
	PostPackageMaster(masterpayloads.PackageMasterResponse) (bool, *exceptions.BaseErrorResponse)
	PostPackageMasterDetailBodyshop(masterpayloads.PackageMasterDetailOperationBodyshop, int) (bool, *exceptions.BaseErrorResponse)
	PostPackageMasterDetailWorkshop(masterpayloads.PackageMasterDetailWorkshop) (bool, *exceptions.BaseErrorResponse)
	ChangeStatusItemPackage(int) (bool, *exceptions.BaseErrorResponse)
	ActivateMultiIdPackageMasterDetail(string, int) (bool, *exceptions.BaseErrorResponse)
	DeactivateMultiIdPackageMasterDetail(string, int) (bool, *exceptions.BaseErrorResponse)
	CopyToOtherModel(int, string, int) (bool, *exceptions.BaseErrorResponse)
}

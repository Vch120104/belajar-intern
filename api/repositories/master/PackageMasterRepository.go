package masterrepository

import (
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type PackageMasterRepository interface {
	GetAllPackageMaster(*gorm.DB, []utils.FilterCondition, pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetAllPackageMasterDetail(*gorm.DB, int, pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetByIdPackageMaster(*gorm.DB, int) (map[string]interface{}, *exceptions.BaseErrorResponse)
	GetByIdPackageMasterDetail(*gorm.DB, int, int, int) (map[string]interface{}, *exceptions.BaseErrorResponse)
	PostpackageMaster(*gorm.DB, masterpayloads.PackageMasterResponse) (bool, *exceptions.BaseErrorResponse)
	PostPackageMasterDetailBodyshop(*gorm.DB, masterpayloads.PackageMasterDetailOperationBodyshop, int) (bool, *exceptions.BaseErrorResponse)
	PostPackageMasterDetailWorkshop(*gorm.DB, masterpayloads.PackageMasterDetailWorkshop) (bool, *exceptions.BaseErrorResponse)
	ChangeStatusItemPackage(*gorm.DB, int) (bool, *exceptions.BaseErrorResponse)
	DeactivateMultiIdPackageMasterDetail(*gorm.DB, string, int) (bool, *exceptions.BaseErrorResponse)
	ActivateMultiIdPackageMasterDetail(*gorm.DB, string, int) (bool, *exceptions.BaseErrorResponse)
	CopyToOtherModel(*gorm.DB, int, string, int) (bool, *exceptions.BaseErrorResponse)
}

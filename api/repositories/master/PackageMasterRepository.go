package masterrepository

import (
	masterentities "after-sales/api/entities/master"
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type PackageMasterRepository interface {
	GetAllPackageMaster(*gorm.DB, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetAllPackageMasterDetail(*gorm.DB, int, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetByIdPackageMaster(*gorm.DB, int) (map[string]interface{}, *exceptions.BaseErrorResponse)
	GetByIdPackageMasterDetail(*gorm.DB, int) (map[string]interface{}, *exceptions.BaseErrorResponse)
	GetByCodePackageMaster(*gorm.DB, string) (masterentities.PackageMaster, *exceptions.BaseErrorResponse)
	PostpackageMaster(*gorm.DB, masterpayloads.PackageMasterResponse) (masterentities.PackageMaster, *exceptions.BaseErrorResponse)
	PostPackageMasterDetail(tx *gorm.DB, req masterpayloads.PackageMasterDetail, id int) (masterentities.PackageMasterDetail, *exceptions.BaseErrorResponse)
	ChangeStatusItemPackage(*gorm.DB, int) (masterentities.PackageMaster, *exceptions.BaseErrorResponse)
	DeactivateMultiIdPackageMasterDetail(*gorm.DB, string) (bool, *exceptions.BaseErrorResponse)
	ActivateMultiIdPackageMasterDetail(*gorm.DB, string) (bool, *exceptions.BaseErrorResponse)
	CopyToOtherModel(*gorm.DB, int, string, int) (int, *exceptions.BaseErrorResponse)
}

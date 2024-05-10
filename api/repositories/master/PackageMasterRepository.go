package masterrepository

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type PackageMasterRepository interface {
	GetAllPackageMaster(*gorm.DB, []utils.FilterCondition, pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse)
	GetAllPackageMasterDetail(*gorm.DB, int, pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse)
	GetByIdPackageMaster(*gorm.DB, int) (map[string]interface{}, *exceptionsss_test.BaseErrorResponse)
	GetByIdPackageMasterDetail(*gorm.DB, int, int, int) (map[string]interface{}, *exceptionsss_test.BaseErrorResponse)
	PostpackageMaster(*gorm.DB, masterpayloads.PackageMasterResponse) (bool, *exceptionsss_test.BaseErrorResponse)
	PostPackageMasterDetailBodyshop(*gorm.DB, masterpayloads.PackageMasterDetailOperationBodyshop, int)(bool,*exceptionsss_test.BaseErrorResponse)
	PostPackageMasterDetailWorkshop(*gorm.DB, masterpayloads.PackageMasterDetailWorkshop)(bool,*exceptionsss_test.BaseErrorResponse)
	ChangeStatusItemPackage(*gorm.DB, int) (bool, *exceptionsss_test.BaseErrorResponse)
	DeactivateMultiIdPackageMasterDetail(*gorm.DB, string, int) (bool, *exceptionsss_test.BaseErrorResponse)
	ActivateMultiIdPackageMasterDetail(*gorm.DB, string, int) (bool, *exceptionsss_test.BaseErrorResponse)
	CopyToOtherModel(*gorm.DB, int, string, int) (bool, *exceptionsss_test.BaseErrorResponse)
}
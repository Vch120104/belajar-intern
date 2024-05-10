package masterservice

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type PackageMasterService interface {
	GetAllPackageMaster([]utils.FilterCondition, pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse)
	GetAllPackageMasterDetail( pagination.Pagination, int)([]map[string]interface{},int,int,*exceptionsss_test.BaseErrorResponse)
	GetByIdPackageMaster(int)(map[string]interface{}, *exceptionsss_test.BaseErrorResponse)
	GetByIdPackageMasterDetail(int,int,int)(map[string]interface{},*exceptionsss_test.BaseErrorResponse)
	PostPackageMaster( masterpayloads.PackageMasterResponse)(bool,*exceptionsss_test.BaseErrorResponse)
	PostPackageMasterDetailBodyshop(masterpayloads.PackageMasterDetailOperationBodyshop, int)(bool,*exceptionsss_test.BaseErrorResponse)
	PostPackageMasterDetailWorkshop(masterpayloads.PackageMasterDetailWorkshop)(bool,*exceptionsss_test.BaseErrorResponse)
	ChangeStatusItemPackage(int)(bool,*exceptionsss_test.BaseErrorResponse)
	ActivateMultiIdPackageMasterDetail(string,int)(bool,*exceptionsss_test.BaseErrorResponse)
	DeactivateMultiIdPackageMasterDetail(string,int)(bool,*exceptionsss_test.BaseErrorResponse)
	CopyToOtherModel(int, string, int)(bool,*exceptionsss_test.BaseErrorResponse)
}
package masterrepository

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type AgreementRepository interface {
	GetAgreementById(*gorm.DB, int) (masterpayloads.AgreementRequest, *exceptionsss_test.BaseErrorResponse)
	SaveAgreement(*gorm.DB, masterpayloads.AgreementResponse) (bool, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusAgreement(*gorm.DB, int) (bool, *exceptionsss_test.BaseErrorResponse)
	GetAllAgreement(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse)
}

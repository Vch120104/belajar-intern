package masterservice

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type AgreementService interface {
	GetAgreementById(int) (masterpayloads.AgreementResponse, *exceptionsss_test.BaseErrorResponse)
	SaveAgreement(masterpayloads.AgreementResponse) (bool, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusAgreement(int) (bool, *exceptionsss_test.BaseErrorResponse)
	GetAllAgreement(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse)
}

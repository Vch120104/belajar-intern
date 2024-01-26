package masteritemrepository

import (
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type MarkupRateRepository interface {
	WithTrx(trxHandle *gorm.DB) MarkupRateRepository
	GetMarkupRateById(Id int) (masteritempayloads.MarkupRateResponse, error)
	SaveMarkupRate(request masteritempayloads.MarkupRateRequest) (bool, error)
	GetAllMarkupRate(filterCondition []utils.FilterCondition) ([]map[string]interface{}, error)
	ChangeStatusMarkupRate(Id int) (bool, error)
}

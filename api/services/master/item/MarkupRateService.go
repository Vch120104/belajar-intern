package masteritemservice

import (
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type MarkupRateService interface {
	WithTrx(trxHandle *gorm.DB) MarkupRateService
	GetMarkupRateById(id int) (masteritempayloads.MarkupRateResponse, error)
	SaveMarkupRate(req masteritempayloads.MarkupRateRequest) (bool, error)
	GetAllMarkupRate(filterCondition []utils.FilterCondition) ([]map[string]interface{}, error)
	ChangeStatusMarkupRate(Id int) (bool, error)
}

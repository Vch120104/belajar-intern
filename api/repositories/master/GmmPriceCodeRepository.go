package masterrepository

import (
	masterentities "after-sales/api/entities/master"
	"after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"

	"gorm.io/gorm"
)

type GmmPriceCodeRepository interface {
	GetAllGmmPriceCode(tx *gorm.DB) ([]masterpayloads.GmmPriceCodeResponse, *exceptions.BaseErrorResponse)
	GetGmmPriceCodeById(tx *gorm.DB, id int) (masterpayloads.GmmPriceCodeResponse, *exceptions.BaseErrorResponse)
	GetGmmPriceCodeDropdown(tx *gorm.DB) ([]masterpayloads.GmmPriceCodeDropdownResponse, *exceptions.BaseErrorResponse)
	SaveGmmPriceCode(tx *gorm.DB, req masterpayloads.GmmPriceCodeSaveRequest) (masterentities.GmmPriceCode, *exceptions.BaseErrorResponse)
	UpdateGmmPriceCode(tx *gorm.DB, id int, req masterpayloads.GmmPriceCodeUpdateRequest) (masterentities.GmmPriceCode, *exceptions.BaseErrorResponse)
	ChangeStatusGmmPriceCode(tx *gorm.DB, id int) (masterentities.GmmPriceCode, *exceptions.BaseErrorResponse)
	DeleteGmmPriceCode(tx *gorm.DB, id int) (bool, *exceptions.BaseErrorResponse)
}

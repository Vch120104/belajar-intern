package masterservice

import (
	masterentities "after-sales/api/entities/master"
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
)

type GmmPriceCodeService interface {
	GetAllGmmPriceCode() ([]masterpayloads.GmmPriceCodeResponse, *exceptions.BaseErrorResponse)
	GetGmmPriceCodeById(id int) (masterpayloads.GmmPriceCodeResponse, *exceptions.BaseErrorResponse)
	GetGmmPriceCodeDropdown() ([]masterpayloads.GmmPriceCodeDropdownResponse, *exceptions.BaseErrorResponse)
	SaveGmmPriceCode(req masterpayloads.GmmPriceCodeSaveRequest) (masterentities.GmmPriceCode, *exceptions.BaseErrorResponse)
	UpdateGmmPriceCode(id int, req masterpayloads.GmmPriceCodeUpdateRequest) (masterentities.GmmPriceCode, *exceptions.BaseErrorResponse)
	ChangeStatusGmmPriceCode(id int) (masterentities.GmmPriceCode, *exceptions.BaseErrorResponse)
	DeleteGmmPriceCode(id int) (bool, *exceptions.BaseErrorResponse)
}

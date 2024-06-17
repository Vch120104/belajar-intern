package masterservice

import (
	masterentities "after-sales/api/entities/master"
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type AgreementService interface {
	GetAgreementById(int) (masterpayloads.AgreementRequest, *exceptions.BaseErrorResponse)
	SaveAgreement(masterpayloads.AgreementRequest) (bool, *exceptions.BaseErrorResponse)
	ChangeStatusAgreement(int) (masterentities.Agreement, *exceptions.BaseErrorResponse)
	GetAllAgreement(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	AddDiscountGroup(int, masterpayloads.DiscountGroupRequest) *exceptions.BaseErrorResponse
	UpdateDiscountGroup(int, int, masterpayloads.DiscountGroupRequest) *exceptions.BaseErrorResponse
	DeleteDiscountGroup(int, int) *exceptions.BaseErrorResponse
	AddItemDiscount(int, masterpayloads.ItemDiscountRequest) *exceptions.BaseErrorResponse
	UpdateItemDiscount(int, int, masterpayloads.ItemDiscountRequest) *exceptions.BaseErrorResponse
	DeleteItemDiscount(int, int) *exceptions.BaseErrorResponse
	AddDiscountValue(int, masterpayloads.DiscountValueRequest) *exceptions.BaseErrorResponse
	UpdateDiscountValue(int, int, masterpayloads.DiscountValueRequest) *exceptions.BaseErrorResponse
	DeleteDiscountValue(int, int) *exceptions.BaseErrorResponse
	GetAllDiscountGroup(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetAllItemDiscount(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetAllDiscountValue(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetDiscountGroupAgreementById(int, int) (masterpayloads.DiscountGroupRequest, *exceptions.BaseErrorResponse)
	GetDiscountItemAgreementById(int, int) (masterpayloads.ItemDiscountRequest, *exceptions.BaseErrorResponse)
	GetDiscountValueAgreementById(int, int) (masterpayloads.DiscountValueRequest, *exceptions.BaseErrorResponse)
}

package masterservice

import (
	masterentities "after-sales/api/entities/master"
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type AgreementService interface {
	GetAgreementById(int) (masterpayloads.AgreementResponse, *exceptions.BaseErrorResponse)
	GetAgreementByCode(string) (masterpayloads.AgreementResponse, *exceptions.BaseErrorResponse)
	SaveAgreement(masterpayloads.AgreementRequest) (masterentities.Agreement, *exceptions.BaseErrorResponse)
	UpdateAgreement(int, masterpayloads.AgreementRequest) (masterentities.Agreement, *exceptions.BaseErrorResponse)
	ChangeStatusAgreement(int) (masterentities.Agreement, *exceptions.BaseErrorResponse)
	GetAllAgreement(internalFilter []utils.FilterCondition, externalFilter []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	AddDiscountGroup(int, masterpayloads.DiscountGroupRequest) (masterentities.AgreementDiscountGroupDetail, *exceptions.BaseErrorResponse)
	UpdateDiscountGroup(int, int, masterpayloads.DiscountGroupRequest) (masterentities.AgreementDiscountGroupDetail, *exceptions.BaseErrorResponse)
	DeleteDiscountGroup(int, int) *exceptions.BaseErrorResponse
	AddItemDiscount(int, masterpayloads.ItemDiscountRequest) (masterentities.AgreementItemDetail, *exceptions.BaseErrorResponse)
	UpdateItemDiscount(int, int, masterpayloads.ItemDiscountRequest) (masterentities.AgreementItemDetail, *exceptions.BaseErrorResponse)
	DeleteItemDiscount(int, int) *exceptions.BaseErrorResponse
	AddDiscountValue(int, masterpayloads.DiscountValueRequest) (masterentities.AgreementDiscount, *exceptions.BaseErrorResponse)
	UpdateDiscountValue(int, int, masterpayloads.DiscountValueRequest) (masterentities.AgreementDiscount, *exceptions.BaseErrorResponse)
	DeleteDiscountValue(int, int) *exceptions.BaseErrorResponse
	GetAllDiscountGroup(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetAllItemDiscount(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetAllDiscountValue(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetDiscountGroupAgreementById(int, int) (masterpayloads.DiscountGroupRequest, *exceptions.BaseErrorResponse)
	GetDiscountItemAgreementById(int, int) (masterpayloads.ItemDiscountRequest, *exceptions.BaseErrorResponse)
	GetDiscountValueAgreementById(int, int) (masterpayloads.DiscountValueRequest, *exceptions.BaseErrorResponse)
	GetDiscountGroupAgreementByHeaderId(id int, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetDiscountItemAgreementByHeaderId(id int, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetDiscountValueAgreementByHeaderId(id int, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	DeleteMultiIdDiscountGroup(agreementID int, intIds []int) (bool, *exceptions.BaseErrorResponse)
	DeleteMultiIdItemDiscount(agreementID int, intIds []int) (bool, *exceptions.BaseErrorResponse)
	DeleteMultiIdDiscountValue(agreementID int, intIds []int) (bool, *exceptions.BaseErrorResponse)
}

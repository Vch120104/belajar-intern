package masterrepository

import (
	masterentities "after-sales/api/entities/master"
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type AgreementRepository interface {
	GetAgreementById(*gorm.DB, int) (masterpayloads.AgreementRequest, *exceptions.BaseErrorResponse)
	GetAgreementByCode(*gorm.DB, string) (masterpayloads.AgreementResponse, *exceptions.BaseErrorResponse)
	SaveAgreement(*gorm.DB, masterpayloads.AgreementRequest) (masterentities.Agreement, *exceptions.BaseErrorResponse)
	UpdateAgreement(*gorm.DB, int, masterpayloads.AgreementRequest) (masterentities.Agreement, *exceptions.BaseErrorResponse)
	ChangeStatusAgreement(*gorm.DB, int) (masterentities.Agreement, *exceptions.BaseErrorResponse)
	GetAllAgreement(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	AddDiscountGroup(*gorm.DB, int, masterpayloads.DiscountGroupRequest) (masterentities.AgreementDiscountGroupDetail, *exceptions.BaseErrorResponse)
	UpdateDiscountGroup(*gorm.DB, int, int, masterpayloads.DiscountGroupRequest) (masterentities.AgreementDiscountGroupDetail, *exceptions.BaseErrorResponse)
	DeleteDiscountGroup(*gorm.DB, int, int) *exceptions.BaseErrorResponse
	AddItemDiscount(*gorm.DB, int, masterpayloads.ItemDiscountRequest) (masterentities.AgreementItemDetail, *exceptions.BaseErrorResponse)
	UpdateItemDiscount(*gorm.DB, int, int, masterpayloads.ItemDiscountRequest) (masterentities.AgreementItemDetail, *exceptions.BaseErrorResponse)
	DeleteItemDiscount(*gorm.DB, int, int) *exceptions.BaseErrorResponse
	AddDiscountValue(*gorm.DB, int, masterpayloads.DiscountValueRequest) (masterentities.AgreementDiscount, *exceptions.BaseErrorResponse)
	UpdateDiscountValue(*gorm.DB, int, int, masterpayloads.DiscountValueRequest) (masterentities.AgreementDiscount, *exceptions.BaseErrorResponse)
	DeleteDiscountValue(*gorm.DB, int, int) *exceptions.BaseErrorResponse
	GetDiscountGroupAgreementById(*gorm.DB, int, int) (masterpayloads.DiscountGroupRequest, *exceptions.BaseErrorResponse)
	GetDiscountItemAgreementById(*gorm.DB, int, int) (masterpayloads.ItemDiscountRequest, *exceptions.BaseErrorResponse)
	GetDiscountValueAgreementById(*gorm.DB, int, int) (masterpayloads.DiscountValueRequest, *exceptions.BaseErrorResponse)
	GetAllDiscountGroup(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetAllItemDiscount(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetAllDiscountValue(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
}

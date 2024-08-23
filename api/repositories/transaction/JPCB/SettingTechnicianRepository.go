package transactionjpcbrepository

import (
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionjpcbpayloads "after-sales/api/payloads/transaction/JPCB"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type SettingTechnicianRepository interface {
	GetAllSettingTechnician(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetAllSettingTechnicianDetail(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetSettingTechnicianById(tx *gorm.DB, settingTechnicianId int) (transactionjpcbpayloads.SettingTechnicianGetByIdResponse, *exceptions.BaseErrorResponse)
	GetSettingTechnicianDetailById(tx *gorm.DB, settingTechnicianDetailId int) (transactionjpcbpayloads.SettingTechnicianDetailGetByIdResponse, *exceptions.BaseErrorResponse)
	SaveSettingTechnician(tx *gorm.DB, CompanyId int) (transactionjpcbpayloads.SettingTechnicianGetByIdResponse, *exceptions.BaseErrorResponse)
	SaveSettingTechnicianDetail(tx *gorm.DB, req transactionjpcbpayloads.SettingTechnicianDetailSaveRequest) (transactionjpcbpayloads.SettingTechnicianDetailGetByIdResponse, *exceptions.BaseErrorResponse)
	UpdateSettingTechnicianDetail(tx *gorm.DB, settingTechnicianDetailId int, req transactionjpcbpayloads.SettingTechnicianDetailUpdateRequest) (transactionjpcbpayloads.SettingTechnicianDetailGetByIdResponse, *exceptions.BaseErrorResponse)
}

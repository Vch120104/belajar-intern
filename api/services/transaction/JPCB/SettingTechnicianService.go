package transactionjpcbservice

import (
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionjpcbpayloads "after-sales/api/payloads/transaction/JPCB"
	"after-sales/api/utils"
)

type SettingTechnicianService interface {
	GetAllSettingTechnician(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetAllSettingTechnicianDetail(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetSettingTechnicianById(settingTechnicianId int) (transactionjpcbpayloads.SettingTechnicianGetByIdResponse, *exceptions.BaseErrorResponse)
	GetSettingTechnicianDetailById(settingTechnicianDetailId int) (transactionjpcbpayloads.SettingTechnicianDetailGetByIdResponse, *exceptions.BaseErrorResponse)
	SaveSettingTechnician(CompanyId int) (transactionjpcbpayloads.SettingTechnicianGetByIdResponse, *exceptions.BaseErrorResponse)
	SaveSettingTechnicianDetail(req transactionjpcbpayloads.SettingTechnicianDetailSaveRequest) (transactionjpcbpayloads.SettingTechnicianDetailGetByIdResponse, *exceptions.BaseErrorResponse)
	UpdateSettingTechnicianDetail(settingTechnicianDetailId int, req transactionjpcbpayloads.SettingTechnicianDetailUpdateRequest) (transactionjpcbpayloads.SettingTechnicianDetailGetByIdResponse, *exceptions.BaseErrorResponse)
}

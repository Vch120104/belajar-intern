package transactionworkshoprepository

import (
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type WorkOrderRepository interface {
	New(tx *gorm.DB, request transactionworkshoppayloads.WorkOrderRequest) (bool, *exceptions.BaseErrorResponse)

	NewStatus(tx *gorm.DB, filter []utils.FilterCondition) ([]transactionworkshopentities.WorkOrderMasterStatus, *exceptions.BaseErrorResponse)
	NewType(tx *gorm.DB, filter []utils.FilterCondition) ([]transactionworkshopentities.WorkOrderMasterType, *exceptions.BaseErrorResponse)
	NewBill(tx *gorm.DB) ([]transactionworkshoppayloads.WorkOrderBillable, *exceptions.BaseErrorResponse)
	NewDropPoint(tx *gorm.DB) ([]transactionworkshoppayloads.WorkOrderDropPoint, *exceptions.BaseErrorResponse)
	NewVehicleBrand(tx *gorm.DB) ([]transactionworkshoppayloads.WorkOrderVehicleBrand, *exceptions.BaseErrorResponse)
	NewVehicleModel(tx *gorm.DB, brandId int) ([]transactionworkshoppayloads.WorkOrderVehicleModel, *exceptions.BaseErrorResponse)

	GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetById(tx *gorm.DB, Id int) (transactionworkshoppayloads.WorkOrderRequest, *exceptions.BaseErrorResponse)
	Save(tx *gorm.DB, request transactionworkshoppayloads.WorkOrderRequest, workOrderId int) (bool, *exceptions.BaseErrorResponse)
	Submit(tx *gorm.DB, Id int) (bool, string, *exceptions.BaseErrorResponse)
	Void(tx *gorm.DB, workOrderId int) (bool, *exceptions.BaseErrorResponse)
	CloseOrder(tx *gorm.DB, Id int) (bool, *exceptions.BaseErrorResponse)

	VehicleLookup(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	CampaignLookup(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)

	GetAllRequest(*gorm.DB, []utils.FilterCondition, pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetRequestById(*gorm.DB, int, int) (transactionworkshoppayloads.WorkOrderServiceRequest, *exceptions.BaseErrorResponse)
	UpdateRequest(*gorm.DB, int, int, transactionworkshoppayloads.WorkOrderServiceRequest) *exceptions.BaseErrorResponse
	AddRequest(*gorm.DB, int, transactionworkshoppayloads.WorkOrderServiceRequest) *exceptions.BaseErrorResponse
	DeleteRequest(*gorm.DB, int, int) *exceptions.BaseErrorResponse

	GetAllVehicleService(*gorm.DB, []utils.FilterCondition, pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetVehicleServiceById(*gorm.DB, int, int) (transactionworkshoppayloads.WorkOrderServiceVehicleRequest, *exceptions.BaseErrorResponse)
	UpdateVehicleService(*gorm.DB, int, int, transactionworkshoppayloads.WorkOrderServiceVehicleRequest) *exceptions.BaseErrorResponse
	AddVehicleService(*gorm.DB, int, transactionworkshoppayloads.WorkOrderServiceVehicleRequest) *exceptions.BaseErrorResponse
	DeleteVehicleService(*gorm.DB, int, int) *exceptions.BaseErrorResponse
	GenerateDocumentNumber(tx *gorm.DB, workOrderId int) (string, *exceptions.BaseErrorResponse)
}

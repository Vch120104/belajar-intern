package transactionworkshoprepositoryimpl

import (
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshoprepository "after-sales/api/repositories/transaction/workshop"

	exceptionsss_test "after-sales/api/expectionsss"
	"after-sales/api/utils"
	"log"

	"gorm.io/gorm"
)

type WorkOrderRepositoryImpl struct {
	DB *gorm.DB
}

func OpenWorkOrderRepositoryImpl() transactionworkshoprepository.WorkOrderRepository {
	return &WorkOrderRepositoryImpl{}
}

func (r *WorkOrderRepositoryImpl) WithTrx(Trxhandle *gorm.DB) transactionworkshoprepository.WorkOrderRepository {
	if Trxhandle == nil {
		log.Println("Transaction Database Not Found")
		return r
	}

	// Wrap the provided database handle with a transaction interceptor
	trx := Trxhandle.Session(&gorm.Session{SkipDefaultTransaction: true})

	// Register a callback for after each transaction commit
	trx.Callback().Create().After("gorm:commit_or_rollback_transaction").Register("log_transaction_commit", func(tx *gorm.DB) {
		// Log the transaction details or perform any other action
		log.Println("Transaction Committed Successfully")
	})

	r.DB = trx
	return r
}

func (r *WorkOrderRepositoryImpl) Save(request transactionworkshoppayloads.WorkOrderRequest) (bool, error) {
	var WorkOrderEntities = transactionworkshopentities.WorkOrder{}
	rows, err := r.DB.Model(&WorkOrderEntities).Save(&WorkOrderEntities).Rows()
	if err != nil {
		return false, err
	}
	defer rows.Close()
	return true, nil
}

func (r *WorkOrderRepositoryImpl) GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse) {
	var entities []transactionworkshopentities.WorkOrder
	// Query to retrieve all work order entities based on the request
	query := tx.Model(&transactionworkshopentities.WorkOrder{})
	if len(filterCondition) > 0 {
		query = query.Where(filterCondition)
	}
	err := query.Find(&entities).Error
	if err != nil {
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{Message: "Failed to retrieve work orders from the database"}
	}

	var workOrderResponses []map[string]interface{}

	// Loop through each entity and copy its data to the response
	for _, entity := range entities {
		workOrderData := make(map[string]interface{})
		// Copy data from entity to response
		workOrderData["BatchSystemNumber"] = entity.BatchSystemNumber
		workOrderData["WorkOrderSystemNumber"] = entity.WorkOrderSystemNumber
		workOrderData["WorkOrderDate"] = entity.WorkOrderDate
		workOrderData["WorkOrderTypeId"] = entity.WorkOrderTypeId
		workOrderData["WorkOrderServiceSite"] = entity.WorkOrderServiceSite
		workOrderData["BrandId"] = entity.BrandId
		workOrderData["ModelId"] = entity.ModelId
		workOrderData["VariantId"] = entity.VariantId
		workOrderData["VehicleId"] = entity.VehicleId
		workOrderData["BilltoCustomerId"] = entity.BilltoCustomerId
		workOrderData["WorkOrderStatusEra"] = entity.WorkOrderStatusEra
		workOrderData["WorkOrderEraNo"] = entity.WorkOrderEraNo
		workOrderData["WorkOrderEraExpiredDate"] = entity.WorkOrderEraExpiredDate
		workOrderData["QueueSystemNumber"] = entity.QueueSystemNumber
		workOrderData["WorkOrderArrivalTime"] = entity.WorkOrderArrivalTime
		workOrderData["WorkOrderCurrentMileage"] = entity.WorkOrderCurrentMileage
		workOrderData["WorkOrderStatusStoring"] = entity.WorkOrderStatusStoring
		workOrderData["WorkOrderRemark"] = entity.WorkOrderRemark
		workOrderData["WorkOrderStatusUnregistered"] = entity.WorkOrderStatusUnregistered
		workOrderData["WorkOrderProfitCenter"] = entity.WorkOrderProfitCenter
		workOrderData["WorkOrderDealerRepCode"] = entity.WorkOrderDealerRepCode
		workOrderData["CampaignId"] = entity.CampaignId
		workOrderData["AgreementId"] = entity.AgreementId
		workOrderData["ServiceRequestSystemNumber"] = entity.ServiceRequestSystemNumber
		workOrderData["EstimationSystemNumber"] = entity.EstimationSystemNumber
		workOrderData["ContractSystemNumber"] = entity.ContractSystemNumber
		workOrderData["CompanyId"] = entity.CompanyId
		workOrderData["DealerRepresentativeId"] = entity.DealerRepresentativeId
		workOrderData["Titleprefix"] = entity.Titleprefix
		workOrderData["NameCust"] = entity.NameCust
		workOrderData["PhoneCust"] = entity.PhoneCust
		workOrderData["MobileCust"] = entity.MobileCust
		workOrderData["MobileCustAlternative"] = entity.MobileCustAlternative
		workOrderData["MobileCustDriver"] = entity.MobileCustDriver
		workOrderData["ContactVia"] = entity.ContactVia
		workOrderData["WorkOrderStatusInsurance"] = entity.WorkOrderStatusInsurance
		workOrderData["WorkOrderInsurancePolicyNo"] = entity.WorkOrderInsurancePolicyNo
		workOrderData["WorkOrderInsuranceExpiredDate"] = entity.WorkOrderInsuranceExpiredDate
		workOrderData["WorkOrderInsuranceClaimNo"] = entity.WorkOrderInsuranceClaimNo
		workOrderData["WorkOrderInsurancePic"] = entity.WorkOrderInsurancePic
		workOrderData["WorkOrderInsuranceWONumber"] = entity.WorkOrderInsuranceWONumber
		workOrderData["WorkOrderInsuranceOwnRisk"] = entity.WorkOrderInsuranceOwnRisk

		workOrderResponses = append(workOrderResponses, workOrderData)
	}

	// Paginate the response data
	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(workOrderResponses, &pages)

	return paginatedData, totalPages, totalRows, nil
}

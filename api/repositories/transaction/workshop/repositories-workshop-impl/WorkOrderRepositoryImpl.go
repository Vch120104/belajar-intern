package transactionworkshoprepositoryimpl

import (
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshoprepository "after-sales/api/repositories/transaction/workshop"

	exceptionsss_test "after-sales/api/expectionsss"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type WorkOrderRepositoryImpl struct {
	DB *gorm.DB
}

func OpenWorkOrderRepositoryImpl() transactionworkshoprepository.WorkOrderRepository {
	return &WorkOrderRepositoryImpl{}
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

func (r *WorkOrderRepositoryImpl) New(tx *gorm.DB, request transactionworkshoppayloads.WorkOrderRequest) (bool, *exceptionsss_test.BaseErrorResponse) {
	// Create a new instance of WorkOrderRepositoryImpl
	repo := &WorkOrderRepositoryImpl{
		DB: tx,
	}

	// Save the work order
	success, err := repo.Save(request)
	if err != nil {
		return false, &exceptionsss_test.BaseErrorResponse{Message: "Failed to save work order"}
	}

	return success, nil
}

func (r *WorkOrderRepositoryImpl) GetById(tx *gorm.DB, Id int) (transactionworkshoppayloads.WorkOrderRequest, *exceptionsss_test.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrder
	err := tx.Model(&transactionworkshopentities.WorkOrder{}).Where("id = ?", Id).First(&entity).Error
	if err != nil {
		return transactionworkshoppayloads.WorkOrderRequest{}, &exceptionsss_test.BaseErrorResponse{Message: "Failed to retrieve work order from the database"}
	}

	// Convert entity to payload
	payload := transactionworkshoppayloads.WorkOrderRequest{
		BatchSystemNumber:             entity.BatchSystemNumber,
		WorkOrderSystemNumber:         entity.WorkOrderSystemNumber,
		WorkOrderDate:                 entity.WorkOrderDate,
		WorkOrderTypeId:               entity.WorkOrderTypeId,
		WorkOrderServiceSite:          entity.WorkOrderServiceSite,
		BrandId:                       entity.BrandId,
		ModelId:                       entity.ModelId,
		VariantId:                     entity.VariantId,
		VehicleId:                     entity.VehicleId,
		BilltoCustomerId:              entity.BilltoCustomerId,
		WorkOrderStatusEra:            entity.WorkOrderStatusEra,
		WorkOrderEraNo:                entity.WorkOrderEraNo,
		WorkOrderEraExpiredDate:       entity.WorkOrderEraExpiredDate,
		QueueSystemNumber:             entity.QueueSystemNumber,
		WorkOrderArrivalTime:          entity.WorkOrderArrivalTime,
		WorkOrderCurrentMileage:       entity.WorkOrderCurrentMileage,
		WorkOrderStatusStoring:        entity.WorkOrderStatusStoring,
		WorkOrderRemark:               entity.WorkOrderRemark,
		WorkOrderStatusUnregistered:   entity.WorkOrderStatusUnregistered,
		WorkOrderProfitCenter:         entity.WorkOrderProfitCenter,
		WorkOrderDealerRepCode:        entity.WorkOrderDealerRepCode,
		CampaignId:                    entity.CampaignId,
		AgreementId:                   entity.AgreementId,
		ServiceRequestSystemNumber:    entity.ServiceRequestSystemNumber,
		EstimationSystemNumber:        entity.EstimationSystemNumber,
		ContractSystemNumber:          entity.ContractSystemNumber,
		CompanyId:                     entity.CompanyId,
		DealerRepresentativeId:        entity.DealerRepresentativeId,
		Titleprefix:                   entity.Titleprefix,
		NameCust:                      entity.NameCust,
		PhoneCust:                     entity.PhoneCust,
		MobileCust:                    entity.MobileCust,
		MobileCustAlternative:         entity.MobileCustAlternative,
		MobileCustDriver:              entity.MobileCustDriver,
		ContactVia:                    entity.ContactVia,
		WorkOrderStatusInsurance:      entity.WorkOrderStatusInsurance,
		WorkOrderInsurancePolicyNo:    entity.WorkOrderInsurancePolicyNo,
		WorkOrderInsuranceExpiredDate: entity.WorkOrderInsuranceExpiredDate,
		WorkOrderInsuranceClaimNo:     entity.WorkOrderInsuranceClaimNo,
		WorkOrderInsurancePic:         entity.WorkOrderInsurancePic,
		WorkOrderInsuranceWONumber:    entity.WorkOrderInsuranceWONumber,
		WorkOrderInsuranceOwnRisk:     entity.WorkOrderInsuranceOwnRisk,
	}

	return payload, nil
}

func (r *WorkOrderRepositoryImpl) Save(request transactionworkshoppayloads.WorkOrderRequest) (bool, error) {
	var WorkOrderEntities = transactionworkshopentities.WorkOrder{}
	err := r.DB.Model(&WorkOrderEntities).Create(&WorkOrderEntities).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *WorkOrderRepositoryImpl) Submit(tx *gorm.DB, Id int) *exceptionsss_test.BaseErrorResponse {
	// Retrieve the work order by Id
	var entity transactionworkshopentities.WorkOrder
	err := tx.Model(&transactionworkshopentities.WorkOrder{}).Where("id = ?", Id).First(&entity).Error
	if err != nil {
		return &exceptionsss_test.BaseErrorResponse{Message: "Failed to retrieve work order from the database"}
	}

	// Perform the necessary operations to submit the work order
	// ...

	// Save the updated work order
	err = tx.Save(&entity).Error
	if err != nil {
		return &exceptionsss_test.BaseErrorResponse{Message: "Failed to save the updated work order"}
	}

	return nil
}

func (r *WorkOrderRepositoryImpl) Void(tx *gorm.DB, Id int) *exceptionsss_test.BaseErrorResponse {
	// Retrieve the work order by Id
	var entity transactionworkshopentities.WorkOrder
	err := tx.Model(&transactionworkshopentities.WorkOrder{}).Where("id = ?", Id).First(&entity).Error
	if err != nil {
		return &exceptionsss_test.BaseErrorResponse{Message: "Failed to retrieve work order from the database"}
	}

	// Perform the necessary operations to void the work order
	// ...

	// Save the updated work order
	err = tx.Save(&entity).Error
	if err != nil {
		return &exceptionsss_test.BaseErrorResponse{Message: "Failed to save the updated work order"}
	}

	return nil
}

func (r *WorkOrderRepositoryImpl) CloseOrder(tx *gorm.DB, Id int) *exceptionsss_test.BaseErrorResponse {
	// Retrieve the work order by Id
	var entity transactionworkshopentities.WorkOrder
	err := tx.Model(&transactionworkshopentities.WorkOrder{}).Where("id = ?", Id).First(&entity).Error
	if err != nil {
		return &exceptionsss_test.BaseErrorResponse{Message: "Failed to retrieve work order from the database"}
	}

	// Perform the necessary operations to close the work order
	// ...

	// Save the updated work order
	err = tx.Save(&entity).Error
	if err != nil {
		return &exceptionsss_test.BaseErrorResponse{Message: "Failed to save the updated work order"}
	}

	return nil
}

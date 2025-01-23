package transactionworkshoprepositoryimpl

import (
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshoprepository "after-sales/api/repositories/transaction/workshop"
	"after-sales/api/utils"
	aftersalesserviceapiutils "after-sales/api/utils/aftersales-service"
	generalserviceapiutils "after-sales/api/utils/general-service"
	salesserviceapiutils "after-sales/api/utils/sales-service"
	"errors"
	"fmt"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type AtpmClaimRegistrationRepositoryImpl struct {
}

func OpenAtpmClaimRegistrationRepositoryImpl() transactionworkshoprepository.AtpmClaimRegistrationRepository {
	return &AtpmClaimRegistrationRepositoryImpl{}
}

// uspg_atAtpmVehicleClaim0_Select
// IF @Option = 0
func (r *AtpmClaimRegistrationRepositoryImpl) GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {

	var entities []transactionworkshopentities.AtpmClaimVehicle

	tx = utils.ApplyFilter(tx.Model(&transactionworkshopentities.AtpmClaimVehicle{}), filterCondition)

	tx.Scopes(pagination.Paginate(&pages, tx)).Find(&entities)

	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Data not found",
				Err:        tx.Error,
			}
		}
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get data",
			Err:        tx.Error,
		}
	}

	if len(entities) == 0 {
		pages.Rows = []map[string]interface{}{}
		return pages, nil
	}

	var results []map[string]interface{}
	for _, entity := range entities {
		brandResponses, brandErr := salesserviceapiutils.GetUnitBrandById(entity.BrandId)
		if brandErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: brandErr.StatusCode,
				Message:    "Failed to fetch brand data from external service",
				Err:        brandErr.Err,
			}
		}

		modelResponses, modelErr := salesserviceapiutils.GetUnitModelById(entity.ModelId)
		if modelErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: modelErr.StatusCode,
				Message:    "Failed to fetch model data from external service",
				Err:        modelErr.Err,
			}
		}

		variantResponses, variantErr := salesserviceapiutils.GetUnitVariantById(entity.VariantId)
		if variantErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch variant data from external service",
				Err:        variantErr.Err,
			}
		}

		companyResponses, companyErr := generalserviceapiutils.GetCompanyDataById(entity.CompanyId)
		if companyErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: companyErr.StatusCode,
				Message:    "Failed to fetch company data from internal service",
				Err:        companyErr.Err,
			}
		}

		// vehicleResponse, vehicleErr := salesserviceapiutils.GetVehicleById(entity.VehicleId)
		// if vehicleErr != nil {
		// 	return pages, &exceptions.BaseErrorResponse{
		// 		StatusCode: vehicleErr.StatusCode,
		// 		Message:    "Failed to fetch vehicle data from external service",
		// 		Err:        vehicleErr.Err,
		// 	}
		// }

		// fetch claim type
		claimTypeResponse, claimTypeErr := generalserviceapiutils.GetClaimTypeById(entity.ClaimTypeId)
		if claimTypeErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: claimTypeErr.StatusCode,
				Message:    "Failed to fetch claim type data from internal service",
				Err:        claimTypeErr.Err,
			}
		}

		result := map[string]interface{}{
			"vehicle_id":                 entity.VehicleId,
			"vehicle_chassis_number":     "vehicleResponse.Data.Master.VehicleChassisNumber",
			"work_order_system_number":   entity.WorkOrderSystemNumber,
			"work_order_document_number": entity.WorkOrderDocumentNumber,
			"work_order_date":            entity.WorkOrderDate,
			"company_id":                 entity.CompanyId,
			"company_name":               companyResponses.CompanyName,
			"claim_system_number":        entity.ClaimSystemNumber,
			"claim_number":               entity.ClaimNumber,
			"claim_date":                 entity.ClaimDate,
			"claim_type_id":              entity.ClaimTypeId,
			"claim_type_description":     claimTypeResponse.ClaimTypeDescription,
			"claim_status_id":            entity.ClaimStatusId,
			"brand_id":                   entity.BrandId,
			"brand_name":                 brandResponses.BrandName,
			"model_id":                   entity.ModelId,
			"model_description":          modelResponses.ModelName,
			"variant_id":                 entity.VariantId,
			"variant_description":        variantResponses.VariantDescription,
			"claim_from":                 entity.ClaimFrom,
			"claim_to":                   entity.ClaimTo,
		}

		results = append(results, result)
	}

	pages.Rows = results
	return pages, nil
}

func (r *AtpmClaimRegistrationRepositoryImpl) GetById(tx *gorm.DB, id int, pages pagination.Pagination) (transactionworkshoppayloads.AtpmClaimRegistrationResponse, *exceptions.BaseErrorResponse) {

	var entity transactionworkshopentities.AtpmClaimVehicle
	err := tx.Model(&transactionworkshopentities.AtpmClaimVehicle{}).
		Where("claim_system_number = ?", id).
		First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return transactionworkshoppayloads.AtpmClaimRegistrationResponse{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Data not found",
				Err:        err,
			}
		}
		return transactionworkshoppayloads.AtpmClaimRegistrationResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get data",
			Err:        err,
		}
	}

	// Get company data
	companyResponses, companyErr := generalserviceapiutils.GetCompanyDataById(entity.CompanyId)
	if companyErr != nil {
		return transactionworkshoppayloads.AtpmClaimRegistrationResponse{}, companyErr
	}

	// Get brand data
	brandResponses, brandErr := salesserviceapiutils.GetUnitBrandById(entity.BrandId)
	if brandErr != nil {
		return transactionworkshoppayloads.AtpmClaimRegistrationResponse{}, brandErr
	}

	// Get model data
	modelResponses, modelErr := salesserviceapiutils.GetUnitModelById(entity.ModelId)
	if modelErr != nil {
		return transactionworkshoppayloads.AtpmClaimRegistrationResponse{}, modelErr
	}

	// Get variant data
	variantResponses, variantErr := salesserviceapiutils.GetUnitVariantById(entity.VariantId)
	if variantErr != nil {
		return transactionworkshoppayloads.AtpmClaimRegistrationResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch variant data from external service",
			Err:        variantErr.Err,
		}
	}

	// Claim Type
	claimTypeResponse, claimTypeErr := generalserviceapiutils.GetClaimTypeById(entity.ClaimTypeId)
	if claimTypeErr != nil {
		return transactionworkshoppayloads.AtpmClaimRegistrationResponse{}, claimTypeErr
	}

	// Get vehicle data
	// vehicleResponse, vehicleErr := salesserviceapiutils.GetVehicleById(entity.VehicleId)
	// if vehicleErr != nil {
	// 	return transactionworkshoppayloads.AtpmClaimRegistrationResponse{}, vehicleErr
	// }

	response := transactionworkshoppayloads.AtpmClaimRegistrationResponse{
		ClaimSystemNumber:       entity.ClaimSystemNumber,
		CompanyId:               entity.CompanyId,
		CompanyName:             companyResponses.CompanyName,
		BrandId:                 entity.BrandId,
		BrandName:               brandResponses.BrandName,
		ClaimTypeId:             entity.ClaimTypeId,
		ClaimTypeDescription:    claimTypeResponse.ClaimTypeDescription,
		ClaimNumber:             entity.ClaimNumber,
		ClaimDate:               entity.ClaimDate,
		WorkOrderDocumentNumber: entity.WorkOrderDocumentNumber,
		WorkOrderDate:           entity.WorkOrderDate,
		VehicleId:               entity.VehicleId,
		VehicleChassisNumber:    "vehicleResponse.Data.Master.VehicleChassisNumber",
		VehicleEngineNumber:     "vehicleResponse.Data.Master.VehicleEngineNumber",
		ModelId:                 entity.ModelId,
		ModelDescription:        modelResponses.ModelName,
		VariantId:               entity.VariantId,
		VariantDescription:      variantResponses.VariantDescription,
	}

	return response, nil
}

// uspg_atAtpmVehicleClaim0_Insert
// IF @Option = 0
func (r *AtpmClaimRegistrationRepositoryImpl) New(tx *gorm.DB, request transactionworkshoppayloads.AtpmClaimRegistrationRequest) (transactionworkshopentities.AtpmClaimVehicle, *exceptions.BaseErrorResponse) {
	// Get approval draft from external service
	approvalDraft, approvalDraftErr := generalserviceapiutils.GetApprovalStatusByCode("10")
	if approvalDraftErr != nil {
		return transactionworkshopentities.AtpmClaimVehicle{}, &exceptions.BaseErrorResponse{
			StatusCode: approvalDraftErr.StatusCode,
			Message:    "Failed to fetch approval draft data from external service",
			Err:        approvalDraftErr.Err,
		}
	}

	// Get work order drom external service
	workOrder, workOrderErr := aftersalesserviceapiutils.GetWorkOrderById(request.WorkOrderSystemNumber)
	if workOrderErr != nil {
		return transactionworkshopentities.AtpmClaimVehicle{}, &exceptions.BaseErrorResponse{
			StatusCode: workOrderErr.StatusCode,
			Message:    "Failed to fetch work order data from external service",
			Err:        workOrderErr.Err,
		}
	}

	fmt.Println("workOrder.Data.WorkOrderDocumentNumber", workOrder.WorkOrderDocumentNumber)
	fmt.Println("workOrder.Data.WorkOrderDate", workOrder.WorkOrderDate)

	// if workOrder.Data.WorkOrderInformation.WorkOrderDocumentNumber == "" || workOrder.Data.WorkOrderInformation.WorkOrderDate.IsZero() {
	// 	return transactionworkshopentities.AtpmClaimVehicle{}, &exceptions.BaseErrorResponse{
	// 		StatusCode: http.StatusBadRequest,
	// 		Message:    "Incomplete work order data from external service",
	// 		Err:        errors.New("missing required work order fields"),
	// 	}
	// }

	var existingEntity transactionworkshopentities.AtpmClaimVehicle
	if err := tx.Where("claim_system_number = ?", request.ClaimSystemNumber).
		First(&existingEntity).Error; err == nil {
		return transactionworkshopentities.AtpmClaimVehicle{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Message:    "Data already exists",
			Err:        errors.New("duplicate entry"),
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return transactionworkshopentities.AtpmClaimVehicle{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to check existing data",
			Err:        err,
		}
	}

	entity := transactionworkshopentities.AtpmClaimVehicle{
		CompanyId:               request.CompanyId,
		BrandId:                 request.BrandId,
		ModelId:                 request.ModelId,
		VariantId:               request.VariantId,
		ClaimTypeId:             request.ClaimTypeId,
		ClaimStatusId:           approvalDraft.ApprovalStatusId,
		CustomerComplaint:       request.CustomerComplaint,
		TechnicianDiagnostic:    request.TechnicianDiagnostic,
		Countermeasure:          request.Countermeasure,
		ClaimDate:               request.ClaimDate,
		RepairEndDate:           request.RepairEndDate,
		WorkOrderSystemNumber:   request.WorkOrderSystemNumber,
		WorkOrderDocumentNumber: workOrder.WorkOrderDocumentNumber,
		WorkOrderDate:           workOrder.WorkOrderDate,

		// Other data
		Fuel:       request.Fuel,
		CustomerId: request.CustomerId,
		Vdn:        request.VDN,

		// Claim Header, Symptom, Trouble Code
		ClaimHeader: request.ClaimHeader,
	}

	if err := tx.Create(&entity).Error; err != nil {
		return transactionworkshopentities.AtpmClaimVehicle{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to create data",
			Err:        err,
		}
	}

	// logic ClaimTypeId = FSI (1)
	if request.ClaimTypeId == 1 {

		// Subquery untuk mendapatkan TOTAL_AFTER_DISCOUNT
		subQuery := tx.Model(&transactionworkshopentities.AtpmWarranty{}).
			Select("COALESCE(total_after_discount, 0)").
			Where("brand_id = ? AND model_id = ? AND variant_id = ? AND fsp_category_id = ? AND effective_date = (?)",
				request.BrandId, request.ModelId, request.VariantId, request.FspCategoryId,
				tx.Model(&transactionworkshopentities.AtpmWarranty{}).
					Select("MAX(effective_date)").
					Where("brand_id = ? AND model_id = ? AND variant_id = ? AND fsp_category_id = ? AND effective_date <= ?",
						request.BrandId, request.ModelId, request.VariantId, request.FspCategoryId, request.ClaimDate),
			)

		// Update FspAmountClaimStandard dengan hasil subquery
		if err := tx.Model(&transactionworkshopentities.AtpmClaimVehicle{}).
			Where("claim_system_number = ?", request.ClaimSystemNumber).
			Update("fsp_amount_claim_standard", subQuery).Error; err != nil {
			return transactionworkshopentities.AtpmClaimVehicle{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to update FSP amount",
				Err:        err,
			}
		}
	}
	return entity, nil
}

// uspg_atAtpmVehicleClaim0_Update
// IF @Option = 0 / 2
func (r *AtpmClaimRegistrationRepositoryImpl) Save(tx *gorm.DB, id int, request transactionworkshoppayloads.AtpmClaimRegistrationRequestSave) (transactionworkshopentities.AtpmClaimVehicle, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.AtpmClaimVehicle

	approvalDraft, approvalDraftErr := generalserviceapiutils.GetApprovalStatusByCode("10")
	if approvalDraftErr != nil {
		return transactionworkshopentities.AtpmClaimVehicle{}, &exceptions.BaseErrorResponse{
			StatusCode: approvalDraftErr.StatusCode,
			Message:    "Failed to fetch approval draft data from external service",
			Err:        approvalDraftErr.Err,
		}
	}

	err := tx.Where("claim_system_number = ?", id).First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return transactionworkshopentities.AtpmClaimVehicle{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Data not found",
				Err:        err,
			}
		}
		return transactionworkshopentities.AtpmClaimVehicle{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch data",
			Err:        err,
		}
	}

	if entity.ClaimStatusId != approvalDraft.ApprovalStatusId {
		return transactionworkshopentities.AtpmClaimVehicle{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusForbidden,
			Message:    "Update header failed, claim document is already submitted",
			Err:        errors.New("update header failed, claim document is already submitted"),
		}
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		updateData := map[string]interface{}{
			"symptom_code": request.SymptomCode,
			"trouble_code": request.TroubleCode,
		}

		if err := tx.Model(&transactionworkshopentities.AtpmClaimVehicle{}).
			Where("claim_system_number = ?", id).
			Updates(updateData).Error; err != nil {
			return transactionworkshopentities.AtpmClaimVehicle{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to update SymptomCode and TroubleCode",
				Err:        err,
			}
		}

		if err := tx.Where("claim_system_number = ?", id).First(&entity).Error; err != nil {
			return transactionworkshopentities.AtpmClaimVehicle{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch updated data",
				Err:        err,
			}
		}

		return entity, nil
	}

	entity.CustomerComplaint = request.CustomerComplaint
	entity.TechnicianDiagnostic = request.TechnicianDiagnostic
	entity.Countermeasure = request.Countermeasure
	entity.RepairEndDate = request.RepairEndDate
	entity.Fuel = request.Fuel
	entity.CustomerId = request.CustomerId
	entity.Vdn = request.VDN
	entity.ClaimHeader = request.ClaimHeader

	if err := tx.Save(&entity).Error; err != nil {
		return transactionworkshopentities.AtpmClaimVehicle{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update data",
			Err:        err,
		}
	}

	return entity, nil
}

// uspg_atAtpmVehicleClaim0_Update
// IF @Option = 1
func (r *AtpmClaimRegistrationRepositoryImpl) Submit(tx *gorm.DB, id int) (bool, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.AtpmClaimVehicle
	var entitywo transactionworkshopentities.WorkOrderDetail

	// Step 1: Retrieve claim details
	if err := tx.Where("claim_system_number = ?", id).First(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Data not found",
				Err:        err,
			}
		}
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch data",
			Err:        err,
		}
	}

	// Step 2: Validate claim status and draft
	approvalDraft, approvalDraftErr := generalserviceapiutils.GetApprovalStatusByCode("10")
	if approvalDraftErr != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: approvalDraftErr.StatusCode,
			Message:    "Failed to fetch approval draft data from external service",
			Err:        approvalDraftErr.Err,
		}
	}

	if entity.ClaimStatusId != approvalDraft.ApprovalStatusId {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusForbidden,
			Message:    "Claim document is already submitted",
			Err:        errors.New("claim document is already submitted"),
		}
	}

	// Step 3: Validate mandatory attachments in AtpmClaimVehicleAttachmentType
	if err := tx.Model(&transactionworkshopentities.AtpmClaimVehicleAttachmentType{}).
		Joins("INNER JOIN trx_atpm_claim_vehicle ON trx_atpm_claim_vehicle.claim_to = trx_atpm_claim_vehicle_attachment_type.atpm_code").
		Joins("LEFT JOIN trx_atpm_claim_vehicle_detail ON trx_atpm_claim_vehicle_detail.claim_system_number = trx_atpm_claim_vehicle.claim_system_number").
		Where("trx_atpm_claim_vehicle_attachment_type.mandatory = 1").
		Where("trx_atpm_claim_vehicle.claim_system_number = ?", id).
		Where("trx_atpm_claim_vehicle_detail.claim_system_number IS NULL").
		Find(&transactionworkshopentities.AtpmClaimVehicleAttachmentType{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusPreconditionFailed,
				Message:    "Mandatory attachments are not complete",
				Err:        err,
			}
		}
	}

	// Step 4: Check existence of Use DMS
	// Req RPS/07/21/00336
	var useDmsExist int
	if err := tx.Table("dms_microservices_general_dev.dbo.mtr_company_reference AS ref").
		Select("COUNT(1)").
		Joins("LEFT JOIN dms_microservices_aftersales_dev.dbo.trx_atpm_claim_vehicle AS atpm ON atpm.claim_from = ref.company_id").
		Where("atpm.claim_system_number = ? AND ref.use_dms = 1", id).
		Scan(&useDmsExist).Error; err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to check use DMS",
			Err:        err,
		}
	}

	if useDmsExist > 0 {
		var count int64

		// Check WO detail is not QC Passed)
		if err := tx.Table("trx_work_order_detail AS w2").
			Select("COUNT(1)").
			Joins("LEFT JOIN trx_service_log sl ON sl.work_order_system_number = w2.work_order_system_number AND sl.operation_item_id = w2.operation_item_id AND ISNULL(sl.service_status_id, '') IN (?, ?, ?)", utils.SrvStatStop, utils.SrvStatQcPass, utils.SrvStatTransfer).
			Joins("INNER JOIN trx_atpm_claim_vehicle_detail cl ON w2.work_order_system_number = cl.work_order_system_number").
			Where("cl.claim_system_number = ? AND w2.line_type_id = 2 AND w2.transaction_type_id IN ('8', '10') AND ISNULL(w2.BYPASS, '') <> '1' AND ISNULL(sl.service_status_id, '') = ''", id). // 8 = F Free Service, 10 = W Warranty
			Count(&count).Error; err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to check if WO detail is not QC Passed",
				Err:        err,
			}
		}

		if count > 0 {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "WO detail is not QC Passed",
			}
		}

		// Check claim without item/service detail)
		var existsClaim bool
		if err := tx.Table("trx_atpm_claim_vehicle_detail").
			Select("1").
			Where("claim_system_number = ? AND line_type_id = '2'", id).
			Limit(1).
			Scan(&existsClaim).Error; err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to check if claim has details",
				Err:        err,
			}
		}

		if !existsClaim {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Claim does not have item or service details",
			}
		}

		var woDate, claimDate time.Time
		if err := tx.Table("trx_atpm_claim_vehicle A").
			Select("A.work_order_date, A.claim_date").
			Joins("INNER JOIN trx_work_order B ON A.work_order_system_number = B.work_order_system_number").
			Where("A.claim_system_number = ?", id).
			First(&woDate).Error; err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch work order and claim dates",
				Err:        err,
			}
		}

		// Check if the claim date is more than 10 days after QC passed date for claims after the cutoff date
		if woDate.After(time.Date(2021, 10, 8, 0, 0, 0, 0, time.UTC)) {
			if claimDate.Sub(woDate).Hours()/24 > 10 {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusBadRequest,
					Message:    "Claim Date More than 10 days from QC Passed date",
				}
			}
		}

		// Check for the QC Passed Date in WO details
		var startDate time.Time
		if err := tx.Table("trx_work_order_detail w2").
			Select("w2.quality_control_pass_datetime").
			Joins("INNER JOIN trx_atpm_claim_vehicle_detail cl ON w2.work_order_system_number = cl.work_order_system_number").
			Where("cl.claim_system_number = ? AND w2.line_type_id = 2 AND w2.transaction_type_id IN ('8', '10')", id). // 8 = F Free Service, 10 = W Warranty
			Order("w2.quality_control_pass_datetime DESC").
			Limit(1).
			Scan(&startDate).Error; err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to check WO QC passed date",
				Err:        err,
			}
		}

		// If no QC passed date
		if startDate.IsZero() {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "WO detail is not QC Passed",
			}
		}

		// Check claim date is more than 10 days after the QC Passed Date
		if claimDate.Sub(startDate).Hours()/24 > 10 {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Claim Date More than 10 days from QC Passed date",
			}
		}

		var exists bool
		var billingCustCodeExists bool
		var servBookNo string

		// Check if the vehicle has a service book number
		if servBookNo == "" {
			if err := tx.Table("trx_atpm_claim_vehicle_detail A").
				Joins("INNER JOIN trx_work_order_detail B ON A.claim_system_number = B.claim_system_number").
				Where("A.claim_system_number = ? AND B.transaction_type_id IN ('8', '10')", id). // 8 = F Free Service, 10 = W Warranty
				Limit(1).
				Scan(&exists).Error; err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to check if vehicle has service book number",
					Err:        err,
				}
			}

			if exists {
				if entity.CompanyId != 0 {
					if err := tx.Table("dms_microservices_general_dev.dbo.mtr_company").
						Where("company_id = ?", entity.CompanyId).
						Limit(1).
						Scan(&billingCustCodeExists).Error; err != nil {
						return false, &exceptions.BaseErrorResponse{
							StatusCode: http.StatusInternalServerError,
							Message:    "Failed to check if Billing Customer Code exists",
							Err:        err,
						}
					}
				}

				if entity.CompanyId != 0 && !billingCustCodeExists {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusBadRequest,
						Message:    "This Vehicle has no Service Book No...",
					}
				}
			}
		} else {

			var exists bool

			if servBookNo == "" {
				if err := tx.Table("trx_atpm_claim_vehicle_detail A").
					Joins("INNER JOIN trx_work_order_detail B ON A.work_order_system_number = B.work_order_system_number AND A.work_order_line_number = B.work_order_operation_item_line").
					Where("A.claim_system_number = ? AND B.transaction_type_id IN ('8', '10')", id). // 8 = F Free Service, 10 = W Warranty
					Limit(1).
					Scan(&exists).Error; err != nil {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to check if service book number exists",
						Err:        err,
					}
				}

				if exists {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusBadRequest,
						Message:    "This Vehicle has no Service Book No...",
					}
				}
			}

		}

	}

	if entity.ClaimFrom == "151" {
		var exists bool
		if err := tx.Table("trx_atpm_claim_vehicle A").
			Joins("INNER JOIN trx_atpm_claim_vehicle_detail B ON A.claim_system_number = B.claim_system_number").
			Joins("INNER JOIN trx_work_order_detail C ON B.work_order_system_number = C.work_order_system_number AND B.work_order_line_number = C.work_order_operation_item_line").
			Where("A.claim_system_number = ? AND (ISNULL(C.atpm_claim_number,'') <> '' AND ISNULL(A.claim_number,'') <> ISNULL(C.atpm_claim_number,''))", id).
			Limit(1).
			Scan(&exists).Error; err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to check if referenced WO detail is already claimed",
				Err:        err,
			}
		}

		if exists {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Referenced WO detail is already claimed",
			}
		} else {
			if err := tx.Table("trx_work_order_detail C").
				Joins("INNER JOIN trx_atpm_claim_vehicle_detail B ON A.claim_system_number = B.claim_system_number").
				Joins("INNER JOIN trx_atpm_claim_vehicle A ON A.claim_system_number = B.claim_system_number").
				Where("A.claim_system_number = ?", id).
				Update("C.atpm_claim_number", entitywo.AtpmClaimNumber).
				Update("C.atpm_claim_date", entitywo.AtpmClaimDate).
				Update("C.claim_system_number", id).Error; err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to update work order details with claim info",
					Err:        err,
				}
			}
		}
	}

	if entitywo.AtpmClaimNumber == "" {
		// Call to external procedure (Dummy document number update)
		// Generate Dummy Document Number (Claim Doc No)
		// Call dbo.uspg_gmSrcDoc1_Update here
		// TODO: Implement logic for dbo.uspg_gmSrcDoc1_Update
		// @Option = 0, -- int (specifies the option)
		// @COMPANY_CODE = @Claim_To, -- varchar (company code of the claim)
		// @TRANSACTION_DATE = @Change_Datetime, -- datetime (date of the transaction)
		// @SOURCE_CODE = @Src_Doc_Type, -- varchar (source document type)
		// @VEHICLE_BRAND = @Vehicle_Brand, -- varchar (brand of the vehicle)
		// @PROFIT_CENTER_CODE = '', -- varchar (profit center code, currently empty)
		// @TRANSACTION_CODE = '', -- varchar (transaction code, currently empty)
		// @BANK_ACC_CODE = '', -- varchar (bank account code, currently empty)
		// @Change_User_Id = @Change_User_Id, -- varchar (ID of the user making the change)
		// @Last_Doc_No = @CLAIM_NO OUTPUT -- varchar (outputs the last document number)
		// --End Generate Dummy Document Number--
		var lastDocNo string
		// if err := tx.Exec("EXEC uspg_gmSrcDoc1_Update @Option = 0, @COMPANY_CODE = ?, @TRANSACTION_DATE = ?, @SOURCE_CODE = ?, @VEHICLE_BRAND = ?, @PROFIT_CENTER_CODE = ?, @TRANSACTION_CODE = ?, @BANK_ACC_CODE = ?, @Change_User_Id = ?, @Last_Doc_No = ?",
		// 	claimTo, changeDatetime, SrcDocType, vehicleBrand, "", "", "", changeUserId, &lastDocNo).Error; err != nil {
		// 	return transactionworkshopentities.AtpmClaimVehicle{}, &exceptions.BaseErrorResponse{
		// 		StatusCode: http.StatusInternalServerError,
		// 		Message:    "Failed to update dummy document number",
		// 		Err:        err,
		// 	}
		// }

		// Set CLAIM_NO to the last document number
		entity.ClaimNumber = lastDocNo
	}

	if entitywo.AtpmClaimNumber == "" {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Document Master Is not Valid",
		}
	}

	//--ADD APPROVAL_REQ HERE
	// 1. Generate Approval Code

	// 2. Generate Approval Request Remark

	// 3. Call usp_comApprovalReq_Insert
	// Generate Approval Request
	// Call usp_comApprovalReq_Insert stored procedure
	// TODO: Implement logic for usp_comApprovalReq_Insert
	// @Company_Code = claimTo
	// @Approval_Code = apvCode
	// @Src_Doc_Type = srcDocType
	// @Src_Sys_No = claimSysNo
	// @Module_Code = 'GR'
	// @Src_Doc_Date = claimDate
	// @Src_Doc_Amount = totalAfterDisc
	// @Req_Remark = apvReqRemark
	// @Req_No = apvReqNo OUTPUT
	// @Change_No = 0
	// @Creation_User_Id = changeUserId
	// @Change_User_Id = changeUserId

	// 4. Update atpm.atAtpmVehicleClaim0

	return true, nil
}

// uspg_atAtpmVehicleClaim0_Delete
// IF @Option = 0
func (r *AtpmClaimRegistrationRepositoryImpl) Void(tx *gorm.DB, id int) (bool, *exceptions.BaseErrorResponse) {

	var claim struct {
		ClaimStatusId int
		ClaimNumber   string
	}

	err := tx.Model(&transactionworkshopentities.AtpmClaimVehicle{}).
		Select("claim_status_id, claim_number").
		Where("claim_system_number = ? AND claim_status_id IN (?)", id, 1). // 1 = Draft
		Limit(1).
		Scan(&claim).Error

	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to check claim status",
			Err:        err,
		}
	}

	if claim.ClaimStatusId != 1 {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Void ATPM Claim failed, claim document is already submitted",
			Err:        errors.New("claim document is already submitted"),
		}
	}

	if claim.ClaimNumber == "" {
		if err := tx.Where("claim_system_number = ?", id).
			Delete(&transactionworkshopentities.AtpmClaimVehicleDetail{}).Error; err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to delete from trx_atpm_claim_vehicle_detail",
				Err:        err,
			}
		}

		if err := tx.Where("claim_system_number = ?", id).
			Delete(&transactionworkshopentities.AtpmClaimVehicle{}).Error; err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to delete from trx_atpm_claim_vehicle",
				Err:        err,
			}
		}
	} else {
		if err := tx.Model(&transactionworkshopentities.AtpmClaimVehicle{}).
			Where("claim_system_number = ?", id).
			Update("claim_status_id", 4).Error; // 4 = Canceled
		err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to update claim status to canceled",
				Err:        err,
			}
		}
	}

	return true, nil
}

// uspg_atAtpmVehicleClaim0_Select
// IF @Option = 6
func (r *AtpmClaimRegistrationRepositoryImpl) GetAllServiceHistory(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {

	var results []map[string]interface{}

	// Fetch approval status by code
	approvalStatus, approvalStatusErr := generalserviceapiutils.GetApprovalStatusByCode("20")
	if approvalStatusErr != nil {
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: approvalStatusErr.StatusCode,
			Message:    "Failed to fetch approval status data from external service",
			Err:        approvalStatusErr.Err,
		}
	}

	baseQuery := tx.Table("trx_work_order A").
		Select("DISTINCT A.work_order_date, A.service_mileage, A.work_order_document_number, C.account_receivable_work_order_document_number, C.total_after_vat_amount, A.company_id, E.work_order_service_request_line").
		Joins("INNER JOIN trx_work_order_detail B ON A.work_order_system_number = B.work_order_system_number").
		Joins("LEFT JOIN dms_microservices_finance_dev.dbo.trx_account_receivable_work_order C ON B.invoice_system_number = C.invoice_receipt_system_number AND C.account_receivable_work_order_approval_status_id = ?", approvalStatus.ApprovalStatusId).
		Joins("LEFT JOIN trx_work_order_service_request E ON E.work_order_system_number = A.work_order_system_number").
		Where("B.job_type_id = ?", 7). // 7 = Periodical Maintenance
		Order("A.work_order_date DESC").
		Limit(pages.GetLimit()).
		Offset(pages.GetOffset())

	tx = utils.ApplyFilter(baseQuery, filterCondition)

	if result := tx.Scopes(pagination.Paginate(&pages, tx)).Find(&results); result.RowsAffected == 0 {
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Data not found",
			Err:        nil,
		}
	}

	if len(results) == 0 {
		pages.Rows = []map[string]interface{}{}
		return pages, nil
	}

	for _, entity := range results {

		companyId, ok := entity["company_id"].(int)
		if !ok {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Invalid company_id type",
				Err:        nil,
			}
		}

		companyResponses, companyErr := generalserviceapiutils.GetCompanyDataById(companyId)
		if companyErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: companyErr.StatusCode,
				Message:    "Failed to fetch company data from internal service",
				Err:        companyErr.Err,
			}
		}

		result := map[string]interface{}{
			"work_order_date":                               entity["work_order_date"],
			"service_mileage":                               entity["service_mileage"],
			"work_order_document_number":                    entity["work_order_document_number"],
			"account_receivable_work_order_document_number": entity["account_receivable_work_order_document_number"],
			"total_after_vat_amount":                        entity["total_after_vat_amount"],
			"company_id":                                    entity["company_id"],
			"company_name":                                  companyResponses.CompanyName,
			"work_order_service_request_line":               entity["work_order_service_request_line"],
		}

		results = append(results, result)
	}

	pages.Rows = results
	return pages, nil
}

// uspg_atAtpmVehicleClaim0_Select
// IF @Option = 7
func (r *AtpmClaimRegistrationRepositoryImpl) GetAllClaimHistory(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {

	var results []map[string]interface{}

	// Fetch approval status by code
	approvalStatus, approvalStatusErr := generalserviceapiutils.GetApprovalStatusByCode("10")
	if approvalStatusErr != nil {
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: approvalStatusErr.StatusCode,
			Message:    "Failed to fetch approval status data from external service",
			Err:        approvalStatusErr.Err,
		}
	}

	baseQuery := tx.Table("trx_atpm_claim_vehicle A").
		Select("A.claim_system_number, A.claim_number, A.claim_date, A.claim_status_id, A.company_id, A.pfp").
		Joins("LEFT JOIN trx_atpm_claim_vehicle_detail B ON A.claim_system_number = B.claim_system_number").
		Where("A.claim_status_id <> ?", approvalStatus.ApprovalStatusId).
		Order("A.claim_date DESC").
		Limit(pages.GetLimit()).
		Offset(pages.GetOffset())

	tx = utils.ApplyFilter(baseQuery, filterCondition)

	if result := tx.Scopes(pagination.Paginate(&pages, tx)).Find(&results); result.RowsAffected == 0 {
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Data not found",
			Err:        nil,
		}
	}

	if len(results) == 0 {
		pages.Rows = []map[string]interface{}{}
		return pages, nil
	}

	for _, entity := range results {

		companyId, ok := entity["company_id"].(int)
		if !ok {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Invalid company_id type",
				Err:        nil,
			}
		}

		companyResponses, companyErr := generalserviceapiutils.GetCompanyDataById(companyId)
		if companyErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: companyErr.StatusCode,
				Message:    "Failed to fetch company data from internal service",
				Err:        companyErr.Err,
			}
		}

		result := map[string]interface{}{
			"claim_system_number": entity["claim_system_number"],
			"claim_number":        entity["claim_number"],
			"claim_date":          entity["claim_date"],
			"claim_status_id":     entity["claim_status_id"],
			"company_id":          entity["company_id"],
			"company_name":        companyResponses.CompanyName,
		}

		results = append(results, result)
	}

	pages.Rows = results
	return pages, nil
}

// uspg_atAtpmVehicleClaim1_Select
// IF @Option = 0
func (r *AtpmClaimRegistrationRepositoryImpl) GetAllDetail(tx *gorm.DB, id int, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var results []map[string]interface{}

	baseQuery := tx.Model(&transactionworkshopentities.AtpmClaimVehicleDetail{}).
		Select(`
		trx_atpm_claim_vehicle_detail.claim_system_number,
		trx_atpm_claim_vehicle_detail.company_id,
		trx_atpm_claim_vehicle_detail.claim_line_number,
		trx_atpm_claim_vehicle_detail.work_order_system_number,
		trx_atpm_claim_vehicle_detail.work_order_line_number,
		trx_atpm_claim_vehicle_detail.line_type_id,
		trx_atpm_claim_vehicle_detail.item_id,
		CASE
			WHEN COALESCE(trx_atpm_claim_vehicle_detail.line_type_id, 0) = 0 THEN E.package_name
			WHEN COALESCE(trx_atpm_claim_vehicle_detail.line_type_id, 0) = 1 THEN D.operation_name
			ELSE C.item_name
		END AS item_name,
		trx_atpm_claim_vehicle_detail.frt_quantity,
		trx_atpm_claim_vehicle_detail.item_price,
		trx_atpm_claim_vehicle_detail.discount_percent,
		trx_atpm_claim_vehicle_detail.discount_amount,
		trx_atpm_claim_vehicle_detail.total_after_discount,
		trx_atpm_claim_vehicle_detail.recall_number,
		COALESCE(trx_atpm_claim_vehicle_detail.part_request, 0) AS part_request,
		COALESCE(trx_atpm_claim_vehicle_detail.incident_part_received, 0) AS incident_part_received
	`).
		Joins("LEFT JOIN mtr_item AS C ON trx_atpm_claim_vehicle_detail.item_id = C.item_id").
		Joins("LEFT JOIN mtr_operation_code AS D ON trx_atpm_claim_vehicle_detail.item_id = D.operation_id").
		Joins("LEFT JOIN mtr_package AS E ON trx_atpm_claim_vehicle_detail.item_id = E.package_id").
		Where("trx_atpm_claim_vehicle_detail.claim_system_number = ?", id).
		Order("trx_atpm_claim_vehicle_detail.claim_line_number ASC").
		Limit(pages.GetLimit()).
		Offset(pages.GetOffset())

	tx = utils.ApplyFilter(baseQuery, filterCondition)

	result := tx.Scopes(pagination.Paginate(&pages, tx)).Find(&results)
	if result.Error != nil {
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve data",
			Err:        result.Error,
		}
	}

	if len(results) == 0 {
		pages.Rows = []map[string]interface{}{}
		return pages, nil
	}

	pages.Rows = results
	return pages, nil
}

// uspg_atAtpmVehicleClaim1_Select
// IF @Option = 1
func (r *AtpmClaimRegistrationRepositoryImpl) GetDetailById(tx *gorm.DB, claimsysno int, detailid int, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var results []map[string]interface{}
	var workOrderSysNo int

	if result := tx.Model(&transactionworkshopentities.AtpmClaimVehicle{}).
		Where("claim_system_number = ?", claimsysno).
		Pluck("work_order_system_number", &workOrderSysNo); result.Error != nil || result.RowsAffected == 0 {
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Data not found",
			Err:        result.Error,
		}
	}

	baseQuery := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Select(`
			work_order_system_number,
			work_order_operation_item_line,
			line_type_id,
			operation_item_id,
			frt_quantity
		`).
		Where("work_order_system_number = ?", workOrderSysNo).
		Order("work_order_operation_item_line ASC").
		Limit(pages.GetLimit()).
		Offset(pages.GetOffset())

	// check claim type id in trx_atpm_claim_vehicle
	var claimTypeId int
	if result := tx.Model(&transactionworkshopentities.AtpmClaimVehicle{}).
		Where("claim_system_number = ?", claimsysno).
		Pluck("claim_type_id", &claimTypeId); result.Error != nil {
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve data",
			Err:        result.Error,
		}
	}

	if claimTypeId == 3 { //CLAIMTYPE_PDI
		baseQuery = baseQuery.Where("transaction_type_id = 10") // 10 = Warranty
	} else if claimTypeId == 1 { //CLAIMTYPE_FREESERVICE
		baseQuery = baseQuery.Where("transaction_type_id = 8") // 8 = Free Service
	} else if claimTypeId == 5 { //CLAIMTYPE_WARRANTY
		baseQuery = baseQuery.Where("transaction_type_id = 10") // 10 = Warranty
	} else if claimTypeId == 6 { //CLAIMTYPE_WARRANTY_PART
		baseQuery = baseQuery.Where("transaction_type_id = 10") // 10 = Warranty
	}

	result := baseQuery.Scopes(pagination.Paginate(&pages, baseQuery)).Find(&results)
	if result.Error != nil {
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve data",
			Err:        result.Error,
		}
	}

	if len(results) == 0 {
		pages.Rows = []map[string]interface{}{}
		return pages, nil
	}

	for i, row := range results {
		lineTypeId, ok := row["line_type_id"].(int)
		if !ok {
			continue // Skip if line_type_id is not present or not an int
		}

		lineType, errResp := generalserviceapiutils.GetLineTypeById(lineTypeId)
		if errResp != nil {
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve line type information",
				Err:        errResp.Err,
			}
		}

		results[i]["line_type_code"] = lineType.LineTypeCode
	}

	pages.Rows = results
	return pages, nil
}

// uspg_atAtpmVehicleClaim1_Insert
// IF @Option = 2
func (r *AtpmClaimRegistrationRepositoryImpl) AddDetail(tx *gorm.DB, id int, request transactionworkshoppayloads.AtpmClaimDetailRequest) (transactionworkshopentities.AtpmClaimVehicleDetail, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.AtpmClaimVehicleDetail

	// Check if claim system number exists
	var claimExists bool
	if err := tx.Model(&transactionworkshopentities.AtpmClaimVehicle{}).
		Where("claim_system_number = ?", id).
		Select("1").
		Limit(1).
		Scan(&claimExists).Error; err != nil {
		return entity, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to check if claim system number exists",
			Err:        err,
		}
	}

	if !claimExists {
		return entity, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Claim system number not found",
			Err:        nil,
		}
	}

	// Get the next claim line number
	var nextClaimLineNumber int
	if err := tx.Model(&transactionworkshopentities.AtpmClaimVehicleDetail{}).
		Select("COALESCE(MAX(claim_line_number), 0) + 1").
		Where("claim_system_number = ?", id).
		Scan(&nextClaimLineNumber).Error; err != nil {
		return entity, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get next claim line number",
			Err:        err,
		}
	}

	// Fetch required work order details
	var workOrder transactionworkshopentities.WorkOrderDetail
	if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Where("work_order_system_number = ? AND work_order_line_number = ? AND claim_system_number IS NULL",
			request.WorkOrderSystemNumber, request.WorkOrderLineNumber).
		First(&workOrder).Error; err != nil {
		return entity, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Work order details not found or already claimed",
			Err:        err,
		}
	}

	// Set entity values based on work order details
	entity.ClaimSystemNumber = id
	entity.ClaimLineNumber = nextClaimLineNumber
	entity.WorkOrderSystemNumber = request.WorkOrderSystemNumber
	entity.WorkOrderLineNumber = request.WorkOrderLineNumber
	entity.LineTypeId = workOrder.LineTypeId
	entity.ItemId = workOrder.OperationItemId
	entity.FrtQuantity = workOrder.FrtQuantity
	entity.DiscountPercent = request.DiscountPercent
	entity.DiscountAmount = request.DiscountAmount
	entity.ItemPrice = workOrder.OperationItemPrice
	entity.TotalAfterDiscount = (workOrder.OperationItemPrice - request.DiscountAmount) * workOrder.FrtQuantity
	entity.PartRequest = 0
	entity.IncidentPartReceived = 0

	if err := tx.Create(&entity).Error; err != nil {
		return entity, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to insert data",
			Err:        err,
		}
	}

	if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Where("work_order_system_number = ? AND work_order_line_number = ?", request.WorkOrderSystemNumber, request.WorkOrderLineNumber).
		Update("claim_system_number", id).Error; err != nil {
		return entity, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update work order with claim system number",
			Err:        err,
		}
	}

	// Update total calculations in ATPM claim
	if err := tx.Model(&transactionworkshopentities.AtpmClaimVehicle{}).
		Where("claim_system_number = ?", id).
		Updates(map[string]interface{}{
			"total_after_discount": tx.Model(&transactionworkshopentities.AtpmClaimVehicleDetail{}).
				Select("COALESCE(SUM(total_after_discount), 0)").Where("claim_system_number = ?", id),
			"total_labour": tx.Model(&transactionworkshopentities.AtpmClaimVehicleDetail{}).
				Select("COALESCE(SUM(total_after_discount), 0)").Where("claim_system_number = ? AND line_type_id = 1", id),
			"total_part": tx.Model(&transactionworkshopentities.AtpmClaimVehicleDetail{}).
				Select("COALESCE(SUM(total_after_discount), 0)").Where("claim_system_number = ? AND line_type_id <> 1", id),
			"total_frt_qty": tx.Model(&transactionworkshopentities.AtpmClaimVehicleDetail{}).
				Select("COALESCE(SUM(frt_quantity), 0)").Where("claim_system_number = ? AND line_type_id = 1", id),
		}).Error; err != nil {
		return entity, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update claim totals",
			Err:        err,
		}
	}

	return entity, nil
}

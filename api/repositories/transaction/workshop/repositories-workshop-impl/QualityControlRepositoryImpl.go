package transactionworkshoprepositoryimpl

import (
	"after-sales/api/config"
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshoprepository "after-sales/api/repositories/transaction/workshop"
	"after-sales/api/utils"
	generalserviceapiutils "after-sales/api/utils/general-service"
	salesserviceapiutils "after-sales/api/utils/sales-service"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type QualityControlRepositoryImpl struct {
}

func OpenQualityControlRepositoryImpl() transactionworkshoprepository.QualityControlRepository {
	return &QualityControlRepositoryImpl{}
}

// uspg_wtWorkOrder0_Select
// IF @Option = 7
// USE IN MODUL : AWS - 006 QUALITY CONTROL PAGE 1 REQ: ???
func (r *QualityControlRepositoryImpl) GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {

	var entities []transactionworkshoppayloads.QualityControlRequest

	joinTable := utils.CreateJoinSelectStatement(tx, transactionworkshoppayloads.QualityControlRequest{})
	whereQuery := utils.ApplyFilterSearch(joinTable, filterCondition)
	whereQuery = whereQuery.Where("work_order_status_id = ?", utils.WoStatStop) // 40 Stop

	if err := whereQuery.Find(&entities).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Work order not found",
				Err:        err,
			}
		}
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch entity",
			Err:        err,
		}
	}

	if len(entities) == 0 {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "No data found",
			Err:        errors.New("no data found"),
		}
	}

	var convertedResponses []transactionworkshoppayloads.QualityControlResponse

	for _, entity := range entities {
		// Fetch data model from external services
		modelResponses, modelErr := salesserviceapiutils.GetUnitModelById(entity.ModelId)
		if modelErr != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch model data from external service",
				Err:        modelErr.Err,
			}
		}

		// Fetch data variant from external services
		variantResponses, variantErr := salesserviceapiutils.GetUnitVariantById(entity.VariantId)
		if variantErr != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch variant data from external service",
				Err:        variantErr.Err,
			}
		}

		// Fetch data vehicle from external API
		vehicleResponses, vehicleErr := salesserviceapiutils.GetVehicleById(entity.VehicleId)
		if vehicleErr != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve vehicle data from the external API",
				Err:        vehicleErr.Err,
			}
		}

		// Fetch data customer from external API
		customerResponses, customerErr := generalserviceapiutils.GetCustomerMasterDetailById(entity.CustomerId)
		if customerErr != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve customer data from the external API",
				Err:        customerErr.Err,
			}
		}

		// Fetch data work order from external API
		WorkOrderUrl := config.EnvConfigs.AfterSalesServiceUrl + "work-order/normal/" + strconv.Itoa(entity.WorkOrderSystemNumber)
		var workOrderResponses transactionworkshoppayloads.WorkOrderResponse
		errWorkOrder := utils.Get(WorkOrderUrl, &workOrderResponses, nil)
		if errWorkOrder != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve work order data from the external API",
				Err:        errWorkOrder,
			}
		}

		// Append converted response
		convertedResponses = append(convertedResponses, transactionworkshoppayloads.QualityControlResponse{
			WorkOrderDocumentNumber: workOrderResponses.WorkOrderDocumentNumber,
			WorkOrderDate:           workOrderResponses.WorkOrderDate.Format(time.RFC3339),
			VehicleCode:             vehicleResponses.VehicleChassisNumber,
			VehicleTnkb:             vehicleResponses.VehicleRegistrationCertificateTNKB,
			CustomerName:            customerResponses.CustomerName,
			WorkOrderSystemNumber:   entity.WorkOrderSystemNumber,
			VarianCode:              variantResponses.VariantCode,
			ModelCode:               modelResponses.ModelCode,
		})
	}

	var mapResponses []map[string]interface{}
	for _, response := range convertedResponses {
		responseMap := map[string]interface{}{
			"work_order_document_number":            response.WorkOrderDocumentNumber,
			"work_order_date":                       response.WorkOrderDate,
			"model_code":                            response.ModelCode,
			"varian_code":                           response.VarianCode,
			"vehicle_chassis_number":                response.VehicleCode,
			"vehicle_registration_certificate_tnkb": response.VehicleTnkb,
			"customer_name":                         response.CustomerName,
			"work_order_system_number":              response.WorkOrderSystemNumber,
		}
		mapResponses = append(mapResponses, responseMap)
	}

	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(mapResponses, &pages)

	return paginatedData, totalPages, totalRows, nil
}

func (r *QualityControlRepositoryImpl) GetById(tx *gorm.DB, id int, filterCondition []utils.FilterCondition, pages pagination.Pagination) (transactionworkshoppayloads.QualityControlIdResponse, *exceptions.BaseErrorResponse) {
	var entity transactionworkshoppayloads.QualityControlRequest

	joinTable := utils.CreateJoinSelectStatement(tx, entity)
	whereQuery := utils.ApplyFilter(joinTable, filterCondition)
	whereQuery = whereQuery.Where("work_order_system_number = ? AND work_order_status_id IN (?,?)", id, utils.WoStatStop, utils.WoStatOngoing)

	if err := whereQuery.Find(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return transactionworkshoppayloads.QualityControlIdResponse{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Work order not found",
				Err:        err,
			}
		}
		return transactionworkshoppayloads.QualityControlIdResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch entity",
			Err:        err,
		}
	}

	// Fetch data brand from external API
	brandResponses, brandErr := salesserviceapiutils.GetUnitBrandById(entity.BrandId)
	if brandErr != nil {
		return transactionworkshoppayloads.QualityControlIdResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch brand data from external service",
			Err:        brandErr.Err,
		}
	}

	// Fetch data model from external services
	modelResponses, modelErr := salesserviceapiutils.GetUnitModelById(entity.ModelId)
	if modelErr != nil {
		return transactionworkshoppayloads.QualityControlIdResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch model data from external service",
			Err:        modelErr.Err,
		}
	}

	// Fetch data variant from external services
	variantResponses, variantErr := salesserviceapiutils.GetUnitVariantById(entity.VariantId)
	if variantErr != nil {
		return transactionworkshoppayloads.QualityControlIdResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch variant data from external service",
			Err:        variantErr.Err,
		}
	}

	// Fetch data vehicle from external API
	vehicleResponses, vehicleErr := salesserviceapiutils.GetVehicleById(entity.VehicleId)
	if vehicleErr != nil {
		return transactionworkshoppayloads.QualityControlIdResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve vehicle data from the external API",
			Err:        vehicleErr.Err,
		}
	}

	// Fetch data colour from external API
	colourUrl := config.EnvConfigs.SalesServiceUrl + "unit-color-dropdown/" + strconv.Itoa(entity.BrandId)
	var colourResponses []transactionworkshoppayloads.WorkOrderVehicleColour
	errColour := utils.GetArray(colourUrl, &colourResponses, nil)
	if errColour != nil || len(colourResponses) == 0 {
		return transactionworkshoppayloads.QualityControlIdResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve colour data from the external API",
			Err:        errColour,
		}
	}

	// Fetch data customer from external API
	CustomerUrl := config.EnvConfigs.SalesServiceUrl + "customer/" + strconv.Itoa(entity.CustomerId)
	var customerResponses transactionworkshoppayloads.CustomerResponse
	errCustomer := utils.Get(CustomerUrl, &customerResponses, nil)
	if errCustomer != nil {
		return transactionworkshoppayloads.QualityControlIdResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve customer data from the external API",
			Err:        errCustomer,
		}
	}

	// Fetch foreman data from external API
	ForemanUrl := config.EnvConfigs.GeneralServiceUrl + "user-detail/" + strconv.Itoa(entity.ForemanId)
	var foremanResponses transactionworkshoppayloads.QualityControlForemanResponse
	errForeman := utils.Get(ForemanUrl, &foremanResponses, nil)
	if errForeman != nil {
		return transactionworkshoppayloads.QualityControlIdResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve foreman data from the external API",
			Err:        errForeman,
		}
	}

	// fetch service advisor data from external API
	ServiceAdvisorUrl := config.EnvConfigs.GeneralServiceUrl + "user-detail/" + strconv.Itoa(entity.ServiceAdvisorId)
	var serviceAdvisorResponses transactionworkshoppayloads.QualityControlServiceAdvisorResponse
	errServiceAdvisor := utils.Get(ServiceAdvisorUrl, &serviceAdvisorResponses, nil)
	if errServiceAdvisor != nil {
		return transactionworkshoppayloads.QualityControlIdResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve service advisor data from the external API",
			Err:        errServiceAdvisor,
		}
	}

	// Fetch work order users
	var workorderUsers []transactionworkshoppayloads.WorkOrderCurrentUserResponse
	if err := tx.Table("dms_microservices_general_dev.dbo.mtr_customer AS c").
		Select(`
		c.customer_id AS customer_id,
		c.customer_name AS customer_name,
		c.customer_code AS customer_code,
		c.id_address_id AS address_id,
		a.address_street_1 AS address_street_1,
		a.address_street_2 AS address_street_2,
		a.address_street_3 AS address_street_3,
		a.village_id AS village_id,
		v.village_name AS village_name,
		v.district_id AS district_id,
		d.district_name AS district_name,
		d.city_id AS city_id,
		ct.city_name AS city_name,
		ct.province_id AS province_id,
		p.province_name AS province_name,
		v.village_zip_code AS zip_code,
		td.npwp_no AS current_user_npwp
	`).
		Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_address AS a ON c.id_address_id = a.address_id").
		Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_village AS v ON a.village_id = v.village_id").
		Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_district AS d ON v.district_id = d.district_id").
		Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_city AS ct ON d.city_id = ct.city_id").
		Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_province AS p ON ct.province_id = p.province_id").
		Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_tax_data AS td ON c.tax_customer_id = td.tax_id").
		Where("c.customer_id = ?", entity.CustomerId).
		Find(&workorderUsers).Error; err != nil {
		fmt.Println("Error executing query:", err)
		return transactionworkshoppayloads.QualityControlIdResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order users from the database",
			Err:        err,
		}
	}

	// Fetch work order detail vehicles
	var workorderVehicleDetails []transactionworkshoppayloads.WorkOrderVehicleDetailResponse
	if err := tx.Table("dms_microservices_sales_dev.dbo.mtr_vehicle AS v").
		Select(`
		v.vehicle_id AS vehicle_id,
        v.vehicle_chassis_number AS vehicle_chassis_number,
		vrc.vehicle_registration_certificate_tnkb AS vehicle_registration_certificate_tnkb,
		vrc.vehicle_registration_certificate_owner_name AS vehicle_registration_certificate_owner_name,
		v.vehicle_production_year AS vehicle_production_year,
		CONCAT(vv.variant_code , ' - ', vv.variant_description) AS vehicle_variant,
		v.option_id AS vehicle_option,
		CONCAT(vm.colour_code , ' - ', vm.colour_commercial_name) AS vehicle_colour,
		v.vehicle_sj_date AS vehicle_sj_date,
        v.vehicle_last_service_date AS vehicle_last_service_date,
        v.vehicle_last_km AS vehicle_last_km
		`).
		Joins("INNER JOIN dms_microservices_sales_dev.dbo.mtr_vehicle_registration_certificate AS vrc ON v.vehicle_id = vrc.vehicle_id").
		Joins("INNER JOIN dms_microservices_sales_dev.dbo.mtr_unit_variant AS vv ON v.vehicle_variant_id = vv.variant_id").
		Joins("INNER JOIN dms_microservices_sales_dev.dbo.mtr_colour AS vm ON v.vehicle_colour_id = vm.colour_id").
		Where("v.vehicle_id = ? AND v.vehicle_brand_id = ? and v.vehicle_variant_id = ?", entity.VehicleId, entity.BrandId, entity.VariantId).
		Find(&workorderVehicleDetails).Error; err != nil {
		return transactionworkshoppayloads.QualityControlIdResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order vehicles from the database",
			Err:        err,
		}
	}

	// Fetch data work order from external API
	WorkOrderUrl := config.EnvConfigs.AfterSalesServiceUrl + "work-order/normal/" + strconv.Itoa(entity.WorkOrderSystemNumber)
	var workOrderResponses transactionworkshoppayloads.WorkOrderResponse
	errWorkOrder := utils.Get(WorkOrderUrl, &workOrderResponses, nil)
	if errWorkOrder != nil {
		return transactionworkshoppayloads.QualityControlIdResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order data from the external API",
			Err:        errWorkOrder,
		}
	}

	var qualitycontrolDetails []transactionworkshoppayloads.QualityControlDetailResponse
	var totalRows int64
	totalRowsQuery := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Where("work_order_system_number = ?", entity.WorkOrderSystemNumber).
		Count(&totalRows).Error
	if totalRowsQuery != nil {
		return transactionworkshoppayloads.QualityControlIdResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count quality control details",
			Err:        totalRowsQuery,
		}
	}

	// Fetch paginated qc details
	query := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Joins("INNER JOIN mtr_item AS WTA ON trx_work_order_detail.operation_item_id = WTA.item_id").
		Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_service_status AS MSS ON trx_work_order_detail.service_status_id = MSS.service_status_id").
		Select("WTA.item_id as operation_item_id,WTA.item_code as operation_item_code, WTA.item_name as operation_item_name, trx_work_order_detail.frt_quantity as frt, trx_work_order_detail.service_status_id, MSS.service_status_description as service_status_name").
		Where("trx_work_order_detail.work_order_system_number = ?", id).
		Offset(pages.GetOffset()).
		Limit(pages.GetLimit())

	if err := query.Find(&qualitycontrolDetails).Error; err != nil {
		return transactionworkshoppayloads.QualityControlIdResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get service details",
			Err:        err,
		}
	}

	// Check if the service_status_id is valid
	if len(qualitycontrolDetails) == 0 {
		return transactionworkshoppayloads.QualityControlIdResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Quality control details not found",
			Err:        errors.New("quality control details not found"),
		}
	}

	validStatuses := map[int]bool{
		utils.SrvStatStop:    true,
		utils.SrvStatQcPass:  true,
		utils.SrvStatReOrder: true,
	}
	for _, detail := range qualitycontrolDetails {
		if !validStatuses[detail.ServiceStatusId] {
			return transactionworkshoppayloads.QualityControlIdResponse{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Operation Status is not valid",
				Err:        errors.New("operation status is not valid"),
			}
		}
	}

	response := transactionworkshoppayloads.QualityControlIdResponse{
		WorkOrderDocumentNumber: workOrderResponses.WorkOrderDocumentNumber,
		WorkOrderDate:           workOrderResponses.WorkOrderDate.Format(time.RFC3339),
		BrandName:               brandResponses.BrandName,
		ModelName:               modelResponses.ModelName,
		VariantName:             variantResponses.VariantName,
		ColourName:              colourResponses[0].VariantColourName,
		VehicleCode:             vehicleResponses.VehicleChassisNumber,
		VehicleTnkb:             vehicleResponses.VehicleRegistrationCertificateTNKB,
		CustomerName:            workorderUsers[0].CustomerName,
		Address0:                workorderUsers[0].Address_Street_1,
		Address1:                workorderUsers[0].Address_Street_2,
		RTRW:                    workorderUsers[0].Address_Street_3,
		LastMilage:              workorderVehicleDetails[0].VehicleLastKm,
		CurrentMilage:           workOrderResponses.WorkOrderCurrentMileage,
		Phone:                   workOrderResponses.PhoneCust,
		ForemanName:             foremanResponses.ForemanName,
		ServiceAdvisorName:      serviceAdvisorResponses.ServiceAdvisorName,
		OrderDateTime:           workOrderResponses.WorkOrderDate.Format(time.RFC3339),
		EstimatedFinished:       workOrderResponses.PromiseTime.Format("15:04:05"),
		QualityControlDetails: transactionworkshoppayloads.QualityControlDetailsResponse{
			Page:       pages.GetPage(),
			Limit:      pages.GetLimit(),
			TotalPages: int(math.Ceil(float64(totalRows) / float64(pages.GetLimit()))),
			TotalRows:  int(totalRows),
			Data:       qualitycontrolDetails,
		},
	}

	return response, nil
}

// uspg_wtWorkOrder2_Update
// IF @Option = 2
// USE IN MODUL : AWS - 006  UPDATE DATA BY KEY (QC PASS) - GENERAL REPAIR
func (r *QualityControlRepositoryImpl) Qcpass(tx *gorm.DB, id int, iddet int) (transactionworkshoppayloads.QualityControlUpdateResponse, *exceptions.BaseErrorResponse) {
	var (
		currentStatus     int
		techAllocSysNo    int
		lineTypeOperation = utils.LinetypeOperation
	)

	var maxWoOprItemLine int
	if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Select("COALESCE(MAX(work_order_operation_item_line), 0)").
		Where("work_order_system_number = ?", id).
		Scan(&maxWoOprItemLine).Error; err != nil {
		return transactionworkshoppayloads.QualityControlUpdateResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve maximum work order operation item line",
			Err:        err,
		}
	}
	// Check the current WO_OPR_STATUS
	err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Select("service_status_id").
		Where("work_order_system_number = ? AND work_order_operation_item_line = ?", id, maxWoOprItemLine).
		First(&currentStatus).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return transactionworkshoppayloads.QualityControlUpdateResponse{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Operation Status is not valid",
				Err:        err,
			}
		}
		return transactionworkshoppayloads.QualityControlUpdateResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch operation status",
			Err:        err,
		}
	}

	// Validate the status
	if currentStatus != utils.SrvStatStop {
		return transactionworkshoppayloads.QualityControlUpdateResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "The current status of the work order is not valid",
			Err:        errors.New("the current status of the work order is not valid"),
		}
	}

	// Fetch work order details
	var details struct {
		VehicleId   int    `gorm:"column:vehicle_id"`
		BrandId     int    `gorm:"column:brand_id"`
		CompanyId   int    `gorm:"column:company_id"`
		OprItemCode string `gorm:"column:operation_item_code"`
		WoStatus    int    `gorm:"column:work_order_status_id"`
	}

	err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Select("trx_work_order.vehicle_id, trx_work_order.company_id, trx_work_order_detail.operation_item_code, trx_work_order.work_order_status_id").
		Joins("JOIN trx_work_order ON trx_work_order_detail.work_order_system_number = trx_work_order.work_order_system_number").
		Where("trx_work_order_detail.work_order_system_number = ? AND trx_work_order_detail.work_order_operation_item_line = ?", id, lineTypeOperation).
		Scan(&details).Error
	if err != nil {
		return transactionworkshoppayloads.QualityControlUpdateResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch work order details",
			Err:        err,
		}
	}

	// Fetch vehicle master data
	vehicleResponses, vehicleErr := salesserviceapiutils.GetVehicleById(details.VehicleId)
	if vehicleErr != nil {
		return transactionworkshoppayloads.QualityControlUpdateResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve vehicle data from the external API",
			Err:        vehicleErr.Err,
		}
	}

	// Fetch the latest TechAllocSystemNumber
	err = tx.Model(&transactionworkshopentities.WorkOrderAllocation{}).
		Select("ISNULL(MAX(technician_allocation_system_number), 0)").
		Where("work_order_system_number = ?", id).
		Where("work_order_line = ?", lineTypeOperation).
		Where("brand_id = ?", vehicleResponses.VehicleBrandID).
		Where("company_id = ?", details.CompanyId).
		Where("operation_code = ?", details.OprItemCode).
		Scan(&techAllocSysNo).Error
	if err != nil {
		return transactionworkshoppayloads.QualityControlUpdateResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch the latest TechAllocSystemNumber",
			Err:        err,
		}
	}

	// Update WorkOrderDetail
	err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Where("work_order_system_number = ? AND work_order_operation_item_line = ? and work_order_detail_id = ?", id, lineTypeOperation, iddet).
		Update("service_status_id", utils.SrvStatQcPass).
		Error
	if err != nil {
		return transactionworkshoppayloads.QualityControlUpdateResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update WorkOrderDetail",
			Err:        err,
		}
	}

	// Update WorkOrderAllocation
	err = tx.Model(&transactionworkshopentities.WorkOrderAllocation{}).
		Where("work_order_system_number = ? AND work_order_line = ?", id, lineTypeOperation).
		Update("service_status_id", utils.SrvStatQcPass).
		Error
	if err != nil {
		return transactionworkshoppayloads.QualityControlUpdateResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update WorkOrderAllocation",
			Err:        err,
		}
	}

	// Check if all related items are updated
	var statusCount int64
	err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Where("work_order_system_number = ? AND service_status_id != ?", id, utils.SrvStatQcPass).
		Count(&statusCount).Error
	if err != nil {
		return transactionworkshoppayloads.QualityControlUpdateResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count non-QC pass items",
			Err:        err,
		}
	}

	if statusCount == 0 {
		// Update WorkOrder if all related WorkOrderDetail have service_status_id as utils.SrvStatQcPass
		err = tx.Model(&transactionworkshopentities.WorkOrder{}).
			Where("work_order_system_number = ?", id).
			Update("work_order_status_id", utils.WoStatQC).
			Error
		if err != nil {
			return transactionworkshoppayloads.QualityControlUpdateResponse{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to update WorkOrder",
				Err:        err,
			}
		}
	}

	// Return response
	response := transactionworkshoppayloads.QualityControlUpdateResponse{
		WorkOrderSystemNumber: id,
		WorkOrderDetailId:     iddet,
		WorkOrderStatusId:     utils.SrvStatQcPass,
		WorkOrderStatusName:   "QC Passed",
	}

	return response, nil
}

// uspg_wtWorkOrder2_Update
// IF @Option = 1
// USE IN MODUL : AWS-006 SHEET: RE-ORDER
func (r *QualityControlRepositoryImpl) Reorder(tx *gorm.DB, id int, iddet int, payload transactionworkshoppayloads.QualityControlReorder) (transactionworkshoppayloads.QualityControlUpdateResponse, *exceptions.BaseErrorResponse) {
	var (
		lineTypeOperation = 1
		woLine            = 1
	)

	// Check if the current WO_OPR_STATUS is valid
	var currentStatus int
	err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Select("service_status_id").
		Where("work_order_system_number = ? AND work_order_operation_item_line = ? AND work_order_detail_id = ?", id, lineTypeOperation, iddet). // Assuming 1 is the value for Wo_Opr_Item_Line
		First(&currentStatus).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return transactionworkshoppayloads.QualityControlUpdateResponse{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Operation Status is not valid",
				Err:        err,
			}
		}
		return transactionworkshoppayloads.QualityControlUpdateResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch operation status",
			Err:        err,
		}
	}

	if currentStatus != utils.SrvStatStop {
		return transactionworkshoppayloads.QualityControlUpdateResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Operation Status is not valid",
		}
	}

	// Update atWOTechAlloc
	err = tx.Model(&transactionworkshopentities.WorkOrderAllocation{}).
		Where("work_order_system_number = ? AND work_order_line = ?", id, woLine). // Assuming 1 is the value for Wo_Opr_Item_Line
		Updates(map[string]interface{}{
			"re_order": true,
		}).Error
	if err != nil {
		return transactionworkshoppayloads.QualityControlUpdateResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update atWOTechAlloc",
			Err:        err,
		}
	}

	// Update wtWorkOrder2
	err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Where("work_order_system_number = ? AND work_order_operation_item_line = ? and work_order_detail_id = ?", id, lineTypeOperation, iddet). // Assuming 1 is the value for Wo_Opr_Item_Line
		Updates(map[string]interface{}{
			"service_status_id":               utils.SrvStatReOrder,
			"quality_control_extra_frt":       payload.ExtraTime,
			"quality_control_total_extra_frt": gorm.Expr("quality_control_total_extra_frt + ?", payload.ExtraTime),
			"quality_control_extra_reason":    payload.Reason,
			"reorder_number":                  gorm.Expr("ISNULL(reorder_number, 0) + 1"),
		}).Error
	if err != nil {
		return transactionworkshoppayloads.QualityControlUpdateResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update wtWorkOrder2",
			Err:        err,
		}
	}

	// Update wtWorkOrder0
	err = tx.Model(&transactionworkshopentities.WorkOrder{}).
		Where("work_order_system_number = ?", id).
		Updates(map[string]interface{}{
			"work_order_status_id": utils.WoStatOngoing,
		}).Error
	if err != nil {
		return transactionworkshoppayloads.QualityControlUpdateResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update wtWorkOrder0",
			Err:        err,
		}
	}

	//return a response
	response := transactionworkshoppayloads.QualityControlUpdateResponse{
		WorkOrderSystemNumber: id,
		WorkOrderDetailId:     iddet,
		WorkOrderStatusId:     utils.SrvStatReOrder,
		WorkOrderStatusName:   "ReOrder",
	}

	return response, nil
}

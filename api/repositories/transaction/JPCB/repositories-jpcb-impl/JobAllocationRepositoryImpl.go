package transactionjpcbrepositoryimpl

import (
	"after-sales/api/config"
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionjpcbpayloads "after-sales/api/payloads/transaction/JPCB"
	transactionjpcbrepository "after-sales/api/repositories/transaction/JPCB"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type JobAllocationRepositoryImpl struct {
}

func StartJobAllocationRepositoryImpl() transactionjpcbrepository.JobAllocationRepository {
	return &JobAllocationRepositoryImpl{}
}

func (r *JobAllocationRepositoryImpl) GetAllJobAllocation(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	entities := transactionworkshopentities.WorkOrderAllocation{}
	payloads := []transactionjpcbpayloads.GetAllJobAllocationPayload{}
	responses := []transactionjpcbpayloads.GetAllJobAllocationResponse{}

	itemGroupUrl := config.EnvConfigs.GeneralServiceUrl + "filter-item-group?item_group_code=OJ"
	itemGroupPayloads := []transactionjpcbpayloads.ItemGroupPayload{}
	if err := utils.GetArray(itemGroupUrl, &itemGroupPayloads, nil); err != nil || len(itemGroupPayloads) == 0 {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("fail to retrieve item group data"),
		}
	}

	baseModelQuery := tx.Model(&entities).
		Select(`
			trx_work_order_allocation.technician_allocation_system_number,
			trx_work_order_allocation.technician_id,
			trx_work_order_allocation.service_status_id,
			trx_work_order_allocation.sequence_number,
			CASE WHEN ISNULL(trx_work_order_allocation.work_order_system_number, 0) = 0
				THEN 'BOOKING'
				ELSE 'WORK ORDER'
			END reference_document_type,
			CASE WHEN ISNULL(trx_work_order_allocation.work_order_system_number, 0) = 0
				THEN trx_work_order_allocation.booking_document_number
				ELSE trx_work_order_allocation.work_order_document_number
			END reference_document_number,
			tbe.vehicle_id,
			CASE WHEN ISNULL(moc.operation_id, 0) != 0
				THEN moc.operation_name
				ELSE mi.item_name
			END operation,
			trx_work_order_allocation.frt,
			trx_work_order_allocation.factor_x,
			trx_work_order_allocation.tech_alloc_total_time AS frt_jpcb,
			trx_work_order_allocation.tech_alloc_last_start_time,
			trx_work_order_allocation.tech_alloc_last_end_time,
			trx_work_order_allocation.is_express`).
		Joins("LEFT JOIN trx_booking_estimation tbe ON tbe.booking_system_number = trx_work_order_allocation.booking_system_number").
		Joins("LEFT JOIN mtr_operation_code moc ON moc.operation_code = trx_work_order_allocation.operation_code").
		Joins("LEFT JOIN mtr_item mi ON mi.item_code = trx_work_order_allocation.operation_code AND mi.item_group_id = ?", itemGroupPayloads[0].ItemGroupId).
		Where("trx_work_order_allocation.technician_id IS NOT NULL AND trx_work_order_allocation.technician_id != 0").
		Where("trx_work_order_allocation.operation_code IS NOT NULL AND trx_work_order_allocation.operation_code != ''").
		Where("moc.operation_name IS NOT NULL OR mi.item_name IS NOT NULL")
	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)
	err := whereQuery.Scopes(pagination.Paginate(&entities, &pages, whereQuery)).Scan(&payloads).Error

	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	for _, result := range payloads {
		var TNKB *string
		if result.VehicleId != nil && *result.VehicleId != 0 {
			vehicleUrl := config.EnvConfigs.SalesServiceUrl + "vehicle-master/" + strconv.Itoa(*result.VehicleId)
			vehiclePayloads := transactionjpcbpayloads.VehiclePayload{}
			if err := utils.Get(vehicleUrl, &vehiclePayloads, nil); err != nil || vehiclePayloads.Master.VehicleId == 0 {
				return pages, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        errors.New("fail to retrieve vehicle data"),
				}
			}
			TNKB = &vehiclePayloads.Stnk.VehicleRegistrationCertificateTnkb
		}

		var technicianName *string
		if result.TechnicianId != nil {
			userDetailsUrl := config.EnvConfigs.GeneralServiceUrl + "user-detail/" + strconv.Itoa(*result.TechnicianId)
			userDetailsPayload := transactionjpcbpayloads.UserDetailsPayload{}
			if err := utils.Get(userDetailsUrl, &userDetailsPayload, nil); err != nil || userDetailsPayload.UserEmployeeId == 0 {
				return pages, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        errors.New("fail to retrieve technician data"),
				}
			}
			technicianName = &userDetailsPayload.EmployeeName
		}

		var serviceStatus *string
		if result.ServiceStatusId != nil {
			serviceStatusUrl := config.EnvConfigs.GeneralServiceUrl + "service-status/" + strconv.Itoa(*result.ServiceStatusId)
			serviceStatusPayload := transactionjpcbpayloads.ServiceStatusPayload{}
			if err := utils.Get(serviceStatusUrl, &serviceStatusPayload, nil); err != nil || serviceStatusPayload.ServiceStatusId == 0 {
				return pages, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        errors.New("fail to retrieve service status data"),
				}
			}
			serviceStatus = &serviceStatusPayload.ServiceStatusDescription
		}

		finalResult := transactionjpcbpayloads.GetAllJobAllocationResponse{
			TechnicianAllocationSystemNumber: result.TechnicianAllocationSystemNumber,
			TechnicianName:                   technicianName,
			ServiceStatus:                    serviceStatus,
			SequenceNumber:                   result.SequenceNumber,
			ReferenceDocumentType:            result.ReferenceDocumentType,
			ReferenceDocumentNumber:          result.ReferenceDocumentNumber,
			VehicleTNBK:                      TNKB,
			Operation:                        result.Operation,
			Frt:                              result.Frt,
			FactorX:                          result.FactorX,
			FrtJPCB:                          result.FrtJPCB,
			TechAllocLastStartTime:           result.TechAllocLastStartTime,
			TechAllocLastEndTime:             result.TechAllocLastEndTime,
			IsExpress:                        result.IsExpress,
		}
		responses = append(responses, finalResult)
	}

	pages.Rows = responses

	return pages, nil
}

func (r *JobAllocationRepositoryImpl) GetJobAllocationById(tx *gorm.DB, technicianAllocationSystemNumber int) (transactionjpcbpayloads.GetJobAllocationByIdResponse, *exceptions.BaseErrorResponse) {
	entities := transactionworkshopentities.WorkOrderAllocation{}
	payloads := transactionjpcbpayloads.GetJobAllocationByIdPayload{}
	responses := transactionjpcbpayloads.GetJobAllocationByIdResponse{}

	itemGroupUrl := config.EnvConfigs.GeneralServiceUrl + "filter-item-group?item_group_code=OJ"
	itemGroupPayloads := []transactionjpcbpayloads.ItemGroupPayload{}
	if err := utils.GetArray(itemGroupUrl, &itemGroupPayloads, nil); err != nil || len(itemGroupPayloads) == 0 {
		return responses, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("fail to retrieve item group data"),
		}
	}

	err := tx.Model(&entities).
		Select(`
			trx_work_order_allocation.technician_allocation_system_number,
			trx_work_order_allocation.company_id,
			trx_work_order_allocation.technician_id,
			trx_work_order_allocation.sequence_number,
			CASE WHEN ISNULL(trx_work_order_allocation.work_order_system_number, 0) = 0
				THEN trx_work_order_allocation.booking_document_number
				ELSE trx_work_order_allocation.work_order_document_number
			END reference_document_number,
			CASE WHEN ISNULL(moc.operation_id, 0) != 0
				THEN moc.operation_name
				ELSE mi.item_name
			END operation,
			trx_work_order_allocation.frt,
			trx_work_order_allocation.factor_x,
			trx_work_order_allocation.work_order_system_number`).
		Where("trx_work_order_allocation.technician_allocation_system_number = ?", technicianAllocationSystemNumber).
		Joins("LEFT JOIN mtr_operation_code moc ON moc.operation_code = trx_work_order_allocation.operation_code").
		Joins("LEFT JOIN mtr_item mi ON mi.item_code = trx_work_order_allocation.operation_code AND mi.item_group_id = ?", itemGroupPayloads[0].ItemGroupId).
		Where("trx_work_order_allocation.technician_id IS NOT NULL AND trx_work_order_allocation.technician_id != 0").
		Where("trx_work_order_allocation.operation_code IS NOT NULL AND trx_work_order_allocation.operation_code != ''").
		Where("moc.operation_name IS NOT NULL OR mi.item_name IS NOT NULL").
		First(&payloads).
		Error

	if err != nil {
		return responses, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	var technicianName *string
	if payloads.TechnicianId != nil {
		userDetailsUrl := config.EnvConfigs.GeneralServiceUrl + "user-detail/" + strconv.Itoa(*payloads.TechnicianId)
		userDetailsPayload := transactionjpcbpayloads.UserDetailsPayload{}
		if err := utils.Get(userDetailsUrl, &userDetailsPayload, nil); err != nil || userDetailsPayload.UserEmployeeId == 0 {
			return transactionjpcbpayloads.GetJobAllocationByIdResponse{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        errors.New("fail to retrieve"),
			}
		}
		technicianName = &userDetailsPayload.EmployeeName
	}

	entityServiceLog := transactionworkshopentities.ServiceLog{}
	var progress float64

	err = tx.Model(&entityServiceLog).
		Select("CONVERT(decimal(5,2), ROUND((SUM(actual_time)/(trx_service_log.frt)) * 100, 2)) AS progress").
		Joins("INNER JOIN work_order_operation woo on trx_service_log.work_order_operation_id = woo.work_order_operation_id").
		Joins("INNER JOIN mtr_operation_model_mapping momm on momm.operation_model_mapping_id = woo.operation_id").
		Joins("INNER JOIN mtr_operation_code moc on moc.operation_id = momm.operation_id").
		Joins("LEFT JOIN trx_work_order_allocation twoa ON twoa.work_order_system_number = trx_service_log.work_order_system_number").
		Where("trx_service_log.work_order_system_number = ?", payloads.WorkOrderSystemNumber).
		Where("twoa.operation_code = moc.operation_code").
		Group("twoa.technician_id, trx_service_log.work_order_system_number, trx_service_log.frt, twoa.factor_x, trx_service_log.service_log_system_number").
		Pluck("progress", &progress).
		Error

	if err != nil {
		return responses, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	responses.TechnicianAllocationSystemNumber = payloads.TechnicianAllocationSystemNumber
	responses.Operation = payloads.Operation
	responses.Frt = payloads.Frt
	responses.FactorX = payloads.FactorX
	responses.TechnicianName = technicianName
	responses.SequenceNumber = payloads.SequenceNumber
	responses.Progress = &progress

	return responses, nil
}

func (r *JobAllocationRepositoryImpl) UpdateJobAllocation(tx *gorm.DB, technicianAllocationSystemNumber int, req transactionjpcbpayloads.JobAllocationUpdateRequest) (transactionworkshopentities.WorkOrderAllocation, *exceptions.BaseErrorResponse) {
	entities := transactionworkshopentities.WorkOrderAllocation{}

	userDetailsUrl := config.EnvConfigs.GeneralServiceUrl + "user-detail/" + strconv.Itoa(req.TechnicianId)
	userDetailsPayload := transactionjpcbpayloads.UserDetailsPayload{}
	if err := utils.Get(userDetailsUrl, &userDetailsPayload, nil); err != nil || userDetailsPayload.UserEmployeeId == 0 {
		return entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("fail to retrieve technician data"),
		}
	}

	err := tx.Model(&entities).
		Where(transactionworkshopentities.WorkOrderAllocation{TechAllocSystemNumber: technicianAllocationSystemNumber}).
		First(&entities).
		Error
	if err != nil {
		return entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	entities.TechnicianId = req.TechnicianId
	entities.SequenceNumber = req.SequenceNumber

	err = tx.Save(&entities).Error
	if err != nil {
		return entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return entities, nil
}

func (r *JobAllocationRepositoryImpl) DeleteJobAllocation(tx *gorm.DB, technicianAllocationSystemNumber int) (bool, *exceptions.BaseErrorResponse) {
	entities := transactionworkshopentities.WorkOrderAllocation{}

	err := tx.Model(&entities).Where(transactionworkshopentities.WorkOrderAllocation{TechAllocSystemNumber: technicianAllocationSystemNumber}).First(&entities).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	entities.TechnicianId = 0
	entities.SequenceNumber = 0

	err = tx.Save(&entities).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return true, nil
}

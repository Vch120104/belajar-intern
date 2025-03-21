package transactionsparepartrepositoryimpl

import (
	masteritementities "after-sales/api/entities/master/item"
	masterwarehouseentities "after-sales/api/entities/master/warehouse"
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	"after-sales/api/exceptions"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	masterwarehouserepository "after-sales/api/repositories/master/warehouse"
	masterwarehouserepositoryimpl "after-sales/api/repositories/master/warehouse/repositories-warehouse-impl"
	transactionsparepartrepository "after-sales/api/repositories/transaction/sparepart"
	"after-sales/api/utils"
	generalserviceapiutils "after-sales/api/utils/general-service"
	"errors"
	"net/http"

	"gorm.io/gorm"
)

func NewItemWarehouseTransferRequestRepositoryImpl() transactionsparepartrepository.ItemWarehouseTransferRequestRepository {
	return &ItemWarehouseTransferRequestRepositoryImpl{}
}

type ItemWarehouseTransferRequestRepositoryImpl struct {
}

// GetByCodeTransferRequest implements transactionsparepartrepository.ItemWarehouseTransferRequestRepository.
func (*ItemWarehouseTransferRequestRepositoryImpl) GetByCodeTransferRequest(tx *gorm.DB, code string) (transactionsparepartpayloads.GetByIdItemWarehouseTransferRequestResponse, *exceptions.BaseErrorResponse) {
	var entities transactionsparepartentities.ItemWarehouseTransferRequest
	var warehouseFrom masterwarehouseentities.WarehouseMaster
	var warehouseGroupFrom masterwarehouseentities.WarehouseGroup
	var warehouseTo masterwarehouseentities.WarehouseMaster
	var warehouseGroupTo masterwarehouseentities.WarehouseGroup
	var response transactionsparepartpayloads.GetByIdItemWarehouseTransferRequestResponse

	errGetEntities := tx.Model(&entities).Where(transactionsparepartentities.ItemWarehouseTransferRequest{TransferRequestDocumentNumber: code}).First(&entities).Error
	if errGetEntities != nil {
		if errors.Is(errGetEntities, gorm.ErrRecordNotFound) {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "item claim with that code is not found please check input",
				Err:        errGetEntities,
			}
		}
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errGetEntities,
			Message:    "failed to get item claim please check input",
		}
	}

	errGetWarehouseFrom := tx.Model(&warehouseFrom).Where(masterwarehouseentities.WarehouseMaster{WarehouseId: entities.RequestFromWarehouseId}).First(&warehouseFrom).Error
	if errGetWarehouseFrom != nil {
		if errors.Is(errGetEntities, gorm.ErrRecordNotFound) {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "warehouse from with that id is not found please check input",
				Err:        errGetWarehouseFrom,
			}
		}
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errGetWarehouseFrom,
			Message:    "failed to get warehouse from please check input",
		}
	}

	errGetWarehouseGroupFrom := tx.Model(&warehouseGroupFrom).Where(masterwarehouseentities.WarehouseGroup{WarehouseGroupId: warehouseFrom.WarehouseGroupId}).First(&warehouseGroupFrom).Error
	if errGetWarehouseGroupFrom != nil {
		if errors.Is(errGetEntities, gorm.ErrRecordNotFound) {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "warehouse group from with that id is not found please check input",
				Err:        errGetWarehouseGroupFrom,
			}
		}
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errGetWarehouseGroupFrom,
			Message:    "failed to get warehouse group from please check input",
		}
	}

	errGetWarehouseTo := tx.Model(&warehouseTo).Where(masterwarehouseentities.WarehouseMaster{WarehouseId: entities.RequestToWarehouseId}).First(&warehouseTo).Error
	if errGetWarehouseTo != nil {
		if errors.Is(errGetWarehouseTo, gorm.ErrRecordNotFound) {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "warehouse to with that id is not found please check input",
				Err:        errGetWarehouseTo,
			}
		}
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errGetWarehouseTo,
			Message:    "failed to get warehouse to please check input",
		}
	}

	errGetWarehouseGroupTo := tx.Model(&warehouseGroupTo).Where(masterwarehouseentities.WarehouseGroup{WarehouseGroupId: warehouseTo.WarehouseGroupId}).First(&warehouseGroupTo).Error
	if errGetWarehouseGroupTo != nil {
		if errors.Is(errGetEntities, gorm.ErrRecordNotFound) {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "warehouse group to with that id is not found please check input",
				Err:        errGetWarehouseGroupTo,
			}
		}
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errGetWarehouseGroupTo,
			Message:    "failed to get warehouse group to please check input",
		}
	}

	var transferRequestStatus masteritementities.ItemTransferStatus
	errGetTransferStatus := tx.Model(&transferRequestStatus).Where(masteritementities.ItemTransferStatus{ItemTransferStatusId: entities.TransferRequestStatusId}).First(&transferRequestStatus).Error
	if errGetTransferStatus != nil {
		if errors.Is(errGetTransferStatus, gorm.ErrRecordNotFound) {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "transfer status with that id is not found please check input",
				Err:        errGetTransferStatus,
			}
		}
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errGetTransferStatus,
			Message:    "failed to get transfer status please check input",
		}
	}

	if *entities.TransferRequestById != 0 {
		getUser, errGetUser := generalserviceapiutils.GetUserDetailsByID(*entities.TransferRequestById)
		if errGetUser != nil {
			if errors.Is(errGetUser, gorm.ErrRecordNotFound) {
				return response, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "user with that id is not found please check input",
					Err:        errGetUser,
				}
			}
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        errGetUser,
				Message:    "failed to get user please check input",
			}
		}

		response.CreatedById = getUser.UserId
		response.CreatedByName = getUser.Username + " - " + getUser.EmployeeName
	}

	if entities.ModifiedById != 0 {
		getUserModified, errGetUserModified := generalserviceapiutils.GetUserDetailsByID(entities.ModifiedById)
		if errGetUserModified != nil {
			if errors.Is(errGetUserModified, gorm.ErrRecordNotFound) {
				return response, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "user with that id is not found please check input",
					Err:        errGetUserModified,
				}
			}
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        errGetUserModified,
				Message:    "failed to get user please check input",
			}
		}

		response.ModifiedById = getUserModified.UserId
		response.ModifiedByName = getUserModified.Username + " - " + getUserModified.EmployeeName
	}

	response.TransferRequestSystemNumber = entities.TransferRequestSystemNumber
	response.TransferRequestDocumentNumber = entities.TransferRequestDocumentNumber
	response.TransferRequestStatusId = entities.TransferRequestStatusId
	response.TransferRequestStatusCode = transferRequestStatus.ItemTransferStatusCode
	response.TransferRequestStatusDescription = transferRequestStatus.ItemTransferStatusDescription
	response.TransferRequestDate = *entities.TransferRequestDate
	response.RequestFromWarehouseId = entities.RequestFromWarehouseId
	response.RequestFromWarehouseCode = warehouseFrom.WarehouseCode
	response.RequestFromWarehouseName = warehouseFrom.WarehouseName
	response.RequestFromWarehouseGroupId = warehouseFrom.WarehouseGroupId
	response.RequestFromWarehouseGroupCode = warehouseGroupFrom.WarehouseGroupCode
	response.RequestFromWarehouseGroupName = warehouseGroupFrom.WarehouseGroupName
	response.RequestToWarehouseId = entities.RequestToWarehouseId
	response.RequestToWarehouseCode = warehouseTo.WarehouseCode
	response.RequestToWarehouseName = warehouseTo.WarehouseName
	response.RequestToWarehouseGroupId = warehouseTo.WarehouseGroupId
	response.RequestToWarehouseGroupCode = warehouseGroupTo.WarehouseGroupCode
	response.RequestToWarehouseGroupName = warehouseGroupTo.WarehouseGroupName
	response.RequestToWarehouseName = warehouseGroupTo.WarehouseGroupName
	response.Purpose = entities.Purpose
	response.ApprovalById = entities.ApprovalById
	response.ApprovalDate = entities.ApprovalDate
	response.ApprovalRemark = entities.ApprovalRemark

	return response, nil
}

// GetTransferRequestDetailLookUp implements transactionsparepartrepository.ItemWarehouseTransferRequestRepository.
func (*ItemWarehouseTransferRequestRepositoryImpl) GetTransferRequestDetailLookUp(tx *gorm.DB, number int, pages pagination.Pagination, filter []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var responses []transactionsparepartpayloads.GetAllItemWarehouseDetailLookUp

	joinTable := tx.Table("trx_item_warehouse_transfer_request_detail as det").
		Select(
			"transfer_request_detail_system_number",
			"det.transfer_request_system_number",
			"det.item_id",
			"it.item_code",
			"it.item_name",
			"uom.uom_code unit_of_measurement",
			"request_quantity",
		).
		Where("det.transfer_request_system_number = ?", number).
		Joins("INNER JOIN trx_item_warehouse_transfer_request req on req.transfer_request_system_number = det.transfer_request_system_number").
		Joins("LEFT JOIN mtr_item it on it.item_id = det.item_id").
		Joins("LEFT JOIN mtr_uom uom on uom.uom_id = it.unit_of_measurement_stock_id")

	whereQuery := utils.ApplyFilter(joinTable, filter)

	err := whereQuery.Scopes(pagination.Paginate(&pages, whereQuery)).Find(&responses).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	pages.Rows = responses
	return pages, nil
}

// GetTransferRequestLookUp implements transactionsparepartrepository.ItemWarehouseTransferRequestRepository.
func (*ItemWarehouseTransferRequestRepositoryImpl) GetTransferRequestLookUp(tx *gorm.DB, pages pagination.Pagination, filter []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	// var entities transactionsparepartentities.ItemWarehouseTransferRequest
	var responses []transactionsparepartpayloads.GetAllItemWarehouseLookUp
	var status masteritementities.ItemTransferStatus
	var status2 masteritementities.ItemTransferStatus

	errGetStatus := tx.Model(&status).Where("item_transfer_status_code = ?", 20).Find(&status)
	if errGetStatus.Error != nil {
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errGetStatus.Error,
		}
	}

	errGetStatus2 := tx.Model(&status2).Where("item_transfer_status_code = ?", 40).Find(&status2)
	if errGetStatus2.Error != nil {
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errGetStatus2.Error,
		}
	}

	joinTable := tx.Table("trx_item_warehouse_transfer_request as b").
		Select(
			"b.transfer_request_system_number",
			"b.transfer_request_document_number",
			"b.transfer_request_date",
			"b.transfer_request_by_id transfer_request_by_id",
			"b.request_from_warehouse_id",
			"wmf.warehouse_name request_from_warehouse_name",
			"wmf.warehouse_code request_from_warehouse_code",
			"wmf.warehouse_group_id request_from_warehouse_group_id",
			"wgf.warehouse_group_name request_from_warehouse_group_name",
			"wgf.warehouse_group_code request_from_warehouse_group_code",
		).
		Where("transfer_request_status_id in (? , ?)", status.ItemTransferStatusId, status2.ItemTransferStatusId).
		Joins("LEFT JOIN mtr_warehouse_master wmf on wmf.warehouse_id = b.request_from_warehouse_id").
		Joins("LEFT JOIN mtr_warehouse_group wgf on wgf.warehouse_group_id = wmf.warehouse_group_id")

	whereQuery := utils.ApplyFilter(joinTable, filter)

	err := whereQuery.Scopes(pagination.Paginate(&pages, whereQuery)).Find(&responses).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if len(responses) == 0 {
		pages.Rows = []map[string]interface{}{}
		return pages, nil
	}

	for i, respon := range responses {
		get, errUser := generalserviceapiutils.GetUserDetailsByID(respon.TransferRequestById)
		if errUser != nil {
			return pages, nil
		}
		responses[i].TransferRequestByName = get.EmployeeName
	}

	pages.Rows = responses
	return pages, nil
}

// GetByIdTransferRequestDetail implements transactionsparepartrepository.ItemWarehouseTransferRequestRepository.
func (*ItemWarehouseTransferRequestRepositoryImpl) GetByIdTransferRequestDetail(tx *gorm.DB, number int) (transactionsparepartpayloads.GetByIdItemWarehouseTransferRequestDetailResponse, *exceptions.BaseErrorResponse) {
	var entities transactionsparepartentities.ItemWarehouseTransferRequestDetail
	var response transactionsparepartpayloads.GetByIdItemWarehouseTransferRequestDetailResponse
	var item masteritementities.Item
	var uom masteritementities.Uom

	errGetEntities := tx.Model(&entities).Where(transactionsparepartentities.ItemWarehouseTransferRequestDetail{TransferRequestDetailSystemNumber: number}).First(&entities).Error
	if errGetEntities != nil {
		if errors.Is(errGetEntities, gorm.ErrRecordNotFound) {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "item claim with that id is not found please check input",
				Err:        errGetEntities,
			}
		}
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errGetEntities,
			Message:    "failed to get item claim please check input",
		}
	}

	errGetItem := tx.Model(&item).Where(masteritementities.Item{ItemId: *entities.ItemId}).First(&item).Error
	if errGetItem != nil {
		if errors.Is(errGetItem, gorm.ErrRecordNotFound) {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "item claim with that id is not found please check input",
				Err:        errGetItem,
			}
		}
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errGetItem,
			Message:    "failed to get item claim please check input",
		}
	}

	errGetUom := tx.Model(&uom).Where(masteritementities.Uom{UomId: *item.UnitOfMeasurementStockId}).First(&uom).Error
	if errGetUom != nil {
		if errors.Is(errGetUom, gorm.ErrRecordNotFound) {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "item claim with that id is not found please check input",
				Err:        errGetUom,
			}
		}
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errGetUom,
			Message:    "failed to get item claim please check input",
		}
	}

	response.ItemId = *entities.ItemId
	response.ItemCode = item.ItemCode
	response.StockUom = uom.UomCode
	response.Quantity = entities.RequestQuantity

	return response, nil
}

// UpdateWhTransferRequestDetail implements transactionsparepartrepository.ItemWarehouseTransferRequestRepository.
// Exec uspg_atTrfReq1_Update @option = 0
func (*ItemWarehouseTransferRequestRepositoryImpl) UpdateWhTransferRequestDetail(tx *gorm.DB, request transactionsparepartpayloads.UpdateItemWarehouseTransferRequestDetailRequest, number int) (transactionsparepartentities.ItemWarehouseTransferRequestDetail, *exceptions.BaseErrorResponse) {
	var entities transactionsparepartentities.ItemWarehouseTransferRequest
	var entitiesDetail transactionsparepartentities.ItemWarehouseTransferRequestDetail

	errDetail := tx.Model(&entitiesDetail).Where(transactionsparepartentities.ItemWarehouseTransferRequestDetail{TransferRequestDetailSystemNumber: number}).
		First(&entitiesDetail).Error
	if errDetail != nil {
		if errors.Is(errDetail, gorm.ErrRecordNotFound) {
			return entitiesDetail, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errDetail,
				Message:    "transfer request with that id is not found",
			}
		}
		return entitiesDetail, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errDetail,
			Message:    "failed to get transfer request entity",
		}
	}

	err := tx.Model(&entities).Where(transactionsparepartentities.ItemWarehouseTransferRequest{TransferRequestSystemNumber: entitiesDetail.TransferRequestSystemNumberId}).
		First(&entities).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entitiesDetail, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
				Message:    "transfer request with that id is not found",
			}
		}
		return entitiesDetail, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			Message:    "failed to get transfer request entity",
		}
	}

	getQuantity, errQuantity := masterwarehouserepository.NewLocationStockRepositoryImpl().GetAvailableQuantity(tx, masterwarehousepayloads.GetAvailableQuantityPayload{
		CompanyId:   entities.CompanyId,
		PeriodDate:  *entities.TransferRequestDate,
		WarehouseId: entities.RequestFromWarehouseId,
		ItemId:      *entitiesDetail.ItemId,
	})

	if errQuantity != nil {
		return transactionsparepartentities.ItemWarehouseTransferRequestDetail{}, errQuantity
	}

	if request.RequestQuantity > getQuantity.QuantityAvailable {
		return transactionsparepartentities.ItemWarehouseTransferRequestDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("qty for transfer request is not available"),
		}
	}

	if request.RequestQuantity != 0.0 {
		entitiesDetail.RequestQuantity = request.RequestQuantity
	}

	//save detail
	err = tx.Save(&entitiesDetail).Scan(&entitiesDetail).Error
	if err != nil {
		return entitiesDetail, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			Message:    "failed to update transfer request",
		}
	}

	//save header
	entities.ModifiedById = request.ModifiedById

	err = tx.Save(&entities).Error
	if err != nil {
		return entitiesDetail, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			Message:    "failed to update transfer request",
		}
	}

	return entitiesDetail, nil
}

// DeleteDetail implements transactionsparepartrepository.ItemWarehouseTransferRequestRepository.
// Exec uspg_atTrfReq1_Delete @option = 0
func (*ItemWarehouseTransferRequestRepositoryImpl) DeleteDetail(tx *gorm.DB, number int, request transactionsparepartpayloads.DeleteDetailItemWarehouseTransferRequest) (bool, *exceptions.BaseErrorResponse) {
	var entitiesDetail transactionsparepartentities.ItemWarehouseTransferRequestDetail
	var entities transactionsparepartentities.ItemWarehouseTransferRequest

	errGet := tx.Find(&entitiesDetail, number).Error
	if errGet != nil {
		if errors.Is(errGet, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "transfer detail with that id is not found please check input",
				Err:        errGet,
			}
		}
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errGet,
			Message:    "failed to get transfer detail please check input",
		}
	}

	errGetE := tx.Find(&entities, entitiesDetail.TransferRequestSystemNumberId).Scan(&entities).Error
	if errGetE != nil {
		if errors.Is(errGet, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "transfer request with that id is not found please check input",
				Err:        errGetE,
			}
		}
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errGetE,
			Message:    "failed to get transfer request please check input",
		}
	}

	errDeleteDetail := tx.Model(&entitiesDetail).Where("transfer_request_detail_system_number = ?", number).Delete(&entitiesDetail)
	if errDeleteDetail.Error != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errDeleteDetail.Error,
		}
	}

	//save header
	entities.ModifiedById = request.ModifiedById

	err := tx.Save(&entities).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			Message:    "failed to update transfer request",
		}
	}

	return true, nil
}

// DeleteHeaderTransferRequest implements transactionsparepartrepository.ItemWarehouseTransferRequestRepository.
// Exec uspg_atTrfReq0_Delete @option = 1
func (*ItemWarehouseTransferRequestRepositoryImpl) DeleteHeaderTransferRequest(tx *gorm.DB, number int) (bool, *exceptions.BaseErrorResponse) {
	var entitiesDetail transactionsparepartentities.ItemWarehouseTransferRequestDetail
	var entities transactionsparepartentities.ItemWarehouseTransferRequest

	errDeleteDetail := tx.Model(&entitiesDetail).Where("transfer_request_system_number = ?", number).Delete(&entitiesDetail)
	if errDeleteDetail.Error != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errDeleteDetail.Error,
		}
	}

	errDelete := tx.Model(&entities).Where("transfer_request_system_number = ?", number).Delete(&entities)
	if errDelete.Error != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errDelete.Error,
		}
	}

	return true, nil
}

// GetAllWhTransferRequest implements transactionsparepartrepository.ItemWarehouseTransferRequestRepository.
// Exec uspg_atTrfReq0_Select @option = 2
func (*ItemWarehouseTransferRequestRepositoryImpl) GetAllWhTransferRequest(tx *gorm.DB, pages pagination.Pagination, filter []utils.FilterCondition, dateParams map[string]string) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var entities transactionsparepartentities.ItemWarehouseTransferRequest
	var responses []transactionsparepartpayloads.GetAllItemWarehouseTransferRequestResponse

	joinTable := tx.Model(&entities).
		Select(
			"transfer_request_system_number",
			"transfer_request_document_number",
			"transfer_request_status_id",
			"item_transfer_status_code transfer_request_status_code",
			"item_transfer_status_description transfer_request_status_description",
			"transfer_request_date",
			"request_from_warehouse_id",
			"wmf.warehouse_name request_from_warehouse_name",
			"wmf.warehouse_group_id request_from_warehouse_group_id",
			"wgf.warehouse_group_name request_from_warehouse_group_name",
			"transfer_request_by_id",
			"request_to_warehouse_id",
			"wmt.warehouse_name request_to_warehouse_name",
			"wmt.warehouse_group_id request_to_warehouse_group_id",
			"wgt.warehouse_group_name request_to_warehouse_group_name",
		).
		Joins("LEFT JOIN mtr_warehouse_master wmf on wmf.warehouse_id = request_from_warehouse_id").
		Joins("LEFT JOIN mtr_item_transfer_status stat on stat.item_transfer_status_id = transfer_request_status_id").
		Joins("LEFT JOIN mtr_warehouse_master wmt on wmt.warehouse_id = request_to_warehouse_id").
		Joins("LEFT JOIN mtr_warehouse_group wgf on wgf.warehouse_group_id = wmf.warehouse_group_id").
		Joins("LEFT JOIN mtr_warehouse_group wgt on wgt.warehouse_group_id = wmt.warehouse_group_id")

	whereQuery := utils.ApplyFilter(joinTable, filter)

	var strDateFilter string
	if dateParams["transfer_request_date_from"] == "" {
		dateParams["transfer_request_date_from"] = "19000101"
	}
	if dateParams["transfer_request_date_to"] == "" {
		dateParams["transfer_request_date_to"] = "99991212"
	}
	strDateFilter = "transfer_request_date >='" + dateParams["transfer_request_date_from"] + "' AND transfer_request_date <= '" + dateParams["transfer_request_date_to"] + "'"

	whereQuery = whereQuery.Where(strDateFilter)

	err := whereQuery.Scopes(pagination.Paginate(&pages, whereQuery)).Find(&responses).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if len(responses) == 0 {
		pages.Rows = []map[string]interface{}{}
		return pages, nil
	}

	for i, respon := range responses {
		get, errUser := generalserviceapiutils.GetUserDetailsByID(respon.TransferRequestById)
		if errUser != nil {
			return pages, nil
		}
		responses[i].TransferRequestByName = get.EmployeeName
	}

	pages.Rows = responses
	return pages, nil
}

// GetAllTransferRequestById implements transactionsparepartrepository.ItemWarehouseTransferRequestRepository.
// Exec uspg_atTrfReq0_Select @option = 6
func (*ItemWarehouseTransferRequestRepositoryImpl) GetByIdTransferRequest(tx *gorm.DB, number int) (transactionsparepartpayloads.GetByIdItemWarehouseTransferRequestResponse, *exceptions.BaseErrorResponse) {
	var entities transactionsparepartentities.ItemWarehouseTransferRequest
	var warehouseFrom masterwarehouseentities.WarehouseMaster
	var warehouseGroupFrom masterwarehouseentities.WarehouseGroup
	var warehouseTo masterwarehouseentities.WarehouseMaster
	var warehouseGroupTo masterwarehouseentities.WarehouseGroup
	var response transactionsparepartpayloads.GetByIdItemWarehouseTransferRequestResponse

	errGetEntities := tx.Model(&entities).Where(transactionsparepartentities.ItemWarehouseTransferRequest{TransferRequestSystemNumber: number}).First(&entities).Error
	if errGetEntities != nil {
		if errors.Is(errGetEntities, gorm.ErrRecordNotFound) {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "item claim with that id is not found please check input",
				Err:        errGetEntities,
			}
		}
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errGetEntities,
			Message:    "failed to get item claim please check input",
		}
	}

	errGetWarehouseFrom := tx.Model(&warehouseFrom).Where(masterwarehouseentities.WarehouseMaster{WarehouseId: entities.RequestFromWarehouseId}).First(&warehouseFrom).Error
	if errGetWarehouseFrom != nil {
		if errors.Is(errGetEntities, gorm.ErrRecordNotFound) {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "warehouse from with that id is not found please check input",
				Err:        errGetWarehouseFrom,
			}
		}
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errGetWarehouseFrom,
			Message:    "failed to get warehouse from please check input",
		}
	}

	errGetWarehouseGroupFrom := tx.Model(&warehouseGroupFrom).Where(masterwarehouseentities.WarehouseGroup{WarehouseGroupId: warehouseFrom.WarehouseGroupId}).First(&warehouseGroupFrom).Error
	if errGetWarehouseGroupFrom != nil {
		if errors.Is(errGetEntities, gorm.ErrRecordNotFound) {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "warehouse group from with that id is not found please check input",
				Err:        errGetWarehouseGroupFrom,
			}
		}
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errGetWarehouseGroupFrom,
			Message:    "failed to get warehouse group from please check input",
		}
	}

	errGetWarehouseTo := tx.Model(&warehouseTo).Where(masterwarehouseentities.WarehouseMaster{WarehouseId: entities.RequestToWarehouseId}).First(&warehouseTo).Error
	if errGetWarehouseTo != nil {
		if errors.Is(errGetWarehouseTo, gorm.ErrRecordNotFound) {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "warehouse to with that id is not found please check input",
				Err:        errGetWarehouseTo,
			}
		}
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errGetWarehouseTo,
			Message:    "failed to get warehouse to please check input",
		}
	}

	errGetWarehouseGroupTo := tx.Model(&warehouseGroupTo).Where(masterwarehouseentities.WarehouseGroup{WarehouseGroupId: warehouseTo.WarehouseGroupId}).First(&warehouseGroupTo).Error
	if errGetWarehouseGroupTo != nil {
		if errors.Is(errGetEntities, gorm.ErrRecordNotFound) {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "warehouse group to with that id is not found please check input",
				Err:        errGetWarehouseGroupTo,
			}
		}
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errGetWarehouseGroupTo,
			Message:    "failed to get warehouse group to please check input",
		}
	}

	var transferRequestStatus masteritementities.ItemTransferStatus
	errGetTransferStatus := tx.Model(&transferRequestStatus).Where(masteritementities.ItemTransferStatus{ItemTransferStatusId: entities.TransferRequestStatusId}).First(&transferRequestStatus).Error
	if errGetTransferStatus != nil {
		if errors.Is(errGetTransferStatus, gorm.ErrRecordNotFound) {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "transfer status with that id is not found please check input",
				Err:        errGetTransferStatus,
			}
		}
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errGetTransferStatus,
			Message:    "failed to get transfer status please check input",
		}
	}

	getUser, errGetUser := generalserviceapiutils.GetUserDetailsByID(*entities.TransferRequestById)
	if errGetUser != nil {
		if errors.Is(errGetUser, gorm.ErrRecordNotFound) {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "user with that id is not found please check input",
				Err:        errGetUser,
			}
		}
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errGetUser,
			Message:    "failed to get user please check input",
		}
	}

	if entities.ModifiedById != 0 {
		getUserModified, errGetUserModified := generalserviceapiutils.GetUserDetailsByID(entities.ModifiedById)
		if errGetUserModified != nil {
			if errors.Is(errGetUserModified, gorm.ErrRecordNotFound) {
				return response, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "user with that id is not found please check input",
					Err:        errGetUserModified,
				}
			}
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        errGetUserModified,
				Message:    "failed to get user please check input",
			}
		}
		response.ModifiedById = getUserModified.UserId
		response.ModifiedByName = getUserModified.Username + " - " + getUserModified.EmployeeName
	}

	response.TransferRequestSystemNumber = number
	response.TransferRequestDocumentNumber = entities.TransferRequestDocumentNumber
	response.TransferRequestStatusId = entities.TransferRequestStatusId
	response.TransferRequestStatusCode = transferRequestStatus.ItemTransferStatusCode
	response.TransferRequestStatusDescription = transferRequestStatus.ItemTransferStatusDescription
	response.TransferRequestDate = *entities.TransferRequestDate
	response.RequestFromWarehouseId = entities.RequestFromWarehouseId
	response.RequestFromWarehouseCode = warehouseFrom.WarehouseCode
	response.RequestFromWarehouseName = warehouseFrom.WarehouseName
	response.RequestFromWarehouseGroupId = warehouseFrom.WarehouseGroupId
	response.RequestFromWarehouseGroupCode = warehouseGroupFrom.WarehouseGroupCode
	response.RequestFromWarehouseGroupName = warehouseGroupFrom.WarehouseGroupName
	response.RequestToWarehouseId = entities.RequestToWarehouseId
	response.RequestToWarehouseCode = warehouseTo.WarehouseCode
	response.RequestToWarehouseName = warehouseTo.WarehouseName
	response.RequestToWarehouseGroupId = warehouseTo.WarehouseGroupId
	response.RequestToWarehouseGroupCode = warehouseGroupTo.WarehouseGroupCode
	response.RequestToWarehouseGroupName = warehouseGroupTo.WarehouseGroupName
	response.RequestToWarehouseName = warehouseGroupTo.WarehouseGroupName
	response.Purpose = entities.Purpose
	response.ApprovalById = entities.ApprovalById
	response.ApprovalDate = entities.ApprovalDate
	response.ApprovalRemark = entities.ApprovalRemark
	response.CreatedById = getUser.UserId
	response.CreatedByName = getUser.Username + " - " + getUser.EmployeeName

	return response, nil
}

// GetAllDetailTransferRequest implements transactionsparepartrepository.ItemWarehouseTransferRequestRepository.
// Exec uspg_atTrfReq1_Select @option = 4
func (*ItemWarehouseTransferRequestRepositoryImpl) GetAllDetailTransferRequest(tx *gorm.DB, warehouseNumber int, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var entitiesDetail transactionsparepartentities.ItemWarehouseTransferRequestDetail
	var responses []transactionsparepartpayloads.GetAllDetailItemWarehouseTransferRequestResponse

	whereQuery := tx.Model(&entitiesDetail).
		Select(
			"transfer_request_detail_system_number",
			"transfer_request_system_number",
			"trx_item_warehouse_transfer_request_detail.item_id",
			"it.item_code",
			"it.item_name",
			"uom.uom_code unit_of_measurement",
			"request_quantity",
			"ISNULL(location_id_from, 0) location_id_from",
			"ISNULL(location_id_to, 0) location_id_to",
		).
		Joins("LEFT JOIN mtr_item it on it.item_id = trx_item_warehouse_transfer_request_detail.item_id").
		Joins("LEFT JOIN mtr_uom uom on uom.uom_id = it.unit_of_measurement_stock_id").
		Where(transactionsparepartentities.ItemWarehouseTransferRequestDetail{TransferRequestSystemNumberId: warehouseNumber})

	err := whereQuery.Scopes(pagination.Paginate(&pages, whereQuery)).Find(&responses).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if len(responses) == 0 {
		pages.Rows = []map[string]interface{}{}
		return pages, nil
	}

	pages.Rows = responses
	return pages, nil
}

// InsertWhTransferRequestDetail implements transactionsparepartrepository.ItemWarehouseTransferRequestRepository.
func (*ItemWarehouseTransferRequestRepositoryImpl) InsertWhTransferRequestDetail(tx *gorm.DB, request transactionsparepartpayloads.InsertItemWarehouseTransferDetailRequest) (transactionsparepartentities.ItemWarehouseTransferRequestDetail, *exceptions.BaseErrorResponse) {
	var entitiesDetail transactionsparepartentities.ItemWarehouseTransferRequestDetail
	var entities transactionsparepartentities.ItemWarehouseTransferRequest

	err := tx.Model(&entities).Where(transactionsparepartentities.ItemWarehouseTransferRequest{TransferRequestSystemNumber: request.TransferRequestSystemNumberId}).
		First(&entities).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entitiesDetail, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
				Message:    "transfer request with that id is not found",
			}
		}
		return entitiesDetail, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			Message:    "failed to get transfer request entity",
		}
	}

	getQuantity, errQuantity := masterwarehouserepository.NewLocationStockRepositoryImpl().GetAvailableQuantity(tx, masterwarehousepayloads.GetAvailableQuantityPayload{
		CompanyId:   entities.CompanyId,
		PeriodDate:  *entities.TransferRequestDate,
		WarehouseId: entities.RequestFromWarehouseId,
		ItemId:      *request.ItemId,
	})

	if errQuantity != nil {
		return transactionsparepartentities.ItemWarehouseTransferRequestDetail{}, errQuantity
	}

	if request.RequestQuantity > getQuantity.QuantityAvailable {
		return transactionsparepartentities.ItemWarehouseTransferRequestDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("qty for transfer request is not available"),
			Message:    "qty for transfer request is not available",
		}
	}

	entitiesDetail.ItemId = request.ItemId
	entitiesDetail.RequestQuantity = request.RequestQuantity
	entitiesDetail.TransferRequestSystemNumberId = request.TransferRequestSystemNumberId

	errDetail := tx.Create(&entitiesDetail).Scan(&entitiesDetail).Error

	if errDetail != nil {
		return transactionsparepartentities.ItemWarehouseTransferRequestDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errDetail,
		}
	}

	entities.ModifiedById = request.ModifiedById

	errSave := tx.Save(&entities).Error
	if errSave != nil {
		return transactionsparepartentities.ItemWarehouseTransferRequestDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			Message:    "failed to update transfer request",
		}
	}

	return entitiesDetail, nil
}

// InsertWhTransferRequestHeader implements transactionsparepartrepository.ItemWarehouseTransferRequestRepository.
// Exec uspg_atTrfReq0_Insert @option = 0
func (*ItemWarehouseTransferRequestRepositoryImpl) InsertWhTransferRequestHeader(tx *gorm.DB, request transactionsparepartpayloads.InsertItemWarehouseTransferRequest) (transactionsparepartentities.ItemWarehouseTransferRequest, *exceptions.BaseErrorResponse) {
	var entities transactionsparepartentities.ItemWarehouseTransferRequest

	var status masteritementities.ItemTransferStatus

	errGetStatus := tx.Model(&status).Where("item_transfer_status_code = ?", 10).Find(&status)
	if errGetStatus.Error != nil {
		return transactionsparepartentities.ItemWarehouseTransferRequest{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errGetStatus.Error,
		}
	}

	entities.CompanyId = request.CompanyId
	entities.TransferRequestStatusId = status.ItemTransferStatusId
	entities.RequestFromWarehouseId = request.RequestFromWarehouseId
	entities.RequestToWarehouseId = request.RequestToWarehouseId
	entities.TransferRequestDate = request.TransferRequestDate
	entities.Purpose = request.Purpose
	entities.TransferRequestById = request.TransferRequestById
	entities.ModifiedById = *request.TransferRequestById
	entities.ApprovalDate = nil

	errCreate := tx.Create(&entities).Scan(&entities).Error
	if errCreate != nil {
		return transactionsparepartentities.ItemWarehouseTransferRequest{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errCreate,
		}
	}

	return entities, nil
}

// SubmitWhTransferRequest implements transactionsparepartrepository.ItemWarehouseTransferRequestRepository.
func (*ItemWarehouseTransferRequestRepositoryImpl) SubmitWhTransferRequest(tx *gorm.DB, number int, request transactionsparepartpayloads.SubmitItemWarehouseTransferRequest) (transactionsparepartentities.ItemWarehouseTransferRequest, *exceptions.BaseErrorResponse) {
	var entities transactionsparepartentities.ItemWarehouseTransferRequest
	var status masteritementities.ItemTransferStatus
	var statusReady masteritementities.ItemTransferStatus

	errGetStatus := tx.Model(&status).Where("item_transfer_status_code = ?", 10).Find(&status)
	if errGetStatus.Error != nil {
		return transactionsparepartentities.ItemWarehouseTransferRequest{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errGetStatus.Error,
		}
	}

	errGetStatusReady := tx.Model(&statusReady).Where("item_transfer_status_code = ?", 15).Find(&statusReady)
	if errGetStatusReady.Error != nil {
		return transactionsparepartentities.ItemWarehouseTransferRequest{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errGetStatusReady.Error,
		}
	}

	err := tx.Model(&entities).Where(transactionsparepartentities.ItemWarehouseTransferRequest{TransferRequestSystemNumber: number}).
		First(&entities).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
				Message:    "transfer request with that id is not found",
			}
		}
		return entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			Message:    "failed to get transfer request entity",
		}
	}

	if entities.TransferRequestStatusId != status.ItemTransferStatusId {
		return transactionsparepartentities.ItemWarehouseTransferRequest{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("transfer request status is not draft"),
		}
	}

	// EXEC uspg_gmSrcDoc1_Update
	// @Option = 0 ,
	// @COMPANY_CODE = @Company_Code ,
	// @SOURCE_CODE = @Src_Code ,
	// @VEHICLE_BRAND = '' ,
	// @PROFIT_CENTER_CODE = '' ,
	// @TRANSACTION_CODE = '' ,
	// @BANK_ACC_CODE = '' ,
	// @TRANSACTION_DATE = @Trfreq_Date ,
	// @Last_Doc_No = @Trfreq_Doc_No OUTPUT

	// getCompany, errComp := generalserviceapiutils.GetCompanyDataById(entities.CompanyId)
	// if errComp != nil {
	// 	return transactionsparepartentities.ItemWarehouseTransferRequest{}, errComp
	// }

	// ss := strconv.Itoa(entities.TransferRequestDate.Year()) + strconv.Itoa(int(entities.TransferRequestDate.Month())) + "01"
	// // ss := strconv.Itoa(entities.TransferRequestDate.Year()) + strconv.Itoa(int(entities.TransferRequestDate.Month())) + "01"
	// var checkEntities []transactionsparepartentities.ItemWarehouseTransferRequest
	// errGetCheck := tx.Where("transfer_request_date >= ?", ss).
	// 	Find(&checkEntities)
	// if errGetCheck.Error != nil {
	// 	return transactionsparepartentities.ItemWarehouseTransferRequest{}, &exceptions.BaseErrorResponse{
	// 		StatusCode: http.StatusInternalServerError,
	// 		Err:        errGetCheck.Error,
	// 	}
	// }

	// month := entities.TransferRequestDate.Month()
	// year := entities.TransferRequestDate.Year() % 100

	// docNum := fmt.Sprintf("SPTR/%s/%02d/%02d/%05d", getCompany.CompanyCode, month, year, len(checkEntities))

	warehouseFrom, errFrom := masterwarehouserepositoryimpl.OpenWarehouseMasterImpl().GetById(tx, entities.RequestFromWarehouseId)
	if errFrom != nil {
		return transactionsparepartentities.ItemWarehouseTransferRequest{}, errFrom
	}

	warehouseTo, errTo := masterwarehouserepositoryimpl.OpenWarehouseMasterImpl().GetById(tx, entities.RequestToWarehouseId)
	if errTo != nil {
		return transactionsparepartentities.ItemWarehouseTransferRequest{}, errTo
	}

	if warehouseFrom.WarehouseCostingTypeId != warehouseTo.WarehouseCostingTypeId {
		return transactionsparepartentities.ItemWarehouseTransferRequest{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("costing type for source warehouse is different than destination warehouse"),
		}
	}

	check := 0
	errGet := tx.Model(&transactionsparepartentities.ItemWarehouseTransferRequestDetail{}).
		Select(
			"TOP 1 1",
		).Where("transfer_request_system_number = ?", number).Group("item_id").Having("count(item_id) > 1").Scan(&check)

	if errGet.Error != nil {
		return transactionsparepartentities.ItemWarehouseTransferRequest{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errGet.Error,
		}
	}

	if check != 0 {
		return transactionsparepartentities.ItemWarehouseTransferRequest{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("duplicate item exist"),
		}
	}

	entities.TransferRequestStatusId = statusReady.ItemTransferStatusId
	// entities.TransferRequestDocumentNumber =

	entities.ModifiedById = request.ModifiedById
	// entities.TransferRequestDocumentNumber = docNum

	err = tx.Save(&entities).Error
	if err != nil {
		return entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			Message:    "failed to update transfer request",
		}
	}

	return entities, nil
}

// UpdateWhTransferRequest implements transactionsparepartrepository.ItemWarehouseTransferRequestRepository.
// Exec uspg_atTrfReq0_Update @option = 3
func (*ItemWarehouseTransferRequestRepositoryImpl) UpdateWhTransferRequest(tx *gorm.DB, request transactionsparepartpayloads.UpdateItemWarehouseTransferRequest, tranferId int) (transactionsparepartentities.ItemWarehouseTransferRequest, *exceptions.BaseErrorResponse) {
	var entities transactionsparepartentities.ItemWarehouseTransferRequest

	err := tx.Model(&entities).Where(transactionsparepartentities.ItemWarehouseTransferRequest{TransferRequestSystemNumber: tranferId}).
		First(&entities).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
				Message:    "transfer request with that id is not found",
			}
		}
		return entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			Message:    "failed to get transfer request entity",
		}
	}

	if request.RequestFromWarehouseId != nil {
		entities.RequestFromWarehouseId = *request.RequestFromWarehouseId
	}
	if request.RequestToWarehouseId != nil {
		entities.RequestToWarehouseId = *request.RequestToWarehouseId
	}
	if request.Purpose != "" {
		entities.Purpose = request.Purpose
	}
	if request.TransferRequestById != nil {
		entities.TransferRequestById = request.TransferRequestById
	}

	//save header
	entities.ModifiedById = request.ModifiedById

	err = tx.Save(&entities).Error
	if err != nil {
		return entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			Message:    "failed to update transfer request",
		}
	}

	return entities, nil
}

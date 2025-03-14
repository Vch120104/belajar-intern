package transactionsparepartrepositoryimpl

import (
	masterentities "after-sales/api/entities/master"
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartrepository "after-sales/api/repositories/transaction/sparepart"
	"after-sales/api/utils"
	"errors"
	"fmt"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type StockOpnameRepositoryImpl struct {
}

func NewStockOpnameRepositoryImpl() transactionsparepartrepository.StockOpnameRepository {
	return &StockOpnameRepositoryImpl{}
}

func (r *StockOpnameRepositoryImpl) GetAllStockOpname(
	tx *gorm.DB, filteredCondition []utils.FilterCondition, pages pagination.Pagination,
	dateParams map[string]interface{}) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var entities transactionsparepartentities.StockOpname
	var responses []transactionsparepartpayloads.GetAllStockOpnameResponse

	query := tx.Model(&entities).
		Select(`trx_stock_opname.company_id ,
		trx_stock_opname.stock_opname_system_number ,
		trx_stock_opname.stock_opname_document_number,
		trx_stock_opname.stock_opname_status_id as stock_opname_status ,
		B.warehouse_group_name,
		C.warehouse_name,
		trx_stock_opname.location_range_from_id ,
		trx_stock_opname.location_range_to_id ,
		trx_stock_opname.show_detail,
		trx_stock_opname.person_in_charge_id,
		trx_stock_opname.remark ,
		trx_stock_opname.execution_date_from,
		trx_stock_opname.execution_date_to,
		trx_stock_opname.include_zero_onhand`).
		Joins("LEFT OUTER JOIN dbo.mtr_warehouse_group B ON	B.warehouse_group_id = trx_stock_opname.warehouse_group_id").
		Joins("LEFT OUTER JOIN dbo.mtr_warehouse_master C ON C.warehouse_id = trx_stock_opname.warehouse_id AND C.company_id = trx_stock_opname.company_id")

	if query.Error != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: 500,
			Message:    "error when selecting data",
			Err:        query.Error,
		}
	}

	whereQ := utils.ApplyFilter(query, filteredCondition)

	for key, value := range dateParams {
		whereQ = whereQ.Where(key, value)
	}

	paginatedQuery := whereQ.Scopes(pagination.Paginate(&pages, query))
	err := paginatedQuery.Scan(&responses).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: 500,
			Message:    "error when scanning data",
			Err:        err,
		}
	}

	pages.Rows = responses
	return pages, nil

}

func (r *StockOpnameRepositoryImpl) GetAllStockOpnameDetail(tx *gorm.DB, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var entities transactionsparepartentities.StockOpnameDetail
	var responses []transactionsparepartpayloads.GetAllStockOpnameDetailResponse

	query := tx.Model(&entities).
		Select(`trx_stock_opname_detail.stock_opname_detail_system_number,
	trx_stock_opname_detail.stock_opname_system_number,
	trx_stock_opname_detail.stock_opname_line,
	trx_stock_opname_detail.item_id,
	b.item_name,
	trx_stock_opname_detail.location_id,
	trx_stock_opname_detail.system_quantity,
	trx_stock_opname_detail.found_quantity,
	trx_stock_opname_detail.broken_quantity,
	trx_stock_opname_detail.need_adjustment,
	trx_stock_opname_detail.remark`).
		Joins("left outer join mtr_item b on b.item_id = trx_stock_opname_detail.item_id")

	if query.Error != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: 500,
			Message:    "error when selecting data",
			Err:        query.Error,
		}
	}

	paginatedQuery := query.Scopes(pagination.Paginate(&pages, query))
	err := paginatedQuery.Scan(&responses).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: 500,
			Message:    "error when scanning data",
			Err:        err,
		}
	}
	pages.Rows = responses
	return pages, nil
}

func (r *StockOpnameRepositoryImpl) GetStockOpnameByStockOpnameSystemNumber(tx *gorm.DB, stockOpnameSystemNumber int) (
	[]transactionsparepartpayloads.GetStockOpnameByStockOpnameSystemNumberResponse, *exceptions.BaseErrorResponse) {
	var entities transactionsparepartentities.StockOpname
	var responses []transactionsparepartpayloads.GetStockOpnameByStockOpnameSystemNumberResponse

	query := tx.Model(&entities).
		Select(`trx_stock_opname.company_id ,
	trx_stock_opname.stock_opname_system_number ,
	trx_stock_opname.stock_opname_document_number,
	trx_stock_opname.stock_opname_status_id ,
	trx_stock_opname.item_group_id as item_group,
	B.warehouse_group_name,
	C.warehouse_name,
	trx_stock_opname.location_range_from_id ,
	trx_stock_opname.location_range_to_id ,
	trx_stock_opname.show_detail,
	trx_stock_opname.person_in_charge_id as person_in_charge,
	trx_stock_opname.remark ,
	trx_stock_opname.execution_date_from,
	trx_stock_opname.execution_date_to,
	trx_stock_opname.include_zero_onhand`).
		Joins("LEFT OUTER JOIN dbo.mtr_warehouse_group B ON	B.warehouse_group_id = trx_stock_opname.warehouse_group_id").
		Joins("LEFT OUTER JOIN dbo.mtr_warehouse_master C ON C.warehouse_id = trx_stock_opname.warehouse_id AND C.company_id = trx_stock_opname.company_id").
		Where("trx_stock_opname.stock_opname_system_number = ?", stockOpnameSystemNumber)

	if query.Error != nil {
		return responses, &exceptions.BaseErrorResponse{
			StatusCode: 500,
			Message:    "error when selecting data",
			Err:        query.Error,
		}
	}

	err := query.Scan(&responses).Error
	if err != nil {
		return responses, &exceptions.BaseErrorResponse{
			StatusCode: 500,
			Message:    "error when scanning data",
			Err:        err,
		}
	}

	return responses, nil
}

func (r *StockOpnameRepositoryImpl) GetStockOpnameAllDetailByStockOpnameSystemNumber(tx *gorm.DB, stockOpnameSystemNumber int,
	pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var entities transactionsparepartentities.StockOpnameDetail
	var responses []transactionsparepartpayloads.GetAllStockOpnameDetailResponse

	query := tx.Model(&entities).
		Select(`trx_stock_opname_detail.stock_opname_detail_system_number,
	trx_stock_opname_detail.stock_opname_system_number,
	trx_stock_opname_detail.stock_opname_line,
	trx_stock_opname_detail.item_id,
	b.item_name,
	trx_stock_opname_detail.location_id as location,
	trx_stock_opname_detail.system_quantity,
	trx_stock_opname_detail.found_quantity,
	trx_stock_opname_detail.broken_quantity,
	trx_stock_opname_detail.need_adjustment,
	trx_stock_opname_detail.remark`).
		Joins("left outer join mtr_item b on b.item_id = trx_stock_opname_detail.item_id").
		Where("trx_stock_opname_detail.stock_opname_system_number = ?", stockOpnameSystemNumber)

	if query.Error != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: 500,
			Message:    "error when selecting data",
			Err:        query.Error,
		}
	}

	paginatedQuery := query.Scopes(pagination.Paginate(&pages, query))
	err := paginatedQuery.Scan(&responses).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: 500,
			Message:    "error when scanning data",
			Err:        err,
		}
	}
	pages.Rows = responses
	return pages, nil
}

func (r *StockOpnameRepositoryImpl) InsertStockOpname(tx *gorm.DB,
	request transactionsparepartpayloads.StockOpnameInsertRequest) (bool, *exceptions.BaseErrorResponse) {
	var entities transactionsparepartentities.StockOpname
	var status masterentities.StockOpnameStatus

	getStatus := tx.Model(&status).Where("stock_opname_status_code = ?", 01).First(&status)
	if getStatus.Error != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: 500,
			Message:    "error when selecting data",
			Err:        getStatus.Error,
		}
	}

	fmt.Println("status ", status.StockOpnameStatusId)

	entities.CompanyID = &request.CompanyId
	entities.StockOpnameDocumentNumber = request.StockOpnameDocumentNumber
	entities.StockOpnameStatusId = &status.StockOpnameStatusId
	entities.WarehouseGroupId = &request.WarehouseGroup
	entities.WarehouseId = &request.WarehouseCode
	entities.LocationRangeFromId = &request.FromLocation
	entities.LocationRangeToId = &request.ToLocation
	entities.ItemGroupId = &request.ItemGroup
	entities.ShowDetail = request.ShowDetail
	entities.PersonInChargeId = &request.PersonInCharge
	entities.Remark = request.Remark
	entities.ExecutionDateFrom = &request.ExecutionDateFrom
	entities.ExecutionDateTo = &request.ExecutionDateTo
	entities.IncludeZeroOnhand = request.IncludeZeroOnhand
	entities.ApprovalRequestedById = &request.ApprovalRequestedById
	entities.ApprovalRequestedDate = nil
	entities.ApprovalById = nil
	entities.ApprovalDate = nil

	insert := tx.Create(&entities)

	if insert.Error != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: 500,
			Message:    "error when inserting data",
			Err:        insert.Error,
		}
	}
	return true, nil
}

func (r *StockOpnameRepositoryImpl) SubmitStockOpname(tx *gorm.DB, systemNumber int, request transactionsparepartpayloads.StockOpnameSubmitRequest) (bool, *exceptions.BaseErrorResponse) {
	var entities transactionsparepartentities.StockOpname
	var status masterentities.StockOpnameStatus
	var approvalStatus masterentities.StockOpnameStatus

	getStatus := tx.Model(&status).Where("stock_opname_status_code = ?", 01).First(&status)
	if getStatus.Error != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: 500,
			Message:    "error when selecting data",
			Err:        getStatus.Error,
		}
	}

	getApprovalStatus := tx.Model(&approvalStatus).Where("stock_opname_status_code = ?", 02).First(&approvalStatus)
	if getApprovalStatus.Error != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: 500,
			Message:    "error when selecting data",
			Err:        getApprovalStatus.Error,
		}
	}

	query := tx.Model(&entities).Where(transactionsparepartentities.StockOpname{StockOpnameSystemNumber: systemNumber}).First(&entities)

	if query.Error != nil {
		if errors.Is(query.Error, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        query.Error,
				Message:    "stock opname with that id is not found",
			}
		}
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        query.Error,
			Message:    "failed to get stock opname entity",
		}
	}

	entities.StockOpnameStatusId = &approvalStatus.StockOpnameStatusId
	entities.ApprovalRequestedById = &request.StockOpnameApprovalRequestId
	now := time.Now()
	entities.ApprovalRequestedDate = &now

	save := tx.Save(&entities)
	if save.Error != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: 500,
			Message:    "error when updating data",
			Err:        save.Error,
		}
	}
	return true, nil
}

func (r *StockOpnameRepositoryImpl) InsertStockOpnameDetail(tx *gorm.DB,
	request transactionsparepartpayloads.StockOpnameInsertDetailRequest) (bool, *exceptions.BaseErrorResponse) {
	var entities transactionsparepartentities.StockOpname

	query := tx.Model(&entities).Where(transactionsparepartentities.StockOpname{StockOpnameSystemNumber: request.StockOpnameSystemNumber}).First(&entities)
	if query.Error != nil {
		if errors.Is(query.Error, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        query.Error,
				Message:    "stock opname with that id is not found",
			}
		}
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        query.Error,
			Message:    "failed to get stock opname entity",
		}
	}

	lineNumber := 1

	for _, itemId := range request.ItemId {
		detailEntities := transactionsparepartentities.StockOpnameDetail{
			StockOpnameSystemNumber: entities.StockOpnameSystemNumber,
			ItemId:                  &itemId,
			StockOpnameLine:         lineNumber,
			WarehouseId:             entities.WarehouseId,
		}

		insertQuery := tx.Create(&detailEntities)
		if insertQuery.Error != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: 500,
				Message:    "error when inserting data",
				Err:        insertQuery.Error,
			}
		}
		lineNumber++
	}
	return true, nil
}

func (r *StockOpnameRepositoryImpl) UpdateStockOpname(tx *gorm.DB,
	request transactionsparepartpayloads.StockOpnameInsertRequest, systemNumber int) (bool, *exceptions.BaseErrorResponse) {
	var entities transactionsparepartentities.StockOpname
	var status masterentities.StockOpnameStatus

	getSystemNumber := tx.Model(&entities).Where("stock_opname_system_number = ?", systemNumber).First(&entities)
	if getSystemNumber.Error != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: 500,
			Message:    "error when selecting data",
			Err:        getSystemNumber.Error,
		}
	}

	getStatus := tx.Model(&status).Where("stock_opname_status_code = ?", 01).First(&status)
	if getStatus.Error != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: 500,
			Message:    "error when selecting data",
			Err:        getStatus.Error,
		}
	}

	entities.CompanyID = &request.CompanyId
	entities.StockOpnameStatusId = &status.StockOpnameStatusId
	entities.StockOpnameDocumentNumber = request.StockOpnameDocumentNumber
	entities.WarehouseGroupId = &request.WarehouseGroup
	entities.WarehouseId = &request.WarehouseCode
	entities.LocationRangeFromId = &request.FromLocation
	entities.LocationRangeToId = &request.ToLocation
	entities.ItemGroupId = &request.ItemGroup
	entities.ShowDetail = request.ShowDetail
	entities.PersonInChargeId = &request.PersonInCharge
	entities.Remark = request.Remark
	entities.ExecutionDateFrom = &request.ExecutionDateFrom
	entities.ExecutionDateTo = &request.ExecutionDateTo
	entities.IncludeZeroOnhand = request.IncludeZeroOnhand
	entities.ApprovalRequestedById = &request.ApprovalRequestedById
	entities.ApprovalRequestedDate = nil
	entities.ApprovalById = nil
	entities.ApprovalDate = nil

	updateQuery := tx.Save(&entities)

	if updateQuery.Error != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: 500,
			Message:    "error when inserting data",
			Err:        updateQuery.Error,
		}
	}
	return true, nil
}

func (r *StockOpnameRepositoryImpl) UpdateStockOpnameDetail(tx *gorm.DB,
	request transactionsparepartpayloads.StockOpnameUpdateDetailRequest, systemNumber int) (bool, *exceptions.BaseErrorResponse) {
	var detailEntities transactionsparepartentities.StockOpnameDetail
	var entities transactionsparepartentities.StockOpname

	query := tx.Model(&entities).Where(transactionsparepartentities.StockOpnameDetail{
		StockOpnameSystemNumber: systemNumber,
		StockOpnameLine:         request.StockOpnameLine,
	}).First(&entities)
	if query.Error != nil {
		if errors.Is(query.Error, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        query.Error,
				Message:    "stock opname with that id is not found",
			}
		}
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        query.Error,
			Message:    "failed to get stock opname entity",
		}
	}

	detailEntities.ItemId = &request.ItemId

	insertQuery := tx.Save(&detailEntities)
	if insertQuery.Error != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: 500,
			Message:    "error when inserting data",
			Err:        insertQuery.Error,
		}
	}
	return true, nil
}

func (r *StockOpnameRepositoryImpl) DeleteStockOpname(tx *gorm.DB, systemNumber int) (bool, *exceptions.BaseErrorResponse) {
	var entities transactionsparepartentities.StockOpname
	var status masterentities.StockOpnameStatus

	getStatus := tx.Model(&status).Where("stock_opname_status_code = ?", 01).First(&status)
	if getStatus.Error != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: 500,
			Message:    "error when selecting data",
			Err:        getStatus.Error,
		}
	}

	query := tx.Model(&entities).Where("stock_opname_system_number = ? and stock_opname_status_id = ?", systemNumber,
		&status.StockOpnameStatusId).First(&entities)
	if query.Error != nil {
		if errors.Is(query.Error, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        query.Error,
				Message:    "stock opname with that id is not found",
			}
		}
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        query.Error,
			Message:    "failed to get stock opname entity",
		}
	}

	deleteQuery := tx.Delete(&entities)
	if deleteQuery.Error != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: 500,
			Message:    "error when deleting data",
			Err:        deleteQuery.Error,
		}
	}
	return true, nil
}

func (r *StockOpnameRepositoryImpl) DeleteStockOpnameDetailByLineNumber(tx *gorm.DB, systemNumber int, lineNumber int) (bool, *exceptions.BaseErrorResponse) {
	var entities transactionsparepartentities.StockOpnameDetail
	var entitiesBySystemNumber []transactionsparepartentities.StockOpnameDetail

	getSystemNumber := tx.Model(&entities).Where(transactionsparepartentities.StockOpnameDetail{
		StockOpnameSystemNumber: systemNumber,
		StockOpnameLine:         lineNumber,
	}).First(&entities)

	if getSystemNumber.Error != nil {
		if errors.Is(getSystemNumber.Error, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        getSystemNumber.Error,
				Message:    "stock opname with that id is not found",
			}
		}
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        getSystemNumber.Error,
			Message:    "failed to get stock opname entity",
		}
	}

	deleteQuery := tx.Delete(&entities)
	if deleteQuery.Error != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: 500,
			Message:    "error when deleting data",
			Err:        deleteQuery.Error,
		}
	}

	getEntitiesBySystemNumber := tx.Model(&entitiesBySystemNumber).Where(transactionsparepartentities.StockOpnameDetail{
		StockOpnameSystemNumber: systemNumber,
	}).Order("stock_opname_line asc").Scan(&entitiesBySystemNumber)

	if getEntitiesBySystemNumber.Error != nil {
		if errors.Is(getEntitiesBySystemNumber.Error, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        getEntitiesBySystemNumber.Error,
				Message:    "stock opname with that id is not found",
			}
		}
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        getEntitiesBySystemNumber.Error,
			Message:    "failed to get stock opname entity",
		}
	}

	newLine := 1
	for _, detail := range entitiesBySystemNumber {
		if detail.StockOpnameLine != newLine {
			detail.StockOpnameLine = newLine
			updateQuery := tx.Save(&detail)

			if updateQuery.Error != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: 500,
					Message:    "error when updating data",
					Err:        updateQuery.Error,
				}
			}
		}
		newLine++
	}

	return true, nil
}

func (r *StockOpnameRepositoryImpl) DeleteStockOpnameDetailBySystemNumber(tx *gorm.DB, systemNumber int) (bool, *exceptions.BaseErrorResponse) {
	var entities transactionsparepartentities.StockOpnameDetail

	deleteQuery := tx.Model(&entities).Where("stock_opname_system_number = ?", systemNumber).Delete(&entities)
	if deleteQuery.Error != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: 500,
			Message:    "error when deleting data",
			Err:        deleteQuery.Error,
		}
	}
	return true, nil
}

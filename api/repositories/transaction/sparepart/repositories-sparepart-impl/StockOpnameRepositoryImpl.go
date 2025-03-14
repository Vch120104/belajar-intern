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
	b.item_type_name,
	trx_stock_opname_detail.location_id,
	trx_stock_opname_detail.system_quantity,
	trx_stock_opname_detail.found_quantity,
	trx_stock_opname_detail.broken_quantity,
	trx_stock_opname_detail.need_adjustment,
	trx_stock_opname_detail.remark`).
		Joins("left outer join mtr_item_type b on b.item_type_id = trx_stock_opname_detail.item_id").
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

	// if entities.StockOpnameStatusId != &status.StockOpnameStatusId {
	// 	return false, &exceptions.BaseErrorResponse{
	// 		StatusCode: http.StatusBadRequest,
	// 		Message:    "stock opname status is not draft",
	// 	}
	// }

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
	}).Scan(&entities)

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
			updateQuery := tx.Model(&entities).Where(transactionsparepartentities.StockOpnameDetail{
				StockOpnameDetailSystemNumber: detail.StockOpnameDetailSystemNumber,
			}).Update("stock_opname_line", newLine)
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
	panic("implement me")
}

// func (r *StockOpnameRepositoryImpl) GetItemList(tx *gorm.DB, pages pagination.Pagination, whsCode string, itemGroup string) (pagination.Pagination, *exceptions.BaseErrorResponse) {
// 	var datas []transactionsparepartpayloads.GetItemListResponse

// 	query := tx.Model(&transactionsparepartentities.AtStockOpname1{}).
// 		Select(`atstockopname1.record_status,
// 	atstockopname1.stock_opname_sys_no,
// 	atstockopname1.stock_opname_line,
// 	atstockopname1.line_status,
// 	atstockopname1.item_code,
// 	b.item_name,
// 	atstockopname1.loc_code as location,
// 	atstockopname1.sys_qty,
// 	atstockopname1.found_qty,
// 	atstockopname1.broke_qty,
// 	atstockopname1.broken_loc_code,
// 	atstockopname1.need_adjustment,
// 	atstockopname1.remark,
// 	atstockopname1.change_no,
// 	atstockopname1.creation_user_id,
// 	atstockopname1.creation_datetime,
// 	atstockopname1.change_user_id,
// 	atstockopname1.change_datetime`).
// 		Joins("left outer join gmitem0 b on b.item_code = atstockopname1.item_code").
// 		Where("atstockopname1.whs_code = ? and b.item_group = ?", whsCode, itemGroup)

// 	if query.Error != nil {
// 		return pages, &exceptions.BaseErrorResponse{
// 			StatusCode: 500,
// 			Message:    "error when selecting data",
// 			Err:        query.Error,
// 		}
// 	}

// 	paginatedQ := query.Scopes(pagination.Paginate(&pages, query)).Order("atstockopname1.loc_code")
// 	err := paginatedQ.Scan(&datas).Error
// 	if err != nil {
// 		return pages, &exceptions.BaseErrorResponse{
// 			StatusCode: 500,
// 			Message:    "error when scanning data",
// 			Err:        err,
// 		}
// 	}
// 	pages.Rows = datas
// 	return pages, nil

// }

// func (r *StockOpnameRepositoryImpl) GetAllStockOpname(tx *gorm.DB, filteredCondition []utils.FilterCondition, pages pagination.Pagination, companyCode float64, dateParams map[string]interface{}) (pagination.Pagination, *exceptions.BaseErrorResponse) {
// 	var datas []transactionsparepartpayloads.GetAllStockOpnameResponse

// 	query := tx.Model(&transactionsparepartentities.AtStockOpname0{}).
// 		Select(`atstockopname0.COMPANY_CODE ,
// 	atstockopname0.STOCK_OPNAME_SYS_NO ,
// 	atstockopname0.STOCK_OPNAME_DOC_NO as stock_opname_no,
// 	atstockopname0.STOCK_OPNAME_STATUS ,
// 	B.DESCRIPTION as warehouse_group,
// 	C.WAREHOUSE_NAME as warehouse_name,
// 	atstockopname0.LOC_RANGE_FROM ,
// 	atstockopname0.LOC_RANGE_TO ,
// 	atstockopname0.SHOW_DETAIL ,
// 	CASE atstockopname0.STOCK_OPNAME_STATUS
// 		WHEN ?  THEN  ?
// 		WHEN ? 	THEN  ?
// 		WHEN  ?  THEN  ?
// 		WHEN  ?  THEN  ?
// 	END AS status,
// 	atstockopname0.PIC ,
// 	atstockopname0.REMARK ,
// 	atstockopname0.EXEC_DATE_FROM as stock_opname_from,
// 	atstockopname0.EXEC_DATE_TO as stock_opname_to,
// 	atstockopname0.Include_Zero_Onhand`,
// 			"01", "Draft", "02", "In Progress", "15", "Wait Approve", "20", "Approved").
// 		Joins("LEFT OUTER JOIN dbo.gmLoc0 B ON	B.WAREHOUSE_GROUP = atstockopname0.WHS_GROUP").
// 		Joins("LEFT OUTER JOIN dbo.gmLoc1 C ON	C.WAREHOUSE_CODE = atstockopname0.WHS_CODE AND C.COMPANY_CODE = atstockopname0.COMPANY_CODE").
// 		Where("atstockopname0.record_status = ? AND atstockopname0.company_code = ?", "A", companyCode)

// 	if query.Error != nil {
// 		return pages, &exceptions.BaseErrorResponse{
// 			StatusCode: 500,
// 			Message:    "error when selecting data",
// 			Err:        query.Error,
// 		}
// 	}
// 	whereQ := utils.ApplyFilter(query, filteredCondition)

// 	for key, value := range dateParams {
// 		whereQ = whereQ.Where(key, value)
// 	}

// 	paginatedQ := whereQ.Scopes(pagination.Paginate(&pages, query))

// 	err := paginatedQ.Scan(&datas).Error
// 	if err != nil {
// 		return pages, &exceptions.BaseErrorResponse{
// 			StatusCode: 500,
// 			Message:    "error when scanning data",
// 			Err:        err,
// 		}
// 	}

// 	pages.Rows = datas
// 	fmt.Printf("Retrieved data: %+v\n", datas)
// 	return pages, nil

// }

// func (r *StockOpnameRepositoryImpl) GetLocationList(tx *gorm.DB, filteredCondition []utils.FilterCondition, pages pagination.Pagination,
// 	companyCode float64, warehouseGroup string, warehouseCode string) (pagination.Pagination, *exceptions.BaseErrorResponse) {
// 	var datas []transactionsparepartpayloads.GetAllLocationList

// 	query := tx.Model(&transactionsparepartentities.GmLoc2{}).
// 		Select(`LOCATION_CODE as location_code, LOCATION_NAME as location_name,
// 			CASE record_status
// 			WHEN  ? THEN ?
// 			WHEN ? THEN ?
// 			END as status`, "A", "Active", "D", "Deactive").
// 		Where("company_code = ? and warehouse_group = ? and warehouse_code = ?", companyCode, warehouseGroup, warehouseCode)

// 	if query.Error != nil {
// 		return pages, &exceptions.BaseErrorResponse{
// 			StatusCode: 500,
// 			Message:    "error when selecting data",
// 			Err:        query.Error,
// 		}
// 	}

// 	whereQ := utils.ApplyFilter(query, filteredCondition)
// 	paginatedQ := whereQ.Scopes(pagination.Paginate(&pages, whereQ)).Order("LOCATION_CODE")

// 	err := paginatedQ.Find(&datas).Error
// 	if err != nil {
// 		return pages, &exceptions.BaseErrorResponse{
// 			StatusCode: 500,
// 			Message:    "error when scanning data",
// 			Err:        err,
// 		}
// 	}
// 	pages.Rows = datas
// 	return pages, nil
// }

// func (r *StockOpnameRepositoryImpl) GetOnGoingStockOpname(tx *gorm.DB, companyCode float64, sysNo float64) ([]transactionsparepartpayloads.GetOnGoingStockOpnameResponse, *exceptions.BaseErrorResponse) {
// 	var datas []transactionsparepartpayloads.GetOnGoingStockOpnameResponse
// 	var list []transactionsparepartpayloads.GetItemListResponse

// 	lists := tx.Model(&transactionsparepartentities.AtStockOpname1{}).
// 		Select(`atstockopname1.record_status,
// 		atstockopname1.stock_opname_sys_no,
// 		atstockopname1.stock_opname_line,
// 		atstockopname1.line_status,
// 		atstockopname1.item_code,
// 		b.item_name,
// 		atstockopname1.loc_code as location,
// 		atstockopname1.sys_qty,
// 		atstockopname1.found_qty,
// 		atstockopname1.broke_qty,
// 		atstockopname1.broken_loc_code,
// 		atstockopname1.need_adjustment,
// 		atstockopname1.remark,
// 		atstockopname1.change_no,
// 		atstockopname1.creation_user_id,
// 		atstockopname1.creation_datetime,
// 		atstockopname1.change_user_id,
// 		atstockopname1.change_datetime`).
// 		Joins("left outer join gmitem0 b on b.item_code = atstockopname1.item_code").
// 		Where("atstockopname1.stock_opname_sys_no = ?", sysNo)

// 	if lists.Error != nil {
// 		return datas, &exceptions.BaseErrorResponse{
// 			StatusCode: 500,
// 			Message:    "error when selecting data",
// 			Err:        lists.Error,
// 		}
// 	}

// 	err := lists.Scan(&list).Error
// 	if err != nil {
// 		return datas, &exceptions.BaseErrorResponse{
// 			StatusCode: 500,
// 			Message:    "error when scanning data",
// 			Err:        err,
// 		}
// 	}

// 	query := tx.Model(&transactionsparepartentities.AtStockOpname0{}).
// 		Select(`atstockopname0.COMPANY_CODE ,
// 	atstockopname0.STOCK_OPNAME_SYS_NO as stock_opname_sys_no,
// 	atstockopname0.STOCK_OPNAME_DOC_NO,
// 	atstockopname0.STOCK_OPNAME_STATUS ,
// 	B.DESCRIPTION as warehouse_group,
// 	C.WAREHOUSE_NAME as warehouse_code,
// 	atstockopname0.LOC_RANGE_FROM as from_location,
// 	atstockopname0.LOC_RANGE_TO as to_location ,
// 	atstockopname0.SHOW_DETAIL ,
// 	CASE atstockopname0.STOCK_OPNAME_STATUS
// 		WHEN ?  THEN  ?
// 		WHEN ? 	THEN  ?
// 		WHEN  ?  THEN  ?
// 		WHEN  ?  THEN  ?
// 	END AS status,
// 	atstockopname0.PIC as person_in_charge,
// 	atstockopname0.REMARK ,
// 	atstockopname0.ITEM_GROUP as item_group,
// 	atstockopname0.EXEC_DATE_FROM as stock_opname_date_from,
// 	atstockopname0.EXEC_DATE_TO as stock_opname_date_to,
// 	atstockopname0.Include_Zero_Onhand`,
// 			"01", "Draft", "02", "In Progress", "15", "Wait Approve", "20", "Approved").
// 		Joins("LEFT OUTER JOIN dbo.gmLoc0 B ON	B.WAREHOUSE_GROUP = atstockopname0.WHS_GROUP").
// 		Joins("LEFT OUTER JOIN dbo.gmLoc1 C ON	C.WAREHOUSE_CODE = atstockopname0.WHS_CODE AND C.COMPANY_CODE = atstockopname0.COMPANY_CODE").
// 		Where("atstockopname0.record_status = ? AND atstockopname0.company_code = ? and atstockopname0.stock_opname_sys_no = ?", "A", companyCode, sysNo)

// 	if query.Error != nil {
// 		return datas, &exceptions.BaseErrorResponse{
// 			StatusCode: 500,
// 			Message:    "error when selecting data",
// 			Err:        query.Error,
// 		}
// 	}

// 	err = query.Scan(&datas).Error
// 	if err != nil {
// 		return datas, &exceptions.BaseErrorResponse{
// 			StatusCode: 500,
// 			Message:    "error when scanning data",
// 			Err:        err,
// 		}
// 	}

// 	for i := range datas {
// 		datas[i].GetItemListResponse = list
// 	}

// 	fmt.Printf("Retrieved data: %+v\n", datas)
// 	return datas, nil

// }

// func (r *StockOpnameRepositoryImpl) InsertNewStockOpname(tx *gorm.DB, request transactionsparepartpayloads.InsertNewStockOpnameRequest) (bool, *exceptions.BaseErrorResponse) {

// 	datas := transactionsparepartentities.AtStockOpname0{
// 		RecordStatus:      request.Status,
// 		CompanyCode:       request.CompanyCode,
// 		StockOpnameStatus: request.StockOpnameStatus,
// 		WhsGroup:          request.WarehouseGroup,
// 		WhsCode:           request.WarehouseCode,
// 		LocRangeFrom:      request.FromLocation,
// 		LocRangeTo:        request.ToLocation,
// 		ExecDateFrom:      request.StockOpnameDateFrom,
// 		ExecDateTo:        request.StockOpnameDateTo,
// 		Pic:               request.PersonInCharge,
// 		ItemGroup:         request.ItemGroup,
// 		Remark:            request.Remark,
// 		CreationUserId:    request.UserIdCreated,
// 		CreationDatetime:  time.Now(),
// 		TotalAdjCost:      request.TotalAdjCost,
// 		ChangeUserId:      request.UserIdCreated,
// 		ChangeDatetime:    time.Now(),
// 	}

// 	fmt.Println("date", datas.CreationDatetime)

// 	query := tx.Create(&datas)
// 	if query.Error != nil {
// 		return false, &exceptions.BaseErrorResponse{
// 			StatusCode: 500,
// 			Message:    "error when inserting data",
// 			Err:        query.Error,
// 		}
// 	}

// 	return true, nil
// }

// func (r *StockOpnameRepositoryImpl) UpdateOnGoingStockOpname(tx *gorm.DB, sysNo float64, request transactionsparepartpayloads.InsertNewStockOpnameRequest) (bool, *exceptions.BaseErrorResponse) {
// 	query := tx.Model(&transactionsparepartentities.AtStockOpname0{}).
// 		Where("stock_opname_sys_no = ?", sysNo).
// 		Updates(map[string]interface{}{
// 			"record_status":       request.Status,
// 			"company_code":        request.CompanyCode,
// 			"stock_opname_doc_no": request.StockOpnameDocNo,
// 			"stock_opname_status": request.StockOpnameStatus,
// 			"whs_group":           request.WarehouseGroup,
// 			"whs_code":            request.WarehouseCode,
// 			"loc_range_from":      request.FromLocation,
// 			"loc_range_to":        request.ToLocation,
// 			"exec_date_from":      request.StockOpnameDateFrom,
// 			"exec_date_to":        request.StockOpnameDateTo,
// 			"pic":                 request.PersonInCharge,
// 			"item_group":          request.ItemGroup,
// 			"remark":              request.Remark,
// 			"total_adj_cost":      request.TotalAdjCost,
// 			"change_user_id":      request.UserIdCreated,
// 			"change_datetime":     time.Now(),
// 			"change_no":           gorm.Expr("change_no + ?", 1),
// 		})

// 	if query.Error != nil {
// 		return false, &exceptions.BaseErrorResponse{
// 			StatusCode: 500,
// 			Message:    "error when updating data",
// 			Err:        query.Error,
// 		}
// 	}

// 	return true, nil
// }

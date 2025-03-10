package transactionsparepartrepositoryimpl

import (
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartrepository "after-sales/api/repositories/transaction/sparepart"
	"after-sales/api/utils"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type StockOpnameRepositoryImpl struct {
}

func NewStockOpnameRepositoryImpl() transactionsparepartrepository.StockOpnameRepository {
	return &StockOpnameRepositoryImpl{}
}

// func (r *StockOpnameRepositoryImpl) GetAllStockOpname(tx *gorm.DB, filteredCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
// 	var datas []transactionsparepartpayloads.GetAllStockOpnameResponse

// 	query := tx.Model(&transactionsparepartentities.AtStockOpname0{}).
// 		Select("RECORD_STATUS",
// 			"COMPANY_CODE",
// 			"STOCK_OPNAME_SYS_NO as stock_opname_no",
// 			"STOCK_OPNAME_DOC_NO ",
// 			"STOCK_OPNAME_STATUS as stock_opname_status",
// 			"WHS_GROUP as warehouse_group",
// 			"WHS_CODE as warehouse_code",
// 			"LOC_RANGE_FROM",
// 			"LOC_RANGE_TO",
// 			"PROFIT_CENTER",
// 			"TRX_TYPE",
// 			"SHOW_DETAIL",
// 			"PIC",
// 			"ITEM_GROUP",
// 			"REMARK",
// 			"EXEC_DATE_FROM as stock_opname_from",
// 			"EXEC_DATE_TO as stock_opname_to",
// 			"ADJUST_DATE",
// 			"APPROVAL_STATUS",
// 			"APPROVAL_REQ_BY",
// 			"APPROVAL_REQ_DATE",
// 			"APPROVAL_BY",
// 			"APPROVAL_DATE",
// 			"CHANGE_NO",
// 			"CREATION_USER_ID",
// 			"CREATION_DATETIME",
// 			"CHANGE_USER_ID",
// 			"CHANGE_DATETIME",
// 			"Include_Zero_Onhand")

// 	if query.Error != nil {
// 		return pages, &exceptions.BaseErrorResponse{
// 			StatusCode: 500,
// 			Message:    "error when selecting data",
// 			Err:        query.Error,
// 		}
// 	}

// 	whereQ := utils.ApplyFilter(query, filteredCondition)
// 	paginatedQ := whereQ.Scopes(pagination.Paginate(&pages, whereQ))

// 	err := paginatedQ.Find(&datas).Error
// 	if err != nil {
// 		return pages, &exceptions.BaseErrorResponse{
// 			StatusCode: 500,
// 			Message:    "error when scanning data",
// 			Err:        err,
// 		}
// 	}

// 	fmt.Printf("Retrieved data: %+v\n", datas)
// 	pages.Rows = datas
// 	return pages, nil
// }

func (r *StockOpnameRepositoryImpl) GetAllStockOpname(tx *gorm.DB, filteredCondition []utils.FilterCondition, pages pagination.Pagination, companyCode float64, dateParams map[string]interface{}) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var datas []transactionsparepartpayloads.GetAllStockOpnameResponse

	query := tx.Model(&transactionsparepartentities.AtStockOpname0{}).
		Select(`atstockopname0.COMPANY_CODE ,																			
	atstockopname0.STOCK_OPNAME_SYS_NO ,																			
	atstockopname0.STOCK_OPNAME_DOC_NO as stock_opname_no,		
	atstockopname0.STOCK_OPNAME_STATUS ,																
	B.DESCRIPTION as warehouse_group,																	
	C.WAREHOUSE_NAME as warehouse_name,																			
	atstockopname0.LOC_RANGE_FROM ,																			
	atstockopname0.LOC_RANGE_TO ,
	atstockopname0.SHOW_DETAIL ,																	
	CASE atstockopname0.STOCK_OPNAME_STATUS																						
		WHEN ?  THEN  ?			
		WHEN ? 	THEN  ?	
		WHEN  ?  THEN  ?			
		WHEN  ?  THEN  ?		
	END AS status,																																						
	atstockopname0.PIC ,																			
	atstockopname0.REMARK ,																			
	atstockopname0.EXEC_DATE_FROM as stock_opname_from,																			
	atstockopname0.EXEC_DATE_TO as stock_opname_to,
	atstockopname0.Include_Zero_Onhand`,
			"01", "Draft", "02", "In Progress", "15", "Wait Approve", "20", "Approved").
		Joins("LEFT OUTER JOIN dbo.gmLoc0 B ON	B.WAREHOUSE_GROUP = atstockopname0.WHS_GROUP").
		Joins("LEFT OUTER JOIN dbo.gmLoc1 C ON	C.WAREHOUSE_CODE = atstockopname0.WHS_CODE AND C.COMPANY_CODE = atstockopname0.COMPANY_CODE").
		Where("atstockopname0.record_status = ? AND atstockopname0.company_code = ?", "A", companyCode)

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

	paginatedQ := whereQ.Scopes(pagination.Paginate(&pages, query))

	err := paginatedQ.Scan(&datas).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: 500,
			Message:    "error when scanning data",
			Err:        err,
		}
	}

	pages.Rows = datas
	fmt.Printf("Retrieved data: %+v\n", datas)
	return pages, nil

}

func (r *StockOpnameRepositoryImpl) GetLocationList(tx *gorm.DB, filteredCondition []utils.FilterCondition, pages pagination.Pagination,
	companyCode float64, warehouseGroup string, warehouseCode string) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var datas []transactionsparepartpayloads.GetAllLocationList

	query := tx.Model(&transactionsparepartentities.GmLoc2{}).
		Select(`LOCATION_CODE as location_code, LOCATION_NAME as location_name,
			CASE record_status
			WHEN  ? THEN ? 
			WHEN ? THEN ? 
			END as status`, "A", "Active", "D", "Deactive").
		Where("company_code = ? and warehouse_group = ? and warehouse_code = ?", companyCode, warehouseGroup, warehouseCode)

	if query.Error != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: 500,
			Message:    "error when selecting data",
			Err:        query.Error,
		}
	}

	whereQ := utils.ApplyFilter(query, filteredCondition)
	paginatedQ := whereQ.Scopes(pagination.Paginate(&pages, whereQ)).Order("LOCATION_CODE")

	err := paginatedQ.Find(&datas).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: 500,
			Message:    "error when scanning data",
			Err:        err,
		}
	}
	pages.Rows = datas
	return pages, nil
}

func (r *StockOpnameRepositoryImpl) GetPersonInChargeList(tx *gorm.DB, filteredCondition []utils.FilterCondition, pages pagination.Pagination, companyCode float64) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var datas []transactionsparepartpayloads.GetPersonInChargeResponse

	query := tx.Model(&transactionsparepartentities.GmEmp{}).
		Distinct(`gmemp.employee_no as employee_no, gmemp.employee_name as employee_name, b.description as position,
	CASE gmemp.record_status
	WHEN ? THEN ?
	WHEN ? THEN ?
	end as status`, "A", "Active", "D", "Deactive").
		Joins("left join gmemp1 as c on c.employee_no = gmemp.employee_no and c.record_status = ?", "gmemp").
		Joins("left join comgentable1 as b on b.table_code = ? and b.table_key0 = gmemp.job_position and b.company_code = ?", "JOBPOS", "0").
		Where("gmemp.company_code = ?", companyCode)

	if query.Error != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: 500,
			Message:    "error when selecting data",
			Err:        query.Error,
		}
	}

	whereQ := utils.ApplyFilter(query, filteredCondition)
	paginatedQ := whereQ.Scopes(pagination.Paginate(&pages, whereQ)).Order("gmemp.EMPLOYEE_NO")

	err := paginatedQ.Find(&datas).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: 500,
			Message:    "error when scanning data",
			Err:        err,
		}
	}
	pages.Rows = datas
	return pages, nil
}

func (r *StockOpnameRepositoryImpl) GetItemList(tx *gorm.DB, pages pagination.Pagination, whsCode string, itemGroup string) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var datas []transactionsparepartpayloads.GetItemListResponse

	query := tx.Model(&transactionsparepartentities.AtStockOpname1{}).
		Select(`atstockopname1.record_status,
	atstockopname1.stock_opname_sys_no,
	atstockopname1.stock_opname_line,
	atstockopname1.line_status,
	atstockopname1.item_code,
	b.item_name,
	atstockopname1.loc_code as location,
	atstockopname1.sys_qty,
	atstockopname1.found_qty,
	atstockopname1.broke_qty,
	atstockopname1.broken_loc_code,
	atstockopname1.need_adjustment,
	atstockopname1.remark,
	atstockopname1.change_no,
	atstockopname1.creation_user_id,
	atstockopname1.creation_datetime,
	atstockopname1.change_user_id,
	atstockopname1.change_datetime`).
		Joins("left outer join gmitem0 b on b.item_code = atstockopname1.item_code").
		Where("atstockopname1.whs_code = ? and b.item_group = ?", whsCode, itemGroup)

	if query.Error != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: 500,
			Message:    "error when selecting data",
			Err:        query.Error,
		}
	}

	paginatedQ := query.Scopes(pagination.Paginate(&pages, query)).Order("atstockopname1.loc_code")
	err := paginatedQ.Scan(&datas).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: 500,
			Message:    "error when scanning data",
			Err:        err,
		}
	}
	pages.Rows = datas
	return pages, nil

}

// func (r *StockOpnameRepositoryImpl) GetOnGoingStockOpname(tx *gorm.DB, sysNo string) *exceptions.BaseErrorResponse {
// 	var list []transactionsparepartpayloads.GetOnGoingStockOpnameResponse
// 	query := tx.Model(&transactionsparepartentities.AtStockOpname0{}).
// 		Select(`atstockopname0.record_status,
// 	atstockopname0.company_code,
// 	atstockopname0.stock_opname_sys_no,
// 	atstockopname0.stock_opname_doc_no,
// 	atstockopname0.stock_opname_status,
// 	atstockopname0.whs_group,
// 	atstockopname0.whs_code,
// 	atstockopname0.loc_range_from,
// 	b.location_name as loc_range_from,
// 	atstockopname0.loc_range_to,
// 	c.location_name as loc_range_to,
// 	atstockopname0.profit_center,
// 	atstockopname0.trx_type,
// 	atstockopname0.show_detail,
// 	atstockopname0.pic,
// 	atstockopname0.item_group,
// 	atstockopname0.remark,
// 	atstockopname0.exec_date_from,
// 	atstockopname0.exec_date_to,
// 	atstockopname0.adjust_date,
// 	atstockopname0.approval_status,
// 	atstockopname0.approval_req_by,
// 	atstockopname0.approval_req_date,
// 	atstockopname0.approval_by,
// 	atstockopname0.approval_date,
// 	atstockopname0.change_no,
// 	atstockopname0.creation_user_id,
// 	atstockopname0.creation_datetime,
// 	atstockopname0.change_user_id,
// 	atstockopname0.change_datetime,
// 	CASE atstockopname0.stock_opname_status
// 		WHEN ? THEN ?
// 		WHEN ? THEN ?
// 		WHEN ? THEN ?
// 		WHEN ? THEN ?
// 	END AS stock_opname_status_desc,
// 	atstockopname0.include_zero_onhand`, "01", "Draft", "02", "In Progress", "15", "Wait Approve", "20", "Approved").
// 		Joins("left join gmloc2 b on atstockopname0.company_code = b.company_code and atstockopname0.whs_code = b.warehouse_code and atstockopname0.loc_range_from = b.location_code").
// 		Joins("left join gmloc2 c on atstockopname0.company_code = c.company_code and atstockopname0.whs_code = c.warehouse_code and atstockopname0.loc_range_to = c.location_code").
// 		Where("atstockopname0.stock_opname_sys_no = ?", sysNo)

// 	if query.Error != nil {
// 		return &exceptions.BaseErrorResponse{
// 			StatusCode: 500,
// 			Message:    "error when selecting data",
// 			Err:        query.Error,
// 		}
// 	}

// 	err := query.Scan(&list).Error
// 	if err != nil {
// 		return &exceptions.BaseErrorResponse{
// 			StatusCode: 500,
// 			Message:    "error when scanning data",
// 			Err:        err,
// 		}
// 	}
// 	return nil
// }

// func (r *StockOpnameRepositoryImpl) GetListForOnGoing(tx *gorm.DB, sysNo string) *exceptions.BaseErrorResponse {
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
// 		return &exceptions.BaseErrorResponse{
// 			StatusCode: 500,
// 			Message:    "error when selecting data",
// 			Err:        lists.Error,
// 		}
// 	}
// 	// paginatedQ := lists.Scopes(pagination.Paginate(&pages, lists))
// 	err := lists.Scan(&list).Error
// 	if err != nil {
// 		return &exceptions.BaseErrorResponse{
// 			StatusCode: 500,
// 			Message:    "error when scanning data",
// 			Err:        err,
// 		}
// 	}
// 	return nil
// }

func (r *StockOpnameRepositoryImpl) GetOnGoingStockOpname(tx *gorm.DB, companyCode float64, sysNo float64) ([]transactionsparepartpayloads.GetOnGoingStockOpnameResponse, *exceptions.BaseErrorResponse) {
	var datas []transactionsparepartpayloads.GetOnGoingStockOpnameResponse
	var list []transactionsparepartpayloads.GetItemListResponse

	lists := tx.Model(&transactionsparepartentities.AtStockOpname1{}).
		Select(`atstockopname1.record_status,
		atstockopname1.stock_opname_sys_no,
		atstockopname1.stock_opname_line,
		atstockopname1.line_status,
		atstockopname1.item_code,
		b.item_name,
		atstockopname1.loc_code as location,
		atstockopname1.sys_qty,
		atstockopname1.found_qty,
		atstockopname1.broke_qty,
		atstockopname1.broken_loc_code,
		atstockopname1.need_adjustment,
		atstockopname1.remark,
		atstockopname1.change_no,
		atstockopname1.creation_user_id,
		atstockopname1.creation_datetime,
		atstockopname1.change_user_id,
		atstockopname1.change_datetime`).
		Joins("left outer join gmitem0 b on b.item_code = atstockopname1.item_code").
		Where("atstockopname1.stock_opname_sys_no = ?", sysNo)

	if lists.Error != nil {
		return datas, &exceptions.BaseErrorResponse{
			StatusCode: 500,
			Message:    "error when selecting data",
			Err:        lists.Error,
		}
	}

	err := lists.Scan(&list).Error
	if err != nil {
		return datas, &exceptions.BaseErrorResponse{
			StatusCode: 500,
			Message:    "error when scanning data",
			Err:        err,
		}
	}

	query := tx.Model(&transactionsparepartentities.AtStockOpname0{}).
		Select(`atstockopname0.COMPANY_CODE ,																			
	atstockopname0.STOCK_OPNAME_SYS_NO as stock_opname_sys_no,																			
	atstockopname0.STOCK_OPNAME_DOC_NO,		
	atstockopname0.STOCK_OPNAME_STATUS ,																
	B.DESCRIPTION as warehouse_group,																	
	C.WAREHOUSE_NAME as warehouse_code,																			
	atstockopname0.LOC_RANGE_FROM as from_location,																			
	atstockopname0.LOC_RANGE_TO as to_location ,
	atstockopname0.SHOW_DETAIL ,																	
	CASE atstockopname0.STOCK_OPNAME_STATUS																						
		WHEN ?  THEN  ?			
		WHEN ? 	THEN  ?	
		WHEN  ?  THEN  ?			
		WHEN  ?  THEN  ?		
	END AS status,																																						
	atstockopname0.PIC as person_in_charge,																			
	atstockopname0.REMARK ,		
	atstockopname0.ITEM_GROUP as item_group,																	
	atstockopname0.EXEC_DATE_FROM as stock_opname_date_from,																			
	atstockopname0.EXEC_DATE_TO as stock_opname_date_to,
	atstockopname0.Include_Zero_Onhand`,
			"01", "Draft", "02", "In Progress", "15", "Wait Approve", "20", "Approved").
		Joins("LEFT OUTER JOIN dbo.gmLoc0 B ON	B.WAREHOUSE_GROUP = atstockopname0.WHS_GROUP").
		Joins("LEFT OUTER JOIN dbo.gmLoc1 C ON	C.WAREHOUSE_CODE = atstockopname0.WHS_CODE AND C.COMPANY_CODE = atstockopname0.COMPANY_CODE").
		Where("atstockopname0.record_status = ? AND atstockopname0.company_code = ? and atstockopname0.stock_opname_sys_no = ?", "A", companyCode, sysNo)

	if query.Error != nil {
		return datas, &exceptions.BaseErrorResponse{
			StatusCode: 500,
			Message:    "error when selecting data",
			Err:        query.Error,
		}
	}

	err = query.Scan(&datas).Error
	if err != nil {
		return datas, &exceptions.BaseErrorResponse{
			StatusCode: 500,
			Message:    "error when scanning data",
			Err:        err,
		}
	}

	for i := range datas {
		datas[i].GetItemListResponse = list
	}

	fmt.Printf("Retrieved data: %+v\n", datas)
	return datas, nil

}

func (r *StockOpnameRepositoryImpl) InsertNewStockOpname(tx *gorm.DB, request transactionsparepartpayloads.InsertNewStockOpnameRequest) (bool, *exceptions.BaseErrorResponse) {

	datas := transactionsparepartentities.AtStockOpname0{
		RecordStatus:      request.Status,
		CompanyCode:       request.CompanyCode,
		StockOpnameStatus: request.StockOpnameStatus,
		WhsGroup:          request.WarehouseGroup,
		WhsCode:           request.WarehouseCode,
		LocRangeFrom:      request.FromLocation,
		LocRangeTo:        request.ToLocation,
		ExecDateFrom:      request.StockOpnameDateFrom,
		ExecDateTo:        request.StockOpnameDateTo,
		Pic:               request.PersonInCharge,
		ItemGroup:         request.ItemGroup,
		Remark:            request.Remark,
		CreationUserId:    request.UserIdCreated,
		CreationDatetime:  time.Now(),
		TotalAdjCost:      request.TotalAdjCost,
		ChangeUserId:      request.UserIdCreated,
		ChangeDatetime:    time.Now(),
	}

	fmt.Println("date", datas.CreationDatetime)

	query := tx.Create(&datas)
	if query.Error != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: 500,
			Message:    "error when inserting data",
			Err:        query.Error,
		}
	}

	return true, nil
}

func (r *StockOpnameRepositoryImpl) UpdateOnGoingStockOpname(tx *gorm.DB, sysNo float64, request transactionsparepartpayloads.InsertNewStockOpnameRequest) (bool, *exceptions.BaseErrorResponse) {
	query := tx.Model(&transactionsparepartentities.AtStockOpname0{}).
		Where("stock_opname_sys_no = ?", sysNo).
		Updates(map[string]interface{}{
			"record_status":       request.Status,
			"company_code":        request.CompanyCode,
			"stock_opname_doc_no": request.StockOpnameDocNo,
			"stock_opname_status": request.StockOpnameStatus,
			"whs_group":           request.WarehouseGroup,
			"whs_code":            request.WarehouseCode,
			"loc_range_from":      request.FromLocation,
			"loc_range_to":        request.ToLocation,
			"exec_date_from":      request.StockOpnameDateFrom,
			"exec_date_to":        request.StockOpnameDateTo,
			"pic":                 request.PersonInCharge,
			"item_group":          request.ItemGroup,
			"remark":              request.Remark,
			"total_adj_cost":      request.TotalAdjCost,
			"change_user_id":      request.UserIdCreated,
			"change_datetime":     time.Now(),
			"change_no":           gorm.Expr("change_no + ?", 1),
		})

	if query.Error != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: 500,
			Message:    "error when updating data",
			Err:        query.Error,
		}
	}

	return true, nil
}

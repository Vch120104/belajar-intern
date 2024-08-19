package transactionsparepartrepositoryimpl

import (
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartrepository "after-sales/api/repositories/transaction/sparepart"
	"after-sales/api/utils"
	"fmt"
	"gorm.io/gorm"
	"net/http"
)

type PurchaseOrderRepositoryImpl struct {
}

func NewPurchaseOrderRepositoryImpl() transactionsparepartrepository.PurchaseOrderRepository {
	return &PurchaseOrderRepositoryImpl{}
}

func (repo *PurchaseOrderRepositoryImpl) GetAllPurchaseOrder(db *gorm.DB, filter []utils.FilterCondition, page pagination.Pagination, DateParams map[string]string) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var payloadsresdb []transactionsparepartpayloads.GetAllDBResponses
	var Result []transactionsparepartpayloads.GetAllResponses
	entities := transactionsparepartentities.PurchaseOrderEntities{}
	var strfilter string
	if DateParams["PurchaseRequestDocNo"] == "" {
		strfilter = "1=1"
	} else {
		strfilter = DateParams["PurchaseRequestDocNo"]
		fmt.Println(strfilter)
	}
	JoinTable := db.Table("trx_item_purchase_order as A").
		Select("A.purchase_order_system_number,A.purchase_order_document_number,A.purchase_order_document_date,A.purchase_order_status_id,A.purchase_order_type_id,A.warehouse_id,A.supplier_id,C.purchase_request_document_number").
		Joins("LEFT JOIN trx_item_purchase_order_detail B ON A.purchase_order_system_number = B.purchase_order_system_number LEFT JOIN trx_purchase_request C ON B.purchase_request_system_number = C.purchase_request_system_number").
		Where(strfilter)
	whereQuery := utils.ApplyFilter(JoinTable, filter)
	var strDateFilter string
	if DateParams["purchase_order_date_from"] == "" {
		DateParams["purchase_order_date_from"] = "19000101"
	}
	if DateParams["purchase_order_date_to"] == "" {
		DateParams["purchase_order_date_to"] = "99991212"
	}
	strDateFilter = "purchase_order_document_date >='" + DateParams["purchase_order_date_from"] + "' AND purchase_order_document_date <= '" + DateParams["purchase_order_date_to"] + "'"

	err := whereQuery.Scopes(pagination.Paginate(&entities, &page, JoinTable)).Order("A.purchase_order_document_date desc").Where(strDateFilter).Scan(&payloadsresdb).Error
	if err != nil {
		return page, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	if len(payloadsresdb) == 0 {
		return page, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	for _, i := range payloadsresdb {
		tempres := transactionsparepartpayloads.GetAllResponses{
			PurchaseOrderSystemNumber:   i.PurchaseOrderSystemNumber,
			PurchaseOrderDocumentNumber: i.PurchaseOrderDocumentNumber,
			PurchaseOrderDocumentDate:   i.PurchaseOrderDocumentDate,
			PurchaseOrderStatus:         "",
			OrderType:                   "",
			WarehouseName:               "",
			SupplierName:                "",
			PurchaseRequestSystemNumber: "",
		}
		Result = append(Result, tempres)
	}
	page.Rows = Result
	return page, nil
}

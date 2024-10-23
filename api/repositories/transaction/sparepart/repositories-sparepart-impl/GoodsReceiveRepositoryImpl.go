package transactionsparepartrepositoryimpl

import (
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartrepository "after-sales/api/repositories/transaction/sparepart"
	"after-sales/api/utils"
	"errors"
	"gorm.io/gorm"
	"net/http"
)

type GoodsReceiveRepository struct {
}

func NewGoodsReceiveRepositoryImpl() transactionsparepartrepository.GoodsReceiveRepository {
	return &GoodsReceiveRepository{}
}

func (repository *GoodsReceiveRepository) GetAllGoodsReceive(db *gorm.DB, filter []utils.FilterCondition, paginations pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var responses []transactionsparepartpayloads.GoodsReceivesGetAllPayloads
	Entities := transactionsparepartentities.GoodsReceive{}
	JoinTable := db.Table("trx_goods_receive IG").
		Joins(`LEFT OUTER JOIN trx_goods_receive_detail IG1 ON IG.goods_receive_system_number = ig1.goods_receive_system_number`).
		Joins("LEFT OUTER JOIN mtr_item_group itemgroup ON IG.item_group_id = itemgroup.item_group_id").
		Joins(`INNER JOIN mtr_reference_type_goods_receive reftype ON reftype.reference_type_good_receive_id = ig.reference_type_good_receive_id`).
		Select(`
						ig.goods_receive_system_number,
						ig.goods_receive_document_number,
						itemgroup.item_group_name,
						ig.goods_receive_document_date,
						ig.reference_document_number,
						ig.supplier_id,
						ig.goods_receive_status_id,
						ig.journal_system_number,
						ig.supplier_delivery_order_number,
						SUM(ISNULL(ig1.quantity_goods_receive,0)) AS quantity_goods_receive,
						SUM(ISNULL(ig1.quantity_goods_receive,0) * ISNULL(ig1.item_price,0)) AS total_amount
					`).
		Group(`	ig.goods_receive_system_number,
						ig.goods_receive_document_number,
						itemgroup.item_group_name,
						ig.goods_receive_document_date,
						ig.reference_document_number,
						ig.supplier_id,
						ig.goods_receive_status_id,
						ig.journal_system_number,
						ig.supplier_delivery_order_number`)
	WhereQuery := utils.ApplyFilter(JoinTable, filter)
	//for i, res := range responses {
	//	var SupplierData generalservicepayloads.SupplierMasterCrossServicePayloads
	//	SupplierDataUrl := config.EnvConfigs.GeneralServiceUrl + "supplier/" + strconv.Itoa(res.SupplierId)
	//	if err := utils.Get(SupplierDataUrl, &SupplierData, nil); err != nil {
	//		return paginations, &exceptions.BaseErrorResponse{
	//			StatusCode: http.StatusInternalServerError,
	//			Message:    "Failed to fetch Supplier data from external service" + err.Error(),
	//			Err:        err,
	//		}
	//	}
	//responses[i].SupplierName = SupplierData.SupplierName
	//}
	err := WhereQuery.Scopes(pagination.Paginate(&Entities, &paginations, WhereQuery)).Scan(&responses).Error
	if err != nil {
		return paginations, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error On Paginate Goods Receive",
		}
	}
	paginations.Rows = responses
	return paginations, nil
}
func (repository *GoodsReceiveRepository) GetGoodsReceiveById(db *gorm.DB, GoodsReceiveId int) (transactionsparepartpayloads.GoodsReceivesGetByIdResponses, *exceptions.BaseErrorResponse) {
	var response transactionsparepartpayloads.GoodsReceivesGetByIdResponses
	err := db.Table("trx_goods_receive A").
		Joins("LEFT OUTER JOIN mtr_warehouse_master D ON D.warehouse_id = A.warehouse_id AND A.company_id = D.company_id").
		Joins("LEFT OUTER JOIN mtr_warehouse_master E ON E.warehouse_id = A.warehouse_id AND E.company_id = A.company_id").
		Select(`
	 A.goods_receive_system_number,
       A.goods_receive_status_id,
       A.goods_receive_document_number,
       A.goods_receive_document_date,
       A.item_group_id,
       a.reference_type_good_receive_id,
       a.reference_system_number,
       a.reference_document_number,
       a.affiliated_purchase_order,
       a.via_binning,
       a.back_order,
       a.set_order,
       a.brand_id,
       a.cost_center_id,
       a.profit_center_id,
       a.transaction_type_id,
       a.event_id,
       a.supplier_id,
       a.supplier_delivery_order_number,
       a.supplier_invoice_number,
       a.supplier_tax_invoice_number,
       a.warehouse_id,
       D.warehouse_code,
       D.warehouse_name,
       a.warehouse_claim_id,
       E.warehouse_code AS warehouse_claim_code,
       E.warehouse_name AS warehouse_claim_name,
       a.item_class_id
		`).
		Where("A.goods_receive_system_number = ?", GoodsReceiveId).
		Scan(&response).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err:        errors.New("header Not Found"),
			}
		}
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error On Get Goods Receive By Id",
		}
	}
	return response, nil
}

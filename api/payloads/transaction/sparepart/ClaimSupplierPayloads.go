package transactionsparepartpayloads

import "time"

type ClaimSupplierInsertPayload struct {
	CompanyId                int    `gorm:"column:company_id;null;size:30" json:"company_id"`
	ClaimDocumentNumber      string `gorm:"column:claim_document_number;null;size:25" json:"claim_document_number"`
	ClaimTypeId              int    `gorm:"column:claim_type_id;null;size:30" json:"claim_type_id"`
	TransactionTypeId        int    `gorm:"column:transaction_type_id;null;size:30" json:"transaction_type_id"`
	EventId                  int    `gorm:"column:event_id;null;size:30" json:"event_id"`
	WarehouseGroupId         int    `gorm:"column:warehouse_group_id;null;size:30" json:"warehouse_group_id"`
	WarehouseId              int    `gorm:"column:warehouse_id;null;size:30" json:"warehouse_id"`
	GoodsReceiveSystemNumber int    `gorm:"column:goods_receive_system_number;null;size:30" json:"goods_receive_system_number"`
}

type ClaimSupplierInsertDetailPayload struct {
	LocationId int `json:"location_id"`
	//ItemDiscountPercentage         float64 `gorm:"column:item_discount_percentage;null"        json:"item_discount_percentage"`
	//ItemDiscountAmount             float64 `gorm:"column:item_discount_amount;null"        json:"item_discount_amount"`
	Remark                         string `gorm:"column:remark;null;size:256"        json:"remark"`
	GoodsReceiveDetailSystemNumber int    `gorm:"column:goods_receive_detail_system_number;not null;primaryKey;size:30"        json:"goods_receive_detail_system_number"`
	GoodsReceiveSystemNumber       int    `gorm:"column:goods_receive_system_number;not null;size:30"        json:"goods_receive_system_number"`
	ClaimSystemNumber              int    `gorm:"column:claim_system_number;not null;primaryKey;size:30" json:"claim_system_number"`
}

type ClaimSupplierGetByIdResponse struct {
	CompanyId                     int       `json:"company_id"`
	CompanyCode                   string    `json:"company_code"`
	CompanyName                   string    `json:"company_name"`
	ClaimSystemNumber             int       `json:"claim_system_number"`
	ClaimStatusId                 int       `json:"claim_status_id"`
	ClaimStatusDesc               string    `json:"claim_status"`
	ClaimDocumentNumber           string    `json:"claim_document_number"`
	ClaimDate                     time.Time `json:"claim_date"`
	ClaimType                     string    `json:"claim_type"`
	ClaimTypeId                   int       `json:"claim_type_id"`
	GoodsReceiveSystemNumber      int       `json:"goods_receive_system_number"`
	GoodsReceiveDocumentNumber    string    `json:"goods_receive_document_number"`
	VehicleId                     int       `json:"vehicle_id"`
	VehicleBrand                  string    `json:"vehicle_brand"`
	CostCenterId                  int       `json:"cost_center_id"`
	ProfitCenterId                int       `json:"profit_center_id"`
	TransactionTypeId             int       `json:"transaction_type_id"`
	TransactionTypeDesc           string    `json:"transaction_type"`
	EventId                       int       `json:"event_id"`
	SupplierId                    int       `json:"supplier_id"`
	SupplierName                  string    `json:"supplier_name"`
	SupplierDoNo                  string    `json:"supplier_do_no"`
	ReferenceTypeGoodsReceiveId   int       `json:"reference_type_goods_receive_id"`
	ReferenceTypeGoodsReceiveDesc string    `json:"reference_type_goods_receive"`
	ReferenceSystemNumber         int       `json:"reference_system_number"`
	ReferenceDocumentNumber       string    `json:"reference_document_number"`
	WarehouseGroupId              int       `json:"warehouse_group_id"`
	WarehouseGroupCode            string    `json:"warehouse_group_code"`
	WarehouseId                   int       `json:"warehouse_id"`
	WarehouseCode                 string    `json:"warehouse_code"`
	WarehouseName                 string    `json:"warehouse_name"`
}

package masteritempayloads

import "time"

type BomMasterListResponseNew struct {
	BomId         int       `gorm:"column:bom_id;size:30;not null;primaryKey" json:"bom_id"`
	ItemCode      string    `json:"item_code" parent_entity:"mtr_item"`
	ItemName      string    `json:"item_name" parent_entity:"mtr_item"`
	EffectiveDate time.Time `gorm:"column:effective_date;not null;type:datetime" json:"effective_date"`
	Qty           float64   `gorm:"column:qty;size:30;not null" json:"qty"`
	UomCode       string    `json:"uom_code" parent_entity:"mtr_uom"`
	IsActive      bool      `gorm:"column:is_active;size:1;not null" json:"is_active"`
}

type BomMasterResponse struct {
	BomMasterId            int       `json:"bom_master_id"`
	IsActive               bool      `json:"is_active"`
	BomMasterCode          string    `json:"bom_master_code"`
	BomMasterSeq           int       `json:"bom_master_seq"`
	BomMasterQty           float64   `json:"bom_master_qty"`
	BomMasterEffectiveDate time.Time `json:"bom_master_effective_date"`
	BomMasterChangeNumber  int       `json:"bom_master_change_number"`
	ItemId                 int       `json:"item_id"`
	ItemCode               string    `json:"item_code"`
	ItemName               string    `json:"item_name"`
	BomMasterUom           string    `json:"bom_master_uom"`
}

type BomMasterListResponse struct {
	IsActive               bool      `json:"is_active" parent_entity:"mtr_bom"`
	BomMasterId            int       `json:"bom_master_id" parent_entity:"mtr_bom" main_table:"mtr_bom"`
	BomMasterSeq           int       `json:"bom_master_seq" parent_entity:"mtr_bom"`
	BomMasterQty           float64   `json:"bom_master_qty" parent_entity:"mtr_bom"`
	BomMasterEffectiveDate time.Time `json:"bom_master_effective_date" parent_entity:"mtr_bom"`
	BomMasterChangeNumber  int       `json:"bom_master_change_number" parent_entity:"mtr_bom"`
	ItemCode               string    `json:"item_code" parent_entity:"mtr_item"`
	ItemName               string    `json:"item_name" parent_entity:"mtr_item"`
	ItemId                 int       `json:"item_id" parent_entity:"mtr_item" references:"mtr_item" `
	UomId                  int       `json:"uom_id" parent_entity:"mtr_uom" `
	UomDescription         string    `json:"uom_description" parent_entity:"mtr_uom"`
}

type BomMasterRequest struct {
	BomMasterId            int       `json:"bom_master_id"`
	IsActive               bool      `json:"is_active"`
	BomMasterQty           float64   `json:"bom_master_qty"`
	BomMasterEffectiveDate time.Time `json:"bom_master_effective_date"`
	BomMasterChangeNumber  int       `json:"bom_master_change_number"`
	ItemId                 int       `json:"item_id"`
}

type BomDetailsResponse struct {
	Page       int                     `json:"page"`
	Limit      int                     `json:"limit"`
	TotalPages int                     `json:"total_pages"`
	TotalRows  int                     `json:"total_rows"`
	Data       []BomDetailListResponse `json:"data"`
}

type BomMasterResponseDetail struct {
	BomMasterId            int                `json:"bom_master_id"`
	IsActive               bool               `json:"is_active"`
	BomMasterQty           float64            `json:"bom_master_qty"`
	BomMasterEffectiveDate time.Time          `json:"bom_master_effective_date"`
	BomMasterChangeNumber  int                `json:"bom_master_change_number"`
	ItemId                 int                `json:"item_id"`
	ItemCode               string             `json:"item_code"`
	ItemName               string             `json:"item_name"`
	BomDetails             BomDetailsResponse `json:"bom_details"`
}

type BomItemNameResponse struct {
	ItemId   int    `json:"item_id"`
	ItemCode string `json:"item_code"`
	ItemName string `json:"item_name"`
}

type BomDetail struct {
	BomDetailId             int     `json:"bom_detail_id"`
	BomDetailSeq            int     `json:"bom_detail_seq"`
	BomDetailQty            float64 `json:"bom_detail_qty"`
	BomDetailUom            string  `json:"bom_detail_uom"`
	BomDetailRemark         string  `json:"bom_detail_remark"`
	BomDetailCostingPercent float64 `json:"bom_detail_costing_percent"`
}

type BomDetailResponse struct {
	BomDetailId             int     `json:"bom_detail_id"`
	BomDetailSeq            int     `json:"bom_detail_seq"`
	BomDetailQty            float64 `json:"bom_detail_qty"`
	BomDetailUom            string  `json:"bom_detail_uom"`
	BomDetailRemark         string  `json:"bom_detail_remark"`
	BomDetailCostingPercent float64 `json:"bom_detail_costing_percent"`
}

type BomDetailRequest struct {
	BomDetailId             int     `json:"bom_detail_id"`
	BomMasterId             int     `json:"bom_master_id"`
	BomDetailQty            float64 `json:"bom_detail_qty"`
	BomDetailRemark         string  `json:"bom_detail_remark"`
	BomDetailCostingPercent float64 `json:"bom_detail_costing_percent"`
	BomDetailTypeId         int     `json:"bom_detail_type_id"`
	BomDetailMaterialId     int     `json:"bom_detail_material_id"`
}

type BomDetailListResponse struct {
	BomMasterId             int     `json:"bom_master_id" parent_entity:"mtr_bom_detail"`
	ItemCode                string  `json:"item_code"`
	ItemName                string  `json:"item_name"`
	LineTypeName            string  `json:"line_type_name"`
	BomDetailTypeId         int     `json:"bom_detail_type_id" parent_entity:"mtr_bom_detail"`
	BomDetailId             int     `json:"bom_detail_id" parent_entity:"mtr_bom_detail"`
	BomDetailSeq            int     `json:"bom_detail_seq" parent_entity:"mtr_bom_detail"`
	BomDetailQty            float64 `json:"bom_detail_qty" parent_entity:"mtr_bom_detail"`
	BomDetailRemark         string  `json:"bom_detail_remark" parent_entity:"mtr_bom_detail"`
	BomDetailCostingPercent float64 `json:"bom_detail_costing_percent" parent_entity:"mtr_bom_detail"`
	UomDescription          string  `json:"uom_description" parent_entity:"mtr_uom"`
}

type BomItemLookup struct {
	IsActive       bool   `json:"is_active" parent_entity:"mtr_item"`
	ItemId         int    `json:"item_id" parent_entity:"mtr_item" main_table:"mtr_item"`
	ItemCode       string `json:"item_code" parent_entity:"mtr_item"`
	ItemName       string `json:"item_name" parent_entity:"mtr_item"`
	ItemTypeId     int    `json:"item_type_id" parent_entity:"mtr_item"`
	ItemGroupId    int    `json:"item_group_id" parent_entity:"mtr_item"`                                   //fk luar mtr_item_group -> item_group_name
	ItemClassId    int    `json:"item_class_id" parent_entity:"mtr_item_class" references:"mtr_item_class"` //fk dalam item_class_id -> ItemClassName
	ItemClassCode  string `json:"item_class_code" parent_entity:"mtr_item_class"`
	UomId          int    `json:"unit_of_measurement_type_id" parent_entity:"mtr_item" references:"mtr_uom"`
	UomDescription string `json:"uom_description" parent_entity:"mtr_uom"`
}

type BomUomLookup struct {
	UomId          int    `json:"uom_id" parent_entity:"mtr_uom" `
	UomDescription string `json:"uom_description" parent_entity:"mtr_uom"`
}

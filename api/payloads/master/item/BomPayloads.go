package masteritempayloads

import "time"

type BomMasterResponse struct {
	BomMasterId            int       `json:"bom_master_id"`
	IsActive               bool      `json:"is_active"`
	BomMasterCode          string    `json:"bom_master_code"`
	BomMasterSeq           int       `json:"bom_master_seq"`
	BomMasterQty           int       `json:"bom_master_qty"`
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
	BomMasterQty           int       `json:"bom_master_qty" parent_entity:"mtr_bom"`
	BomMasterEffectiveDate time.Time `json:"bom_master_effective_date" parent_entity:"mtr_bom"`
	BomMasterChangeNumber  int       `json:"bom_master_change_number" parent_entity:"mtr_bom"`
	ItemId                 int       `json:"item_id" parent_entity:"mtr_item" references:"mtr_item"`
	ItemCode               string    `json:"item_code" parent_entity:"mtr_item"`
	ItemName               string    `json:"item_name" parent_entity:"mtr_item"`
	BomMasterUom           string    `json:"bom_master_uom" parent_entity:"mtr_bom"`
}

type BomMasterRequest struct {
	BomMasterId            int       `json:"bom_master_id"`
	IsActive               bool      `json:"is_active"`
	BomMasterSeq           int       `json:"bom_master_seq"`
	BomMasterQty           int       `json:"bom_master_qty"`
	BomMasterUom           string    `json:"bom_master_uom"`
	BomMasterEffectiveDate time.Time `json:"bom_master_effective_date"`
	BomMasterChangeNumber  int       `json:"bom_master_change_number"`
	ItemId                 int       `json:"item_id"`
}

type BomItemNameResponse struct {
	ItemId   int    `json:"item_id"`
	ItemCode string `json:"item_code"`
	ItemName string `json:"item_name"`
}

type BomDetail struct {
	BomDetailId             int    `json:"bom_detail_id"`
	BomDetailSeq            int    `json:"bom_detail_seq"`
	BomDetailQty            int    `json:"bom_detail_qty"`
	BomDetailUom            string `json:"bom_detail_uom"`
	BomDetailRemark         string `json:"bom_detail_remark"`
	BomDetailCostingPercent int    `json:"bom_detail_costing_percent"`
}

type BomDetailResponse struct {
	BomDetailId             int    `json:"bom_detail_id"`
	BomDetailSeq            int    `json:"bom_detail_seq"`
	BomDetailQty            int    `json:"bom_detail_qty"`
	BomDetailUom            string `json:"bom_detail_uom"`
	BomDetailRemark         string `json:"bom_detail_remark"`
	BomDetailCostingPercent int    `json:"bom_detail_costing_percent"`
}

type BomDetailRequest struct {
	BomDetailId             int    `json:"bom_detail_id"`
	BomMasterId             int    `json:"bom_master_id"`
	BomDetailSeq            int    `json:"bom_detail_seq"`
	BomDetailQty            int    `json:"bom_detail_qty"`
	BomDetailUom            string `json:"bom_detail_uom"`
	BomDetailRemark         string `json:"bom_detail_remark"`
	BomDetailCostingPercent int    `json:"bom_detail_costing_percent"`
}

type BomDetailListResponse struct {
	BomDetailId             int    `json:"bom_detail_id" parent_entity:"mtr_bom_detail" main_table:"mtr_bom_detail"`
	BomMasterId             int    `json:"bom_master_id" parent_entity:"bom_master_id"`
	BomDetailSeq            int    `json:"bom_detail_seq" parent_entity:"mtr_bom_detail"`
	BomDetailQty            int    `json:"bom_detail_qty" parent_entity:"mtr_bom_detail"`
	BomDetailUom            string `json:"bom_detail_uom" parent_entity:"mtr_bom_detail"`
	BomDetailRemark         string `json:"bom_detail_remark" parent_entity:"mtr_bom_detail"`
	BomDetailCostingPercent int    `json:"bom_detail_costing_percent" parent_entity:"mtr_bom_detail"`
}

type BomItemLookup struct {
	IsActive      bool   `json:"is_active" parent_entity:"mtr_item"`
	ItemId        int    `json:"item_id" parent_entity:"mtr_item" main_table:"mtr_item"`
	ItemCode      string `json:"item_code" parent_entity:"mtr_item"`
	ItemName      string `json:"item_name" parent_entity:"mtr_item"`
	ItemType      string `json:"item_type" parent_entity:"mtr_item"`
	ItemGroupId   int    `json:"item_group_id" parent_entity:"mtr_item"`                                   //fk luar mtr_item_group -> item_group_name
	ItemClassId   int    `json:"item_class_id" parent_entity:"mtr_item_class" references:"mtr_item_class"` //fk dalam item_class_id -> ItemClassName
	ItemClassCode string `json:"item_class_code" parent_entity:"mtr_item_class"`
	//UomId          int    `json:"unit_of_measurement_selling_id" parent_entity:"mtr_item" references:"mtr_uom"`
	//UomDescription string `json:"uom_description" parent_entity:"mtr_uom"`
}

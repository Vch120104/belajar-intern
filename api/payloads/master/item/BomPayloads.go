package masteritempayloads

import "time"

type BomMasterResponse struct {
	BomMasterId            int       `json:"bom_master_id"`
	IsActive               bool      `json:"is_active"`
	BomMasterCode          string    `json:"bom_master_code"`
	BomMasterSeq           int       `json:"bom_master_seq"`
	BomMasterQty           int       `json:"bom_master_qty"`
	BomMasterUom           string    `json:"bom_master_uom"`
	BomMasterEffectiveDate time.Time `json:"bom_master_effective_date"`
	BomMasterChange        string    `json:"bom_master_change"`
	ItemId                 int       `json:"item_id"`
	ItemCode               string    `json:"item_code"`
	ItemName               string    `json:"item_name"`
}

type BomMasterListResponse struct {
	IsActive               bool      `json:"is_active" parent_entity:"mtr_bom"`
	BomMasterId            int       `json:"bom_master_id" parent_entity:"mtr_bom" main_table:"mtr_bom"`
	BomMasterQty           int       `json:"bom_master_qty" parent_entity:"mtr_bom"`
	BomMasterUom           string    `json:"bom_master_uom" parent_entity:"mtr_bom"`
	BomMasterEffectiveDate time.Time `json:"bom_master_effective_date" parent_entity:"mtr_bom"`
	ItemId                 int       `json:"item_id" parent_entity:"mtr_item" references:"mtr_item"`
	ItemCode               string    `json:"item_code" parent_entity:"mtr_item"`
	ItemName               string    `json:"item_name" parent_entity:"mtr_item"`
}

type BomMasterRequest struct {
	BomMasterId            int       `json:"bom_master_id"`
	ItemId                 int       `json:"item_id"`
	BomDetailId            int       `json:"bom_detail_id"`
	BomMasterQty           int       `json:"bom_master_qty"`
	BomMasterUom           string    `json:"bom_master_uom"`
	BomMasterEffectiveDate time.Time `json:"bom_master_effective_date"`
}

type BomItemNameResponse struct {
	ItemId   int    `json:"item_id"`
	ItemCode string `json:"item_code"`
	ItemName string `json:"item_name"`
}

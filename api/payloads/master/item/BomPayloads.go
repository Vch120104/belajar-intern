package masteritempayloads

import "time"

type BomListResponse struct { // View multiple bom master
	BomId         int       `json:"bom_id"`
	ItemCode      string    `json:"item_code" parent_entity:"mtr_item"`
	ItemName      string    `json:"item_name" parent_entity:"mtr_item"`
	EffectiveDate time.Time `json:"effective_date"`
	Qty           float64   `json:"qty"`
	UomCode       string    `json:"uom_code" parent_entity:"mtr_uom"`
	IsActive      bool      `json:"is_active"`
}

type BomResponse struct { // View one bom master
	BomId         int       `json:"bom_id"`
	IsActive      bool      `json:"is_active"`
	ItemId        int       `json:"item_id" parent_entity:"mtr_item"`
	Qty           float64   `json:"qty"`
	EffectiveDate time.Time `json:"effective_date"`
}

type BomDetailListResponse struct { // View multiple bom detail
	BomId             int     `json:"bom_id"`
	BomDetailId       int     `json:"bom_detail_id"`
	IsActive          bool    `json:"is_active"`
	Seq               int     `json:"seq"`
	ItemId            int     `json:"item_id"`
	ItemClassName     string  `json:"item_class_name"`
	ItemCode          string  `json:"item_code"`
	ItemName          string  `json:"item_name"`
	Qty               float64 `json:"qty"`
	UomCode           string  `json:"uom_code" parent_entity:"mtr_uom"`
	CostingPercentage float64 `json:"costing_percentage"`
	Remark            string  `json:"remark"`
}

type BomDetailResponse struct {
	IsActive          bool    `gorm:"column:is_active;size:1;not null" json:"is_active"` // Naturally, this will always be `true`
	BomDetailId       int     `gorm:"column:bom_detail_id;size:30;not null;primaryKey" json:"bom_detail_id"`
	BomId             int     `gorm:"column:bom_id;size:30;not null;index:,unique,composite:un" json:"bom_id"`
	Seq               int     `gorm:"column:seq;size:30;not null" json:"seq"`
	ItemId            int     `gorm:"column:item_id;size:30;not null;index:,unique,composite:un" json:"item_id"`
	ItemClassId       int     `gorm:"column:item_class_id;size:30;not null" json:"item_class_id"`
	Qty               float64 `gorm:"column:qty;size:30;not null" json:"qty"`
	Remark            string  `gorm:"column:remark;size:512" json:"remark"`
	CostingPercentage float64 `gorm:"column:costing_percentage;size:30;not null" json:"costing_percentage"`
}

type BomPercentageResponse struct {
	Total      float64 `json:"total"`
	IsComplete bool    `json:"is_complete"`
}

type BomMaxSeqResponse struct {
	Curr int `json:"curr"`
	Next int `json:"next"`
}

type BomDetailRequest struct {
	BomId          int     `json:"bom_id"`
	Seq            int     `json:"seq"`             // detail
	ItemId         int     `json:"item_id"`         // detail
	Qty            float64 `json:"qty"`             // detail
	Remark         string  `json:"remark"`          // detail
	CostingPercent float64 `json:"costing_percent"` // detail
	// Below are used only if BomId = 0
	BomQty           float64   `json:"bom_qty"`
	BomEffectiveDate time.Time `json:"bom_effective_date"`
	BomItemId        int       `json:"bom_item_id"`
}

type BomMasterSaveRequest struct {
	Qty float64 `json:"qty"`
}

type BomMasterNewRequest struct {
	Qty           float64   `json:"qty"`
	EffectiveDate time.Time `json:"effective_date"`
	ItemId        int       `json:"item_id"`
}

type BomDetailTemplate struct {
	ItemCode                   string    `json:"item_code"`
	EffectiveDate              time.Time `json:"effective_date"`
	Qty                        float64   `json:"qty"`
	BomDetailItemCode          string    `json:"bom_detail_item_code" parent_entity:"mtr_bom_detail"`
	BomDetailSeq               int       `json:"bom_detail_seq" parent_entity:"mtr_bom_detail"`
	BomDetailQty               float64   `json:"bom_detail_qty" parent_entity:"mtr_bom_detail"`
	BomDetailRemark            string    `json:"bom_detail_remark" parent_entity:"mtr_bom_detail"`
	BomDetailCostingPercentage float64   `json:"bom_detail_costing_percentage" parent_entity:"mtr_bom_detail"`
	Validation                 string    `json:"validation"`
}

type BomDetailUpload struct {
	BomDetails []BomDetailTemplate `json:"bom_details"`
}

type BomItemNameResponse struct {
	ItemId   int    `json:"item_id"`
	ItemCode string `json:"item_code"`
	ItemName string `json:"item_name"`
}

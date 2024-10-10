package transactionworkshopentities

import "time"

const TableNameRecall = "trx_recall"

type Recall struct {
	RecallSystemNumber   int                `gorm:"column:recall_system_number;size:30;primaryKey" json:"recall_system_number"`
	RecallDocumentNumber string             `gorm:"column:recall_document_number;" json:"recall_document_number"`
	RecallStatusId       int                `gorm:"column:recall_status_id;size:30;" json:"recall_status_id"`
	IsActive             bool               `gorm:"column:is_active;" json:"is_active"`
	RecallName           string             `gorm:"column:recall_name;" json:"recall_name"`
	RecallPeriodFrom     time.Time          `gorm:"column:recall_period_from;" json:"recall_period_from"`
	RecallPeriodTo       time.Time          `gorm:"column:recall_period_to;" json:"recall_period_to"`
	NeverExpired         bool               `gorm:"column:never_expired;" json:"never_expired"`
	BrandId              int                `gorm:"column:brand_id;" json:"brand_id"`
	RecallRemarkPopup    string             `gorm:"column:recall_remark_popup;" json:"recall_remark_popup"`
	RecallRemarkInvoice  string             `gorm:"column:recall_remark_invoice;" json:"recall_remark_invoice"`
	IsCritical           bool               `gorm:"column:is_critical;" json:"is_critical"`
	IsExtendWaranty      bool               `gorm:"column:is_extend_waranty;" json:"is_extend_waranty"`
	RecallDetail         []RecallDetail     `gorm:"foreignKey:RecallSystemNumber;references:RecallSystemNumber" json:"recall_detail"`
	RecallDetailType     []RecallDetailType `gorm:"foreignKey:RecallSystemNumber;references:RecallSystemNumber" json:"recall_detail_type"`
}

func (*Recall) TableName() string {
	return TableNameRecall
}

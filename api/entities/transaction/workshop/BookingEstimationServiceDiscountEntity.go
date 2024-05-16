package transactionworkshopentities

import "time"

const TableNameBookingEstimationServiceDiscount = "trx_booking_estimation_service_discount"

type BookingEstimationServiceDiscount struct {
	EstimationSystemNumber           int                       `gorm:"column:estimation_system_number;size:30;primaryKey" json:"estimation_system_number"`
	BatchSystemNumber                int                       `gorm:"column:batch_system_number;size:30;default:null" json:"batch_system_number"`
	DocumentStatusID                 int                       `gorm:"column:document_status_id;size:30;default:null" json:"document_status_id"`
	EstimationDiscountApprovalStatus int                       `gorm:"column:estimation_discount_approval_status;size:30;default:null" json:"estimation_discount_approval_status"`
	CompanyID                        int                       `gorm:"column:company_id;size:30;default:null" json:"company_id"`
	ApprovalRequestNumber            int                       `gorm:"column:approval_request_number;size:30;default:null" json:"approval_request_number"`
	EstimationDocumentNumber         string                    `gorm:"column:estimation_document_number;type:varchar(25);not null" json:"estimation_document_number"`
	EstimationDate                   *time.Time                `gorm:"column:estimation_date;default:null" json:"estimation_date"`
	TotalPricePackage                float32                   `gorm:"column:total_price_package;default:null" json:"total_price_package"`
	TotalPriceOperation              float32                   `gorm:"column:total_price_operation;default:null" json:"total_price_operation"`
	TotalPricePart                   float32                   `gorm:"column:total_price_part;default:null" json:"total_price_part"`
	TotalPriceOil                    float32                   `gorm:"column:total_price_oil;default:null" json:"total_price_oil"`
	TotalPriceMaterial               float32                   `gorm:"column:total_price_material;default:null" json:"total_price_material"`
	TotalPriceConsumableMaterial     float32                   `gorm:"column:total_price_consumable_material;default:null" json:"total_price_consumable_material"`
	TotalSublet                      float32                   `gorm:"column:total_sublet;default:null" json:"total_sublet"`
	TotalPriceAccessories            float32                   `gorm:"column:total_price_accessories;default:null" json:"total_price_accessories"`
	TotalDiscount                    float32                   `gorm:"column:total_discount;default:null" json:"total_discount"`
	TotalVAT                         float32                   `gorm:"column:total_vat;default:null" json:"total_vat"`
	TotalAfterVAT                    float32                   `gorm:"column:total_after_vat;default:null" json:"total_after_vat"`
	AdditionalDiscountRequestPercent float32                   `gorm:"column:additional_discount_request_percent;default:null" json:"additional_discount_request_percent"`
	AdditionalDiscountRequestAmount  float32                   `gorm:"column:additional_discount_request_amount;default:null" json:"additional_discount_request_amount"`
	VATTaxRate                       float32                   `gorm:"column:vat_tax_rate;default:null" json:"vat_tax_rate"`
	DiscountApprovalBy               string                    `gorm:"column:discount_approval_by;type:varchar(10);not null" json:"discount_approval_by"`
	DiscountApprovalDate             *time.Time                `gorm:"column:discount_approval_date;default:null" json:"discount_approval_date"`
	TotalAfterDiscount               float32                   `gorm:"column:total_after_discount;default:null" json:"total_after_discount"`
	BookingEstimationDetail          []BookingEstimationDetail `gorm:"foreignKey:EstimationSystemNumber;references:EstimationSystemNumber" json:"booking_estimation_service_discount_batch"`
}

func (*BookingEstimationServiceDiscount) TableName() string {
	return TableNameBookingEstimationServiceDiscount
}

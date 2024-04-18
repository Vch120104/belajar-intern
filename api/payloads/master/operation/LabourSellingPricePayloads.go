package masteroperationpayloads

import "time"

type LabourSellingPriceResponse struct {
	IsActive             bool      `json:"is_active"`
	LabourSellingPriceId int       `json:"labour_selling_price_id"`
	CompanyId            int       `json:"company_id"`
	BrandId              int       `json:"brand_id"`
	JobTypeId            int       `json:"JobTypeId"`
	EffectiveDate        time.Time `json:"effective_date"`
	BillableTo           string    `json:"billable_to"`
	Description          string    `json:"description"`
}

type LabourSellingPriceRequest struct {
	CompanyId     int       `json:"company_id"`
	BrandId       int       `json:"brand_id"`
	JobTypeId     int       `json:"job_type_id"`
	EffectiveDate time.Time `json:"effective_date"`
	BillableTo    string    `json:"billable_to"`
	Description   string    `json:"description"`
}

type LabourSellingPriceDetailRequest struct {
	LabourSellingPriceId int     `json:"labour_selling_price_id"`
	ModelId              int     `json:"model_id"`
	VariantId            int     `json:"variant_id"`
	SellingPrice         float64 `json:"selling_price"`
}

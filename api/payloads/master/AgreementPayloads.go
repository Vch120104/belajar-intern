package masterpayloads

import "time"

type AgreementRequest struct {
	AgreementId       int       `json:"agreement_id" parent_entity:"mtr_agreement" main_table:"mtr_agreement"`
	AgreementCode     string    `json:"agreement_code" parent_entity:"mtr_agreement"`
	IsActive          bool      `json:"is_active" parent_entity:"mtr_agreement"`
	BrandId           int       `json:"brand_id" parent_entity:"mtr_agreement"`
	CustomerId        int       `json:"customer_id" parent_entity:"mtr_agreement"`
	ProfitCenterId    int       `json:"profit_center_id"  parent_entity:"mtr_agreement"`
	AgreementDateFrom time.Time `json:"agreement_date_from" parent_entity:"mtr_agreement"`
	AgreementDateTo   time.Time `json:"agreement_date_to" parent_entity:"mtr_agreement"`
	DealerId          int       `json:"company_id" parent_entity:"mtr_agreement"`
	TopId             int       `json:"top_id" parent_entity:"mtr_agreement"`
	AgreementRemark   string    `json:"agreement_remark" parent_entity:"mtr_agreement"`
}

type AgreementResponse struct {
	AgreementId       int       `json:"agreement_id"`
	AgreementCode     string    `json:"agreement_code"`
	IsActive          bool      `json:"is_active"`
	BrandId           int       `json:"brand_id"`
	CustomerId        int       `json:"customer_id"`
	CustomerCode      string    `json:"customer_code"`
	CustomerName      string    `json:"customer_name"`
	CustomerType      string    `json:"customer_type"`
	ProfitCenterId    int       `json:"profit_center_id"`
	AgreementDateFrom time.Time `json:"agreement_date_from"`
	AgreementDateTo   time.Time `json:"agreement_date_to"`
	DealerId          int       `json:"company_id"`
	DealerName        string    `json:"company_name"`
	DealerCode        string    `json:"company_code"`
	TopId             int       `json:"top_id"`
	AgreementRemark   string    `json:"agreement_remark"`
}

type AgreementCustomerResponse struct {
	CustomerId   int    `json:"customer_id"`
	CustomerCode string `json:"customer_code"`
	CustomerName string `json:"customer_name"`
	CustomerType string `json:"customer_type"`
}

type AgreementCompanyResponse struct {
	CompanyId   int    `json:"company_id"`
	CompanyCode string `json:"company_code"`
	CompanyName string `json:"company_name"`
	CompanyType string `json:"company_type"`
}

type DiscountGroupRequest struct {
	AgreementDiscountGroupId int     `json:"agreement_discount_group_id" parent_entity:"mtr_agreement_discount_group_detail" main_table:"mtr_agreement_discount_group_detail"`
	AgreementId              int     `json:"agreement_id" parent_entity:"mtr_agreement_discount_group_detail"`
	AgreementSelection       string  `json:"agreement_selection" parent_entity:"mtr_agreement_discount_group_detail"`
	AgreementLineTypeId      string  `json:"agreement_order_type" parent_entity:"mtr_agreement_discount_group_detail"`
	AgreementDiscountMarkup  int     `json:"agreement_discount_markup_id" parent_entity:"mtr_agreement_discount_group_detail"`
	AgreementDiscount        float32 `json:"agreement_discount" parent_entity:"mtr_agreement_discount_group_detail"`
	AgreementDetailRemaks    string  `json:"agreement_detail_remarks" parent_entity:"mtr_agreement_discount_group_detail"`
}

type DiscountGroupResponse struct {
	AgreementDiscountGroupId int     `json:"agreement_discount_group_id"`
	AgreementId              int     `json:"agreement_id"`
	AgreementSelection       string  `json:"agreement_selection"`
	AgreementLineTypeId      string  `json:"agreement_order_type"`
	AgreementDiscountMarkup  int     `json:"agreement_discount_markup_id"`
	AgreementDiscount        float32 `json:"agreement_discount"`
	AgreementDetailRemaks    string  `json:"agreement_detail_remarks"`
}

type ItemDiscountRequest struct {
	ItemDiscountId     int    `json:"item_discount_id" parent_entity:"mtr_item_discount" main_table:"mtr_item_discount"`
	DiscountGroupId    int    `json:"discount_group_id" parent_entity:"mtr_item_discount"`
	ItemId             int    `json:"item_id" parent_entity:"mtr_item_discount"`
	ItemDiscountRemark string `json:"item_discount_remark" parent_entity:"mtr_item_discount"`
}

type ItemDiscountResponse struct {
	ItemDiscountId     int    `json:"item_discount_id"`
	DiscountGroupId    int    `json:"discount_group_id"`
	ItemId             int    `json:"item_id"`
	ItemCode           string `json:"item_code"`
	ItemName           string `json:"item_name"`
	ItemDiscountRemark string `json:"item_discount_remark"`
}

type DiscountValueRequest struct {
	DiscountValueId     int     `json:"discount_value_id" parent_entity:"mtr_discount_value" main_table:"mtr_discount_value"`
	DiscountGroupId     int     `json:"discount_group_id" parent_entity:"mtr_discount_value"`
	DiscountValue       float64 `json:"discount_value" parent_entity:"mtr_discount_value"`
	DiscountValueRemark string  `json:"discount_value_remark" parent_entity:"mtr_discount_value"`
}

type DiscountValueResponse struct {
	DiscountValueId     int     `json:"discount_value_id"`
	DiscountGroupId     int     `json:"discount_group_id"`
	DiscountValue       float64 `json:"discount_value"`
	DiscountValueRemark string  `json:"discount_value_remark"`
}

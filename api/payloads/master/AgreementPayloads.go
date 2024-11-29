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
	AddressStreet1    string    `json:"address_street_1"`
	AddressStreet2    string    `json:"address_street_2"`
	AddressStreet3    string    `json:"address_street_3"`
	VillageName       string    `json:"village_name"`
	VillageZipCode    string    `json:"village_zip_code"`
	DistrictName      string    `json:"district_name"`
	CityName          string    `json:"city_name"`
	ProvinceName      string    `json:"province_name"`
	CountryName       string    `json:"country_name"`
	ProfitCenterId    int       `json:"profit_center_id"`
	ProfitCenterName  string    `json:"profit_center_name"`
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
	AgreementSelection       int     `json:"agreement_selection_id" parent_entity:"mtr_agreement_discount_group_detail"`
	AgreementOrderTypeId     int     `json:"agreement_order_type_id" parent_entity:"mtr_agreement_discount_group_detail"`
	AgreementDiscountMarkup  int     `json:"agreement_discount_markup_id" parent_entity:"mtr_agreement_discount_group_detail"`
	AgreementDiscount        float32 `json:"agreement_discount" parent_entity:"mtr_agreement_discount_group_detail"`
	AgreementDetailRemaks    string  `json:"agreement_detail_remarks" parent_entity:"mtr_agreement_discount_group_detail"`
}

type DiscountGroupResponse struct {
	AgreementDiscountGroupId int     `json:"agreement_discount_group_id"`
	AgreementId              int     `json:"agreement_id"`
	AgreementSelection       int     `json:"agreement_selection_id"`
	AgreementOrderTypeId     int     `json:"agreement_order_type_id"`
	AgreementDiscountMarkup  int     `json:"agreement_discount_markup_id"`
	AgreementDiscount        float32 `json:"agreement_discount"`
	AgreementDetailRemaks    string  `json:"agreement_detail_remarks"`
}

type ItemDiscountRequest struct {
	AgreementItemId          int     `json:"agreement_item_id" parent_entity:"mtr_agreement_item_detail" main_table:"mtr_agreement_item_detail"`
	AgreementId              int     `json:"agreement_id" parent_entity:"mtr_agreement_item_detail"`
	LineTypeId               int     `json:"line_type_id" parent_entity:"mtr_agreement_item_detail"`
	AgreementItemOperationId int     `json:"agreement_item_operation_id" parent_entity:"mtr_agreement_item_detail"`
	DiscountPercent          float32 `json:"discount_percent" parent_entity:"mtr_agreement_item_detail"`
	MinValue                 int     `json:"min_value" parent_entity:"mtr_agreement_item_detail"`
	AgreementRemark          string  `json:"agreement_remark" parent_entity:"mtr_agreement_item_detail"`
}

type ItemDiscountResponse struct {
	AgreementItemId            int     `json:"agreement_item_id"`
	AgreementId                int     `json:"agreement_id"`
	LineTypeId                 int     `json:"line_type_id"`
	AgreementItemOperationId   int     `json:"agreement_item_operation_id"`
	AgreementItemOperationCode string  `json:"agreement_item_operation_code"`
	AgreementItemOperationName string  `json:"agreement_item_operation_name"`
	DiscountPercent            float32 `json:"discount_percent"`
	MinValue                   int     `json:"min_value"`
	AgreementRemark            string  `json:"agreement_remark"`
}

type DiscountValueRequest struct {
	AgreementDiscountId int     `json:"agreement_discount_id" parent_entity:"mtr_agreement_discount_detail" main_table:"mtr_agreement_discount_detail"`
	AgreementId         int     `json:"agreement_id" parent_entity:"mtr_agreement_discount_detail"`
	LineTypeId          int     `json:"line_type_id" parent_entity:"mtr_agreement_discount_detail"`
	MinValue            int     `json:"min_value" parent_entity:"mtr_agreement_discount_detail"`
	DiscountPercent     float32 `json:"discount_percent" parent_entity:"mtr_agreement_discount_detail"`
	DiscountRemarks     string  `json:"discount_remarks" parent_entity:"mtr_agreement_discount_detail"`
}

type DiscountValueResponse struct {
	AgreementDiscountId int     `json:"agreement_discount_id"`
	AgreementId         int     `json:"agreement_id"`
	LineTypeId          int     `json:"line_type_id"`
	MinValue            int     `json:"min_value"`
	DiscountPercent     float32 `json:"discount_percent"`
	DiscountRemarks     string  `json:"discount_remarks"`
}

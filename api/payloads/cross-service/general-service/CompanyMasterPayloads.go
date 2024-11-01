package generalservicepayloads

type GetCompanyByIdResponses struct {
	CompanyName            string      `json:"company_name"`
	RegionId               int         `json:"region_id"`
	TermOfPaymentId        int         `json:"term_of_payment_id"`
	TaxCompanyId           int         `json:"tax_company_id"`
	FinanceAreaId          int         `json:"finance_area_id"`
	BusinessScopeId        int         `json:"business_scope_id"`
	IsActive               bool        `json:"is_active"`
	CompanyId              int         `json:"company_id"`
	AreaId                 int         `json:"area_id"`
	BusinessCategoryId     int         `json:"business_category_id"`
	CompanyPhoneNumber     string      `json:"company_phone_number"`
	IncentiveGroupId       int         `json:"incentive_group_id"`
	CompanyNoOfStall       int         `json:"company_no_of_stall"`
	CompanyCode            string      `json:"company_code"`
	CompanyTypeId          int         `json:"company_type_id"`
	CompanyFaxNumber       string      `json:"company_fax_number"`
	AftersalesAreaId       int         `json:"aftersales_area_id"`
	CompanyDealerKiaCode   string      `json:"company_dealer_kia_code"`
	CompanyTypeSellingId   int         `json:"company_type_selling_id"`
	CompanyAbbreviation    string      `json:"company_abbreviation"`
	CompanyEmail           string      `json:"company_email"`
	VatCompanyId           int         `json:"vat_company_id"`
	CompanyOwnershipId     int         `json:"company_ownership_id"`
	HeadOfficeCompanyId    interface{} `json:"head_office_company_id"`
	CompanyOfficeAddressId int         `json:"company_office_address_id"`
}

package transactionunitentities

import "time"

type ComGenVariable struct {
	CompanyCode float64 `gorm:"column:COMPANY_CODE" json:"company_code"`
	Variable    string  `gorm:"column:VARIABLE" json:"variable"`
	Value       string  `gorm:"column:VALUE" json:"value"`
	Description string  `gorm:"column:DESCRIPTION" json:"description"`
	OrderNo     float64 `gorm:"column:ORDER_NO" json:"order_no"`
}

func (cg *ComGenVariable) TableName() string {
	return "comGenVariable"
}

type GmComp0 struct {
	RecordStatus         string    `gorm:"column:RECORD_STATUS" json:"status"`
	CompanyCode          float64   `gorm:"column:COMPANY_CODE; unique" json:"company_code"`
	BizScope             string    `gorm:"column:BIZ_SCOPE" json:"biz_scope"`
	SalesAreaCode        string    `gorm:"column:SALES_AREA_CODE" json:"sales_area_code"`
	OwnershipType        string    `gorm:"column:OWNERSHIP_TYPE" json:"ownership_type"`
	RegionalCode         string    `gorm:"column:REGIONAL_CODE" json:"regional_code"`
	TopCode              string    `gorm:"column:TOP_CODE" json:"top_code"`
	CustCode             string    `gorm:"column:CUST_CODE" json:"cust_code"`
	SupplierCode         string    `gorm:"column:SUPPLIER_CODE" json:"supplier_code"`
	CompanyHo            float64   `gorm:"column:COMPANY_HO" json:"company_ho"`
	BranchDesc           string    `gorm:"column:BRANCH_DESC" json:"branch_desc"`
	CompanyType          string    `gorm:"column:COMPANY_TYPE" json:"company_type"`
	CompanyName          string    `gorm:"column:COMPANY_NAME" json:"company_name"`
	CompAbbreviation     string    `gorm:"column:COMP_ABBREVIATION" json:"comp_abbreviation"`
	CompAddress1         string    `gorm:"column:COMP_ADDRESS_1" json:"comp_address_1"`
	CompAddress2         string    `gorm:"column:COMP_ADDRESS_2" json:"comp_address_2"`
	CompAddress3         string    `gorm:"column:COMP_ADDRESS_3" json:"comp_address_3"`
	CompVillageCode      string    `gorm:"column:COMP_VILLAGE_CODE" json:"comp_village_code"`
	CompSubdistrictCode  string    `gorm:"column:COMP_SUBDISTRICT_CODE" json:"comp_subdistrict_code"`
	CompMunicipalityCode string    `gorm:"column:COMP_MUNICIPALITY_CODE" json:"comp_municipality_code"`
	CompProvinceCode     string    `gorm:"column:COMP_PROVINCE_CODE" json:"comp_province_code"`
	CompCityCode         string    `gorm:"column:COMP_CITY_CODE" json:"comp_city_code"`
	CompZipCode          string    `gorm:"column:COMP_ZIP_CODE" json:"comp_zip_code"`
	CompPhoneNo          string    `gorm:"column:COMP_PHONE_NO" json:"comp_phone_no"`
	CompFaxNo            string    `gorm:"column:COMP_FAX_NO" json:"comp_fax_no"`
	CompEmailAddress     string    `gorm:"column:COMP_EMAIL_ADDRESS" json:"comp_email_address"`
	TaxCompany           float64   `gorm:"column:TAX_COMPANY" json:"tax_company"`
	TaxBranchCode        string    `gorm:"column:TAX_BRANCH_CODE" json:"tax_branch_code"`
	TaxRegNo             string    `gorm:"column:TAX_REG_NO" json:"tax_reg_no"`
	TaxRegDate           time.Time `gorm:"column:TAX_REG_DATE" json:"tax_reg_date"`
	TaxName              string    `gorm:"column:TAX_NAME" json:"tax_name"`
	TaxAddress1          string    `gorm:"column:TAX_ADDRESS_1" json:"tax_address_1"`
	TaxAddress2          string    `gorm:"column:TAX_ADDRESS_2" json:"tax_address_2"`
	TaxAddress3          string    `gorm:"column:TAX_ADDRESS_3" json:"tax_address_3"`
	TaxVillageCode       string    `gorm:"column:TAX_VILLAGE_CODE" json:"tax_village_code"`
	TaxSubdistrictCode   string    `gorm:"column:TAX_SUBDISTRICT_CODE" json:"tax_subdistrict_code"`
	TaxMunicipalityCode  string    `gorm:"column:TAX_MUNICIPALITY_CODE" json:"tax_municipality_code"`
	TaxProvinceCode      string    `gorm:"column:TAX_PROVINCE_CODE" json:"tax_province_code"`
	TaxCityCode          string    `gorm:"column:TAX_CITY_CODE" json:"tax_city_code"`
	TaxZipCode           string    `gorm:"column:TAX_ZIP_CODE" json:"tax_zip_code"`
	PkpType              string    `gorm:"column:PKP_TYPE" json:"pkp_type"`
	PkpDate              time.Time `gorm:"column:PKP_DATE" json:"pkp_date"`
	PkpNo                string    `gorm:"column:PKP_NO" json:"pkp_no"`
	KppCode              string    `gorm:"column:KPP_CODE" json:"kpp_code"`
	VatCompany           float64   `gorm:"column:VAT_COMPANY" json:"vat_company"`
	VatRegNo             string    `gorm:"column:VAT_REG_NO" json:"vat_reg_no"`
	VatRegDate           time.Time `gorm:"column:VAT_REG_DATE" json:"vat_reg_date"`
	VatName              string    `gorm:"column:VAT_NAME" json:"vat_name"`
	VatAddress1          string    `gorm:"column:VAT_ADDRESS_1" json:"vat_address_1"`
	VatAddress2          string    `gorm:"column:VAT_ADDRESS_2" json:"vat_address_2"`
	VatAddress3          string    `gorm:"column:VAT_ADDRESS_3" json:"vat_address_3"`
	VatVillageCode       string    `gorm:"column:VAT_VILLAGE_CODE" json:"vat_village_code"`
	VatSubdistrictCode   string    `gorm:"column:VAT_SUBDISTRICT_CODE" json:"vat_subdistrict_code"`
	VatMunicipalityCode  string    `gorm:"column:VAT_MUNICIPALITY_CODE" json:"vat_municipality_code"`
	VatProvinceCode      string    `gorm:"column:VAT_PROVINCE_CODE" json:"vat_province_code"`
	VatCityCode          string    `gorm:"column:VAT_CITY_CODE" json:"vat_city_code"`
	VatZipCode           string    `gorm:"column:VAT_ZIP_CODE" json:"vat_zip_code"`
	VatPkpType           string    `gorm:"column:VAT_PKP_TYPE" json:"vat_pkp_type"`
	VatPkpDate           time.Time `gorm:"column:VAT_PKP_DATE" json:"vat_pkp_date"`
	VatPkpNo             string    `gorm:"column:VAT_PKP_NO" json:"vat_pkp_no"`
	VatKppCode           string    `gorm:"column:VAT_KPP_CODE" json:"vat_kpp_code"`
	ChangeNo             float64   `gorm:"column:CHANGE_NO" json:"change_no"`
	ChangeUserId         string    `gorm:"column:CHANGE_USER_ID" json:"change_user_id"`
	ChangeDatetime       time.Time `gorm:"column:CHANGE_DATETIME" json:"change_datetime"`
	CreationUserId       string    `gorm:"column:CREATION_USER_ID" json:"creation_user_id"`
	CreationDatetime     time.Time `gorm:"column:CREATION_DATETIME" json:"creation_datetime"`
	VatTrxCode           string    `gorm:"column:VAT_TRX_CODE" json:"vat_trx_code"`
	FinAreaCode          string    `gorm:"column:FIN_AREA_CODE" json:"fin_area_code"`
	VatReserve           string    `gorm:"column:VAT_RESERVE" json:"vat_reserve"`
	IsDistributor        bool      `gorm:"column:IS_DISTRIBUTOR" json:"is_distributor"`
	IncentiveZone        string    `gorm:"column:INCENTIVE_ZONE" json:"incentive_zone"`
	VatDatabase          string    `gorm:"column:VAT_DATABASE" json:"vat_database"`
	BizCategory          string    `gorm:"column:BIZ_CATEGORY" json:"biz_category"`
	NoOfStall            float64   `gorm:"column:NO_OF_STALL" json:"no_of_stall"`
	CoaGroupCode         string    `gorm:"column:COA_GROUP_CODE" json:"coa_group_code"`
	BizRule              string    `gorm:"column:BIZ_RULE" json:"biz_rule"`
	PrintFormName        string    `gorm:"column:PRINT_FORM_NAME" json:"print_form_name"`
	Latitude             string    `gorm:"column:LATITUDE" json:"latitude"`
	Longitude            string    `gorm:"column:LONGITUDE" json:"longitude"`
	AfsArea              string    `gorm:"column:AFS_AREA" json:"afs_area"`
	AfsIncentiveCode     string    `gorm:"column:AFS_INCENTIVE_CODE" json:"afs_incentive_code"`
	VaMap                string    `gorm:"column:VA_MAP" json:"va_map"`
	IsAtpm               bool      `gorm:"column:IS_ATPM" json:"is_atpm"`
	DealerKiaCode        string    `gorm:"column:DEALER_KIA_CODE" json:"dealer_kia_code"`
	UseEspk              bool      `gorm:"column:USE_ESPK" json:"use_espk"`
	ID                   int       `gorm:"column:id; primary_key" json:"id"`
}

func (*GmComp0) TableName() string {
	return "gmComp0"
}

type GmEmp struct {
	RecordStatus         string    `gorm:"column: RECORD_STATUS" json:"record_status"`
	EmployeeNo           string    `gorm:"column: EMPLOYEE_NO" json:"employee_no"`
	CompanyCode          float64   `gorm:"column: COMPANY_CODE" json:"company_code"`
	EmployeeName         string    `gorm:"column: EMPLOYEE_NAME" json:"employee_name"`
	EmpNickName          string    `gorm:"column: EMP_NICK_NAME" json:"emp_nick_name"`
	IdType               string    `gorm:"column: ID_TYPE" json:"id_type"`
	IdNo                 string    `gorm:"column: ID_NO" json:"id_no"`
	JobTitle             string    `gorm:"column: JOB_TITLE" json:"job_title"`
	JobPosition          string    `gorm:"column: JOB_POSITION" json:"job_position"`
	CostCenter           string    `gorm:"column: COST_CENTER" json:"cost_center"`
	ProfitCenter         string    `gorm:"column: PROFIT_CENTER" json:"profit_center"`
	SkillLvl             string    `gorm:"column: SKILL_LVL" json:"skill_lvl"`
	AccountNo            string    `gorm:"column: ACCOUNT_NO" json:"account_no"`
	AccountName          string    `gorm:"column: ACCOUNT_NAME" json:"account_name"`
	BankCode             string    `gorm:"column: BANK_CODE" json:"bank_code"`
	BankBranchCode       string    `gorm:"column: BANK_BRANCH_CODE" json:"bank_branch_code"`
	Department           string    `gorm:"column: DEPARTMENT" json:"department"`
	Manager              string    `gorm:"column: MANAGER" json:"manager"`
	OfficePhoneNo        string    `gorm:"column: OFFICE_PHONE_NO" json:"office_phone_no"`
	EmailAddress         string    `gorm:"column: EMAIL_ADDRESS" json:"email_address"`
	HomeAddress1         string    `gorm:"column: HOME_ADDRESS_1" json:"home_address_1"`
	HomeAddress2         string    `gorm:"column: HOME_ADDRESS_2" json:"home_address_2"`
	HomeAddress3         string    `gorm:"column: HOME_ADDRESS_3" json:"home_address_3"`
	HomeVillageCode      string    `gorm:"column: HOME_VILLAGE_CODE" json:"home_village_code"`
	HomeSubdistrictCode  string    `gorm:"column: HOME_SUBDISTRICT_CODE" json:"home_subdistrict_code"`
	HomeMunicipalityCode string    `gorm:"column: HOME_MUNICIPALITY_CODE" string:"home_municipality_code"`
	HomeProvinceCode     string    `gorm:"column: HOME_PROVINCE_CODE" json:"home_province_code"`
	HomeCityCode         string    `gorm:"column: HOME_CITY_CODE" json:"home_city_code"`
	HomeZipCode          string    `gorm:"column: HOME_ZIP_CODE" json:"home_zip_code"`
	HomePhoneNo          string    `gorm:"column: HOME_PHONE_NO" json:"home_phone_no"`
	MobilePhone          string    `gorm:"column: MOBILE_PHONE" json:"mobile_phone"`
	Gender               string    `gorm:"column: GENDER" json:"gender"`
	DateOfBirth          time.Time `gorm:"column: DATE_OF_BIRTH" json:"date_of_birth"`
	CityOfBirth          string    `gorm:"column: CITY_OF_BIRTH" json:"city_of_birth"`
	MaritalStatus        string    `gorm:"column: MARITAL_STATUS" json:"marital_status"`
	NoOFChildern         float64   `gorm:"column: NO_OF_CHILDREN" json:"no_of_children"`
	Citizenship          string    `gorm:"column: CITIZENSHIP" json:"citizenship"`
	StartDate            time.Time `gorm:"column: START_DATE" json:"start_date"`
	TerminationDate      time.Time `gorm:"column: TERMINATION_DATE" json:"termination_date"`
	LastEducation        string    `gorm:"column: LAST_EDUCATION" json:"last_education"`
	LastEmployment       string    `gorm:"column: LAST_EMPLOYMENT" json:"last_employment"`
	BrGroupCode          string    `gorm:"column: BR_GROUP_CODE" json:"br_group_code"`
	KeyPass              string    `gorm:"column: KEY_PASS" json:"key_pass"`
	HassPass             string    `gorm:"column: HASS_PASS" json:"hass_pass"`
	ChangeNo             float64   `gorm:"column: CHANGE_NO" json:"change_no"`
	ChangeUserId         string    `gorm:"column: CHANGE_USER_ID" json:"change_user_id"`
	ChangeDatetime       time.Time `gorm:"column: CHANGE_DATETIME" json:"change_datetime"`
	CreationUserId       string    `gorm:"column: CREATION_USER_ID" json:"creation_user_id"`
	CreationDatetime     time.Time `gorm:"column: CREATION_DATETIME" json:"creation_datetime"`
	FactorX              float64   `gorm:"column: FACTOR_X" json:"factor_x"`
	SkillLevelCode       string    `gorm:"column: SKILL_LEVEL_CODE" json:"skill_level_code"`
	TokenNotification    string    `gorm:"column: TOKEN_NOTIFICATION" json:"token_notification"`
	TokenInstanceid      string    `gorm:"column: TOKEN_INSTANCEID" json:"token_instanceid"`
	ForceExpire          string    `gorm:"column: FORCE_EXPIRE" json:"force_expire"`
}

func (ge *GmEmp) TableName() string {
	return "gmEmp"
}

type GmEmp1 struct {
	RecordStatus     string    `gorm:"column:RECORD_STATUS" json:"record_status"`
	EmployeeNo       string    `gorm:"column:EMPLOYEE_NO; primaryKey" json:"employee_no"`
	CompanyCode      float64   `gorm:"column:COMPANY_CODE; primaryKey" json:"company_code"`
	ChangeNo         float64   `gorm:"column:CHANGE_NO" json:"change_no"`
	ChangeUserId     string    `gorm:"column:CHANGE_USER_ID" json:"change_user_id"`
	ChangeDatetime   time.Time `gorm:"column:CHANGE_DATETIME" json:"change_datetime"`
	CreationUserId   string    `gorm:"column:CREATION_USER_ID" json:"creation_user_id"`
	CreationDatetime time.Time `gorm:"column:CREATION_DATETIME" json:"creation_datetime"`
}

func (ge *GmEmp1) TableName() string {
	return "gmEmp1"
}

type GmRef struct {
	CompanyCode               float64   `gorm:"column:COMPANY_CODE;primaryKey" json:"company_code"`
	CcyCode                   string    `gorm:"column:CCY_CODE" json:"ccy_code"`
	SpMarginPercentOrderKpp   float64   `gorm:"column:SP_MARGIN_PERCENT_ORDER_KPP" json:"sp_margin_percent_order_kpp"`
	StkOpnameAdjRsn           string    `gorm:"column:STK_OPNAME_ADJ_RSN" json:"stk_opname_adj_rsn"`
	LeadTimeUnitEtd           float64   `gorm:"column:LEAD_TIME_UNIT_ETD" json:"lead_time_unit_etd"`
	BankAccReceive            string    `gorm:"column:BANK_ACC_RECEIVE" json:"bank_acc_receive"`
	VatCode                   string    `gorm:"column:VAT_CODE" json:"vat_code"`
	ItemBrokenWhs             string    `gorm:"column:ITEM_BROKEN_WHS" json:"item_broken_whs"`
	ItemBrokenLoc             string    `gorm:"column:ITEM_BROKEN_LOC" json:"item_broken_loc"`
	UnitWhsCode               string    `gorm:"column:UNIT_WHS_CODE" json:"unit_whs_code"`
	UseDms                    bool      `gorm:"column:USE_DMS" json:"use_dms"`
	TimeDiff                  float64   `gorm:"column:TIME_DIFF" json:"time_diff"`
	ChangeNo                  float64   `gorm:"column:CHANGE_NO" json:"change_no"`
	CreationUserId            string    `gorm:"column:CREATION_USER_ID" json:"creation_user_id"`
	CreationDatetime          time.Time `gorm:"column:CREATION_DATETIME" json:"creation_datetime"`
	ChangeUserId              string    `gorm:"column:CHANGE_USER_ID" json:"change_user_id"`
	ChangeDatetime            time.Time `gorm:"column:CHANGE_DATETIME" json:"change_datetime"`
	OprDiscPercentOuterKpp    float64   `gorm:"column:OPR_DISC_PERCENT_OUTER_KPP" json:"opr_disc_percent_outer_kpp"`
	CheckMonthEnd             bool      `gorm:"column:CHECK_MONTH_END" json:"check_month_end"`
	DatabaseLocation          string    `gorm:"column:DATABASE_LOCATION" json:"database_location"`
	CoaGroupCode              string    `gorm:"column:COA_GROUP_CODE" json:"coa_group_code"`
	WithVat                   bool      `gorm:"column:WITH_VAT" json:"with_vat"`
	ApprovalSpm               string    `gorm:"column:APPROVAL_SPM" json:"approval_spm"`
	IsUseTaxIndustry          bool      `gorm:"column:IS_USE_TAX_INDUSTRY" json:"is_use_tax_industry"`
	MarkupPercentage          float64   `gorm:"column:MARKUP_PERCENTAGE" json:"markup_percentage"`
	IsExternalPdi             bool      `gorm:"column:IS_EXTERNAL_PDI" json:"is_external_pdi"`
	HideCost                  bool      `gorm:"column:HIDE_COST" json:"hide_cost"`
	UsePriceCode              bool      `gorm:"column:USE_PRICE_CODE" json:"use_price_code"`
	DisableEditDraftSoinvoice bool      `gorm:"column:DISABLE_EDIT_DRAFT_SOINVOICE" json:"disable_edit_draft_soinvoice"`
	ActualClose               bool      `gorm:"column:ACTUAL_CLOSE" json:"actual_close"`
	ActualCloseDate           time.Time `gorm:"column:ACTUAL_CLOSE_DATE" json:"actual_close_date"`
}

func (gr *GmRef) TableName() string {
	return "gmRef"
}

type RtInvoice0 struct {
	RecordStatus                       string    `gorm:"column:RECORD_STATUS" json:"record_status"`
	CompanyCode                        float64   `gorm:"column:COMPANY_CODE" json:"company_code"`
	InvStatus                          string    `gorm:"column:INV_STATUS" json:"inv_status"`
	InvSysNo                           float64   `gorm:"column:INV_SYS_NO; primary_key" json:"inv_sys_no"`
	InvDocNo                           string    `gorm:"column:INV_DOC_NO" json:"inv_doc_no"`
	InvType                            string    `gorm:"column:INV_TYPE" json:"inv_type"`
	InvDate                            time.Time `gorm:"column:INV_DATE" json:"inv_date"`
	InvDueDate                         time.Time `gorm:"column:INV_DUE_DATE" json:"inv_due_date"`
	Remark                             string    `gorm:"column:REMARK" json:"remark"`
	VehicleBrand                       string    `gorm:"column:VEHICLE_BRAND" json:"vehicle_brand"`
	CpcCode                            string    `gorm:"column:CPC_CODE" json:"cpc_code"`
	TrxType                            string    `gorm:"column:TRX_TYPE" json:"trx_type"`
	EventNo                            string    `gorm:"column:EVENT_NO" json:"event_no"`
	RevisiEventNo                      string    `gorm:"column:REVISI_EVENT_NO" json:"revisi_event_no"`
	BillCode                           string    `gorm:"column:BILL_CODE" json:"bill_code"`
	CustType                           string    `gorm:"column:CUST_TYPE" json:"cust_type"`
	CustCode                           string    `gorm:"column:CUST_CODE" json:"cust_code"`
	FundType                           string    `gorm:"column:FUND_TYPE" json:"fund_type"`
	SalesRepCode                       string    `gorm:"column:SALES_REP_CODE" json:"sales_rep_code"`
	CcyCode                            string    `gorm:"column:CCY_CODE" json:"ccy_code"`
	CcyExchRateType                    string    `gorm:"column:CCY_EXCH_RATE_TYPE" json:"ccy_exch_rate_type"`
	CcyExchRateDate                    time.Time `gorm:"column:CCY_EXCH_RATE_DATE" json:"ccy_exch_rate_date"`
	CcyExchRate                        float64   `gorm:"column:CCY_EXCH_RATE" json:"ccy_exch_rate"`
	TaxExchRateType                    string    `gorm:"column:TAX_EXCH_RATE_TYPE" json:"tax_exch_rate_type"`
	TaxExchRateDate                    time.Time `gorm:"column:TAX_EXCH_RATE_DATE" json:"tax_exch_rate_date"`
	TaxExchRate                        float64   `gorm:"column:TAX_EXCH_RATE" json:"tax_exch_rate"`
	LeasingSupplierCode                string    `gorm:"column:LEASING_SUPPLIER_CODE" json:"leasing_supplier_code"`
	PoSysNo                            float64   `gorm:"column:PO_SYS_NO" json:"po_sys_no"`
	PoDocNo                            string    `gorm:"column:PO_DOC_NO" json:"po_doc_no"`
	RefType                            string    `gorm:"column:REF_TYPE" json:"ref_type"`
	RefSysNo                           float64   `gorm:"column:REF_SYS_NO" json:"ref_sys_no"`
	RefDocNo                           string    `gorm:"column:REF_DOC_NO" json:"ref_doc_no"`
	RefTnkb                            string    `gorm:"column:REF_TNKB" json:"ref_tnkb"`
	RefInsPolicyNo                     string    `gorm:"column:REF_INS_POLICY_NO" json:"ref_ins_policy_no"`
	RefServAdvisor                     string    `gorm:"column:REF_SERV_ADVISOR" json:"ref_serv_advisor"`
	RefContractServCode                string    `gorm:"column:REF_CONTRACT_SERV_CODE" json:"ref_contract_serv_code"`
	RefDocDate                         time.Time `gorm:"column:REF_DOC_DATE" json:"ref_doc_date"`
	RefInvDueDate                      time.Time `gorm:"column:REF_INV_DUE_DATE" json:"ref_inv_due_date"`
	RefTotal                           float64   `gorm:"column:REF_TOTAL" json:"ref_total"`
	RefTotalBaseAmount                 float64   `gorm:"column:REF_TOTAL_BASE_AMOUNT" json:"ref_total_base_amount"`
	Ref2Type                           string    `gorm:"column:REF2_TYPE" json:"ref2_type"`
	Ref2SysNo                          float64   `gorm:"column:REF2_SYS_NO" json:"ref2_sys_no"`
	Ref2DocNo                          string    `gorm:"column:REF2_DOC_NO" json:"ref2_doc_no"`
	BillableTo                         string    `gorm:"column:BILLABLE_TO" json:"billable_to"`
	BillToCustType                     string    `gorm:"column:BILL_TO_CUST_TYPE" json:"bill_to_cust_type"`
	BillToCustCode                     string    `gorm:"column:BILL_TO_CUST_CODE" json:"bill_to_cust_code"`
	BillToTitlePrefix                  string    `gorm:"column:BILL_TO_TITLE_PREFIX" json:"bill_to_title_prefix"`
	BillToName                         string    `gorm:"column:BILL_TO_NAME" json:"bill_to_name"`
	BillToTitleSuffix                  string    `gorm:"column:BILL_TO_TITLE_SUFFIX" json:"bill_to_title_suffix"`
	BillToIdType                       string    `gorm:"column:BILL_TO_ID_TYPE" json:"bill_to_id_type"`
	BillToIdNo                         string    `gorm:"column:BILL_TO_ID_NO" json:"bill_to_id_no"`
	BillToAddress1                     string    `gorm:"column:BILL_TO_ADDRESS_1" json:"bill_to_address1"`
	BillToAddress2                     string    `gorm:"column:BILL_TO_ADDRESS_2" json:"bill_to_address2"`
	BillToAddress3                     string    `gorm:"column:BILL_TO_ADDRESS_3" json:"bill_to_address3"`
	BillToVillageCode                  string    `gorm:"column:BILL_TO_VILLAGE_CODE" json:"bill_to_village_code"`
	BillToSubdistrictCode              string    `gorm:"column:BILL_TO_SUBDISTRICT_CODE" json:"bill_to_subdistrict_code"`
	BillToMunicipalityCode             string    `gorm:"column:BILL_TO_MUNICIPALITY_CODE" json:"bill_to_municipality_code"`
	BillToProvinceCode                 string    `gorm:"column:BILL_TO_PROVINCE_CODE" json:"bill_to_province_code"`
	BillToCityCode                     string    `gorm:"column:BILL_TO_CITY_CODE" json:"bill_to_city_code"`
	BillToZipCode                      string    `gorm:"column:BILL_TO_ZIP_CODE" json:"bill_to_zip_code"`
	BillToPhoneNo                      string    `gorm:"column:BILL_TO_PHONE_NO" json:"bill_to_phone_no"`
	BillToFax                          string    `gorm:"column:BILL_TO_FAX" json:"bill_to_fax_no"`
	BillToTaxNo                        string    `gorm:"column:BILL_TO_TAX_NO" json:"bill_to_tax_no"`
	BillToRegDate                      time.Time `gorm:"column:BILL_TO_REG_DATE" json:"bill_to_reg_date"`
	TopCode                            string    `gorm:"column:TOP_CODE" json:"top_code"`
	PayType                            string    `gorm:"column:PAY_TYPE" json:"pay_type"`
	VatCode                            string    `gorm:"column:VAT_CODE" json:"vat_code"`
	VatPercent                         float64   `gorm:"column:VAT_PERCENT" json:"vat_percent"`
	VatTaxType                         string    `gorm:"column:VAT_TAX_TYPE" json:"vat_tax_type"`
	VatTaxServCode                     string    `gorm:"column:VAT_TAX_SERV_CODE" json:"vat_tax_serv_code"`
	PkpType                            string    `gorm:"column:PKP_TYPE" json:"pkp_type"`
	PkpNo                              string    `gorm:"column:PKP_NO" json:"pkp_no"`
	PkpDate                            time.Time `gorm:"column:PKP_DATE" json:"pkp_date"`
	TaxName                            string    `gorm:"column:TAX_NAME" json:"tax_name"`
	TaxAddress1                        string    `gorm:"column:TAX_ADDRESS_1" json:"tax_address1"`
	TaxAddress2                        string    `gorm:"column:TAX_ADDRESS_2" json:"tax_address2"`
	TaxAddress3                        string    `gorm:"column:TAX_ADDRESS_3" json:"tax_address3"`
	TaxVillageCode                     string    `gorm:"column:TAX_VILLAGE_CODE" json:"tax_village_code"`
	TaxSubdistrictCode                 string    `gorm:"column:TAX_SUBDISTRICT_CODE" json:"tax_subdistrict_code"`
	TaxMunicipalityCode                string    `gorm:"column:TAX_MUNICIPALITY_CODE" json:"tax_municipality_code"`
	TaxProvinceCode                    string    `gorm:"column:TAX_PROVINCE_CODE" json:"tax_province_code"`
	TaxCityCode                        string    `gorm:"column:TAX_CITY_CODE" json:"tax_city_code"`
	TaxZipCode                         string    `gorm:"column:TAX_ZIP_CODE" json:"tax_zip_code"`
	TaxInvSysNo                        float64   `gorm:"column:TAX_INV_SYS_NO" json:"tax_inv_sys_no"`
	TaxInvDocNo                        string    `gorm:"column:TAX_INV_DOC_NO" json:"tax_inv_doc_no"`
	TaxInvDate                         time.Time `gorm:"column:TAX_INV_DATE" json:"tax_inv_date"`
	TaxInvType                         string    `gorm:"column:TAX_INV_TYPE" json:"tax_inv_type"`
	RefTaxInvSysNo                     float64   `gorm:"column:REF_TAX_INV_SYS_NO" json:"ref_tax_inv_sys_no"`
	RefTaxInvDocNo                     string    `gorm:"column:REF_TAX_INV_DOC_NO" json:"ref_tax_inv_doc_no"`
	RefTaxInvDate                      time.Time `gorm:"column:REF_TAX_INV_DATE" json:"ref_tax_inv_date"`
	IvrSysNo                           float64   `gorm:"column:IVR_SYS_NO" json:"ivr_sys_no"`
	TotalPpnbm                         float64   `gorm:"column:TOTAL_PPNBM" json:"total_ppnbm"`
	TotalMediatorFee                   float64   `gorm:"column:TOTAL_MEDIATOR_FEE" json:"total_mediator_fee"`
	TotalInsurence                     float64   `gorm:"column:TOTAL_INSURENCE" json:"total_insurence"`
	TotalServ                          float64   `gorm:"column:TOTAL_SERV" json:"total_serv"`
	TotalOpr                           float64   `gorm:"column:TOTAL_OPR" json:"total_opr"`
	TotalPart                          float64   `gorm:"column:TOTAL_PART" json:"total_part"`
	TotalOil                           float64   `gorm:"column:TOTAL_OIL" json:"total_oil"`
	TotalMaterial                      float64   `gorm:"column:TOTAL_MATERIAL" json:"total_material"`
	TotalConsumableMaterial            float64   `gorm:"column:TOTAL_CONSUMABLE_MATERIAL" json:"total_consumable_material"`
	TotalSublet                        float64   `gorm:"column:TOTAL_SUBLET" json:"total_sublet"`
	TotalAccs                          float64   `gorm:"column:TOTAL_ACCS" json:"total_accs"`
	TotalPph                           float64   `gorm:"column:TOTAL_PPH" json:"total_pph"`
	TotalPphBaseAmount                 float64   `gorm:"column:TOTAL_PPH_BASE_AMOUNT" json:"total_pph_base_amount"`
	TotalBbn                           float64   `gorm:"column:TOTAL_BBN" json:"total_bbn"`
	TotalCn                            float64   `gorm:"column:TOTAL_CN" json:"total_cn"`
	TotalCnBaseAmount                  float64   `gorm:"column:TOTAL_CN_BASE_AMOUNT" json:"total_cn_base_amount"`
	TotalCnAllocated                   float64   `gorm:"column:TOTAL_CN_ALLOCATED" json:"total_cn_allocated"`
	TotalCnAllocatedBaseAmount         float64   `gorm:"column:TOTAL_CN_ALLOCATED_BASE_AMOUNT" json:"total_cn_allocated_base_amount"`
	TotalCnDp                          float64   `gorm:"column:TOTAL_CN_DP" json:"total_cn_dp"`
	TotalCnDpBaseAmount                float64   `gorm:"column:TOTAL_CN_DP_BASE_AMOUNT" json:"total_cn_dp_base_amount"`
	TotalCnBbn                         float64   `gorm:"column:TOTAL_CN_BBN" json:"total_cn_bbn"`
	TotalCnBbnBaseAmount               float64   `gorm:"column:TOTAL_CN_BBN_BASE_AMOUNT" json:"total_cn_bbn_base_amount"`
	TotalDn                            float64   `gorm:"column:TOTAL_DN" json:"total_dn"`
	TotalDnBaseAmount                  float64   `gorm:"column:TOTAL_DN_BASE_AMOUNT" json:"total_dn_base_amount"`
	TotalMinDp                         float64   `gorm:"column:TOTAL_MIN_DP" json:"total_min_dp"`
	TotalDp                            float64   `gorm:"column:TOTAL_DP" json:"total_dp"`
	TotalDpBaseAmount                  float64   `gorm:"column:TOTAL_DP_BASE_AMOUNT" json:"total_dp_base_amount"`
	TotalDpVat                         float64   `gorm:"column:TOTAL_DP_VAT" json:"total_dp_vat"`
	TotalDpVatBaseAmount               float64   `gorm:"column:TOTAL_DP_VAT_BASE_AMOUNT" json:"total_dp_vat_base_amount"`
	TotalDpAfterVat                    float64   `gorm:"column:TOTAL_DP_AFTER_VAT" json:"total_dp_after_vat"`
	TotalDpAfterVatBaseAmount          float64   `gorm:"column:TOTAL_DP_AFTER_VAT_BASE_AMOUNT" json:"total_dp_after_vat_base_amount"`
	Total                              float64   `gorm:"column:TOTAL" json:"total"`
	TotalBaseAmount                    float64   `gorm:"column:TOTAL_BASE_AMOUNT" json:"total_base_amount"`
	TotalDiscLine                      float64   `gorm:"column:TOTAL_DISC_LINE" json:"total_disc_line"`
	AddDiscPercent                     float64   `gorm:"column:ADD_DISC_PERCENT" json:"add_disc_percent"`
	AddDiscAmount                      float64   `gorm:"column:ADD_DISC_AMOUNT" json:"add_disc_amount"`
	TotalCogs                          float64   `gorm:"column:TOTAL_COGS" json:"total_cogs"`
	TotalUnitCogs                      float64   `gorm:"column:TOTAL_UNIT_COGS" json:"total_unit_cogs"`
	TotalStdAccsCogs                   float64   `gorm:"column:TOTAL_STD_ACCS_COGS" json:"total_std_accs_cogs"`
	TotalFreeAccsCogs                  float64   `gorm:"column:TOTAL_FREE_ACCS_COGS" json:"total_free_accs_cogs"`
	TotalTransportCogs                 float64   `gorm:"column:TOTAL_TRANSPORT_COGS" json:"total_transport_cogs"`
	TotalAccrued                       float64   `gorm:"column:TOTAL_ACCRUED" json:"total_accrued"`
	TotalFreeAccsAccrued               float64   `gorm:"column:TOTAL_FREE_ACCS_ACCRUED" json:"total_free_accs_accrued"`
	TotalTransportAccrued              float64   `gorm:"column:TOTAL_TRANSPORT_ACCRUED" json:"total_transport_accrued"`
	TotalDisc                          float64   `gorm:"column:TOTAL_DISC" json:"total_disc"`
	TotalDiscBaseAmount                float64   `gorm:"column:TOTAL_DISC_BASE_AMOUNT" json:"total_disc_base_amount"`
	TotalAfterDisc                     float64   `gorm:"column:TOTAL_AFTER_DISC" json:"total_after_disc"`
	TotalAfterDiscBaseAmount           float64   `gorm:"column:TOTAL_AFTER_DISC_BASE_AMOUNT" json:"total_after_disc_base_amount"`
	TotalVat                           float64   `gorm:"column:TOTAL_VAT" json:"total_vat"`
	TotalVatBaseAmount                 float64   `gorm:"column:TOTAL_VAT_BASE_AMOUNT" json:"total_vat_base_amount"`
	TotalAfterVat                      float64   `gorm:"column:TOTAL_AFTER_VAT" json:"total_after_vat"`
	TotalAfterVatBaseAmount            float64   `gorm:"column:TOTAL_AFTER_VAT_BASE_AMOUNT" json:"total_after_vat_base_amount"`
	TotalPayment                       float64   `gorm:"column:TOTAL_PAYMENT" json:"total_payment"`
	TotalPaymentBaseAmount             float64   `gorm:"column:TOTAL_PAYMENT_BASE_AMOUNT" json:"total_payment_base_amount"`
	TotalPaymentAllocated              float64   `gorm:"column:TOTAL_PAYMENT_ALLOCATED" json:"total_payment_allocated"`
	TotalPaymentAllocatedBaseAmount    float64   `gorm:"column:TOTAL_PAYMENT_ALLOCATED_BASE_AMOUNT" json:"total_payment_allocated_base_amount"`
	CnSysNo                            float64   `gorm:"column:CN_SYS_NO" json:"cn_sys_no"`
	CnDocNo                            string    `gorm:"column:CN_DOC_NO" json:"cn_doc_no"`
	CnEventNo                          string    `gorm:"column:CN_EVENT_NO" json:"cn_event_no"`
	CnRevisiEventNo                    string    `gorm:"column:CN_REVISI_EVENT_NO" json:"cn_revisi_event_no"`
	CnBbnSyNo                          float64   `gorm:"column:CN_BBN_SY_NO" json:"cn_bbn_sy_no"`
	CnBbnDocNo                         string    `gorm:"column:CN_BBN_DOC_NO" json:"cn_bbn_doc_no"`
	CnBbnEventNo                       string    `gorm:"column:CN_BBN_EVENT_NO" json:"cn_bbn_event_no"`
	CnBbnRevisiEventNo                 string    `gorm:"column:CN_BBN_REVISI_EVENT_NO" json:"cn_bbn_revisi_event_no"`
	DnSysNo                            float64   `gorm:"column:DN_SYS_NO" json:"dn_sys_no"`
	DnDocNo                            string    `gorm:"column:DN_DOC_NO" json:"dn_doc_no"`
	DnEventNo                          string    `gorm:"column:DN_EVENT_NO" json:"dn_event_no"`
	DnRevisiEventNo                    string    `gorm:"column:DN_REVISI_EVENT_NO" json:"dn_revisi_event_no"`
	RevisiInvSysNo                     float64   `gorm:"column:REVISI_INV_SYS_NO" json:"revisi_inv_sys_no"`
	RevisiInvDocNo                     string    `gorm:"column:REVISI_INV_DOC_NO" json:"revisi_inv_doc_no"`
	KwSysNo                            float64   `gorm:"column:KW_SYS_NO" json:"kw_sys_no"`
	KwDocNo                            string    `gorm:"column:KW_DOC_NO" json:"kw_doc_no"`
	JournalSysNo                       float64   `gorm:"column:JOURNAL_SYS_NO" json:"journal_sys_no"`
	JournalRevisiSysNo                 float64   `gorm:"column:JOURNAL_REVISI_SYS_NO" json:"journal_revisi_sys_no"`
	PphBuktiPtgSysNo                   float64   `gorm:"column:PPH_BUKTI_PTG_SYS_NO" json:"pph_bukti_ptg_sys_no"`
	PphBuktiPtgDocNo                   string    `gorm:"column:PPH_BUKTI_PTG_DOC_NO" json:"pph_bukti_ptg_doc_no"`
	ApprovalReqNo                      float64   `gorm:"column:APPROVAL_REQ_NO" json:"approval_req_no"`
	ApprovalReqBy                      string    `gorm:"column:APPROVAL_REQ_BY" json:"approval_req_by"`
	ApprovalReqDate                    time.Time `gorm:"column:APPROVAL_REQ_DATE" json:"approval_req_date"`
	ApprovalLastBy                     string    `gorm:"column:APPROVAL_LAST_BY" json:"approval_last_by"`
	ApprovalLastDate                   time.Time `gorm:"column:APPROVAL_LAST_DATE" json:"approval_last_date"`
	ApprovalRemark                     string    `gorm:"column:APPROVAL_REMARK" json:"approval_remark"`
	VoidReqNo                          float64   `gorm:"column:VOID_REQ_NO" json:"void_req_no"`
	VoidReqBy                          string    `gorm:"column:VOID_REQ_BY" json:"void_req_by"`
	VoidReqDate                        time.Time `gorm:"column:VOID_REQ_DATE" json:"void_req_date"`
	VoidLastBy                         string    `gorm:"column:VOID_LAST_BY" json:"void_last_by"`
	VoidLastDate                       time.Time `gorm:"column:VOID_LAST_DATE" json:"void_last_date"`
	VoidRemark                         string    `gorm:"column:VOID_REMARK" json:"void_remark"`
	VoidPaymentEventNo                 string    `gorm:"column:VOID_PAYMENT_EVENT_NO" json:"void_payment_event_no"`
	PrintingNo                         float64   `gorm:"column:PRINTING_NO" json:"printing_no"`
	LastPrintBy                        string    `gorm:"column:LAST_PRINT_BY" json:"last_print_by"`
	AtpmReplyDate                      time.Time `gorm:"column:ATPM_REPLY_DATE" json:"atpm_reply_date"`
	AtpmReplyNotes                     string    `gorm:"column:ATPM_REPLY_NOTES" json:"atpm_reply_notes"`
	AtpmReplyStatus                    string    `gorm:"column:ATPM_REPLY_STATUS" json:"atpm_reply_status"`
	AtpmReplyUserId                    string    `gorm:"column:ATPM_REPLY_USER_ID" json:"atpm_reply_user_id"`
	ChangeNo                           float64   `gorm:"column:CHANGE_NO" json:"change_no"`
	CreationUserId                     string    `gorm:"column:CREATION_USER_ID" json:"creation_user_id"`
	CreationDateTime                   time.Time `gorm:"column:CREATION_DATE_TIME" json:"creation_date_time"`
	ChangeUserId                       string    `gorm:"column:CHANGE_USER_ID" json:"change_user_id"`
	ChangeDateTime                     time.Time `gorm:"column:CHANGE_DATE_TIME" json:"change_date_time"`
	VatTrxCode                         string    `gorm:"column:VAT_TRX_CODE" json:"vat_trx_code"`
	TotalCogsReturn                    float64   `gorm:"column:TOTAL_COGS_RETURN" json:"total_cogs_return"`
	RefTrxType                         string    `gorm:"column:REF_TRX_TYPE" json:"ref_trx_type"`
	Movement                           bool      `gorm:"column:MOVEMENT" json:"movement"`
	MovementCancel                     bool      `gorm:"column:MOVEMENT_CANCEL" json:"movement_cancel"`
	TotalReturnPayment                 float64   `gorm:"column:TOTAL_RETURN_PAYMENT" json:"total_return_payment"`
	TotalReturnPaymentBaseAmount       float64   `gorm:"column:TOTAL_RETURN_PAYMENT_BASE_AMOUNT" json:"total_return_payment_base_amount"`
	DnCorrectionSysNo                  float64   `gorm:"column:DN_CORRECTION_SYS_NO" json:"dn_correction_sys_no"`
	TotalDnCorrection                  float64   `gorm:"column:TOTAL_DN_CORRECTION" json:"total_dn_correction"`
	TotalDnCorrectionBaseAmount        float64   `gorm:"column:TOTAL_DN_CORRECTION_BASE_AMOUNT" json:"total_dn_correction_base_amount"`
	WithVat                            bool      `gorm:"column:WITH_VAT" json:"with_vat"`
	DealerRepCode                      string    `gorm:"column:DEALER_REP_CODE" json:"dealer_rep_code"`
	TotalRobbingCogs                   float64   `gorm:"column:TOTAL_ROBBING_COGS" json:"total_robbing_cogs"`
	TotalRobbingAccrued                float64   `gorm:"column:TOTAL_ROBBING_ACCRUED" json:"total_robbing_accrued"`
	TaxIndDocNo                        string    `gorm:"column:TAX_IND_DOC_NO" json:"tax_ind_doc_no"`
	TaxIndDocDate                      time.Time `gorm:"column:TAX_IND_DOC_DATE" json:"tax_ind_doc_date"`
	TaxIndAmount                       float64   `gorm:"column:TAX_IND_AMOUNT" json:"tax_ind_amount"`
	CutOffDate                         time.Time `gorm:"column:CUT_OFF_DATE" json:"cut_off_date"`
	SelisihOfftrVsAfterDisc            float64   `gorm:"column:SELISIH_OFFTR_VS_AFTER_DISC" json:"selisih_offtr_vs_after_disc"`
	Rounding                           float64   `gorm:"column:ROUNDING" json:"rounding"`
	SelisihDppDpVsInv                  float64   `gorm:"column:SELISIH_DPP_DP_VS_INV" json:"selisih_dpp_dp_vs_inv"`
	SelisihVatDpVsInv                  float64   `gorm:"column:SELISIH_VAT_DP_VS_INV" json:"selisih_vat_dp_vs_inv"`
	FreightAmount                      float64   `gorm:"column:FREIGHT_AMOUNT" json:"freight_amount"`
	FreightBaseAmount                  float64   `gorm:"column:FREIGHT_BASE_AMOUNT" json:"freight_base_amount"`
	InvExportSysNo                     float64   `gorm:"column:INV_EXPORT_SYS_NO" json:"inv_export_sys_no"`
	TotalMaterai                       float64   `gorm:"column:TOTAL_MATERAI" json:"total_materai"`
	TotalMateraiBaseAmount             float64   `gorm:"column:TOTAL_MATERAI_BASE_AMOUNT" json:"total_materai_base_amount"`
	InsuranceAmount                    float64   `gorm:"column:INSURANCE_AMOUNT" json:"insurance_amount"`
	InsuranceBaseAmount                float64   `gorm:"column:INSURANCE_BASE_AMOUNT" json:"insurance_base_amount"`
	TotalPaymentDpp                    float64   `gorm:"column:TOTAL_PAYMENT_DPP" json:"total_payment_dpp"`
	TotalPaymentDppBaseAmount          float64   `gorm:"column:TOTAL_PAYMENT_DPP_BASE_AMOUNT" json:"total_payment_dpp_base_amount"`
	TotalPaymentDppAllocated           float64   `gorm:"column:TOTAL_PAYMENT_DPP_ALLOCATED" json:"total_payment_dpp_allocated"`
	TotalPaymentDppAllocatedBaseAmount float64   `gorm:"column:TOTAL_PAYMENT_DPP_ALLOCATED_BASE_AMOUNT" json:"total_payment_dpp_allocated_base_amount"`
	TotalPaymentVat                    float64   `gorm:"column:TOTAL_PAYMENT_VAT" json:"total_payment_vat"`
	TotalPaymentVatBaseAmount          float64   `gorm:"column:TOTAL_PAYMENT_VAT_BASE_AMOUNT" json:"total_payment_vat_base_amount"`
	TotalPaymentVatAllocated           float64   `gorm:"column:TOTAL_PAYMENT_VAT_ALLOCATED" json:"total_payment_vat_allocated"`
	TotalPaymentVatAllocatedBaseAmount float64   `gorm:"column:TOTAL_PAYMENT_VAT_ALLOCATED_BASE_AMOUNT" json:"total_payment_vat_allocated_base_amount"`
	TaxSupportingDoc                   string    `gorm:"column:TAX_SUPPORTING_DOC" json:"tax_supporting_doc"`
	GovernmentSubsidy                  float64   `gorm:"column:GOVERNMENT_SUBSIDY" json:"government_subsidy"`
	Nitku                              string    `gorm:"column:NITKU" json:"nitku"`
}

func (r0 *RtInvoice0) TableName() string {
	return "rtInvoice0"
}

type RtInvoice1 struct {
	RecordStatus             string    `gorm:"column:RECORD_STATUS" json:"record_status"`
	InvSysNo                 float64   `gorm:"column:INV_SYS_NO" json:"inv_sys_no"`
	InvLineNo                float64   `gorm:"column:INV_LINE_NO" json:"inv_line_no"`
	InvLineStatus            string    `gorm:"column:INV_LINE_STATUS" json:"inv_line_status"`
	RefType                  string    `gorm:"column:REF_TYPE" json:"ref_type"`
	RefSysNo                 float64   `gorm:"column:REF_SYS_NO" json:"ref_sys_no"`
	RefDocNo                 string    `gorm:"column:REF_DOC_NO" json:"ref_doc_no"`
	RefLineNo                float64   `gorm:"column:REF_LINE_NO" json:"ref_line_no"`
	Ref2Type                 string    `gorm:"column:REF2_TYPE" json:"ref2_type"`
	Ref2SysNo                float64   `gorm:"column:REF2_SYS_NO" json:"ref2_sys_no"`
	Ref2DocNo                string    `gorm:"column:REF2_DOC_NO" json:"ref2_doc_no"`
	Ref2LineNo               float64   `gorm:"column:REF2_LINE_NO" json:"ref2_line_no"`
	VehicleChassisNo         string    `gorm:"column:VEHICLE_NO" json:"vehicle_chassis_no"`
	VehicleEngineNo          string    `gorm:"column:VEHICLE_ENGINE_NO" json:"vehicle_engine_no"`
	VehicleBrand             string    `gorm:"column:VEHICLE_BRAND" json:"vehicle_brand"`
	ModelCode                string    `gorm:"column:MODEL_CODE" json:"model_code"`
	VariantCode              string    `gorm:"column:VARIANT_CODE" json:"variant_code"`
	ColourCode               string    `gorm:"column:COLOUR_CODE" json:"colour_code"`
	VehicleShortDesc         string    `gorm:"column:VEHICLE_SHORT_DESC" json:"vehicle_short_desc"`
	VehicleLongDesc          string    `gorm:"column:VEHICLE_LONG_DESC" json:"vehicle_long_desc"`
	SalesRepCode             string    `gorm:"column:SALES_REP_CODE" json:"sales_rep_code"`
	MediatorCode             string    `gorm:"column:MEDIATOR_CODE" json:"mediator_code"`
	MediatorFeeAmount        float64   `gorm:"column:MEDIATOR_FEE_AMOUNT" json:"mediator_fee_amount"`
	AccYearRemark            string    `gorm:"column:ACC_YEAR_REMARK" json:"acc_year_remark"`
	InsuranceSupplierCode    string    `gorm:"column:INSURANCE_SUPPLIER_CODE" json:"insurance_supplier_code"`
	InsuranceAmount          float64   `gorm:"column:INSURANCE_AMOUNT" json:"insurance_amount"`
	BbnAmount                float64   `gorm:"column:BBN_AMOUNT" json:"bbn_amount"`
	OfftrAmount              float64   `gorm:"column:OFFTR_AMOUNT" json:"offtr_amount"`
	OfftrNetAmount           float64   `gorm:"column:OFFTR_NET_AMOUNT" json:"offtr_net_amount"`
	DiscAmount               float64   `gorm:"column:DISC_AMOUNT" json:"disc_amount"`
	OntrAmount               float64   `gorm:"column:ONTR_AMOUNT" json:"ontr_amount"`
	VatAmount                float64   `gorm:"column:VAT_AMOUNT" json:"vat_amount"`
	EventNo                  string    `gorm:"column:EVENT_NO" json:"event_no"`
	CostGroupCode            string    `gorm:"column:COST_GROUP_CODE" json:"cost_group_code"`
	MinDp                    float64   `gorm:"column:MIN_DP" json:"min_dp"`
	JobType                  string    `gorm:"column:JOB_TYPE" json:"job_type"`
	ItemGroup                string    `gorm:"column:ITEM_GROUP" json:"item_group"`
	ItemLineType             string    `gorm:"column:ITEM_LINE_TYPE" json:"item_line_type"`
	ItemCode                 string    `gorm:"column:ITEM_CODE" json:"item_code"`
	ItemDesc                 string    `gorm:"column:ITEM_DESC" json:"item_desc"`
	ItemQty                  float64   `gorm:"column:ITEM_QTY" json:"item_qty"`
	ItemQtyReturn            float64   `gorm:"column:ITEM_QTY_RETURN" json:"item_qty_return"`
	ItemUom                  string    `gorm:"column:ITEM_UOM" json:"item_uom"`
	ItemPrice                float64   `gorm:"column:ITEM_PRICE" json:"item_price"`
	ItemDiscPercent          float64   `gorm:"column:ITEM_DISC_PERCENT" json:"item_disc_percent"`
	ItemDiscAmount           float64   `gorm:"column:ITEM_DISC_AMOUNT" json:"item_disc_amount"`
	ItemCogs                 float64   `gorm:"column:ITEM_COGS" json:"item_cogs"`
	ItemCogsReturn           float64   `gorm:"column:ITEM_COGS_RETURN" json:"item_cogs_return"`
	TotalItemCogs            float64   `gorm:"column:TOTAL_ITEM_COGS" json:"total_item_cogs"`
	TotalItemCogsReturn      float64   `gorm:"column:TOTAL_ITEM_COGS_RETURN" json:"total_item_cogs_return"`
	UnitCogs                 float64   `gorm:"column:UNIT_COGS" json:"unit_cogs"`
	StdAccsCogs              float64   `gorm:"column:STD_ACCS_COGS" json:"std_accs_cogs"`
	FreeAccsCogs             float64   `gorm:"column:FREE_ACCS_COGS" json:"free_accs_cogs"`
	TransportCogs            float64   `gorm:"column:TRANSPORT_COGS" json:"transport_cogs"`
	FreeAccsAccrued          float64   `gorm:"column:FREE_ACCS_ACCRUED" json:"free_accs_accrued"`
	TransportAccrued         float64   `gorm:"column:TRANSPORT_ACCRUED" json:"transport_accrued"`
	PphTaxCode               string    `gorm:"column:PPH_TAX_CODE" json:"pph_tax_code"`
	PphTaxType               string    `gorm:"column:PPH_TAX_TYPE" json:"pph_tax_type"`
	PphTaxServCode           string    `gorm:"column:PPH_TAX_SERV_CODE" json:"pph_tax_serv_code"`
	PphTaxPercent            float64   `gorm:"column:PPH_TAX_PERCENT" json:"pph_tax_percent"`
	PphAmount                float64   `gorm:"column:PPH_AMOUNT" json:"pph_amount"`
	PoSysNo                  float64   `gorm:"column:PO_SYS_NO" json:"po_sys_no"`
	PoLine                   float64   `gorm:"column:PO_LINE" json:"po_line"`
	IvrSysNo                 float64   `gorm:"column:IVR_SYS_NO" json:"ivr_sys_no"`
	IvrLineNo                float64   `gorm:"column:IVR_LINE_NO" json:"ivr_line_no"`
	Ppnbm                    float64   `gorm:"column:PPNBM" json:"ppnbm"`
	TotalCn                  float64   `gorm:"column:TOTAL_CN" json:"total_cn"`
	TotalCnBaseAmount        float64   `gorm:"column:TOTAL_CN_BASE_AMOUNT" json:"total_cn_base_amount"`
	TotalCnDp                float64   `gorm:"column:TOTAL_CN_DP" json:"total_cn_dp"`
	TotalCnDpBaseAmount      float64   `gorm:"column:TOTAL_CN_DP_BASE_AMOUNT" json:"total_cn_dp_base_amount"`
	TOtalCnBbn               float64   `gorm:"column:TOTAL_CN_BBN" json:"total_cn_bbn"`
	TotalCnBbnBaseAmount     float64   `gorm:"column:TOTAL_CN_BBN_BASE_AMOUNT" json:"total_cn_bbn_base_amount"`
	TotalDn                  float64   `gorm:"column:TOTAL_DN" json:"total_dn"`
	TotalDnBaseAmount        float64   `gorm:"column:TOTAL_DN_BASE_AMOUNT" json:"total_dn_base_amount"`
	TotalDp                  float64   `gorm:"column:TOTAL_DP" json:"total_dp"`
	TotalDpBaseAmount        float64   `gorm:"column:TOTAL_DP_BASE_AMOUNT" json:"total_dp_base_amount"`
	TotalDpVat               float64   `gorm:"column:TOTAL_DP_VAT" json:"total_dp_vat"`
	TotalDpVatBaseAmount     float64   `gorm:"column:TOTAL_DP_VAT_BASE_AMOUNT" json:"total_dp_vat_base_amount"`
	Total                    float64   `gorm:"column:TOTAL" json:"total"`
	TotalBaseAmount          float64   `gorm:"column:TOTAL_BASE_AMOUNT" json:"total_base_amount"`
	TotalDisc                float64   `gorm:"column:TOTAL_DISC" json:"total_disc"`
	TotalDiscBaseAmount      float64   `gorm:"column:TOTAL_DISC_BASE_AMOUNT" json:"total_disc_base_amount"`
	TotalAfterDisc           float64   `gorm:"column:TOTAL_AFTER_DISC" json:"total_after_disc"`
	TotalAfterDiscBaseAmount float64   `gorm:"column:TOTAL_AFTER_DISC_BASE_AMOUNT" json:"total_after_disc_base_amount"`
	TotalVat                 float64   `gorm:"column:TOTAL_VAT" json:"total_vat"`
	TotalVatBaseAmount       float64   `gorm:"column:TOTAL_VAT_BASE_AMOUNT" json:"total_vat_base_amount"`
	TotalAfterVat            float64   `gorm:"column:TOTAL_AFTER_VAT" json:"total_after_vat"`
	TotalAfterVatBaseAmount  float64   `gorm:"column:TOTAL_AFTER_VAT_BASE_AMOUNT" json:"total_after_vat_base_amount"`
	TotalPayment             float64   `gorm:"column:TOTAL_PAYMENT" json:"total_payment"`
	TotalPaymentAllocated    float64   `gorm:"column:TOTAL_PAYMENT_ALLOCATED" json:"total_payment_allocated"`
	ChangeNo                 float64   `gorm:"column:CHANGE_NO" json:"change_no"`
	CreationUserId           string    `gorm:"column:CREATION_USER_ID" json:"creation_user_id"`
	CreationDatetime         time.Time `gorm:"column:CREATION_DATETIME" json:"creation_datetime"`
	ChangeUserId             string    `gorm:"column:CHANGE_USER_ID" json:"change_user_id"`
	ChangeDatetime           time.Time `gorm:"column:CHANGE_DATETIME" json:"change_datetime"`
	WarehouseBrand           string    `gorm:"column:WAREHOUSE_BRAND" json:"warehouse_brand"`
	IsReturn                 bool      `gorm:"column:IS_RETURN" json:"is_return"`
	RobbingCogs              float64   `gorm:"column:ROBBING_COGS" json:"robbing_cogs"`
	RobbingAccrued           float64   `gorm:"column:ROBBING_ACCRUED" json:"robbing_accrued"`
	VarianceDpBbn            float64   `gorm:"column:VARIANCE_DP_BBN" json:"variance_dp_bbn"`
	OptionCode               string    `gorm:"column:OPTION_CODE" json:"option_code"`
	PaymentRetur             float64   `gorm:"column:PAYMENT_RETUR" json:"payment_retur"`
	PaymentReturBbn          float64   `gorm:"column:PAYMENT_RETUR_BBN" json:"payment_retur_bbn"`
	TotalDiscRetur           float64   `gorm:"column:TOTAL_DISC_RETUR" json:"total_disc_retur"`
	TotalVatRetur            float64   `gorm:"column:TOTAL_VAT_RETUR" json:"total_vat_retur"`
	TotalBebanPajak          float64   `gorm:"column:TOTAL_BEBAN_PAJAK" json:"total_beban_pajak"`
	FlagRetur                bool      `gorm:"column:FLAG_RETUR" json:"flag_retur"`
	TotalAfterDiscInvDp      float64   `gorm:"column:TOTAL_AFTER_DISC_INV_DP" json:"total_after_disc_inv_dp"`
	TotalVatInvDp            float64   `gorm:"column:TOTAL_VAT_INV_DP" json:"total_vat_inv_dp"`
	TotalAfterVatInvDp       float64   `gorm:"column:TOTAL_AFTER_VAT_INV_DP" json:"total_after_vat_inv_dp"`
	RemarkDetail             string    `gorm:"column:REMARK_DETAIL" json:"remark_detail"`
}

func (r1 *RtInvoice1) TableName() string {
	return "rtInvoice1"
}

type UtPointProspecting struct {
	CompanyCode                float64   `gorm:"column: COMPANY_CODE" json:"company_code"`
	PeriodYear                 string    `gorm:"column: PERIOD_YEAR" json:"period_year"`
	PeriodMonth                string    `gorm:"column: PERIOD_MONTH" json:"period_month"`
	SalesRepCode               string    `gorm:"column: SALES_REP_CODE" json:"sales_rep_code"`
	ProspectSystemNo           float64   `gorm:"column: PROSPECT_SYSTEM_NO" json:"prospect_system_no"`
	SpmSystemNo                float64   `gorm:"column: SPMSYSTEM_NO" json:"spm_system_no"`
	VehicleBrand               string    `gorm:"column: VEHICLE_BRAND" json:"vehicle_brand"`
	ProspectDate               int       `gorm:"column: PROSPECT_DATE" json:"prospect_date"`
	ProspectStage              int       `gorm:"column: PROSPECT_STAGE" json:"prospect_stage"`
	ProspectTitlePRefix        int       `gorm:"column: PROSPECT_TITLE_PREFIX" json:"prospect_title_prefix"`
	ProspectName               int       `gorm:"column: PROSPECT_NAME" json:"prospect_name"`
	ProspectTitleSuffix        int       `gorm:"column: PROSPECT_TITLE_SUFFIX" json:"prospect_title_suffix"`
	ProspectSrcCode            int       `gorm:"column: PROSPECT_SRC_CODE" json:"prospect_src_code"`
	ProspectType               int       `gorm:"column: PROSPECT_TYPE" json:"prospect_type"`
	ProspectRef                int       `gorm:"column: PROSPECT_REF" json:"prospect_ref"`
	ProspectNotes              int       `gorm:"column: PROSPECT_NOTES" json:"prospect_notes"`
	BuyingBudget               int       `gorm:"column: BUYING_BUDGET" json:"buying_budget"`
	BuyingPlan                 int       `gorm:"column: BUYING_PLAN" json:"buying_plan"`
	FundType                   int       `gorm:"column: FUND_TYPE" json:"fund_type"`
	ModelCode                  int       `gorm:"column: MODEL_CODE" json:"model_code"`
	VariantCode                int       `gorm:"column: VARIANT_CODE" json:"variant_code"`
	ProspectAddress            int       `gorm:"column: PROSPECT_ADDRESS" json:"prospect_address"`
	ProspectArea               int       `gorm:"column: PROSPECT_AREA" json:"prospect_area"`
	ProspectMobilePhone        int       `gorm:"column: PROSPECT_MOBILE_PHONE" json:"prospect_mobile_phone"`
	ProspectEmailAddress       int       `gorm:"column: PROSPECT_EMAIL_ADDRESS" json:"prospect_email_address"`
	ProspectWebsite            int       `gorm:"column: PROSPECT_WEBSITE" json:"prospect_website"`
	ProspectPhoneNo            int       `gorm:"column: PROSPECT_PHONE_NO" json:"prospect_phone_no"`
	ProspectFaxNo              int       `gorm:"column: PROSPECT_FAX_NO" json:"prospect_fax_no"`
	ProspectGender             int       `gorm:"column: PROSPECT_GENDER" json:"prospect_gender"`
	BizType                    int       `gorm:"column: BIZ_TYPE" json:"biz_type"`
	BizGroup                   int       `gorm:"column: BIZ_GROUP" json:"biz_group"`
	ContactPerson              int       `gorm:"column: CONTACT_PERSON" json:"contact_person"`
	ContactGender              int       `gorm:"column: CONTACT_GENDER" json:"contact_gender"`
	ContactJobTitle            int       `gorm:"column: CONTACT_JOB_TITLE" json:"contact_job_title"`
	ContactMobilePhone         int       `gorm:"column: CONTACT_MOBILE_PHONE" json:"contact_mobile_phone"`
	ContactEmailAddress        int       `gorm:"column: CONTACT_EMAIL_ADDRESS" json:"contact_email_address"`
	TestDrvDateSchedule        int       `gorm:"column: TEST_DRV_DATE_SCHEDULE" json:"test_drv_date_schedule"`
	TestDrvDateActual          int       `gorm:"column: TEST_DRV_DATE_ACTUAL" json:"test_drv_date_actual"`
	CompetitorModel            int       `gorm:"column: COMPETITOR_MODEL" json:"competitor_model"`
	OptionCode                 int       `gorm:"column: OPTION_CODE" json:"option_code"`
	CcCh                       int       `gorm:"column: CC_CH" json:"cc_ch"`
	ChP                        int       `gorm:"column: CH_P" json:"ch_p"`
	PHp                        int       `gorm:"column: P_HP" json:"p_hp"`
	HpDo                       int       `gorm:"column: HP_DO" json:"hp_do"`
	FollowUpNotes              int       `gorm:"column: FOLLOW_UP_NOTES" json:"follow_up_notes"`
	FollowUpResult             int       `gorm:"column: FOLLOW_UP_RESULT" json:"follow_up_result"`
	SpmStageDate               int       `gorm:"column: SPM_STAGE_DATE" json:"spm_stage_date"`
	SpmRemark                  int       `gorm:"column: SPM_REMARK" json:"spm_remark"`
	OrderByTaxRegNo            int       `gorm:"column: ORDER_BY_TAX_REG_NO" json:"order_by_tax_reg_no"`
	OrderByTaxRegDate          int       `gorm:"column: ORDER_BY_TAX_REG_DATE" json:"order_by_tax_reg_date"`
	OrderByTaxName             int       `gorm:"column: ORDER_BY_TAX_NAME" json:"order_by_tax_name"`
	OrderByTaxAddress1         int       `gorm:"column: ORDER_BY_TAX_ADDRESS_1" json:"order_by_tax_address_1"`
	OrderByTaxAddress2         int       `gorm:"column: ORDER_BY_TAX_ADDRESS_2" json:"order_by_tax_address_2"`
	OrderByTaxAddress3         int       `gorm:"column: ORDER_BY_TAX_ADDRESS_3" json:"order_by_tax_address_3"`
	OrderByTaxVillageCode      int       `gorm:"column: ORDER_BY_TAX_VILLAGE_CODE" json:"order_by_tax_village_code"`
	OrderByTaxSubdistrictCode  int       `gorm:"column: ORDER_BY_TAX_SUBDISTRICT_CODE" json:"order_by_tax_subdistrict_code"`
	OrderByTaxMunicipalityCode int       `gorm:"column: ORDER_BY_TAX_MUNICIPALITY_CODE" json:"order_by_tax_municipality_code"`
	OrderByTaxProvinceCode     int       `gorm:"column: ORDER_BY_TAX_PROVINCE_CODE" json:" order_by_tax_province_code"`
	OrderByTaxCityCode         int       `gorm:"column: ORDER_BY_TAX_CITY_CODE" json:"order_by_tax_city_code"`
	OrderByTaxZipCode          int       `gorm:"column: ORDER_BY_TAX_ZIP_CODE" json:"order_by_tax_zip_code"`
	PkpStatus                  int       `gorm:"column: PKP_STATUS" json:"pkp_status"`
	PkpDate                    int       `gorm:"column: PKP_DATE" json:"pkp_date"`
	PkpNo                      int       `gorm:"column: PKP_NO" json:"pkp_no"`
	PkpType                    int       `gorm:"column: PKP_TYPE" json:"pkp_type"`
	CorpBizType                int       `gorm:"column: CORP_BIZ_TYPE" json:"corp_biz_type"`
	CorpBizGroup               int       `gorm:"column: CORP_BIZ_GROUP" json:"corp_biz_group"`
	CorpWebSite                int       `gorm:"column: CORP_WEB_SITE" json:"corp_web_site"`
	CorpPoNo                   int       `gorm:"column: CORP_PO_NO" json:"corp_po_no"`
	CorpPoDate                 int       `gorm:"column: CORP_PO_DATE" json:"corp_po_date"`
	CorpContactName            int       `gorm:"column: CORP_CONTACT_NAME" json:"corp_contact_name"`
	CorpContactGender          int       `gorm:"column: CORP_CONTACT_GENDER" json:"corp_contact_gender"`
	CorpContactJobTitle        int       `gorm:"column: CORP_CONTACT_JOB_TITLE" json:"corp_contact_job_title"`
	CorpMobilePhone            int       `gorm:"column: CORP_MOBILE_PHONE" json:"corp_mobile_phone"`
	CorpEmailAddress           int       `gorm:"column: CORP_EMAIL_ADDRESS" json:"corp_email_address"`
	CorrespondencePrefix       int       `gorm:"column: CORRESPONDENCE_PREFIX" json:"correspondence_prefix"`
	CorrespondenceName         int       `gorm:"column: CORRESPONDENCE_NAME" json:"correspondence_name"`
	CorrespondenceSuffix       int       `gorm:"column: CORRESPONDENCE_SUFFIX" json:"correspondence_suffix"`
	CorrespondenceGender       int       `gorm:"column: CORRESPONDENCE_GENDER" json:"correspondence_gender"`
	CorrespondenceJobTitle     int       `gorm:"column: CORRESPONDENCE_JOB_TITLE" json:"correspondence_job_title"`
	CorrespondenceAddress      int       `gorm:"column: CORRESPONDENCE_ADDRESS" json:"correspondence_address"`
	CorrespondenceArea         int       `gorm:"column: CORRESPONDENCE_AREA" json:"correspondence_area"`
	CorrespondencePhoneNo      int       `gorm:"column: CORRESPONDENCE_PHONE_NO" json:"correspondence_phone_no"`
	CorrespondenceFaxNo        int       `gorm:"column: CORRESPONDENCE_FAX_NO" json:"correspondence_fax_no"`
	CorrespondenceMobilePhone  int       `gorm:"column: CORRESPONDENCE_MOBILE_PHONE" json:"correspondence_mobile_phone"`
	CorrespondenceEmailAddress int       `gorm:"column: CORRESPONDENCE_EMAIL_ADDRESS" json:"correspondence_email_address"`
	Remark                     int       `gorm:"column: REMARK" json:"remark"`
	UserBirthday               int       `gorm:"column: USER_BIRTHDAY" json:"user_birthday"`
	UserReligion               int       `gorm:"column: USER_RELIGION" json:"user_religion"`
	ChangeNo                   float64   `gorm:"column: CHANGE_NO" json:"change_no"`
	CreationUserId             string    `gorm:"column: CREATION_USER_ID" json:"creation_user_id"`
	CreationDatetime           time.Time `gorm:"column: CREATION_DATETIME" json:"creation_datetime"`
	ChangeUserId               string    `gorm:"column: CHANGE_USER_ID" json:"change_user_id"`
	ChangeDatetime             time.Time `gorm:"column: CHANGE_DATETIME" json:"change_datetime"`
}

func (UtPointProspecting) TableName() string {
	return "utPointProspecting"
}

type UTProspect0 struct {
	RecordStatus             string    `gorm:"column:RECORD_STATUS" json:"record_status"`
	CompanyCode              float64   `gorm:"column:COMPANY_CODE" json:"company_code"`
	ProspectSystemNo         float64   `gorm:"column:PROSPECT_SYSTEM_NO" json:"prospect_system_no"`
	ProspectStatus           string    `gorm:"column:PROSPECT_STATUS" json:"prospect_status"`
	ProspectDocNo            string    `gorm:"column:PROSPECT_DOC_NO" json:"prospect_doc_no"`
	ProspectDate             time.Time `gorm:"column:PROSPECT_DATE" json:"prospect_date"`
	ProspectStage            string    `gorm:"column:PROSPECT_STAGE" json:"prospect_stage"`
	ProspectTitlePrefix      string    `gorm:"column:PROSPECT_TITLE_PREFIX" json:"prospect_title_prefix"`
	ProspectName             string    `gorm:"column:PROSPECT_NAME" json:"prospect_name"`
	ProspectTitleSuffix      string    `gorm:"column:PROSPECT_TITLE_SUFFIX" json:"prospect_title_suffix"`
	ProspectSrcCode          string    `gorm:"column:PROSPECT_SRC_CODE" json:"prospect_src_code"`
	ProspectType             string    `gorm:"column:PROSPECT_TYPE" json:"prospect_type"`
	ProspectRef              string    `gorm:"column:PROSPECT_REF" json:"prospect_ref"`
	ProspectNotes            string    `gorm:"column:PROSPECT_NOTES" json:"prospect_notes"`
	ProspectCode             string    `gorm:"column:PROSPECT_CODE" json:"prospect_code"`
	CustomerCode             string    `gorm:"column:CUSTOMER_CODE" json:"customer_code"`
	BuyingBudget             float64   `gorm:"column:BUYING_BUDGET" json:"buying_budget"`
	BuyingPlan               time.Time `gorm:"column:BUYING_PLAN" json:"buying_plan"`
	FundType                 string    `gorm:"column:FUND_TYPE" json:"fund_type"`
	VehicleBrand             string    `gorm:"column:VEHICLE_BRAND" json:"vehicle_brand"`
	ModelCode                string    `gorm:"column:MODEL_CODE" json:"model_code"`
	VariantCode              string    `gorm:"column:VARIANT_CODE" json:"variant_code"`
	ProspectAddress1         string    `gorm:"column:PROSPECT_ADDRESS_1" json:"prospect_address1"`
	ProspectAddress2         string    `gorm:"column:PROSPECT_ADDRESS_2" json:"prospect_address2"`
	ProspectAddress3         string    `gorm:"column:PROSPECT_ADDRESS_3" json:"prospect_address3"`
	ProspectVillageCode      string    `gorm:"column:PROSPECT_VILLAGE_CODE" json:"prospect_village_code"`
	ProspectSubdistrictCode  string    `gorm:"column:PROSPECT_SUBDISTRICT_CODE" json:"prospect_subdistrict_code"`
	ProspectMunicipalityCode string    `gorm:"column:PROSPECT_MUNICIPALITY_CODE" json:"prospect_municipality_code"`
	ProspectProvinceCode     string    `gorm:"column:PROSPECT_PROVINCE_CODE" json:"prospect_province_code"`
	ProspectCityCode         string    `gorm:"column:PROSPECT_CITY_CODE" json:"prospect_city_code"`
	ProspectZipCode          string    `gorm:"column:PROSPECT_ZIP_CODE" json:"prospect_zip_code"`
	ProspectMobilePhone      string    `gorm:"column:PROSPECT_MOBILE_PHONE" json:"prospect_mobile_phone"`
	ProspectEmailAddress     string    `gorm:"column:PROSPECT_EMAIL_ADDRESS" json:"prospect_email_address"`
	ProspectWebsite          string    `gorm:"column:PROSPECT_WEBSITE" json:"prospect_website"`
	ProspectPhoneNo          string    `gorm:"column:PROSPECT_PHONE_NO" json:"prospect_phone_no"`
	ProspectFaxNo            string    `gorm:"column:PROSPECT_FAX_NO" json:"prospect_fax_no"`
	ProspectGender           string    `gorm:"column:PROSPECT_GENDER" json:"prospect_gender"`
	BizType                  string    `gorm:"column:BIZ_TYPE" json:"biz_type"`
	BizGroup                 string    `gorm:"column:BIZ_GROUP" json:"biz_group"`
	ContactPerson            string    `gorm:"column:CONTACT_PERSON" json:"contact_person"`
	ContactGender            string    `gorm:"column:CONTACT_GENDER" json:"contact_gender"`
	ContactJobTitle          string    `gorm:"column:CONTACT_JOB_TITLE" json:"contact_job_title"`
	ContactMobilePhone       string    `gorm:"column:CONTACT_MOBILE_PHONE" json:"contact_mobile_phone"`
	ContactEmailAddress      string    `gorm:"column:CONTACT_EMAIL_ADDRESS" json:"contact_email_address"`
	StageDateCc              time.Time `gorm:"column:STAGE_DATE_CC" json:"stage_date_cc"`
	StageDateCh              time.Time `gorm:"column:STAGE_DATE_CH" json:"stage_date_ch"`
	StageDateSpm             time.Time `gorm:"column:STAGE_DATE_SPM" json:"stage_date_spm"`
	SpmSysNo                 float64   `gorm:"column:SPM_SYS_NO" json:"spm_sys_no"`
	SpmDocNo                 string    `gorm:"column:SPM_DOC_NO" json:"spm_doc_no"`
	SalesRepCode1St          string    `gorm:"column:SALES_REP_CODE_1ST" json:"sales_rep_code_1st"`
	SalesRepCode             string    `gorm:"column:SALES_REP_CODE" json:"sales_rep_code"`
	ConductTestDrv           string    `gorm:"column:CONDUCT_TEST_DRV" json:"conduct_test_drv"`
	TestDrvDateSchedule      time.Time `gorm:"column:TEST_DRV_DATE_SCHEDULE" json:"test_drv_date_schedule"`
	TestDrvDateActual        time.Time `gorm:"column:TEST_DRV_DATE_ACTUAL" json:"test_drv_date_actual"`
	CompetitorModel          string    `gorm:"column:COMPETITOR_MODEL" json:"competitor_model"`
	DroppedDate              time.Time `gorm:"column:DROPPED_DATE" json:"dropped_date"`
	DroppedReasonCode        string    `gorm:"column:DROPPED_REASON_CODE" json:"dropped_reason_code"`
	DroppedRemark            string    `gorm:"column:DROPPED_REMARK" json:"dropped_remark"`
	TotalQty                 int       `gorm:"column:TOTAL_QTY" json:"total_qty"`
	Keyword                  string    `gorm:"column:KEYWORD" json:"keyword"`
	LastFollowUpDate         time.Time `gorm:"column:LAST_FOLLOW_UP_DATE" json:"last_follow_up_date"`
	ChangeNo                 float64   `gorm:"column:CHANGE_NO" json:"change_no"`
	CreationUserId           string    `gorm:"column:CREATION_USER_ID" json:"creation_user_id"`
	CreationDatetime         time.Time `gorm:"column:CREATION_DATETIME" json:"creation_datetime"`
	ChangeUserId             string    `gorm:"column:CHANGE_USER_ID" json:"change_user_id"`
	ChangeDatetime           time.Time `gorm:"column:CHANGE_DATETIME" json:"change_datetime"`
	IsFromMobile             bool      `gorm:"column:IS_FROM_MOBILE" json:"is_from_mobile"`
	OptionCode               string    `gorm:"column:OPTION_CODE" json:"option_code"`
	BussinessCategory        string    `gorm:"column:BUSSINESS_CATEGORY" json:"bussiness_category"`
	SubBussinessCategory     string    `gorm:"column:SUB_BUSSINESS_CATEGORY" json:"sub_bussiness_category"`
	Tenor                    float64   `gorm:"column:TENOR" json:"tenor"`
	Retailprice              float64   `gorm:"column:RETAILPRICE" json:"retailprice"`
	Requestdiscount          float64   `gorm:"column:REQUESTDISCOUNT" json:"requestdiscount"`
}

func (*UTProspect0) TableName() string {
	return "utProspect0"
}

type UTProspect1 struct {
	RecordStatus         string    `gorm:"column:RECORD_STATUS" json:"record_status"`
	ProspectSystemNo     float64   `gorm:"column:PROSPECT_SYSTEM_NO" json:"prospect_system_no"`
	ProspectLine         float64   `gorm:"column:PROSPECT_LINE" json:"prospect_line"`
	LineStatus           string    `gorm:"column:LINE_STATUS" json:"line_status"`
	ModelCode            string    `gorm:"column:MODEL_CODE" json:"model_code"`
	VariantCode          string    `gorm:"column:VARIANT_CODE" json:"variant_code"`
	ColourCode           string    `gorm:"column:COLOUR_CODE" json:"colour_code"`
	ProspectQty          float64   `gorm:"column:PROSPECT_QTY" json:"prospect_qty"`
	SpmQty               float64   `gorm:"column:SPM_QTY" json:"spm_qty"`
	CanceledQty          float64   `gorm:"column:CANCELED_QTY" json:"canceled_qty"`
	CanceledReason       string    `gorm:"column:CANCELED_REASON" json:"canceled_reason"`
	ExpectedDeliveryDate time.Time `gorm:"column:EXPECTED_DELIVERY_DATE" json:"expected_delivery_date"`
	ChangeNo             float64   `gorm:"column:CHANGE_NO" json:"change_no"`
	CreationUserId       string    `gorm:"column:CREATION_USER_ID" json:"creation_user_id"`
	CreationDatetime     time.Time `gorm:"column:CREATION_DATETIME" json:"creation_datetime"`
	ChangeUserId         string    `gorm:"column:CHANGE_USER_ID" json:"change_user_id"`
	ChangeDatetime       time.Time `gorm:"column:CHANGE_DATETIME" json:"change_datetime"`
	OptionCode           string    `gorm:"column:OPTION_CODE" json:"option_code"`
	OptionCodeEffDate    time.Time `gorm:"column:OPTION_CODE_EFF_DATE" json:"option_code_eff_date"`
}

func (*UTProspect1) TableName() string {
	return "utProspect1"
}

type UTProspect2 struct {
	RecordStatus     string    `gorm:"column:RECORD_STATUS" json:"record_status"`
	ProspectSystemNo float64   `gorm:"column:PROSPECT_SYSTEM_NO" json:"prospect_system_no"`
	ProspectLine     float64   `gorm:"column:PROSPECT_LINE" json:"prospect_line"`
	FollowUpCode     string    `gorm:"column:FOLLOW_UP_CODE" json:"follow_up_code"`
	FollowUpDate     time.Time `gorm:"column:FOLLOW_UP_DATE" json:"follow_up_date"`
	FollowUpNote     string    `gorm:"column:FOLLOW_UP_NOTE" json:"follow_up_note"`
	ResultDate       time.Time `gorm:"column:RESULT_DATE" json:"result_date"`
	ResultNote       string    `gorm:"column:RESULT_NOTE" json:"result_note"`
	ChangeNo         float64   `gorm:"column:CHANGE_NO" json:"change_no"`
	CreationUserId   string    `gorm:"column:CREATION_USER_ID" json:"creation_user_id"`
	CreationDatetime time.Time `gorm:"column:CREATION_DATETIME" json:"creation_datetime"`
	ChangeUserId     string    `gorm:"column:CHANGE_USER_ID" json:"change_user_id"`
	ChangeDatetime   time.Time `gorm:"column:CHANGE_DATETIME" json:"change_datetime"`
}

func (*UTProspect2) TableName() string {
	return "UT_PROSPECT2"
}

type UtSpm0 struct {
	RecordStatus               string    `gorm:"column:RECORD_STATUS" json:"record_status"`
	CompanyCode                float64   `gorm:"column:COMPANY_CODE" json:"company_code"`
	SpmSystemNo                float64   `gorm:"column:SPM_SYSTEM_NO" json:"spm_system_no"`
	SpmStatus                  string    `gorm:"column:SPM_STATUS" json:"spm_status"`
	SpmDocNo                   string    `gorm:"column:SPM_DOC_NO" json:"spm_doc_no"`
	SpmType                    string    `gorm:"column:SPM_TYPE" json:"spm_type"`
	SpmHighestLineStage        string    `gorm:"column:SPM_HIGHEST_LINE_STAGE" json:"spm_highest_line_stage"`
	SpmLowestLineStage         string    `gorm:"column:SPM_LOWEST_LINE_STAGE" json:"spm_lowest_line_stage"`
	SpmStageDate               time.Time `gorm:"column:SPM_STAGE_DATE" json:"spm_stage_date"`
	SpmRemark                  string    `gorm:"column:SPM_REMARK" json:"spm_remark"`
	OrderStatus                string    `gorm:"column:ORDER_STATUS" json:"order_status"`
	ProspectSystemNo           float64   `gorm:"column:PROSPECT_SYSTEM_NO" json:"prospect_system_no"`
	ProspectDocNo              string    `gorm:"column:PROSPECT_DOC_NO" json:"prospect_doc_no"`
	ProspectDate               time.Time `gorm:"column:PROSPECT_DATE" json:"prospect_date"`
	ProspectSrcCode            string    `gorm:"column:PROSPECT_SRC_CODE" json:"prospect_src_code"`
	UnitTrxType                string    `gorm:"column:UNIT_TRX_TYPE" json:"unit_trx_type"`
	TopCode                    string    `gorm:"column:TOP_CODE" json:"top_code"`
	SalesRepCode               string    `gorm:"column:SALES_REP_CODE" json:"sales_rep_code"`
	Keyword                    string    `gorm:"column:KEYWORD" json:"keyword"`
	CustomerType               string    `gorm:"column:CUSTOMER_TYPE" json:"customer_type"`
	CustomerClass              string    `gorm:"column:CUSTOMER_CLASS" json:"customer_class"`
	OrderByCustCode            string    `gorm:"column:ORDER_BY_CUST_CODE" json:"order_by_cust_code"`
	OrderByIdType              string    `gorm:"column:ORDER_BY_ID_TYPE" json:"order_by_id_type"`
	OrderByIdNo                string    `gorm:"column:ORDER_BY_ID_NO" json:"order_by_id_no"`
	OrderByPrefix              string    `gorm:"column:ORDER_BY_PREFIX" json:"order_by_prefix"`
	OrderByName                string    `gorm:"column:ORDER_BY_NAME" json:"order_by_name"`
	OrderBySuffix              string    `gorm:"column:ORDER_BY_SUFFIX" json:"order_by_suffix"`
	OrderByGender              string    `gorm:"column:ORDER_BY_GENDER" json:"order_by_gender"`
	OrderByAddress1            string    `gorm:"column:ORDER_BY_ADDRESS_1" json:"order_by_address1"`
	OrderByAddress2            string    `gorm:"column:ORDER_BY_ADDRESS_2" json:"order_by_address2"`
	OrderByAddress3            string    `gorm:"column:ORDER_BY_ADDRESS_3" json:"order_by_address3"`
	OrderByVillageCode         string    `gorm:"column:ORDER_BY_VILLAGE_CODE" json:"order_by_village_code"`
	OrderBySubistrictCode      string    `gorm:"column:ORDER_BY_SUBDISTRICT_CODE" json:"order_by_subdistrict_code"`
	OrderByMunicipalityCode    string    `gorm:"column:ORDER_BY_MUNICIPALITY_CODE" json:"order_by_municipality_code"`
	OrderByProvinceCode        string    `gorm:"column:ORDER_BY_PROVINCE_CODE" json:"order_by_province_code"`
	OrderByCityCode            string    `gorm:"column:ORDER_BY_CITY_CODE" json:"order_by_city_code"`
	OrderByZipCode             string    `gorm:"column:ORDER_BY_ZIP_CODE" json:"order_by_zip_code"`
	OrderByPhoneNo             string    `gorm:"column:ORDER_BY_PHONE_NO" json:"order_by_phone_no"`
	OrderByFaxNo               string    `gorm:"column:ORDER_BY_FAX_NO" json:"order_by_fax_no"`
	OrderByMobilePhone         string    `gorm:"column:ORDER_BY_MOBILE_PHONE" json:"order_by_mobile_phone"`
	OrderByEmailAddress        string    `gorm:"column:ORDER_BY_EMAIL_ADDRESS" json:"order_by_email_address"`
	OrderByTaxInvType          string    `gorm:"column:ORDER_BY_TAX_INV_TYPE" json:"order_by_tax_inv_type"`
	OrderByTaxRegNo            string    `gorm:"column:ORDER_BY_TAX_REG_NO" json:"order_by_tax_reg_no"`
	OrderByTaxRegDate          time.Time `gorm:"column:ORDER_BY_TAX_REG_DATE" json:"order_by_tax_reg_date"`
	OrderByTaxName             string    `gorm:"column:ORDER_BY_TAX_NAME" json:"order_by_tax_name"`
	OrderByTaxAddress1         string    `gorm:"column:ORDER_BY_TAX_ADDRESS_1" json:"order_by_tax_address1"`
	OrderByTaxAddress2         string    `gorm:"column:ORDER_BY_TAX_ADDRESS_2" json:"order_by_tax_address2"`
	OrderByTaxAddress3         string    `gorm:"column:ORDER_BY_TAX_ADDRESS_3" json:"order_by_tax_address3"`
	OrderByTaxVillageCode      string    `gorm:"column:ORDER_BY_TAX_VILLAGE_CODE" json:"order_by_tax_village_code"`
	OrderByTaxSubdistrictCode  string    `gorm:"column:ORDER_BY_TAX_SUBDISTRICT_CODE" json:"order_by_tax_subdistrict_code"`
	OrderByTaxMunicipalityCode string    `gorm:"column:ORDER_BY_TAX_MUNICIPALITY_CODE" json:"order_by_tax_municipality_code"`
	OrderByTaxProvinceCode     string    `gorm:"column:ORDER_BY_TAX_PROVINCE_CODE" json:"order_by_tax_province_code"`
	OrderByTaxCityCode         string    `gorm:"column:ORDER_BY_TAX_CITY_CODE" json:"order_by_tax_city_code"`
	OrderByTaxZipCode          string    `gorm:"column:ORDER_BY_TAX_ZIP_CODE" json:"order_by_tax_zip_code"`
	PkpStatus                  bool      `gorm:"column:PKP_STATUS" json:"pkp_status"`
	PkpDate                    time.Time `gorm:"column:PKP_DATE" json:"pkp_date"`
	PkpNo                      string    `gorm:"column:PKP_NO" json:"pkp_no"`
	PkpType                    string    `gorm:"column:PKP_TYPE" json:"pkp_type"`
	CorpBizType                string    `gorm:"column:CORP_BIZ_TYPE" json:"corp_biz_type"`
	CorpBizGroup               string    `gorm:"column:CORP_BIZ_GROUP" json:"corp_biz_group"`
	CorpWebSite                string    `gorm:"column:CORP_WEB_SITE" json:"corp_web_site"`
	CorpPoNo                   string    `gorm:"column:CORP_PO_NO" json:"corp_po_no"`
	CorpPoDate                 time.Time `gorm:"column:CORP_PO_DATE" json:"corp_po_date"`
	CorpContactName            string    `gorm:"column:CORP_CONTACT_NAME" json:"corp_contact_name"`
	CorpContactGender          string    `gorm:"column:CORP_CONTACT_GENDER" json:"corp_contact_gender"`
	CorpContactJobTitle        string    `gorm:"column:CORP_CONTACT_JOB_TITLE" json:"corp_contact_job_title"`
	CorpMobilePhone            string    `gorm:"column:CORP_MOBILE_PHONE" json:"corp_mobile_phone"`
	CorpEmailAddress           string    `gorm:"column:CORP_EMAIL_ADDRESS" json:"corp_email_address"`
	StnkAddressEqual           string    `gorm:"column:STNK_ADDRESS_EQUAL" json:"stnk_address_equal"`
	StnkCustCode               string    `gorm:"column:STNK_CUST_CODE" json:"stnk_cust_code"`
	StnkIdType                 string    `gorm:"column:STNK_ID_TYPE" json:"stnk_id_type"`
	StnkIdNo                   string    `gorm:"column:STNK_ID_NO" json:"stnk_id_no"`
	StnkPrefix                 string    `gorm:"column:STNK_PREFIX" json:"stnk_prefix"`
	StnkName                   string    `gorm:"column:STNK_NAME" json:"stnk_name"`
	StnkSuffix                 string    `gorm:"column:STNK_SUFFIX" json:"stnk_suffix"`
	StnkGender                 string    `gorm:"column:STNK_GENDER" json:"stnk_gender"`
	StnkAddress1               string    `gorm:"column:STNK_ADDRESS_1" json:"stnk_address1"`
	StnkAddress2               string    `gorm:"column:STNK_ADDRESS_2" json:"stnk_address2"`
	StnkAddress3               string    `gorm:"column:STNK_ADDRESS_3" json:"stnk_address3"`
	StnkVillageCode            string    `gorm:"column:STNK_VILLAGE_CODE" json:"stnk_village_code"`
	StnkSubdistrictCode        string    `gorm:"column:STNK_SUBDISTRICT_CODE" json:"stnk_subdistrict_code"`
	StnkMunicipalityCode       string    `gorm:"column:STNK_MUNICIPALITY_CODE" json:"stnk_municipality_code"`
	StnkProvinceCode           string    `gorm:"column:STNK_PROVINCE_CODE" json:"stnk_province_code"`
	StnkCityCode               string    `gorm:"column:STNK_CITY_CODE" json:"stnk_city_code"`
	StnkZipCode                string    `gorm:"column:STNK_ZIP_CODE" json:"stnk_zip_code"`
	StnkPhoneNo                string    `gorm:"column:STNK_PHONE_NO" json:"stnk_phone_no"`
	StnkFaxNo                  string    `gorm:"column:STNK_FAX_NO" json:"stnk_fax_no"`
	StnkMobilePhone            string    `gorm:"column:STNK_MOBILE_PHONE" json:"stnk_mobile_phone"`
	StnkEmailAddress           string    `gorm:"column:STNK_EMAIL_ADDRESS" json:"stnk_email_address"`
	CorrAddressEqual           string    `gorm:"column:CORR_ADDRESS_EQUAL" json:"corr_address_equal"`
	CorrPrefix                 string    `gorm:"column:CORR_PREFIX" json:"corr_prefix"`
	CorrName                   string    `gorm:"column:CORR_NAME" json:"corr_name"`
	CorrSuffix                 string    `gorm:"column:CORR_SUFFIX" json:"corr_suffix"`
	CorrGender                 string    `gorm:"column:CORR_GENDER" json:"corr_gender"`
	CorrJobTitle               string    `gorm:"column:CORR_JOB_TITLE" json:"corr_job_title"`
	CorrAddress1               string    `gorm:"column:CORR_ADDRESS_1" json:"corr_address1"`
	CorrAddress2               string    `gorm:"column:CORR_ADDRESS_2" json:"corr_address2"`
	CorrAddress3               string    `gorm:"column:CORR_ADDRESS_3" json:"corr_address3"`
	CorrVillageCode            string    `gorm:"column:CORR_VILLAGE_CODE" json:"corr_village_code"`
	CorrSubdistrictCode        string    `gorm:"column:CORR_SUBDISTRICT_CODE" json:"corr_subdistrict_code"`
	CorrMunicipalityCode       string    `gorm:"column:CORR_MUNICIPALITY_CODE" json:"corr_municipality_code"`
	CorrProvinceCode           string    `gorm:"column:CORR_PROVINCE_CODE" json:"corr_province_code"`
	CorrCityCode               string    `gorm:"column:CORR_CITY_CODE" json:"corr_city_code"`
	CorrZipCode                string    `gorm:"column:CORR_ZIP_CODE" json:"corr_zip_code"`
	CorrPhoneNo                string    `gorm:"column:CORR_PHONE_NO" json:"corr_phone_no"`
	CorrFaxNo                  string    `gorm:"column:CORR_FAX_NO" json:"corr_fax_no"`
	CorrMobilePhone            string    `gorm:"column:CORR_MOBILE_PHONE" json:"corr_mobile_phone"`
	CorrEmailAddress           string    `gorm:"column:CORR_EMAIL_ADDRESS" json:"corr_email_address"`
	VatCode                    string    `gorm:"column:VAT_CODE" json:"vat_code"`
	VarPercentage              float64   `gorm:"column:VAR_PERCENTAGE" json:"var_percentage"`
	VehicleBrand               string    `gorm:"column:VEHICLE_BRAND" json:"vehicle_brand"`
	ModelCode                  string    `gorm:"column:MODEL_CODE" json:"model_code"`
	VariantCode                string    `gorm:"column:VARIANT_CODE" json:"variant_code"`
	PriceCode                  string    `gorm:"column:PRICE_CODE" json:"price_code"`
	PriceEffDate               time.Time `gorm:"column:PRICE_EFF_DATE" json:"price_eff_date"`
	SubmissionNo               float64   `gorm:"column:SUBMISSION_NO" json:"submission_no"`
	AtpmOrderNo                string    `gorm:"column:ATPM_ORDER_NO" json:"atpm_order_no"`
	AtpmOrderDate              time.Time `gorm:"column:ATPM_ORDER_DATE" json:"atpm_order_date"`
	AtpmPoReadBy               string    `gorm:"column:ATPM_PO_READ_BY" json:"atpm_po_read_by"`
	AtpmPoReadDate             time.Time `gorm:"column:ATPM_PO_READ_DATE" json:"atpm_po_read_date"`
	AtpmPoReplyBy              string    `gorm:"column:ATPM_PO_REPLY_BY" json:"atpm_po_reply_by"`
	AtpmPoReplyDate            time.Time `gorm:"column:ATPM_PO_REPLY_DATE" json:"atpm_po_reply_date"`
	AtpmPoReplyNotes           string    `gorm:"column:ATPM_PO_REPLY_NOTES" json:"atpm_po_reply_notes"`
	AtpmIndentBookingNo        string    `gorm:"column:ATPM_INDENT_BOOKING_NO" json:"atpm_indent_booking_no"`
	DocApprovalStatus          string    `gorm:"column:DOC_APPROVAL_STATUS" json:"doc_approval_status"`
	DocEdittingDate            time.Time `gorm:"column:DOC_EDITTING_DATE" json:"doc_editting_date"`
	DocEdittingBy              string    `gorm:"column:DOC_EDITTING_BY" json:"doc_editting_by"`
	IncentiveFleet             bool      `gorm:"column:INCENTIVE_FLEET" json:"incentive_fleet"`
	ChangeNo                   float64   `gorm:"column:CHANGE_NO" json:"change_no"`
	CreationUserId             string    `gorm:"column:CREATION_USER_ID" json:"creation_user_id"`
	CreationDatetime           time.Time `gorm:"column:CREATION_DATETIME" json:"creation_datetime"`
	ChangeUserId               string    `gorm:"column:CHANGE_USER_ID" json:"change_user_id"`
	ChangeDatetime             time.Time `gorm:"column:CHANGE_DATETIME" json:"change_datetime"`
	Remark                     string    `gorm:"column:REMARK" json:"remark"`
	CostCenterCode             string    `gorm:"column:COST_CENTER_CODE" json:"cost_center_code"`
	CorpPoLineNo               float64   `gorm:"column:CORP_PO_LINE_NO" json:"corp_po_line_no"`
	IsFromMobile               bool      `gorm:"column:IS_FROM_MOBILE" json:"is_from_mobile"`
	SalesMitra                 string    `gorm:"column:SALES_MITRA" json:"sales_mitra"`
	WarnaPlat                  string    `gorm:"column:WARNA_PLAT" json:"warna_plat"`
	BodyType                   string    `gorm:"column:BODY_TYPE" json:"body_type"`
	KeySysNo                   string    `gorm:"column:KEY_SYS_NO" json:"key_sys_no"`
	HassSysNo                  string    `gorm:"column:HASS_SYS_NO" json:"hass_sys_no"`
	ConfirmDate                time.Time `gorm:"column:CONFIRM_DATE" json:"confirm_date"`
	BuyingType                 string    `gorm:"column:BUYING_TYPE" json:"buying_type"`
	Nitku                      string    `gorm:"column:NITKU" json:"nitku"`
}

func (u0 *UtSpm0) TableName() string {
	return "utSPM0"
}

type UtSpm1 struct {
	RecordStatus               string    `gorm:"column:RECORD_STATUS" json:"record_status"`
	CompanyCode                float64   `gorm:"column:COMPANY_CODE" json:"company_code"`
	SpmSystemNo                float64   `gorm:"column:SPM_SYSTEM_NO" json:"spm_system_no"`
	SpmLineNo                  float64   `gorm:"column:SPM_LINE_NO" json:"spm_line_no"`
	SpmLineStatus              string    `gorm:"column:SPM_LINE_STATUS" json:"spm_line_status"`
	SpmStageDate               time.Time `gorm:"column:SPM_STAGE_DATE" json:"spm_stage_date"`
	OrderStatus                string    `gorm:"column:ORDER_STATUS" json:"order_status"`
	VehicleBrand               string    `gorm:"column:VEHICLE_BRAND" json:"vehicle_brand"`
	ModelCode                  string    `gorm:"column:MODEL_CODE" json:"model_code"`
	VariantCode                string    `gorm:"column:VARIANT_CODE" json:"variant_code"`
	ColourCode                 string    `gorm:"column:COLOUR_CODE" json:"colour_code"`
	UserAddressEqual           string    `gorm:"column:USER_ADDRESS_EQUAL" json:"user_address_equal"`
	UserCustCode               string    `gorm:"column:USER_CUST_CODE" json:"user_cust_code"`
	UserTitlePrefix            string    `gorm:"column:USER_TITLE_PREFIX" json:"user_title_prefix"`
	UserName                   string    `gorm:"column:USER_NAME" json:"user_name"`
	UserTitleSuffix            string    `gorm:"column:USER_TITLE_SUFFIX" json:"user_title_suffix"`
	UserIdType                 string    `gorm:"column:USER_ID_TYPE" json:"user_id_type"`
	UserIdNo                   string    `gorm:"column:USER_ID_NO" json:"user_id_no"`
	UserAddress1               string    `gorm:"column:USER_ADDRESS_1" json:"user_address1"`
	UserAddress2               string    `gorm:"column:USER_ADDRESS_2" json:"user_address2"`
	UserAddress3               string    `gorm:"column:USER_ADDRESS_3" json:"user_address3"`
	UserVillageCode            string    `gorm:"column:USER_VILLAGE_CODE" json:"user_village_code"`
	UserSubdistrictCode        string    `gorm:"column:USER_SUBDISTRICT_CODE" json:"user_subdistrict_code"`
	UserMunicipalityCode       string    `gorm:"column:USER_MUNICIPALITY_CODE" json:"user_municipality_code"`
	UserProvinceCode           string    `gorm:"column:USER_PROVINCE_CODE" json:"user_province_code"`
	UserCityCode               string    `gorm:"column:USER_CITY_CODE" json:"user_city_code"`
	UserZipCode                string    `gorm:"column:USER_ZIP_CODE" json:"user_zip_code"`
	UserPhoneNo                string    `gorm:"column:USER_PHONE_NO" json:"user_phone_no"`
	UserFaxNo                  string    `gorm:"column:USER_FAX_NO" json:"user_fax_no"`
	UserMobilePhone            string    `gorm:"column:USER_MOBILE_PHONE" json:"user_mobile_phone"`
	UserEmailAddress           string    `gorm:"column:USER_EMAIL_ADDRESS" json:"user_email_address"`
	UserJobTitle               string    `gorm:"column:USER_JOB_TITLE" json:"user_job_title"`
	UserGender                 string    `gorm:"column:USER_GENDER" json:"user_gender"`
	UserReligion               string    `gorm:"column:USER_RELIGION" json:"user_religion"`
	UserHobby                  string    `gorm:"column:USER_HOBBY" json:"user_hobby"`
	UserBirthdayPlace          string    `gorm:"column:USER_BIRTHDAY_PLACE" json:"user_birthday_place"`
	UserBirthday               time.Time `gorm:"column:USER_BIRTHDAY" json:"user_birthday"`
	UnitPrice                  float64   `gorm:"column:UNIT_PRICE" json:"unit_price"`
	ApUnitPrice                float64   `gorm:"column:AP_UNIT_PRICE" json:"ap_unit_price"`
	ApVariance                 float64   `gorm:"column:AP_VARIANCE" json:"ap_variance"`
	TransportAmount            float64   `gorm:"column:TRANSPORT_AMOUNT" json:"transport_amount"`
	TransportInsAmount         float64   `gorm:"column:TRANSPORT_INS_AMOUNT" json:"transport_ins_amount"`
	WindowFilmAmount           float64   `gorm:"column:WINDOW_FILM_AMOUNT" json:"window_film_amount"`
	Margin01Amount             float64   `gorm:"column:MARGIN01_AMOUNT" json:"margin_01_amount"`
	Margin02Amount             float64   `gorm:"column:MARGIN02_AMOUNT" json:"margin_02_amount"`
	OfftrPriceOther01          float64   `gorm:"column:OFFTR_PRICE_OTHER01" json:"offtr_price_other01"`
	OfftrPriceOther02          float64   `gorm:"column:OFFTR_PRICE_OTHER02" json:"offtr_price_other02"`
	OfftrPriceOther03          float64   `gorm:"column:OFFTR_PRICE_OTHER03" json:"offtr_price_other03"`
	OfftrPriceOther04          float64   `gorm:"column:OFFTR_PRICE_OTHER04" json:"offtr_price_other04"`
	OfftrPriceOther05          float64   `gorm:"column:OFFTR_PRICE_OTHER05" json:"offtr_price_other05"`
	OfftrPrice                 float64   `gorm:"column:OFFTR_PRICE" json:"offtr_price"`
	OfftrPriceDppOriginal      float64   `gorm:"column:OFFTR_PRICE_DPP_ORIGINAL" json:"offtr_price_dpp_original"`
	OfftrPriceDpp              float64   `gorm:"column:OFFTR_PRICE_DPP" json:"offtr_price_dpp"`
	CashDisc                   float64   `gorm:"column:CASH_DISC" json:"cash_disc"`
	CashDp                     float64   `gorm:"column:CASH_DP" json:"cash_dp"`
	BbnAmount                  float64   `gorm:"column:BBN_AMOUNT" json:"bbn_amount"`
	OntrPrice                  float64   `gorm:"column:ONTR_PRICE" json:"ontr_price"`
	PriceOntrOther01           float64   `gorm:"column:PRICE_ONTR_OTHER01" json:"price_ontr_other01"`
	PriceOntrOther02           float64   `gorm:"column:PRICE_ONTR_OTHER02" json:"price_ontr_other02"`
	PriceOntrOther03           float64   `gorm:"column:PRICE_ONTR_OTHER03" json:"price_ontr_other03"`
	PriceOntrOther04           float64   `gorm:"column:PRICE_ONTR_OTHER04" json:"price_ontr_other04"`
	PriceOntrOther05           float64   `gorm:"column:PRICE_ONTR_OTHER05" json:"price_ontr_other05"`
	AddChargesAccYear          float64   `gorm:"column:ADD_CHARGES_ACC_YEAR" json:"add_charges_acc_year"`
	AddChargesChooseNo         float64   `gorm:"column:ADD_CHARGES_CHOOSE_NO" json:"add_charges_choose_no"`
	AddChargesCrossArea        float64   `gorm:"column:ADD_CHARGES_CROSS_AREA" json:"add_charges_cross_area"`
	AddChargesTransport        float64   `gorm:"column:ADD_CHARGES_TRANSPORT" json:"add_charges_transport"`
	AddChargesRsv01            float64   `gorm:"column:ADD_CHARGES_RSV01" json:"add_charges_rsv01"`
	AddChargesRsv02            float64   `gorm:"column:ADD_CHARGES_RSV02" json:"add_charges_rsv02"`
	AddChargesRsv03            float64   `gorm:"column:ADD_CHARGES_RSV03" json:"add_charges_rsv03"`
	AddChargesRsv04            float64   `gorm:"column:ADD_CHARGES_RSV04" json:"add_charges_rsv04"`
	AddChargesRsv05            float64   `gorm:"column:ADD_CHARGES_RSV05" json:"add_charges_rsv05"`
	AddDealerCostBbn           bool      `gorm:"column:ADD_DEALER_COST_BBN" json:"add_dealer_cost_bbn"`
	AddDealerCostAccYear       bool      `gorm:"column:ADD_DEALER_COST_ACC_YEAR" json:"add_dealer_cost_acc_year"`
	AddDealerCostChooseNo      bool      `gorm:"column:ADD_DEALER_COST_CHOOSE_NO" json:"add_dealer_cost_choose_no"`
	AddDealerCostCrossArea     bool      `gorm:"column:ADD_DEALER_COST_CROSS_AREA" json:"add_dealer_cost_cross_area"`
	AddDealerCostTransport     bool      `gorm:"column:ADD_DEALER_COST_TRANSPORT" json:"add_dealer_cost_transport"`
	AddDealerCostRsv01         bool      `gorm:"column:ADD_DEALER_COST_RSV01" json:"add_dealer_cost_rsv01"`
	AddDealerCostRsv02         bool      `gorm:"column:ADD_DEALER_COST_RSV02" json:"add_dealer_cost_rsv02"`
	AddDealerCostRsv03         bool      `gorm:"column:ADD_DEALER_COST_RSV03" json:"add_dealer_cost_rsv03"`
	AddDealerCostRsv04         bool      `gorm:"column:ADD_DEALER_COST_RSV04" json:"add_dealer_cost_rsv04"`
	AddDealerCostRsv05         bool      `gorm:"column:ADD_DEALER_COST_RSV05" json:"add_dealer_cost_rsv05"`
	InsuranceType              string    `gorm:"column:INSURANCE_TYPE" json:"insurance_type"`
	SubmissionNo               float64   `gorm:"column:SUBMISSION_NO" json:"submission_no"`
	ProspectSystemNo           float64   `gorm:"column:PROSPECT_SYSTEM_NO" json:"prospect_system_no"`
	ProspectLine               float64   `gorm:"column:PROSPECT_LINE" json:"prospect_line"`
	RefundCode                 string    `gorm:"column:REFUND_CODE" json:"refund_code"`
	MediatorCode               string    `gorm:"column:MEDIATOR_CODE" json:"mediator_code"`
	MediatorFeeAmount          float64   `gorm:"column:MEDIATOR_FEE_AMOUNT" json:"mediator_fee_amount"`
	MediatorIvrSysNo           float64   `gorm:"column:MEDIATOR_IVR_SYS_NO" json:"mediator_ivr_sys_no"`
	FundType                   string    `gorm:"column:FUND_TYPE" json:"fund_type"`
	PriceCode                  string    `gorm:"column:PRICE_CODE" json:"price_code"`
	SalPriceEffDate            time.Time `gorm:"column:SAL_PRICE_EFF_DATE" json:"sal_price_eff_date"`
	PurcPriceEffDate           time.Time `gorm:"column:PURC_PRICE_EFF_DATE" json:"purc_price_eff_date"`
	VatAmount                  float64   `gorm:"column:VAT_AMOUNT" json:"vat_amount"`
	AddAccsGross               float64   `gorm:"column:ADD_ACCS_GROSS" json:"add_accs_gross"`
	AddAccsDisc                float64   `gorm:"column:ADD_ACCS_DISC" json:"add_accs_disc"`
	AddAccsVat                 float64   `gorm:"column:ADD_ACCS_VAT" json:"add_accs_vat"`
	InsuranceSupplierCode      string    `gorm:"column:INSURANCE_SUPPLIER_CODE" json:"insurance_supplier_code"`
	InsuranceAmount            float64   `gorm:"column:INSURANCE_AMOUNT" json:"insurance_amount"`
	FreeAccsProg               float64   `gorm:"column:FREE_ACCS_PROG" json:"free_accs_prog"`
	FreeAccsDeal               float64   `gorm:"column:FREE_ACCS_DEAL" json:"free_accs_deal"`
	DpPayment                  float64   `gorm:"column:DP_PAYMENT" json:"dp_payment"`
	DpPaymentAllocated         float64   `gorm:"column:DP_PAYMENT_ALLOCATED" json:"dp_payment_allocated"`
	DpPaymentVat               float64   `gorm:"column:DP_PAYMENT_VAT" json:"dp_payment_vat"`
	DrawdownAmount             float64   `gorm:"column:DRAWDOWN_AMOUNT" json:"drawdown_amount"`
	LeasingSupplierCode        string    `gorm:"column:LEASING_SUPPLIER_CODE" json:"leasing_supplier_code"`
	LeasingDpPercentage        float64   `gorm:"column:LEASING_DP_PERCENTAGE" json:"leasing_dp_percentage"`
	LeasingDpAmount            float64   `gorm:"column:LEASING_DP_AMOUNT" json:"leasing_dp_amount"`
	LeasingTenorInMonth        float64   `gorm:"column:LEASING_TENOR_IN_MONTH" json:"leasing_tenor_in_month"`
	LeasingInterestFlatPerYear float64   `gorm:"column:LEASING_INTEREST_FLAT_PER_YEAR" json:"leasing_interest_flat_per_year"`
	LeasingPaymentTermCode     string    `gorm:"column:LEASING_PAYMENT_TERM_CODE" json:"leasing_payment_term_code"`
	LeasingPoNo                string    `gorm:"column:LEASING_PO_NO" json:"leasing_po_no"`
	LeasingPoDate              time.Time `gorm:"column:LEASING_PO_DATE" json:"leasing_po_date"`
	RefundBySystem             float64   `gorm:"column:REFUND_BY_SYSTEM" json:"refund_by_system"`
	RefundFromThirdParty       float64   `gorm:"column:REFUND_FROM_THIRD_PARTY" json:"refund_from_third_party"`
	RefundPayment              float64   `gorm:"column:REFUND_PAYMENT" json:"refund_payment"`
	RefundPaymentAllocated     float64   `gorm:"column:REFUND_PAYMENT_ALLOCATED" json:"refund_payment_allocated"`
	SubsidiFromAtpm            float64   `gorm:"column:SUBSIDI_FROM_ATPM" json:"subsidi_from_atpm"`
	LastStageStatus            string    `gorm:"column:LAST_STAGE_STATUS" json:"last_stage_status"`
	LastStageDate              time.Time `gorm:"column:LAST_STAGE_DATE" json:"last_stage_date"`
	StageCcDate                time.Time `gorm:"column:STAGE_CC_DATE" json:"stage_cc_date"`
	StageChDate                time.Time `gorm:"column:STAGE_CH_DATE" json:"stage_ch_date"`
	StagePDate                 time.Time `gorm:"column:STAGE_P_DATE" json:"stage_p_date"`
	StagePoDate                time.Time `gorm:"column:STAGE_PO_DATE" json:"stage_po_date"`
	StageGrpoDate              time.Time `gorm:"column:STAGE_GRPO_DATE" json:"stage_grpo_date"`
	StageHpDate                time.Time `gorm:"column:STAGE_HP_DATE" json:"stage_hp_date"`
	StageDoDate                time.Time `gorm:"column:STAGE_DO_DATE" json:"stage_do_date"`
	StageFakPolDate            time.Time `gorm:"column:STAGE_FAK_POL_DATE" json:"stage_fak_pol_date"`
	StageBbnPoIssue            time.Time `gorm:"column:STAGE_BBN_PO_ISSUE" json:"stage_bbn_po_issue"`
	StageStnkReceipt           time.Time `gorm:"column:STAGE_STNK_RECEIPT" json:"stage_stnk_receipt"`
	StagePdi                   time.Time `gorm:"column:STAGE_PDI" json:"stage_pdi"`
	StageBpk                   time.Time `gorm:"column:STAGE_BPK" json:"stage_bpk"`
	StageBpkbReceipt           time.Time `gorm:"column:STAGE_BPKB_RECEIPT" json:"stage_bpkb_receipt"`
	StageBpkbIssue             time.Time `gorm:"column:STAGE_BPKB_ISSUE" json:"stage_bpkb_issue"`
	AllocationChassisNo        string    `gorm:"column:ALLOCATION_CHASSIS_NO" json:"allocation_chassis_no"`
	AllocationEngineNo         string    `gorm:"column:ALLOCATION_ENGINE_NO" json:"allocation_engine_no"`
	AllocationBy               string    `gorm:"column:ALLOCATION_BY" json:"allocation_by"`
	AllocationDate             time.Time `gorm:"column:ALLOCATION_DATE" json:"allocation_date"`
	MinDp                      float64   `gorm:"column:MIN_DP" json:"min_dp"`
	MinDpAmount                float64   `gorm:"column:MIN_DP_AMOUNT" json:"min_dp_amount"`
	BookingFee                 float64   `gorm:"column:BOOKING_FEE" json:"booking_fee"`
	TotalCollection            float64   `gorm:"column:TOTAL_COLLECTION" json:"total_collection"`
	TotalAddChargesCustomer    float64   `gorm:"column:TOTAL_ADD_CHARGES_CUSTOMER" json:"total_add_charges_customer"`
	TotalAddChargesDealer      float64   `gorm:"column:TOTAL_ADD_CHARGES_DEALER" json:"total_add_charges_dealer"`
	AddChargesInvSysNo         float64   `gorm:"column:ADD_CHARGES_INV_SYS_NO" json:"add_charges_inv_sys_no"`
	AddChargesInvDocNo         string    `gorm:"column:ADD_CHARGES_INV_DOC_NO" json:"add_charges_inv_doc_no"`
	SalesInvoiceDate           time.Time `gorm:"column:SALES_INVOICE_DATE" json:"sales_invoice_date"`
	SalesInvoiceNo             string    `gorm:"column:SALES_INVOICE_NO" json:"sales_invoice_no"`
	SalesInvoiceSysNo          float64   `gorm:"column:SALES_INVOICE_SYS_NO" json:"sales_invoice_sys_no"`
	SubsidiInvDocNo            string    `gorm:"column:SUBSIDI_INV_DOC_NO" json:"subsidi_inv_doc_no"`
	SubsidiInvSysNo            float64   `gorm:"column:SUBSIDI_INV_SYS_NO" json:"subsidi_inv_sys_no"`
	ExhiInvDocNo               string    `gorm:"column:EXHI_INV_DOC_NO" json:"exhi_inv_doc_no"`
	ExhiInvSysNo               string    `gorm:"column:EXHI_INV_SYS_NO" json:"exhi_inv_sys_no"`
	IndentBookingNo            string    `gorm:"column:INDENT_BOOKING_NO" json:"indent_booking_no"`
	IndentBookingValidation    float64   `gorm:"column:INDENT_BOOKING_VALIDATION" json:"indent_booking_validation"`
	FakPolNo                   string    `gorm:"column:FAK_POL_NO" json:"fak_pol_no"`
	AtpmFakPolReqNo            string    `gorm:"column:ATPM_FAK_POL_REQ_NO" json:"atpm_fak_pol_req_no"`
	AtpmFakPolReqType          string    `gorm:"column:ATPM_FAK_POL_REQ_TYPE" json:"atpm_fak_pol_req_type"`
	AtpmOrderNo                string    `gorm:"column:ATPM_ORDER_NO" json:"atpm_order_no"`
	AtpmOrderLine              float64   `gorm:"column:ATPM_ORDER_LINE" json:"atpm_order_line"`
	AtpmFakPolDate             time.Time `gorm:"column:ATPM_FAK_POL_DATE" json:"atpm_fak_pol_date"`
	AtpmRetailSalesDate        time.Time `gorm:"column:ATPM_RETAIL_SALES_DATE" json:"atpm_retail_sales_date"`
	SupplyMethod               float64   `gorm:"column:SUPPLY_METHOD" json:"supply_method"`
	PoSystemNo                 float64   `gorm:"column:PO_SYSTEM_NO" json:"po_system_no"`
	PoLineNo                   float64   `gorm:"column:PO_LINE_NO" json:"po_line_no"`
	PoNoUnit                   string    `gorm:"column:PO_NO_UNIT" json:"po_no_unit"`
	PoNoTransport              string    `gorm:"column:PO_NO_TRANSPORT" json:"po_no_transport"`
	IncPeriodMonth             string    `gorm:"column:INC_PERIOD_MONTH" json:"inc_period_month"`
	IncPeriodYear              string    `gorm:"column:INC_PERIOD_YEAR" json:"inc_period_year"`
	ReqProdYear                string    `gorm:"column:REQ_PROD_YEAR" json:"req_prod_year"`
	Vat                        bool      `gorm:"column:VAT" json:"vat"`
	IncProcess                 bool      `gorm:"column:INC_PROCESS" json:"inc_process"`
	TrfType                    string    `gorm:"column:TRF_TYPE" json:"trf_type"`
	TrfDriverCode              string    `gorm:"column:TRF_DRIVER_CODE" json:"trf_driver_code"`
	TrfSupplierCode            string    `gorm:"column:TRF_SUPPLIER_CODE" json:"trf_supplier_code"`
	TrfAmount                  float64   `gorm:"column:TRF_AMOUNT" json:"trf_amount"`
	FaStandardAmount           float64   `gorm:"column:FA_STANDARD_AMOUNT" json:"fa_standard_amount"`
	FaChoiceCashAmount         float64   `gorm:"column:FA_CHOICE_CASH_AMOUNT" json:"fa_choice_cash_amount"`
	FaChoiceAccsAmount         float64   `gorm:"column:FA_CHOICE_ACCS_AMOUNT" json:"fa_choice_accs_amount"`
	FaDealAmount               float64   `gorm:"column:FA_DEAL_AMOUNT" json:"fa_deal_amount"`
	ApprovalAmount             float64   `gorm:"column:APPROVAL_AMOUNT" json:"approval_amount"`
	ExpDeliveryDate            time.Time `gorm:"column:EXP_DELIVERY_DATE" json:"exp_delivery_date"`
	SoSysNo                    float64   `gorm:"column:SO_SYS_NO" json:"so_sys_no"`
	SoDocNo                    string    `gorm:"column:SO_DOC_NO" json:"so_doc_no"`
	ChangeNo                   float64   `gorm:"column:CHANGE_NO" json:"change_no"`
	CreationUserId             string    `gorm:"column:CREATION_USER_ID" json:"creation_user_id"`
	CreationDatetime           time.Time `gorm:"column:CREATION_DATETIME" json:"creation_datetime"`
	ChangeUserId               string    `gorm:"column:CHANGE_USER_ID" json:"change_user_id"`
	ChangeDatetime             time.Time `gorm:"column:CHANGE_DATETIME" json:"change_datetime"`
	UniqueSpecPrice            float64   `gorm:"column:UNIQUE_SPEC_PRICE" json:"unique_spec_price"`
	UniqueSpecPriceBeforeVat   float64   `gorm:"column:UNIQUE_SPEC_PRICE_BEFORE_VAT" json:"unique_spec_price_before_vat"`
	OfftrPriceOriginal         float64   `gorm:"column:OFFTR_PRICE_ORIGINAL" json:"offtr_price_original"`
	DpOverpay                  float64   `gorm:"column:DP_OVERPAY" json:"dp_overpay"`
	UnitPriceSubsidi           float64   `gorm:"column:UNIT_PRICE_SUBSIDI" json:"unit_price_subsidi"`
	Remark                     string    `gorm:"column:REMARK" json:"remark"`
	SubsidiAddtional           float64   `gorm:"column:SUBSIDI_ADDTIONAL" json:"subsidi_addtional"`
	SpmProgress                string    `gorm:"column:SPM_PROGRESS" json:"spm_progress"`
	PendingReason              string    `gorm:"column:PENDING_REASON" json:"pending_reason"`
	InvDpSysNo                 float64   `gorm:"column:INV_DP_SYS_NO" json:"inv_dp_sys_no"`
	TotalAccsPurcPrice         float64   `gorm:"column:TOTAL_ACCS_PURC_PRICE" json:"total_accs_purc_price"`
	OptionCode                 string    `gorm:"column:OPTION_CODE" json:"option_code"`
	PaymentType                string    `gorm:"column:PAYMENT_TYPE" json:"payment_type"`
}

func (u1 *UtSpm1) TableName() string {
	return "utSPM1"
}

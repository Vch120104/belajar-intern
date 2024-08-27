package transactionsparepartpayloads

import "time"

type GetAllDBResponses struct {
	PurchaseOrderSystemNumber int `json:"purchase_order_system_number" parent_entity:"trx_item_purchase_order"`
	//WarehouseId int `json:"warehouse_id" parent_entity:"trx_work_order_detail"`
	PurchaseOrderDocumentNumber string     `json:"purchase_order_document_number" parent_entity:"trx_item_purchase_order"`
	PurchaseOrderDocumentDate   *time.Time `json:"purchase_order_document_date" parent_entity:"trx_item_purchase_order"`
	PurchaseOrderStatusId       int        `json:"purchase_order_status_id" parent_entity:"trx_item_purchase_order"`
	OrderTypeId                 int        `json:"order_type_id" parent_entity:"trx_item_purchase_order"`
	WarehouseId                 int        `json:"warehouse_id" parent_entity:"trx_item_purchase_order"`
	SupplierId                  int        `json:"supplier_id" parent_entity:"trx_item_purchase_order"`
	PurchaseRequestSystemNumber int        `json:"purchase_request_system_number" parent_entity:"trx_purchase_request"`
}

type GetAllPurchaseOrderResponses struct {
	PurchaseOrderSystemNumber int `json:"purchase_order_system_number" parent_entity:"trx_item_purchase_order"`
	//WarehouseId int `json:"warehouse_id" parent_entity:"trx_work_order_detail"`
	PurchaseOrderDocumentNumber string     `json:"purchase_order_document_number" parent_entity:"trx_item_purchase_order"`
	PurchaseOrderDocumentDate   *time.Time `json:"purchase_order_document_date" parent_entity:"trx_item_purchase_order"`
	PurchaseOrderStatus         string     `json:"purchase_order_status" parent_entity:"trx_item_purchase_order"`
	OrderType                   string     `json:"order_type" parent_entity:"trx_item_purchase_order"`
	WarehouseName               string     `json:"warehouse_name" parent_entity:"trx_item_purchase_order"`
	SupplierName                string     `json:"supplier_name" parent_entity:"trx_item_purchase_order"`
	PurchaseRequestDocNo        string     `json:"purchase_request_doc_no" parent_entity:"trx_purchase_request"`
}

//	type PurchaseOrderEntities struct {
//		CompanyId                           int        `gorm:"column:company_id;size:30;" json:"company_id"`
//		PurchaseOrderSystemNumber           int        `gorm:"column:purchase_order_system_number;size:30;not null;primaryKey;" json:"purchase_order_system_number"`
//		PurchaseOrderDocumentNumber         string     `gorm:"column:purchase_order_document_number;size:30;" json:"purchase_order_document_number"`
//		PurchaseOrderDocumentDate           *time.Time `gorm:"column:purchase_order_document_date;size:30;" json:"purchase_order_document_date"`
//		PurchaseOrderStatusId               int        `gorm:"column:purchase_order_status_id;size:30;" json:"purchase_order_status_id"`
//		BrandId                             int        `gorm:"column:brand_id;size:30;" json:"brand_id"`
//		ItemGroupId                         int        `gorm:"column:item_group_id;size:30;" json:"item_group_id"`
//		OrderTypeId                         int        `gorm:"column:order_type_id;size:30;" json:"order_type_id"`
//		SupplierId                          int        `gorm:"column:supplier_id;size:30;" json:"supplier_id"`
//		SupplierPicId                       int        `gorm:"column:supplier_pic_id;size:30;" json:"supplier_pic_id"`
//		WarehouseId                         int        `gorm:"column:warehouse_id;size:30;" json:"warehouse_id"`
//		CostCenterId                        int        `gorm:"column:cost_center_id;size:2;" json:"cost_center_id"`
//		ProfitType                          string     `gorm:"column:profit_type;size:30;" json:"profit_type"`
//		ProfitCenterId                      int        `gorm:"column:profit_center_id;size:30;" json:"profit_center_id"`
//		AffiliatedPurchaseOrder             bool       `gorm:"column:affiliated_purchase_order" json:"affiliated_purchase_order"`
//		CurrencyId                          int        `gorm:"column:currency_id;size:30;" json:"currency_id"`
//		BackOrder                           bool       `gorm:"column:back_order;" json:"back_order"`
//		SetOrder                            bool       `gorm:"set_order;" json:"set_order"`
//		ViaBinning                          bool       `gorm:"via_binning;" json:"via_binning"`
//		VatCode                             string     `gorm:"column:vat_code;size:30;" json:"vat_code"`
//		TotalDiscount                       *float64   `gorm:"column:total_discount;" json:"total_discount"`
//		TotalAmount                         *float64   `gorm:"column:total_amount;" json:"total_amount"`
//		TotalVat                            *float64   `gorm:"column:total_vat;" json:"total_vat"`
//		TotalAfterVat                       *float64   `gorm:"column:total_after_vat;" json:"total_after_vat"`
//		LastTotalDiscount                   *float64   `gorm:"column:last_total_discount;" json:"last_total_discount"`
//		LastTotalAmount                     *float64   `gorm:"column:last_total_amount;" json:"last_total_amount"`
//		LastTotalVat                        *float64   `gorm:"column:last_total_vat;" json:"last_total_vat"`
//		LastTotalAfterVat                   *float64   `gorm:"column:last_total_after_vat;" json:"last_total_after_vat"`
//		TotalAmountConfirm                  *float64   `gorm:"column:total_amount_confirm;" json:"total_amount_confirm"`
//		PurchaseOrderRemark                 string     `gorm:"column:purchase_order_remark;size:256;" json:"purchase_order_remark"`
//		DpRequest                           *float64   `gorm:"column:dp_request;" json:"dp_request"`
//		DpPayment                           *float64   `gorm:"column:dp_payment;" json:"dp_payment"`
//		DpPaymentAllocated                  *float64   `gorm:"column:dp_payment_allocated;" json:"dp_payment_allocated"`
//		DpPaymentAllocatedInvoice           *float64   `gorm:"column:dp_payment_allocated_invoice;" json:"dp_payment_allocated_invoice"`
//		DpPaymentAllocatedPpn               *float64   `gorm:"column:dp_payment_allocated_ppn;" json:"dp_payment_allocated_ppn"`
//		DpPaymentAllocatedRequestForPayment *float64   `gorm:"column:dp_payment_allocated_request_for_payment;" json:"dp_payment_allocated_request_for_payment"`
//		DeliveryId                          int        `gorm:"column:delivery_id;" json:"delivery_id"`
//		ExpectedDeliveryDate                *time.Time `gorm:"column:expected_delivery_date;" json:"expected_delivery_date"`
//		ExpectedArrivalDate                 *time.Time `gorm:"column:expected_arrival_date;" json:"expected_arrival_date"`
//		EstimatedDeliveryDate               *time.Time `gorm:"column:estimated_delivery_date;" json:"estimated_delivery_date"`
//		EstimatedDeliveryTime               string     `gorm:"column:estimated_delivery_time;size:5;" json:"estimated_delivery_time"`
//		SalesOrderSystemNumber              int        `gorm:"column:sales_order_system_number;" json:"sales_order_system_number"`
//		SalesOrderDocumentNumber            string     `gorm:"column:sales_order_document_number;size:25;" json:"sales_order_document_number"`
//		LastPrintById                       int        `gorm:"column:last_print_by_id;" json:"last_print_by_id"`
//		ApprovalRequestById                 int        `gorm:"column:approval_request_by_id;" json:"approval_request_by_id"`
//		ApprovalRequestNumber               int        `gorm:"column:approval_request_number;" json:"approval_request_number"`
//		ApprovalRequestDate                 *time.Time `gorm:"column:approval_request_date;" json:"approval_request_date"`
//		ApprovalRemark                      string     `gorm:"column:approval_remark;size:256;" json:"approval_remark"`
//		ApprovalLastById                    int        `gorm:"column:approval_last_by_id;" json:"approval_last_by_id"`
//		ApprovalLastDate                    *time.Time `gorm:"column:approval_last_date;" json:"approval_last_date"`
//		TotalInvoiceDownPayment             *float64   `gorm:"column:total_invoice_down_payment;" json:"total_invoice_down_payment"`
//		TotalInvoiceDownPaymentVat          *float64   `gorm:"column:total_invoice_down_payment_vat;" json:"total_invoice_down_payment_vat"`
//		TotalInvoiceDownPaymentAfterVat     *float64   `gorm:"column:total_invoice_down_payment_after_vat;" json:"total_invoice_down_payment_after_vat"`
//		DownPaymentReturn                   *float64   `gorm:"column:down_payment_return;" json:"down_payment_return"`
//		JournalSystemNumber                 int        `gorm:"column:journal_system_number;" json:"journal_system_number"`
//		EventNumber                         string     `gorm:"column:event_number;size:10;" json:"event_number"`
//		ItemClassId                         int        `gorm:"column:item_class_id;" json:"item_class_id"`
//		APMIsDirectShipment                    string     `gorm:"column:is_direct_shipment;size:1;" json:"is_direct_shipment"`
//		DirectShipmentCustomerId                          int        `gorm:"column:customer_id;" json:"customer_id"`
//		ExternalPurchaseOrderNumber         string     `gorm:"column:external_purchase_order_number;size:10;" json:"external_purchase_order_number"`
//		PurchaseOrderTypeId                 int        `gorm:"column:purchase_order_type_id;" json:"purchase_order_type_id"`
//		CurrencyExchangeRate                *float64   `gorm:"column:currency_exchange_rate;" json:"currency_exchange_rate"`
//		//PurchaseOrderDetail                 []PurchaseOrderDetailEntities `gorm:"foreignKey:PurchaseOrderSystemNumber;references:PurchaseOrderSystemNumber" json:"work_order_detail"`
//	}
type OrderTypeStatusResponse struct {
	OrderTypeId   int    `json:"order_type_id"`
	OrderTypeCode string `json:"order_type_code"`
	OrderTypeName string `json:"order_type_name"`
}
type WarehouseResponsesPurchaseOrder struct {
	IsActive                      bool   `json:"is_active"`
	WarehouseId                   int    `json:"warehouse_id"`
	WarehouseCostingType          string `json:"warehouse_costing_type"`
	WarehouseKaroseri             bool   `json:"warehouse_karoseri"`
	WahouseNegativeStock          bool   `json:"wahouse_negative_stock"`
	WarehouseReplishmentIndicator bool   `json:"warehouse_replishment_indicator"`
	WarehouseContact              string `json:"warehouse_contact"`
	WarehouseCode                 string `json:"warehouse_code"`
	AddressId                     int    `json:"address_id"`
	BrandId                       int    `json:"brand_id"`
	SupplierId                    int    `json:"supplier_id"`
	UserId                        int    `json:"user_id"`
	WarehouseSalesAllow           bool   `json:"warehouse_sales_allow"`
	WarehouseInTransit            bool   `json:"warehouse_in_transit"`
	WarehouseName                 string `json:"warehouse_name"`
	WarehouseDetailName           string `json:"warehouse_detail_name"`
	WarehouseTransitDefault       string `json:"warehouse_transit_default"`
	WarehouseGroupId              int    `json:"warehouse_group_id"`
	WarehousePhoneNumber          string `json:"warehouse_phone_number"`
	WarehouseFaxNumber            string `json:"warehouse_fax_number"`
	AddressDetails                struct {
		AddressId      int    `json:"address_id"`
		AddressStreet1 string `json:"address_street_1"`
		AddressStreet2 string `json:"address_street_2"`
		AddressStreet3 string `json:"address_street_3"`
		VillageId      int    `json:"village_id"`
	} `json:"address_details"`
	BrandDetails struct {
		BrandId   int    `json:"brand_id"`
		BrandCode string `json:"brand_code"`
		BrandName string `json:"brand_name"`
	} `json:"brand_details"`
	SupplierDetails struct {
		SupplierId   int    `json:"supplier_id"`
		SupplierCode string `json:"supplier_code"`
		SupplierName string `json:"supplier_name"`
	} `json:"supplier_details"`
	UserDetails struct {
		UserId        int    `json:"user_id"`
		EmployeeName  string `json:"employee_name"`
		JobPositionId int    `json:"job_position_id"`
	} `json:"user_details"`
	JobPositionDetails struct {
		JobPositionId   int    `json:"job_position_id"`
		JobPositionName string `json:"job_position_name"`
	} `json:"job_position_details"`
	VillageDetails struct {
		VillageId      int    `json:"village_id"`
		VillageName    string `json:"village_name"`
		DistrictCode   string `json:"district_code"`
		DistrictName   string `json:"district_name"`
		CityName       string `json:"city_name"`
		ProvinceName   string `json:"province_name"`
		CountryName    string `json:"country_name"`
		VillageZipCode string `json:"village_zip_code"`
		CityPhoneArea  string `json:"city_phone_area"`
	} `json:"village_details"`
}
type CompanyDetailResponses struct {
	IsActive      bool   `json:"is_active"`
	CompanyId     int    `json:"company_id"`
	CompanyCode   string `json:"company_code"`
	CompanyTypeId int    `json:"company_type_id"`
	CompanyName   string `json:"company_name"`
	VatCompany    struct {
		NpwpNo             string `json:"npwp_no"`
		NpwpDate           string `json:"npwp_date"`
		PkpType            bool   `json:"pkp_type"`
		PkpNo              string `json:"pkp_no"`
		PkpDate            string `json:"pkp_date"`
		TaxTransactionId   int    `json:"tax_transaction_id"`
		Name               string `json:"name"`
		AddressStreet1     string `json:"address_street_1"`
		AddressStreet2     string `json:"address_street_2"`
		AddressStreet3     string `json:"address_street_3"`
		VillageId          int    `json:"village_id"`
		TaxServiceOfficeId int    `json:"tax_service_office_id"`
	} `json:"vat_company"`
	TaxCompanyId int `json:"tax_company_id"`
	TaxCompany   struct {
		NpwpNo             string `json:"npwp_no"`
		NpwpDate           string `json:"npwp_date"`
		PkpType            bool   `json:"pkp_type"`
		PkpNo              string `json:"pkp_no"`
		PkpDate            string `json:"pkp_date"`
		TaxTransactionId   int    `json:"tax_transaction_id"`
		Name               string `json:"name"`
		AddressStreet1     string `json:"address_street_1"`
		AddressStreet2     string `json:"address_street_2"`
		AddressStreet3     string `json:"address_street_3"`
		VillageId          int    `json:"village_id"`
		TaxServiceOfficeId int    `json:"tax_service_office_id"`
	} `json:"tax_company"`
}
type SupplierResponsesAPI struct {
	IsActive              bool        `json:"is_active"`
	SupplierId            int         `json:"supplier_id"`
	CompanyId             int         `json:"company_id"`
	EffectiveDate         string      `json:"effective_date"`
	SupplierCode          string      `json:"supplier_code"`
	SupplierTitlePrefix   string      `json:"supplier_title_prefix"`
	SupplierName          string      `json:"supplier_name"`
	SupplierTitleSuffix   string      `json:"supplier_title_suffix"`
	ClientTypeId          int         `json:"client_type_id"`
	TermOfPaymentId       int         `json:"term_of_payment_id"`
	DefaultCurrencyId     int         `json:"default_currency_id"`
	ViaBinning            bool        `json:"via_binning"`
	IsImportSupplier      bool        `json:"is_import_supplier"`
	SupplierInvoiceTypeId int         `json:"supplier_invoice_type_id"`
	SupplierUniqueId      string      `json:"supplier_unique_id"`
	SupplierNitkuId       interface{} `json:"supplier_nitku_id"`
	SupplierAddressId     int         `json:"supplier_address_id"`
	SupplierAddress       struct {
		AddressStreet1 string `json:"address_street_1"`
		AddressStreet2 string `json:"address_street_2"`
		AddressStreet3 string `json:"address_street_3"`
		VillageId      int    `json:"village_id"`
	} `json:"supplier_address"`
	SupplierPhoneNo      string   `json:"supplier_phone_no"`
	SupplierFaxNo        string   `json:"supplier_fax_no"`
	SupplierMobilePhone  string   `json:"supplier_mobile_phone"`
	SupplierEmailAddress string   `json:"supplier_email_address"`
	MinimumDownPayment   *float64 `json:"minimum_down_payment"`
	BehaviourId          int      `json:"behaviour_id"`
	SupplierCategoryId   int      `json:"supplier_category_id"`
	//TaxIndustry          string `json:"tax_industry"`
	VatSupplierId int `json:"vat_supplier_id"`
	VatSupplier   struct {
		NpwpNo             string `json:"npwp_no"`
		NpwpDate           string `json:"npwp_date"`
		PkpType            bool   `json:"pkp_type"`
		PkpNo              string `json:"pkp_no"`
		PkpDate            string `json:"pkp_date"`
		TaxTransactionId   int    `json:"tax_transaction_id"`
		Name               string `json:"name"`
		AddressStreet1     string `json:"address_street_1"`
		AddressStreet2     string `json:"address_street_2"`
		AddressStreet3     string `json:"address_street_3"`
		VillageId          int    `json:"village_id"`
		TaxServiceOfficeId int    `json:"tax_service_office_id"`
	} `json:"vat_supplier"`
	TaxSupplierId int `json:"tax_supplier_id"`
	TaxSupplier   struct {
		NpwpNo             string `json:"npwp_no"`
		NpwpDate           string `json:"npwp_date"`
		PkpType            bool   `json:"pkp_type"`
		PkpNo              string `json:"pkp_no"`
		PkpDate            string `json:"pkp_date"`
		TaxTransactionId   int    `json:"tax_transaction_id"`
		Name               string `json:"name"`
		AddressStreet1     string `json:"address_street_1"`
		AddressStreet2     string `json:"address_street_2"`
		AddressStreet3     string `json:"address_street_3"`
		VillageId          int    `json:"village_id"`
		TaxServiceOfficeId int    `json:"tax_service_office_id"`
	} `json:"tax_supplier"`
	SupplierContact struct {
		Page      int `json:"page"`
		PageLimit int `json:"page_limit"`
		Npages    int `json:"npages"`
		Nrows     int `json:"nrows"`
		Data      []struct {
			ClientContactId int         `json:"client_contact_id"`
			ContactName     string      `json:"contact_name"`
			DivisionName    string      `json:"division_name"`
			JobTitleName    string      `json:"job_title_name"`
			PhoneNumber     string      `json:"phone_number"`
			GenderId        int         `json:"gender_id"`
			EmailAddress    interface{} `json:"email_address"`
			IsActive        bool        `json:"is_active"`
		} `json:"data"`
	} `json:"supplier_contact"`
	SupplierBankAccount struct {
		Page      int `json:"page"`
		PageLimit int `json:"page_limit"`
		Npages    int `json:"npages"`
		Nrows     int `json:"nrows"`
		Data      []struct {
			BankAccountId     int    `json:"bank_account_id"`
			IsActive          bool   `json:"is_active"`
			BankId            int    `json:"bank_id"`
			BankAccountTypeId int    `json:"bank_account_type_id"`
			CurrencyId        int    `json:"currency_id"`
			BankAccountNumber string `json:"bank_account_number"`
			BankAccountName   string `json:"bank_account_name"`
		} `json:"data"`
	} `json:"supplier_bank_account"`
}
type TaxRateResponseApi struct {
	TaxPercent *float64 `json:"tax_percent"`
}

type PurchaseRequestResponse struct {
	PurchaseRequestSystemNumber   int    `json:"purchase_request_system_number"`
	PurchaseRequestDocumentNumber string `json:"purchase_request_document_number"`
}

type PurchaseOrderGetByIdResponses struct {
	//CompanyId                   int        `gorm:"column:company_id;size:30;" json:"company_id"`
	//PurchaseOrderSystemNumber   int        `gorm:"column:purchase_order_system_number;size:30;not null;primaryKey;" json:"purchase_order_system_number"`
	//PurchaseOrderDocumentNumber string     `gorm:"column:purchase_order_document_number;size:30;" json:"purchase_order_document_number"`
	//PurchaseOrderDocumentDate   *time.Time `gorm:"column:purchase_order_document_date;size:30;" json:"purchase_order_document_date"`
	//PurchaseOrderStatusId       int        `gorm:"column:purchase_order_status_id;size:30;" json:"purchase_order_status_id"`
	//BrandId                     int        `gorm:"column:brand_id;size:30;" json:"brand_id"`
	//ItemGroupId                 int        `gorm:"column:item_group_id;size:30;" json:"item_group_id"`
	//OrderTypeId                 int        `gorm:"column:order_type_id;size:30;" json:"order_type_id"`
	//SupplierId                  int        `gorm:"column:supplier_id;size:30;" json:"supplier_id"`
	//SupplierPicId               int        `gorm:"column:supplier_pic_id;size:30;" json:"supplier_pic_id"`
	//WarehouseId                 int        `gorm:"column:warehouse_id;size:30;" json:"warehouse_id"`
	//WarehouseGroupId            int        `gorm:"column:warehouse_group_id;size:30;" json:"warehouse_group_id"`
	//CostCenterId                int        `gorm:"column:cost_center_id;size:2;" json:"cost_center_id"`
	//ProfitType                  string     `gorm:"column:profit_type;size:30;" json:"profit_type"`
	//ProfitCenterId              int        `gorm:"column:profit_center_id;size:30;" json:"profit_center_id"`
	//AffiliatedPurchaseOrder     bool       `gorm:"column:affiliated_purchase_order" json:"affiliated_purchase_order"`
	//CurrencyId                  int        `gorm:"column:currency_id;size:30;" json:"currency_id"`
	//BackOrder                   bool       `gorm:"column:back_order;" json:"back_order"`
	//SetOrder                    bool       `gorm:"set_order;" json:"set_order"`
	//ViaBinning                  bool       `gorm:"via_binning;" json:"via_binning"`
	//VatCode                     string     `gorm:"column:vat_code;size:30;" json:"vat_code"`
	//PphCode                     string     `gorm:"column:pph_code;size:30;" json:"pph_code"`
	//TotalDiscount               *float64   `gorm:"column:total_discount;" json:"total_discount"`
	//TotalAmount                 *float64   `gorm:"column:total_amount;" json:"total_amount"`
	//TotalVat                    *float64   `gorm:"column:total_vat;" json:"total_vat"`
	//TotalAfterVat               *float64   `gorm:"column:total_after_vat;" json:"total_after_vat"`
	//LastTotalDiscount           *float64   `gorm:"column:last_total_discount;" json:"last_total_discount"`
	//LastTotalAmount             *float64   `gorm:"column:last_total_amount;" json:"last_total_amount"`
	//LastTotalVat                *float64   `gorm:"column:last_total_vat;" json:"last_total_vat"`
	//LastTotalAfterVat           *float64   `gorm:"column:last_total_after_vat;" json:"last_total_after_vat"`

	//CompanyId                           int        `json:"company_id" parent_entity:"trx_item_purchase_order"`
	PurchaseOrderSystemNumber   int        `json:"purchase_order_system_number" parent_entity:"trx_item_purchase_order" gorm:"not null;primaryKey;"`
	PurchaseOrderDocumentNumber string     `json:"purchase_order_document_number" parent_entity:"trx_item_purchase_order"`
	PurchaseOrderDocumentDate   *time.Time `json:"purchase_order_document_date" parent_entity:"trx_item_purchase_order"`
	ExternalPurchaseOrderNumber string     `json:"external_purchase_order_number" parent_entity:"trx_item_purchase_order"`
	PurchaseOrderStatusId       int        `json:"purchase_order_status_id" parent_entity:"trx_item_purchase_order"`
	BrandId                     int        `json:"brand_id" parent_entity:"trx_item_purchase_order"`
	ItemGroupId                 int        `json:"item_group_id" parent_entity:"trx_item_purchase_order"`
	SupplierId                  int        `json:"supplier_id" parent_entity:"trx_item_purchase_order"`
	SupplierPicId               int        `json:"supplier_pic_id" parent_entity:"trx_item_purchase_order"`
	WarehouseId                 int        `json:"warehouse_id" parent_entity:"trx_item_purchase_order"`
	WarehouseGroupId            int        `json:"warehouse_group_id" parent_entity:"trx_item_purchase_order"`
	CostCenterId                int        `json:"cost_center_id" parent_entity:"trx_item_purchase_order"`
	ProfitCenterId              int        `json:"profit_center_id" parent_entity:"trx_item_purchase_order"`
	AffiliatedPurchaseOrder     bool       `json:"affiliated_purchase_order" parent_entity:"trx_item_purchase_order"`
	CurrencyId                  int        `json:"currency_id" parent_entity:"trx_item_purchase_order"`
	BackOrder                   bool       `json:"back_order" parent_entity:"trx_item_purchase_order"`
	SetOrder                    bool       `json:"set_order" parent_entity:"trx_item_purchase_order"`
	ViaBinning                  bool       `json:"via_binning" parent_entity:"trx_item_purchase_order"`
	PurchaseOrderRemark         string     `json:"purchase_order_remark" parent_entity:"trx_item_purchase_order"`
	DpRequest                   *float64   `json:"dp_request" parent_entity:"trx_item_purchase_order"`
	DeliveryId                  int        `json:"delivery_id" parent_entity:"trx_item_purchase_order"`
	ExpectedDeliveryDate        *time.Time `json:"expected_delivery_date" parent_entity:"trx_item_purchase_order"`
	ExpectedArrivalDate         *time.Time `json:"expected_arrival_date" parent_entity:"trx_item_purchase_order"`
	PurchaseOrderTypeId         int        `json:"purchase_order_type_id" parent_entity:"trx_item_purchase_order"`
	CreatedByUserId             int        `json:"created_by_user_id" parent_entity:"trx_item_purchase_order"`
	CreatedDate                 *time.Time `json:"created_date" parent_entity:"trx_item_purchase_order"`
	UpdatedByUserId             int        `json:"updated_by_user_id" parent_entity:"trx_item_purchase_order"`
	UpdatedDate                 *time.Time `json:"updated_date" parent_entity:"trx_item_purchase_order"`
	ChangeNo                    int        `json:"change_no" parent_entity:"trx_item_purchase_order"`
	CustomerId                  int        `json:"customer_id" parent_entity:"customer_id"`
	TotalDiscount               *float64   `json:"total_discount" parent_entity:"customer_id"`
	TotalAmount                 *float64   `json:"total_amount" parent_entity:"customer_id"`
	TotalVat                    *float64   `json:"total_vat" parent_entity:"customer_id"`
	TotalAfterVat               *float64   `json:"total_after_vat" parent_entity:"customer_id"`
	APMIsDirectShipment         string     `json:"apm_is_direct_shipment" parent_entity:"trx_item_purchase_order"`
}
type PurchaseOrderGetDetail struct {
	PurchaseOrderDetailSystemNumber int      `gorm:"column:purchase_order_detail_system_number;"  parent_entity:"trx_item_purchase_order_detail"`
	Snp                             *float64 `gorm:"column:snp;"  parent_entity:"trx_item_purchase_order_detail"`
	ItemDiscountAmount              *float64 `gorm:"column:item_discount_amount;" json:"item_discount_amount"`
	ItemPrice                       *float64 `gorm:"column:item_price;" json:"item_price"`
	ItemQuantity                    *float64 `gorm:"column:item_quantity;" json:"item_quantity"`
	ItemUnitOfMeasurement           string   `gorm:"column:item_unit_of_measurement;size:1;" json:"item_unit_of_measurement"`
	UnitOfMeasurementRate           *float64 `gorm:"column:unit_of_measurement_rate;" json:"unit_of_measurement_rate"`
	ItemCode                        string   `gorm:"column:item_code;" json:"item_code"`
	ItemName                        string   `gorm:"column:item_name;" json:"item_name"`
	//
	PurchaseOrderSystemNumber     int      `gorm:"column:purchase_order_system_number;size:30;" json:"purchase_order_system_number"`
	PurchaseOrderLineNumber       int      `gorm:"column:purchase_order_line_number;" json:"purchase_order_line_number"`
	ItemTotal                     *float64 `gorm:"column:item_total;" json:"item_total"`
	PurchaseRequestSystemNumber   int      `gorm:"column:purchase_order_system_number;size:30;" json:"purchase_request_system_number"`
	PurchaseRequestLineNumber     int      `gorm:"column:purchase_request_line_number;" json:"purchase_request_line_number"`
	PurchaseRequestDocumentNumber string   `gorm:"column:purchase_request_document_number;" json:"purchase_request_document_number"`
	StockOnHand                   *float64 `gorm:"column:stock_on_hand;" json:"stock_on_hand"`
	ItemRemark                    string   `gorm:"column:item_remark;size:255;" json:"item_remark"`
}
type PurchaseOrderNewPurchaseOrderPayloads struct {
	CompanyId                   int        `json:"company_id" parent_entity:"trx_item_purchase_order"`
	PurchaseOrderSystemNumber   int        `json:"purchase_order_system_number" parent_entity:"trx_item_purchase_order" gorm:"not null;primaryKey;"`
	PurchaseOrderDocumentNumber string     `json:"purchase_order_document_number" parent_entity:"trx_item_purchase_order"`
	PurchaseOrderDocumentDate   *time.Time `json:"purchase_order_document_date" parent_entity:"trx_item_purchase_order"`
	ExternalPurchaseOrderNumber string     `json:"external_purchase_order_number" parent_entity:"trx_item_purchase_order"`
	PurchaseOrderStatusId       int        `json:"purchase_order_status_id" parent_entity:"trx_item_purchase_order"`
	BrandId                     int        `json:"brand_id" parent_entity:"trx_item_purchase_order"`
	ItemGroupId                 int        `json:"item_group_id" parent_entity:"trx_item_purchase_order"`
	SupplierId                  int        `json:"supplier_id" parent_entity:"trx_item_purchase_order"`
	SupplierPicId               int        `json:"supplier_pic_id" parent_entity:"trx_item_purchase_order"`
	WarehouseId                 int        `json:"warehouse_id" parent_entity:"trx_item_purchase_order"`
	WarehouseGroupId            int        `json:"warehouse_group_id" parent_entity:"trx_item_purchase_order"`
	CostCenterId                int        `json:"cost_center_id" parent_entity:"trx_item_purchase_order"`
	ProfitCenterId              int        `json:"profit_center_id" parent_entity:"trx_item_purchase_order"`
	AffiliatedPurchaseOrder     bool       `json:"affiliated_purchase_order" parent_entity:"trx_item_purchase_order"`
	CurrencyId                  int        `json:"currency_id" parent_entity:"trx_item_purchase_order"`
	BackOrder                   bool       `json:"back_order" parent_entity:"trx_item_purchase_order"`
	SetOrder                    bool       `json:"set_order" parent_entity:"trx_item_purchase_order"`
	ViaBinning                  bool       `json:"via_binning" parent_entity:"trx_item_purchase_order"`
	PurchaseOrderRemark         string     `json:"purchase_order_remark" parent_entity:"trx_item_purchase_order"`
	DpRequest                   *float64   `json:"dp_request" parent_entity:"trx_item_purchase_order"`
	DeliveryId                  int        `json:"delivery_id" parent_entity:"trx_item_purchase_order"`
	ExpectedDeliveryDate        *time.Time `json:"expected_delivery_date" parent_entity:"trx_item_purchase_order"`
	ExpectedArrivalDate         *time.Time `json:"expected_arrival_date" parent_entity:"trx_item_purchase_order"`
	PurchaseOrderTypeId         int        `json:"purchase_order_type_id" parent_entity:"trx_item_purchase_order"`
	CreatedByUserId             int        `json:"created_by_user_id" parent_entity:"trx_item_purchase_order"`
	CreatedDate                 *time.Time `json:"created_date" parent_entity:"trx_item_purchase_order"`
	UpdatedByUserId             int        `json:"updated_by_user_id" parent_entity:"trx_item_purchase_order"`
	UpdatedDate                 *time.Time `json:"updated_date" parent_entity:"trx_item_purchase_order"`
	ChangeNo                    int        `json:"change_no" parent_entity:"trx_item_purchase_order"`
	CustomerId                  int        `json:"customer_id" parent_entity:"customer_id"`
	TotalDiscount               *float64   `json:"total_discount" parent_entity:"customer_id"`
	TotalAmount                 *float64   `json:"total_amount" parent_entity:"customer_id"`
	TotalVat                    *float64   `json:"total_vat" parent_entity:"customer_id"`
	TotalAfterVat               *float64   `json:"total_after_vat" parent_entity:"customer_id"`
	APMIsDirectShipment         string     `json:"apm_is_direct_shipment" parent_entity:"trx_item_purchase_order"`
}
type PurchaseOrderNewPurchaseOrderResponses struct {
	CompanyId                   int        `json:"company_id" parent_entity:"trx_item_purchase_order"`
	PurchaseOrderSystemNumber   int        `json:"purchase_order_system_number" parent_entity:"trx_item_purchase_order" gorm:"not null;primaryKey;"`
	PurchaseOrderDocumentNumber string     `json:"purchase_order_document_number" parent_entity:"trx_item_purchase_order"`
	PurchaseOrderDocumentDate   *time.Time `json:"purchase_order_document_date" parent_entity:"trx_item_purchase_order"`
	ExternalPurchaseOrderNumber string     `json:"external_purchase_order_number" parent_entity:"trx_item_purchase_order"`
	PurchaseOrderStatusId       int        `json:"purchase_order_status_id" parent_entity:"trx_item_purchase_order"`
	BrandId                     int        `json:"brand_id" parent_entity:"trx_item_purchase_order"`
	ItemGroupId                 int        `json:"item_group_id" parent_entity:"trx_item_purchase_order"`
	SupplierId                  int        `json:"supplier_id" parent_entity:"trx_item_purchase_order"`
	SupplierPicId               int        `json:"supplier_pic_id" parent_entity:"trx_item_purchase_order"`
	WarehouseId                 int        `json:"warehouse_id" parent_entity:"trx_item_purchase_order"`
	WarehouseGroupId            int        `json:"warehouse_group_id" parent_entity:"trx_item_purchase_order"`
	CostCenterId                int        `json:"cost_center_id" parent_entity:"trx_item_purchase_order"`
	ProfitCenterId              int        `json:"profit_center_id" parent_entity:"trx_item_purchase_order"`
	AffiliatedPurchaseOrder     bool       `json:"affiliated_purchase_order" parent_entity:"trx_item_purchase_order"`
	CurrencyId                  int        `json:"currency_id" parent_entity:"trx_item_purchase_order"`
	BackOrder                   bool       `json:"back_order" parent_entity:"trx_item_purchase_order"`
	SetOrder                    bool       `json:"set_order" parent_entity:"trx_item_purchase_order"`
	ViaBinning                  bool       `json:"via_binning" parent_entity:"trx_item_purchase_order"`
	PurchaseOrderRemark         string     `json:"purchase_order_remark" parent_entity:"trx_item_purchase_order"`
	DpRequest                   *float64   `json:"dp_request" parent_entity:"trx_item_purchase_order"`
	DeliveryId                  int        `json:"delivery_id" parent_entity:"trx_item_purchase_order"`
	ExpectedDeliveryDate        *time.Time `json:"expected_delivery_date" parent_entity:"trx_item_purchase_order"`
	ExpectedArrivalDate         *time.Time `json:"expected_arrival_date" parent_entity:"trx_item_purchase_order"`
	PurchaseOrderTypeId         int        `json:"purchase_order_type_id" parent_entity:"trx_item_purchase_order"`
	CreatedByUserId             int        `json:"created_by_user_id" parent_entity:"trx_item_purchase_order"`
	CreatedDate                 *time.Time `json:"created_date" parent_entity:"trx_item_purchase_order"`
	UpdatedByUserId             int        `json:"updated_by_user_id" parent_entity:"trx_item_purchase_order"`
	UpdatedDate                 *time.Time `json:"updated_date" parent_entity:"trx_item_purchase_order"`
	ChangeNo                    int        `json:"change_no" parent_entity:"trx_item_purchase_order"`
	CustomerId                  int        `json:"customer_id" parent_entity:"customer_id"`
	TotalDiscount               *float64   `json:"total_discount" parent_entity:"customer_id"`
	TotalAmount                 *float64   `json:"total_amount" parent_entity:"customer_id"`
	TotalVat                    *float64   `json:"total_vat" parent_entity:"customer_id"`
	TotalAfterVat               *float64   `json:"total_after_vat" parent_entity:"customer_id"`
	APMIsDirectShipment         string     `json:"apm_is_direct_shipment" parent_entity:"trx_item_purchase_order"`
}

type PurchaseOrderDetailPayloads struct {
	PurchaseOrderSystemNumber         int        `json:"purchase_order_system_number" parent_entity:"trx_item_purchase_order" gorm:"not null;primaryKey;"`
	PurchaseOrderLine                 int        `json:"purchase_order_line"`
	ItemId                            int        `json:"item_id"`
	ItemUnitOfMeasurement             string     `json:"item_unit_of_measurement"`
	UnitOfMeasurementRate             *float64   `json:"unit_of_measurement_rate"`
	ItemQuantity                      *float64   `json:"item_quantity"`
	ItemPrice                         *float64   `json:"item_price"`
	ItemTotal                         *float64   `json:"item_total"`
	PurchaseRequestDetailSystemNumber int        `json:"purchase_request_detail_system_number"`
	PurchaseRequestSystemNumber       int        `json:"purchase_request_system_number"`
	PurchaseRequestLineNumber         int        `json:"purchase_request_line_number"`
	StockOnHand                       *float64   `json:"stock_on_hand"`
	OldPurchaseOrderSystemNo          int        `json:"old_purchase_order_system_no"`
	ItemRemark                        string     `json:"item_remark"`
	OldPurchaseOrderLineNumber        int        `json:"old_purchase_order_line_number"`
	CreatedByUserId                   int        `json:"created_by_user_id"`
	CreatedDate                       *time.Time `json:"created_date"`
	UpdatedByUserId                   int        `json:"updated_by_user_id"`
	UpdatedDate                       *time.Time `json:"updated_date"`
	Snp                               *float64   `json:"snp"`
	ItemDiscountPercentage            *float64   `gorm:"column:item_discount_percentage;" json:"item_discount_percentage"`
	ItemDiscountAmount                *float64   `gorm:"column:item_discount_amount;" json:"item_discount_amount"`
	//ItemTotal                         *float64   `gorm:"column:item_total;" json:"item_total"`

	//totaldiscountnyadjabb nggak kagak kagak kaga kaga kagak bmsananananan
	//cek coba
}
type PurchaseOrderSaveDetailPayloads struct {
	//PurchaseOrderDetailSystemNumber
	PurchaseOrderDetailSystemNumber   int        `json:"purchase_order_detail_system_number"`
	PurchaseOrderSystemNumber         int        `json:"purchase_order_system_number" parent_entity:"trx_item_purchase_order" gorm:"not null;primaryKey;"`
	PurchaseOrderLine                 int        `json:"purchase_order_line"`
	ItemId                            int        `json:"item_id"`
	ItemUnitOfMeasurement             string     `json:"item_unit_of_measurement"`
	UnitOfMeasurementRate             *float64   `json:"unit_of_measurement_rate"`
	ItemQuantity                      *float64   `json:"item_quantity"`
	ItemPrice                         *float64   `json:"item_price"`
	ItemTotal                         *float64   `json:"item_total"`
	PurchaseRequestDetailSystemNumber int        `json:"purchase_request_detail_system_number"`
	PurchaseRequestSystemNumber       int        `json:"purchase_request_system_number"`
	PurchaseRequestLineNumber         int        `json:"purchase_request_line_number"`
	StockOnHand                       *float64   `json:"stock_on_hand"`
	OldPurchaseOrderSystemNo          int        `json:"old_purchase_order_system_no"`
	ItemRemark                        string     `json:"item_remark"`
	OldPurchaseOrderLineNumber        int        `json:"old_purchase_order_line_number"`
	CreatedByUserId                   int        `json:"created_by_user_id"`
	CreatedDate                       *time.Time `json:"created_date"`
	UpdatedByUserId                   int        `json:"updated_by_user_id"`
	UpdatedDate                       *time.Time `json:"updated_date"`
	Snp                               *float64   `json:"snp"`
	ItemDiscountPercentage            *float64   `gorm:"column:item_discount_percentage;" json:"item_discount_percentage"`
	ItemDiscountAmount                *float64   `gorm:"column:item_discount_amount;" json:"item_discount_amount"`
	//ItemTotal                         *float64   `gorm:"column:item_total;" json:"item_total"`

	//totaldiscountnya cek coba
}

//CompanyId                           int        `json:"company_id" parent_entity:"trx_item_purchase_order"`
//PurchaseOrderSystemNumber   int        `json:"purchase_order_system_number" parent_entity:"trx_item_purchase_order" gorm:"not null;primaryKey;"`
//PurchaseOrderDocumentNumber string     `json:"purchase_order_document_number" parent_entity:"trx_item_purchase_order"`
//PurchaseOrderDocumentDate   *time.Time `json:"purchase_order_document_date" parent_entity:"trx_item_purchase_order"`
////pr cari external po dari mana
//PurchaseOrderStatusId int `json:"purchase_order_status_id" parent_entity:"trx_item_purchase_order"`
//BrandId               int `json:"brand_id" parent_entity:"trx_item_purchase_order"`
//ItemGroupId           int `json:"item_group_id" parent_entity:"trx_item_purchase_order"`
////OrderTypeId                         int        `json:"order_type_id" parent_entity:"trx_item_purchase_order"`
//SupplierId int `json:"supplier_id" parent_entity:"trx_item_purchase_order"`
////SupplierPicId                       int        `json:"supplier_pic_id" parent_entity:"trx_item_purchase_order"`
//WarehouseId      int `json:"warehouse_id" parent_entity:"trx_item_purchase_order"`
//WarehouseGroupId int `json:"warehouse_group_id" parent_entity:"trx_item_purchase_order"`
//CostCenterId     int `json:"cost_center_id" parent_entity:"trx_item_purchase_order"`
////ProfitType                          string     `json:"profit_type" parent_entity:"trx_item_purchase_order"`
//ProfitCenterId          int  `json:"profit_center_id" parent_entity:"trx_item_purchase_order"`
//AffiliatedPurchaseOrder bool `json:"affiliated_purchase_order" parent_entity:"trx_item_purchase_order"`
//CurrencyId              int  `json:"currency_id" parent_entity:"trx_item_purchase_order"`
//BackOrder               bool `json:"back_order" parent_entity:"trx_item_purchase_order"`
//SetOrder                bool `json:"set_order" parent_entity:"trx_item_purchase_order"`
//ViaBinning              bool `json:"via_binning" parent_entity:"trx_item_purchase_order"`
////VatCode                             string     `json:"vat_code" parent_entity:"trx_item_purchase_order"`
////TotalDiscount                       *float64   `json:"total_discount" parent_entity:"trx_item_purchase_order"`
////TotalAmount                         *float64   `json:"total_amount" parent_entity:"trx_item_purchase_order"`
////TotalVat                            *float64   `json:"total_vat" parent_entity:"trx_item_purchase_order"`
////TotalAfterVat                       *float64   `json:"total_after_vat" parent_entity:"trx_item_purchase_order"`
////LastTotalDiscount                   *float64   `json:"last_total_discount" parent_entity:"trx_item_purchase_order"`
////LastTotalAmount                     *float64   `json:"last_total_amount" parent_entity:"trx_item_purchase_order"`
////LastTotalVat                        *float64   `json:"last_total_vat" parent_entity:"trx_item_purchase_order"`
////LastTotalAfterVat                   *float64   `json:"last_total_after_vat" parent_entity:"trx_item_purchase_order"`
////TotalAmountConfirm                  *float64   `json:"total_amount_confirm" parent_entity:"trx_item_purchase_order"`
//PurchaseOrderRemark string   `json:"purchase_order_remark" parent_entity:"trx_item_purchase_order"`
//DpRequest           *float64 `json:"dp_request" parent_entity:"trx_item_purchase_order"`
////DpPayment                           *float64   `json:"dp_payment" parent_entity:"trx_item_purchase_order"`
////DpPaymentAllocated                  *float64   `json:"dp_payment_allocated" parent_entity:"trx_item_purchase_order"`
////DpPaymentAllocatedInvoice           *float64   `json:"dp_payment_allocated_invoice" parent_entity:"trx_item_purchase_order"`
////DpPaymentAllocatedPpn               *float64   `json:"dp_payment_allocated_ppn" parent_entity:"trx_item_purchase_order"`
////DpPaymentAllocatedRequestForPayment *float64   `json:"dp_payment_allocated_request_for_payment" parent_entity:"trx_item_purchase_order"`
//DeliveryId           int        `json:"delivery_id" parent_entity:"trx_item_purchase_order"`
//ExpectedDeliveryDate *time.Time `json:"expected_delivery_date" parent_entity:"trx_item_purchase_order"`
//ExpectedArrivalDate  *time.Time `json:"expected_arrival_date" parent_entity:"trx_item_purchase_order"`
////EstimatedDeliveryDate               *time.Time `json:"estimated_delivery_date" parent_entity:"trx_item_purchase_order"`
////EstimatedDeliveryTime               string     `json:"estimated_delivery_time" parent_entity:"trx_item_purchase_order"`
////SalesOrderSystemNumber              int        `json:"sales_order_system_number" parent_entity:"trx_item_purchase_order"`
////SalesOrderDocumentNumber            string     `json:"sales_order_document_number" parent_entity:"trx_item_purchase_order"`
////LastPrintById                       int        `json:"last_print_by_id" parent_entity:"trx_item_purchase_order"`
////ApprovalRequestById                 int        `json:"approval_request_by_id" parent_entity:"trx_item_purchase_order"`
////ApprovalRequestNumber               int        `json:"approval_request_number" parent_entity:"trx_item_purchase_order"`
////ApprovalRequestDate                 *time.Time `json:"approval_request_date" parent_entity:"trx_item_purchase_order"`
////ApprovalRemark                      string     `json:"approval_remark" parent_entity:"trx_item_purchase_order"`
////ApprovalLastById                    int        `json:"approval_last_by_id" parent_entity:"trx_item_purchase_order"`
////ApprovalLastDate                    *time.Time `json:"approval_last_date" parent_entity:"trx_item_purchase_order"`
////TotalInvoiceDownPayment             *float64   `json:"total_invoice_down_payment" parent_entity:"trx_item_purchase_order"`
////TotalInvoiceDownPaymentVat          *float64   `json:"total_invoice_down_payment_vat" parent_entity:"trx_item_purchase_order"`
////TotalInvoiceDownPaymentAfterVat     *float64   `json:"total_invoice_down_payment_after_vat" parent_entity:"trx_item_purchase_order"`
////DownPaymentReturn                   *float64   `json:"down_payment_return" parent_entity:"trx_item_purchase_order"`
////JournalSystemNumber                 int        `json:"journal_system_number" parent_entity:"trx_item_purchase_order"`
////EventNumber                         string     `json:"event_number" parent_entity:"trx_item_purchase_order"`
////ItemClassId                         int        `json:"item_class_id" parent_entity:"trx_item_purchase_order"`
////APMIsDirectShipment                    string     `json:"is_direct_shipment" parent_entity:"trx_item_purchase_order"`
////DirectShipmentCustomerId                          int        `json:"customer_id" parent_entity:"trx_item_purchase_order"`
////ExternalPurchaseOrderNumber         string     `json:"external_purchase_order_number" parent_entity:"trx_item_purchase_order"`
//PurchaseOrderTypeId int `json:"purchase_order_type_id" parent_entity:"trx_item_purchase_order"`
////CurrencyExchangeRate                *float64   `json:"currency_exchange_rate" parent_entity:"trx_item_purchase_order"`
//CreatedByUserId int        `json:"created_by_user_id" parent_entity:"trx_item_purchase_order"`
//CreatedDate     *time.Time `json:"created_date" parent_entity:"trx_item_purchase_order"`
//UpdatedByUserId int        `json:"updated_by_user_id" parent_entity:"trx_item_purchase_order"`
//UpdatedDate     *time.Time `json:"updated_date" parent_entity:"trx_item_purchase_order"`
//ChangeNo        int        `json:"change_no" parent_entity:"trx_item_purchase_order"`

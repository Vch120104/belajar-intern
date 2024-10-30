package generalservicepayloads

type SupplierAddress struct {
	AddressStreet1 string `json:"address_street_1"`
	AddressStreet2 string `json:"address_street_2"`
	AddressStreet3 string `json:"address_street_3"`
	VillageId      int    `json:"village_id"`
}
type VatSupplier struct {
	NpwpNo             string `json:"npwp_no"`
	NpwpDate           string `json:"npwp_date"`
	PkpType            string `json:"pkp_type"`
	PkpNo              string `json:"pkp_no"`
	PkpDate            string `json:"pkp_date"`
	TaxTransactionId   int    `json:"tax_transaction_id"`
	TaxBranchCode      string `json:"tax_branch_code"`
	Reserve            string `json:"reserve"`
	Name               string `json:"name"`
	AddressStreet1     string `json:"address_street_1"`
	AddressStreet2     string `json:"address_street_2"`
	AddressStreet3     string `json:"address_street_3"`
	VillageId          int    `json:"village_id"`
	TaxServiceOfficeId int    `json:"tax_service_office_id"`
}
type TaxSupplier struct {
	NpwpNo             string `json:"npwp_no"`
	NpwpDate           string `json:"npwp_date"`
	PkpType            string `json:"pkp_type"`
	PkpNo              string `json:"pkp_no"`
	PkpDate            string `json:"pkp_date"`
	TaxTransactionId   int    `json:"tax_transaction_id"`
	TaxBranchCode      string `json:"tax_branch_code"`
	Reserve            string `json:"reserve"`
	Name               string `json:"name"`
	AddressStreet1     string `json:"address_street_1"`
	AddressStreet2     string `json:"address_street_2"`
	AddressStreet3     string `json:"address_street_3"`
	VillageId          int    `json:"village_id"`
	TaxServiceOfficeId int    `json:"tax_service_office_id"`
}
type SupplierContactData struct {
	ClientContactId int    `json:"client_contact_id"`
	ContactName     string `json:"contact_name"`
	DivisionName    string `json:"division_name"`
	JobTitleName    string `json:"job_title_name"`
	PhoneNumber     string `json:"phone_number"`
	GenderId        int    `json:"gender_id"`
	EmailAddress    string `json:"email_address"`
	IsActive        bool   `json:"is_active"`
}
type SupplierContact struct {
	Page      int                   `json:"page"`
	PageLimit int                   `json:"page_limit"`
	Npages    int                   `json:"npages"`
	Nrows     int                   `json:"nrows"`
	Data      []SupplierContactData `json:"data"`
}
type SupplierBankAccountData struct {
	BankAccountId              int    `json:"bank_account_id"`
	IsActive                   bool   `json:"is_active"`
	BankId                     int    `json:"bank_id"`
	BankName                   string `json:"bank_name"`
	BankAccountTypeDescription string `json:"bank_account_type_description"`
	CurrencyCode               string `json:"currency_code"`
	BankAccountNumber          string `json:"bank_account_number"`
	BankAccountName            string `json:"bank_account_name"`
}
type SupplerBankAccount struct {
	Page      int                       `json:"page"`
	PageLimit int                       `json:"page_limit"`
	Npages    int                       `json:"npages"`
	Nrows     int                       `json:"nrows"`
	Data      []SupplierBankAccountData `json:"data"`
}
type SupplierMasterCrossServicePayloads struct {
	IsActive              bool               `json:"is_active"`
	SupplierId            int                `json:"supplier_id"`
	CompanyId             int                `json:"company_id"`
	EffectiveDate         string             `json:"effective_date"`
	SupplierCode          string             `json:"supplier_code"`
	SupplierTitlePrefix   string             `json:"supplier_title_prefix"`
	SupplierName          string             `json:"supplier_name"`
	SupplierTitleSuffix   string             `json:"supplier_title_suffix"`
	ClientTypeId          int                `json:"client_type_id"`
	TermOfPaymentId       int                `json:"term_of_payment_id"`
	DefaultCurrencyId     int                `json:"default_currency_id"`
	ViaBinning            bool               `json:"via_binning"`
	IsImportSupplier      bool               `json:"is_import_supplier"`
	IsSnpSupplier         bool               `json:"is_snp_supplier"`
	SupplierInvoiceTypeId int                `json:"supplier_invoice_type_id"`
	SupplierUniqueId      string             `json:"supplier_unique_id"`
	SupplierNitkuId       string             `json:"supplier_nitku_id"`
	SupplierAddressId     int                `json:"supplier_address_id"`
	SupplierAddress       SupplierAddress    `json:"supplier_address"`
	SupplierPhoneNo       string             `json:"supplier_phone_no"`
	SupplierFaxNo         string             `json:"supplier_fax_no"`
	SupplierMobilePhone   string             `json:"supplier_mobile_phone"`
	SupplierEmailAddress  string             `json:"supplier_email_address"`
	MinimumDownPayment    int                `json:"minimum_down_payment"`
	BehaviourId           int                `json:"behaviour_id"`
	SupplierCategoryId    int                `json:"supplier_category_id"`
	TaxIndustry           *float64           `json:"tax_industry"`
	SupplierStatusId      int                `json:"supplier_status_id"`
	SupplierStatus        string             `json:"supplier_status"`
	VatSupplierId         int                `json:"vat_supplier_id"`
	VatSupplier           VatSupplier        `json:"vat_supplier"`
	TaxSupplierId         int                `json:"tax_supplier_id"`
	TaxSupplier           TaxSupplier        `json:"tax_supplier"`
	SupplierContact       SupplierContact    `json:"supplier_contact"`
	SupplierBankAccount   SupplerBankAccount `json:"supplier_bank_account"`
}

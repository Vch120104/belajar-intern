package generalservicepayloads

import "encoding/json"

type CompanyReferenceBetByIdResponse struct {
	CurrencyId                int         `json:"currency_id"`
	CoaGroupId                int         `json:"coa_group_id"`
	OperationDiscountOuterKpp json.Number `json:"operation_discount_outer_kpp"`
	MarginOuterKpp            json.Number `json:"margin_outer_kpp"`
	AdjustmentReasonId        int         `json:"adjustment_reason_id"`
	LeadTimeUnitEtd           int         `json:"lead_time_unit_etd"`
	BankAccReceiveCompanyId   int         `json:"bank_acc_receive_company_id"`
	UnitWarehouseId           int         `json:"unit_warehouse_id"`
	TimeDifference            int         `json:"time_difference"`
	UseDms                    bool        `json:"use_dms"`
	UseJpcb                   bool        `json:"use_jpcb"`
	CheckMonthEnd             bool        `json:"check_month_end"`
	IsDistributor             bool        `json:"is_distributor"`
	WithVat                   bool        `json:"with_vat"`
}

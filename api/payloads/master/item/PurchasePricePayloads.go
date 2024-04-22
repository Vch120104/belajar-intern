package masteritempayloads

import "time"

type PurchasePriceResponse struct {
	PurchasePriceId            int       `json:"purchase_price_id"`
	SupplierId                 int       `json:"supplier_id"`
	SupplierCode               string    `json:"supplier_code"`
	SupplierName               string    `json:"supplier_name"`
	CurrencyId                 int       `json:"currency_id"`
	CurrencyCode               string    `json:"currency_code"`
	CurrencyName               string    `json:"currency_name"`
	PurchasePriceEffectiveDate time.Time `json:"purchase_price_effective_date"`
}

type PurchasePriceRequest struct {
	PurchasePriceId            int       `json:"purchase_price_id" parent_entity:"mtr_purchase_price" main_table:"mtr_purchase_price"`
	SupplierId                 int       `json:"supplier_id" parent_entity:"mtr_purchase_price"`
	CurrencyId                 int       `json:"currency_id" parent_entity:"mtr_purchase_price"`
	PurchasePriceEffectiveDate time.Time `json:"purchase_price_effective_date" parent_entity:"mtr_purchase_price"`
}

type PurchasePriceSupplierResponse struct {
	SupplierId   int    `json:"supplier_id"`
	SupplierCode string `json:"supplier_code"`
	SupplierName string `json:"supplier_name"`
}

type CurrencyResponse struct {
	CurrencyId   int    `json:"currency_id"`
	CurrencyCode string `json:"currency_code"`
	CurrencyName string `json:"currency_name"`
}

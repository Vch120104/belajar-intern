package masteritempayloads

import (
	"time"
)

type PurchasePriceDetailsResponse struct {
	Page       int                           `json:"page"`
	Limit      int                           `json:"limit"`
	TotalPages int                           `json:"total_pages"`
	TotalRows  int                           `json:"total_rows"`
	Data       []PurchasePriceDetailResponse `json:"data"`
}

type PurchasePriceResponse struct {
	PurchasePriceId            int                          `json:"purchase_price_id"`
	SupplierId                 int                          `json:"supplier_id"`
	SupplierCode               string                       `json:"supplier_code"`
	SupplierName               string                       `json:"supplier_name"`
	CurrencyId                 int                          `json:"currency_id"`
	CurrencyCode               string                       `json:"currency_code"`
	CurrencyName               string                       `json:"currency_name"`
	PurchasePriceEffectiveDate time.Time                    `json:"purchase_price_effective_date"`
	IsActive                   bool                         `json:"is_active"`
	IdentitySysNumber          int                          `json:"identity_system_number"`
	PurchasePriceDetails       PurchasePriceDetailsResponse `json:"purchase_price_details"`
}

type PurchasePriceRequest struct {
	PurchasePriceId            int       `json:"purchase_price_id" parent_entity:"mtr_purchase_price" main_table:"mtr_purchase_price"`
	SupplierId                 int       `json:"supplier_id" parent_entity:"mtr_purchase_price"`
	CurrencyId                 int       `json:"currency_id" parent_entity:"mtr_purchase_price"`
	PurchasePriceEffectiveDate time.Time `json:"purchase_price_effective_date" parent_entity:"mtr_purchase_price"`
	IsActive                   bool      `json:"is_active" parent_entity:"mtr_purchase_price"`
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

type PurchasePriceDetailRequest struct {
	PurchasePriceDetailId int     `json:"purchase_price_detail_id" parent_entity:"mtr_purchase_price_detail" main_table:"mtr_purchase_price_detail"`
	PurchasePriceId       int     `json:"purchase_price_id" parent_entity:"mtr_purchase_price_detail"`
	ItemId                int     `json:"item_id" parent_entity:"mtr_purchase_price_detail"`
	IsActive              bool    `json:"is_active" parent_entity:"mtr_purchase_price_detail"`
	PurchasePrice         float64 `json:"purchase_price" parent_entity:"mtr_purchase_price_detail"`
}

type PurchasePriceDetailResponses struct {
	PurchasePriceDetailId int     `json:"purchase_price_detail_id"`
	PurchasePriceId       int     `json:"purchase_price_id"`
	ItemId                int     `json:"item_id"`
	ItemCode              string  `json:"item_code"`
	ItemName              string  `json:"item_name"`
	IsActive              bool    `json:"is_active"`
	PurchasePrice         float64 `json:"purchase_price"`
}

type PurchasePriceDetailResponse struct {
	PurchasePriceDetailId int     `json:"purchase_price_detail_id"`
	PurchasePriceId       int     `json:"purchase_price_id"`
	ItemId                int     `json:"item_id"`
	ItemCode              string  `json:"item_code"`
	ItemName              string  `json:"item_name"`
	IsActive              bool    `json:"is_active"`
	PurchasePrice         float64 `json:"purchase_price"`
}

type PurchasePriceItemResponse struct {
	ItemId   int    `json:"item_id"`
	ItemCode string `json:"item_code"`
	ItemName string `json:"item_name"`
}

type PurchasePriceSubDetailResponse struct {
	PurchasePriceId            int                         `json:"purchase_price_id"`
	SupplierId                 int                         `json:"supplier_id"`
	SupplierCode               string                      `json:"supplier_code"`
	SupplierName               string                      `json:"supplier_name"`
	CurrencyId                 int                         `json:"currency_id"`
	CurrencyCode               string                      `json:"currency_code"`
	CurrencyName               string                      `json:"currency_name"`
	PurchasePriceEffectiveDate time.Time                   `json:"purchase_price_effective_date"`
	IsActive                   bool                        `json:"is_active"`
	PurchasePriceDetail        PurchasePriceDetailResponse `json:"purchase_price_detail"`
}

type PurchasePriceByIdResponse struct {
	PurchasePriceId int    `json:"purchase_price_id"`
	ItemId          int    `json:"item_id"`
	ItemCode        string `json:"item_code"`
	ItemName        string `json:"item_name"`
	PurchasePrice   int    `json:"purchase_price"`
}

type UploadRequest struct {
	Data []PurchasePriceDetailResponses `json:"data"`
}

type PurchasePriceDownloadResponse struct {
	PurchasePriceId            int                         `json:"purchase_price_id"`
	SupplierCode               string                      `json:"supplier_code"`
	SupplierName               string                      `json:"supplier_name"`
	CurrencyCode               string                      `json:"currency_code"`
	CurrencyName               string                      `json:"currency_name"`
	PurchasePriceEffectiveDate time.Time                   `json:"purchase_price_effective_date"`
	IsActive                   bool                        `json:"is_active"`
	PurchasePriceDetail        PurchasePriceDetailResponse `json:"purchase_price_detail"`
}

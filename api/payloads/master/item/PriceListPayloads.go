package masteritempayloads

import "time"

type PriceListResponse struct {
	IsActive            bool      `json:"is_active"`
	PriceListId         int       `json:"price_list_id"`
	PriceListCode       string    `json:"price_list_code"`
	CompanyId           int       `json:"company_id"`
	BrandId             int       `json:"brand_id"`
	CurrencyId          int       `json:"currency_id"`
	EffectiveDate       time.Time `json:"effective_date"`
	ItemId              int       `json:"item_id"`
	ItemGroupId         int       `json:"item_group_id"`
	ItemClassId         int       `json:"item_class_id"`
	PriceListAmount     float64   `json:"price_list_amount"`
	PriceListModifiable bool      `json:"price_list_modifiable"`
	AtpmSyncronize      bool      `json:"atpm_syncronize"`
	AtpmSyncronizeTime  time.Time `json:"atpm_syncronize_time"`
}

type PriceListRequest struct {
	IsActive            bool      `json:"is_active"`
	PriceListCode       string    `json:"price_list_code"`
	CompanyId           int       `json:"company_id"`
	BrandId             int       `json:"brand_id"`
	CurrencyId          int       `json:"currency_id"`
	EffectiveDate       time.Time `json:"effective_date"`
	ItemId              int       `json:"item_id"`
	ItemGroupId         int       `json:"item_group_id"`
	ItemClassId         int       `json:"item_class_id"`
	PriceListAmount     float64   `json:"price_list_amount"`
	PriceListModifiable bool      `json:"price_list_modifiable"`
	AtpmSyncronize      bool      `json:"atpm_syncronize"`
	AtpmSyncronizeTime  time.Time `json:"atpm_syncronize_time"`
}

type PriceListGetAllRequest struct {
	IsActive            string    `json:"is_active"`
	PriceListCode       string    `json:"price_list_code"`
	CompanyId           int       `json:"company_id"`
	BrandId             int       `json:"brand_id"`
	CurrencyId          int       `json:"currency_id"`
	EffectiveDate       time.Time `json:"effective_date"`
	ItemId              int       `json:"item_id"`
	ItemGroupId         int       `json:"item_group_id"`
	ItemClassId         int       `json:"item_class_id"`
	PriceListAmount     float64   `json:"price_list_amount"`
	PriceListModifiable string    `json:"price_list_modifiable"`
	AtpmSyncronize      string    `json:"atpm_syncronize"`
	AtpmSyncronizeTime  time.Time `json:"atpm_syncronize_time"`
}

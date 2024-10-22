package masteritempayloads

import "time"

type PriceListResponse struct {
	IsActive            bool      `json:"is_active"`
	PriceListId         int       `json:"price_list_id"`
	PriceListCodeId     int       `json:"price_list_code_id"`
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

type PriceListItemResponses struct {
	PriceListId     int     `json:"price_list_id"`
	ItemCode        string  `json:"item_code"`
	ItemName        string  `json:"item_name"`
	PriceListAmount float64 `json:"price_list_amount"`
	IsActive        bool    `json:"is_active"`
	ItemId          int     `json:"item_id"`
	ItemClassId     int     `json:"item_class_id"`
}

type PriceListRequest struct {
	IsActive            bool      `json:"is_active"`
	PriceListId         int       `json:"price_list_id"`
	PriceListCodeId     int       `json:"price_list_code_id"`
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

type SavePriceListMultiple struct {
	BrandId         int       `json:"brand_id"`
	CurrencyId      int       `json:"currency_id"`
	EffectiveDate   time.Time `json:"effective_date"`
	ItemGroupId     int       `json:"item_group_id"`
	CompanyId       int       `json:"company_id"`
	PriceListCodeId int       `json:"price_list_code_id"`

	Detail []PriceListItemResponses `json:"detail"`
}

type PriceListUploadDataRequest struct {
	BrandName       string `json:"brand_name" validate:"required"`
	BrandId         int    `json:"brand_id" validate:"required"`
	ItemGroupCode   string `json:"item_group_code" validate:"required"`
	ItemGroupId     int    `json:"item_group_id" validate:"required"`
	CurrencyCode    string `json:"currency_code" validate:"required"`
	CurrencyId      int    `json:"currency_id" validate:"required"`
	Date            string `json:"date" validate:"required"`
	PriceListCodeId int    `json:"price_list_code_id" validate:"required"`
	CompanyCode     string `json:"company_code" validate:"required"`
	CompanyId       int    `json:"company_id" `
}

type PriceListProcessdDataRequest struct {
	BrandId         int     `json:"brand_id" validate:"required"`
	ItemGroupId     int     `json:"item_group_id" validate:"required"`
	CurrencyId      int     `json:"currency_id" validate:"required"`
	Date            string  `json:"date" validate:"required"`
	PriceListCodeId int     `json:"price_list_code_id" validate:"required"`
	CompanyId       int     `json:"company_id" `
	ItemId          int     `json:"item_id"`
	ItemClassId     int     `json:"item_class_id"`
	PriceListAmount float64 `json:"price_list_amount"`
}

type PriceListGetAllRequest struct {
	IsActive            string    `json:"is_active" parent_entity:"mtr_item_price_list"`
	PriceListId         int       `json:"price_list_id" parent_entity:"mtr_item_price_list" main_table:"mtr_item_price_list"`
	PriceListCode       string    `json:"price_list_code" parent_entity:"mtr_item_price_list"`
	CompanyId           int       `json:"company_id" parent_entity:"mtr_item_price_list"`
	BrandId             int       `json:"brand_id" parent_entity:"mtr_item_price_list"`
	CurrencyId          int       `json:"currency_id" parent_entity:"mtr_item_price_list"`
	EffectiveDate       time.Time `json:"effective_date" parent_entity:"mtr_item_price_list"`
	ItemId              int       `json:"item_id" parent_entity:"mtr_item"`
	ItemName            string    `json:"item_name"  parent_entity:"mtr_item"`
	ItemCode            string    `json:"item_code"  parent_entity:"mtr_item"`
	ItemGroupId         int       `json:"item_group_id" parent_entity:"mtr_item_group"`
	ItemClassId         int       `json:"item_class_id" parent_entity:"mtr_item_class"`
	ItemClassName       string    `json:"item_class_name" parent_entity:"mtr_item_class"`
	PriceListAmount     float64   `json:"price_list_amount" parent_entity:"mtr_item_price_list"`
	PriceListModifiable string    `json:"price_list_modifiable" parent_entity:"mtr_item_price_list"`
	AtpmSyncronize      string    `json:"atpm_syncronize" parent_entity:"mtr_item_price_list"`
	AtpmSyncronizeTime  time.Time `json:"atpm_syncronize_time" parent_entity:"mtr_item_price_list"`
}

type PriceListGetAllResponse struct {
	IsActive            string  `json:"is_active" parent_entity:"mtr_item_price_list"`
	PriceListId         int     `json:"price_list_id" parent_entity:"mtr_item_price_list" main_table:"mtr_item_price_list"`
	PriceListCode       string  `json:"price_list_code" parent_entity:"mtr_item_price_list"`
	CompanyId           int     `json:"company_id" parent_entity:"mtr_item_price_list"`
	BrandId             int     `json:"brand_id" parent_entity:"mtr_item_price_list"`
	CurrencyId          int     `json:"currency_id" parent_entity:"mtr_item_price_list"`
	EffectiveDate       string  `json:"effective_date" parent_entity:"mtr_item_price_list"`
	ItemId              int     `json:"item_id" parent_entity:"mtr_item"`
	ItemName            string  `json:"item_name"  parent_entity:"mtr_item"`
	ItemCode            string  `json:"item_code"  parent_entity:"mtr_item"`
	ItemGroupId         int     `json:"item_group_id" parent_entity:"mtr_item_group"`
	ItemClassId         int     `json:"item_class_id" parent_entity:"mtr_item_class"`
	ItemClassName       string  `json:"item_class_name" parent_entity:"mtr_item_class"`
	ItemPriceCode       string  `json:"item_price_code" parent_entity:"mtr_item_price_code"`
	PriceListAmount     float64 `json:"price_list_amount" parent_entity:"mtr_item_price_list"`
	PriceListModifiable string  `json:"price_list_modifiable" parent_entity:"mtr_item_price_list"`
	AtpmSyncronize      string  `json:"atpm_syncronize" parent_entity:"mtr_item_price_list"`
	AtpmSyncronizeTime  string  `json:"atpm_syncronize_time" parent_entity:"mtr_item_price_list"`
}

type PriceListGetbyId struct {
	IsActive            string  `json:"is_active" parent_entity:"mtr_item_price_list"`
	PriceListId         int     `json:"price_list_id" parent_entity:"mtr_item_price_list" main_table:"mtr_item_price_list"`
	PriceListCode       string  `json:"price_list_code" parent_entity:"mtr_item_price_list"`
	PriceListCodeId     int     `json:"price_list_code_id"`
	CompanyId           int     `json:"company_id" parent_entity:"mtr_item_price_list"`
	BrandId             int     `json:"brand_id" parent_entity:"mtr_item_price_list"`
	BrandName           string  `json:"brand_name"`
	CurrencyId          int     `json:"currency_id" parent_entity:"mtr_item_price_list"`
	CurrencyCode        string  `json:"currency_code"`
	EffectiveDate       string  `json:"effective_date" parent_entity:"mtr_item_price_list"`
	ItemId              int     `json:"item_id" parent_entity:"mtr_item"`
	ItemName            string  `json:"item_name"  parent_entity:"mtr_item"`
	ItemCode            string  `json:"item_code"  parent_entity:"mtr_item"`
	ItemGroupId         int     `json:"item_group_id" parent_entity:"mtr_item_group"`
	ItemGroupName       string  `json:"item_group_name"`
	ItemClassId         int     `json:"item_class_id" parent_entity:"mtr_item_class"`
	ItemClassName       string  `json:"item_class_name" parent_entity:"mtr_item_class"`
	PriceListAmount     float64 `json:"price_list_amount" parent_entity:"mtr_item_price_list"`
	PriceListModifiable string  `json:"price_list_modifiable" parent_entity:"mtr_item_price_list"`
	AtpmSyncronize      string  `json:"atpm_syncronize" parent_entity:"mtr_item_price_list"`
	AtpmSyncronizeTime  string  `json:"atpm_syncronize_time" parent_entity:"mtr_item_price_list"`
}

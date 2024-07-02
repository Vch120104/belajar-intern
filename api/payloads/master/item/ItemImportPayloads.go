package masteritempayloads

import (
	masteritementities "after-sales/api/entities/master/item"
	"encoding/json"
)

type ItemImportResponse struct {
	ItemImportId int    `json:"item_import_id"`
	ItemId       int    `json:"item_id"`
	ItemCode     string `json:"item_code"`
	ItemName     string `json:"item_name"`
	SupplierId   int    `json:"supplier_id"`
}

type ItemImportByIdResponse struct {
	ItemImportId       int     `json:"item_import_id"`
	ItemId             int     `json:"item_id"`
	ItemCode           string  `json:"item_code"`
	ItemName           string  `json:"item_name"`
	SupplierId         int     `json:"supplier_id"`
	SupplierName       string  `json:"supplier_name"`
	SupplierCode       string  `json:"supplier_code"`
	ItemAliasCode      string  `json:"item_alias_code"`
	ItemAliasName      string  `json:"item_alias_name"`
	RoyaltyFlag        string  `json:"royalty_flag"`
	OrderConversion    float64 `json:"order_conversion"`
	OrderQtyMultiplier float64 `json:"order_qty_multiplier"`
}

type SupplierResponse struct {
	SupplierId   int    `json:"supplier_id"`
	SupplierName string `json:"supplier_name"`
	SupplierCode string `json:"supplier_code"`
}

type ItemImportUploadResponse struct {
	ItemCode           string  `json:"part_number"`
	SupplierCode       string  `json:"supplier_code"`
	ItemAliasCode      string  `json:"part_number_alias"`
	ItemAliasName      string  `json:"part_name_alias"`
	OrderQtyMultiplier float64 `json:"moq"`
	RoyaltyFlag        string  `json:"royalty"`
	OrderConversion    float64 `json:"order_conversion"`
}

type ItemImportUploadRequest struct {
	Data []masteritementities.ItemImport `json:"data"`
}

func ConvertItemImportMapToStruct(maps []map[string]any) ([]ItemImportByIdResponse, error) {
	var result []ItemImportByIdResponse
	// Marshal the maps into JSON
	jsonData, err := json.Marshal(maps)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON into the struct
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

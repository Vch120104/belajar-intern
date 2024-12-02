package masteritempayloads

import (
	masteritementities "after-sales/api/entities/master/item"
	"encoding/json"
	"fmt"
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

	// Ensure that the maps slice is not nil or empty
	if len(maps) == 0 {
		return result, nil // Return empty result if no data is provided
	}

	// Marshal the maps into JSON format
	jsonData, err := json.Marshal(maps)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal maps to JSON: %w", err)
	}

	// Unmarshal the JSON data into the result slice
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON into ItemImportByIdResponse: %w", err)
	}

	return result, nil
}

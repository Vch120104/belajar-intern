package masteritempayloads

type ItemImportResponse struct {
	ItemImportId int    `json:"item_import_id"`
	ItemId       int    `json:"item_id"`
	ItemCode     string `json:"item_code"`
	ItemName     string `json:"item_name"`
	SupplierId   int    `json:"supplier_id"`
}

type ItemImportByIdResponse struct {
	ItemImportId            int     `json:"item_import_id"`
	ItemId                  int     `json:"item_id"`
	ItemCode                string  `json:"item_code"`
	ItemName                string  `json:"item_name"`
	SupplierId              int     `json:"supplier_id"`
	SupplierName            string  `json:"supplier_name"`
	SupplierCode            string  `json:"supplier_code"`
	ItemAliasCode           string  `json:"item_alias_code"`
	ItemAliasName           string  `json:"item_alias_name"`
	Royalty                 string  `json:"royalty_flag"`
	OrderConversion         float64 `json:"order_conversion"`
	OrderQuantityMultiplier float64 `json:"order_qty_multiplier"`
}
type SupplierResponse struct {
	SupplierId   int    `json:"supplier_id"`
	SupplierName string `json:"supplier_name"`
	SupplierCode string `json:"supplier_code"`
}

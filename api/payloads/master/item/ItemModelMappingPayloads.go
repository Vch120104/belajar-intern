package masteritempayloads

type CreateItemModelMapping struct {
	ItemDetailId int     `json:"item_detail_id"`
	IsActive     bool    `json:"is_active"`
	ItemId       int     `json:"item_id"`
	BrandId      int     `json:"brand_id"`
	ModelId      int     `json:"model_id"`
	VariantId    int     `json:"variant_id"`
	MileageEvery float64 `json:"mileage_every"`
	ReturnEvery  float64 `json:"return_every"`
}

type ItemModelMappingReponses struct {
	ItemDetailId int     `json:"item_detail_id"`
	IsActive     bool    `json:"is_active"`
	ItemId       int     `json:"item_id"`
	BrandId      int     `json:"brand_id"`
	ModelId      int     `json:"model_id"`
	VariantId    int     `json:"variant_id"`
	MileageEvery float64 `json:"mileage_every"`
	ReturnEvery  float64 `json:"return_every"`
}

type UnitBrandResponses struct {
	BrandId   int    `json:"brand_id"`
	BrandName string `json:"brand_name"`
}

type UnitModelResponses struct {
	ModelId          int    `json:"model_id"`
	ModelCode        string `json:"model_code"`
	ModelDescription string `json:"model_description"`
}

type UnitVariantResponses struct {
	VariantId          int    `json:"variant_id"`
	VariantCode        string `json:"variant_code"`
	VariantDescription string `json:"variant_description"`
}

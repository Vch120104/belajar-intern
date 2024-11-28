package transactionworkshoppayloads

type LicenseOwnerChange struct {
	BrandId         int    `json:"brand_id"`
	ModelId         int    `json:"model_id"`
	VehicleId       int    `json:"vehicle_id"`
	ChangeType      string `json:"change_type"`
	TnkbOld         string `json:"tnkb_old"`
	TnkbNew         string `json:"tnkb_new"`
	OwnerNameOld    string `json:"owner_name_old"`
	OwnerNameNew    string `json:"owner_name_new"`
	OwnerAddressOld string `json:"owner_address_old"`
	OwnerAdressNew  string `json:"owner_address_new"`
}

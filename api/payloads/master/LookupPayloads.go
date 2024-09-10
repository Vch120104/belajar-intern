package masterpayloads

type ItemOprCodeResponse struct {
	ItemOprCodeId int    `json:"item_opr_code_id"`
	ItemOprCode   string `json:"item_opr_code"`
	ItemOprDesc   string `json:"item_opr_desc"`
}

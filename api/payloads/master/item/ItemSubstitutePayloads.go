package masteritempayloads

import "time"

type ItemSubstitutePayloads struct {
	SubstituteTypeCode string    `json:"substitute_type_code"`
	ItemSubstituteId   int       `json:"item_substitute_is"`
	EffectiveDate      time.Time `json:"effective_date"`
	ItemId             int       `json:"item_id"`
}
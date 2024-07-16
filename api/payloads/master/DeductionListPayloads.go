package masterpayloads

import (
	"time"
)

type DeductionListResponse struct {
	IsActive      bool      `json:"is_active"`
	DeductionId   int       `json:"deduction_id"`
	DeductionCode string    `json:"deduction_code"`
	DeductionName string    `json:"deduction_name"`
	EffectiveDate time.Time `json:"effective_date"`
}

type DeductionListPostResponse struct {
	DeductionName string    `json:"deduction_name"`
	DeductionCode string    `json:"deduction_code"`
	EffectiveDate time.Time `json:"effective_date"`
}

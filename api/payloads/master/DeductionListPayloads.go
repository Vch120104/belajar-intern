package masterpayloads

import (
	"time"
)

type DeductionListResponse struct {
	IsActive        bool      `json:"is_active"`
	DeductionListId int       `json:"deduction_list_id"`
	DeductionName   string    `json:"deduction_name"`
	EffectiveDate   time.Time `json:"effective_date"`
}

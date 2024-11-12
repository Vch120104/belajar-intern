package financeserviceapiutils

import (
	config2 "after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"net/http"
)

type JournalEntryApprovalUpdatePayload struct {
	IsVoid         bool   `json:"is_void"`
	IsApprove      bool   `json:"is_approve"`
	ApprovalLastBy *int   `json:"approval_last_by"`
	ApprovalRemark string `json:"approval_remark"`
}

func UpdateApprovalJournalEntry(payload JournalEntryApprovalUpdatePayload, journalId string) (bool, *exceptions.BaseErrorResponse) {
	updateApprovalJournalEntryUrl := config2.EnvConfigs.FinanceServiceUrl + "journal-entry/approval/" + journalId
	err := utils.Put(updateApprovalJournalEntryUrl, payload, nil)
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	return true, nil
}

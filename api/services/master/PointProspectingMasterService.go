package masterservice

import (
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type PointProspectingService interface {
	CreatePointProspecting(masterpayloads.PointProspectingRequest) (bool, *exceptions.BaseErrorResponse)
	UpdatePointProspectingStatus(string, int, masterpayloads.PointProspectingUpdateStatus) (bool, *exceptions.BaseErrorResponse)
	GetAllPointProspecting([]utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetOnePointProspecting(string, int) (masterpayloads.PointProspectingResponse, *exceptions.BaseErrorResponse)
	UpdatePointProspectingData(string, int, masterpayloads.PointProspectingUpdateRequest) (bool, *exceptions.BaseErrorResponse)
}

package masterrepository

import (
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type PointProspectingRepository interface {
	CreatePointProspecting(*gorm.DB, masterpayloads.PointProspectingRequest) (bool, *exceptions.BaseErrorResponse)
	UpdatePointProspectingStatus(*gorm.DB, string, int, masterpayloads.PointProspectingUpdateStatus) (bool, *exceptions.BaseErrorResponse)
	GetAllPointProspecting(*gorm.DB, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetOnePointProspecting(*gorm.DB, string, int) (masterpayloads.PointProspectingResponse, *exceptions.BaseErrorResponse)
	UpdatePointProspectingData(*gorm.DB, string, int, masterpayloads.PointProspectingUpdateRequest) (bool, *exceptions.BaseErrorResponse)
}

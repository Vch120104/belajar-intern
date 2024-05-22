package masterrepository

import (
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type CampaignMasterRepository interface {
	PostCampaignMaster(*gorm.DB, masterpayloads.CampaignMasterPost) (bool, *exceptions.BaseErrorResponse)
	PostCampaignDetailMaster(*gorm.DB, masterpayloads.CampaignMasterDetailPayloads) (bool, *exceptions.BaseErrorResponse)
	PostCampaignMasterDetailFromHistory(*gorm.DB, int, int) (bool, *exceptions.BaseErrorResponse)
	ChangeStatusCampaignMaster(*gorm.DB, int) (bool, *exceptions.BaseErrorResponse)
	ActivateCampaignMasterDetail(*gorm.DB, string, int) (bool, *exceptions.BaseErrorResponse)
	DeactivateCampaignMasterDetail(*gorm.DB, string, int) (bool, *exceptions.BaseErrorResponse)
	GetByIdCampaignMaster(*gorm.DB, int) ([]map[string]interface{}, *exceptions.BaseErrorResponse)
	GetByIdCampaignMasterDetail(*gorm.DB, int, int) (map[string]interface{}, *exceptions.BaseErrorResponse)
	GetAllCampaignMasterCodeAndName(*gorm.DB, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetAllCampaignMaster(*gorm.DB, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetAllCampaignMasterDetail(*gorm.DB, pagination.Pagination, int) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	UpdateCampaignMasterDetail(*gorm.DB, int, masterpayloads.CampaignMasterDetailPayloads) (bool, *exceptions.BaseErrorResponse)
	GetAllPackageMasterToCopy(tx *gorm.DB, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	SelectFromPackageMaster(tx *gorm.DB, id int, idhead int) (bool, *exceptions.BaseErrorResponse)
}

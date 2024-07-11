package masterrepository

import (
	mastercampaignmasterentities "after-sales/api/entities/master/campaign_master"
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type CampaignMasterRepository interface {
	PostCampaignMaster(*gorm.DB, masterpayloads.CampaignMasterPost) (mastercampaignmasterentities.CampaignMaster, *exceptions.BaseErrorResponse)
	PostCampaignDetailMaster(*gorm.DB, masterpayloads.CampaignMasterDetailPayloads) (int, *exceptions.BaseErrorResponse)
	PostCampaignMasterDetailFromHistory(*gorm.DB, int, int) (int, *exceptions.BaseErrorResponse)
	ChangeStatusCampaignMaster(*gorm.DB, int) (bool, *exceptions.BaseErrorResponse)
	ActivateCampaignMasterDetail(*gorm.DB, string, int) (bool,int, *exceptions.BaseErrorResponse)
	DeactivateCampaignMasterDetail(*gorm.DB, string, int) (bool,int, *exceptions.BaseErrorResponse)
	GetByIdCampaignMaster(*gorm.DB, int) ([]map[string]interface{}, *exceptions.BaseErrorResponse)
	GetByIdCampaignMasterDetail(*gorm.DB, int, int) (map[string]interface{}, *exceptions.BaseErrorResponse)
	GetAllCampaignMasterCodeAndName(*gorm.DB, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetAllCampaignMaster(*gorm.DB, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetAllCampaignMasterDetail(*gorm.DB, pagination.Pagination, int) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	UpdateCampaignMasterDetail(*gorm.DB, int,int, masterpayloads.CampaignMasterDetailPayloads) (int, *exceptions.BaseErrorResponse)
	GetAllPackageMasterToCopy(tx *gorm.DB, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	SelectFromPackageMaster(tx *gorm.DB, id int, idhead int) (int, *exceptions.BaseErrorResponse)
}

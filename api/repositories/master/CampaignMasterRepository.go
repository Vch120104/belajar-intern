package masterrepository

import (
	masterentities "after-sales/api/entities/master"
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type CampaignMasterRepository interface {
	PostCampaignMaster(*gorm.DB, masterpayloads.CampaignMasterPost) (masterentities.CampaignMaster, *exceptions.BaseErrorResponse)
	PostCampaignDetailMaster(*gorm.DB, masterpayloads.CampaignMasterDetailPayloads, int) (masterentities.CampaignMasterDetail, *exceptions.BaseErrorResponse)
	PostCampaignMasterDetailFromHistory(*gorm.DB, int, int) (int, *exceptions.BaseErrorResponse)
	PostCampaignMasterDetailFromPackage(*gorm.DB, masterpayloads.CampaignMasterDetailPostFromPackageRequest) (masterentities.CampaignMasterDetail, *exceptions.BaseErrorResponse)
	ChangeStatusCampaignMaster(*gorm.DB, int) (bool, *exceptions.BaseErrorResponse)
	ActivateCampaignMasterDetail(*gorm.DB, string) (bool, *exceptions.BaseErrorResponse)
	DeactivateCampaignMasterDetail(*gorm.DB, string) (bool, *exceptions.BaseErrorResponse)
	GetByIdCampaignMaster(*gorm.DB, int) (masterpayloads.CampaignMasterResponse, *exceptions.BaseErrorResponse)
	GetByIdCampaignMasterDetail(*gorm.DB, int) (map[string]interface{}, *exceptions.BaseErrorResponse)
	GetByCodeCampaignMaster(*gorm.DB, string) (masterpayloads.CampaignMasterResponse, *exceptions.BaseErrorResponse)
	GetAllCampaignMasterCodeAndName(*gorm.DB, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetAllCampaignMaster(*gorm.DB, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetAllCampaignMasterDetail(*gorm.DB, pagination.Pagination, int) (pagination.Pagination, *exceptions.BaseErrorResponse)
	UpdateCampaignMasterDetail(*gorm.DB, int, masterpayloads.CampaignMasterDetailPayloads) (int, *exceptions.BaseErrorResponse)
	GetAllPackageMasterToCopy(tx *gorm.DB, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	SelectFromPackageMaster(tx *gorm.DB, id int, idhead int) (int, *exceptions.BaseErrorResponse)
}

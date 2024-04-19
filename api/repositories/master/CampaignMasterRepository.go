package masterrepository

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type CampaignMasterRepository interface {
	PostCampaignMaster(*gorm.DB,masterpayloads.CampaignMasterPost)(bool,*exceptionsss_test.BaseErrorResponse)
	PostCampaignDetailMaster(*gorm.DB,masterpayloads.CampaignMasterDetailPayloads)(bool,*exceptionsss_test.BaseErrorResponse)
	PostCampaignMasterDetailFromHistory(*gorm.DB,int,int)(bool,*exceptionsss_test.BaseErrorResponse)
	ChangeStatusCampaignMaster(*gorm.DB,int)(bool,*exceptionsss_test.BaseErrorResponse)
	ActivateCampaignMasterDetail(*gorm.DB,string)(bool,*exceptionsss_test.BaseErrorResponse)
	DeactivateCampaignMasterDetail(*gorm.DB,string)(bool,*exceptionsss_test.BaseErrorResponse)
	GetByIdCampaignMaster(*gorm.DB,int)([]map[string]interface{},*exceptionsss_test.BaseErrorResponse)
	GetByIdCampaignMasterDetail(*gorm.DB,int)(masterpayloads.CampaignMasterDetailPayloads,*exceptionsss_test.BaseErrorResponse)
	GetAllCampaignMasterCodeAndName(*gorm.DB,pagination.Pagination)(pagination.Pagination,*exceptionsss_test.BaseErrorResponse)
	GetAllCampaignMaster(*gorm.DB,[]utils.FilterCondition,pagination.Pagination)(pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	GetAllCampaignMasterDetail(*gorm.DB,pagination.Pagination,int)(pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	UpdateCampaignMasterDetail(*gorm.DB,int,masterpayloads.CampaignMasterDetailPayloads)(bool,*exceptionsss_test.BaseErrorResponse)
}
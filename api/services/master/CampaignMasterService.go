package masterservice

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type CampaignMasterService interface {
	PostCampaignMaster(masterpayloads.CampaignMasterPost)(bool,*exceptionsss_test.BaseErrorResponse)
	PostCampaignDetailMaster(masterpayloads.CampaignMasterDetailPayloads)(bool,*exceptionsss_test.BaseErrorResponse)
	PostCampaignMasterDetailFromHistory(int,int)(bool,*exceptionsss_test.BaseErrorResponse)
	ChangeStatusCampaignMaster(int)(bool,*exceptionsss_test.BaseErrorResponse)
	ActivateCampaignMasterDetail(string)(bool,*exceptionsss_test.BaseErrorResponse)
	DeactivateCampaignMasterDetail(string)(bool,*exceptionsss_test.BaseErrorResponse)
	GetByIdCampaignMaster(int)([]map[string]interface{},*exceptionsss_test.BaseErrorResponse)
	GetByIdCampaignMasterDetail(int)(masterpayloads.CampaignMasterDetailPayloads,*exceptionsss_test.BaseErrorResponse)
	GetAllCampaignMasterCodeAndName(pagination.Pagination)(pagination.Pagination,*exceptionsss_test.BaseErrorResponse)
	GetAllCampaignMaster([]utils.FilterCondition,pagination.Pagination)(pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	GetAllCampaignMasterDetail(pagination.Pagination,int)(pagination.Pagination,*exceptionsss_test.BaseErrorResponse)
	UpdateCampaignMasterDetail(int,masterpayloads.CampaignMasterDetailPayloads)(bool,*exceptionsss_test.BaseErrorResponse)
	GetAllPackageMasterToCopy(pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	SelectFromPackageMaster(id int, idhead int) (bool, *exceptionsss_test.BaseErrorResponse)
}
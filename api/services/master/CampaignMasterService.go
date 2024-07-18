package masterservice

import (
	mastercampaignmasterentities "after-sales/api/entities/master/campaign_master"
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type CampaignMasterService interface {
	PostCampaignMaster(masterpayloads.CampaignMasterPost) (mastercampaignmasterentities.CampaignMaster, *exceptions.BaseErrorResponse)
	PostCampaignDetailMaster(masterpayloads.CampaignMasterDetailPayloads) (int, *exceptions.BaseErrorResponse)
	PostCampaignMasterDetailFromHistory(int, int) (int, *exceptions.BaseErrorResponse)
	ChangeStatusCampaignMaster(int) (bool, *exceptions.BaseErrorResponse)
	ActivateCampaignMasterDetail(ids string, id int) (bool,int, *exceptions.BaseErrorResponse)
	DeactivateCampaignMasterDetail(ids string, id int) (bool,int, *exceptions.BaseErrorResponse)
	GetByIdCampaignMaster(int) ([]map[string]interface{}, *exceptions.BaseErrorResponse)
	GetByIdCampaignMasterDetail(id int, idhead int) (map[string]interface{}, *exceptions.BaseErrorResponse)
	GetAllCampaignMasterCodeAndName(pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetAllCampaignMaster([]utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetAllCampaignMasterDetail(pages pagination.Pagination, id int) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	UpdateCampaignMasterDetail(int,int, masterpayloads.CampaignMasterDetailPayloads) (int, *exceptions.BaseErrorResponse)
	GetAllPackageMasterToCopy(pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	SelectFromPackageMaster(id int, idhead int) (int, *exceptions.BaseErrorResponse)
}

package masterservice

import (
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type CampaignMasterService interface {
	PostCampaignMaster(masterpayloads.CampaignMasterPost) (bool, *exceptions.BaseErrorResponse)
	PostCampaignDetailMaster(masterpayloads.CampaignMasterDetailPayloads) (bool, *exceptions.BaseErrorResponse)
	PostCampaignMasterDetailFromHistory(int, int) (bool, *exceptions.BaseErrorResponse)
	ChangeStatusCampaignMaster(int) (bool, *exceptions.BaseErrorResponse)
	ActivateCampaignMasterDetail(ids string, id int) (bool, *exceptions.BaseErrorResponse)
	DeactivateCampaignMasterDetail(ids string, id int) (bool, *exceptions.BaseErrorResponse)
	GetByIdCampaignMaster(int) ([]map[string]interface{}, *exceptions.BaseErrorResponse)
	GetByIdCampaignMasterDetail(id int, idhead int) (map[string]interface{}, *exceptions.BaseErrorResponse)
	GetAllCampaignMasterCodeAndName(pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetAllCampaignMaster([]utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetAllCampaignMasterDetail(pages pagination.Pagination, id int) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	UpdateCampaignMasterDetail(int, masterpayloads.CampaignMasterDetailPayloads) (bool, *exceptions.BaseErrorResponse)
	GetAllPackageMasterToCopy(pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	SelectFromPackageMaster(id int, idhead int) (bool, *exceptions.BaseErrorResponse)
}

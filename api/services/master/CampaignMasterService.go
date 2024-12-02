package masterservice

import (
	masterentities "after-sales/api/entities/master"
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type CampaignMasterService interface {
	PostCampaignMaster(masterpayloads.CampaignMasterPost) (masterentities.CampaignMaster, *exceptions.BaseErrorResponse)
	PostCampaignDetailMaster(masterpayloads.CampaignMasterDetailPayloads, int) (masterentities.CampaignMasterDetail, *exceptions.BaseErrorResponse)
	PostCampaignMasterDetailFromHistory(int, int) (int, *exceptions.BaseErrorResponse)
	PostCampaignMasterDetailFromPackage(masterpayloads.CampaignMasterDetailPostFromPackageRequest) (masterentities.CampaignMasterDetail, *exceptions.BaseErrorResponse)
	ChangeStatusCampaignMaster(int) (bool, *exceptions.BaseErrorResponse)
	ActivateCampaignMasterDetail(ids string) (bool, *exceptions.BaseErrorResponse)
	DeactivateCampaignMasterDetail(ids string) (bool, *exceptions.BaseErrorResponse)
	GetByIdCampaignMaster(int) (map[string]interface{}, *exceptions.BaseErrorResponse)
	GetByIdCampaignMasterDetail(id int) (map[string]interface{}, *exceptions.BaseErrorResponse)
	GetByCodeCampaignMaster(code string) (map[string]interface{}, *exceptions.BaseErrorResponse)
	GetAllCampaignMasterCodeAndName(pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetAllCampaignMaster([]utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetAllCampaignMasterDetail(pages pagination.Pagination, id int) (pagination.Pagination, *exceptions.BaseErrorResponse)
	UpdateCampaignMasterDetail(int, masterpayloads.CampaignMasterDetailPayloads) (int, *exceptions.BaseErrorResponse)
	GetAllPackageMasterToCopy(pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	SelectFromPackageMaster(id int, idhead int) (int, *exceptions.BaseErrorResponse)
}

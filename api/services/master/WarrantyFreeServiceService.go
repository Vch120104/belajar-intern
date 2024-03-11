package masterservice

import masterpayloads "after-sales/api/payloads/master"

type WarrantyFreeServiceService interface {
	GetWarrantyFreeServiceById(Id int) map[string]interface{}
	SaveWarrantyFreeService(req masterpayloads.WarrantyFreeServiceRequest) bool
}

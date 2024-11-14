package generalserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"
)

type WorkOrderJobType struct {
	JobTypeId   int    `json:"work_order_job_type_id"`
	JobTypeCode string `json:"work_order_job_type_code"`
	JobTypeName string `json:"work_order_job_type_name"`
}

func GetJobTransactionTypeByID(id int) (WorkOrderJobType, *exceptions.BaseErrorResponse) {
	var jobType WorkOrderJobType
	url := config.EnvConfigs.GeneralServiceUrl + "work-order-job-type/" + strconv.Itoa(id)

	err := utils.CallAPI("GET", url, nil, &jobType)
	if err != nil {
		return jobType, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve job type",
			Err:        errors.New("error consuming external API while getting job type by ID"),
		}
	}
	return jobType, nil
}

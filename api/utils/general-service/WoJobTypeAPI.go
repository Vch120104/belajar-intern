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
	JobTypeName string `json:"work_order_job_type_description"`
}

func GetJobTransactionTypeByID(id int) (WorkOrderJobType, *exceptions.BaseErrorResponse) {
	var jobType WorkOrderJobType
	url := config.EnvConfigs.GeneralServiceUrl + "work-order-job-type/" + strconv.Itoa(id)

	err := utils.CallAPI("GET", url, nil, &jobType)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve job type due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "job type service is temporarily unavailable"
		}

		return jobType, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting job type by ID"),
		}
	}
	return jobType, nil
}

package financeserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"
)

type EventResponse struct {
	IsActive                   bool   `json:"is_active"`
	EventId                    int    `json:"event_id"`
	EventNo                    string `json:"event_no"`
	EventDescription           string `json:"event_description"`
	TransactionTypeId          int    `json:"transaction_type_id"`
	TransactionTypeDescription string `json:"transaction_type_description"`
	ProcessId                  int    `json:"process_id"`
	ProcessDescription         string `json:"process_description"`
	ContentFilter1             string `json:"content_filter1"`
	ContentFilter2             string `json:"content_filter2"`
}

func GetEventById(id int) (EventResponse, *exceptions.BaseErrorResponse) {
	EventUrl := config.EnvConfigs.FinanceServiceUrl + "event/" + strconv.Itoa(id)
	event := EventResponse{}
	err := utils.CallAPI("GET", EventUrl, nil, &event)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve event due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "event service is temporarily unavailable"
		}

		return event, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting event by id"),
		}
	}
	return event, nil
}

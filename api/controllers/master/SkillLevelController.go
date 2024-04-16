package mastercontroller

import (
	exceptionsss_test "after-sales/api/expectionsss"
	helper_test "after-sales/api/helper_testt"
	jsonchecker "after-sales/api/helper_testt/json/json-checker"
	"after-sales/api/payloads"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"
	"after-sales/api/validation"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type SkillLevelController interface {
	GetAllSkillLevel(writer http.ResponseWriter, request *http.Request)
	GetSkillLevelByID(writer http.ResponseWriter, request *http.Request)
	SaveSkillLevel(writer http.ResponseWriter, request *http.Request)
	ChangeStatusSkillLevel(writer http.ResponseWriter, request *http.Request)
}
type SkillLevelControllerImpl struct {
	SkillLevelService masterservice.SkillLevelService
}

func NewSkillLevelController(skillLevelService masterservice.SkillLevelService) SkillLevelController {
	return &SkillLevelControllerImpl{
		SkillLevelService: skillLevelService,
	}
}

func (r *SkillLevelControllerImpl) GetAllSkillLevel(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	queryParams := map[string]string{
		"is_active":               query.Get("is_active"),
		"skill_level_code":        query.Get("skill_level_code"),
		"skill_level_description": query.Get("skill_level_description"),
	}

	pagination := pagination.Pagination{
		Limit:  utils.NewGetQueryInt(query, "limit"),
		Page:   utils.NewGetQueryInt(query, "page"),
		SortOf: query.Get("sort_of"),
		SortBy: query.Get("sort_by"),
	}

	filterCondition := utils.BuildFilterCondition(queryParams)

	result, err := r.SkillLevelService.GetAllSkillLevel(filterCondition, pagination)
	if err != nil {
		exceptionsss_test.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}

func (r *SkillLevelControllerImpl) GetSkillLevelByID(writer http.ResponseWriter, request *http.Request) {

	skillLevelId, _ := strconv.Atoi(chi.URLParam(request, "skill_level_id"))

	result, err := r.SkillLevelService.GetSkillLevelById(skillLevelId)
	if err != nil {
		exceptionsss_test.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, utils.ModifyKeysInResponse(result), "Get Data Successfully!", http.StatusOK)
}

func (r *SkillLevelControllerImpl) SaveSkillLevel(writer http.ResponseWriter, request *http.Request) {

	var formRequest masterpayloads.SkillLevelResponse
	err := jsonchecker.ReadFromRequestBody(request, &formRequest)
	var message string

	if err != nil {
		exceptionsss_test.NewEntityException(writer, request, err)
		return
	}
	err = validation.ValidationForm(writer, request, formRequest)
	if err != nil {
		exceptionsss_test.NewBadRequestException(writer, request, err)
		return
	}

	create, err := r.SkillLevelService.SaveSkillLevel(formRequest)

	if err != nil {
		helper_test.ReturnError(writer, request, err)
		return
	}

	if formRequest.SkillLevelId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

func (r *SkillLevelControllerImpl) ChangeStatusSkillLevel(writer http.ResponseWriter, request *http.Request) {
	skillLevelId, _ := strconv.Atoi(chi.URLParam(request, "skill_level_id"))

	response, err := r.SkillLevelService.ChangeStatusSkillLevel(int(skillLevelId))

	if err != nil {
		exceptionsss_test.NewBadRequestException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}
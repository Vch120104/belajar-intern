package mastercontroller

import (

	// "after-sales/api/helper"

	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	jsonchecker "after-sales/api/helper/json/json-checker"
	"after-sales/api/payloads"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"
	"after-sales/api/validation"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	// "github.com/julienschmidt/httprouter"
)

type SkillLevelController interface {
	GetAllSkillLevel(writer http.ResponseWriter, request *http.Request)
	GetSkillLevelById(writer http.ResponseWriter, request *http.Request)
	GetSkillLevelByCode(writer http.ResponseWriter, request *http.Request)
	SaveSkillLevel(writer http.ResponseWriter, request *http.Request)
	ChangeStatusSkillLevel(writer http.ResponseWriter, request *http.Request)
	UpdateSkillLevel(writer http.ResponseWriter, request *http.Request)
}

type SkillLevelControllerImpl struct {
	SkillLevelService masterservice.SkillLevelService
}

func NewSkillLevelController(SkillLevelService masterservice.SkillLevelService) SkillLevelController {
	return &SkillLevelControllerImpl{
		SkillLevelService: SkillLevelService,
	}
}

// @Summary Get All Skill Level
// @Description REST API Skill Level
// @Accept json
// @Produce json
// @Tags Master : Skill Level
// @Security BearerAuth
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Param is_active query string false "is_active" Enums(true, false)
// @Param skill_level_code query string false "skill_level_code"
// @Param skill_level_description query string false "skill_level_description"
// @Param sort_by query string false "sort_by"
// @Param sort_of query string false "sort_of"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/skill-level/ [get]
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
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccessPagination(writer, result.Rows, "Get Data Successfully!", 200, result.Limit, result.Page, result.TotalRows, result.TotalPages)
}

// @Summary Skill Level By Id
// @Description REST API Skill Level
// @Accept json
// @Produce json
// @Tags Master : Skill Level
// @Security BearerAuth
// @param skill_level_id path int true "skill_level_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/skill-level/{skill_level_id} [get]
func (r *SkillLevelControllerImpl) GetSkillLevelById(writer http.ResponseWriter, request *http.Request) {
	skillLevelId, errA := strconv.Atoi(chi.URLParam(request, "skill_level_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	result, err := r.SkillLevelService.GetSkillLevelById(skillLevelId)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, utils.ModifyKeysInResponse(result), "Get Data Successfully!", http.StatusOK)
}

// @Summary Save Skill Level
// @Description REST API Skill Level
// @Accept json
// @Produce json
// @Tags Master : Skill Level
// @Security BearerAuth
// @param reqBody body masterpayloads.SkillLevelResponse true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/skill-level/ [post]
func (r *SkillLevelControllerImpl) SaveSkillLevel(writer http.ResponseWriter, request *http.Request) {

	var formRequest masterpayloads.SkillLevelResponse
	err := jsonchecker.ReadFromRequestBody(request, &formRequest)
	var message string

	if err != nil {
		exceptions.NewEntityException(writer, request, err)
		return
	}
	err = validation.ValidationForm(writer, request, formRequest)
	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	create, err := r.SkillLevelService.SaveSkillLevel(formRequest)

	if err != nil {
		helper.ReturnError(writer, request, err)
		return
	}

	if formRequest.SkillLevelId == 0 {
		message = "Create Data Successfully!"
	} else {
		message = "Update Data Successfully!"
	}

	payloads.NewHandleSuccess(writer, create, message, http.StatusOK)
}

// @Summary Change Status Skill Level
// @Description REST API Skill Level
// @Accept json
// @Produce json
// @Tags Master : Skill Level
// @Security BearerAuth
// @param skill_level_id path int true "skill_level_id"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/skill-level/{skill_level_id} [patch]
func (r *SkillLevelControllerImpl) ChangeStatusSkillLevel(writer http.ResponseWriter, request *http.Request) {
	SkillLevelId, errA := strconv.Atoi(chi.URLParam(request, "skill_level_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}

	response, err := r.SkillLevelService.ChangeStatusSkillLevel(int(SkillLevelId))

	if err != nil {
		exceptions.NewBadRequestException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, response, "Update Data Successfully!", http.StatusOK)
}

// @Summary Update Skill Level
// @Description REST API Skill Level
// @Accept json
// @Produce json
// @Tags Master : Skill Level
// @Security BearerAuth
// @param skill_level_id path int true "skill_level_id"
// @param reqBody body masterpayloads.SkillLevelResponse true "Form Request"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/skill-level/{skill_level_id} [put]
func (r *SkillLevelControllerImpl) UpdateSkillLevel(writer http.ResponseWriter, request *http.Request) {
	skill_level_id, errA := strconv.Atoi(chi.URLParam(request, "skill_level_id"))

	if errA != nil {
		exceptions.NewBadRequestException(writer, request, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Err: errors.New("failed to read request param, please check your param input")})
		return
	}
	var formRequest masterpayloads.SkillLevelResponse
	helper.ReadFromRequestBody(request, &formRequest)
	result, err := r.SkillLevelService.UpdateSkillLevel(formRequest, skill_level_id)
	if err != nil {
		exceptions.NewConflictException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, result, "Update Data Successfully!", http.StatusOK)
}

// @Summary Get Skill Level By Code
// @Description REST API Skill Level
// @Accept json
// @Produce json
// @Tags Master : Skill Level
// @Security BearerAuth
// @param skill_level_code path string true "skill_level_code"
// @Success 200 {object} payloads.Response
// @Failure 500,400,401,404,403,422 {object} exceptions.BaseErrorResponse
// @Router /v1/skill-level/{skill_level_code} [get]
func (r *SkillLevelControllerImpl) GetSkillLevelByCode(writer http.ResponseWriter, request *http.Request) {
	skillLevelCode := chi.URLParam(request, "skill_level_code")

	result, err := r.SkillLevelService.GetSkillLevelByCode(skillLevelCode)
	if err != nil {
		exceptions.NewNotFoundException(writer, request, err)
		return
	}

	payloads.NewHandleSuccess(writer, utils.ModifyKeysInResponse(result), "Get Data Successfully!", http.StatusOK)
}

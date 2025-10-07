package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"term-service/internal/term/dto/request"
	"term-service/internal/term/mappers"
	"term-service/internal/term/model"
	"term-service/internal/term/service"
	"term-service/pkg/helper"
)

type TermHandler struct {
	service service.TermService
}

func NewHandler(s service.TermService) *TermHandler {
	return &TermHandler{service: s}
}

func (h *TermHandler) CreateTerm(c *gin.Context) {
	var input request.CreateTermReqDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	startDate, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}
	endDate, err := time.Parse("2006-01-02", input.EndDate)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	if !helper.ValidateDateRange(startDate, endDate) {
		helper.SendError(c, http.StatusBadRequest, nil, "start_date must be before or equal to end_date")
		return
	}

	term := &model.Term{
		Title:     input.Title,
		StartDate: startDate,
		EndDate:   endDate,
	}

	res, err := h.service.CreateTerm(c.Request.Context(), term)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusCreated, "Success", res)
}

func (h *TermHandler) GetTerms4Web(c *gin.Context) {
	terms, err := h.service.GetTerms4Web(c.Request.Context())
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Success", terms)
}

func (h *TermHandler) GetTermByID(c *gin.Context) {
	id := c.Param("id")

	term, err := h.service.GetTermByID(c.Request.Context(), id)
	if err != nil {
		helper.SendError(c, http.StatusNotFound, err, helper.ErrNotFount)
		return
	}

	res := mappers.MapTermToResDTO(term)
	helper.SendSuccess(c, http.StatusOK, "Success", res)
}

func (h *TermHandler) GetCurrentTerm(c *gin.Context) {

	organizationID := c.Query("organization_id")
	// if organizationID == "" {
	// 	helper.SendError(c, http.StatusBadRequest, fmt.Errorf("missing organization_id in"), helper.ErrInvalidOperation)
	// 	return
	// }

	term, err := h.service.GetCurrentTermByOrg(c.Request.Context(), organizationID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInternal)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Success", term)
}

func (h *TermHandler) UploadTerm(c *gin.Context) {
	var req request.UploadTermRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	if err := h.service.UploadTerms(c.Request.Context(), req); err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, err.Error())
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Upload terms successfully", nil)
}

func (h *TermHandler) GetTermsByOrgID(c *gin.Context) {
	orgID := c.Param("organization_id")
	if orgID == "" {
		helper.SendError(c, http.StatusBadRequest, fmt.Errorf("missing organization_id"), helper.ErrInvalidOperation)
		return
	}

	terms, err := h.service.GetTermsByOrgID(c.Request.Context(), orgID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Success", terms)
}

func (h *TermHandler) GetTermsByStudent(c *gin.Context) {
	studentID := c.Param("student_id")
	if studentID == "" {
		helper.SendError(c, http.StatusBadRequest, fmt.Errorf("missing studentID"), helper.ErrInvalidOperation)
		return
	}

	terms, err := h.service.GetTermsByStudent(c.Request.Context(), studentID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Success", terms)
}

func (h *TermHandler) GetTerms4App(c *gin.Context) {
	organizationID := c.Query("organization_id")
	if organizationID == "" {
		helper.SendError(c, http.StatusBadRequest, fmt.Errorf("missing organization_id in"), helper.ErrInvalidOperation)
		return
	}

	terms, err := h.service.GetTerms4App(c.Request.Context(), organizationID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInternal)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Success", terms.Terms)
}

func (h *TermHandler) GetTerm4Gw(c *gin.Context) {
	termID := c.Param("term_id")
	if termID == "" {
		helper.SendError(c, http.StatusBadRequest, fmt.Errorf("missing term_id in"), helper.ErrInvalidOperation)
		return
	}
	res, err := h.service.GetTerm4Gw(c.Request.Context(), termID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInternal)
		return
	}
	helper.SendSuccess(c, http.StatusOK, "Success", res)
}

func (h *TermHandler) GetTermsByOrg4App(c *gin.Context) {
	orgID := c.Param("organization_id")
	if orgID == "" {
		helper.SendError(c, http.StatusBadRequest, fmt.Errorf("missing organization_id"), helper.ErrInvalidOperation)
		return
	}

	terms, err := h.service.GetTermsByOrg4App(c.Request.Context(), orgID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Success", terms)
}

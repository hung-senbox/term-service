package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"term-service/helper"
	"term-service/internal/term/dto/request"
	"term-service/internal/term/mappers"
	"term-service/internal/term/model"
	"term-service/internal/term/service"
	pkg_helpder "term-service/pkg/helper"
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

	if !pkg_helpder.ValidateDateRange(startDate, endDate) {
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

func (h *TermHandler) ListTerms(c *gin.Context) {
	terms, err := h.service.ListTerms(c.Request.Context())
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

func (h *TermHandler) UpdateTerm(c *gin.Context) {
	id := c.Param("id")

	var input request.UpdateTermReqDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	// Fetch existing term
	existing, err := h.service.GetTermByID(c.Request.Context(), id)
	if err != nil {
		helper.SendError(c, http.StatusNotFound, err, helper.ErrNotFount)
		return
	}

	// Only update provided fields
	if input.Title != nil {
		existing.Title = *input.Title
	}

	if input.StartDate != nil {
		startDate, err := time.Parse("2006-01-02", *input.StartDate)
		if err != nil {
			helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
			return
		}
		existing.StartDate = startDate
	}

	if input.EndDate != nil {
		endDate, err := time.Parse("2006-01-02", *input.EndDate)
		if err != nil {
			helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
			return
		}
		existing.EndDate = endDate
	}

	// ✅ Validate date range
	if !pkg_helpder.ValidateDateRange(existing.StartDate, existing.EndDate) {
		helper.SendError(c, http.StatusBadRequest, nil, "start_date must be before or equal to end_date")
		return
	}

	// ✅ Save updates
	if err := h.service.UpdateTerm(c.Request.Context(), id, existing); err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	// Fetch updated term
	updated, err := h.service.GetTermByID(c.Request.Context(), id)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInternal)
		return
	}

	res := mappers.MapTermToResDTO(updated)
	helper.SendSuccess(c, http.StatusOK, "Updated successfully", res)
}

func (h *TermHandler) DeleteTerm(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.DeleteTerm(c.Request.Context(), id); err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *TermHandler) GetCurrentTerm(c *gin.Context) {
	term, err := h.service.GetCurrentTerm(c.Request.Context())
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInternal)
		return
	}
	if term == nil {
		helper.SendError(c, http.StatusNotFound, nil, "No current term found")
		return
	}

	res := mappers.MapTermToCurrentResDTO(term)

	helper.SendSuccess(c, http.StatusOK, "Success", res)
}

func (h *TermHandler) UploadTerm(c *gin.Context) {
	var terms []request.UploadTermItem
	if err := c.ShouldBindJSON(&terms); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	if err := h.service.UploadTerms(c.Request.Context(), terms); err != nil {
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

func (h *TermHandler) GetTerms4Student(c *gin.Context) {
	studentID := c.Param("student_id")
	if studentID == "" {
		helper.SendError(c, http.StatusBadRequest, fmt.Errorf("missing studentID"), helper.ErrInvalidOperation)
		return
	}

	terms, err := h.service.GetTerms4Student(c.Request.Context(), studentID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Success", terms)
}

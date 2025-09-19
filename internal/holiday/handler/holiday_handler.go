package handler

import (
	"net/http"
	"term-service/internal/holiday/dto/request"
	"term-service/internal/holiday/service"
	"term-service/pkg/helper"

	"github.com/gin-gonic/gin"
)

type HolidayHandler struct {
	service service.HolidayService
}

func NewHandler(s service.HolidayService) *HolidayHandler {
	return &HolidayHandler{service: s}
}

func (h *HolidayHandler) GetHolidays4Web(c *gin.Context) {
	holidays, err := h.service.GetHolidays4Web(c.Request.Context())
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Success", holidays)
}

func (h *HolidayHandler) UploadHolidays(c *gin.Context) {
	var req request.UploadHolidayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	if err := h.service.UploadHolidays(c.Request.Context(), req); err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, err.Error())
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Upload holidays successfully", nil)
}

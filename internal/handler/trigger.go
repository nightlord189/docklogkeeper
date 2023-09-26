package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/nightlord189/docklogkeeper/internal/entity"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
)

// GetTriggers godoc
// @Tags trigger
// @Accept  json
// @Produce json
// @Param trigger_id query int false "trigger ID"
// @Success 200 {object} GetTriggersResponse
// @Failure 400 {object} GenericResponse
// @Failure 401 {object} GenericResponse
// @Failure 500 {object} GenericResponse
// @Router /api/trigger [Get]
// @BasePath /
func (h *Handler) GetTriggers(c *gin.Context) {
	var triggerID int64
	triggerIDStr := c.Query("trigger_id")
	if triggerIDStr != "" {
		parsed, err := strconv.ParseInt(triggerIDStr, 10, 64)
		if err != nil {
			log.Ctx(c.Request.Context()).Err(err).Msg("parse trigger_id error")
			c.JSON(http.StatusBadRequest, GenericErrorf("invalid trigger_id: %v", err))
		}
		triggerID = parsed
	}
	triggers, err := h.Usecase.GetTriggers(triggerID)
	if err != nil {
		log.Ctx(c.Request.Context()).Err(err).Msg("get triggers error")
		c.JSON(http.StatusInternalServerError, GenericError(err.Error()))
	}
	c.JSON(http.StatusOK, GetTriggersResponse{Records: triggers})
}

// CreateTrigger godoc
// @Tags trigger
// @Accept  json
// @Produce json
// @Param data body CreateTriggerRequest true "Input model"
// @Success 201 {object} entity.TriggerDB
// @Failure 401 {object} GenericResponse
// @Failure 500 {object} GenericResponse
// @Router /api/trigger [Post]
// @BasePath /
func (h *Handler) CreateTrigger(c *gin.Context) {
	var req CreateTriggerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, GenericErrorf("parse json error: %v", err.Error()))
		return
	}

	dbEntity := req.ToDB()
	if err := dbEntity.IsValid(); err != nil {
		c.JSON(http.StatusBadRequest, GenericError(err.Error()))
		return
	}

	if err := h.Usecase.CreateTrigger(c.Request.Context(), &dbEntity); err != nil {
		log.Ctx(c.Request.Context()).Err(err).Msg("create trigger error")
		c.JSON(http.StatusInternalServerError, GenericError(err.Error()))
	}

	c.JSON(http.StatusCreated, dbEntity)
}

// UpdateTrigger godoc
// @Tags trigger
// @Accept  json
// @Produce json
// @Param id path int true "Trigger ID"
// @Param data body entity.TriggerDB true "Input model"
// @Success 200 {object} entity.TriggerDB
// @Failure 401 {object} GenericResponse
// @Failure 500 {object} GenericResponse
// @Router /api/trigger/{id} [Put]
// @BasePath /
func (h *Handler) UpdateTrigger(c *gin.Context) {
	id, err := getParamInt64(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, GenericError(err.Error()))
		return
	}

	var req entity.TriggerDB
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, GenericErrorf("parse json error: %v", err.Error()))
		return
	}

	req.ID = id

	if err := req.IsValid(); err != nil {
		c.JSON(http.StatusBadRequest, GenericError(err.Error()))
		return
	}

	if err := h.Usecase.UpdateTrigger(c.Request.Context(), &req); err != nil {
		log.Ctx(c.Request.Context()).Err(err).Msg("create trigger error")
		c.JSON(http.StatusInternalServerError, GenericError(err.Error()))
	}

	c.JSON(http.StatusOK, req)
}

// DeleteTrigger godoc
// @Tags trigger
// @Accept  json
// @Produce json
// @Param id path int true "Trigger ID"
// @Success 200 {object} GenericResponse
// @Failure 401 {object} GenericResponse
// @Failure 500 {object} GenericResponse
// @Router /api/trigger/{id} [Delete]
// @BasePath /
func (h *Handler) DeleteTrigger(c *gin.Context) {
	id, err := getParamInt64(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, GenericError(err.Error()))
		return
	}

	if err := h.Usecase.DeleteTrigger(c.Request.Context(), id); err != nil {
		log.Ctx(c.Request.Context()).Err(err).Msg("delete trigger error")
		c.JSON(http.StatusInternalServerError, GenericError(err.Error()))
	}

	c.JSON(http.StatusOK, GenericError("success"))
}

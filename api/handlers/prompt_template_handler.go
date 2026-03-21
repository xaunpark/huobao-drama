package handlers

import (
	"net/http"
	"strconv"

	"github.com/drama-generator/backend/application/services"
	"github.com/drama-generator/backend/pkg/logger"
	"github.com/drama-generator/backend/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PromptTemplateHandler struct {
	service *services.PromptTemplateService
	log     *logger.Logger
}

func NewPromptTemplateHandler(db *gorm.DB, log *logger.Logger) *PromptTemplateHandler {
	return &PromptTemplateHandler{
		service: services.NewPromptTemplateService(db, log),
		log:     log,
	}
}

func (h *PromptTemplateHandler) ListTemplates(c *gin.Context) {
	templates, err := h.service.List()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, templates)
}

func (h *PromptTemplateHandler) GetTemplate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid template ID")
		return
	}

	tpl, err := h.service.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, tpl)
}

func (h *PromptTemplateHandler) CreateTemplate(c *gin.Context) {
	var req services.CreatePromptTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	tpl, err := h.service.Create(&req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Created(c, tpl)
}

func (h *PromptTemplateHandler) UpdateTemplate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid template ID")
		return
	}

	var req services.UpdatePromptTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	tpl, err := h.service.Update(uint(id), &req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, tpl)
}

func (h *PromptTemplateHandler) DeleteTemplate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid template ID")
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		if err.Error() == "prompt template not found" {
			response.NotFound(c, err.Error())
		} else if err.Error() == "template is in use by one or more projects, cannot delete" {
			response.Error(c, http.StatusConflict, "CONFLICT", err.Error())
		} else {
			response.InternalError(c, err.Error())
		}
		return
	}
	response.SuccessWithMessage(c, "template deleted", nil)
}

func (h *PromptTemplateHandler) DuplicateTemplate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid template ID")
		return
	}

	tpl, err := h.service.Duplicate(uint(id))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Created(c, tpl)
}

func (h *PromptTemplateHandler) GetDefaultPrompts(c *gin.Context) {
	defaults := h.service.GetDefaultPrompts()
	response.Success(c, defaults)
}

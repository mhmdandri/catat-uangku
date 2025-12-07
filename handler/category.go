package handler

import (
	"catatan-keuangan/modules/categories"
	"net/http"

	"github.com/gin-gonic/gin"
)

type categoryHandler struct {
	categoryService categories.Service
}

func NewCategory(categoryService categories.Service) *categoryHandler {
	return &categoryHandler{categoryService}
}

func (h *categoryHandler) CreateCategoryHandler(c *gin.Context) {
	var categoryRequest categories.CategoryRequest
	if err := c.ShouldBindJSON(&categoryRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": formatValidationError(err),
		})
		return
	}
	NewCategory, err := h.categoryService.Create(categoryRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	categoryResponse := categories.FormatCategoryResponse(NewCategory)
	c.JSON(http.StatusCreated, gin.H{
		"data": categoryResponse,
	})
}

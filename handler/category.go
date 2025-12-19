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

// CreateCategoryHandler godoc
// @Summary Buat kategori
// @Tags Categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body categories.CategoryRequest true "Data kategori"
// @Success 201 {object} CategoryDataResponse
// @Failure 400 {object} ErrorResponse
// @Router /category [post]
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

func (h *categoryHandler) GetAllCategoriesHandler(c *gin.Context) {
	categoriesData, err := h.categoryService.GetAllCategories()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	categoryResponses := categories.FormatCategoryResponses(categoriesData)
	c.JSON(http.StatusOK, gin.H{
		"data": categoryResponses,
	})
}

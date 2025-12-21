package handler

import (
	"catatan-keuangan/modules/categories"
	"catatan-keuangan/modules/common"
	"errors"
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
// @Failure 403 {object} ErrorResponse
// @Router /category [post]
func (h *categoryHandler) CreateCategoryHandler(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		return
	}
	var categoryRequest categories.CategoryRequest
	if err := c.ShouldBindJSON(&categoryRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": formatValidationError(err),
		})
		return
	}
	NewCategory, err := h.categoryService.Create(userID, categoryRequest)
	if err != nil {
		if errors.Is(err, common.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "akses tidak diizinkan"})
			return
		}
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

// GetAllCategoriesHandler godoc
// @Summary Daftar kategori
// @Tags Categories
// @Security BearerAuth
// @Produce json
// @Success 200 {object} CategoriesDataResponse
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /categories [get]
func (h *categoryHandler) GetAllCategoriesHandler(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		return
	}
	categoriesData, err := h.categoryService.GetAllCategories(userID)
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

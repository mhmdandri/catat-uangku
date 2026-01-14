package handler

import (
	"catatan-keuangan/modules/common"
	"catatan-keuangan/modules/groups"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type groupHandler struct {
	groupService groups.Service
}

func NewGroupHandler(groupService groups.Service) *groupHandler {
	return &groupHandler{groupService}
}

// CreateGroupHandler godoc
// @Summary Buat group baru
// @Tags Groups
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body groups.GroupRequest true "Data group"
// @Success 201 {object} GroupDataResponse
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /groups [post]
func (h *groupHandler) CreateGroupHandler(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		return
	}
	var groupRequest groups.GroupRequest
	if err := c.ShouldBindJSON(&groupRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Semua form wajib diisi",
		})
		return
	}
	newGroup, err := h.groupService.Create(userID, groupRequest)
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
	formatNewGroup := groups.FormatGroupResponse(newGroup)
	c.JSON(http.StatusCreated, gin.H{
		"data": formatNewGroup,
	})
}

// GetGroupByIDHandler godoc
// @Summary Detail group
// @Tags Groups
// @Security BearerAuth
// @Produce json
// @Param id path string true "Group ID"
// @Success 200 {object} GroupDataResponse
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /groups/{id} [get]
func (h *groupHandler) GetGroupByIDHandler(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "id tidak valid",
		})
		return
	}
	groupData, err := h.groupService.FindByID(userID, id)
	if err != nil {
		if errors.Is(err, common.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "akses tidak diizinkan"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Group tidak di temukan",
		})
		return
	}
	grouResponse := groups.FormatGroupResponse(groupData)
	c.JSON(http.StatusOK, gin.H{
		"data": grouResponse,
	})
}

// GetGroupByUserID godoc
// @Summary Group yang dimiliki user
// @Tags Groups
// @Security BearerAuth
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} GroupsDataResponse
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /groups/user/{id} [get]
func (h *groupHandler) GetGroupByUserID(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id tidak valid"})
		return
	}
	if id != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "akses tidak diizinkan"})
		return
	}
	groupsData, err := h.groupService.FindGroupByUserID(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Gagal mengambil group untuk user tersebut",
		})
		return
	}
	groupIDs := make([]uuid.UUID, 0, len(groupsData))
	for _, group := range groupsData {
		groupIDs = append(groupIDs, group.ID)
	}
	totalTransactions, err := h.groupService.CountGroupTransactions(groupIDs)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Gagal mengambil ringkasan grup",
		})
		return
	}
	groupsResponse := groups.FormatGroupResponses(groupsData)
	summary := groups.BuildGroupSummary(groupsData, totalTransactions)
	c.JSON(http.StatusOK, gin.H{
		"summary": summary,
		"data":    groupsResponse,
	})
}

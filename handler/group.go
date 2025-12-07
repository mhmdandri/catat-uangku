package handler

import (
	"catatan-keuangan/modules/groups"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type groupHandler struct {
	groupService groups.Service
}

func NewGroupHandler(groupService groups.Service) *groupHandler {
	return &groupHandler{groupService}
}

func (h *groupHandler) CreateGroupHandler(c *gin.Context) {
	var groupRequest groups.GroupRequest
	if err := c.ShouldBindJSON(&groupRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Semua form wajib diisi",
		})
		return
	}
	newGroup, err := h.groupService.Create(groupRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"data": newGroup,
	})
}

func (h *groupHandler) GetGroupByIDHandler(c *gin.Context) {
	idString := c.Param("id")
	id, _ := strconv.Atoi(idString)
	groupData, err := h.groupService.FindByID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Group tidak di temukan",
		})
		return
	}
	grouResponse := groups.FormatGroupResponse(groupData)
	c.JSON(http.StatusOK, gin.H{
		"data": grouResponse,
	})
}

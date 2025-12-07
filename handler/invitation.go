package handler

import (
	"catatan-keuangan/modules/invitations"
	"net/http"

	"github.com/gin-gonic/gin"
)

type invitationHandler struct {
	invitationService invitations.Service
}

func NewInvitationHandler(invitationService invitations.Service) *invitationHandler {
	return &invitationHandler{invitationService}
}

func (h *invitationHandler) SendInvitationGroupHandler(c *gin.Context) {
	var invitationRequest invitations.InvitationRequest
	if err := c.ShouldBindJSON(&invitationRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Semua form wajib diisi",
		})
		return
	}
	newInvitation, err := h.invitationService.SendInvitationGroup(invitationRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"data": newInvitation,
	})
}

func (h *invitationHandler) AcceptInvitationGroupHandler(c *gin.Context) {
	var tokenRequest invitations.AcceptInvitationRequest
	if err := c.ShouldBindJSON(&tokenRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Semua form wajib diisi",
		})
		return
	}
	acceptedInvitation, err := h.invitationService.AcceptInvitationGroup(tokenRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": acceptedInvitation,
	})
}

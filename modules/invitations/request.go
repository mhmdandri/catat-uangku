package invitations

import (
	"github.com/google/uuid"
)

type InvitationRequest struct {
	GroupID       uuid.UUID `json:"group_id" binding:"required"`
	InviterUserID uuid.UUID `json:"inviter_user_id" binding:"required"`
	InviteeEmail  string    `json:"invitee_email" binding:"required,email"`
	// Token			 string `json:"token"`
	Role string `json:"role" binding:"required"`
	// Status      	 string `json:"status"`
	//ExpiredAt time.Time `json:"expired_at"`
}

type AcceptInvitationRequest struct {
	GroupID      uuid.UUID `json:"group_id" binding:"required"`
	Token        string    `json:"token" binding:"required"`
	InviteeEmail string    `json:"invitee_email" binding:"required,email"`
	UserID       uuid.UUID `json:"user_id" binding:"required"`
}

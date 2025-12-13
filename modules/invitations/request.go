package invitations

import (
	"github.com/google/uuid"
)

type InvitationRequest struct {
	GroupID       uuid.UUID `json:"group_id" binding:"required"`
	InviterUserID uuid.UUID `json:"inviter_user_id" binding:"required"`
	InviteeEmail  string    `json:"invitee_email" binding:"required,email"`
	Role          string    `json:"role" binding:"required"`
}

type AcceptInvitationRequest struct {
	GroupID      uuid.UUID `json:"group_id" binding:"required"`
	Token        string    `json:"token" binding:"required"`
	InviteeEmail string    `json:"invitee_email" binding:"required,email"`
	UserID       uuid.UUID `json:"user_id" binding:"required"`
}

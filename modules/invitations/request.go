package invitations

import "time"

type InvitationRequest struct {
	GroupID       int    `json:"group_id"`
	InviterUserID int    `json:"inviter_user_id"`
	InviteeEmail  string `json:"invitee_email"`
	// Token			 string `json:"token"`
	Role string `json:"role"`
	// Status      	 string `json:"status"`
	ExpiredAt time.Time `json:"expired_at"`
}

type AcceptInvitationRequest struct {
	GroupID      int    `json:"group_id" binding:"required"`
	Token        string `json:"token" binding:"required"`
	InviteeEmail string `json:"invitee_email" binding:"required,email"`
	UserID       int    `json:"user_id" binding:"required"`
}

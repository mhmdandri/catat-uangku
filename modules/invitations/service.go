package invitations

import (
	"errors"
	"time"

	groupmembers "catatan-keuangan/modules/group_members"

	"github.com/google/uuid"
)

type Service interface {
	SendInvitationGroup(invitationRequest InvitationRequest) (Invitation, error)
	AcceptInvitationGroup(acceptRequest AcceptInvitationRequest) (Invitation, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) SendInvitationGroup(invitationRequest InvitationRequest) (Invitation, error) {
	token := uuid.New().String()
	expired := time.Now().Add(48 * time.Hour)

	invitation := Invitation{
		GroupID:       invitationRequest.GroupID,
		InviterUserID: invitationRequest.InviterUserID,
		InviteeEmail:  invitationRequest.InviteeEmail,
		Token:         token,
		Role:          invitationRequest.Role,
		Status:        "pending",
		ExpiredAt:     expired,
	}
	newInvitation, err := s.repository.SendInvitationGroup(invitation)
	return newInvitation, err
}

func (s *service) AcceptInvitationGroup(acceptRequest AcceptInvitationRequest) (Invitation, error) {
	invitation, err := s.repository.FindByToken(acceptRequest.Token)
	if err != nil {
		return Invitation{}, err
	}
	if invitation.GroupID != acceptRequest.GroupID {
		return Invitation{}, errors.New("invitation invalid")
	}
	if invitation.Status != "pending" {
		return Invitation{}, errors.New("invitation invalid")
	}
	if acceptRequest.InviteeEmail != invitation.InviteeEmail {
		return Invitation{}, errors.New("unauthorized email")
	}
	if time.Now().After(invitation.ExpiredAt) {
		return Invitation{}, errors.New("invitation expired")
	}
	memberExists, err := s.repository.IsUserAlreadyMember(invitation.GroupID, acceptRequest.UserID)
	if err != nil {
		return Invitation{}, err
	}
	if memberExists {
		return Invitation{}, errors.New("user already joined this group")
	}
	invitation.Status = "accepted"
	member := groupmembers.GroupMember{
		GroupID:  invitation.GroupID,
		UserID:   acceptRequest.UserID,
		Role:     invitation.Role,
		JoinedAt: time.Now(),
		IsActive: true,
	}
	acceptedInvitation, err := s.repository.AcceptInvitationGroup(invitation, member)
	return acceptedInvitation, err
}

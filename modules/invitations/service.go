package invitations

import (
	"errors"
	"strings"
	"time"

	"catatan-keuangan/modules/common"
	groupmembers "catatan-keuangan/modules/group_members"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service interface {
	SendInvitationGroup(userID uuid.UUID, invitationRequest InvitationRequest) (Invitation, error)
	AcceptInvitationGroup(userID uuid.UUID, acceptRequest AcceptInvitationRequest) (Invitation, error)
}

type service struct {
	repository Repository
	db         *gorm.DB
}

func NewService(repository Repository, db *gorm.DB) *service {
	return &service{repository: repository, db: db}
}

func (s *service) SendInvitationGroup(userID uuid.UUID, invitationRequest InvitationRequest) (Invitation, error) {
	isMember, err := s.repository.IsUserAlreadyMember(invitationRequest.GroupID, userID)
	if err != nil {
		return Invitation{}, err
	}
	if !isMember {
		return Invitation{}, common.ErrForbidden
	}
	token := uuid.New().String()
	expired := time.Now().Add(48 * time.Hour)

	invitation := Invitation{
		GroupID:       invitationRequest.GroupID,
		InviterUserID: userID,
		InviteeEmail:  invitationRequest.InviteeEmail,
		Token:         token,
		Role:          invitationRequest.Role,
		Status:        "pending",
		ExpiredAt:     expired,
	}
	newInvitation, err := s.repository.SendInvitationGroup(invitation)
	return newInvitation, err
}

func (s *service) AcceptInvitationGroup(userID uuid.UUID, acceptRequest AcceptInvitationRequest) (Invitation, error) {
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
	if time.Now().After(invitation.ExpiredAt) {
		return Invitation{}, errors.New("invitation expired")
	}
	var email string
	if err := s.db.Table("users").Select("email").Where("id = ?", userID).Scan(&email).Error; err != nil {
		return Invitation{}, err
	}
	if email == "" || !strings.EqualFold(email, invitation.InviteeEmail) {
		return Invitation{}, common.ErrForbidden
	}
	memberExists, err := s.repository.IsUserAlreadyMember(invitation.GroupID, userID)
	if err != nil {
		return Invitation{}, err
	}
	if memberExists {
		return Invitation{}, errors.New("user already joined this group")
	}
	invitation.Status = "accepted"
	member := groupmembers.GroupMember{
		GroupID:  invitation.GroupID,
		UserID:   userID,
		Role:     invitation.Role,
		JoinedAt: time.Now(),
		IsActive: true,
	}
	acceptedInvitation, err := s.repository.AcceptInvitationGroup(invitation, member)
	return acceptedInvitation, err
}

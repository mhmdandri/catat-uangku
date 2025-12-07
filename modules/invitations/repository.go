package invitations

import (
	groupmembers "catatan-keuangan/modules/group_members"

	"gorm.io/gorm"
)

type Repository interface {
	SendInvitationGroup(invitation Invitation) (Invitation, error)
	AcceptInvitationGroup(invitation Invitation, member groupmembers.GroupMember) (Invitation, error)
	FindByToken(token string) (Invitation, error)
	IsUserAlreadyMember(groupID, userID int) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) SendInvitationGroup(invitation Invitation) (Invitation, error) {
	err := r.db.Create(&invitation).Error
	return invitation, err
}

func (r *repository) AcceptInvitationGroup(invitation Invitation, member groupmembers.GroupMember) (Invitation, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&invitation).Error; err != nil {
			return err
		}
		if err := tx.Create(&member).Error; err != nil {
			return err
		}
		return nil
	})
	return invitation, err
}

func (r *repository) FindByToken(token string) (Invitation, error) {
	var invitation Invitation
	err := r.db.Where("token = ?", token).First(&invitation).Error
	return invitation, err
}

func (r *repository) IsUserAlreadyMember(groupID, userID int) (bool, error) {
	var count int64
	err := r.db.Model(&groupmembers.GroupMember{}).
		Where("group_id = ? AND user_id = ?", groupID, userID).
		Count(&count).Error
	return count > 0, err
}

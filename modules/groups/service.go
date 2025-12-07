package groups

import (
	groupmembers "catatan-keuangan/modules/group_members"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service interface {
	Create(groupRequest GroupRequest) (Group, error)
	FindByID(ID uuid.UUID) (Group, error)
	FindGroupByUserID(userID uuid.UUID) ([]Group, error)
}

type service struct {
	repository Repository
	db         *gorm.DB
}

func NewService(repository Repository, db *gorm.DB) *service {
	return &service{
		repository: repository,
		db:         db,
	}
}

func (s *service) Create(groupRequest GroupRequest) (Group, error) {
	var newGroup Group
	err := s.db.Transaction(func(tx *gorm.DB) error {
		group := Group{
			Name: groupRequest.Name,
			Type: groupRequest.Type,
		}
		txRepo := NewRepository(tx)
		createdGroup, err := txRepo.Create(group)
		if err != nil {
			return err
		}
		creatorMember := groupmembers.GroupMember{
			GroupID:  createdGroup.ID,
			UserID:   groupRequest.CreatorUserID,
			Role:     "owner",
			JoinedAt: time.Now(),
			IsActive: true,
		}
		if err := tx.Create(&creatorMember).Error; err != nil {
			return err
		}
		newGroup = createdGroup
		return nil
	})
	return newGroup, err
}

func (s *service) FindByID(ID uuid.UUID) (Group, error) {
	group, err := s.repository.FindByID(ID)
	return group, err
}

func (s *service) FindGroupByUserID(userID uuid.UUID) ([]Group, error) {
	groups, err := s.repository.FindGroupByUserID(userID)
	return groups, err
}

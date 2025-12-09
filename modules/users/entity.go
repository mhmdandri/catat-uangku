package users

import (
	"catatan-keuangan/modules/accounts"
	groupmembers "catatan-keuangan/modules/group_members"
	"catatan-keuangan/modules/invitations"
	userprofile "catatan-keuangan/modules/user_profile"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	Email     string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"type:varchar(255);not null" json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Accounts      []accounts.Account         `gorm:"foreignKey:OwnerUserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"accounts,omitempty"`
	GroupMembers  []groupmembers.GroupMember `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"group_members,omitempty"`
	InviterUserID []invitations.Invitation   `gorm:"foreignKey:InviterUserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"invitations,omitempty"`
	Profile       userprofile.UserProfile    `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"profile,omitempty"`
}

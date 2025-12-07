package users

import (
	"catatan-keuangan/modules/accounts"
	groupmembers "catatan-keuangan/modules/group_members"
	"catatan-keuangan/modules/invitations"

	"time"
)

type User struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	Email     string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"type:varchar(255);not null" json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Accounts      []accounts.Account         `gorm:"foreignKey:OwnerUserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"accounts,omitempty"`
	GroupMembers  []groupmembers.GroupMember `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"group_members,omitempty"`
	InviterUserID []invitations.Invitation   `gorm:"foreignKey:InviterUserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"invitations,omitempty"`
}

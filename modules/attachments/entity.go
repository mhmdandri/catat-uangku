package attachments

import (
	"time"

	"github.com/google/uuid"
)

type Attachment struct {
	ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	TransactionID uuid.UUID `gorm:"type:uuid;index;not null" json:"transaction_id"`
	FilePath      string    `gorm:"type:varchar(255);not null" json:"file_path"`
	FileName      string    `gorm:"type:varchar(100);not null" json:"file_name"`
	UploadedAt    time.Time `json:"uploaded_at"`
}

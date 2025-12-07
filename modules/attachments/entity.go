package attachments

import "time"

type Attachment struct {
	ID            int       `gorm:"primaryKey;autoIncrement" json:"id"`
	TransactionID int       `gorm:"index;not null" json:"transaction_id"`
	FilePath      string    `gorm:"type:varchar(255);not null" json:"file_path"`
	FileName      string    `gorm:"type:varchar(100);not null" json:"file_name"`
	UploadedAt    time.Time `json:"uploaded_at"`
}

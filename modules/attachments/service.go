package attachments

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service interface {
	SaveFiles(tx *gorm.DB, transactionID uuid.UUID, files []*multipart.FileHeader) ([]Attachment, error)
	BuildResponse([]Attachment, string) []AttachmentResponse
}

type service struct {
	repository Repository
	uploadDir  string
	baseURL    string
}

func NewService(repository Repository, uploadDir, baseURL string) *service {
	return &service{repository, uploadDir, baseURL}
}

func (s *service) SaveFiles(tx *gorm.DB, transactionID uuid.UUID, files []*multipart.FileHeader) ([]Attachment, error) {
	if len(files) == 0 {
		return []Attachment{}, nil
	}
	if err := os.MkdirAll(s.uploadDir, 0755); err != nil {
		return nil, err
	}
	var out []Attachment
	for _, fileHeader := range files {
		src, err := fileHeader.Open()
		if err != nil {
			return nil, err
		}
		defer src.Close()
		filename := fmt.Sprintf("%s_%d_%s", transactionID.String(), time.Now().UnixNano(), filepath.Base(fileHeader.Filename))
		dstPath := filepath.Join(s.uploadDir, filename)
		dst, err := os.Create(dstPath)
		if err != nil {
			return nil, err
		}
		if _, err := io.Copy(dst, src); err != nil {
			dst.Close()
			return nil, err
		}
		dst.Close()
		attact := Attachment{
			TransactionID: transactionID,
			FilePath:      dstPath,
			FileName:      fileHeader.Filename,
			UploadedAt:    time.Now(),
		}
		if err := tx.Create(&attact).Error; err != nil {
			return nil, err
		}
		out = append(out, attact)
	}
	return out, nil
}

func (s *service) BuildResponse(attachment []Attachment, base string) []AttachmentResponse {
	out := make([]AttachmentResponse, 0, len(attachment))
	for _, attach := range attachment {
		out = append(out, AttachmentResponse{
			ID:            attach.ID,
			TransactionID: attach.TransactionID,
			FileName:      attach.FileName,
			FileURL:       fmt.Sprintf("%s/%s", base, filepath.Base(attach.FilePath)),
			UploadedAt:    attach.UploadedAt.Format(time.RFC3339),
		})
	}
	return out
}

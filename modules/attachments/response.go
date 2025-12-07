package attachments

type AttachmentResponse struct {
	ID            int    `json:"id"`
	TransactionID int    `json:"transaction_id"`
	FileName      string `json:"file_name"`
	FileURL       string `json:"file_url"`
	UploadedAt    string `json:"uploaded_at"`
}

func FormatAttachmentResponse(att Attachment) AttachmentResponse {
	return AttachmentResponse{
		ID:            att.ID,
		TransactionID: att.TransactionID,
		FileName:      att.FileName,
		FileURL:       att.FilePath,
		UploadedAt:    att.UploadedAt.Format("2006-01-02 15:04:05"),
	}
}

func FormatAttachmentResponses(atts []Attachment) []AttachmentResponse {
	formatted := make([]AttachmentResponse, 0, len(atts))
	for _, att := range atts {
		formatted = append(formatted, FormatAttachmentResponse(att))
	}
	return formatted
}

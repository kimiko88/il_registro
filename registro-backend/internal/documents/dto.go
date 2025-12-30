package documents

type UploadDocumentRequest struct {
	Title       string `form:"title" binding:"required"`
	Description string `form:"description"`
	Type        string `form:"type" binding:"required"`
}

type DocumentResponse struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Type        string `json:"type"`
	Date        string `json:"date"`
}

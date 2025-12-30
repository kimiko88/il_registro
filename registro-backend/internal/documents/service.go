package documents

type Service interface {
	GetAllDocuments() ([]DocumentResponse, error)
	UploadDocument(userID uint, req UploadDocumentRequest) error
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) GetAllDocuments() ([]DocumentResponse, error) {
	return []DocumentResponse{}, nil
}

func (s *service) UploadDocument(userID uint, req UploadDocumentRequest) error {
	return nil
}

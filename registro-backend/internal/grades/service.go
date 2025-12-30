package grades

type Service interface {
	GetStudentGrades(studentID uint) ([]GradeResponse, error)
	AddGrade(teacherID uint, req CreateGradeRequest) error
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) GetStudentGrades(studentID uint) ([]GradeResponse, error) {
	return []GradeResponse{}, nil
}

func (s *service) AddGrade(teacherID uint, req CreateGradeRequest) error {
	return nil
}

package attendance

type Service interface {
	GetClassAttendance(classID uint, date string) ([]AttendanceResponse, error)
	MarkAttendance(req CreateAttendanceRequest) error
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) GetClassAttendance(classID uint, date string) ([]AttendanceResponse, error) {
	return []AttendanceResponse{}, nil
}

func (s *service) MarkAttendance(req CreateAttendanceRequest) error {
	return nil
}

package scheduling

type Service interface {
	GetClassSchedule(classID uint) ([]ScheduleResponse, error)
	AddSchedule(req CreateScheduleRequest) error
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) GetClassSchedule(classID uint) ([]ScheduleResponse, error) {
	return []ScheduleResponse{}, nil
}

func (s *service) AddSchedule(req CreateScheduleRequest) error {
	return nil
}

package homework_submissions

import (
	"errors"
)

type Service interface {
	SubmitHomework(s *HomeworkSubmission) error
	GetSubmission(id string) (*HomeworkSubmission, error)
	GetStudentSubmission(homeworkID, studentID string) (*HomeworkSubmission, error)
	ListHomeworkSubmissions(homeworkID string) ([]HomeworkSubmission, error)
	ListStudentSubmissions(studentID string) ([]HomeworkSubmission, error)
	GradeSubmission(id string, req GradeSubmissionRequest, teacherID string) error
	SyncGoogleClassroom(id, submissionID string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) SubmitHomework(sub *HomeworkSubmission) error {
	if sub.HomeworkID == "" || sub.StudentID == "" || sub.FileURL == "" {
		return errors.New("homeworkID, studentID e fileURL sono obbligatori")
	}
	if sub.FileType == "" {
		sub.FileType = "pdf"
	}
	return s.repo.UpsertSubmission(sub)
}

func (s *service) GetSubmission(id string) (*HomeworkSubmission, error) {
	return s.repo.GetSubmission(id)
}

func (s *service) GetStudentSubmission(homeworkID, studentID string) (*HomeworkSubmission, error) {
	return s.repo.GetSubmissionByHomeworkAndStudent(homeworkID, studentID)
}

func (s *service) ListHomeworkSubmissions(homeworkID string) ([]HomeworkSubmission, error) {
	if homeworkID == "" {
		return nil, errors.New("homeworkID obbligatorio")
	}
	return s.repo.ListSubmissionsByHomework(homeworkID)
}

func (s *service) ListStudentSubmissions(studentID string) ([]HomeworkSubmission, error) {
	if studentID == "" {
		return nil, errors.New("studentID obbligatorio")
	}
	return s.repo.ListSubmissionsByStudent(studentID)
}

func (s *service) GradeSubmission(id string, req GradeSubmissionRequest, teacherID string) error {
	if id == "" || teacherID == "" {
		return errors.New("id sottomissione e teacherID obbligatori")
	}
	if req.TeacherGrade != nil && (*req.TeacherGrade < 1.0 || *req.TeacherGrade > 10.0) {
		return errors.New("il voto deve essere compreso tra 1 e 10")
	}
	return s.repo.SaveGrading(id, req, teacherID)
}

func (s *service) SyncGoogleClassroom(id, submissionID string) error {
	if id == "" || submissionID == "" {
		return errors.New("id e submissionID obbligatori per la sincronizzazione Classroom")
	}
	return s.repo.SetClassroomSync(id, submissionID)
}

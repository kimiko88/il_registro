package grades

import "errors"

func ValidateGradeValue(value float64) error {
	if value < 0 || value > 10 {
		return errors.New("grade must be between 0 and 10")
	}
	return nil
}

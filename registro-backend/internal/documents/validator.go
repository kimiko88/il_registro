package documents

import "errors"

type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) ValidateDocument(d *Document, content string) error {
	if !d.IsValid() {
		return errors.New("invalid document metadata")
	}
	if len(content) < 10 {
		return errors.New("content too short")
	}
	return nil
}

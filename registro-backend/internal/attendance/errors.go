package attendance

import "errors"

var (
	ErrForbidden        = errors.New("forbidden")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrNotFound         = errors.New("not found")
	ErrAlreadyProcessed = errors.New("già stata elaborata")
	ErrOverlapping      = errors.New("overlapping justification")
	ErrFutureDate       = errors.New("future date not allowed")
	ErrInvalidInput     = errors.New("invalid input")
)

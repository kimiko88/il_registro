package admin

import "errors"

var (
	// Permission errors
	ErrSuperAdminRequired = errors.New("superadmin access required")
	ErrAdminRequired      = errors.New("admin or superadmin access required")
	ErrUnauthorized       = errors.New("unauthorized access")
	ErrSchoolNotFound     = errors.New("school not found")
	ErrInvalidSchoolID    = errors.New("invalid school ID")

	// Admin user errors
	ErrAdminNotFound    = errors.New("admin user not found")
	ErrAdminExists      = errors.New("admin user already exists")
	ErrCannotDeleteSelf = errors.New("cannot delete your own admin account")

	// Import/Export errors
	ErrInvalidCSVFormat = errors.New("invalid CSV format")
	ErrImportFailed     = errors.New("import operation failed")
	ErrExportFailed     = errors.New("export operation failed")
)

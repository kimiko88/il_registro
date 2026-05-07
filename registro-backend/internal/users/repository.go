package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrEmailExists  = errors.New("email already exists")
	ErrFiscalCode   = errors.New("fiscal code already exists")
)

type Repository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error  // Soft delete
	Restore(ctx context.Context, id string) error // Restore
	List(ctx context.Context, filter UserFilter) ([]User, int, error)
	ListByIDs(ctx context.Context, ids []string) ([]User, error)

	// Audit & Bulk
	LogAudit(ctx context.Context, log *AuditLog) error
	GetAuditLogs(ctx context.Context, userID string, limit, offset int) ([]AuditLog, int, error)
	BulkCreate(ctx context.Context, users []User) (int, []string, error) // Returns count, errors
	BulkDelete(ctx context.Context, ids []string) (int, error)

	// GDPR
	HardDelete(ctx context.Context, id string) error // Actual DB delete

	// Relationships
	IsGuardian(ctx context.Context, parentUserID string, studentUserID string) (bool, error)
	GetChildren(ctx context.Context, parentUserID string) ([]StudentChild, error)
	GetStudentsByClass(ctx context.Context, classID string) ([]User, error)
	GetStudentProfile(ctx context.Context, userID string) (string, error)
	GetParentProfile(ctx context.Context, userID string) (string, error)
	AddGuardian(ctx context.Context, studentProfileID, parentProfileID, relationship string) error
	RemoveGuardian(ctx context.Context, studentProfileID, parentProfileID string) error
	GetGuardians(ctx context.Context, studentProfileID string) ([]GuardianInfo, error)

	// Guardianship
	IsActive(ctx context.Context, id string) (bool, error)
}



type PostgresRepository struct {
	db *sql.DB
}

// ... existing NewRepository ...

// ... existing methods ...

func (r *PostgresRepository) GetChildren(ctx context.Context, parentUserID string) ([]StudentChild, error) {
	query := `
		SELECT s.id, u.id, u.first_name, u.last_name, COALESCE(c.name, 'N/A'), sc.name
		FROM student_parents sp
		JOIN parents p ON sp.parent_id = p.id
		JOIN students s ON sp.student_id = s.id
		JOIN users u ON s.user_id = u.id
		LEFT JOIN classes c ON s.class_id = c.id
		JOIN schools sc ON u.school_id = sc.id
		WHERE p.user_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, parentUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var children []StudentChild
	for rows.Next() {
		var c StudentChild
		if err := rows.Scan(&c.ID, &c.UserID, &c.FirstName, &c.LastName, &c.Class, &c.SchoolName); err != nil {
			return nil, err
		}
		children = append(children, c)
	}
	return children, nil
}

func NewRepository(db *sql.DB) Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, user *User) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() // Safe: no-op after Commit()

	query := `
		INSERT INTO users (
			id, email, password_hash, first_name, last_name, fiscal_code,
			role, school_id, is_active, email_verified, mfa_enabled, mfa_secret,
			phone_number, job_title, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10, $11, $12,
			$13, $14, $15, $16
		)
	`
	_, err = tx.ExecContext(ctx, query,
		user.ID, user.Email, user.PasswordHash, user.FirstName, user.LastName, user.FiscalCode,
		user.Role, user.SchoolID, user.IsActive, user.EmailVerified, user.MFAEnabled, user.MFASecret,
		user.PhoneNumber, user.JobTitle, user.CreatedAt, user.UpdatedAt,
	)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // Unique violation
				if strings.Contains(pqErr.Message, "email") {
					return ErrEmailExists
				}
				if strings.Contains(pqErr.Message, "fiscal_code") {
					return ErrFiscalCode
				}
			}
		}
		return err
	}

	// Create Role-Specific Profiles
	if user.SchoolID != nil {
		switch user.Role {
		case "student":
			studentQuery := `INSERT INTO students (user_id, school_id, class_id) VALUES ($1, $2, $3)`
			_, err = tx.ExecContext(ctx, studentQuery, user.ID, *user.SchoolID, user.ClassID)
		case "parent":
			parentQuery := `INSERT INTO parents (user_id, school_id) VALUES ($1, $2)`
			_, err = tx.ExecContext(ctx, parentQuery, user.ID, *user.SchoolID)
		case "teacher":
			teacherQuery := `INSERT INTO teachers (user_id, school_id) VALUES ($1, $2)`
			_, err = tx.ExecContext(ctx, teacherQuery, user.ID, *user.SchoolID)
		}
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*User, error) {
	query := `
		SELECT u.id, u.email, u.password_hash, u.first_name, u.last_name, COALESCE(u.fiscal_code, ''),
		       u.role, u.school_id, u.is_active, u.email_verified, u.mfa_enabled, COALESCE(u.phone_number, ''), COALESCE(u.job_title, ''),
		       u.created_at, u.updated_at, u.last_login, u.deleted_at, u.pseudonymized_at,
		       s.class_id, c.name, c.section
		FROM users u
		LEFT JOIN students s ON u.id = s.user_id
		LEFT JOIN classes c ON s.class_id = c.id
		WHERE u.id = $1
	`
	var u User
	var classID, className, classSection *string
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.FiscalCode,
		&u.Role, &u.SchoolID, &u.IsActive, &u.EmailVerified, &u.MFAEnabled, &u.PhoneNumber, &u.JobTitle,
		&u.CreatedAt, &u.UpdatedAt, &u.LastLogin, &u.DeletedAt, &u.PseudonymizedAt,
		&classID, &className, &classSection,
	)
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	u.ClassID = classID
	if className != nil {
		u.ClassName = className
	} else if classSection != nil {
		u.ClassName = classSection
	}
	return &u, err
}

func (r *PostgresRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `SELECT id, email, password_hash, role, is_active, mfa_enabled, school_id FROM users WHERE email = $1`
	var u User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.IsActive, &u.MFAEnabled, &u.SchoolID,
	)
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	return &u, err
}

func (r *PostgresRepository) Update(ctx context.Context, user *User) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		UPDATE users SET
			first_name = $1, last_name = $2, phone_number = $3, job_title = $4,
			is_active = $5, role = $6, school_id = $7, updated_at = $8,
			password_hash = $9, mfa_enabled = $10,
			fiscal_code = COALESCE(NULLIF($11, ''), fiscal_code)
		WHERE id = $12
	`
	res, err := tx.ExecContext(ctx, query,
		user.FirstName, user.LastName, user.PhoneNumber, user.JobTitle,
		user.IsActive, user.Role, user.SchoolID, time.Now(),
		user.PasswordHash, user.MFAEnabled,
		user.FiscalCode,
		user.ID,
	)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return ErrUserNotFound
	}

	// Update Student Class if role is student and ClassID is provided (or nil to unassign)
	if user.Role == "student" && user.SchoolID != nil {
		// We expect the service to pass the updated ClassID in the user struct
		// Check if student record exists
		// Upsert logic for students table
		studentQuery := `
			INSERT INTO students (user_id, school_id, class_id, updated_at)
			VALUES ($1, $2, $3, NOW())
			ON CONFLICT (user_id, school_id) 
			DO UPDATE SET class_id = $3, updated_at = NOW()
		`
		_, err = tx.ExecContext(ctx, studentQuery, user.ID, *user.SchoolID, user.ClassID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE users SET deleted_at = $1 WHERE id = $2`
	res, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *PostgresRepository) BulkDelete(ctx context.Context, ids []string) (int, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids)+1)
	args[0] = time.Now()
	
	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+2)
		args[i+1] = id
	}
	
	query := fmt.Sprintf("UPDATE users SET deleted_at = $1 WHERE id IN (%s)", strings.Join(placeholders, ","))
	
	res, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	
	rows, _ := res.RowsAffected()
	return int(rows), nil
}

func (r *PostgresRepository) Restore(ctx context.Context, id string) error {
	query := `UPDATE users SET deleted_at = NULL WHERE id = $1`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *PostgresRepository) List(ctx context.Context, filter UserFilter) ([]User, int, error) {
	baseQuery := `
		SELECT u.id, u.email, u.first_name, u.last_name, u.role, u.school_id, u.is_active, u.created_at, u.deleted_at, u.pseudonymized_at,
		       s.class_id, c.name, c.section
		FROM users u
		LEFT JOIN students s ON u.id = s.user_id
		LEFT JOIN classes c ON s.class_id = c.id
		WHERE 1=1
	`
	// Apply filters
	countQuery := `SELECT COUNT(*) FROM users u WHERE 1=1`
	var args []interface{}
	argCount := 1

	// Filters
	if !filter.IsDeleted {
		baseQuery += fmt.Sprintf(" AND u.deleted_at IS NULL")
		countQuery += fmt.Sprintf(" AND u.deleted_at IS NULL")
	}
	if filter.Role != "" {
		baseQuery += fmt.Sprintf(" AND u.role = $%d", argCount)
		countQuery += fmt.Sprintf(" AND u.role = $%d", argCount)
		args = append(args, filter.Role)
		argCount++
	}
	if len(filter.ExcludeRoles) > 0 {
		phs := make([]string, len(filter.ExcludeRoles))
		for i := range filter.ExcludeRoles {
			phs[i] = fmt.Sprintf("$%d", argCount)
			args = append(args, filter.ExcludeRoles[i])
			argCount++
		}
		baseQuery += fmt.Sprintf(" AND u.role NOT IN (%s)", strings.Join(phs, ","))
		countQuery += fmt.Sprintf(" AND u.role NOT IN (%s)", strings.Join(phs, ","))
	}
	if filter.SchoolID != nil {
		baseQuery += fmt.Sprintf(" AND u.school_id = $%d", argCount)
		countQuery += fmt.Sprintf(" AND u.school_id = $%d", argCount)
		args = append(args, *filter.SchoolID)
		argCount++
	}
	if filter.IsActive != nil {
		baseQuery += fmt.Sprintf(" AND u.is_active = $%d", argCount)
		countQuery += fmt.Sprintf(" AND u.is_active = $%d", argCount)
		args = append(args, *filter.IsActive)
		argCount++
	}
	if filter.Query != "" {
		q := "%" + filter.Query + "%"
		baseQuery += fmt.Sprintf(" AND (u.email ILIKE $%d OR u.first_name ILIKE $%d OR u.last_name ILIKE $%d OR u.fiscal_code ILIKE $%d)", argCount, argCount+1, argCount+2, argCount+3)
		countQuery += fmt.Sprintf(" AND (u.email ILIKE $%d OR u.first_name ILIKE $%d OR u.last_name ILIKE $%d OR u.fiscal_code ILIKE $%d)", argCount, argCount+1, argCount+2, argCount+3)
		args = append(args, q, q, q, q)
		argCount += 4
	}

	// Count total
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Sort and Paginate
	sortBy := "u.created_at"
	if filter.SortBy != "" {
		// Map sortBy fields to columns with table alias
		switch filter.SortBy {
		case "last_name":
			sortBy = "u.last_name"
		case "email":
			sortBy = "u.email"
		case "role":
			sortBy = "u.role"
		case "class_name": // Special case
			sortBy = "c.section"
		}
	}
	sortOrder := "DESC"
	if strings.ToUpper(filter.SortOrder) == "ASC" {
		sortOrder = "ASC"
	}

	baseQuery += fmt.Sprintf(" ORDER BY %s %s LIMIT $%d OFFSET $%d", sortBy, sortOrder, argCount, argCount+1)
	args = append(args, filter.PageSize, (filter.Page-1)*filter.PageSize)

	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		var classID, className, classSection *string

		if err := rows.Scan(
			&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.Role, &u.SchoolID, &u.IsActive, &u.CreatedAt, &u.DeletedAt, &u.PseudonymizedAt,
			&classID, &className, &classSection,
		); err != nil {
			return nil, 0, err
		}

		u.ClassID = classID

		// Use Name if available, fallback to section
		if className != nil {
			u.ClassName = className
		} else if classSection != nil {
			u.ClassName = classSection
		}

		users = append(users, u)
	}

	return users, total, nil
}

func (r *PostgresRepository) ListByIDs(ctx context.Context, ids []string) ([]User, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	// Construct IN clause
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}
	query := fmt.Sprintf("SELECT id, email, first_name, last_name, role FROM users WHERE id IN (%s)", strings.Join(placeholders, ","))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.Role); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *PostgresRepository) HardDelete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)
	return err
}

func (r *PostgresRepository) LogAudit(ctx context.Context, log *AuditLog) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO audit_logs (id, user_id, actor_id, action, details, ip_address, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		log.ID, log.UserID, log.ActorID, log.Action, log.Details, log.IPAddress, log.CreatedAt,
	)
	return err
}

func (r *PostgresRepository) GetAuditLogs(ctx context.Context, userID string, limit, offset int) ([]AuditLog, int, error) {
	// Simple implementation
	var logs []AuditLog
	var total int

	r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM audit_logs WHERE user_id = $1", userID).Scan(&total)

	rows, err := r.db.QueryContext(ctx,
		"SELECT id, user_id, actor_id, action, details, ip_address, created_at FROM audit_logs WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3",
		userID, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var l AuditLog
		rows.Scan(&l.ID, &l.UserID, &l.ActorID, &l.Action, &l.Details, &l.IPAddress, &l.CreatedAt)
		logs = append(logs, l)
	}
	return logs, total, nil
}

func (r *PostgresRepository) BulkCreate(ctx context.Context, users []User) (int, []string, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, nil, err
	}
	defer tx.Rollback()

	// Prepare COPY statement
	stmt, err := tx.PrepareContext(ctx, pq.CopyIn("users", "id", "email", "password_hash", "first_name", "last_name", "fiscal_code", "role", "created_at", "updated_at"))
	if err != nil {
		return 0, nil, err
	}
	defer stmt.Close()

	var errs []string
	count := 0

	for _, u := range users {
		_, err := stmt.ExecContext(ctx, u.ID, u.Email, u.PasswordHash, u.FirstName, u.LastName, u.FiscalCode, u.Role, u.CreatedAt, u.UpdatedAt)
		if err != nil {
			errs = append(errs, fmt.Sprintf("Failed to stage %s: %v", u.Email, err))
			continue
		}
		count++
	}

	if _, err := stmt.ExecContext(ctx); err != nil {
		return 0, nil, err
	}

	if err := tx.Commit(); err != nil {
		return 0, nil, err
	}

	return count, errs, nil
}

func (r *PostgresRepository) IsGuardian(ctx context.Context, parentUserID string, studentID string) (bool, error) {
	// The table is `student_parents`, linking `student_id` (from `students` table) and `parent_id` (from `parents` table).
	// However, our input is `parentUserID` (from JWT) and `studentID` (could be student USER ID or student PROFILE ID??).
	// The Grades Service refers to `StudentID`.
	// In `003_academic_structure.sql`, `students` table has `id` and `user_id`.
	// The `grades` table links to `students.id` (profile ID) usually, but key name is `student_id`.
	// Let's assume the API inputs typically use the Profile IDs (UUIDs from `students`/`parents` tables), NOT User IDs.
	// BUT, the JWT usually contains the USER ID.
	// So we need to map ParentUserID -> ParentID first, OR check via join.
	//
	// Query: Check if there is a row in student_parents where parent's user_id matches and student's id matches.
	//
	// Join: student_parents sp JOIN parents p ON sp.parent_id = p.id
	// WHERE p.user_id = $1 AND sp.student_id = $2
	//
	// Note: `studentID` param: Grades Service usually uses the Student Profile ID (`students.id`).
	// If the `studentID` param passed forces UserID, we'd adjust.
	// The prompt says "Verifica che ParentID sia genitore di StudentID in DB".
	// Assuming logic:
	// 1. Resolve Parent Profile ID from Parent User ID.
	// 2. Check relationship with Target Student ID.

	query := `
		SELECT EXISTS (
			SELECT 1 
			FROM student_parents sp
			JOIN parents p ON sp.parent_id = p.id
			WHERE p.user_id = $1 AND sp.student_id = $2
		)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, parentUserID, studentID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *PostgresRepository) IsActive(ctx context.Context, id string) (bool, error) {
	query := `SELECT is_active FROM users WHERE id = $1 AND deleted_at IS NULL`
	var isActive bool
	err := r.db.QueryRowContext(ctx, query, id).Scan(&isActive)
	if err == sql.ErrNoRows {
		return false, nil // User not found effectively means not active
	}
	if err != nil {
		return false, err
	}
	return isActive, nil
}
func (r *PostgresRepository) GetStudentsByClass(ctx context.Context, classID string) ([]User, error) {
	query := `SELECT u.id, u.email, u.first_name, u.last_name, u.fiscal_code, u.role, u.school_id, u.is_active, u.created_at, u.updated_at, s.id as student_id
	          FROM users u
	          JOIN students s ON u.id = s.user_id
	          WHERE s.class_id = $1 AND u.role = 'student' AND u.deleted_at IS NULL 
	          ORDER BY u.last_name, u.first_name`
	rows, err := r.db.QueryContext(ctx, query, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.FiscalCode, &u.Role, &u.SchoolID, &u.IsActive, &u.CreatedAt, &u.UpdatedAt, &u.StudentID); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}
func (r *PostgresRepository) GetStudentProfile(ctx context.Context, userID string) (string, error) {
	var id string
	err := r.db.QueryRowContext(ctx, "SELECT id FROM students WHERE user_id = $1", userID).Scan(&id)
	return id, err
}

func (r *PostgresRepository) GetParentProfile(ctx context.Context, userID string) (string, error) {
	var id string
	err := r.db.QueryRowContext(ctx, "SELECT id FROM parents WHERE user_id = $1", userID).Scan(&id)
	return id, err
}

func (r *PostgresRepository) AddGuardian(ctx context.Context, studentProfileID, parentProfileID, relationship string) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO student_parents (student_id, parent_id, relationship_type) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`, studentProfileID, parentProfileID, relationship)
	return err
}

func (r *PostgresRepository) RemoveGuardian(ctx context.Context, studentProfileID, parentProfileID string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM student_parents WHERE student_id = $1 AND parent_id = $2`, studentProfileID, parentProfileID)
	return err
}

func (r *PostgresRepository) GetGuardians(ctx context.Context, studentProfileID string) ([]GuardianInfo, error) {
	query := `
		SELECT sp.parent_id, u.id as user_id, u.first_name, u.last_name, u.email, sp.relationship_type
		FROM student_parents sp
		JOIN parents p ON sp.parent_id = p.id
		JOIN users u ON p.user_id = u.id
		WHERE sp.student_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, studentProfileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var guardians []GuardianInfo
	for rows.Next() {
		var g GuardianInfo
		if err := rows.Scan(&g.ID, &g.ParentUserID, &g.FirstName, &g.LastName, &g.Email, &g.RelationshipType); err != nil {
			return nil, err
		}
		guardians = append(guardians, g)
	}
	return guardians, nil
}

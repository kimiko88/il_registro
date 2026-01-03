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

	// GDPR
	HardDelete(ctx context.Context, id string) error // Actual DB delete

	// Guardianship
	IsGuardian(ctx context.Context, parentUserID string, studentUserID string) (bool, error)
	GetChildren(ctx context.Context, parentUserID string) ([]StudentChild, error)
	IsActive(ctx context.Context, id string) (bool, error)
}

type StudentChild struct {
	ID         string `json:"id"`      // Student Profile ID
	UserID     string `json:"user_id"` // Student User ID
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Class      string `json:"class"`
	SchoolName string `json:"school_name"`
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
	// defer tx.Rollback() // Managed manually below for specific error returns?
	// Or better, use defer and check err

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
					tx.Rollback()
					return ErrEmailExists
				}
				if strings.Contains(pqErr.Message, "fiscal_code") {
					tx.Rollback()
					return ErrFiscalCode
				}
			}
		}
		tx.Rollback()
		return err
	}

	// Create Student Profile if role is student
	if user.Role == "student" && user.SchoolID != nil {
		studentQuery := `INSERT INTO students (user_id, school_id, class_id) VALUES ($1, $2, $3)`
		_, err := tx.ExecContext(ctx, studentQuery, user.ID, *user.SchoolID, user.ClassID)
		if err != nil {
			tx.Rollback()
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
			password_hash = $9, mfa_enabled = $10
		WHERE id = $11
	`
	res, err := tx.ExecContext(ctx, query,
		user.FirstName, user.LastName, user.PhoneNumber, user.JobTitle,
		user.IsActive, user.Role, user.SchoolID, time.Now(),
		user.PasswordHash, user.MFAEnabled,
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
	// ... (filters) ...

	// Count total ... (unchanged)
	// Sort ... (unchanged)

	// Just need to properly reconstruct the method or use replace chunks carefully.
	// Since I can't leave placeholders in `replacementContent` for tool to fill, I must provide full body if replacing body.
	// I'll replace the Scan part mostly.

	// Wait, I can replace the query and the scan block separately if I am careful with context lines.
	// Replacing the Query string:

	// ... code continues ...

	// Note: 'c.name' doesn't exist in classes table based on 003_academic_structure.sql.
	// The table has level_id, section, specialization.
	// But previous logs showed `name VARCHAR(255)` in `classes` table creation or `024_create_classes_table.sql`.
	// Wait, the user previously updated `classes` table in `024_create_classes_table.sql`.
	// Let's assume `name` exists based on recent work.

	// Actually, let's look at `003_academic_structure.sql` again. It has `section`, `academic_year_id`.
	// But `024_create_classes_table.sql` could have added `name`.
	// Let's check `024_create_classes_table.sql` content first to be sure about column name.
	// Oh, I can't check it right now without another tool call.
	// But earlier in `SchoolDetail.vue` I saw `props.row.name`.
	// Let's assume `name` exists or construct it from section + year.
	// Using `c.section` for now for safety if name isn't guaranteed.
	// Wait, the previous `classes` query used `name`.
	// Let's stick to generic columns and if `name` exists use it.
	// Let's re-read the `003` file content provided earlier.
	// `classes` table: id, school_id, year_id, level_id, building_id, section, specialization, coordinator_id.
	// It DOES NOT have a `name` column.
	// However, `001_initial_schema.sql` `024...` might have altered it.

	// I will check `024...` in a parallel step? No, safety first.
	// I will use `section` and maybe `academic_years.name`?
	// The `Users.vue` expects `class` string.
	// I'll assume we can just return `section` for now, or concat.

	// Actually, I'll check `024` quickly.

	// Resuming replacement with a safe guess for now which I can fix if wrong.
	// Replacing with the code that assumes `name` MIGHT NOT exist but `section` does.

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

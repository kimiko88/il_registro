package payments

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, p *SchoolPayment) error
	GetByID(ctx context.Context, id string) (*SchoolPayment, error)
	ListByStudent(ctx context.Context, schoolID, studentID string) ([]*SchoolPayment, error)
	ListByParent(ctx context.Context, schoolID, parentUserID string) ([]*SchoolPayment, error)
	ListBySchool(ctx context.Context, schoolID string) ([]*SchoolPayment, error)
	Pay(ctx context.Context, id, payerUserID, method, transactionID, receiptNumber string) error
}

type PostgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, p *SchoolPayment) error {
	p.ID = uuid.New().String()
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	p.Status = "pending"

	query := `
		INSERT INTO school_payments (id, school_id, student_id, title, description, amount, due_date, status, created_at, updated_at)
		VALUES ($1, $2::uuid, $3::uuid, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.ExecContext(ctx, query,
		p.ID, p.SchoolID, p.StudentID, p.Title, p.Description, p.Amount, p.DueDate, p.Status, p.CreatedAt, p.UpdatedAt,
	)
	return err
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*SchoolPayment, error) {
	query := `
		SELECT sp.id, sp.school_id, sp.student_id, sp.title, COALESCE(sp.description, ''), sp.amount, sp.due_date,
		       sp.status, sp.paid_at, COALESCE(sp.payment_method, ''), COALESCE(sp.transaction_id, ''),
		       sp.payer_user_id, COALESCE(sp.receipt_number, ''), sp.created_at, sp.updated_at,
		       COALESCE(u.first_name || ' ' || u.last_name, '') AS student_name
		FROM school_payments sp
		LEFT JOIN users u ON sp.student_id = u.id
		WHERE sp.id = $1::uuid
	`
	p := &SchoolPayment{}
	var payerID sql.NullString
	var paidAt sql.NullTime
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.SchoolID, &p.StudentID, &p.Title, &p.Description, &p.Amount, &p.DueDate,
		&p.Status, &paidAt, &p.PaymentMethod, &p.TransactionID,
		&payerID, &p.ReceiptNumber, &p.CreatedAt, &p.UpdatedAt,
		&p.StudentName,
	)
	if err != nil {
		return nil, err
	}
	if paidAt.Valid {
		p.PaidAt = &paidAt.Time
	}
	if payerID.Valid {
		p.PayerUserID = &payerID.String
	}
	return p, nil
}

func (r *PostgresRepository) ListByStudent(ctx context.Context, schoolID, studentID string) ([]*SchoolPayment, error) {
	query := `
		SELECT sp.id, sp.school_id, sp.student_id, sp.title, COALESCE(sp.description, ''), sp.amount, sp.due_date,
		       sp.status, sp.paid_at, COALESCE(sp.payment_method, ''), COALESCE(sp.transaction_id, ''),
		       sp.payer_user_id, COALESCE(sp.receipt_number, ''), sp.created_at, sp.updated_at,
		       COALESCE(u.first_name || ' ' || u.last_name, '') AS student_name
		FROM school_payments sp
		LEFT JOIN users u ON sp.student_id = u.id
		WHERE (sp.school_id = $1::uuid OR $1 = '')
		  AND (
			sp.student_id = $2::uuid OR
			sp.student_id = (SELECT user_id FROM students WHERE id = $2::uuid) OR
			sp.student_id = (SELECT id FROM students WHERE user_id = $2::uuid)
		  )
		ORDER BY sp.due_date DESC
	`
	rows, err := r.db.QueryContext(ctx, query, schoolID, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*SchoolPayment
	for rows.Next() {
		p := &SchoolPayment{}
		var payerID sql.NullString
		var paidAt sql.NullTime
		if err := rows.Scan(
			&p.ID, &p.SchoolID, &p.StudentID, &p.Title, &p.Description, &p.Amount, &p.DueDate,
			&p.Status, &paidAt, &p.PaymentMethod, &p.TransactionID,
			&payerID, &p.ReceiptNumber, &p.CreatedAt, &p.UpdatedAt,
			&p.StudentName,
		); err != nil {
			return nil, err
		}
		if paidAt.Valid {
			p.PaidAt = &paidAt.Time
		}
		if payerID.Valid {
			p.PayerUserID = &payerID.String
		}
		list = append(list, p)
	}
	if list == nil {
		list = []*SchoolPayment{}
	}
	return list, rows.Err()
}

func (r *PostgresRepository) ListByParent(ctx context.Context, schoolID, parentUserID string) ([]*SchoolPayment, error) {
	query := `
		SELECT sp.id, sp.school_id, sp.student_id, sp.title, COALESCE(sp.description, ''), sp.amount, sp.due_date,
		       sp.status, sp.paid_at, COALESCE(sp.payment_method, ''), COALESCE(sp.transaction_id, ''),
		       sp.payer_user_id, COALESCE(sp.receipt_number, ''), sp.created_at, sp.updated_at,
		       COALESCE(u.first_name || ' ' || u.last_name, '') AS student_name
		FROM school_payments sp
		LEFT JOIN users u ON sp.student_id = u.id
		WHERE (sp.school_id = $1::uuid OR $1 = '')
		  AND (
			sp.student_id IN (
				SELECT s.user_id FROM student_parents rel JOIN parents p ON rel.parent_id = p.id JOIN students s ON rel.student_id = s.id WHERE p.user_id = $2::uuid OR p.id = $2::uuid
				UNION
				SELECT s.id FROM student_parents rel JOIN parents p ON rel.parent_id = p.id JOIN students s ON rel.student_id = s.id WHERE p.user_id = $2::uuid OR p.id = $2::uuid
				UNION
				SELECT s.user_id FROM parent_students ps LEFT JOIN parents p ON ps.parent_id = p.id LEFT JOIN students s ON (ps.student_id = s.id OR ps.student_id = s.user_id) WHERE ps.parent_id = $2::uuid OR p.user_id = $2::uuid OR p.id = $2::uuid
				UNION
				SELECT s.id FROM parent_students ps LEFT JOIN parents p ON ps.parent_id = p.id LEFT JOIN students s ON (ps.student_id = s.id OR ps.student_id = s.user_id) WHERE ps.parent_id = $2::uuid OR p.user_id = $2::uuid OR p.id = $2::uuid
				UNION
				SELECT psg.student_id FROM parent_student_guardians psg WHERE psg.parent_id::text = $2::text
			)
			OR sp.payer_user_id = $2::uuid
		  )
		ORDER BY sp.due_date DESC
	`
	rows, err := r.db.QueryContext(ctx, query, schoolID, parentUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*SchoolPayment
	for rows.Next() {
		p := &SchoolPayment{}
		var payerID sql.NullString
		var paidAt sql.NullTime
		if err := rows.Scan(
			&p.ID, &p.SchoolID, &p.StudentID, &p.Title, &p.Description, &p.Amount, &p.DueDate,
			&p.Status, &paidAt, &p.PaymentMethod, &p.TransactionID,
			&payerID, &p.ReceiptNumber, &p.CreatedAt, &p.UpdatedAt,
			&p.StudentName,
		); err != nil {
			return nil, err
		}
		if paidAt.Valid {
			p.PaidAt = &paidAt.Time
		}
		if payerID.Valid {
			p.PayerUserID = &payerID.String
		}
		list = append(list, p)
	}
	if list == nil {
		list = []*SchoolPayment{}
	}
	return list, rows.Err()
}

func (r *PostgresRepository) ListBySchool(ctx context.Context, schoolID string) ([]*SchoolPayment, error) {
	query := `
		SELECT sp.id, sp.school_id, sp.student_id, sp.title, COALESCE(sp.description, ''), sp.amount, sp.due_date,
		       sp.status, sp.paid_at, COALESCE(sp.payment_method, ''), COALESCE(sp.transaction_id, ''),
		       sp.payer_user_id, COALESCE(sp.receipt_number, ''), sp.created_at, sp.updated_at,
		       COALESCE(u.first_name || ' ' || u.last_name, '') AS student_name
		FROM school_payments sp
		LEFT JOIN users u ON sp.student_id = u.id
		WHERE sp.school_id = $1::uuid OR $1 = ''
		ORDER BY sp.due_date DESC
	`
	rows, err := r.db.QueryContext(ctx, query, schoolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*SchoolPayment
	for rows.Next() {
		p := &SchoolPayment{}
		var payerID sql.NullString
		var paidAt sql.NullTime
		if err := rows.Scan(
			&p.ID, &p.SchoolID, &p.StudentID, &p.Title, &p.Description, &p.Amount, &p.DueDate,
			&p.Status, &paidAt, &p.PaymentMethod, &p.TransactionID,
			&payerID, &p.ReceiptNumber, &p.CreatedAt, &p.UpdatedAt,
			&p.StudentName,
		); err != nil {
			return nil, err
		}
		if paidAt.Valid {
			p.PaidAt = &paidAt.Time
		}
		if payerID.Valid {
			p.PayerUserID = &payerID.String
		}
		list = append(list, p)
	}
	if list == nil {
		list = []*SchoolPayment{}
	}
	return list, rows.Err()
}

func (r *PostgresRepository) Pay(ctx context.Context, id, payerUserID, method, transactionID, receiptNumber string) error {
	query := `
		UPDATE school_payments
		SET status = 'paid', paid_at = NOW(), payment_method = $1, transaction_id = $2, payer_user_id = $3::uuid, receipt_number = $4, updated_at = NOW()
		WHERE id = $5::uuid AND status = 'pending'
	`
	res, err := r.db.ExecContext(ctx, query, method, transactionID, payerUserID, receiptNumber, id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("pagamento già effettuato o non trovato")
	}
	return nil
}

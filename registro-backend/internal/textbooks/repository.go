package textbooks

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, t *Textbook) error
	GetByID(ctx context.Context, id string) (*Textbook, error)
	Update(ctx context.Context, t *Textbook) error
	List(ctx context.Context, schoolID string) ([]Textbook, error)
	Delete(ctx context.Context, id string) error

	AssignToClass(ctx context.Context, classID string, subjectID string, textbookID string, optional bool) error
	RemoveFromClass(ctx context.Context, assignmentID string) error
	ListByClass(ctx context.Context, classID string) ([]ClassTextbook, error)

	// AIE & Spending Limits
	UpsertAIECatalog(ctx context.Context, books []AIECatalogBook) (int, error)
	SearchAIECatalog(ctx context.Context, query, subject, schoolOrder string, limit int) ([]AIECatalogBook, error)
	GetSpendingLimit(ctx context.Context, schoolID string, classYear int, schoolOrder, academicYear string) (*SpendingLimit, error)
	UpsertSpendingLimit(ctx context.Context, limit *SpendingLimit) error
	ListClassAdoptions(ctx context.Context, classID string) ([]ClassAdoptionItem, error)
	SaveClassAdoption(ctx context.Context, item *ClassAdoptionItem) error
	DeleteClassAdoption(ctx context.Context, id string) error
	GetClassInfo(ctx context.Context, classID string) (classYear int, schoolOrder string, academicYear string, schoolCode string, className string, err error)
}

type postgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(ctx context.Context, t *Textbook) error {
	t.ID = uuid.New().String()
	query := `INSERT INTO textbooks (id, school_id, title, author, subject, isbn, publisher, price) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.ExecContext(ctx, query, t.ID, t.SchoolID, t.Title, t.Author, t.Subject, t.ISBN, t.Publisher, t.Price)
	return err
}

func (r *postgresRepository) GetByID(ctx context.Context, id string) (*Textbook, error) {
	query := `SELECT id, school_id, title, COALESCE(author, ''), COALESCE(subject, ''), COALESCE(isbn, ''), COALESCE(publisher, ''), COALESCE(price, 0), created_at FROM textbooks WHERE id = $1`
	var t Textbook
	err := r.db.QueryRowContext(ctx, query, id).Scan(&t.ID, &t.SchoolID, &t.Title, &t.Author, &t.Subject, &t.ISBN, &t.Publisher, &t.Price, &t.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *postgresRepository) Update(ctx context.Context, t *Textbook) error {
	query := `UPDATE textbooks SET title = $1, author = $2, subject = $3, isbn = $4, publisher = $5, price = $6 WHERE id = $7`
	_, err := r.db.ExecContext(ctx, query, t.Title, t.Author, t.Subject, t.ISBN, t.Publisher, t.Price, t.ID)
	return err
}

func (r *postgresRepository) List(ctx context.Context, schoolID string) ([]Textbook, error) {
	query := `SELECT id, school_id, title, COALESCE(author, ''), COALESCE(subject, ''), COALESCE(isbn, ''), COALESCE(publisher, ''), COALESCE(price, 0), created_at FROM textbooks WHERE school_id = $1 ORDER BY title`
	rows, err := r.db.QueryContext(ctx, query, schoolID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var res []Textbook
	for rows.Next() {
		var t Textbook
		if err := rows.Scan(&t.ID, &t.SchoolID, &t.Title, &t.Author, &t.Subject, &t.ISBN, &t.Publisher, &t.Price, &t.CreatedAt); err != nil {
			return nil, err
		}
		res = append(res, t)
	}
	return res, rows.Err()
}

func (r *postgresRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM textbooks WHERE id = $1", id)
	return err
}

func (r *postgresRepository) AssignToClass(ctx context.Context, classID string, subjectID string, textbookID string, optional bool) error {
	id := uuid.New().String()
	query := `INSERT INTO class_textbooks (id, class_id, subject_id, textbook_id, is_optional) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query, id, classID, subjectID, textbookID, optional)
	return err
}

func (r *postgresRepository) RemoveFromClass(ctx context.Context, assignmentID string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM class_textbooks WHERE id = $1", assignmentID)
	return err
}

func (r *postgresRepository) ListByClass(ctx context.Context, classID string) ([]ClassTextbook, error) {
	query := `
		SELECT ct.id, ct.class_id, ct.subject_id, s.name, ct.textbook_id, t.title, t.author, ct.is_optional
		FROM class_textbooks ct
		JOIN textbooks t ON ct.textbook_id = t.id
		JOIN subjects s ON ct.subject_id = s.id
		WHERE ct.class_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, classID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var res []ClassTextbook
	for rows.Next() {
		var ct ClassTextbook
		if err := rows.Scan(&ct.ID, &ct.ClassID, &ct.SubjectID, &ct.SubjectName, &ct.TextbookID, &ct.Title, &ct.Author, &ct.IsOptional); err != nil {
			return nil, err
		}
		res = append(res, ct)
	}
	return res, rows.Err()
}

// AIE Catalog & Spending Limit implementations

func (r *postgresRepository) UpsertAIECatalog(ctx context.Context, books []AIECatalogBook) (int, error) {
	if len(books) == 0 {
		return 0, nil
	}
	count := 0
	query := `
		INSERT INTO aie_catalog (isbn, title, authors, publisher, subject, price, volume, edition_year, school_order, is_digital_only)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (isbn) DO UPDATE SET
			title = EXCLUDED.title,
			authors = EXCLUDED.authors,
			publisher = EXCLUDED.publisher,
			subject = EXCLUDED.subject,
			price = EXCLUDED.price,
			volume = EXCLUDED.volume,
			edition_year = EXCLUDED.edition_year,
			school_order = EXCLUDED.school_order,
			is_digital_only = EXCLUDED.is_digital_only
	`
	for _, b := range books {
		_, err := r.db.ExecContext(ctx, query, b.ISBN, b.Title, b.Authors, b.Publisher, b.Subject, b.Price, b.Volume, b.EditionYear, b.SchoolOrder, b.IsDigitalOnly)
		if err == nil {
			count++
		}
	}
	return count, nil
}

func (r *postgresRepository) SearchAIECatalog(ctx context.Context, queryStr, subject, schoolOrder string, limit int) ([]AIECatalogBook, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	var conditions []string
	var args []interface{}
	argIdx := 1

	if queryStr != "" {
		conditions = append(conditions, fmt.Sprintf("(isbn ILIKE $%d OR title ILIKE $%d OR authors ILIKE $%d OR publisher ILIKE $%d)", argIdx, argIdx, argIdx, argIdx))
		args = append(args, "%"+queryStr+"%")
		argIdx++
	}
	if subject != "" {
		conditions = append(conditions, fmt.Sprintf("subject ILIKE $%d", argIdx))
		args = append(args, "%"+subject+"%")
		argIdx++
	}
	if schoolOrder != "" {
		conditions = append(conditions, fmt.Sprintf("school_order = $%d", argIdx))
		args = append(args, schoolOrder)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	q := fmt.Sprintf(`
		SELECT id, isbn, title, authors, publisher, subject, price, volume, COALESCE(edition_year, 0), school_order, is_digital_only, created_at
		FROM aie_catalog
		%s
		ORDER BY title ASC
		LIMIT %d
	`, whereClause, limit)

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var results []AIECatalogBook
	for rows.Next() {
		var b AIECatalogBook
		if err := rows.Scan(&b.ID, &b.ISBN, &b.Title, &b.Authors, &b.Publisher, &b.Subject, &b.Price, &b.Volume, &b.EditionYear, &b.SchoolOrder, &b.IsDigitalOnly, &b.CreatedAt); err != nil {
			return nil, err
		}
		results = append(results, b)
	}
	return results, rows.Err()
}

func (r *postgresRepository) GetSpendingLimit(ctx context.Context, schoolID string, classYear int, schoolOrder, academicYear string) (*SpendingLimit, error) {
	query := `
		SELECT id, school_id, class_year, school_order, max_amount, allowed_tolerance_pct, academic_year, created_at
		FROM textbook_spending_limits
		WHERE school_id = $1 AND class_year = $2 AND school_order = $3 AND academic_year = $4
	`
	var sl SpendingLimit
	err := r.db.QueryRowContext(ctx, query, schoolID, classYear, schoolOrder, academicYear).Scan(
		&sl.ID, &sl.SchoolID, &sl.ClassYear, &sl.SchoolOrder, &sl.MaxAmount, &sl.AllowedTolerancePct, &sl.AcademicYear, &sl.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &sl, nil
}

func (r *postgresRepository) UpsertSpendingLimit(ctx context.Context, limit *SpendingLimit) error {
	query := `
		INSERT INTO textbook_spending_limits (school_id, class_year, school_order, max_amount, allowed_tolerance_pct, academic_year)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (school_id, class_year, school_order, academic_year) DO UPDATE SET
			max_amount = EXCLUDED.max_amount,
			allowed_tolerance_pct = EXCLUDED.allowed_tolerance_pct
	`
	_, err := r.db.ExecContext(ctx, query, limit.SchoolID, limit.ClassYear, limit.SchoolOrder, limit.MaxAmount, limit.AllowedTolerancePct, limit.AcademicYear)
	return err
}

func (r *postgresRepository) ListClassAdoptions(ctx context.Context, classID string) ([]ClassAdoptionItem, error) {
	query := `
		SELECT 
			cta.id, 
			cta.class_id, 
			cta.subject_id, 
			COALESCE(s.name, ''),
			COALESCE(cta.book_id::text, cta.textbook_id::text, ''),
			COALESCE(a.title, t.title, ''),
			COALESCE(a.authors, t.author, ''),
			COALESCE(a.publisher, t.publisher, ''),
			COALESCE(a.isbn, t.isbn, ''),
			COALESCE(a.price, t.price, 0),
			COALESCE(cta.adoption_type, 'nuova_adozione'),
			COALESCE(cta.is_already_owned, false),
			COALESCE(cta.is_monographic, false),
			cta.deliberated_at,
			COALESCE(cta.notes, '')
		FROM class_textbook_adoptions cta
		LEFT JOIN aie_catalog a ON cta.book_id = a.id
		LEFT JOIN textbooks t ON cta.textbook_id = t.id
		LEFT JOIN subjects s ON cta.subject_id = s.id
		WHERE cta.class_id = $1
		ORDER BY s.name ASC, a.title ASC
	`
	rows, err := r.db.QueryContext(ctx, query, classID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var items []ClassAdoptionItem
	for rows.Next() {
		var item ClassAdoptionItem
		if err := rows.Scan(
			&item.ID, &item.ClassID, &item.SubjectID, &item.SubjectName,
			&item.BookID, &item.BookTitle, &item.Authors, &item.Publisher, &item.ISBN,
			&item.Price, &item.AdoptionType, &item.IsAlreadyOwned, &item.IsMonographic,
			&item.DeliberatedAt, &item.Notes,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *postgresRepository) SaveClassAdoption(ctx context.Context, item *ClassAdoptionItem) error {
	if item.ID == "" {
		item.ID = uuid.New().String()
	}

	var isAie bool
	_ = r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM aie_catalog WHERE id = $1)", item.BookID).Scan(&isAie)

	if isAie {
		query := `
			INSERT INTO class_textbook_adoptions (id, class_id, subject_id, book_id, adoption_type, is_already_owned, is_monographic, notes)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`
		_, err := r.db.ExecContext(ctx, query, item.ID, item.ClassID, item.SubjectID, item.BookID, item.AdoptionType, item.IsAlreadyOwned, item.IsMonographic, item.Notes)
		return err
	}

	query := `
		INSERT INTO class_textbook_adoptions (id, class_id, subject_id, textbook_id, adoption_type, is_already_owned, is_monographic, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.ExecContext(ctx, query, item.ID, item.ClassID, item.SubjectID, item.BookID, item.AdoptionType, item.IsAlreadyOwned, item.IsMonographic, item.Notes)
	return err
}

func (r *postgresRepository) DeleteClassAdoption(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM class_textbook_adoptions WHERE id = $1", id)
	return err
}

func (r *postgresRepository) GetClassInfo(ctx context.Context, classID string) (int, string, string, string, string, error) {
	query := `
		SELECT 
			COALESCE(al.sorting_order, 1),
			COALESCE(sc.order_type, 'secondaria_2'),
			COALESCE(ay.name, '2026/2027'),
			COALESCE(sc.code, 'SCUOLA01'),
			COALESCE(c.section, 'A')
		FROM classes c
		LEFT JOIN academic_levels al ON c.level_id = al.id
		LEFT JOIN academic_years ay ON c.academic_year_id = ay.id
		LEFT JOIN schools sc ON c.school_id = sc.id
		WHERE c.id = $1
	`
	var classYear int
	var schoolOrder, academicYear, schoolCode, section string
	err := r.db.QueryRowContext(ctx, query, classID).Scan(&classYear, &schoolOrder, &academicYear, &schoolCode, &section)
	if err != nil {
		return 1, "secondaria_2", "2026/2027", "SCUOLA01", "1A", nil
	}
	return classYear, schoolOrder, academicYear, schoolCode, section, nil
}

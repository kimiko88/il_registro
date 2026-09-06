/**
 * @fileoverview API Type Definitions — Registro Elettronico Scolastico
 *
 * JSDoc type definitions per i principali modelli API del backend Go.
 * Questi tipi migliorano l'autocomplete negli IDE senza richiedere
 * una migrazione completa a TypeScript.
 *
 * Uso:
 *   // Nel file JS che usa il tipo:
 *   /** @type {import('@/types/api').Grade} *\/
 *   const grade = await gradeService.saveGrade(data)
 *
 * @module types/api
 */

// ─────────────────────────────────────────────────────────────────────────────
// Autenticazione
// ─────────────────────────────────────────────────────────────────────────────

/**
 * @typedef {Object} LoginRequest
 * @property {string} email
 * @property {string} password
 * @property {string} [totp_code] - Codice TOTP a 6 cifre per MFA
 */

/**
 * @typedef {Object} LoginResponse
 * @property {string} access_token
 * @property {string} [refresh_token]
 * @property {User} user
 * @property {boolean} [mfa_required]
 */

// ─────────────────────────────────────────────────────────────────────────────
// Utenti
// ─────────────────────────────────────────────────────────────────────────────

/**
 * @typedef {'superadmin'|'admin'|'secretary'|'principal'|'vice_principal'|'staff'|'coordinator'|'teacher'|'student'|'parent'} UserRole
 */

/**
 * @typedef {Object} User
 * @property {string} id - UUID
 * @property {string} email
 * @property {string} first_name
 * @property {string} last_name
 * @property {UserRole} role
 * @property {string} [school_id] - UUID istituto
 * @property {string} [locale] - Locale preferita (es. 'it-IT')
 * @property {string} [avatar_url]
 * @property {boolean} [mfa_enabled]
 * @property {string} [created_at] - ISO 8601
 * @property {string|null} [deleted_at] - ISO 8601 | null
 */

// ─────────────────────────────────────────────────────────────────────────────
// Scuole e classi
// ─────────────────────────────────────────────────────────────────────────────

/**
 * @typedef {Object} School
 * @property {string} id - UUID
 * @property {string} name
 * @property {string} [address]
 * @property {string} [city]
 * @property {string} [phone]
 * @property {string} [email_peo] - Email ordinaria
 * @property {string} [email_pec] - Email PEC
 * @property {string} [mechanical_code] - Codice meccanografico MIM
 * @property {string} [created_at]
 */

/**
 * @typedef {Object} Class
 * @property {string} id - UUID
 * @property {string} school_id - UUID
 * @property {string} name - Es. '2A', '3B Liceo Scientifico'
 * @property {string} [section]
 * @property {string} [year] - Es. '1', '2', '3', '4', '5'
 * @property {string} [address] - Indirizzo di studi
 * @property {string} [academic_year_id] - UUID
 * @property {string} [created_at]
 */

/**
 * @typedef {Object} Subject
 * @property {string} id - UUID
 * @property {string} name - Es. 'Matematica'
 * @property {string} [school_id]
 */

// ─────────────────────────────────────────────────────────────────────────────
// Voti
// ─────────────────────────────────────────────────────────────────────────────

/**
 * @typedef {'Scritto'|'Orale'|'Pratico'|'Progetto'} EvaluationType
 */

/**
 * @typedef {Object} Grade
 * @property {string} id - UUID
 * @property {string} student_id - UUID
 * @property {string} subject_id - UUID
 * @property {string} [class_id] - UUID
 * @property {number} grade_value - Tra 1.0 e 10.0
 * @property {EvaluationType} [evaluation_type]
 * @property {number} [weight] - Default 1.0
 * @property {string} [description] - Note sul voto
 * @property {number} [semester] - 1 o 2
 * @property {boolean} is_published
 * @property {string} [teacher_id] - UUID
 * @property {string} [date] - Data voto ISO 8601
 * @property {string} [created_at]
 * @property {string|null} [deleted_at]
 * @property {string} [compensative_measures] - JSON array misure compensative DSA
 */

/**
 * @typedef {Object} SaveGradeRequest
 * @property {string} student_id
 * @property {string} subject_id
 * @property {string} [class_id]
 * @property {number} grade_value
 * @property {EvaluationType} [evaluation_type]
 * @property {number} [weight]
 * @property {string} [description]
 * @property {number} [semester]
 * @property {string} [date]
 */

/**
 * @typedef {Object} ClassGradesResponse
 * @property {Array<StudentGrades>} students
 */

/**
 * @typedef {Object} StudentGrades
 * @property {string} student_id
 * @property {string} student_name
 * @property {Grade[]} grades
 */

// ─────────────────────────────────────────────────────────────────────────────
// Presenze
// ─────────────────────────────────────────────────────────────────────────────

/**
 * @typedef {'present'|'absent'|'late'|'early_exit'|'out_of_class'} AttendanceStatus
 */

/**
 * @typedef {Object} AttendanceRecord
 * @property {string} id - UUID
 * @property {string} student_id - UUID
 * @property {string} class_id - UUID
 * @property {string} date - ISO 8601 (YYYY-MM-DD)
 * @property {AttendanceStatus} status
 * @property {number} [hour] - Ora di lezione (1-8)
 * @property {boolean} [is_justified]
 * @property {string} [justification_note]
 * @property {string} [teacher_id]
 * @property {string} [created_at]
 * @property {string|null} [deleted_at]
 */

// ─────────────────────────────────────────────────────────────────────────────
// Comunicazioni / Circolari
// ─────────────────────────────────────────────────────────────────────────────

/**
 * @typedef {'circular'|'bacheca'|'direct'} MessageType
 */

/**
 * @typedef {Object} Communication
 * @property {string} id - UUID
 * @property {string} sender_id - UUID
 * @property {string[]} receiver_ids - Array UUID destinatari
 * @property {string} subject
 * @property {string} body
 * @property {MessageType} type
 * @property {boolean} requires_signature
 * @property {boolean} requires_ack
 * @property {string} [attachment_url]
 * @property {string} [school_id]
 * @property {string} [created_at]
 * @property {string|null} [deleted_at]
 */

// ─────────────────────────────────────────────────────────────────────────────
// Colloqui
// ─────────────────────────────────────────────────────────────────────────────

/**
 * @typedef {'confirmed'|'cancelled'|'completed'|'no_show'} BookingStatus
 */

/**
 * @typedef {Object} ColloquioSlot
 * @property {string} id - UUID
 * @property {string} teacher_id - UUID
 * @property {string} school_id - UUID
 * @property {string} date - YYYY-MM-DD
 * @property {string} start_time - HH:MM
 * @property {string} end_time - HH:MM
 * @property {number} max_bookings
 * @property {number} booking_count
 * @property {'individual'|'group'} type
 * @property {string} [location]
 */

/**
 * @typedef {Object} ColloquioBooking
 * @property {string} id - UUID
 * @property {string} slot_id - UUID
 * @property {string} [parent_id] - UUID
 * @property {string} [student_id] - UUID
 * @property {BookingStatus} status
 * @property {string} [notes]
 * @property {string} [created_at]
 */

// ─────────────────────────────────────────────────────────────────────────────
// Pagelle / Scrutinio
// ─────────────────────────────────────────────────────────────────────────────

/**
 * @typedef {'approved'|'approved_with_conditions'|'not_approved'|'suspended'|'pending'} ScrutinyOutcome
 */

/**
 * @typedef {Object} ScrutinyRecord
 * @property {string} id - UUID
 * @property {string} student_id - UUID
 * @property {string} class_id - UUID
 * @property {number} semester - 1 o 2
 * @property {ScrutinyOutcome} outcome
 * @property {string} [conduct_grade] - Voto di condotta
 * @property {boolean} [is_validated]
 * @property {string} [validated_by] - UUID docente coordinatore
 * @property {string} [validated_at]
 * @property {string} [created_at]
 */

// ─────────────────────────────────────────────────────────────────────────────
// Risposta API generica
// ─────────────────────────────────────────────────────────────────────────────

/**
 * @template T
 * @typedef {Object} ApiResponse
 * @property {T} data - Dati della risposta
 * @property {string} [message] - Messaggio opzionale
 * @property {string} [code] - Codice errore strutturato (es. 'AUTH_RATE_LIMIT_EXCEEDED')
 * @property {string} [error] - Descrizione errore
 */

/**
 * @typedef {Object} PaginatedResponse
 * @property {any[]} data
 * @property {number} total
 * @property {number} page
 * @property {number} per_page
 */

export {}

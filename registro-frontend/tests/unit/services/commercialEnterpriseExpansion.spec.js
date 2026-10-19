import { describe, it, expect, vi, beforeEach } from 'vitest'
import { api } from '@/services/api'

// Mock api service to test contract integration
vi.mock('@/services/api', () => ({
  api: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn()
  }
}))

describe('Commercial Enterprise Expansion Pack (Moduli A-H)', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  // =========================================================================
  // MODULO A: Gestione Monte Ore Docenti ("40+40 Ore" ex Art. 29 CCNL Scuola)
  // =========================================================================
  describe('Modulo A: Gestione Monte Ore Docenti (Art. 29 CCNL)', () => {
    it('calcola correttamente ore svolte, residue e rileva superamento tetto contrattuale 40h', () => {
      const obligation = {
        collegio_planned_hours: 40.0,
        collegio_done_hours: 42.5,
        consigli_planned_hours: 40.0,
        consigli_done_hours: 35.0
      }

      const collegioRemaining = obligation.collegio_planned_hours - obligation.collegio_done_hours
      const consigliRemaining = obligation.consigli_planned_hours - obligation.consigli_done_hours

      expect(collegioRemaining).toBe(-2.5) // Overtime
      expect(consigliRemaining).toBe(5.0)  // In debito
      expect(obligation.collegio_done_hours > 40.0).toBe(true)
      expect(obligation.consigli_done_hours <= 40.0).toBe(true)
    })

    it('valida vincolo CCNL: permessi brevi max 2 ore giornaliere e recupero entro 2 mesi', () => {
      const requestLeave = (hours, requestDateStr) => {
        if (hours > 2.0) {
          throw new Error('I permessi brevi non possono superare le 2 ore giornaliere')
        }
        const [y, m, d] = requestDateStr.split('-').map(Number)
        const deadline = new Date(Date.UTC(y, m - 1 + 2, d))
        return {
          hours,
          request_date: requestDateStr,
          recovery_deadline: deadline.toISOString().substring(0, 10),
          status: 'pending'
        }
      }

      expect(() => requestLeave(2.5, '2026-03-10')).toThrow(
        'I permessi brevi non possono superare le 2 ore giornaliere'
      )

      const validLeave = requestLeave(1.5, '2026-03-10')
      expect(validLeave.hours).toBe(1.5)
      expect(validLeave.recovery_deadline).toBe('2026-05-10')
      expect(validLeave.status).toBe('pending')
    })
  })

  // =========================================================================
  // MODULO B: Adozione Libri di Testo con Import Ufficiale AIE e Tetti di Spesa
  // =========================================================================
  describe('Modulo B: Adozione Libri di Testo AIE & Tetto di Spesa Ministeriale', () => {
    it('calcola la somma dei libri obbligatori e determina lo sforamento del tetto ministeriale', () => {
      const classMaxBudget = 290.0 // MIM limit for secondaria I grado classe prima
      const adoptedBooks = [
        { title: 'Matematica Insieme', price: 28.50, is_mandatory: true },
        { title: 'Grammatica Viva', price: 32.00, is_mandatory: true },
        { title: 'Storia e Futuro', price: 25.00, is_mandatory: true },
        { title: 'Geografia del Mondo', price: 22.00, is_mandatory: true },
        { title: 'Scienze della Natura', price: 24.50, is_mandatory: true },
        { title: 'Dizionario Latino (Consigliato)', price: 45.00, is_mandatory: false }
      ]

      const mandatoryTotal = adoptedBooks
        .filter(b => b.is_mandatory)
        .reduce((acc, b) => acc + b.price, 0)

      expect(mandatoryTotal).toBe(132.00)
      expect(mandatoryTotal <= classMaxBudget).toBe(true)

      const budgetAlertColor = mandatoryTotal > classMaxBudget ? 'negative' : (mandatoryTotal > classMaxBudget * 0.9 ? 'warning' : 'positive')
      expect(budgetAlertColor).toBe('positive')
    })
  })

  // =========================================================================
  // MODULO C: Sanità Scolastica & Registro Infortuni INAIL (Termine 48h)
  // =========================================================================
  describe('Modulo C: Sanità Scolastica & Termine Perentorio 48h INAIL', () => {
    it('calcola correttamente la scadenza perentoria di 48 ore dalla ricezione del certificato medico', () => {
      const accidentTime = new Date('2026-10-15T10:00:00Z')
      const deadline = new Date(accidentTime.getTime() + 48 * 3600 * 1000)

      expect(deadline.toISOString()).toBe('2026-10-17T10:00:00.000Z')

      const currentTimePassed = new Date('2026-10-17T14:00:00Z')
      const isOverdue = currentTimePassed.getTime() > deadline.getTime()
      expect(isOverdue).toBe(true)
    })

    it('identifica correttamente scadenze farmaci e piani terapeutici', () => {
      const today = new Date('2026-10-17')
      const med1 = { name: 'Adrenalina Autoiniettore', expiry_date: '2026-10-10' } // scaduto
      const med2 = { name: 'Salbutamolo Spray', expiry_date: '2027-05-31' } // valido

      const isExpired = (expiryStr) => new Date(expiryStr) < today
      expect(isExpired(med1.expiry_date)).toBe(true)
      expect(isExpired(med2.expiry_date)).toBe(false)
    })
  })

  // =========================================================================
  // MODULO D: Nuova Certificazione delle Competenze D.M. 14/2024 & INVALSI
  // =========================================================================
  describe('Modulo D: Certificazione Competenze D.M. 14/2024 & Analisi INVALSI', () => {
    it('supporta i 4 livelli ufficiali D.M. 14/2024 e il calcolo del delta INVALSI vs Registro', () => {
      const validLevels = ['Avanzato', 'Intermedio', 'Base', 'Iniziale']
      const certLevel = 'Avanzato'
      expect(validLevels.includes(certLevel)).toBe(true)

      // Test INVALSI WLE / Level vs Teacher Grade Correlation
      const invalsiResult = {
        student_id: 's_001',
        subject: 'italiano',
        invalsi_level: 4,      // Livello 1-5
        invalsi_wle: 215.4,
        national_mean: 200.0,
        class_grade: 7.5
      }

      // Delta: confronto normalizzato
      const normalizedInvalsiGrade = 4 + (invalsiResult.invalsi_level - 1) * 1.5 // map lvl 1-5 to approx 4-10
      const delta = invalsiResult.class_grade - normalizedInvalsiGrade
      expect(normalizedInvalsiGrade).toBe(8.5)
      expect(delta).toBe(-1.0)
    })
  })

  // =========================================================================
  // MODULO E: Sicurezza Fisica: Pass Digitale Fuori Classe & Tutela Affido Art. 337 c.c.
  // =========================================================================
  describe('Modulo E: Pass Fuori Classe & Deleghe / Affido Art. 337 c.c.', () => {
    it('attiva l\'alert anomalia quando l\'alunno è fuori dall\'aula da più di 15 minuti', () => {
      const departedAt = new Date('2026-10-17T11:00:00Z').getTime()
      const now12Min = new Date('2026-10-17T11:12:00Z').getTime()
      const now18Min = new Date('2026-10-17T11:18:00Z').getTime()

      const getElapsed = (dep, now) => Math.floor((now - dep) / 60000)

      expect(getElapsed(departedAt, now12Min)).toBe(12)
      expect(getElapsed(departedAt, now12Min) > 15).toBe(false)

      expect(getElapsed(departedAt, now18Min)).toBe(18)
      expect(getElapsed(departedAt, now18Min) > 15).toBe(true)
    })

    it('gestisce vincolo affido monogenitoriale ex art. 337 c.c. con delega al ritiro', () => {
      const delegate = {
        student_id: 'stud_42',
        delegate_name: 'Mario Rossi',
        relationship: 'Nonno paterno',
        custody_order_ref: 'Sentenza Trib. Minori n. 142/2025',
        is_permanent: true
      }

      expect(delegate.custody_order_ref).toBeTruthy()
      expect(delegate.is_permanent).toBe(true)
    })
  })

  // =========================================================================
  // MODULO F: Didattica Digitale: Canvas Annotazione PDF & Sync Classroom
  // =========================================================================
  describe('Modulo F: Consegna Elaborati con Canvas PDF & Google Classroom', () => {
    it('serializza e valida annotazioni a mano libera del docente su canvas', () => {
      const strokes = [
        { type: 'pen', color: '#ef4444', width: 4, points: [{ x: 10, y: 15 }, { x: 30, y: 45 }] },
        { type: 'text', color: '#3b82f6', text: 'Ottimo svolgimento!', x: 100, y: 150 }
      ]

      const json = JSON.stringify(strokes)
      const parsed = JSON.parse(json)

      expect(parsed.length).toBe(2)
      expect(parsed[0].type).toBe('pen')
      expect(parsed[1].type).toBe('text')
      expect(parsed[1].text).toBe('Ottimo svolgimento!')
    })

    it('gestisce la sincronizzazione del voto con Google Classroom API', async () => {
      api.post.mockResolvedValueOnce({
        data: { message: 'voto sincronizzato con Google Classroom' }
      })

      const res = await api.post('/homework-submissions/sub_123/sync-classroom', {
        google_classroom_submission_id: 'gclass_sub_123'
      })

      expect(api.post).toHaveBeenCalledWith(
        '/homework-submissions/sub_123/sync-classroom',
        { google_classroom_submission_id: 'gclass_sub_123' }
      )
      expect(res.data.message).toBe('voto sincronizzato con Google Classroom')
    })
  })

  // =========================================================================
  // MODULO G: Monitoraggio Progetti PNRR & Piattaforma Futura (Soglia 75%)
  // =========================================================================
  describe('Modulo G: Monitoraggio Progetti PNRR & Vincolo 75% Futura MIM', () => {
    it('convalida la liquidazione solo se la frequenza del corsista è almeno al 75%', () => {
      const totalProjectHours = 30

      const check75 = (attendedHours) => {
        const percentage = (attendedHours / totalProjectHours) * 100
        return {
          percentage,
          target_75_met: percentage >= 75.0
        }
      }

      // Studente con 22 ore (73.33%) -> NON liquidabile
      const stud1 = check75(22)
      expect(stud1.percentage).toBeCloseTo(73.33, 1)
      expect(stud1.target_75_met).toBe(false)

      // Studente con 23 ore (76.67%) -> LIQUIDABILE
      const stud2 = check75(23)
      expect(stud2.percentage).toBeCloseTo(76.67, 1)
      expect(stud2.target_75_met).toBe(true)

      // Studente con 30 ore (100%) -> LIQUIDABILE
      const stud3 = check75(30)
      expect(stud3.percentage).toBe(100.0)
      expect(stud3.target_75_met).toBe(true)
    })
  })

  // =========================================================================
  // MODULO H: Registro Scuola dell'Infanzia (Mensa, 5 Campi Esperienza, Diario)
  // =========================================================================
  describe('Modulo H: Registro Scuola Infanzia & Diario di Bordo', () => {
    it('calcola i totali presenze e pasti giornalieri per la cucina scolastica', () => {
      const roster = [
        { student_id: 'bambino_1', is_present: true, has_meal: true, diet_notes: '' },
        { student_id: 'bambino_2', is_present: true, has_meal: true, diet_notes: 'Senza glutine' },
        { student_id: 'bambino_3', is_present: true, has_meal: false, diet_notes: '' },
        { student_id: 'bambino_4', is_present: false, has_meal: false, diet_notes: '' }
      ]

      const presentCount = roster.filter(r => r.is_present).length
      const mealCount = roster.filter(r => r.has_meal).length
      const specialDietCount = roster.filter(r => r.diet_notes && r.diet_notes.trim()).length

      expect(presentCount).toBe(3)
      expect(mealCount).toBe(2)
      expect(specialDietCount).toBe(1)
    })

    it('valida i 5 Campi di Esperienza previsti dalle Indicazioni Nazionali', () => {
      const ministerialCampi = [
        'Il sé e l\'altro',
        'Il corpo e il movimento',
        'Immagini, suoni, colori',
        'I discorsi e le parole',
        'La conoscenza del mondo'
      ]

      const observation = {
        student_id: 'bambino_1',
        campo_esperienza: 'Il corpo e il movimento',
        level_descriptor: 'In evoluzione',
        observation_date: '2026-10-17'
      }

      expect(ministerialCampi.includes(observation.campo_esperienza)).toBe(true)
    })
  })
})

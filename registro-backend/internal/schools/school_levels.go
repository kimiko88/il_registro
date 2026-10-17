package schools

import "strings"

// Canonical School Tiers in the Italian Education System
const (
	TierInfanzia               = "infanzia"
	TierPrimaria               = "primaria"
	TierSecondariaPrimoGrado   = "secondaria_primo_grado"
	TierSecondariaSecondoGrado = "secondaria_secondo_grado"
	TierComprensivo            = "comprensivo"
	TierOmnicomprensivo        = "omnicomprensivo"
)

// TierFeatures models the regulatory characteristics and capabilities of a school tier.
type TierFeatures struct {
	Tier                  string   `json:"tier"`
	Name                  string   `json:"name"`
	NormativeReference    string   `json:"normative_reference"`
	EvaluationType        string   `json:"evaluation_type"` // 'campi_esperienza', 'obiettivi_descrittivi', 'numerica_decimale', 'misto'
	HasGrades             bool     `json:"has_grades"`
	HasCampiEsperienza    bool     `json:"has_campi_esperienza"`
	HasDailyDiary         bool     `json:"has_daily_diary"`
	HasPrimaryLevels      bool     `json:"has_primary_levels"`
	HasLearningObjectives bool     `json:"has_learning_objectives"`
	HasDeferredScrutiny   bool     `json:"has_deferred_scrutiny"`
	HasSchoolCredits      bool     `json:"has_school_credits"`
	HasPCTO               bool     `json:"has_pcto"`
	HasInvalsi            bool     `json:"has_invalsi"`
	HasOrientamento       bool     `json:"has_orientamento"`
	LevelsSupported       []string `json:"levels_supported,omitempty"`
	GradeScaleMin         *float64 `json:"grade_scale_min,omitempty"`
	GradeScaleMax         *float64 `json:"grade_scale_max,omitempty"`
	PrimaryLevels         []string `json:"primary_levels,omitempty"`
	Dimensions            []string `json:"dimensions,omitempty"`
	CampiEsperienza       []string `json:"campi_esperienza,omitempty"`
}

// NormalizeTier returns a canonical tier identifier from arbitrary input strings.
func NormalizeTier(raw string) string {
	t := strings.ToLower(strings.TrimSpace(raw))
	switch t {
	case TierInfanzia, "scuola dell'infanzia", "nido", "materna":
		return TierInfanzia
	case TierPrimaria, "scuola primaria", "elementare", "elementari":
		return TierPrimaria
	case TierSecondariaPrimoGrado, "secondaria_1_grado", "secondaria_primo", "media", "medie", "secondaria_1":
		return TierSecondariaPrimoGrado
	case TierComprensivo, "istituto comprensivo", "ic":
		return TierComprensivo
	case TierOmnicomprensivo, "istituto omnicomprensivo":
		return TierOmnicomprensivo
	case TierSecondariaSecondoGrado, "secondaria_2_grado", "secondaria_secondo", "superiore", "superiori", "liceo", "istituto tecnico", "istituto professionale":
		return TierSecondariaSecondoGrado
	default:
		if strings.Contains(t, "infanzia") {
			return TierInfanzia
		}
		if strings.Contains(t, "primaria") {
			return TierPrimaria
		}
		if strings.Contains(t, "comprensivo") && !strings.Contains(t, "omni") {
			return TierComprensivo
		}
		if strings.Contains(t, "omni") {
			return TierOmnicomprensivo
		}
		if strings.Contains(t, "primo") || strings.Contains(t, "media") {
			return TierSecondariaPrimoGrado
		}
		return TierSecondariaSecondoGrado
	}
}

// GetTierFeatures returns the feature set and normative configuration for a school tier.
func GetTierFeatures(rawTier string) TierFeatures {
	tier := NormalizeTier(rawTier)
	minGrade := 1.0
	maxGrade := 10.0

	switch tier {
	case TierInfanzia:
		return TierFeatures{
			Tier:                  TierInfanzia,
			Name:                  "Scuola dell'Infanzia",
			NormativeReference:    "D.M. 254/2012 Indicazioni Nazionali per il curricolo",
			EvaluationType:        "campi_esperienza",
			HasGrades:             false,
			HasCampiEsperienza:    true,
			HasDailyDiary:         true,
			HasPrimaryLevels:      false,
			HasLearningObjectives: false,
			HasDeferredScrutiny:   false,
			HasSchoolCredits:      false,
			HasPCTO:               false,
			HasInvalsi:            false,
			HasOrientamento:       false,
			CampiEsperienza: []string{
				"Il sé e l'altro",
				"Il corpo e il movimento",
				"Immagini, suoni, colori",
				"I discorsi e le parole",
				"La conoscenza del mondo",
			},
		}

	case TierPrimaria:
		return TierFeatures{
			Tier:                  TierPrimaria,
			Name:                  "Scuola Primaria",
			NormativeReference:    "O.M. 172/2020 e D.Lgs. 62/2017 (Giudizi descrittivi per obiettivi)",
			EvaluationType:        "obiettivi_descrittivi",
			HasGrades:             false,
			HasCampiEsperienza:    false,
			HasDailyDiary:         false,
			HasPrimaryLevels:      true,
			HasLearningObjectives: true,
			HasDeferredScrutiny:   false,
			HasSchoolCredits:      false,
			HasPCTO:               false,
			HasInvalsi:            true,
			HasOrientamento:       false,
			PrimaryLevels: []string{
				"Avanzato",
				"Intermedio",
				"Base",
				"In via di prima acquisizione",
			},
			Dimensions: []string{
				"Continuità",
				"Tipologia della situazione (nota/non nota)",
				"Risorse utilizzate",
				"Autonomia",
			},
		}

	case TierSecondariaPrimoGrado:
		return TierFeatures{
			Tier:                  TierSecondariaPrimoGrado,
			Name:                  "Scuola Secondaria di I Grado",
			NormativeReference:    "D.Lgs. 62/2017 e D.M. 14/2024 (Modello Nazionale Certificazione Competenze)",
			EvaluationType:        "numerica_decimale",
			HasGrades:             true,
			HasCampiEsperienza:    false,
			HasDailyDiary:         false,
			HasPrimaryLevels:      false,
			HasLearningObjectives: false,
			HasDeferredScrutiny:   false, // Nella scuola media le carenze non prevedono esami di riparazione a settembre
			HasSchoolCredits:      false,
			HasPCTO:               false,
			HasInvalsi:            true,
			HasOrientamento:       true,
			GradeScaleMin:         &minGrade,
			GradeScaleMax:         &maxGrade,
		}

	case TierComprensivo:
		return TierFeatures{
			Tier:                  TierComprensivo,
			Name:                  "Istituto Comprensivo",
			NormativeReference:    "L. 111/2011 (Infanzia, Primaria e Secondaria di I Grado integrati)",
			EvaluationType:        "misto",
			HasGrades:             true,
			HasCampiEsperienza:    true,
			HasDailyDiary:         true,
			HasPrimaryLevels:      true,
			HasLearningObjectives: true,
			HasDeferredScrutiny:   false,
			HasSchoolCredits:      false,
			HasPCTO:               false,
			HasInvalsi:            true,
			HasOrientamento:       true,
			LevelsSupported: []string{
				TierInfanzia,
				TierPrimaria,
				TierSecondariaPrimoGrado,
			},
			GradeScaleMin: &minGrade,
			GradeScaleMax: &maxGrade,
			PrimaryLevels: []string{
				"Avanzato",
				"Intermedio",
				"Base",
				"In via di prima acquisizione",
			},
			CampiEsperienza: []string{
				"Il sé e l'altro",
				"Il corpo e il movimento",
				"Immagini, suoni, colori",
				"I discorsi e le parole",
				"La conoscenza del mondo",
			},
		}

	case TierOmnicomprensivo:
		return TierFeatures{
			Tier:                  TierOmnicomprensivo,
			Name:                  "Istituto Omnicomprensivo",
			NormativeReference:    "D.Lgs. 297/1994 (Aggregazione di tutti gli ordini scolastici)",
			EvaluationType:        "misto",
			HasGrades:             true,
			HasCampiEsperienza:    true,
			HasDailyDiary:         true,
			HasPrimaryLevels:      true,
			HasLearningObjectives: true,
			HasDeferredScrutiny:   true,
			HasSchoolCredits:      true,
			HasPCTO:               true,
			HasInvalsi:            true,
			HasOrientamento:       true,
			LevelsSupported: []string{
				TierInfanzia,
				TierPrimaria,
				TierSecondariaPrimoGrado,
				TierSecondariaSecondoGrado,
			},
			GradeScaleMin: &minGrade,
			GradeScaleMax: &maxGrade,
			PrimaryLevels: []string{
				"Avanzato",
				"Intermedio",
				"Base",
				"In via di prima acquisizione",
			},
		}

	default: // TierSecondariaSecondoGrado
		return TierFeatures{
			Tier:                  TierSecondariaSecondoGrado,
			Name:                  "Scuola Secondaria di II Grado",
			NormativeReference:    "D.P.R. 89/2010, O.M. 92/2007 (Debiti formativi e Scrutinio Differito) e D.Lgs. 62/2017",
			EvaluationType:        "numerica_decimale",
			HasGrades:             true,
			HasCampiEsperienza:    false,
			HasDailyDiary:         false,
			HasPrimaryLevels:      false,
			HasLearningObjectives: false,
			HasDeferredScrutiny:   true, // O.M. 92/2007: Prove di recupero a fine agosto / inizio settembre e scioglimento riserva
			HasSchoolCredits:      true, // D.Lgs. 62/2017: Credito scolastico triennio (max 40 punti)
			HasPCTO:               true,
			HasInvalsi:            true,
			HasOrientamento:       true,
			GradeScaleMin:         &minGrade,
			GradeScaleMax:         &maxGrade,
		}
	}
}

// GetAllTiers returns all supported canonical tiers and their normative feature configurations.
func GetAllTiers() []TierFeatures {
	return []TierFeatures{
		GetTierFeatures(TierInfanzia),
		GetTierFeatures(TierPrimaria),
		GetTierFeatures(TierSecondariaPrimoGrado),
		GetTierFeatures(TierSecondariaSecondoGrado),
		GetTierFeatures(TierComprensivo),
		GetTierFeatures(TierOmnicomprensivo),
	}
}

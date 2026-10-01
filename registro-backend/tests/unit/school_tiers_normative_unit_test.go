package unit

import (
	"encoding/json"
	"testing"

	"registro-backend/internal/admin"
	"registro-backend/internal/schools"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSchoolTiers_Normalization(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"infanzia", schools.TierInfanzia},
		{"Scuola dell'Infanzia", schools.TierInfanzia},
		{"materna", schools.TierInfanzia},
		{"primaria", schools.TierPrimaria},
		{"Scuola Primaria", schools.TierPrimaria},
		{"elementare", schools.TierPrimaria},
		{"secondaria_primo_grado", schools.TierSecondariaPrimoGrado},
		{"secondaria_1_grado", schools.TierSecondariaPrimoGrado},
		{"media", schools.TierSecondariaPrimoGrado},
		{"scuola media statale", schools.TierSecondariaPrimoGrado},
		{"secondaria_secondo_grado", schools.TierSecondariaSecondoGrado},
		{"Liceo Scientifico", schools.TierSecondariaSecondoGrado},
		{"Istituto Tecnico Industriale", schools.TierSecondariaSecondoGrado},
		{"comprensivo", schools.TierComprensivo},
		{"Istituto Comprensivo G. Rodari", schools.TierComprensivo},
		{"omnicomprensivo", schools.TierOmnicomprensivo},
		{"", schools.TierSecondariaSecondoGrado},
	}

	for _, tc := range tests {
		actual := schools.NormalizeTier(tc.input)
		assert.Equal(t, tc.expected, actual, "Failed normalizing %s", tc.input)
	}
}

func TestSchoolTiers_NormativeFeatures_Infanzia(t *testing.T) {
	features := schools.GetTierFeatures(schools.TierInfanzia)
	assert.Equal(t, schools.TierInfanzia, features.Tier)
	assert.False(t, features.HasGrades, "Infanzia must not use numerical grades")
	assert.True(t, features.HasCampiEsperienza, "Infanzia must use Campi d'Esperienza")
	assert.True(t, features.HasDailyDiary, "Infanzia must support pedagogical logbook")
	assert.False(t, features.HasDeferredScrutiny, "Infanzia does not have deferred scrutiny")
	assert.False(t, features.HasSchoolCredits, "Infanzia does not have school credits")
	assert.Contains(t, features.NormativeReference, "D.M. 254/2012")
	assert.Equal(t, 5, len(features.CampiEsperienza))
	assert.Contains(t, features.CampiEsperienza, "Il sé e l'altro")
	assert.Contains(t, features.CampiEsperienza, "La conoscenza del mondo")
}

func TestSchoolTiers_NormativeFeatures_Primaria(t *testing.T) {
	features := schools.GetTierFeatures(schools.TierPrimaria)
	assert.Equal(t, schools.TierPrimaria, features.Tier)
	assert.False(t, features.HasGrades, "Primaria must not use numerical grades since O.M. 172/2020")
	assert.True(t, features.HasPrimaryLevels, "Primaria must use descriptive levels")
	assert.True(t, features.HasLearningObjectives, "Primaria evaluates by learning objectives")
	assert.False(t, features.HasDeferredScrutiny, "Primaria does not have summer debt recovery exams")
	assert.Contains(t, features.NormativeReference, "O.M. 172/2020")
	assert.Equal(t, 4, len(features.PrimaryLevels))
	assert.Contains(t, features.PrimaryLevels, "Avanzato")
	assert.Contains(t, features.PrimaryLevels, "Intermedio")
	assert.Contains(t, features.PrimaryLevels, "Base")
	assert.Contains(t, features.PrimaryLevels, "In via di prima acquisizione")
	assert.Equal(t, 4, len(features.Dimensions))
	assert.Contains(t, features.Dimensions, "Autonomia")
	assert.Contains(t, features.Dimensions, "Continuità")
}

func TestSchoolTiers_NormativeFeatures_SecondariaPrimoGrado(t *testing.T) {
	features := schools.GetTierFeatures(schools.TierSecondariaPrimoGrado)
	assert.Equal(t, schools.TierSecondariaPrimoGrado, features.Tier)
	assert.True(t, features.HasGrades, "Secondaria I Grado uses numerical 1-10 grades")
	assert.True(t, features.HasInvalsi, "Grade 8 includes INVALSI")
	assert.True(t, features.HasOrientamento, "Grade 8 includes Consiglio Orientativo")
	assert.False(t, features.HasDeferredScrutiny, "Middle school does not have debt recovery exams in September")
	assert.False(t, features.HasSchoolCredits, "Middle school does not have triennium credits")
}

func TestSchoolTiers_NormativeFeatures_SecondariaSecondoGrado(t *testing.T) {
	features := schools.GetTierFeatures(schools.TierSecondariaSecondoGrado)
	assert.Equal(t, schools.TierSecondariaSecondoGrado, features.Tier)
	assert.True(t, features.HasGrades)
	assert.True(t, features.HasDeferredScrutiny, "High school must have deferred scrutiny per O.M. 92/2007")
	assert.True(t, features.HasSchoolCredits, "High school has triennium credits per D.Lgs. 62/2017")
	assert.True(t, features.HasPCTO, "High school triennium includes PCTO")
	assert.Contains(t, features.NormativeReference, "O.M. 92/2007")
}

func TestSchoolTiers_GetAllTiers(t *testing.T) {
	tiers := schools.GetAllTiers()
	assert.Equal(t, 6, len(tiers))
	tierMap := make(map[string]schools.TierFeatures)
	for _, tier := range tiers {
		tierMap[tier.Tier] = tier
	}
	assert.Contains(t, tierMap, schools.TierInfanzia)
	assert.Contains(t, tierMap, schools.TierPrimaria)
	assert.Contains(t, tierMap, schools.TierSecondariaPrimoGrado)
	assert.Contains(t, tierMap, schools.TierSecondariaSecondoGrado)
	assert.Contains(t, tierMap, schools.TierComprensivo)
	assert.Contains(t, tierMap, schools.TierOmnicomprensivo)
}

func TestSchoolTiers_DTOSerialization(t *testing.T) {
	// 1. schools.CreateSchoolRequest
	req1 := schools.CreateSchoolRequest{
		Name:        "Scuola Primaria Gianni Rodari",
		Code:        "RMEE81000X",
		SchoolLevel: schools.TierPrimaria,
		Address:     "Via delle Scuole 12",
	}
	data1, err := json.Marshal(req1)
	require.NoError(t, err)
	assert.Contains(t, string(data1), `"school_level":"primaria"`)

	// 2. admin.CreateSchoolRequest
	req2 := admin.CreateSchoolRequest{
		Name:        "Liceo Statale Morgagni",
		Code:        "RMPS02000X",
		SchoolLevel: schools.TierSecondariaSecondoGrado,
		Address:     "Via Fonteiana 85",
		City:        "Roma",
		Province:    "RM",
		ZipCode:     "00152",
	}
	data2, err := json.Marshal(req2)
	require.NoError(t, err)
	assert.Contains(t, string(data2), `"school_level":"secondaria_secondo_grado"`)

	// 3. admin.SchoolResponse
	resp := admin.SchoolResponse{
		ID:          "s-123",
		Name:        "IC Dante Alighieri",
		Code:        "RMIC800001",
		SchoolLevel: schools.TierComprensivo,
		Type:        schools.TierComprensivo,
	}
	data3, err := json.Marshal(resp)
	require.NoError(t, err)
	assert.Contains(t, string(data3), `"school_level":"comprensivo"`)
	assert.Contains(t, string(data3), `"type":"comprensivo"`)
}

package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAll25RolesValidation tests that all 25 official Italian school governance roles are valid.
func TestAll25RolesValidation(t *testing.T) {
	expectedRoles := []string{
		// 1. Governance & Direzione
		RolePrincipal,
		RoleVicePrincipal,
		RoleCollaboratoreDS,
		// 2. Didattica & Coordinamento
		RoleTeacher,
		RoleCoordinator,
		// 3. Studenti & Famiglie
		RoleStudent,
		RoleParent,
		// 4. Amministrazione & ATA
		RoleDSGA,
		RoleSecretary,
		RoleAssistenteAmministrativo,
		RoleAssistenteAlunni,
		RoleAssistentePersonale,
		RoleAssistenteContabilita,
		RoleAssistenteProtocollo,
		RoleAssistenteSportello,
		RoleCollaboratoreScolastico,
		RoleAssistenteTecnico,
		RoleResponsabileServizio,
		// 5. Garanzia & Tecnici
		RoleSuperAdmin,
		RoleAdmin,
		RoleResponsabileGestioneDocumentale,
		RoleResponsabileConservazione,
		RoleDPO,
		RoleSystemAuditor,
		// Incarichi aggiuntivi come ruoli formali
		RoleCoordinatoreClasse,
		RoleSegretarioConsiglio,
		RoleReferenteProgetto,
		RoleReferenteInclusione,
		RoleResponsabileDipartimento,
		RoleTutorOrientatore,
		RoleAnimatoreDigitale,
	}

	for _, role := range expectedRoles {
		assert.Truef(t, IsValidRole(role), "Expected role '%s' to be valid in IsValidRole", role)
	}

	invalidRoles := []string{"", "invalid_role", "hacker", "super_teacher", "preside_ombra"}
	for _, r := range invalidRoles {
		assert.Falsef(t, IsValidRole(r), "Expected invalid role '%s' to return false", r)
	}
}

// TestIsStaffRole checks classification between school staff vs students/parents.
func TestIsStaffRole(t *testing.T) {
	staffRoles := []string{
		RoleSuperAdmin, RoleAdmin, RolePrincipal, RoleVicePrincipal, RoleSecretary, RoleDSGA,
		RoleTeacher, RoleCoordinator, RoleAssistenteAmministrativo, RoleCollaboratoreScolastico,
		RoleCollaboratoreDS, RoleAssistenteAlunni, RoleAssistentePersonale, RoleAssistenteContabilita,
		RoleAssistenteProtocollo, RoleAssistenteSportello, RoleAssistenteTecnico, RoleResponsabileServizio,
		RoleResponsabileGestioneDocumentale, RoleResponsabileConservazione, RoleDPO, RoleSystemAuditor,
	}

	for _, r := range staffRoles {
		assert.Truef(t, IsStaffRole(r), "Role '%s' should be classified as staff", r)
	}

	nonStaffRoles := []string{RoleStudent, RoleParent, "unknown"}
	for _, r := range nonStaffRoles {
		assert.Falsef(t, IsStaffRole(r), "Role '%s' should not be classified as staff", r)
	}
}

// TestCanAssignDutyGovernanceMatrix tests the Italian school governance authority matrix.
func TestCanAssignDutyGovernanceMatrix(t *testing.T) {
	educationalDuties := []string{
		"coordinatore_classe",
		"segretario_consiglio",
		"referente_inclusione",
		"referente_progetto",
		"responsabile_dipartimento",
		"tutor_orientatore",
		"animatore_digitale",
	}

	ataDuties := []string{
		"responsabile_servizio",
		"addetto_sicurezza",
	}

	// 1. Dirigente Scolastico / Vice Preside: ha piena autorità didattica e istituzionale
	for _, duty := range educationalDuties {
		assert.True(t, CanAssignDuty(RolePrincipal, duty), "Principal must be able to assign %s", duty)
		assert.True(t, CanAssignDuty(RoleVicePrincipal, duty), "Vice Principal must be able to assign %s", duty)
	}

	// 2. Admin e SuperAdmin: delega tecnica di sistema per tutti gli incarichi
	for _, duty := range append(educationalDuties, ataDuties...) {
		assert.True(t, CanAssignDuty(RoleAdmin, duty), "Admin must be able to assign %s", duty)
		assert.True(t, CanAssignDuty(RoleSuperAdmin, duty), "SuperAdmin must be able to assign %s", duty)
	}

	// 3. DSGA: ha autorità sui compiti ATA e servizi, ma NON sulla didattica/coordinamento
	for _, ataDuty := range ataDuties {
		assert.True(t, CanAssignDuty(RoleDSGA, ataDuty), "DSGA must be able to assign ATA duty %s", ataDuty)
	}
	assert.False(t, CanAssignDuty(RoleDSGA, "coordinatore_classe"), "DSGA must NOT be able to assign coordinatore_classe")
	assert.False(t, CanAssignDuty(RoleDSGA, "referente_inclusione"), "DSGA must NOT be able to assign referente_inclusione")
	assert.False(t, CanAssignDuty(RoleDSGA, "animatore_digitale"), "DSGA must NOT be able to assign animatore_digitale")

	// 4. Segreteria: verbalizzazione/registrazione di coordinatori e segretari di classe
	assert.True(t, CanAssignDuty(RoleSecretary, "coordinatore_classe"), "Secretary can record coordinatore_classe")
	assert.True(t, CanAssignDuty(RoleSecretary, "segretario_consiglio"), "Secretary can record segretario_consiglio")
	assert.False(t, CanAssignDuty(RoleSecretary, "referente_inclusione"), "Secretary cannot assign referente_inclusione")
	assert.False(t, CanAssignDuty(RoleSecretary, "responsabile_servizio"), "Secretary cannot assign responsabile_servizio")

	// 5. Docenti, Studenti, Genitori, DPO, Auditor: nessuna facoltà di assegnazione (403)
	forbiddenRoles := []string{
		RoleTeacher,
		RoleCoordinator,
		RoleStudent,
		RoleParent,
		RoleDPO,
		RoleSystemAuditor,
		RoleCollaboratoreScolastico,
		RoleAssistenteTecnico,
	}

	for _, actor := range forbiddenRoles {
		for _, duty := range append(educationalDuties, ataDuties...) {
			assert.Falsef(t, CanAssignDuty(actor, duty), "Role '%s' must not be allowed to assign duty '%s'", actor, duty)
		}
	}
}

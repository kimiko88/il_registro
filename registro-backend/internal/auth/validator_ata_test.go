package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateRegisterRequest_ATARolesRules(t *testing.T) {
	baseReq := func(role string) *RegisterRequest {
		return &RegisterRequest{
			Email:     "utente.test@scuola.it",
			Password:  "Sicura1234!",
			FirstName: "Mario",
			LastName:  "Rossi",
			Role:      role,
		}
	}

	// 1. Segreteria PUÒ creare: assistente_amministrativo, collaboratore_scolastico, collaboratore_ds
	assert.NoError(t, ValidateRegisterRequest(RoleSecretary, baseReq(RoleAssistenteAmministrativo)))
	assert.NoError(t, ValidateRegisterRequest(RoleSecretary, baseReq(RoleCollaboratoreScolastico)))
	assert.NoError(t, ValidateRegisterRequest(RoleSecretary, baseReq(RoleCollaboratoreDS)))
	assert.NoError(t, ValidateRegisterRequest(RoleSecretary, baseReq(RoleTeacher)))
	assert.NoError(t, ValidateRegisterRequest(RoleSecretary, baseReq(RoleStudent)))

	// 2. Segreteria NON PUÒ creare DSGA (solo Admin/SuperAdmin)
	errDSGA := ValidateRegisterRequest(RoleSecretary, baseReq(RoleDSGA))
	assert.ErrorIs(t, errDSGA, ErrCannotCreateRole, "La segreteria non deve poter creare l'utente DSGA")

	// 3. Segreteria NON PUÒ creare Admin o Principal
	assert.ErrorIs(t, ValidateRegisterRequest(RoleSecretary, baseReq(RoleAdmin)), ErrCannotCreateRole)
	assert.ErrorIs(t, ValidateRegisterRequest(RoleSecretary, baseReq(RolePrincipal)), ErrCannotCreateRole)

	// 4. Admin PUÒ creare l'utente DSGA e tutti i ruoli ATA
	assert.NoError(t, ValidateRegisterRequest(RoleAdmin, baseReq(RoleDSGA)))
	assert.NoError(t, ValidateRegisterRequest(RoleAdmin, baseReq(RoleAssistenteAmministrativo)))
	assert.NoError(t, ValidateRegisterRequest(RoleAdmin, baseReq(RoleCollaboratoreDS)))
	assert.NoError(t, ValidateRegisterRequest(RoleAdmin, baseReq(RoleCollaboratoreScolastico)))

	// 5. SuperAdmin PUÒ creare DSGA e tutti i ruoli ATA
	assert.NoError(t, ValidateRegisterRequest(RoleSuperAdmin, baseReq(RoleDSGA)))

	// 6. Docente o Studente non possono creare alcun utente
	assert.ErrorIs(t, ValidateRegisterRequest(RoleTeacher, baseReq(RoleStudent)), ErrInsufficientRole)
	assert.ErrorIs(t, ValidateRegisterRequest(RoleStudent, baseReq(RoleStudent)), ErrInsufficientRole)
}

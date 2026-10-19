package auth

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"registro-backend/pkg/jwt"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func setupTestSPIDRouter(svc *Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	api := r.Group("/api/v1")

	spidH := NewSPIDHandler(svc)
	spidH.RegisterSPIDRoutes(api)

	cieH := NewCIEHandler(svc)
	cieH.RegisterCIERoutes(api)

	return r
}

func createTestSPIDService(mockRepo *MockRepository) *Service {
	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	tm := jwt.NewTokenManager(priv, &priv.PublicKey)
	svc := NewService(mockRepo, tm, nil, nil)
	svc.bcryptCost = bcrypt.MinCost
	return svc
}

func TestSPID_GetProviders(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := createTestSPIDService(mockRepo)
	router := setupTestSPIDRouter(svc)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/auth/spid/providers", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var body map[string][]SPIDProvider
	err := json.Unmarshal(w.Body.Bytes(), &body)
	assert.NoError(t, err)
	assert.NotEmpty(t, body["providers"])
	assert.True(t, len(body["providers"]) >= 9)

	// Check PosteID and InfoCert are present
	foundPoste := false
	for _, p := range body["providers"] {
		if p.ID == "poste" {
			foundPoste = true
			assert.Contains(t, p.SSOEndpoint, "posteid")
		}
	}
	assert.True(t, foundPoste)
}

func TestSPID_InitiateLogin(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := createTestSPIDService(mockRepo)
	router := setupTestSPIDRouter(svc)

	// 1. Valid IdP
	payload := `{"idp_id":"poste"}`
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/spid/login", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp["request_id"])
	assert.NotEmpty(t, resp["authn_request"])
	assert.NotEmpty(t, resp["redirect_url"])

	// Check Base64 XML contents
	b64Req := resp["authn_request"].(string)
	decoded, err := base64.StdEncoding.DecodeString(b64Req)
	assert.NoError(t, err)
	assert.Contains(t, string(decoded), "https://www.spid.gov.it/SpidL2")
	assert.Contains(t, string(decoded), "AuthnRequest")

	// 2. Unknown IdP
	invalidPayload := `{"idp_id":"unknown_idp"}`
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/auth/spid/login", bytes.NewBufferString(invalidPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSPID_Callback_ExistingUser(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := createTestSPIDService(mockRepo)
	router := setupTestSPIDRouter(svc)

	fc := "RSSMRA85M01H501Z"
	user := &User{
		ID:         "user-spid-1",
		Email:      "mario.rossi@spid.it",
		FirstName:  "Mario",
		LastName:   "Rossi",
		Role:       "parent",
		IsActive:   true,
		FiscalCode: &fc,
	}

	mockRepo.On("GetUserByFiscalCode", mock.Anything, fc).Return(user, nil).Once()
	mockRepo.On("CreateRefreshToken", mock.Anything, mock.AnythingOfType("*auth.RefreshToken")).Return(nil).Once()
	mockRepo.On("UpdateLastLogin", mock.Anything, "user-spid-1").Return(nil).Once()

	callbackPayload := `{"fiscal_code":"TINIT-RSSMRA85M01H501Z","first_name":"Mario","last_name":"Rossi"}`
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/spid/callback", bytes.NewBufferString(callbackPayload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var authResp AuthResponse
	err := json.Unmarshal(w.Body.Bytes(), &authResp)
	assert.NoError(t, err)
	assert.Equal(t, "user-spid-1", authResp.User.ID)
	assert.NotEmpty(t, authResp.AccessToken)
	assert.NotEmpty(t, authResp.RefreshToken)

	mockRepo.AssertExpectations(t)
}

func TestSPID_Callback_SAMLAssertion(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := createTestSPIDService(mockRepo)
	router := setupTestSPIDRouter(svc)

	// Auto-provision a new user
	fc := "VRDLGI90A01H501Z"
	mockRepo.On("GetUserByFiscalCode", mock.Anything, fc).Return(nil, ErrUserNotFound).Once()
	mockRepo.On("GetUserByEmail", mock.Anything, "luigi.verdi@posta.it").Return(nil, ErrUserNotFound).Once()
	mockRepo.On("CreateUser", mock.Anything, mock.MatchedBy(func(u *User) bool {
		return u.Role == "parent" && *u.FiscalCode == fc
	})).Return(nil).Once()
	mockRepo.On("CreateRefreshToken", mock.Anything, mock.AnythingOfType("*auth.RefreshToken")).Return(nil).Once()
	mockRepo.On("UpdateLastLogin", mock.Anything, mock.Anything).Return(nil).Once()

	// Simulate SAML XML Assertion with AgID attributes
	samlXML := fmt.Sprintf(`<?xml version="1.0"?>
<saml2p:Response xmlns:saml2p="urn:oasis:names:tc:SAML:2.0:protocol" xmlns:saml="urn:oasis:names:tc:SAML:2.0:assertion">
    <saml:Assertion>
        <saml:AttributeStatement>
            <saml:Attribute Name="fiscalNumber">
                <saml:AttributeValue>TINIT-%s</saml:AttributeValue>
            </saml:Attribute>
            <saml:Attribute Name="name">
                <saml:AttributeValue>Luigi</saml:AttributeValue>
            </saml:Attribute>
            <saml:Attribute Name="familyName">
                <saml:AttributeValue>Verdi</saml:AttributeValue>
            </saml:Attribute>
            <saml:Attribute Name="email">
                <saml:AttributeValue>luigi.verdi@posta.it</saml:AttributeValue>
            </saml:Attribute>
        </saml:AttributeStatement>
    </saml:Assertion>
</saml2p:Response>`, fc)

	b64SAML := base64.StdEncoding.EncodeToString([]byte(samlXML))
	callbackPayload := fmt.Sprintf(`{"saml_response":"%s"}`, b64SAML)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/spid/callback", bytes.NewBufferString(callbackPayload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var authResp AuthResponse
	err := json.Unmarshal(w.Body.Bytes(), &authResp)
	assert.NoError(t, err)
	assert.Equal(t, "Luigi", authResp.User.FirstName)
	assert.Equal(t, "Verdi", authResp.User.LastName)
	assert.NotEmpty(t, authResp.AccessToken)

	mockRepo.AssertExpectations(t)
}

func TestCIE_InfoAndLogin(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := createTestSPIDService(mockRepo)
	router := setupTestSPIDRouter(svc)

	// 1. GET /cie/info
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/auth/cie/info", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Carta d'Identità Elettronica")

	// 2. POST /cie/login
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/auth/cie/login", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "authn_request")

	// 3. POST /cie/callback
	fc := "BNCGNN80A01H501Z"
	user := &User{
		ID:         "user-cie-1",
		Email:      "giovanni.bianchi@cie.it",
		FirstName:  "Giovanni",
		LastName:   "Bianchi",
		Role:       "parent",
		IsActive:   true,
		FiscalCode: &fc,
	}

	mockRepo.On("GetUserByFiscalCode", mock.Anything, fc).Return(user, nil).Once()
	mockRepo.On("CreateRefreshToken", mock.Anything, mock.AnythingOfType("*auth.RefreshToken")).Return(nil).Once()
	mockRepo.On("UpdateLastLogin", mock.Anything, "user-cie-1").Return(nil).Once()

	callbackPayload := fmt.Sprintf(`{"fiscal_code":"%s","first_name":"Giovanni","last_name":"Bianchi"}`, fc)
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/auth/cie/callback", bytes.NewBufferString(callbackPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var authResp AuthResponse
	err := json.Unmarshal(w.Body.Bytes(), &authResp)
	assert.NoError(t, err)
	assert.Equal(t, "user-cie-1", authResp.User.ID)

	mockRepo.AssertExpectations(t)
}

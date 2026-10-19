package auth

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CIEHandler handles CIE (Carta d'Identità Elettronica - Ministero dell'Interno) authentication.
type CIEHandler struct {
	service *Service
}

// NewCIEHandler creates a new CIE handler
func NewCIEHandler(service *Service) *CIEHandler {
	return &CIEHandler{
		service: service,
	}
}

// GetCIEInfo returns metadata and endpoint configuration for CIE.
// GET /auth/cie/info
func (h *CIEHandler) GetCIEInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"name":         "Entra con CIE (Carta d'Identità Elettronica)",
		"entity_id":    "https://idserver.servizicie.interno.gov.it/idp/profile/SAML2/POST/SSO",
		"sso_endpoint": "https://idserver.servizicie.interno.gov.it/idp/profile/SAML2/POST/SSO",
		"logo_url":     "https://www.cartaidentita.interno.gov.it/assets/images/logo_cie_header.svg",
		"authn_level":  "https://www.spid.gov.it/SpidL2",
	})
}

// InitiateCIELogin initiates CIE authentication by producing a SAML 2.0 AuthnRequest.
// POST /auth/cie/login
func (h *CIEHandler) InitiateCIELogin(c *gin.Context) {
	requestID := "_" + uuid.New().String()
	issueInstant := time.Now().UTC().Format(time.RFC3339)
	cieSSOEndpoint := "https://idserver.servizicie.interno.gov.it/idp/profile/SAML2/POST/SSO"
	entityID := "https://registro.scuola.edu.it"
	acsURL := "https://registro.scuola.edu.it/api/v1/auth/cie/callback"

	authnRequestXML := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<samlp:AuthnRequest xmlns:samlp="urn:oasis:names:tc:SAML:2.0:protocol"
                    xmlns:saml="urn:oasis:names:tc:SAML:2.0:assertion"
                    ID="%s"
                    Version="2.0"
                    IssueInstant="%s"
                    Destination="%s"
                    ForceAuthn="true"
                    AssertionConsumerServiceURL="%s">
    <saml:Issuer NameQualifier="%s" Format="urn:oasis:names:tc:SAML:2.0:nameid-format:entity">%s</saml:Issuer>
    <samlp:NameIDPolicy Format="urn:oasis:names:tc:SAML:2.0:nameid-format:transient" AllowCreate="true"/>
    <samlp:RequestedAuthnContext Comparison="exact">
        <saml:AuthnContextClassRef>https://www.spid.gov.it/SpidL2</saml:AuthnContextClassRef>
    </samlp:RequestedAuthnContext>
</samlp:AuthnRequest>`,
		requestID, issueInstant, cieSSOEndpoint, acsURL, entityID, entityID,
	)

	b64XML := base64.StdEncoding.EncodeToString([]byte(authnRequestXML))

	c.JSON(http.StatusOK, gin.H{
		"request_id":    requestID,
		"authn_request": b64XML,
		"redirect_url":  cieSSOEndpoint,
		"relay_state":   uuid.New().String(),
		"cie_level":     "https://www.spid.gov.it/SpidL2",
		"binding":       "HTTP-POST",
	})
}

// HandleCIECallback processes CIE callback with SAMLResponse or direct testing payload.
// POST & GET /auth/cie/callback
func (h *CIEHandler) HandleCIECallback(c *gin.Context) {
	var payload SPIDCallbackPayload
	_ = c.ShouldBind(&payload)

	fiscalCode := payload.FiscalCode
	email := payload.Email
	firstName := payload.FirstName
	lastName := payload.LastName

	if payload.SAMLResponse != "" {
		decodedBytes, err := base64.StdEncoding.DecodeString(payload.SAMLResponse)
		if err == nil {
			xmlStr := string(decodedBytes)
			extractedFC := extractXMLAttribute(xmlStr, "fiscalNumber")
			if extractedFC != "" {
				fiscalCode = extractedFC
			}
			extractedEmail := extractXMLAttribute(xmlStr, "email")
			if extractedEmail != "" && email == "" {
				email = extractedEmail
			}
			extractedName := extractXMLAttribute(xmlStr, "name")
			if extractedName != "" && firstName == "" {
				firstName = extractedName
			}
			extractedFamily := extractXMLAttribute(xmlStr, "familyName")
			if extractedFamily != "" && lastName == "" {
				lastName = extractedFamily
			}
		}
	}

	fiscalCode = CleanSpidFiscalCode(fiscalCode)
	if fiscalCode == "" && email == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_assertion",
			Message: "Risposta CIE priva di codice fiscale o email validi",
		})
		return
	}

	ipAddress := c.ClientIP()
	userAgent := c.Request.UserAgent()

	authResp, err := h.service.AuthenticateFederatedUser(
		c.Request.Context(),
		fiscalCode,
		email,
		firstName,
		lastName,
		"CIE",
		ipAddress,
		userAgent,
	)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "cie_auth_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, authResp)
}

// RegisterCIERoutes registers CIE routes
func (h *CIEHandler) RegisterCIERoutes(router *gin.RouterGroup) {
	cie := router.Group("/auth/cie")
	{
		cie.GET("/info", h.GetCIEInfo)
		cie.POST("/login", h.InitiateCIELogin)
		cie.POST("/callback", h.HandleCIECallback)
		cie.GET("/callback", h.HandleCIECallback)
	}
}

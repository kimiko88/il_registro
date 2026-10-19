package auth

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SPIDProvider represents an AgID-accredited Identity Provider.
type SPIDProvider struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	EntityID    string `json:"entity_id"`
	SSOEndpoint string `json:"sso_endpoint"`
	LogoURL     string `json:"logo_url"`
}

// Certified AgID SPID Identity Providers
var AccreditedSPIDProviders = []SPIDProvider{
	{
		ID:          "poste",
		Name:        "PosteID",
		EntityID:    "https://posteid.poste.it",
		SSOEndpoint: "https://posteid.poste.it/jod-fs/ssoservicepost",
		LogoURL:     "https://assets.spid.gov.it/idp-logos/spid-idp-posteid.svg",
	},
	{
		ID:          "infocert",
		Name:        "InfoCert ID",
		EntityID:    "https://identity.infocert.it",
		SSOEndpoint: "https://identity.infocert.it/spid/samlsso",
		LogoURL:     "https://assets.spid.gov.it/idp-logos/spid-idp-infocertid.svg",
	},
	{
		ID:          "aruba",
		Name:        "Aruba ID",
		EntityID:    "https://loginspid.aruba.it",
		SSOEndpoint: "https://loginspid.aruba.it/ServiceLoginWelcome",
		LogoURL:     "https://assets.spid.gov.it/idp-logos/spid-idp-arubaid.svg",
	},
	{
		ID:          "namirial",
		Name:        "Namirial ID",
		EntityID:    "https://idp.namirialtsp.com/idp",
		SSOEndpoint: "https://idp.namirialtsp.com/idp/profile/SAML2/POST/SSO",
		LogoURL:     "https://assets.spid.gov.it/idp-logos/spid-idp-namirialid.svg",
	},
	{
		ID:          "tim",
		Name:        "TIM ID",
		EntityID:    "https://login.id.tim.it/affwebservices/public/saml2sso",
		SSOEndpoint: "https://login.id.tim.it/affwebservices/public/saml2sso",
		LogoURL:     "https://assets.spid.gov.it/idp-logos/spid-idp-timid.svg",
	},
	{
		ID:          "sielte",
		Name:        "Sielte ID",
		EntityID:    "https://identity.sieltecloud.it",
		SSOEndpoint: "https://identity.sieltecloud.it/simplesaml/saml2/idp/SSOService.php",
		LogoURL:     "https://assets.spid.gov.it/idp-logos/spid-idp-sielteid.svg",
	},
	{
		ID:          "lepida",
		Name:        "Lepida ID",
		EntityID:    "https://id.lepida.it/idp/shibboleth",
		SSOEndpoint: "https://id.lepida.it/idp/shibboleth",
		LogoURL:     "https://assets.spid.gov.it/idp-logos/spid-idp-lepidaid.svg",
	},
	{
		ID:          "intesi",
		Name:        "Intesi Group",
		EntityID:    "https://spid.intesagroup.com",
		SSOEndpoint: "https://spid.intesagroup.com/idp/profile/SAML2/POST/SSO",
		LogoURL:     "https://assets.spid.gov.it/idp-logos/spid-idp-intesigroupspid.svg",
	},
	{
		ID:          "spiditalia",
		Name:        "SPIDItalia",
		EntityID:    "https://spid.register.it",
		SSOEndpoint: "https://spid.register.it/login/sso",
		LogoURL:     "https://assets.spid.gov.it/idp-logos/spid-idp-spiditalia.svg",
	},
	{
		ID:          "demo",
		Name:        "SPID Testbed AgID",
		EntityID:    "https://demo.spid.gov.it",
		SSOEndpoint: "https://demo.spid.gov.it/samlsso",
		LogoURL:     "https://assets.spid.gov.it/idp-logos/spid-idp-validator.svg",
	},
}

// SPIDHandler handles SPID authentication
type SPIDHandler struct {
	service *Service
}

// NewSPIDHandler creates a new SPID handler
func NewSPIDHandler(service *Service) *SPIDHandler {
	return &SPIDHandler{
		service: service,
	}
}

// GetProviders returns list of all accredited SPID IdPs.
// GET /auth/spid/providers
func (h *SPIDHandler) GetProviders(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"providers": AccreditedSPIDProviders,
	})
}

type InitiateSPIDRequest struct {
	IdPID string `json:"idp_id" binding:"required"`
}

// InitiateSPIDLogin generates a SAML 2.0 AuthnRequest for AgID SPID Level 2.
// POST /auth/spid/login
func (h *SPIDHandler) InitiateSPIDLogin(c *gin.Context) {
	var req InitiateSPIDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "idp_id obbligatorio"})
		return
	}

	var selectedProvider *SPIDProvider
	for _, p := range AccreditedSPIDProviders {
		if p.ID == req.IdPID {
			selectedProvider = &p
			break
		}
	}

	if selectedProvider == nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: fmt.Sprintf("provider SPID '%s' non riconosciuto", req.IdPID)})
		return
	}

	requestID := "_" + uuid.New().String()
	issueInstant := time.Now().UTC().Format(time.RFC3339)
	entityID := "https://registro.scuola.edu.it"
	acsURL := "https://registro.scuola.edu.it/api/v1/auth/spid/callback"

	// AgID SPID SAML 2.0 AuthnRequest
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
		requestID, issueInstant, selectedProvider.SSOEndpoint, acsURL, entityID, entityID,
	)

	b64XML := base64.StdEncoding.EncodeToString([]byte(authnRequestXML))

	c.JSON(http.StatusOK, gin.H{
		"request_id":    requestID,
		"provider":      selectedProvider,
		"authn_request": b64XML,
		"redirect_url":  selectedProvider.SSOEndpoint,
		"relay_state":   uuid.New().String(),
		"spid_level":    "https://www.spid.gov.it/SpidL2",
		"binding":       "HTTP-POST",
	})
}

type SPIDCallbackPayload struct {
	SAMLResponse string `json:"saml_response" form:"SAMLResponse"`
	RelayState   string `json:"relay_state" form:"RelayState"`
	FiscalCode   string `json:"fiscal_code"`
	Email        string `json:"email"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
}

// CleanSpidFiscalCode strips the TINIT- prefix specified by AgID
func CleanSpidFiscalCode(raw string) string {
	raw = strings.TrimSpace(raw)
	raw = strings.TrimPrefix(raw, "TINIT-")
	raw = strings.TrimPrefix(raw, "TINIT:")
	raw = strings.TrimPrefix(raw, "tinit-")
	return strings.ToUpper(raw)
}

// HandleSPIDCallback processes SAML Response from the IdP.
// POST & GET /auth/spid/callback
func (h *SPIDHandler) HandleSPIDCallback(c *gin.Context) {
	var payload SPIDCallbackPayload
	_ = c.ShouldBind(&payload)

	fiscalCode := payload.FiscalCode
	email := payload.Email
	firstName := payload.FirstName
	lastName := payload.LastName

	// If SAMLResponse is passed, extract attributes from SAML Assertion
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
			Message: "Risposta SPID priva di codice fiscale o email validi",
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
		"SPID",
		ipAddress,
		userAgent,
	)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "spid_auth_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, authResp)
}

// extractXMLAttribute extracts the text inside <saml:AttributeValue> for an attribute Name.
func extractXMLAttribute(xmlContent, attrName string) string {
	type AttrVal struct {
		Name  string `xml:"Name,attr"`
		Value string `xml:"AttributeValue"`
	}
	type ResponseWrapper struct {
		Attributes []AttrVal `xml:"Assertion>AttributeStatement>Attribute"`
	}

	var parsed ResponseWrapper
	if err := xml.Unmarshal([]byte(xmlContent), &parsed); err == nil {
		for _, a := range parsed.Attributes {
			if a.Name == attrName {
				return strings.TrimSpace(a.Value)
			}
		}
	}

	// Fallback simple string search for robust edge tolerance
	marker := fmt.Sprintf(`Name="%s"`, attrName)
	idx := strings.Index(xmlContent, marker)
	if idx != -1 {
		rest := xmlContent[idx:]
		valStart := strings.Index(rest, "<saml:AttributeValue>")
		valEnd := strings.Index(rest, "</saml:AttributeValue>")
		if valStart != -1 && valEnd != -1 && valEnd > valStart {
			return strings.TrimSpace(rest[valStart+len("<saml:AttributeValue>") : valEnd])
		}
	}
	return ""
}

// RegisterSPIDRoutes registers all SPID endpoints.
func (h *SPIDHandler) RegisterSPIDRoutes(router *gin.RouterGroup) {
	spid := router.Group("/auth/spid")
	{
		spid.GET("/providers", h.GetProviders)
		spid.POST("/login", h.InitiateSPIDLogin)
		spid.POST("/callback", h.HandleSPIDCallback)
		spid.GET("/callback", h.HandleSPIDCallback)
	}
}

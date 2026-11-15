package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CIE authentication is intentionally fail-closed until a complete SAML
// implementation validates signed assertions and request correlation.
type CIEHandler struct{}

func NewCIEHandler(_ *Service) *CIEHandler { return &CIEHandler{} }

func (h *CIEHandler) GetCIEInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"name":         "Entra con CIE (Carta d'Identità Elettronica)",
		"entity_id":    "https://idserver.servizicie.interno.gov.it/idp/profile/SAML2/POST/SSO",
		"sso_endpoint": "https://idserver.servizicie.interno.gov.it/idp/profile/SAML2/POST/SSO",
		"logo_url":     "https://www.cartaidentita.interno.gov.it/assets/images/logo_cie_header.svg",
		"enabled":      false,
	})
}

func (h *CIEHandler) InitiateCIELogin(c *gin.Context) {
	c.JSON(http.StatusServiceUnavailable, ErrorResponse{
		Error:   "cie_temporarily_unavailable",
		Message: "Autenticazione CIE disabilitata: il flusso SAML non è ancora sottoposto a verifica crittografica completa.",
	})
}

func (h *CIEHandler) HandleCIECallback(c *gin.Context) {
	// Direct identity fields and unsigned SAML payloads must never create sessions.
	c.JSON(http.StatusServiceUnavailable, ErrorResponse{
		Error: "cie_temporarily_unavailable",
		Message: "Callback CIE disabilitato finché la risposta SAML non può essere verificata integralmente.",
	})
}

func (h *CIEHandler) RegisterCIERoutes(router *gin.RouterGroup) {
	cie := router.Group("/auth/cie")
	cie.GET("/info", h.GetCIEInfo)
	cie.POST("/login", h.InitiateCIELogin)
	cie.POST("/callback", h.HandleCIECallback)
}

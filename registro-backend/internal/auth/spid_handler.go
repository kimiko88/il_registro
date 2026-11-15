package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SPID providers are listed for UI discovery only. The SAML flow remains
// disabled until signed assertions, issuer/audience, destination, request
// correlation, replay protection and certified metadata are fully verified.
type SPIDProvider struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	EntityID    string `json:"entity_id"`
	SSOEndpoint string `json:"sso_endpoint"`
	LogoURL     string `json:"logo_url"`
}

var AccreditedSPIDProviders = []SPIDProvider{
	{ID: "poste", Name: "PosteID", EntityID: "https://posteid.poste.it", SSOEndpoint: "https://posteid.poste.it/jod-fs/ssoservicepost", LogoURL: "https://assets.spid.gov.it/idp-logos/spid-idp-posteid.svg"},
	{ID: "infocert", Name: "InfoCert ID", EntityID: "https://identity.infocert.it", SSOEndpoint: "https://identity.infocert.it/spid/samlsso", LogoURL: "https://assets.spid.gov.it/idp-logos/spid-idp-infocertid.svg"},
	{ID: "aruba", Name: "Aruba ID", EntityID: "https://loginspid.aruba.it", SSOEndpoint: "https://loginspid.aruba.it/ServiceLoginWelcome", LogoURL: "https://assets.spid.gov.it/idp-logos/spid-idp-arubaid.svg"},
	{ID: "namirial", Name: "Namirial ID", EntityID: "https://idp.namirialtsp.com/idp", SSOEndpoint: "https://idp.namirialtsp.com/idp/profile/SAML2/POST/SSO", LogoURL: "https://assets.spid.gov.it/idp-logos/spid-idp-namirialid.svg"},
	{ID: "tim", Name: "TIM ID", EntityID: "https://login.id.tim.it/affwebservices/public/saml2sso", SSOEndpoint: "https://login.id.tim.it/affwebservices/public/saml2sso", LogoURL: "https://assets.spid.gov.it/idp-logos/spid-idp-timid.svg"},
	{ID: "sielte", Name: "Sielte ID", EntityID: "https://identity.sieltecloud.it", SSOEndpoint: "https://identity.sieltecloud.it/simplesaml/saml2/idp/SSOService.php", LogoURL: "https://assets.spid.gov.it/idp-logos/spid-idp-sielteid.svg"},
	{ID: "lepida", Name: "Lepida ID", EntityID: "https://id.lepida.it/idp/shibboleth", SSOEndpoint: "https://id.lepida.it/idp/shibboleth", LogoURL: "https://assets.spid.gov.it/idp-logos/spid-idp-lepidaid.svg"},
	{ID: "intesi", Name: "Intesi Group", EntityID: "https://spid.intesagroup.com", SSOEndpoint: "https://spid.intesagroup.com/idp/profile/SAML2/POST/SSO", LogoURL: "https://assets.spid.gov.it/idp-logos/spid-idp-intesigroupspid.svg"},
	{ID: "spiditalia", Name: "SPIDItalia", EntityID: "https://spid.register.it", SSOEndpoint: "https://spid.register.it/login/sso", LogoURL: "https://assets.spid.gov.it/idp-logos/spid-idp-spiditalia.svg"},
	{ID: "demo", Name: "SPID Testbed AgID", EntityID: "https://demo.spid.gov.it", SSOEndpoint: "https://demo.spid.gov.it/samlsso", LogoURL: "https://assets.spid.gov.it/idp-logos/spid-idp-validator.svg"},
}

type SPIDHandler struct{}

func NewSPIDHandler(_ *Service) *SPIDHandler { return &SPIDHandler{} }

func (h *SPIDHandler) GetProviders(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"providers": AccreditedSPIDProviders, "enabled": false})
}

func (h *SPIDHandler) InitiateSPIDLogin(c *gin.Context) {
	c.JSON(http.StatusServiceUnavailable, ErrorResponse{
		Error:   "spid_temporarily_unavailable",
		Message: "Autenticazione SPID disabilitata: il flusso SAML non è ancora sottoposto a verifica crittografica completa.",
	})
}

func (h *SPIDHandler) HandleSPIDCallback(c *gin.Context) {
	// Never trust email/fiscal_code from a request body or unsigned XML.
	c.JSON(http.StatusServiceUnavailable, ErrorResponse{
		Error:   "spid_temporarily_unavailable",
		Message: "Callback SPID disabilitato finché la risposta SAML non può essere verificata integralmente.",
	})
}

func (h *SPIDHandler) RegisterSPIDRoutes(router *gin.RouterGroup) {
	spid := router.Group("/auth/spid")
	spid.GET("/providers", h.GetProviders)
	spid.POST("/login", h.InitiateSPIDLogin)
	spid.POST("/callback", h.HandleSPIDCallback)
	// GET callbacks are intentionally not registered: SAML responses use POST.
}

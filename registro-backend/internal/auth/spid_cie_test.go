package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func setupFederatedAuthRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	api := r.Group("/api/v1")
	NewSPIDHandler(nil).RegisterSPIDRoutes(api)
	NewCIEHandler(nil).RegisterCIERoutes(api)
	return r
}

func TestFederatedAuthIsFailClosed(t *testing.T) {
	router := setupFederatedAuthRouter()
	cases := []struct{ method, path string }{
		{http.MethodPost, "/api/v1/auth/spid/login"},
		{http.MethodPost, "/api/v1/auth/spid/callback"},
		{http.MethodPost, "/api/v1/auth/cie/login"},
		{http.MethodPost, "/api/v1/auth/cie/callback"},
	}
	for _, tc := range cases {
		t.Run(tc.path, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, strings.NewReader(`{"email":"victim@example.test","fiscal_code":"RSSMRA85M01H501Z"}`))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)
			require.Equal(t, http.StatusServiceUnavailable, rec.Code)
			require.NotContains(t, rec.Body.String(), "access_token")
		})
	}
}

func TestFederatedProviderMetadataDoesNotEnableAuthentication(t *testing.T) {
	router := setupFederatedAuthRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/spid/providers", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), `"enabled":false`)
	require.Contains(t, rec.Body.String(), "PosteID")
}

func TestCIEInfoReportsDisabled(t *testing.T) {
	router := setupFederatedAuthRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/cie/info", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), `"enabled":false`)
}

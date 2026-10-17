package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"registro-backend/pkg/logger"
)

type CSPReportWrapper struct {
	CSPReport struct {
		DocumentURI       string `json:"document-uri"`
		Referrer          string `json:"referrer"`
		BlockedURI        string `json:"blocked-uri"`
		ViolatedDirective string `json:"violated-directive"`
		OriginalPolicy    string `json:"original-policy"`
		Disposition       string `json:"disposition"`
		StatusCode        int    `json:"status-code"`
	} `json:"csp-report"`
}

// HandleCSPReport ingests browser Content Security Policy violation reports.
func HandleCSPReport(c *gin.Context) {
	body, err := io.ReadAll(io.LimitReader(c.Request.Body, 64*1024)) // limit to 64KB
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	var report CSPReportWrapper
	if err := json.Unmarshal(body, &report); err == nil && report.CSPReport.ViolatedDirective != "" {
		logger.Log.Warnf("[CSP Violation] document=%s blocked=%s directive=%s client_ip=%s",
			report.CSPReport.DocumentURI,
			report.CSPReport.BlockedURI,
			report.CSPReport.ViolatedDirective,
			c.ClientIP(),
		)
	} else if len(body) > 0 {
		logger.Log.Warnf("[CSP Violation] client_ip=%s raw=%s", c.ClientIP(), string(body))
	}

	c.Status(http.StatusNoContent)
}

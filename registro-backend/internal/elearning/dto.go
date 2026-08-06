package elearning

type ProviderStatus struct {
	Configured  bool   `json:"configured"`
	Connected   bool   `json:"connected"`
	ClientID    string `json:"client_id,omitempty"`
	RedirectURI string `json:"redirect_uri,omitempty"`
	TenantID    string `json:"tenant_id,omitempty"`
}

type ProvidersResponse struct {
	Google    ProviderStatus `json:"google"`
	Microsoft ProviderStatus `json:"microsoft"`
}

type ConnectRequest struct {
	Code string `json:"code" binding:"required"`
}

type SyncAssignmentsRequest struct {
	ClassID string `json:"class_id" binding:"required"`
}

type SyncGradesRequest struct {
	ClassID string `json:"class_id" binding:"required"`
}

type SyncResponse struct {
	Message string `json:"message"`
	Count   int    `json:"synced_count"`
}

// NOTE: GetAuditLogs moved from handler direct-repo call to service layer.
// This file replaces the direct h.service.repo.GetAuditLogs call in handler.go.
// The actual service.go for users must expose GetAuditLogs on its interface.
// Below is a patch comment only — see users/handler.go for the updated call.
package users

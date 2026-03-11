package ticket

import (
	"net/http"
	"strings"

	platformauth "support-go/backend/internal/platform/auth"
)

type Role string

const (
	RoleClient Role = "client"
	RoleAgent  Role = "agent"
	RoleAdmin  Role = "admin"
)

func roleFromRequest(request *http.Request) Role {
	claims, ok := platformauth.ClaimsFromContext(request.Context())
	if !ok {
		return ""
	}

	return Role(strings.ToLower(strings.TrimSpace(claims.Role)))
}

func canManageAssignmentsAndStatus(role Role) bool {
	return role == RoleAgent || role == RoleAdmin
}

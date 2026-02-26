package ticket

import (
	"net/http"
	"strings"
)

type Role string

const (
	RoleClient Role = "client"
	RoleAgent  Role = "agent"
	RoleAdmin  Role = "admin"
)

func roleFromRequest(request *http.Request) Role {
	rawRole := request.Header.Get("X-User-Role")
	if rawRole == "" {
		rawRole = request.Header.Get("X-Role")
	}

	return Role(strings.ToLower(strings.TrimSpace(rawRole)))
}

func canManageAssignmentsAndStatus(role Role) bool {
	return role == RoleAgent || role == RoleAdmin
}

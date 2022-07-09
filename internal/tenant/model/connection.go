package model

// NewConnection represents a new tenant connection.
type NewConnection struct {
	UserID   string `json:"userId"`
	TenantID string `json:"tenantId"`
}

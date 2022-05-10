package model

type Tenant struct {
	ID   string `json:"tenantId"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

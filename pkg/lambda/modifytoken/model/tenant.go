package model

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"regexp"
	"strings"
)

type Tenant struct {
	TenantID    string `json:"id"`
	CompanyName string `json:"companyName"`
	Plan        string `json:"plan"`
	Path        string `json:"path"`
}

// TenantConnectionMap represents a valid tenant connection mapping.
type TenantConnectionMap map[string]Tenant

func NewTenants(items []map[string]types.AttributeValue) ([]Tenant, error) {
	var tenants = make([]Tenant, 0, len(items))
	for _, v := range items {
		var item Tenant
		if err := attributevalue.UnmarshalMap(v, &item); err != nil {
			return nil, err
		}
		tenants = append(tenants, item)
	}
	return tenants, nil
}

func NewTenantConnectionMap(tenants []Tenant) TenantConnectionMap {
	var m = make(TenantConnectionMap, len(tenants))
	for _, tenant := range tenants {
		tenantPath := formatPath(tenant.CompanyName)
		m[tenantPath] = Tenant{
			TenantID:    tenant.TenantID,
			CompanyName: tenant.CompanyName,
			Plan:        tenant.Plan,
			Path:        tenantPath,
		}
	}
	return m
}

func formatPath(company string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		panic("regex failed to compile" + err.Error())
	}
	return strings.ToLower(reg.ReplaceAllString(company, ""))
}

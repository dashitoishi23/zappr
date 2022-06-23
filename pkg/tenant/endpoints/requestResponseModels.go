package tenantendpoint

import (
	commonmodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/models"
	tenantmodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/tenant/models"
)

type CreateTenantRequest struct {
	NewTenant tenantmodels.Tenant `json:"newTenant"`
}

type CreateTenantResponse struct {
	NewTenant tenantmodels.Tenant `json:"newTenant"`
	Err       error               `json:"-"`
}

func (c CreateTenantResponse) Failed() error { return c.Err }

type FindFirstTenantRequest struct {
	CurrentTenant tenantmodels.SearchableTenant `json:"currentTenant"`
}

type FindFirstTenantResponse struct {
	CurrentTenant tenantmodels.Tenant `json:"currentTenant"`
	Err           error               `json:"-"`
}

func (f FindFirstTenantResponse) Failed() error { return f.Err }

type GetAllTenantsRequest struct {
}

type GetAllTenantsResponse struct {
	Tenants []tenantmodels.Tenant `json:"tenants"`
	Err     error                 `json:"-"`
}

func (g GetAllTenantsResponse) Failed() error { return g.Err }

type FindTenantsRequest struct {
	CurrentTenant tenantmodels.SearchableTenant `json:"currentTenant"`
}

type FindTenantsResponse struct {
	CurrentTenant []tenantmodels.Tenant `json:"currentTenant"`
	Err           error                 `json:"-"`
}

func (f FindTenantsResponse) Failed() error { return f.Err }

type UpdateTenantRequest struct {
	NewTenant tenantmodels.Tenant `json:"newTenant"`
}

type UpdateTenantResponse struct {
	UpdatedTenant tenantmodels.Tenant `json:"updatedTenant"`
	Err error `json:"-"`
}

func (u UpdateTenantResponse) Failed() error { return u.Err }

type GetPagedTenantResponse struct {
	PagedTenants commonmodels.PagedResponse[tenantmodels.Tenant] `json:"pagedTenants"`
	Err error `json:"-"`
}

func (g GetPagedTenantResponse) Failed() error { return g.Err }



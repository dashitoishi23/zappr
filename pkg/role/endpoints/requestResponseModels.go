package masterroleendpoint

import (
	commonmodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/models"
	masterrolemodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/role/models"
)

type CreateRoleRequest struct {
	NewRole masterrolemodels.Role `json:"newRole"`
}

type CreateRoleResponse struct {
	NewRole masterrolemodels.Role `json:"newRole"`
	Err       error               `json:"-"`
}

func (c CreateRoleResponse) Failed() error { return c.Err }

type FindFirstRoleRequest struct {
	CurrentRole masterrolemodels.SearchableRole `json:"currentRole"`
}

type FindFirstRoleResponse struct {
	CurrentRole masterrolemodels.Role `json:"currentRole"`
	Err           error               `json:"-"`
}

func (f FindFirstRoleResponse) Failed() error { return f.Err }

type GetAllRolesRequest struct {
}

type GetAllRolesResponse struct {
	Roles []masterrolemodels.Role `json:"tenants"`
	Err     error                 `json:"-"`
}

func (g GetAllRolesResponse) Failed() error { return g.Err }

type FindRolesRequest struct {
	CurrentRole masterrolemodels.SearchableRole `json:"currentRole"`
}

type FindRolesResponse struct {
	CurrentRole []masterrolemodels.Role `json:"currentRole"`
	Err           error                 `json:"-"`
}

func (f FindRolesResponse) Failed() error { return f.Err }

type UpdateRoleRequest struct {
	NewRole masterrolemodels.Role `json:"newRole"`
}

type UpdateRoleResponse struct {
	UpdatedRole masterrolemodels.Role `json:"updatedRole"`
	Err error `json:"-"`
}

func (u UpdateRoleResponse) Failed() error { return u.Err }

type GetPagedRoleResponse struct {
	PagedRoles commonmodels.PagedResponse[masterrolemodels.Role] `json:"pagedRoles"`
	Err error `json:"-"`
}

func (g GetPagedRoleResponse) Failed() error { return g.Err }



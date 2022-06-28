package rolesendpoint

import (
	commonmodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/models"
	rolemodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/user/role/models"
)

type CreateRoleRequest struct {
	NewRole rolemodels.Role `json:"newRole"`
}

type CreateRoleResponse struct {
	NewRole rolemodels.Role `json:"newRole"`
	Err       error               `json:"-"`
}

func (c CreateRoleResponse) Failed() error { return c.Err }

type FindFirstRoleRequest struct {
	CurrentRole rolemodels.SearchableRole `json:"currentRole"`
}

type FindFirstRoleResponse struct {
	CurrentRole rolemodels.Role `json:"currentRole"`
	Err           error               `json:"-"`
}

func (f FindFirstRoleResponse) Failed() error { return f.Err }

type GetAllRolesRequest struct {
}

type GetAllRolesResponse struct {
	Roles []rolemodels.Role `json:"Roles"`
	Err     error                 `json:"-"`
}

func (g GetAllRolesResponse) Failed() error { return g.Err }

type FindRolesRequest struct {
	CurrentRole rolemodels.SearchableRole `json:"currentRole"`
}

type FindRolesResponse struct {
	CurrentRole []rolemodels.Role `json:"currentRole"`
	Err           error                 `json:"-"`
}

func (f FindRolesResponse) Failed() error { return f.Err }

type UpdateRoleRequest struct {
	NewRole rolemodels.Role `json:"newRole"`
}

type UpdateRoleResponse struct {
	UpdatedRole rolemodels.Role `json:"updatedRole"`
	Err error `json:"-"`
}

func (u UpdateRoleResponse) Failed() error { return u.Err }

type GetPagedRoleResponse struct {
	PagedRoles commonmodels.PagedResponse[rolemodels.Role] `json:"pagedRoles"`
	Err error `json:"-"`
}

func (g GetPagedRoleResponse) Failed() error { return g.Err }



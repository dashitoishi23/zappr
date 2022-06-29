package masterrolemodels

import "dev.azure.com/technovert-vso/Zappr/_git/Zappr/state"

type SearchableRole struct {
	Identifier       string `json:"identifier"`
	Name             string `json:"name"`
	CreatedOn        string `json:"createdOn"`
	ModifiedOn       string `json:"modifiedOn"`
	TenantIdentifier string `json:"tenantIdentifier"`
}

func (s *SearchableRole) AddTenant() {
	s.TenantIdentifier = state.GetState().GetUserContext().UserTenant
}
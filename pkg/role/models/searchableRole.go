package masterrolemodels

import (
	commonmodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/models"
)

type SearchableRole struct {
	Identifier       string `json:"identifier"`
	Name             string `json:"name"`
	CreatedOn        string `json:"createdOn"`
	ModifiedOn       string `json:"modifiedOn"`
	TenantIdentifier string `json:"tenantIdentifier"`
}

func (s *SearchableRole) AddTenant(requestScope commonmodels.RequestScope) {
	s.TenantIdentifier = requestScope.UserTenant
}
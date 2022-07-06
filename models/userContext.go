package commonmodels

type UserContext struct {
	UserTenant     string
	UserIdentifier string
	UserScopes    []string
}
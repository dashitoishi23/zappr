package usermodels

type UserLogin struct {
	Email            string `json:"email" validate:"nonzero"`
	Password         string `json:"password" validate:"nonzero"`
	TenantIdentifier string `json:"tenantIdentifier" validate:"nonzero"`
}
package usermodels

type UpdateUser struct {
	Identifier        string `json:"identifier" gorm:"primaryKey"`
	Name              string `json:"name"`
	ProfilePictureURL string `json:"profilePictureUrl"`
	TenantIdentifier  string `json:"tenantIdentifier"`
}
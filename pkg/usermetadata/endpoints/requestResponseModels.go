package usermetadataendpoints

import usermetadatamodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/usermetadata/models"

type AddUserMetadataRequest struct {
	NewUserMetadata usermetadatamodels.UserMetadata `json:"newUserMetadata"`
}

type AddUserMetadataResponse struct {
	NewUserMetadata usermetadatamodels.UserMetadata `json:"newUserMetadata"`
	Err error `json:"-"`
}

func (a *AddUserMetadataResponse) Failed() error { return a.Err }

type GetUserMetadataRequest struct {
	Query map[string]interface{} `json:"query"`
	EntityName string `json:"entityName"`
}

type GetUserMetadataResponse struct {
	UserMetadata []usermetadatamodels.UserMetadata `json:"userMetadata"`
	Err error `json:"-"`
}

func (g *GetUserMetadataResponse) Failed() error { return g.Err }
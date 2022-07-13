package usermetadataendpoints

import (
	"encoding/json"

	commonmodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/models"
	usermetadatamodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/usermetadata/models"
)

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

type GetMetadataByEntityRequest struct {

}

type GetMetadataByEntityResponse struct {
	UserMetadata []json.RawMessage `json:"metadata"`
	Err error `json:"-"`
}

func (g *GetMetadataByEntityResponse) Failed() error { return g.Err }

type GetMetadataByEntityPagedRequest struct {
	Query map[string]interface{} `json:"query"`
	EntityName string `json:"entityName"`
	Page   int `json:"page"`
	Size   int `json:"size"`
}

type GetMetadataByEntityPagedResponse struct {
	commonmodels.PagedResponse[json.RawMessage] `json:"userMetadata"`
	Err error `json:"-"`
}

func (g *GetMetadataByEntityPagedResponse) Failed() error { return g.Err }

type UpdateMetadataRequest struct {
	CurrentQuery map[string]interface{} `json:"currentQuery"`
	UpdatedQuery map[string]interface{} `json:"updatedQuery"`
	EntityName string `json:"entityName"`
}

type UpdateMetadataResponse struct {
	UpdatedQuery map[string]interface{} `json:"updatedQuery"`
	Err error `json:"-"`
}

func (u *UpdateMetadataResponse) Failed() error { return u.Err }
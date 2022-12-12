package initsetupendpoints

import initsetupmodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/initsetup/models"

type AddConfigRequest struct {
	NewConfig initsetupmodels.Config `json:"newConfig"`
}

type AddConfigResponse struct {
	Config initsetupmodels.Config `json:"config"`
	Err error `json:"-"`
}

func (a *AddConfigResponse) Failed() error { return a.Err } 

type EditConfigRequest struct {
	NewConfig initsetupmodels.Config `json:"newConfig"`
}

type EditConfigResponse struct {
	IsUpdated bool `json:"isUpdated"`
	Err error `json:"-"`
}

func (a *EditConfigResponse) Failed() error { return a.Err } 
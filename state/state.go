package state

import commonmodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/models"

type GlobalState struct {
	UserContext commonmodels.UserContext
}

var stateAddress *GlobalState

func (g *GlobalState) GetUserContext() commonmodels.UserContext {
	return stateAddress.UserContext
}

func (g *GlobalState) SetUserContext(userContext commonmodels.UserContext) {
	stateAddress.UserContext = userContext
}

func (g *GlobalState) IsAllowedToWrite() bool {
	for _, scope := range g.UserContext.UserScopes {
		if scope == "write" {
			return true
		}
	}

	return false
}

func (g *GlobalState) IsAllowedToDelete() bool {
		for _, scope := range g.UserContext.UserScopes {
		if scope == "delete" {
			return true
		}
	}

	return false
}

func (g *GlobalState) IsAllowedToUpdate() bool {
		for _, scope := range g.UserContext.UserScopes {
		if scope == "update" {
			return true
		}
	}

	return false
}


func GetState() *GlobalState {
	return stateAddress
}

func InitState() {
	stateAddress = &GlobalState{
		UserContext: commonmodels.UserContext{
			UserTenant: "",
			UserIdentifier: "",
			UserScopes: []string{},
		},
	}
}
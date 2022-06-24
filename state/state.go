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

func GetState() *GlobalState {
	return stateAddress
}

func InitState() {
	stateAddress = &GlobalState{
		UserContext: commonmodels.UserContext{
			UserTenant: "",
			UserIdentifier: "",
		},
	}
}
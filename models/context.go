package commonmodels

type RequestScope struct {
	UserTenant     string
	UserIdentifier string
	UserScopes     []string
}

func (r *RequestScope) IsAllowedToWrite() bool {
	for _, scope := range r.UserScopes {
		if scope == "write" {
			return true
		}
	}

	return false
}

func (r *RequestScope) IsAllowedToDelete() bool {
	for _, scope := range r.UserScopes {
		if scope == "delete" {
			return true
		}
	}

	return false
}

func (r *RequestScope) IsAllowedToUpdate() bool {
	for _, scope := range r.UserScopes {
		if scope == "update" {
			return true
		}
	}

	return false
}
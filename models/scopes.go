package commonmodels

func GetScopes(roleType string) []string {
	switch roleType {
	case "Admin":
		return []string{"read", "write", "update", "delete"}

	case "Normal User":
	default:
		return []string{"read"}
	}

	return []string{"read"}
}
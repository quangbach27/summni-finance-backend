package auth

type Claims struct {
	Subject       string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	PreferredUser string `json:"preferred_username"`
	// Keycloak specific roles are often nested under 'realm_access'
	RealmAccess struct {
		Roles []string `json:"roles"`
	} `json:"realm_access"`
}

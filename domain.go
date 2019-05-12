package main

type AccessToken struct {
	Token string
}

type AccessTokenResponse struct {
	RefreshTokenExpiresIn string `json:"refresh_token_expires_in"`
	APIProductList        string `json:"api_product_list"`
	OrganizationName      string `json:"organization_name"`
	DeveloperEmail        string `json:"developer.email"`
	TokenType             string `json:"token_type"`
	IssuedAt              string `json:"issued_at"`
	ClientID              string `json:"client_id"`
	AccessToken           string `json:"access_token"`
	ApplicationName       string `json:"application_name"`
	Scope                 string `json:"scope"`
	ExpiresIn             string `json:"expires_in"`
	RefreshCount          string `json:"refresh_count"`
	Status                string `json:"status"`
}

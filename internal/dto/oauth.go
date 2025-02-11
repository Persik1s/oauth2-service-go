package dto

type OGoogleTokenRequestDto struct {
	ClientId     string `url:"client_id" json:"client_id"`
	ClientSecret string `url:"client_secret" json:"client_secret"`
	RedirectUri  string `url:"redirect_uri" json:"redirect_uri"`
	GrantType    string `url:"grant_type" json:"garnt_type"`
	Code         string `url:"code" json:"code"`
}

type OAuthTokenResponeDto struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	TokenType   string `json:"token_type"`
	IdToken     string `json:"id_token"`
}

type OGoogleUserDto struct {
	Email    string `json:"email"`
	Username string `json:"name"`
}

package dto

type TokenPair struct {
	RefreshToken string `json:"refresh"`
	AccessToken  string `json:"access"`
}

package domain

type AuthToken struct {
	AccessToken  string `mapstructure:"access_token" json:"access_token"`
	RefreshToken string `mapstructure:"refresh_token" json:"refresh_token"`
	TokenType    string `mapstructure:"token_type" json:"token_type"`
	ExpiredIn    int    `mapstructure:"expired_in" json:"expired_in"`
	IDToken      string `mapstructure:"id_token" json:"id_token"`
}

func NewAuthToken(accessToken string, refreshToken string, tokenType string, expiredIn int, idToken string) *AuthToken {
	return &AuthToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    tokenType,
		ExpiredIn:    expiredIn,
		IDToken:      idToken,
	}
}

package domain

type AuthToken struct {
	AccessToken  string `mapstructure:"accessToken" json:"accessToken"`
	RefreshToken string `mapstructure:"refreshToken" json:"refreshToken"`
	TokenType    string `mapstructure:"tokenType" json:"tokenType"`
	ExpiredIn    int    `mapstructure:"expiredIn" json:"expiredIn"`
	IDToken      string `mapstructure:"idToken" json:"idToken"`
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

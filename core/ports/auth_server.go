package ports

type AuthServerService interface {
	// AccessToken get access token
	AccessToken() (string, error)
}

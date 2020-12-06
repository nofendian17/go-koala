package models

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type RefreshTokenResponse struct {
	AccesToken   string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AccessDetails struct {
	AccessUuid string
	CustomerId string
}

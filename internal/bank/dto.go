package bank

type authorizeRequest struct {
	GrantType    string `json:"grant_type"`
	Login        string `json:"login"`
	Password     string `json:"password"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type authorizeResponse struct {
	AccessToken string `json:"access_token"`
}

type discoveryResponse struct {
	Login string `json:"login"`
}

type discoveryAppResponse struct {
	Lift string `json:"lift"`
}

type liftRequest struct {
	QrCodeId string `json:"qr_code_id"`
	Type     string `json:"type"`
}

type liftResponse struct {
	AccessToken string `json:"access_token"`
	Links       struct {
		RevokeToken struct {
			Href string `json:"href"`
		} `json:"revoke_token"`
	} `json:"_links"`
}

package model

type FacebookUserNormal struct {
	Email		string	`json:"email"`
	Id			string	`json:"id"`
	Gender		string	`json:"gender"`
	Name		string	`json:"name"`
	Link		string	`json:"link"`
	Picture		struct {
		Data		struct {
			Url		string	`json:"url"`
		}	`json:"data"`
	}	`json:"picture"`
}

type OAuth_accessTokenResponse struct {
	AccessToken		string	`json:"access_token"`
	TokenType		string	`json:"token_type"`
	ExpiresIn		int64	`json:"expires_in"`
}
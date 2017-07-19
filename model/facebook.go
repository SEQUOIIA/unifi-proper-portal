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
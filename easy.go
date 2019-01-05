package main

//easyjson:json
type EasyAccount struct {
	ID        int         `json:"id"`
	Joined    int         `json:"joined"`
	Birth     int         `json:"birth"`
	Country   string      `json:"country"`
	City      string      `json:"city"`
	Email     string      `json:"email"`
	Sname     string      `json:"sname"`
	Fname     string      `json:"fname"`
	Sex       string      `json:"sex"`
	Status    string      `json:"status"`
	Premium   EasyPremium `json:"premium"`
	Likes     []EasyLike  `json:"likes"`
	Interests []string    `json:"interests"`
}

//easyjson:json
type EasyPremium struct {
	Finish int `json:"finish"`
	Start  int `json:"start"`
}

//easyjson:json
type EasyLike struct {
	ID int `json:"id"`
	Ts int `json:"ts"`
}

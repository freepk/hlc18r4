package main

//easyjson:json
type EasyAccount struct {
	ID      int         `json:"id"`
	Joined  int         `json:"joined"`
	Birth   int         `json:"birth"`
	Premium EasyPremium `json:"premium"`
	Likes   []EasyLike  `json:"likes"`
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

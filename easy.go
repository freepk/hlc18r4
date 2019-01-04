package main

//easyjson:json
type EasyAccount struct {
	ID      int `json:"id"`
	Joined  int `json:"joined"`
	Birth   int `json:"birth"`
	Premium struct {
		Finish int `json:"finish"`
		Start  int `json:"start"`
	} `json:"premium"`
	Likes []struct {
		ID int `json:"id"`
		Ts int `json:"ts"`
	} `json:"likes"`
}

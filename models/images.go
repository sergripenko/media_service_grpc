package models

type Image struct {
	Id        string
	UserId    string
	Filename  string
	Height    int32
	Width     int32
	UniqId    string
	OrigImgId string
	Url       string
}

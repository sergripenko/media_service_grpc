package models

import "github.com/micro/protobuf/ptypes/timestamp"

type Image struct {
	Id        string
	UserId    string
	Filename  string
	Height    int32
	Width     int32
	UniqId    string
	OrigImgId string
	Url       string
	CreatedAt *timestamp.Timestamp
	UpdatedAt *timestamp.Timestamp
}

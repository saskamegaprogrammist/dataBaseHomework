package models

import "time"

type Post struct {
	Id int32 `json:"id"`
	Message string `json:"message"`
	Date time.Time `json:"created"`
	Parent int32 `json:"parent"`
	Edited bool `json:"isEdited"`
	User string `json:"author"`
	Forum string `json:"forum"`
	Thread int32 `json:"thread"`
}

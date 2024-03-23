package domain

import "time"

// User 领域对象，是entity
// Bo
type User struct {
	Id       int64
	Email    string
	Phone    string
	Password string
	NickName string
	Birthday string
	Desc     string
	Created  time.Time
}

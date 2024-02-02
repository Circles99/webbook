package domain

import "time"

// User 领域对象，是entity
// Bo
type User struct {
	Id       int64
	Email    string
	Password string
	Created  time.Time
}

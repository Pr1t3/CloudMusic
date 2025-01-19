package models

import (
	"time"
)

type Claims struct {
	Email  string
	UserId int
	exp    time.Time
	iat    time.Time
}

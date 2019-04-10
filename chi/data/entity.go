package data;

import (
	"time"
)

type User struct {
    IDX			int64	`json:"idx"`
    UserId		string	`json:"user_id"`
	CreateDt	time.Time	`json:"create_dt"`
	UpdateDt	time.Time	`json:"update_dt"`
}
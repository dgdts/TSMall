package bizcontext

import "context"

type User struct {
	UserName string `json:"user_name"`
}

type BizContext struct {
	context.Context

	User *User
}

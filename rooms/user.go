package rooms

import "errors"

type User interface {
	GetId() string
}

type UserServicer interface {
	GetUser(string) (User, error)
}

var (
	ERROR_RETRIEVING_USER = errors.New("ERROR_RETRIEVING_USER")
)

type UserInfo struct {
	CheerImages []*CheerImage
	User        User
}

func NewUserAndInfo(user User) *UserInfo {
	return &UserInfo{
		[]*CheerImage{},
		user,
	}
}

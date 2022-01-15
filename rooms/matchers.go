package rooms

import "github.com/aljorhythm/yumseng/cheers"

type roomMatcher struct {
	matcherFn func(*Room) bool
}

func (r roomMatcher) Matches(x interface{}) bool {
	inputRoom := x.(*Room)
	return r.matcherFn(inputRoom)
}

func (r roomMatcher) String() string {
	return "room matches custom criteria"
}

type userMatcher struct {
	matcherFn func(user User) bool
}

func (r userMatcher) Matches(x interface{}) bool {
	user := x.(User)
	return r.matcherFn(user)
}

func (r userMatcher) String() string {
	return "user matches custom criteria"
}

type cheerMatcher struct {
	matcherFn func(*cheers.Cheer) bool
}

func (r cheerMatcher) Matches(arg interface{}) bool {
	cheer := arg.(*cheers.Cheer)
	return r.matcherFn(cheer)
}

func (r cheerMatcher) String() string {
	return "cheer matches custom criteria"
}

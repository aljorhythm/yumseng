package rooms

import (
	"github.com/aljorhythm/yumseng/cheers"
	"github.com/aljorhythm/yumseng/utils/movingavg"
	"time"
)

type Room struct {
	Cheers []*cheers.Cheer
	Name   string
	movingavg.MovingAverageCalculator
	Users map[string]*UserInfo
}

type CheerItem struct {
	cheers.Cheer
}

func (item CheerItem) GetTime() time.Time {
	return item.Cheer.ClientCreatedAt
}

func (room *Room) AddCheer(cheer *cheers.Cheer) {
	room.Cheers = append(room.Cheers, cheer)
	room.MovingAverageCalculator.AddItem(CheerItem{*cheer})
}

func (room *Room) AddUserIfNotPresent(user User) (bool, error) {
	id := user.GetId()

	if _, ok := room.Users[id]; !ok {
		room.Users[id] = NewUserAndInfo(user)
		return true, nil
	}
	return false, nil
}

func (room *Room) GetUserInfo(user User) (*UserInfo, error) {
	if userInfo, ok := room.Users[user.GetId()]; ok {
		return userInfo, nil
	} else {
		return nil, ERROR_NO_USER_INFO
	}
}

func (room *Room) GetCheerImages(user User) ([]*CheerImage, error) {
	if userInfo, err := room.GetUserInfo(user); err != nil {
		return nil, err
	} else {
		return userInfo.CheerImages, nil
	}
}

func (room *Room) AddCheerImage(user User, cheerImage *CheerImage) error {
	if userInfo, err := room.GetUserInfo(user); err != nil {
		return err
	} else {
		userInfo.CheerImages = append(userInfo.CheerImages, cheerImage)
		return nil
	}
}

func NewRoom(name string) *Room {
	return &Room{
		[]*cheers.Cheer{},
		name,
		movingavg.NewCalculator(movingavg.NowTime{}),
		map[string]*UserInfo{},
	}
}

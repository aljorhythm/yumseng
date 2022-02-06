package rooms

import (
	"errors"
	"fmt"
	"github.com/aljorhythm/yumseng/cheers"
	"github.com/aljorhythm/yumseng/utils/movingavg"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"
)

type Room struct {
	Cheers        []*cheers.Cheer
	Name          string
	calculator    movingavg.MovingAverageCalculator
	Users         map[string]*UserInfo
	usersRWMutex  *sync.RWMutex
	cheersRWMutex *sync.RWMutex
	toSkipCheerFn func(cheer *cheers.Cheer) bool
}

type CheerItem struct {
	*cheers.Cheer
}

func (item CheerItem) GetTime() time.Time {
	return item.Cheer.ClientCreatedAt
}

type CheerUser struct {
	*cheers.Cheer
}

func (c CheerUser) GetId() string {
	return c.UserId
}

func (room *Room) AddCheer(cheer *cheers.Cheer) error {
	room.cheersRWMutex.Lock()
	defer room.cheersRWMutex.Unlock()
	if cheer.ClientCreatedAt.IsZero() {
		return errors.New(fmt.Sprintf("%#v ClientCreatedAt cannot be 0", cheer))
	}

	room.Cheers = append(room.Cheers, cheer)

	cheerUser := CheerUser{cheer}
	_, err := room.AddUserIfNotPresent(cheerUser)

	if err != nil {
		return err
	}

	points := rand.Intn(50)
	room.addUserPoints(cheerUser, points)

	return nil
}

func (room *Room) AddUserIfNotPresent(user User) (bool, error) {
	id := user.GetId()

	room.usersRWMutex.Lock()
	defer room.usersRWMutex.Unlock()

	if _, ok := room.Users[id]; !ok {
		room.Users[id] = NewUserAndInfo(user)
		return true, nil
	}

	return false, nil
}

func (room *Room) addUserPoints(userArg User, points int) error {
	if user, ok := room.Users[userArg.GetId()]; ok {
		user.Points += points
		return nil
	} else {
		return errors.New(fmt.Sprintf("no user found %#v", userArg.GetId()))
	}
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

const expectedMaxPerSecond float32 = 7.2

func (room *Room) Intensity(filterCheer func(cheer cheers.Cheer) bool) float32 {
	items := []movingavg.Item{}
	for _, item := range room.Cheers {
		if filterCheer(*item) {
			items = append(items, CheerItem{item})
		}
	}

	users := room.Users
	usersCount := 0

	// todo remove hardcoding of excluded user
	for _, user := range users {
		if user.User.GetId() != "global-user" {
			usersCount += 1
		}
	}

	if usersCount == 0 {
		return 0
	}

	expectedMax := (float32(usersCount) * expectedMaxPerSecond)
	count := room.calculator.CountFrom((time.Duration(1) * time.Second), items)
	intensity := float32(count) / expectedMax
	return intensity
}

func (room *Room) removeOutdatedCheers() {
	room.cheersRWMutex.Lock()
	defer room.cheersRWMutex.Unlock()

	beforeTime := time.Now().Add(-time.Second * time.Duration(5))

	allCheers := room.Cheers
	filteredCheers := []*cheers.Cheer{}
	removedCheerCount := 0
	for _, cheer := range allCheers {
		if cheer.ClientCreatedAt.After(beforeTime) {
			filteredCheers = append(filteredCheers, cheer)
		} else {
			removedCheerCount++
		}
	}

	room.Cheers = filteredCheers
	log.Printf("[Room#removeOutdatedCheers] room: %s cheers count: %d removed: %d", room.Name, len(room.Cheers), removedCheerCount)
}

func (room *Room) GetUsers() []*UserInfo {
	room.usersRWMutex.Lock()
	defer room.usersRWMutex.Unlock()

	users := []*UserInfo{}
	for _, user := range room.Users {
		users = append(users, user)
	}
	return users
}

func (room *Room) toSkipCheer(cheer *cheers.Cheer) bool {
	if room.toSkipCheerFn == nil {
		return false
	} else {
		return room.toSkipCheerFn(cheer)
	}
}

func (room *Room) ResetPoints() {
	room.usersRWMutex.Lock()
	defer room.usersRWMutex.Unlock()

	for _, user := range room.Users {
		user.Points = 0
	}
}

func (room *Room) DeleteAllUsers() {
	room.usersRWMutex.Lock()
	defer room.usersRWMutex.Unlock()
	room.Users = map[string]*UserInfo{}
}

func ToSkipCheerIfSeng(cheer *cheers.Cheer) bool {
	return strings.ToLower(cheer.Value) == "seng"
}

func (room *Room) HasSkipCheerRule() bool {
	return room.toSkipCheerFn != nil
}

func (room *Room) ClearSkipCheerRule() {
	room.toSkipCheerFn = nil
}

func (room *Room) SkipAddCheerIfSeng() {
	room.toSkipCheerFn = ToSkipCheerIfSeng
}

func NewRoom(name string) *Room {
	return &Room{
		[]*cheers.Cheer{},
		name,
		movingavg.NewCalculator(movingavg.NowTime{}),
		map[string]*UserInfo{},
		&sync.RWMutex{},
		&sync.RWMutex{},
		// todo remove default behavior
		ToSkipCheerIfSeng,
	}
}

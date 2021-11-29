package room

import (
	"github.com/aljorhythm/yumseng/cheers"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestRoomEvents(t *testing.T) {
	t.Run("Should listen, send and receive a cheer", func(t *testing.T) {
		roomEvents := NewRoomEvents()
		room := NewRoom("room-1")
		clientId := "client-1"

		var got cheers.Cheer

		roomEvents.SubscribeCheerAdded(room, clientId, func(args ...interface{}) {
			log.Printf("got a cheer callback args: %#v", args)
			arg0 := args[0]
			got = arg0.(cheers.Cheer)
		})

		cheer := cheers.Cheer{
			Value: "yoohoo",
		}
		roomEvents.PublishCheerAdded(room, cheer)

		assert.Equal(t, cheer, got)
	})

	t.Run("Should listen, send and receive a cheer on different rooms", func(t *testing.T) {
		roomEvents := NewRoomEvents()

		room := NewRoom("room-1")

		clientOneId := "client-1"
		clientTwoId := "client-2"

		var clientOneGot cheers.Cheer
		var clientTwoGot cheers.Cheer

		roomEvents.SubscribeCheerAdded(room, clientOneId, func(args ...interface{}) {
			log.Printf("client one got a cheer callback args: %#v", args)
			arg0 := args[0]
			clientOneGot = arg0.(cheers.Cheer)
		})

		roomEvents.SubscribeCheerAdded(room, clientTwoId, func(args ...interface{}) {
			log.Printf("client two got a cheer callback args: %#v", args)
			arg0 := args[0]
			clientTwoGot = arg0.(cheers.Cheer)
		})

		cheer := cheers.Cheer{
			Value: "yoohoo",
		}
		roomEvents.PublishCheerAdded(room, cheer)

		assert.Equal(t, cheer, clientOneGot)
		assert.Equal(t, cheer, clientTwoGot)
	})
}

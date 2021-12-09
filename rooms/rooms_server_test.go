package rooms

import (
	"github.com/aljorhythm/yumseng/cheers"
	"github.com/aljorhythm/yumseng/objectstorage"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type MockUser struct {
	id string
}

func (m MockUser) GetId() string {
	return m.id
}

type MockUserService struct {
}

func (m MockUserService) GetUser(id string) (User, error) {
	return MockUser{id: id}, nil
}

func TestRoomServer(t *testing.T) {
	roomsServer := NewRoomsServer(mux.NewRouter(), MockUserService{}, objectstorage.NewInmemoryStore(), RoomsServerOpts{})

	t.Run("when we send a cheer it must be broadcasted", func(t *testing.T) {
		server := httptest.NewServer(roomsServer)
		defer server.Close()

		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/events"

		roomName := "test-room"
		userId := "test-user"

		header := http.Header{}

		ws, _, err := websocket.DefaultDialer.Dial(wsURL, header)
		if err != nil {
			t.Fatalf("could not open a ws connection on %s %v", wsURL, err)
		}
		defer ws.Close()

		joinRoomRequest := JoinRoomRequest{
			RoomName: roomName,
			UserId:   userId,
		}

		err = ws.WriteJSON(joinRoomRequest)
		assert.NoError(t, err)

		gotConnected := RoomConnectedMessage{}
		err = ws.ReadJSON(&gotConnected)
		assert.NoError(t, err)

		wantedConnected := RoomConnectedMessage{
			EventName: "EVENT_ROOM_CONNECTED",
			UserId:    userId,
			RoomName:  roomName,
		}

		assert.Equal(t, wantedConnected, gotConnected)
		t.Logf("room connected message %#v", gotConnected)

		gotIntensity := RoomLastSecondsCheerCountMessage{}
		wantedIntensity := RoomLastSecondsCheerCountMessage{
			EventName: "EVENT_LAST_SECONDS_COUNT",
			Count:     0,
		}
		err = ws.ReadJSON(&gotIntensity)

		assert.NoError(t, err)
		assert.Equal(t, wantedIntensity, gotIntensity)

		cheer := cheers.Cheer{
			Value:           "this is a cheer",
			ClientCreatedAt: time.Now().UTC(),
		}

		if err := ws.WriteJSON(cheer); err != nil {
			t.Fatalf("could not send message over ws connection %v", err)
		}

		gotCheer := CheerAddedMessage{}
		err = ws.ReadJSON(&gotCheer)
		assert.NoError(t, err)

		cheer.UserId = userId
		wantedCheer := CheerAddedMessage{cheer, EVENT_CHEER_ADDED.name}

		assert.Equal(t, wantedCheer, gotCheer)

		t.Run("Recurring intensity messages interleaved with cheer event", func(t *testing.T) {
			//todo
		})
	})
}

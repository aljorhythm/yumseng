package rooms

import (
	"bytes"
	"github.com/aljorhythm/yumseng/cheers"
	"github.com/aljorhythm/yumseng/utils"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestRoomServer(t *testing.T) {
	t.Run("when we send a cheer it must be broadcasted", func(t *testing.T) {
		roomsServer := NewRoomsServer(mux.NewRouter())
		server := httptest.NewServer(roomsServer)
		defer server.Close()

		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/events"
		header := http.Header{}
		header.Add("rooms-name", "test-rooms")
		ws, _, err := websocket.DefaultDialer.Dial(wsURL, header)
		if err != nil {
			t.Fatalf("could not open a ws connection on %s %v", wsURL, err)
		}
		defer ws.Close()

		cheer := cheers.Cheer{
			Value:           "this is a cheer",
			ClientCreatedAt: time.Now().UTC(),
		}

		message := utils.MustEncodeJson(cheer)
		if err := ws.WriteMessage(websocket.TextMessage, message); err != nil {
			t.Fatalf("could not send message over ws connection %v", err)
		}

		_, rawMessage, err := ws.ReadMessage()
		if err != nil {
			t.Fatalf("error reading cheer %#v", err)
		}

		got := CheerAddedMessage{}
		utils.DecodeJson(bytes.NewReader(rawMessage), &got)

		wanted := CheerAddedMessage{cheer, string(EVENT_CHEER_ADDED.name)}

		assert.Equal(t, wanted, got)
	})
}

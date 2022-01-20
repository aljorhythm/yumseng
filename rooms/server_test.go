package rooms

import (
	"bytes"
	"context"
	"github.com/aljorhythm/yumseng/cheers"
	"github.com/aljorhythm/yumseng/objectstorage"
	"github.com/aljorhythm/yumseng/utils"
	"github.com/aljorhythm/yumseng/utils/movingavg"
	"github.com/golang/mock/gomock"
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

func TestRoomServerUserImages(t *testing.T) {
	storage := objectstorage.NewInmemoryStore()
	service := NewRoomsService(storage)
	userService := MockUserService{}
	roomsServer := NewRoomsServer(mux.NewRouter(), service, userService, RoomsServerOpts{})

	t.Run("server should respond with error when room does not exist", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/room-1/user/user-1/images", nil)
		roomsServer.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("server should respond with correct status code", func(t *testing.T) {
		t.Run("room with no images for user", func(t *testing.T) {
			room2 := service.GetOrCreateRoom("room-2")
			userId := "user-1"
			user, _ := userService.GetUser(userId)
			service.UserJoinsRoom(context.Background(), room2, user)
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/room-2/user/user-1/images", nil)
			roomsServer.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusOK, recorder.Code, recorder.Body.String())

			t.Run("room with images for user", func(t *testing.T) {
				data := []byte("asd")
				image, err := service.AddCheerImage(context.Background(), room2.Name, user, data)
				assert.NoError(t, err)
				recorder := httptest.NewRecorder()
				request := httptest.NewRequest(http.MethodGet, "/room-2/user/user-1/images", nil)
				roomsServer.ServeHTTP(recorder, request)

				got := []*CheerImage{}
				utils.DecodeJson(recorder.Body, &got)

				want := []*CheerImage{
					image,
				}

				assert.Equal(t, http.StatusOK, recorder.Code)
				assert.Equal(t, want, got)

				t.Run("post image", func(t *testing.T) {
					recorder := httptest.NewRecorder()
					data := []byte("image-data")
					dataReader := bytes.NewReader(data)
					request := httptest.NewRequest(http.MethodPost, "/room-2/user/user-1/images", dataReader)
					roomsServer.ServeHTTP(recorder, request)

					assert.Equal(t, http.StatusOK, recorder.Code)

					response := CheerImage{}
					err := utils.DecodeJson(recorder.Body, &response)
					assert.NoError(t, err)

					objectId := response.ObjectId
					gotData, err := storage.Retrieve(context.Background(), objectId)
					assert.NoError(t, err)
					assert.Equal(t, data, gotData)
				})
			})
		})
	})
}

func TestRoomsServerCheers(t *testing.T) {

	t.Run("add cheer should be successful", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		mockRoomService := NewMockRoomServicer(mockCtrl)
		mockRoomService.EXPECT().GetOrCreateRoom(gomock.Eq("room-2")).Return(&Room{
			Cheers:     nil,
			Name:       "room-2",
			calculator: movingavg.NewCalculator(movingavg.NowTime{}),
			Users:      nil,
		})
		mockRoomService.EXPECT().AddCheer(
			roomMatcher{func(room *Room) bool {
				return room.Name == "room-2"
			}},
			cheerMatcher{matcherFn: func(cheer *cheers.Cheer) bool {
				return cheer.UserId == "dummy-user"
			}},
			userMatcher{matcherFn: func(user User) bool {
				return user.GetId() == "dummy-user"
			}},
		).Return(nil)
		roomsServer := NewRoomsServer(mux.NewRouter(), mockRoomService, MockUserService{}, RoomsServerOpts{})
		server := httptest.NewServer(roomsServer)
		defer server.Close()

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/room-2/user/dummy-user/cheers",
			utils.ToJsonReaderPanic(cheers.Cheer{
				Value:           "",
				ClientCreatedAt: time.Time{},
				UserId:          "dummy-user",
				ImageUrl:        "",
			}))

		roomsServer.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assertAllResponses(t, recorder)
	})
}

func TestRoomServerEventsSocket(t *testing.T) {
	storage := objectstorage.NewInmemoryStore()
	service := NewRoomsService(storage)
	roomsServer := NewRoomsServer(mux.NewRouter(), service, MockUserService{}, RoomsServerOpts{})

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

		wantedIntensity := RoomIntensityMessage{
			EventName: "EVENT_INTENSITY",
			Intensity: 0,
		}
		gotIntensity := RoomIntensityMessage{}
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

package rooms

import (
	"context"
	"github.com/aljorhythm/yumseng/objectstorage"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRoomsService(t *testing.T) {
	storage := objectstorage.NewInmemoryStore()
	service := NewRoomsService(storage)
	ctx := context.Background()
	roomName := "test-room"

	t.Run("add room", func(t *testing.T) {
		room := service.GetOrCreateRoom(roomName)
		userId := "test-user-1"
		user := MockUser{id: userId}

		t.Run("add user", func(t *testing.T) {
			err := service.UserJoinsRoom(ctx, room, user)
			assert.NoError(t, err)

			t.Run("add cheer image", func(t *testing.T) {
				objectId := "cheer-image-1"

				data := []byte("fake-image")

				err := service.AddCheerImage(ctx, room, user, data, objectId)
				assert.NoError(t, err)

				gotCheerImages, err := service.GetCheerImages(ctx, room, user)
				assert.NoError(t, err)

				cheerImages := []*CheerImage{{
					ObjectId: "rooms/test-room/test-user-1/cheer-image-1",
					Url:      "rooms/test-room/test-user-1/cheer-image-1",
				}}
				assert.Equal(t, cheerImages, gotCheerImages)

				t.Run("get from storage", func(t *testing.T) {
					gotData, err := storage.Retrieve(ctx, gotCheerImages[0].ObjectId)

					assert.NoError(t, err)
					assert.Equal(t, gotData, data)
				})
			})
		})

	})
}

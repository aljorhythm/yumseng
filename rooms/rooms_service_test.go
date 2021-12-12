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
				data := []byte("fake-image")

				image, err := service.AddCheerImage(ctx, room.Name, user, data)
				assert.NoError(t, err)

				gotCheerImages, err := service.GetCheerImages(ctx, room.Name, user)
				assert.NoError(t, err)

				wantCheerImages := []*CheerImage{image}
				assert.Equal(t, wantCheerImages, gotCheerImages)

				t.Run("get from storage", func(t *testing.T) {
					gotData, err := storage.Retrieve(ctx, gotCheerImages[0].ObjectId)

					assert.NoError(t, err)
					assert.Equal(t, gotData, data)
				})
			})
		})

	})
}

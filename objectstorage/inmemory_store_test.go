package objectstorage

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInmemoryStore(t *testing.T) {
	ctx := context.Background()
	store := NewInmemoryStore()
	id := MustHaveUUIDString(t)

	data := []byte(MustHaveUUIDString(t))

	_, err := store.Retrieve(ctx, id)
	assert.Error(t, ERROR_DATA_NOT_FOUND)

	err = store.Store(ctx, id, data)
	assert.NoError(t, err)

	got, err := store.Retrieve(ctx, id)
	assert.NoError(t, err)

	assert.Equal(t, data, got)
}

func MustHaveUUIDString(t *testing.T) string {
	return MustHaveUUID(t).String()
}

func MustHaveUUID(t *testing.T) uuid.UUID {
	id, err := uuid.NewUUID()
	if err != nil {
		t.Fatalf("fail to generate uuid")
	}
	return id
}

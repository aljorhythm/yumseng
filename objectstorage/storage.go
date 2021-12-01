package objectstorage

import "context"

type Metadata struct {
	Url string
}

type Storage interface {
	Store(ctx context.Context, id string, data []byte) error
	Retrieve(context context.Context, id string) (data []byte, err error)
	RetrieveObjectMetadata(context context.Context, id string) (metadata *Metadata, err error)
	Delete(context context.Context, id string) error
}

package objectstorage

import "context"

type inmemoryStore struct {
	dataMap map[string][]byte
}

func (i *inmemoryStore) RetrieveObjectMetadata(context context.Context, id string) (metadata *Metadata, err error) {
	return &Metadata{Url: id}, nil
}

func (i *inmemoryStore) Delete(context context.Context, id string) error {
	delete(i.dataMap, id)
	return nil
}

func (i *inmemoryStore) Retrieve(ctx context.Context, id string) ([]byte, error) {
	if data, ok := i.dataMap[id]; ok {
		return data, nil
	} else {
		return nil, ERROR_DATA_NOT_FOUND
	}
}

func (i *inmemoryStore) Store(ctx context.Context, id string, bytes []byte) error {
	i.dataMap[id] = bytes
	return nil
}

func NewInmemoryStore() Storage {
	return &inmemoryStore{
		map[string][]byte{},
	}
}

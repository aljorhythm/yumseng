package objectstorage

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"os"
	"testing"
)

func readTestConfig(t *testing.T) DigitalOceanSpaceConfig {
	dat, err := os.ReadFile(".config.yml")

	if err != nil {
		t.Fatal("unable to read config.yml")
	}

	config := DigitalOceanSpaceConfig{}
	err = yaml.Unmarshal(dat, &config)

	if err != nil {
		t.Fatal("unable to parse config.yml")
	}
	return config
}

func TestDigitalOceanSpaceStorage(t *testing.T) {
	if os.Getenv("EXTERNAL_TEST") == "" {
		t.Skipf("env var EXTERNAL_TEST not set, skipping")
	}

	ctx := context.Background()

	config := readTestConfig(t)

	accessKey := config.AccessKey
	secKey := config.SecretKey
	endpoint := config.Endpoint

	bucketProps := BucketProperties{
		Name:     "yumseng",
		Location: "",
	}

	store, err := NewSpaceStore(ctx, endpoint, accessKey, secKey, bucketProps)

	assert.NoError(t, err)
	assert.NotNil(t, store)

	t.Run("list spaces", func(t *testing.T) {
		input := &s3.ListObjectsInput{
			Bucket: aws.String(bucketProps.Name),
		}

		objects, err := store.ListObjects(input)
		assert.NoError(t, err)

		t.Logf("length objects %d", len(objects.Contents))
		assert.NotZero(t, len(objects.Contents))

		for _, obj := range objects.Contents {
			t.Log(aws.StringValue(obj.Key))
		}
	})

	t.Run("list buckets", func(t *testing.T) {
		spaces, err := store.ListBuckets(nil)
		assert.NoError(t, err)

		t.Logf("buckets size %d", len(spaces.Buckets))
		assert.NotZero(t, len(spaces.Buckets))
		for _, b := range spaces.Buckets {
			t.Logf("bucket %s", (aws.StringValue(b.Name)))
		}
	})

	t.Run("store data", func(t *testing.T) {
		objectId := "testdata/object-1"
		data := []byte("abcde")

		t.Run("store object", func(t *testing.T) {
			err := store.Store(ctx, objectId, data)
			assert.NoError(t, err)

			t.Run("retrieve object", func(t *testing.T) {
				got, err := store.Retrieve(ctx, objectId)
				assert.NoError(t, err)
				assert.Equal(t, data, got)

				meta, err := store.RetrieveObjectMetadata(ctx, objectId)

				url := meta.Url
				t.Logf("url: %s", url)
				assert.NotEmpty(t, "url")

				t.Run("delete object", func(t *testing.T) {
					err := store.Delete(ctx, objectId)
					assert.NoError(t, err)
					assert.Equal(t, data, got)

					t.Run("retrieve object should fail", func(t *testing.T) {
						_, err := store.Retrieve(ctx, objectId)
						assert.Error(t, err)
					})
				})
			})
		})
	})
}

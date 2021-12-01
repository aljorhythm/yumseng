package objectstorage

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
)

type SpaceStore struct {
	*s3.S3
	BucketProperties
}

func (s *SpaceStore) RetrieveObjectMetadata(context context.Context, id string) (*Metadata, error) {
	url := fmt.Sprintf("%s/%s/%s", s.Client.Endpoint, *s.getBucketName(), id)

	metadata := &Metadata{
		Url: url,
	}
	return metadata, nil
}

func (store *SpaceStore) getBucketName() *string {
	return aws.String(store.BucketProperties.Name)
}

func (s *SpaceStore) Delete(ctx context.Context, id string) error {
	params := &s3.DeleteObjectInput{
		Bucket: s.getBucketName(),
		Key:    aws.String(id),
	}
	output, err := s.DeleteObject(params)

	log.Printf("delete object %s from bucket %s output %#v",
		*s.getBucketName(),
		id,
		output)

	return err
}

func (s *SpaceStore) Store(ctx context.Context, id string, data []byte) error {
	object := s3.PutObjectInput{
		Bucket: s.getBucketName(),
		Key:    aws.String(id),
		Body:   bytes.NewReader(data),
		ACL:    aws.String("public-read"),
	}

	output, err := s.PutObject(&object)

	log.Printf("uploaded object %s to bucket %s output %#v",
		*s.getBucketName(),
		id,
		output)
	return err
}

func (s *SpaceStore) Retrieve(ctx context.Context, id string) ([]byte, error) {
	input := &s3.GetObjectInput{
		Bucket: s.getBucketName(),
		Key:    aws.String(id),
	}
	result, err := s.GetObject(input)

	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)

	i, err := buf.ReadFrom(result.Body)

	log.Printf("successfully read object %s of length %d from bucket %s", id, i, *s.getBucketName())

	return buf.Bytes(), err
}

type BucketProperties struct {
	Name     string
	Location string
}

/**
https://docs.digitalocean.com/products/spaces/resources/s3-sdk-examples/
aws region is us-east-1 for digital ocean spaces
*/
func getRegion() *string {
	return aws.String("us-east-1")
}

func NewSpaceStore(ctx context.Context, endpoint string, key string, secret string, bucket BucketProperties) (*SpaceStore, error) {
	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(key, secret, ""),
		Endpoint:    aws.String(endpoint),
		Region:      getRegion(),
	}

	newSession, err := session.NewSession(s3Config)

	if err != nil {
		return nil, err
	}

	client := s3.New(newSession)

	return &SpaceStore{
		client,
		bucket,
	}, nil
}

type DigitalOceanSpaceConfig struct {
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
	Endpoint  string `yaml:"endpoint"`
}

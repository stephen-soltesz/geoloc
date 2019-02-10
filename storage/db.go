package storage

import (
	"bytes"
	"context"
	"io"
	"time"

	// storage "github.com/google/google-api-go-client/storage/v1"
	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

// obj := client.Bucket("downloader-mlab-sandbox").
// 		Object("Maxmind/2019/02/05/20190205T180204Z-GeoLite2-City.tar.gz")

// GetObject downloads the given file from the given bucket.
func GetObject(ctx context.Context, bucket string, filename string, timeout time.Duration) (*bytes.Buffer, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Create a new storage client.
	client, err := storage.NewClient(ctx, option.WithScopes(storage.ScopeReadOnly))
	if err != nil {
		return nil, err
	}

	obj := client.Bucket(bucket).Object(filename)

	// Read the object data.
	r, err := obj.NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	b := new(bytes.Buffer)
	if _, err := io.Copy(b, r); err != nil {
		return nil, err
	}

	return b, nil
}

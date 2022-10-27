package bucket

import (
	"context"
	"io"

	"gocloud.dev/blob"
)

func New(bucket *blob.Bucket) *Bucket {
	return &Bucket{bucket}
}

type Bucket struct {
	*blob.Bucket
}

func (b *Bucket) NewReader(ctx context.Context, key string) (io.ReadCloser, error) {
	return b.Bucket.NewReader(ctx, key, &blob.ReaderOptions{})
}

func (b *Bucket) NewWriter(ctx context.Context, key string) (io.WriteCloser, error) {
	return b.Bucket.NewWriter(ctx, key, &blob.WriterOptions{})
}

func (b *Bucket) GoString() string {
	return "&Bucket{}"
}

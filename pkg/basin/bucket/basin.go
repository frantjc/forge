package bucket

import (
	"context"
	"io"

	"gocloud.dev/blob"
)

func New(ctx context.Context, addr string) (*Bucket, error) {
	bucket, err := blob.OpenBucket(ctx, addr)
	if err != nil {
		return nil, err
	}

	return &Bucket{bucket}, nil
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

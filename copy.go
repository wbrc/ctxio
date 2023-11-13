package ctxio

import (
	"context"
	"io"
)

type readerFunc func(p []byte) (n int, err error)

func (f readerFunc) Read(p []byte) (n int, err error) { return f(p) }

func Copy(ctx context.Context, dst io.Writer, src io.Reader) (written int64, err error) {
	return io.Copy(dst, readerFunc(func(p []byte) (n int, err error) {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
			return src.Read(p)
		}
	}))
}

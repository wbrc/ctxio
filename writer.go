package ctxio

import (
	"context"
	"io"
)

type Writer struct {
	ctx context.Context
	w   io.Writer
}

func NewWriter(ctx context.Context, w io.Writer) Writer {
	return Writer{
		ctx: ctx,
		w:   w,
	}
}

func (w Writer) Write(p []byte) (n int, err error) {
	select {
	case <-w.ctx.Done():
		return 0, w.ctx.Err()
	default:
		return w.w.Write(p)
	}
}

func (w Writer) ReadFrom(r io.Reader) (n int64, err error) {
	select {
	case <-w.ctx.Done():
		return 0, w.ctx.Err()
	default:
		rt, ok := w.w.(io.ReaderFrom)
		if ok {
			return rt.ReadFrom(r)
		}
		return Copy(w.ctx, w.w, r)
	}
}

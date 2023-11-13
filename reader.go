package ctxio

import (
	"context"
	"io"
)

type Reader struct {
	r   io.Reader
	ctx context.Context
}

func NewReader(ctx context.Context, r io.Reader) Reader {
	return Reader{
		r:   r,
		ctx: ctx,
	}
}

func (r Reader) Read(p []byte) (n int, err error) {
	select {
	case <-r.ctx.Done():
		return 0, r.ctx.Err()
	default:
		return r.r.Read(p)
	}
}

func (r Reader) WriteTo(w io.Writer) (n int64, err error) {
	select {
	case <-r.ctx.Done():
		return 0, r.ctx.Err()
	default:
		wt, ok := r.r.(io.WriterTo)
		if ok {
			return wt.WriteTo(w)
		}
		return Copy(r.ctx, w, r.r)
	}
}

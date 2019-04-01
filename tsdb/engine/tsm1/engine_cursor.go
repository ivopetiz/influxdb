package tsm1

import (
	"context"

	"github.com/ivopetiz/influxdb/tsdb"
)

func (e *Engine) CreateCursorIterator(ctx context.Context) (tsdb.CursorIterator, error) {
	return &arrayCursorIterator{e: e}, nil
}

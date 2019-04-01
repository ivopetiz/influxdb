package tsm1

import (
	"context"
	"fmt"

	"github.com/ivopetiz/influxdb/models"
	"github.com/ivopetiz/influxdb/pkg/metrics"
	"github.com/ivopetiz/influxdb/query"
	"github.com/ivopetiz/influxdb/tsdb"
	"github.com/ivopetiz/influxql"
)

type arrayCursorIterator struct {
	e   *Engine
	key []byte

	asc struct {
		Float    *floatArrayAscendingCursor
		Integer  *integerArrayAscendingCursor
		Unsigned *unsignedArrayAscendingCursor
		Boolean  *booleanArrayAscendingCursor
		String   *stringArrayAscendingCursor
	}

	desc struct {
		Float    *floatArrayDescendingCursor
		Integer  *integerArrayDescendingCursor
		Unsigned *unsignedArrayDescendingCursor
		Boolean  *booleanArrayDescendingCursor
		String   *stringArrayDescendingCursor
	}
}

func (q *arrayCursorIterator) Next(ctx context.Context, r *tsdb.CursorRequest) (tsdb.Cursor, error) {
	// Look up fields for measurement.
	mf := q.e.fieldset.Fields(r.Name)
	if mf == nil {
		return nil, nil
	}

	// Find individual field.
	f := mf.Field(r.Field)
	if f == nil {
		// field doesn't exist for this measurement
		return nil, nil
	}

	if grp := metrics.GroupFromContext(ctx); grp != nil {
		grp.GetCounter(numberOfRefCursorsCounter).Add(1)
	}

	var opt query.IteratorOptions
	opt.Ascending = r.Ascending
	opt.StartTime = r.StartTime
	opt.EndTime = r.EndTime

	// Return appropriate cursor based on type.
	switch f.Type {
	case influxql.Float:
		return q.buildFloatArrayCursor(ctx, r.Name, r.Tags, r.Field, opt), nil
	case influxql.Integer:
		return q.buildIntegerArrayCursor(ctx, r.Name, r.Tags, r.Field, opt), nil
	case influxql.Unsigned:
		return q.buildUnsignedArrayCursor(ctx, r.Name, r.Tags, r.Field, opt), nil
	case influxql.String:
		return q.buildStringArrayCursor(ctx, r.Name, r.Tags, r.Field, opt), nil
	case influxql.Boolean:
		return q.buildBooleanArrayCursor(ctx, r.Name, r.Tags, r.Field, opt), nil
	default:
		panic(fmt.Sprintf("unreachable: %T", f.Type))
	}
}

func (q *arrayCursorIterator) seriesFieldKeyBytes(name []byte, tags models.Tags, field string) []byte {
	q.key = models.AppendMakeKey(q.key[:0], name, tags)
	q.key = append(q.key, keyFieldSeparatorBytes...)
	q.key = append(q.key, field...)
	return q.key
}

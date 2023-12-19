package query

import (
	"encoding/json"
	"strings"

	"github.com/google/go-cmp/cmp"
)

type Accessor interface {
	Get(key string) any
}

type MatcherFunc func(q *Query, data Accessor) bool

func eqMatcher(q *Query, data Accessor) bool {
	return cmp.Equal(data.Get(q.Field), q.Value)
}

func likeMatcher(q *Query, data Accessor) bool {
	return strings.Contains(data.Get(q.Field).(string), q.Value.(string))
}

func isnullMatcher(q *Query, data Accessor) bool {
	return data.Get(q.Field) == nil
}

func inverseMatch(fn MatcherFunc) MatcherFunc {
	return func(q *Query, data Accessor) bool {
		return !fn(q, data)
	}
}

type numberC[T comparable] struct {
	val T
	set bool
}

func decodeNumberValue(value any) (fl64 numberC[float64], i64 numberC[int64]) {
	switch v := value.(type) {
	case json.Number:
		if dotCount := strings.Count(v.String(), "."); dotCount > 0 {
			fl64.val, _ = v.Float64()
			fl64.set = true
		} else {
			i64.val, _ = v.Int64()
			i64.set = true
		}
	case float64:
		fl64.val = v
		fl64.set = true
	case int64:
		i64.val = v
		i64.set = true
	case float32:
		fl64.val = float64(v)
		fl64.set = true
	case int:
		i64.val = int64(v)
		i64.set = true
	}
	return
}

func compareNumbersFn(f1 numberC[float64], i1 numberC[int64], f2 numberC[float64], i2 numberC[int64], fn func(a, b float64) bool) bool {
	if f1.set {
		if f2.set {
			return fn(f1.val, f2.val)
		} else if i2.set {
			return fn(f1.val, float64(i2.val))
		}
	} else if i1.set {
		if i2.set {
			return fn(float64(i1.val), float64(i2.val))
		} else if f2.set {
			return fn(float64(i1.val), f2.val)
		}
	}
	return false
}

func compareNumbers(op Operator) MatcherFunc {
	if !op.IsComparative() {
		panic("compareValues called with non-comparative operator")
	}

	return func(q *Query, data Accessor) bool {
		value := data.Get(q.Field)
		if value == nil {
			return false
		}

		fl64FromQ, i64FromQ := decodeNumberValue(q.Value)
		fl64FromV, i64FromV := decodeNumberValue(value)

		return compareNumbersFn(fl64FromQ, i64FromQ, fl64FromV, i64FromV, func(a, b float64) bool {
			switch op {
			case OpGt:
				return a > b
			case OpGte:
				return a >= b
			case OpLt:
				return a < b
			case OpLte:
				return a <= b
			}
			return false
		})
	}
}

var matchers = map[Operator]MatcherFunc{
	OpEq:  eqMatcher,
	OpNeq: inverseMatch(eqMatcher),
	OpGt:  compareNumbers(OpGt),
	OpGte: compareNumbers(OpGte),
	OpLt:  compareNumbers(OpLt),
	OpLte: compareNumbers(OpLte),
	OpIn: func(q *Query, data Accessor) bool {
		for _, v := range q.Value.([]any) {
			if v == data.Get(q.Field) {
				return true
			}
		}
		return false
	},
	OpNin: func(q *Query, data Accessor) bool {
		for _, v := range q.Value.([]any) {
			if v == data.Get(q.Field) {
				return false
			}
		}
		return true
	},
	OpLike:    likeMatcher,
	OpNlike:   inverseMatch(likeMatcher),
	OpIsnull:  isnullMatcher,
	OpNotnull: inverseMatch(isnullMatcher),
}

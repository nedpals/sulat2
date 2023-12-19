package query

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type Operator string

func (op Operator) IsComparative() bool {
	return op == OpEq || op == OpNeq || op == OpGt || op == OpGte || op == OpLt || op == OpLte
}

func (op Operator) IsLogical() bool {
	return op == OpAnd || op == OpOr
}

func (op Operator) IsIn() bool {
	return op == OpIn || op == OpNin
}

func (op Operator) IsLike() bool {
	return op == OpLike || op == OpNlike
}

func (op Operator) IsNull() bool {
	return op == OpIsnull || op == OpNotnull
}

func (op Operator) IsBetween() bool {
	return op == OpBetween || op == OpNbetween
}

func (op Operator) IsSupported() bool {
	return supportedOperators[op]
}

const (
	OpAnd      Operator = "and"      // and
	OpOr       Operator = "or"       // or
	OpEq       Operator = "eq"       // equal to
	OpNeq      Operator = "neq"      // not equal to
	OpGt       Operator = "gt"       // greater than
	OpGte      Operator = "gte"      // greater than or equal to
	OpLt       Operator = "lt"       // less than
	OpLte      Operator = "lte"      // less than or equal to
	OpIn       Operator = "in"       // in
	OpNin      Operator = "nin"      // not in
	OpLike     Operator = "like"     // like
	OpNlike    Operator = "nlike"    // not like
	OpIsnull   Operator = "isnull"   // is null
	OpNotnull  Operator = "notnull"  // is not null
	OpBetween  Operator = "between"  // between
	OpNbetween Operator = "nbetween" // not between
)

// a Query is a collection of conditions
type Query struct {
	Field    string
	Operator Operator
	Value    any
	Options  map[string]any
}

func (q *Query) Match(data Accessor) bool {
	if q.Operator.IsLogical() {
		return q.matchLogical(data)
	}
	return matchers[q.Operator](q, data)
}

func (q *Query) matchLogical(data Accessor) bool {
	if q.Operator == OpAnd {
		for _, query := range q.Value.([]*Query) {
			if !query.Match(data) {
				return false
			}
		}
		return true
	}

	for _, query := range q.Value.([]*Query) {
		if query.Match(data) {
			return true
		}
	}
	return false
}

func (q *Query) String() string {
	sb := &strings.Builder{}
	sb.WriteString(string(q.Operator) + "(")
	if !q.Operator.IsLogical() {
		sb.WriteString(q.Field + " ")
	}
	json.NewEncoder(sb).Encode(q.Value)
	if q.Options != nil {
		sb.WriteByte(',')
		stringifyQueryOptions(q.Options, sb)
	}
	sb.WriteString(")")
	return sb.String()
}

func stringifyQueryOptions(options map[string]any, sb *strings.Builder) {
	sb.WriteString("{")
	i := 0
	for k, v := range options {
		sb.WriteString(k + ":")
		if mp, ok := v.(map[string]any); ok {
			stringifyQueryOptions(mp, sb)
		} else {
			json.NewEncoder(sb).Encode(v)
		}
		if i < len(options)-1 {
			sb.WriteString(",")
		}
		i++
	}
	sb.WriteString("}")
}

var supportedOperators = map[Operator]bool{
	OpAnd:      true,
	OpOr:       true,
	OpEq:       true,
	OpNeq:      true,
	OpGt:       true,
	OpGte:      true,
	OpLt:       true,
	OpLte:      true,
	OpIn:       true,
	OpNin:      true,
	OpLike:     true,
	OpNlike:    true,
	OpIsnull:   true,
	OpNotnull:  true,
	OpBetween:  true,
	OpNbetween: true,
}

// ParseFromRequest parses a query from a request
func ParseFromRequest(r *http.Request) (*Query, error) {
	if !r.URL.Query().Has("q") {
		return nil, errors.New("no query found in request")
	}

	rawQueryValue := r.URL.Query().Get("q")
	if len(rawQueryValue) == 0 {
		return nil, errors.New("empty query")
	}

	return ParseFromString(rawQueryValue)
}

// ParseFromString parses a query from a string
func ParseFromString(rawQuery string) (*Query, error) {
	query := &Query{}
	if len(rawQuery) == 0 {
		return query, nil
	}

	return query, nil
}

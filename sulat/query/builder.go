package query

func buildQuery(op Operator, field string, value any) *Query {
	return &Query{
		Field:    field,
		Operator: op,
		Value:    value,
	}
}

func Eq(field string, value any) *Query {
	return buildQuery(OpEq, field, value)
}

func Neq(field string, value any) *Query {
	return buildQuery(OpNeq, field, value)
}

func Gt(field string, value any) *Query {
	return buildQuery(OpGt, field, value)
}

func Gte(field string, value any) *Query {
	return buildQuery(OpGte, field, value)
}

func Lt(field string, value any) *Query {
	return buildQuery(OpLt, field, value)
}

func Lte(field string, value any) *Query {
	return buildQuery(OpLte, field, value)
}

func In(field string, value any) *Query {
	return buildQuery(OpIn, field, value)
}

func Nin(field string, value any) *Query {
	return buildQuery(OpNin, field, value)
}

func Like(field string, value any) *Query {
	return buildQuery(OpLike, field, value)
}

func Nlike(field string, value any) *Query {
	return buildQuery(OpNlike, field, value)
}

func Isnull(field string, value any) *Query {
	return buildQuery(OpIsnull, field, value)
}

func Notnull(field string, value any) *Query {
	return buildQuery(OpNotnull, field, value)
}

func Between(field string, value []any) *Query {
	return buildQuery(OpBetween, field, value)
}

func Nbetween(field string, value []any) *Query {
	return buildQuery(OpNbetween, field, value)
}

func And(queries ...*Query) *Query {
	return buildQuery(OpAnd, "", queries)
}

func Or(queries ...*Query) *Query {
	return buildQuery(OpOr, "", queries)
}

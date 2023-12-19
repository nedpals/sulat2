package query

import (
	"encoding/json"
	"fmt"
	"strings"
	"text/scanner"

	"golang.org/x/exp/slices"
)

var isParserDebug = false

type Token struct {
	Text     string
	Raw      rune
	Position scanner.Position
}

func (t *Token) debug(labels ...string) {
	if !isParserDebug {
		return
	}

	label := "Token.debug"
	if len(labels) > 0 {
		label = labels[0]
	}

	if t == nil {
		fmt.Printf("[%s] Token: nil\n", label)
		return
	}
	fmt.Printf("[%s] Token: %q | %s\n", label, t.Text, t.Position.String())
}

type Parser struct {
	IsTest        bool
	parentQuery   *Query
	prevToken     *Token
	currentToken  *Token
	nextToken     *Token
	nextNextToken *Token
	sc            *scanner.Scanner
}

func NewParser() *Parser {
	p := &Parser{
		sc: &scanner.Scanner{},
	}
	return p
}

func (p *Parser) Next() *Token {
	p.prevToken = p.currentToken
	p.currentToken = p.nextToken
	p.nextToken = p.nextNextToken

	if p.sc.Peek() != scanner.EOF {
		nextTok := p.sc.Scan()
		p.nextNextToken = &Token{
			Text:     p.sc.TokenText(),
			Raw:      nextTok,
			Position: p.sc.Pos(),
		}
	} else {
		p.nextNextToken = nil
	}

	if p.currentToken != nil {
		p.log("currentToken: %s|skipnl", p.currentToken.Text)
		if p.nextToken != nil {
			p.log(" | nextToken: %s|skipnl", p.nextToken.Text)
		}
		p.log("\n")
	}

	return p.currentToken
}

func (p *Parser) Parse(input string) (*Query, error) {
	isParserDebug = p.IsTest

	p.prevToken = nil
	p.currentToken = nil
	p.nextToken = nil

	p.sc.Init(strings.NewReader(input))
	p.sc.Mode = scanner.ScanIdents | scanner.ScanStrings | scanner.ScanFloats | scanner.ScanInts

	p.Next()
	p.Next()

	defer p.log("===== end of parsing =====")
	return p.parseQuery()
}

func (p *Parser) log(f string, a ...any) {
	if !p.IsTest {
		return
	} else if !strings.HasSuffix(f, "\n") && !strings.HasSuffix(f, "|skipnl") {
		f = f + "\n"
	}
	f = strings.TrimSuffix(f, "|skipnl")
	fmt.Printf(f, a...)
}

var validQueryValues = []rune{'[', '{', scanner.String, scanner.Int, scanner.Float}

func (p *Parser) parseQuery() (*Query, error) {
	token := p.Next()
	if token.Raw != scanner.Ident {
		return nil, fmt.Errorf("expected operator or field, got %s", token.Text)
	}

	query := &Query{
		Operator: Operator(token.Text),
		Options:  nil,
	}

	field := p.currentToken.Text
	// support dot notation-like field
	for p.nextToken.Raw == '.' || (len(p.nextToken.Text) > 1 && p.nextToken.Text[0] == '.') {
		tok := p.Next() // .
		if len(tok.Text) > 1 {
			field += tok.Text
		} else {
			p.Next() // field
			field += "." + p.currentToken.Text
		}
	}

	if p.nextToken.Raw == '(' {
		// Parse nested query entries
		queries, err := p.parseNestedQueries(query)
		if err != nil {
			return nil, err
		}

		if len(queries) > 0 {
			query.Value = queries
		}
	} else if slices.Contains(validQueryValues, p.nextToken.Raw) {
		// Parse query entry
		query = p.parentQuery
		value, err := p.parseJSONValue()
		if err != nil {
			return nil, err
		}

		query.Field = field
		query.Value = value

		return query, nil
	} else if _, ok := supportedOperators[Operator(token.Text)]; !ok {
		return nil, fmt.Errorf("unsupported operator: %s", token.Text)
	} else {
		return nil, fmt.Errorf("expected '(' or ',', got %s", token.Text)
	}

	if p.nextToken != nil {
		if p.nextToken.Raw == ',' && p.nextNextToken.Raw == '{' {
			p.Next() // ,
			options, err := p.parseQueryOptions()
			if err != nil {
				return nil, err
			}
			query.Options = options
		}
	}

	token = p.Next()
	token.debug("expected end parseQuery")

	if token.Raw != ')' {
		return nil, expectedError(')', token)
	}

	return query, nil
}

func (p *Parser) parseNestedQueries(parent *Query) ([]*Query, error) {
	oldParentQuery := p.parentQuery

	defer func() {
		p.log("going back to previous parent...")
		p.parentQuery = oldParentQuery
	}()

	p.parentQuery = parent
	var queries []*Query

	token := p.Next()
	if token.Raw != '(' {
		return nil, expectedError('(', token)
	}

	for p.nextToken.Raw != ')' {
		query, err := p.parseQuery()
		if err != nil {
			return nil, err
		}
		if query != parent {
			queries = append(queries, query)
		}
		if p.nextToken == nil || p.nextToken.Raw != ',' || p.nextToken.Raw == ')' {
			break
		} else if p.nextNextToken.Raw == '{' {
			p.log("break!")
			break
		}
		p.Next() // move to next token
		p.currentToken.debug("loop inside parseNestedQueries")
	}

	// parse query options in the parent query
	if p.nextToken != nil {
		if p.nextToken.Raw == ',' && p.nextNextToken.Raw == '{' {
			p.log("end of nested queries, expecting parsing options")
			return queries, nil
		} else if p.nextToken.Raw != ')' {
			return nil, expectedError(')', token)
		}
	}

	p.log("end of nested queries")

	if len(queries) == 0 {
		return nil, nil
	}
	return queries, nil
}

func (p *Parser) parseQueryOptions() (map[string]any, error) {
	options := make(map[string]any)

	token := p.Next()
	if token.Raw != '{' {
		return nil, expectedError('{', token)
	}

	i := 0

	for p.nextToken.Raw != '}' {
		token = p.Next()
		token.debug("loop inside parseQueryOptions")

		if i > 0 {
			if token.Raw != ',' {
				return nil, expectedError(',', token)
			}
			token = p.Next()
		}

		if token.Raw != scanner.Ident {
			return nil, fmt.Errorf("expected option key, got %s", token.Text)
		}

		key := token.Text

		token = p.Next()
		if token.Raw != ':' {
			return nil, expectedError(':', token)
		}

		value, err := p.parseJSONValue()
		if err != nil {
			return nil, err
		}

		options[key] = value

		i++
	}

	token = p.Next()
	if token.Raw != '}' {
		return nil, expectedError('}', token)
	}

	return options, nil
}

func (p *Parser) parseJSONValue() (any, error) {
	token := p.nextToken
	if p.nextToken.Raw != '{' {
		p.Next()
	}

	p.log("parseJSONValue %s", token.Text)

	switch token.Raw {
	case scanner.String:
		value := token.Text
		// Remove the surrounding quotes from the value
		value = value[1 : len(value)-1]
		return value, nil
	case scanner.Int, scanner.Float:
		var value json.Number
		if err := json.NewDecoder(strings.NewReader(token.Text)).Decode(&value); err != nil {
			return nil, err
		}
		return value, nil
	case '{':
		// Parse object same as parseQueryOptions
		return p.parseQueryOptions()
	case '[':
		// Parse JSON array
		accumulateText := &strings.Builder{}
		for token.Raw != ']' {
			accumulateText.WriteString(token.Text)
			token = p.Next()

			if token.Raw == scanner.EOF {
				return nil, fmt.Errorf("unexpected end of input")
			} else if token.Raw == ']' {
				accumulateText.WriteString(token.Text)
				break
			}
		}
		var arr []any
		if err := json.NewDecoder(strings.NewReader(accumulateText.String())).Decode(&arr); err != nil {
			return nil, err
		}
		return arr, nil
	default:
		return nil, fmt.Errorf("invalid JSON value: %s", token.Text)
	}
}

func expectedError(exp rune, tok *Token) error {
	return fmt.Errorf("expected %q, got %s", exp, tok.Text)
}

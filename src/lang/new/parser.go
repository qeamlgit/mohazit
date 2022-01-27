package new

import (
	"fmt"
	"mohazit/lang"
	"mohazit/lib"
	"strconv"
	"strings"
)

var (
	notIdent = lib.LazyError("parser: unexpected token (got %s, want ident)", "npar_notident")
)

type Statement struct {
	Keyword string
	Args    []*Token
}

type Parser struct {
	lexer *Lexer
}

func NewParser() *Parser {
	return &Parser{NewLexer()}
}

func (p *Parser) Source(src string) {
	p.lexer.Source(src)
}

func (p *Parser) tokens() ([]*Token, error) {
	out := []*Token{}
	var t *Token
	var err error
	for p.lexer.Has() {
		t, err = p.lexer.Next()
		if err != nil {
			return nil, err
		}
		if t == nil {
			continue
		}
		out = append(out, t)
		if t.Type == tLinefeed {
			return out, nil
		}
	}
	return p.TrimSpace(out), nil
}

func (p *Parser) Next() (*Statement, error) {
	raw, err := p.tokens()
	if err != nil {
		return nil, err
	}
	if len(raw) < 1 {
		return nil, nil
	}
	kwToken := raw[0]
	if kwToken.Type != tIdent {
		return nil, notIdent.Get(fmt.Sprint(kwToken.Type))
	}
	kw := strings.ToLower(kwToken.Raw)
	args := []*Token{}
	for i := 1; i < len(raw); i++ {
		args = append(args, raw[i])
	}
	return &Statement{kw, args}, nil
}

func (p *Parser) Args(tkns []*Token) ([]*lang.Object, error) {
	out := []*lang.Object{}
	ctx := ""
	for _, tkn := range tkns {
		switch tkn.Type {
		case tOper:
			if tkn.Raw == "\\" {
				out = append(out, lang.NewStr(strings.TrimSpace(ctx)))
				ctx = ""
			} else {
				ctx += tkn.Raw
			}
		case tLiteral:
			out = append(out, lang.NewStr(strings.TrimSpace(ctx)))
			ctx = ""
			v, err := strconv.Atoi(tkn.Raw)
			if err != nil {
				return nil, err
			}
			out = append(out, lang.NewInt(v))
		default:
			ctx += tkn.Raw
		}
	}
	if ctx != "" {
		out = append(out, lang.NewStr(strings.TrimSpace(ctx)))
		ctx = ""
	}
	return out, nil
}

func (p *Parser) Tokens2object(t []*Token) (*lang.Object, error) {
	t = p.TrimSpace(t)
	switch t[0].Type {
	case tIdent, tInvalid:
		v := lang.NewStr(t[0].Raw)
		for i := 0; i < len(t); i++ {
			tkn := t[i]
			switch tkn.Type {
			case tIdent, tInvalid, tSpace:
				v.StrV += tkn.Raw
			default:
				return lang.NewNil(), unexTkn.Get(tkn.Type.String())
			}
		}
		return v, nil
	case tLiteral:
		v, err := strconv.Atoi(t[0].Raw)
		return lang.NewInt(v), err
	default:
		return lang.NewNil(), unexTkn.Get(t[0].Type.String())
	}
}

func (p *Parser) TrimSpace(t []*Token) []*Token {
	ltrim := []*Token{}
	ignore := true
	for _, tkn := range t {
		if tkn.Type != tSpace && ignore {
			ignore = false
		}
		if !ignore {
			ltrim = append(ltrim, tkn)
		}
	}
	rtrim := []*Token{}
	ignore = true
	for i := len(ltrim) - 1; i >= 0; i-- {
		if ltrim[i].Type != tSpace && ignore {
			ignore = false
		}
		if !ignore {
			rtrim = append([]*Token{ltrim[i]}, rtrim...)
		}
	}
	return rtrim
}

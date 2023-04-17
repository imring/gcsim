package parse

import (
	"errors"
	"fmt"
	"runtime"
	"strconv"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character/profile"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

type Parser struct {
	lex *lexer
	res *ast.ActionList

	//other information tracked as we parse
	chars          map[keys.Char]*profile.CharacterProfile
	charOrder      []keys.Char
	currentCharKey keys.Char

	//lookahead
	token []ast.Token
	pos   int

	//parseFn
	prefixParseFns map[ast.TokenType]func() (ast.Expr, error)
	infixParseFns  map[ast.TokenType]func(ast.Expr) (ast.Expr, error)
}

type parseFn func(*Parser) (parseFn, error)

func New(input string) *Parser {
	p := &Parser{
		chars:          make(map[keys.Char]*profile.CharacterProfile),
		prefixParseFns: make(map[ast.TokenType]func() (ast.Expr, error)),
		infixParseFns:  make(map[ast.TokenType]func(ast.Expr) (ast.Expr, error)),
		token:          make([]ast.Token, 0, 20),
		pos:            -1,
	}
	p.lex = lex(input)
	p.res = &ast.ActionList{
		Program: ast.NewBlockStmt(0),
		Settings: ast.SimulatorSettings{
			EnableHitlag:    true, // default hitlag enabled
			DefHalt:         true, //default defhalt to true
			NumberOfWorkers: 20,   //default 20 workers if none set
			Iterations:      1000, //default 1000 iterations
			Delays: player.Delays{
				Swap: 1, //default swap timer of 1
			},
		},
		PlayerPos: core.Coord{
			R: 0.3, //default player radius 0.3, pos 0,0
		},
	}
	//expr functions
	p.prefixParseFns[ast.ItemIdentifier] = p.parseIdent
	p.prefixParseFns[ast.ItemField] = p.parseField
	p.prefixParseFns[ast.ItemNumber] = p.parseNumber
	p.prefixParseFns[ast.ItemBool] = p.parseBool
	p.prefixParseFns[ast.ItemString] = p.parseString
	p.prefixParseFns[ast.KeywordFn] = p.parseFnLit
	p.prefixParseFns[ast.LogicNot] = p.parseUnaryExpr
	p.prefixParseFns[ast.ItemMinus] = p.parseUnaryExpr
	p.prefixParseFns[ast.ItemLeftParen] = p.parseParen
	p.prefixParseFns[ast.ItemLeftSquareParen] = p.parseMap
	p.infixParseFns[ast.LogicAnd] = p.parseBinaryExpr
	p.infixParseFns[ast.LogicOr] = p.parseBinaryExpr
	p.infixParseFns[ast.ItemPlus] = p.parseBinaryExpr
	p.infixParseFns[ast.ItemMinus] = p.parseBinaryExpr
	p.infixParseFns[ast.ItemForwardSlash] = p.parseBinaryExpr
	p.infixParseFns[ast.ItemAsterisk] = p.parseBinaryExpr
	p.infixParseFns[ast.OpEqual] = p.parseBinaryExpr
	p.infixParseFns[ast.OpNotEqual] = p.parseBinaryExpr
	p.infixParseFns[ast.OpLessThan] = p.parseBinaryExpr
	p.infixParseFns[ast.OpLessThanOrEqual] = p.parseBinaryExpr
	p.infixParseFns[ast.OpGreaterThan] = p.parseBinaryExpr
	p.infixParseFns[ast.OpGreaterThanOrEqual] = p.parseBinaryExpr
	p.infixParseFns[ast.ItemLeftParen] = p.parseCall
	return p
}

// consume returns err if next token does not match expected
// otherwise return next token and nil error
func (p *Parser) consume(i ast.TokenType) (ast.Token, error) {
	n := p.next()
	if n.Typ != i {
		return n, errors.New("unexpected token")
	}
	return n, nil
}

// next returns the next token.
func (p *Parser) next() ast.Token {
	p.pos++
	if p.pos == len(p.token) {
		//grab more from the stream
		n := p.lex.nextItem()
		p.token = append(p.token, n)
	}
	return p.token[p.pos]
}

// backup backs the input stream up one token.
func (p *Parser) backup() {
	p.pos--
	//no op if at beginning
	if p.pos < -1 {
		p.pos = -1
	}
}

// peek returns but does not consume the next token.
func (p *Parser) peek() ast.Token {
	n := p.next()
	p.backup()
	return n
}

func (p *Parser) acceptSeqReturnLast(items ...ast.TokenType) (ast.Token, error) {
	var n ast.Token
	for _, v := range items {
		n = p.next()
		if n.Typ != v {
			_, file, no, _ := runtime.Caller(1)
			return n, fmt.Errorf("(%s#%d) expecting %v, got token %v", file, no, v, n)
		}
	}
	return n, nil
}

func itemNumberToInt(i ast.Token) (int, error) {
	r, err := strconv.Atoi(i.Val)
	return int(r), err
}

func itemNumberToFloat64(i ast.Token) (float64, error) {
	r, err := strconv.ParseFloat(i.Val, 64)
	return r, err
}

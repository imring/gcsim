package ast

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/shortcut"
)

type precedence int

const (
	_ precedence = iota
	Lowest
	LogicalOr
	LogicalAnd // TODO: or make one for && and ||?
	Equals
	LessOrGreater
	Sum
	Product
	Prefix
	Call
)

var precedences = map[TokenType]precedence{
	LogicOr:              LogicalOr,
	LogicAnd:             LogicalAnd,
	OpEqual:              Equals,
	OpNotEqual:           Equals,
	OpLessThan:           LessOrGreater,
	OpGreaterThan:        LessOrGreater,
	OpLessThanOrEqual:    LessOrGreater,
	OpGreaterThanOrEqual: LessOrGreater,
	ItemPlus:             Sum,
	ItemMinus:            Sum,
	ItemForwardSlash:     Product,
	ItemAsterisk:         Product,
	ItemLeftParen:        Call,
}

func (t Token) precedence() precedence {
	if p, ok := precedences[t.Typ]; ok {
		return p
	}
	return Lowest
}

// Parse returns the ActionList and any error that prevents the ActionList from being parsed
func (p *Parser) Parse() (*ActionList, error) {
	var err error
	for state := parseRows; state != nil; {
		state, err = state(p)
		if err != nil {
			return nil, err
		}
	}

	//sanity checks
	if len(p.charOrder) > 4 {
		p.res.Errors = append(p.res.Errors, fmt.Errorf("config contains a total of %v characters; cannot exceed 4", len(p.charOrder)))
	}

	if p.res.InitialChar == keys.NoChar {
		p.res.Errors = append(p.res.Errors, errors.New("config does not contain active char"))
	}

	initialCharFound := false
	for _, v := range p.charOrder {
		p.res.Characters = append(p.res.Characters, *p.chars[v])
		//check if active is part of the team
		if v == p.res.InitialChar {
			initialCharFound = true
		}
		//check number of set
		count := 0
		for _, c := range p.chars[v].Sets {
			count += c
		}
		if count > 5 {
			p.res.Errors = append(p.res.Errors, fmt.Errorf("character %v have more than 5 total set items", v.String()))
		}
	}

	if !initialCharFound {
		p.res.Errors = append(p.res.Errors, fmt.Errorf("active char %v not found in team", p.res.InitialChar))
	}

	if len(p.res.Targets) == 0 {
		p.res.Errors = append(p.res.Errors, errors.New("config does not contain any targets"))
	}

	//set some sane defaults; leave pos default to 0,0
	for i := range p.res.Targets {
		if p.res.Targets[i].Pos.R == 0 {
			p.res.Targets[i].Pos.R = 1
		}
	}

	//check all targets have hp if damage mode
	if p.res.Settings.DamageMode {
		for i, v := range p.res.Targets {
			if v.HP == 0 {
				p.res.Errors = append(p.res.Errors, fmt.Errorf("damage mode is activated; target #%v does not have hp set", i+1))
			}
		}
	}

	//build the err msgs
	p.res.ErrorMsgs = make([]string, 0, len(p.res.Errors))
	for _, v := range p.res.Errors {
		p.res.ErrorMsgs = append(p.res.ErrorMsgs, v.Error())
	}

	return p.res, nil
}

func parseRows(p *Parser) (parseFn, error) {
	switch n := p.peek(); n.Typ {
	case ItemCharacterKey:
		p.next()
		//check if this is character stats etc or an action
		if p.peek().Typ != ItemActionKey {
			//not an ActionStmt
			//set up char and set key
			key, ok := shortcut.CharNameToKey[n.Val]
			if !ok {
				//this would never happen
				return nil, fmt.Errorf("ln%v: unexpected error; invalid char key %v", n.line, n.Val)
			}
			if _, ok := p.chars[key]; !ok {
				p.newChar(key)
			}
			p.currentCharKey = key
			return parseCharacter, nil
		}
		p.backup()
		//parse action item
		// return parseProgram, nil
		node, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		p.res.Program.append(node)
		return parseRows, nil
	case KeywordActive:
		p.next()
		//next should be char then end line
		char, err := p.consume(ItemCharacterKey)
		if err != nil {
			return nil, fmt.Errorf("ln%v: setting active char: invalid char %v", char.line, char.Val)
		}
		p.res.InitialChar = shortcut.CharNameToKey[char.Val]
		n, err := p.consume(ItemTerminateLine)
		if err != nil {
			return nil, fmt.Errorf("ln%v: expecting ; after active <char>, got %v", n.line, n.Val)
		}
		return parseRows, nil
	case KeywordTarget:
		p.next()
		return parseTarget, nil
	case KeywordEnergy:
		p.next()
		return parseEnergy, nil
	case KeywordOptions:
		p.next()
		return parseOptions, nil
	case ItemEOF:
		return nil, nil
	default: //default should be look for gcsl
		node, err := p.parseStatement()
		p.res.Program.append(node)
		if err != nil {
			return nil, err
		}
		return parseRows, nil
	}
}

func (p *Parser) parseStatement() (Node, error) {
	//some statements end in semi, other don't
	hasSemi := true
	stmtType := ""
	var node Node
	var err error
	switch n := p.peek(); n.Typ {
	case KeywordBreak:
		fallthrough
	case KeywordFallthrough:
		fallthrough
	case KeywordContinue:
		stmtType = "continue"
		node, err = p.parseCtrl()
	case KeywordLet:
		stmtType = "let"
		node, err = p.parseLet()
	case ItemCharacterKey:
		stmtType = "char action"
		node, err = p.parseAction()
	case KeywordReturn:
		stmtType = "return"
		node, err = p.parseReturn()
	case KeywordIf:
		node, err = p.parseIf()
		hasSemi = false
	case KeywordSwitch:
		node, err = p.parseSwitch()
		hasSemi = false
	case KeywordFn:
		node, err = p.parseFn(true)
		hasSemi = false
	case KeywordWhile:
		node, err = p.parseWhile()
		hasSemi = false
	case KeywordFor:
		node, err = p.parseFor()
		hasSemi = false
	case ItemLeftBrace:
		node, err = p.parseBlock()
		hasSemi = false
	case ItemIdentifier:
		p.next()
		//check if = after
		if x := p.peek(); x.Typ == ItemAssign {
			p.backup()
			node, err = p.parseAssign()
			break
		}
		//it's an expr if no assign
		p.backup()
		fallthrough
	default:
		node, err = p.parseExpr(Lowest)
	}
	//check if any of the parse error'd
	if err != nil {
		return node, err
	}
	//check for semi
	if hasSemi {
		n, err := p.consume(ItemTerminateLine)
		if err != nil {
			return nil, fmt.Errorf("ln%v: expecting ; at end of %v statement, got %v", n.line, stmtType, n.Val)
		}
	}
	return node, nil
}

// excepting let ident = expr;
func (p *Parser) parseLet() (Stmt, error) {
	n := p.next()

	ident, err := p.consume(ItemIdentifier)
	if err != nil {
		//next token not an identifier
		return nil, fmt.Errorf("ln%v: expecting identifier after let, got %v", ident.line, ident.Val)
	}

	a, err := p.consume(ItemAssign)
	if err != nil {
		//next token not and identifier
		return nil, fmt.Errorf("ln%v: expecting = after identifier in let statement, got %v", a.line, a.Val)
	}

	expr, err := p.parseExpr(Lowest)

	stmt := &LetStmt{
		Pos:   n.pos,
		Ident: ident,
		Val:   expr,
	}

	return stmt, err
}

// expecting ident = expr
func (p *Parser) parseAssign() (Stmt, error) {

	ident, err := p.consume(ItemIdentifier)
	if err != nil {
		//next token not and identifier
		return nil, fmt.Errorf("ln%v: expecting identifier in assign statement, got %v", ident.line, ident.Val)
	}

	a, err := p.consume(ItemAssign)
	if err != nil {
		//next token not and identifier
		return nil, fmt.Errorf("ln%v: expecting = after identifier in assign statement, got %v", a.line, a.Val)
	}

	expr, err := p.parseExpr(Lowest)

	if err != nil {
		return nil, err
	}

	stmt := &AssignStmt{
		Pos:   ident.pos,
		Ident: ident,
		Val:   expr,
	}

	return stmt, nil

}

func (p *Parser) parseIf() (Stmt, error) {
	n := p.next()

	stmt := &IfStmt{
		Pos: n.pos,
	}

	var err error

	stmt.Condition, err = p.parseExpr(Lowest)
	if err != nil {
		return nil, err
	}

	//expecting a { next
	if n := p.peek(); n.Typ != ItemLeftBrace {
		return nil, fmt.Errorf("ln%v: expecting { after if, got %v", n.line, n.Val)
	}

	stmt.IfBlock, err = p.parseBlock() //parse block here
	if err != nil {
		return nil, err
	}

	//stop if no else
	if n := p.peek(); n.Typ != KeywordElse {
		return stmt, nil
	}

	//skip the else keyword
	p.next()

	//expecting another stmt (should be either if or block)
	block, err := p.parseStatement()
	switch block.(type) {
	case *IfStmt, *BlockStmt:
	default:
		stmt.ElseBlock = nil
		return stmt, fmt.Errorf("ln%v: expecting either if or normal block after else", n.line)
	}

	stmt.ElseBlock = block.(Stmt)

	return stmt, err
}

func (p *Parser) parseSwitch() (Stmt, error) {

	//switch expr { }
	n, err := p.consume(KeywordSwitch)
	if err != nil {
		panic("unreachable")
	}

	stmt := &SwitchStmt{
		Pos: n.pos,
	}

	//condition can be optional; if next item is itemLeftBrace then simply set condition to 1
	if n := p.peek(); n.Typ != ItemLeftBrace {
		stmt.Condition, err = p.parseExpr(Lowest)
		if err != nil {
			return nil, err
		}
	} else {
		stmt.Condition = nil
	}

	if n := p.next(); n.Typ != ItemLeftBrace {
		return nil, fmt.Errorf("ln%v: expecting { after switch, got %v", n.line, n.Val)
	}

	//look for cases while not }
	for n := p.next(); n.Typ != ItemRightBrace; n = p.next() {
		var err error
		//expecting case expr: block
		switch n.Typ {
		case KeywordCase:
			cs := &CaseStmt{
				Pos: n.pos,
			}
			cs.Condition, err = p.parseExpr(Lowest)
			if err != nil {
				return nil, err
			}
			//colon, then read until we hit next case
			if n := p.peek(); n.Typ != ItemColon {
				return nil, fmt.Errorf("ln%v: expecting : after case, got %v", n.line, n.Val)
			}
			cs.Body, err = p.parseCaseBody()
			if err != nil {
				return nil, err
			}
			stmt.Cases = append(stmt.Cases, cs)
		case KeywordDefault:
			//colon, then read until we hit next case
			if p.peek().Typ != ItemColon {
				return nil, fmt.Errorf("ln%v: expecting : after default, got %v", n.line, n.Val)
			}
			stmt.Default, err = p.parseCaseBody()
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("ln%v: expecting case or default token, got %v", n.line, n.Val)
		}

	}

	return stmt, nil
}

func (p *Parser) parseCaseBody() (*BlockStmt, error) {
	n := p.next() //start with :
	block := newBlockStmt(n.pos)
	var node Node
	var err error
	//parse line by line until we hit }
	for {
		//make sure we don't get any illegal lines
		switch n := p.peek(); n.Typ {
		case ItemCharacterKey:
			if !p.peekValidCharAction() {
				n := p.next()
				n = p.next()
				return nil, fmt.Errorf("ln%v: expecting action after character token, got %v", n.line, n.Val)
			}
		case KeywordDefault:
			fallthrough
		case KeywordCase:
			fallthrough
		case ItemRightBrace:
			return block, nil
		case ItemEOF:
			return nil, fmt.Errorf("reached end of file without closing }")
		}
		//parse statement here
		node, err = p.parseStatement()
		if err != nil {
			return nil, err
		}
		block.append(node)
	}
}

// while { }
func (p *Parser) parseWhile() (Stmt, error) {
	n := p.next()

	stmt := &WhileStmt{
		Pos: n.pos,
	}

	var err error

	stmt.Condition, err = p.parseExpr(Lowest)
	if err != nil {
		return nil, err
	}

	//expecting a { next
	if n := p.peek(); n.Typ != ItemLeftBrace {
		return nil, fmt.Errorf("ln%v: expecting { after while, got %v", n.line, n.Val)
	}

	stmt.WhileBlock, err = p.parseBlock() //parse block here

	return stmt, err
}

// for <init ;> <cond> <; post> { <body> }
// for { <body > }
func (p *Parser) existVarDecl() bool {
	switch n := p.peek(); n.Typ {
	case KeywordLet:
		return true
	case ItemIdentifier:
		p.next()
		b := p.peek().Typ == ItemAssign
		p.backup()
		return b
	}
	return false
}

func (p *Parser) parseFor() (Stmt, error) {
	n := p.next()

	stmt := &ForStmt{
		Pos: n.pos,
	}

	var err error

	if n := p.peek(); n.Typ == ItemLeftBrace {
		stmt.Body, err = p.parseBlock() //parse block here
		return stmt, err
	}

	//init
	if p.existVarDecl() {
		if n := p.peek(); n.Typ == KeywordLet {
			stmt.Init, err = p.parseLet()
		} else {
			stmt.Init, err = p.parseAssign()
		}
		if err != nil {
			return nil, err
		}

		if n := p.peek(); n.Typ != ItemTerminateLine {
			return nil, fmt.Errorf("ln%v: expecting ; after statement, got %v", n.line, n.Val)
		}
		p.next() //skip ;
	}

	//cond
	stmt.Cond, err = p.parseExpr(Lowest)
	if err != nil {
		return nil, err
	}

	//post
	if n := p.peek(); n.Typ == ItemTerminateLine {
		p.next() //skip ;
		if n := p.peek(); n.Typ != ItemLeftBrace {
			stmt.Post, err = p.parseAssign()
			if err != nil {
				return nil, err
			}
		}
	}

	//expecting a { next
	if n := p.peek(); n.Typ != ItemLeftBrace {
		return nil, fmt.Errorf("ln%v: expecting { after for, got %v", n.line, n.Val)
	}

	stmt.Body, err = p.parseBlock() //parse block here

	return stmt, err
}

func (p *Parser) parseFn(indent bool) (Stmt, error) {
	//fn ident(...ident){ block }
	//consume fn
	n := p.next()
	stmt := &FnStmt{
		Pos: n.pos,
	}

	var err error
	if indent {
		//ident next
		n, err := p.consume(ItemIdentifier)
		if err != nil {
			return nil, fmt.Errorf("ln%v: expecting identifier after fn, got %v", n.line, n.Val)
		}
		stmt.FunVal = n
	}

	if l := p.peek(); l.Typ != ItemLeftParen {
		return nil, fmt.Errorf("ln%v: expecting ( after identifier, got %v", l.line, l.Val)
	}

	stmt.Args, err = p.parseFnArgs()
	if err != nil {
		return nil, err
	}
	stmt.Body, err = p.parseBlock()
	if err != nil {
		return nil, err
	}

	//check that args are not duplicates
	chk := make(map[string]bool)
	for _, v := range stmt.Args {
		if _, ok := chk[v.Value]; ok {
			return nil, fmt.Errorf("fn %v contains duplicated param name %v", stmt.FunVal.Val, v.Value)
		}
		chk[v.Value] = true
	}

	return stmt, nil
}

func (p *Parser) parseFnArgs() ([]*Ident, error) {
	//consume (
	var args []*Ident
	p.next()
	for n := p.next(); n.Typ != ItemRightParen; n = p.next() {
		a := &Ident{}
		//expecting ident, comma
		if n.Typ != ItemIdentifier {
			return nil, fmt.Errorf("ln%v: expecting identifier in param list, got %v", n.line, n.Val)
		}
		a.Pos = n.pos
		a.Value = n.Val

		args = append(args, a)

		//if next token is a comma, then there should be another ident after that
		//otherwise we have a problem
		if l := p.peek(); l.Typ == ItemComma {
			p.next() //consume the comma
			if l = p.peek(); l.Typ != ItemIdentifier {
				return nil, fmt.Errorf("ln%v: expecting another identifier after comma in param list, got %v", n.line, n.Val)
			}
		}
	}
	return args, nil
}

func (p *Parser) parseReturn() (Stmt, error) {
	n := p.next() //return
	stmt := &ReturnStmt{
		Pos: n.pos,
	}
	var err error
	stmt.Val, err = p.parseExpr(Lowest)
	return stmt, err
}

func (p *Parser) parseCtrl() (Stmt, error) {
	n := p.next()
	stmt := &CtrlStmt{
		Pos: n.pos,
	}
	switch n.Typ {
	case KeywordBreak:
		stmt.Typ = CtrlBreak
	case KeywordContinue:
		stmt.Typ = CtrlContinue
	case KeywordFallthrough:
		stmt.Typ = CtrlFallthrough
	default:
		return nil, fmt.Errorf("ln%v: expecting ctrl token, got %v", n.line, n.Val)
	}
	return stmt, nil
}

func (p *Parser) parseCall(fun Expr) (Expr, error) {
	//expecting (params)
	n, err := p.consume(ItemLeftParen)
	if err != nil {
		return nil, fmt.Errorf("expecting ( after ident, got %v", fun.String())
	}
	expr := &CallExpr{
		Pos: n.pos,
		Fun: fun,
	}
	expr.Args, err = p.parseCallArgs()

	return expr, err

}

func (p *Parser) parseCallArgs() ([]Expr, error) {
	var args []Expr

	if p.peek().Typ == ItemRightParen {
		//consume the right paren
		p.next()
		return args, nil
	}

	//next should be an expression
	exp, err := p.parseExpr(Lowest)
	if err != nil {
		return args, err
	}
	args = append(args, exp)

	for p.peek().Typ == ItemComma {
		p.next() //skip the comma
		exp, err = p.parseExpr(Lowest)
		if err != nil {
			return args, err
		}
		args = append(args, exp)
	}

	if n := p.next(); n.Typ != ItemRightParen {
		p.backup()
		return nil, fmt.Errorf("ln%v: expecting ) at end of function call, got: %v", n.line, n.pos)
	}

	return args, nil
}

// check if it's a valid character action, assuming current token is "character"
func (p *Parser) peekValidCharAction() bool {
	p.next()
	//check if this is character stats etc or an action
	if p.peek().Typ != ItemActionKey {
		p.backup()
		//not an ActionStmt
		return false
	}
	p.backup()
	return true
}

// parseBlock return a node contain and BlockStmt
func (p *Parser) parseBlock() (*BlockStmt, error) {
	//should be surronded by {}
	n, err := p.consume(ItemLeftBrace)
	if err != nil {
		return nil, fmt.Errorf("ln%v: expecting {, got %v", n.line, n.Val)
	}
	block := newBlockStmt(n.pos)
	var node Node
	//parse line by line until we hit }
	for {
		//make sure we don't get any illegal lines
		switch n := p.peek(); n.Typ {
		case ItemCharacterKey:
			if !p.peekValidCharAction() {
				n := p.next()
				n = p.next()
				return nil, fmt.Errorf("ln%v: expecting action after character token, got %v", n.line, n.Val)
			}
		case ItemRightBrace:
			p.next() //consume the braces
			return block, nil
		case ItemEOF:
			return nil, fmt.Errorf("reached end of file without closing }")
		}
		//parse statement here
		node, err = p.parseStatement()
		if err != nil {
			return nil, err
		}
		block.append(node)
	}

}
func (p *Parser) parseExpr(pre precedence) (Expr, error) {
	t := p.next()
	prefix := p.prefixParseFns[t.Typ]
	if prefix == nil {
		return nil, nil
	}
	p.backup()
	leftExp, err := prefix()
	if err != nil {
		return nil, err
	}

	for n := p.peek(); n.Typ != ItemTerminateLine && pre < n.precedence(); n = p.peek() {
		infix := p.infixParseFns[n.Typ]
		if infix == nil {
			return leftExp, nil
		}

		leftExp, err = infix(leftExp)
		if err != nil {
			return nil, err
		}
	}

	return leftExp, nil
}

// next is an identifier
func (p *Parser) parseIdent() (Expr, error) {
	n := p.next()
	return &Ident{Pos: n.pos, Value: n.Val}, nil
}

func (p *Parser) parseField() (Expr, error) {
	//next is field, keep parsing as long as it is still fields
	//then concat them all together
	n := p.next()
	fields := make([]string, 0, 5)
	for ; n.Typ == ItemField; n = p.next() {
		fields = append(fields, strings.Trim(n.Val, "."))
	}
	//we would have consumed one too many here
	p.backup()
	return &Field{Pos: n.pos, Value: fields}, nil
}

func (p *Parser) parseString() (Expr, error) {
	n := p.next()
	return &StringLit{Pos: n.pos, Value: n.Val}, nil
}

func (p *Parser) parseFnLit() (Expr, error) {
	n := p.peek()
	stmt, err := p.parseFn(false)
	if err != nil {
		return nil, err
	}

	f := stmt.(*FnStmt)
	return &FuncLit{
		Pos:  n.pos,
		Args: f.Args,
		Body: f.Body,
	}, nil
}

func (p *Parser) parseNumber() (Expr, error) {
	//string, int, float, or bool
	n := p.next()
	num := &NumberLit{Pos: n.pos}
	//try parse int, if not ok then try parse float
	iv, err := strconv.ParseInt(n.Val, 10, 64)
	if err == nil {
		num.IntVal = iv
		num.FloatVal = float64(iv)
	} else {
		fv, err := strconv.ParseFloat(n.Val, 64)
		if err != nil {
			return nil, fmt.Errorf("ln%v: cannot parse %v to number", n.line, n.Val)
		}
		num.IsFloat = true
		num.FloatVal = fv
	}
	return num, nil
}

func (p *Parser) parseBool() (Expr, error) {
	// bool is a number (true = 1, false = 0)
	n := p.next()
	num := &NumberLit{Pos: n.pos}
	switch n.Val {
	case "true":
		num.IntVal = 1
		num.FloatVal = 1
	case "false":
		num.IntVal = 0
		num.FloatVal = 0
	default:
		return nil, fmt.Errorf("ln%v: expecting boolean, got %v", n.line, n.Val)
	}
	return num, nil
}

func (p *Parser) parseUnaryExpr() (Expr, error) {
	n := p.next()
	switch n.Typ {
	case LogicNot:
	case ItemMinus:
	default:
		return nil, fmt.Errorf("ln%v: unrecognized unary operator %v", n.line, n.Val)
	}
	var err error
	expr := &UnaryExpr{
		Pos: n.pos,
		Op:  n,
	}
	expr.Right, err = p.parseExpr(Prefix)
	return expr, err
}

func (p *Parser) parseBinaryExpr(left Expr) (Expr, error) {
	n := p.next()
	expr := &BinaryExpr{
		Pos:  n.pos,
		Op:   n,
		Left: left,
	}
	pr := n.precedence()
	var err error
	expr.Right, err = p.parseExpr(pr)
	return expr, err
}

func (p *Parser) parseParen() (Expr, error) {
	//skip the paren
	p.next()

	exp, err := p.parseExpr(Lowest)
	if err != nil {
		return nil, err
	}

	if n := p.peek(); n.Typ != ItemRightParen {
		return nil, nil
	}
	p.next() // consume the right paren

	return exp, nil
}

func (p *Parser) parseMap() (Expr, error) {
	//skip the paren
	n := p.next()
	expr := &MapExpr{
		Pos:    n.pos,
		Fields: make(map[string]Expr),
	}

	if p.peek().Typ == ItemRightSquareParen { // empty map
		p.next()
		return expr, nil
	}

	//loop until we hit square paren
	for {
		//we're expecting ident = int
		i, err := p.consume(ItemIdentifier)
		if err != nil {
			return nil, fmt.Errorf("ln%v: expecting identifier in map expression, got %v", i.line, i.Val)
		}

		a, err := p.consume(ItemAssign)
		if err != nil {
			return nil, fmt.Errorf("ln%v: expecting = after identifier in map expression, got %v", a.line, a.Val)
		}

		e, err := p.parseExpr(Lowest)
		if err != nil {
			return nil, err
		}
		expr.Fields[i.Val] = e

		//if we hit ], return; if we hit , keep going, other wise error
		n := p.next()
		switch n.Typ {
		case ItemRightSquareParen:
			return expr, nil
		case ItemComma:
			//do nothing, keep going
		default:
			return nil, fmt.Errorf("ln%v: <action param> bad token %v", n.line, n)
		}
	}
}

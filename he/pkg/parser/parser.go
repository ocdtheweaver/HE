package parser

import (
	"fmt"
	"strconv"

	"github.com/user/he/pkg/ast"
	"github.com/user/he/pkg/lexer"
)

// ---------------------
// Parser Structure
// ---------------------
type Parser struct {
	lexer     *lexer.Lexer
	curToken  lexer.Token
	peekToken lexer.Token
	errors    []string
}

// ---------------------
// Parser Constructor
// ---------------------
func New(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l, errors: []string{}}
	p.nextToken()
	p.nextToken()
	return p
}

// ---------------------
// Token Helpers
// ---------------------
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.Next()
}

func (p *Parser) curTokenIs(t lexer.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t lexer.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t lexer.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) peekError(t lexer.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead at %d:%d",
		t, p.peekToken.Type, p.peekToken.Pos.Line, p.peekToken.Pos.Col)
	p.errors = append(p.errors, msg)
}

func (p *Parser) Errors() []string {
	return p.errors
}

// ---------------------
// Program
// ---------------------
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{
		Summons:    []*ast.SummonDecl{},
		Modules:    []*ast.ModuleDecl{},
		Objects:    []*ast.ObjectDecl{},
		Assets:     []*ast.AssetDecl{},
		Statements: []ast.Statement{},
	}

	for !p.curTokenIs(lexer.EOF) {
		switch {
		case p.curTokenIs(lexer.SUMMON):
			if s := p.parseSummon(); s != nil {
				program.Summons = append(program.Summons, s)
			}
		case p.curTokenIs(lexer.MODULE):
			if m := p.parseModule(); m != nil {
				program.Modules = append(program.Modules, m)
			}
		case p.curTokenIs(lexer.MAKE), p.curTokenIs(lexer.CREATE):
			if o := p.parseObject(); o != nil {
				program.Objects = append(program.Objects, o)
			}
		case p.curTokenIs(lexer.WITH):
			program.Assets = append(program.Assets, p.parseAssets()...)
		default:
			if stmt := p.parseStatement(); stmt != nil {
				program.Statements = append(program.Statements, stmt)
			} else {
				p.nextToken()
			}
		}
	}

	return program
}

// ---------------------
// Summon
// ---------------------
func (p *Parser) parseSummon() *ast.SummonDecl {
	pos := p.tokenPos()
	p.nextToken() // skip SUMMON

	if !p.curTokenIs(lexer.STRING) {
		p.errors = append(p.errors, "expected string after summon")
		return nil
	}
	path := p.curToken.Lit
	p.nextToken()

	alias := ""
	if p.curTokenIs(lexer.AS) || p.curTokenIs(lexer.NAMED) {
		p.nextToken()
		if p.curTokenIs(lexer.IDENT) {
			alias = p.curToken.Lit
			p.nextToken()
		}
	}

	return &ast.SummonDecl{Pos: pos, Path: path, As: alias}
}

// ---------------------
// Module
// ---------------------
func (p *Parser) parseModule() *ast.ModuleDecl {
	pos := p.tokenPos()
	p.nextToken() // skip MODULE

	if !p.curTokenIs(lexer.IDENT) {
		p.errors = append(p.errors, "expected module name")
		return nil
	}
	name := p.curToken.Lit
	p.nextToken()

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	module := &ast.ModuleDecl{
		Pos:        pos,
		Name:       name,
		Submodules: []*ast.ModuleDecl{},
		Functions:  []*ast.FunctionDecl{},
		Body:       []ast.Declaration{},
	}

	for !p.curTokenIs(lexer.RBRACE) && !p.curTokenIs(lexer.EOF) {
		switch {
		case p.curTokenIs(lexer.MODULE):
			if sub := p.parseModule(); sub != nil {
				module.Submodules = append(module.Submodules, sub)
			}
		case p.curTokenIs(lexer.FN):
			if fn := p.parseFunction(); fn != nil {
				module.Functions = append(module.Functions, fn)
			}
		default:
			p.nextToken()
		}
	}

	p.expectPeek(lexer.RBRACE)
	return module
}

// ---------------------
// Function
// ---------------------
func (p *Parser) parseFunction() *ast.FunctionDecl {
	pos := p.tokenPos()
	p.nextToken() // skip FN

	if !p.curTokenIs(lexer.IDENT) {
		p.errors = append(p.errors, "expected function name")
		return nil
	}
	name := p.curToken.Lit
	p.nextToken()

	params := []*ast.ParamDecl{}
	if p.curTokenIs(lexer.LPAREN) {
		params = p.parseParamList()
	}

	if !p.expectPeek(lexer.LBRACK) {
		return nil
	}
	body := p.parseBlock()
	return &ast.FunctionDecl{Pos: pos, Name: name, Params: params, Body: body}
}

func (p *Parser) parseParamList() []*ast.ParamDecl {
	params := []*ast.ParamDecl{}
	p.nextToken() // skip LPAREN

	for !p.curTokenIs(lexer.RPAREN) && !p.curTokenIs(lexer.EOF) {
		if !p.curTokenIs(lexer.IDENT) {
			p.errors = append(p.errors, "expected parameter name")
			return nil
		}
		param := &ast.ParamDecl{Name: p.curToken.Lit, Pos: p.tokenPos()}
		p.nextToken()
		if p.curTokenIs(lexer.COLON) {
			p.nextToken()
			if !p.curTokenIs(lexer.IDENT) {
				p.errors = append(p.errors, "expected type after colon")
				return nil
			}
			param.TypeName = p.curToken.Lit
			p.nextToken()
		}
		params = append(params, param)
		if p.curTokenIs(lexer.COMMA) {
			p.nextToken()
		}
	}

	p.expectPeek(lexer.RPAREN)
	return params
}

// ---------------------
// Object
// ---------------------
func (p *Parser) parseObject() *ast.ObjectDecl {
	pos := p.tokenPos()
	p.nextToken() // skip MAKE/CREATE

	if !p.curTokenIs(lexer.IDENT) {
		p.errors = append(p.errors, "expected object name")
		return nil
	}
	name := p.curToken.Lit
	p.nextToken()

	like := ""
	if p.curTokenIs(lexer.LIKE) {
		p.nextToken()
		if !p.curTokenIs(lexer.IDENT) {
			p.errors = append(p.errors, "expected identifier after 'like'")
			return nil
		}
		like = p.curToken.Lit
		p.nextToken()
	}

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	obj := &ast.ObjectDecl{Pos: pos, Name: name, Like: like}
	obj.Properties, obj.Abilities, obj.Reactions, obj.Memories, obj.NestedObjects = p.parseObjectBody()
	return obj
}

func (p *Parser) parseObjectBody() ([]*ast.PropertyDecl, []*ast.AbilityDecl, []*ast.ReactionDecl, []*ast.MemoryDecl, []*ast.ObjectDecl) {
	properties := []*ast.PropertyDecl{}
	abilities := []*ast.AbilityDecl{}
	reactions := []*ast.ReactionDecl{}
	memories := []*ast.MemoryDecl{}
	nested := []*ast.ObjectDecl{}

	for !p.curTokenIs(lexer.RBRACE) && !p.curTokenIs(lexer.EOF) {
		switch {
		case p.curTokenIs(lexer.IDENT):
			switch p.peekToken.Type {
			case lexer.HAS, lexer.OWNS, lexer.CARRIES:
				prop := p.parseProperty()
				if prop != nil {
					properties = append(properties, prop)
				}
			case lexer.CAN, lexer.KNOWS:
				if ab := p.parseAbility(); ab != nil {
					abilities = append(abilities, ab)
				}
			case lexer.REMEMBERS:
				if mem := p.parseMemory(); mem != nil {
					memories = append(memories, mem)
				}
			default:
				p.nextToken()
			}
		case p.curTokenIs(lexer.MAKE), p.curTokenIs(lexer.CREATE):
			if obj := p.parseObject(); obj != nil {
				nested = append(nested, obj)
			}
		case p.curTokenIs(lexer.WITH):
			p.parseAssets()
		default:
			p.nextToken()
		}
	}

	p.expectPeek(lexer.RBRACE)
	return properties, abilities, reactions, memories, nested
}

// ---------------------
// Property
// ---------------------
func (p *Parser) parseProperty() *ast.PropertyDecl {
	pos := p.tokenPos()
	name := p.curToken.Lit
	p.nextToken() // skip identifier
	p.nextToken() // skip IS / STARTS AS
	value := p.parseExpression()
	return &ast.PropertyDecl{Pos: pos, Name: name, Value: value}
}

// ---------------------
// Ability
// ---------------------
func (p *Parser) parseAbility() *ast.AbilityDecl {
	pos := p.tokenPos()
	name := p.curToken.Lit
	p.nextToken() // skip CAN / KNOWS

	var action ast.Expr
	if !p.curTokenIs(lexer.RBRACE) && !p.curTokenIs(lexer.EOF) && !p.isStartOfStatement() {
		action = p.parseExpression()
	}

	return &ast.AbilityDecl{Pos: pos, Name: name, Action: action}
}

// ---------------------
// Memory
// ---------------------
func (p *Parser) parseMemory() *ast.MemoryDecl {
	pos := p.tokenPos()
	p.nextToken() // skip REMEMBERS

	var content ast.Expr
	if !p.curTokenIs(lexer.RBRACE) && !p.curTokenIs(lexer.EOF) {
		content = p.parseExpression()
	}

	return &ast.MemoryDecl{Pos: pos, Content: content}
}

// ---------------------
// Statements
// ---------------------
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case lexer.PRINT, lexer.SAY:
		return p.parsePrint()
	case lexer.WAIT:
		return p.parseWait()
	case lexer.RETURN:
		return p.parseReturn()
	case lexer.IDENT:
		return p.parseAssignOrCall()
	default:
		p.nextToken()
		return nil
	}
}

func (p *Parser) parsePrint() *ast.PrintStmt {
	pos := p.tokenPos()
	p.nextToken()
	expr := p.parseExpression()
	return &ast.PrintStmt{Pos: pos, Expr: expr}
}

func (p *Parser) parseWait() *ast.WaitStmt {
	pos := p.tokenPos()
	p.nextToken()
	expr := p.parseExpression()
	if p.curTokenIs(lexer.SECONDS) {
		p.nextToken()
	}
	return &ast.WaitStmt{Pos: pos, Seconds: expr}
}

func (p *Parser) parseReturn() *ast.ReturnStmt {
	pos := p.tokenPos()
	p.nextToken()
	expr := p.parseExpression()
	return &ast.ReturnStmt{Pos: pos, Expr: expr}
}

func (p *Parser) parseAssignOrCall() ast.Statement {
	pos := p.tokenPos()
	ident := &ast.IdentifierExpr{Name: p.curToken.Lit, Pos: pos}
	p.nextToken()

	// Assignment: x = expr
	if p.curTokenIs(lexer.ASSIGN) {
		p.nextToken()
		value := p.parseExpression()
		return &ast.AssignStmt{Pos: pos, Name: ident, Value: value}
	}

	// Function call: ident(args)
	if p.curTokenIs(lexer.LPAREN) {
		call := p.parseCallExpression(ident)
		return &ast.ExprStmt{Pos: pos, Expr: call}
	}

	// Plain identifier statement
	return &ast.ExprStmt{Pos: pos, Expr: ident}
}

// ---------------------
// Expressions (Pratt Parser)
// ---------------------
const (
	_ int = iota
	LOWEST
	LOGICALOR
	LOGICALAND
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	EXPONENT
	PREFIX
	CALL
)

var precedences = map[lexer.TokenType]int{
	lexer.OR:     LOGICALOR,
	lexer.AND:    LOGICALAND,
	lexer.EQ:     EQUALS,
	lexer.NOT_EQ: EQUALS,
	lexer.LT:     LESSGREATER,
	lexer.GT:     LESSGREATER,
	lexer.LTE:    LESSGREATER,
	lexer.GTE:    LESSGREATER,
	lexer.PLUS:   SUM,
	lexer.MINUS:  SUM,
	lexer.MULT:   PRODUCT,
	lexer.DIV:    PRODUCT,
	lexer.POW:    EXPONENT,
	lexer.LPAREN: CALL,
}

func (p *Parser) parseExpression() ast.Expr {
	return p.parseExprPrecedence(LOWEST)
}

func (p *Parser) parseExprPrecedence(precedence int) ast.Expr {
	var left ast.Expr

	switch p.curToken.Type {
	case lexer.IDENT:
		left = p.parseIdentifier()
	case lexer.NUMBER:
		left = p.parseNumber()
	case lexer.STRING:
		left = p.parseString()
	case lexer.BOOLEAN:
		left = p.parseBoolean()
	case lexer.MINUS, lexer.NOT:
		left = p.parsePrefix()
	case lexer.LPAREN:
		left = p.parseGroupedExpression()
	case lexer.LBRACK:
		left = p.parseArrayLiteral()
	default:
		p.errors = append(p.errors, fmt.Sprintf("unexpected token in expression: %v", p.curToken))
		p.nextToken()
		return nil
	}

	for !p.peekTokenIs(lexer.SEMICOLON) && precedence < p.peekPrecedence() {
		switch p.peekToken.Type {
		case lexer.PLUS, lexer.MINUS, lexer.MULT, lexer.DIV, lexer.POW,
			lexer.EQ, lexer.NOT_EQ, lexer.LT, lexer.GT, lexer.LTE, lexer.GTE,
			lexer.AND, lexer.OR:
			p.nextToken()
			left = p.parseInfix(left)
		case lexer.LPAREN:
			p.nextToken()
			left = p.parseCallExpression(left)
		default:
			return left
		}
	}

	return left
}

// Prefix parsers
func (p *Parser) parseIdentifier() ast.Expr {
	expr := &ast.IdentifierExpr{Name: p.curToken.Lit, Pos: p.tokenPos()}
	p.nextToken()
	return expr
}

func (p *Parser) parseNumber() ast.Expr {
	value, err := strconv.ParseFloat(p.curToken.Lit, 64)
	if err != nil {
		p.errors = append(p.errors, fmt.Sprintf("invalid number: %v", p.curToken.Lit))
		return nil
	}
	expr := &ast.NumberExpr{Value: value, Pos: p.tokenPos()}
	p.nextToken()
	return expr
}

func (p *Parser) parseString() ast.Expr {
	expr := &ast.StringExpr{Value: p.curToken.Lit, Pos: p.tokenPos()}
	p.nextToken()
	return expr
}

func (p *Parser) parseBoolean() ast.Expr {
	val := false
	if p.curToken.Lit == "true" || p.curToken.Lit == "yes" {
		val = true
	}
	expr := &ast.BooleanExpr{Value: val, Pos: p.tokenPos()}
	p.nextToken()
	return expr
}

func (p *Parser) parsePrefix() ast.Expr {
	pos := p.tokenPos()
	operator := p.curToken.Lit
	p.nextToken()
	right := p.parseExprPrecedence(PREFIX)
	return &ast.PrefixExpr{Operator: operator, Right: right, Pos: pos}
}

func (p *Parser) parseGroupedExpression() ast.Expr {
	p.nextToken()
	exp := p.parseExpression()
	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parseArrayLiteral() ast.Expr {
	pos := p.tokenPos()
	elements := []ast.Expr{}
	p.nextToken()

	for !p.curTokenIs(lexer.RBRACK) && !p.curTokenIs(lexer.EOF) {
		elem := p.parseExpression()
		if elem != nil {
			elements = append(elements, elem)
		}
		if p.curTokenIs(lexer.COMMA) {
			p.nextToken()
		}
	}

	p.expectPeek(lexer.RBRACK)
	return &ast.ArrayExpr{Elements: elements, Pos: pos}
}

// Infix parsers
func (p *Parser) parseInfix(left ast.Expr) ast.Expr {
	pos := p.tokenPos()
	operator := p.curToken.Lit
	prec := p.curPrecedence()
	p.nextToken()
	right := p.parseExprPrecedence(prec)
	return &ast.InfixExpr{Left: left, Operator: operator, Right: right, Pos: pos}
}

// Function calls
func (p *Parser) parseCallExpression(fn ast.Expr) ast.Expr {
	pos := p.tokenPos()
	args := []ast.Expr{}

	p.nextToken()
	for !p.curTokenIs(lexer.RPAREN) && !p.curTokenIs(lexer.EOF) {
		arg := p.parseExpression()
		if arg != nil {
			args = append(args, arg)
		}
		if p.curTokenIs(lexer.COMMA) {
			p.nextToken()
		}
	}

	p.expectPeek(lexer.RPAREN)
	return &ast.CallExpr{Function: fn, Args: args, Pos: pos}
}

// Precedence helpers
func (p *Parser) peekPrecedence() int {
	if prec, ok := precedences[p.peekToken.Type]; ok {
		return prec
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if prec, ok := precedences[p.curToken.Type]; ok {
		return prec
	}
	return LOWEST
}

// ---------------------
// Assets
// ---------------------
func (p *Parser) parseAssets() []*ast.AssetDecl {
	p.nextToken() // skip WITH
	for !p.curTokenIs(lexer.EOF) && !p.isStartOfStatement() {
		p.nextToken()
	}
	return []*ast.AssetDecl{}
}

func (p *Parser) isStartOfStatement() bool {
	switch p.curToken.Type {
	case lexer.PRINT, lexer.SAY, lexer.WAIT, lexer.RETURN,
		lexer.MODULE, lexer.SUMMON, lexer.WITH,
		lexer.MAKE, lexer.CREATE, lexer.FN:
		return true
	default:
		return false
	}
}

func (p *Parser) tokenPos() ast.Position {
	return ast.Position{Line: p.curToken.Pos.Line, Column: p.curToken.Pos.Col, Offset: p.curToken.Pos.Offset}
}

// Package parser converts a flat token stream into an AST for Hunter's Engine.
//
// Grammar (simplified):
//
//	program      → stmt* EOF
//	stmt         → classDecl | objectDecl | importStmt
//	             | setStmt | printStmt | inputStmt
//	             | ifStmt | whileStmt | forEachStmt
//	             | returnStmt | breakStmt | continueStmt
//	             | callStmt | exprStmt
//	expr         → logicOr
//	logicOr      → logicAnd ( "or" logicAnd )*
//	logicAnd     → equality ( "and" equality )*
//	equality     → comparison ( ("==" | "!=" | "<>") comparison )*
//	comparison   → term ( ("<" | ">" | "<=" | ">=") term )*
//	term         → factor ( ("+" | "-") factor )*
//	factor       → unary ( ("*" | "/" | "%") unary )*
//	unary        → "not" unary | "-" unary | power
//	power        → postfix ( "^" unary )*
//	postfix      → primary ( "." IDENT | "[" expr "]" | "(" args ")" )*
//	primary      → NUMBER | STRING | "true" | "false" | "nil"
//	             | IDENT | "this" | "super"
//	             | "new" IDENT initFields?
//	             | "(" expr ")"
//	             | "[" args "]"
package parser

import (
	"fmt"
	"strconv"

	"github.com/hunter/he/pkg/ast"
	"github.com/hunter/he/pkg/lexer"
)

// Error is a parse diagnostic.
type Error struct {
	Pos     ast.Pos
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("parse error at %s: %s", e.Pos, e.Message)
}

// Parser consumes tokens and builds an AST.
type Parser struct {
	tokens  []lexer.Token
	pos     int
	errors  []*Error
}

// New creates a Parser from a token slice (as produced by Lexer.Tokenise).
func New(tokens []lexer.Token) *Parser {
	return &Parser{tokens: tokens}
}

// Errors returns all parse diagnostics.
func (p *Parser) Errors() []*Error { return p.errors }

// Parse parses a complete program.
func (p *Parser) Parse() *ast.Program {
	prog := &ast.Program{}
	for !p.isAtEnd() {
		p.skipNewlines()
		if p.isAtEnd() {
			break
		}
		s := p.parseStmt()
		if s != nil {
			prog.Statements = append(prog.Statements, s)
		}
	}
	return prog
}

// ─── token stream helpers ─────────────────────────────────────────────────────

func (p *Parser) peek() lexer.Token {
	if p.pos >= len(p.tokens) {
		return lexer.Token{Type: lexer.TOKEN_EOF}
	}
	return p.tokens[p.pos]
}

func (p *Parser) peekType() lexer.TokenType { return p.peek().Type }

func (p *Parser) peekAt(offset int) lexer.Token {
	idx := p.pos + offset
	if idx >= len(p.tokens) {
		return lexer.Token{Type: lexer.TOKEN_EOF}
	}
	return p.tokens[idx]
}

func (p *Parser) advance() lexer.Token {
	tok := p.tokens[p.pos]
	p.pos++
	return tok
}

func (p *Parser) isAtEnd() bool {
	return p.peekType() == lexer.TOKEN_EOF
}

func (p *Parser) check(t lexer.TokenType) bool { return p.peekType() == t }

func (p *Parser) match(types ...lexer.TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			return true
		}
	}
	return false
}

func (p *Parser) consume(t lexer.TokenType, msg string) lexer.Token {
	if p.check(t) {
		return p.advance()
	}
	tok := p.peek()
	p.addError(tok, msg)
	return tok
}

func (p *Parser) skipNewlines() {
	for p.check(lexer.TOKEN_NEWLINE) {
		p.advance()
	}
}

func (p *Parser) expectNewlineOrSemi() {
	if p.check(lexer.TOKEN_NEWLINE) || p.check(lexer.TOKEN_SEMICOLON) {
		p.advance()
	}
	// tolerate missing newline (best-effort recovery)
}

func (p *Parser) addError(tok lexer.Token, msg string) {
	p.errors = append(p.errors, &Error{
		Pos:     ast.Pos{Line: tok.Line, Col: tok.Col},
		Message: msg,
	})
}

func (p *Parser) tokPos(tok lexer.Token) ast.Pos {
	return ast.Pos{Line: tok.Line, Col: tok.Col}
}

// ─── statement dispatch ───────────────────────────────────────────────────────

func (p *Parser) parseStmt() ast.Stmt {
	tok := p.peek()
	switch tok.Type {
	case lexer.TOKEN_CLASS:
		return p.parseClassDecl()
	case lexer.TOKEN_OBJECT:
		return p.parseObjectDecl()
	case lexer.TOKEN_IMPORT:
		return p.parseImport()
	case lexer.TOKEN_SET:
		return p.parseSetStmt()
	case lexer.TOKEN_PRINT:
		return p.parsePrintStmt()
	case lexer.TOKEN_INPUT:
		return p.parseInputStmt()
	case lexer.TOKEN_IF:
		return p.parseIfStmt()
	case lexer.TOKEN_WHILE:
		return p.parseWhileStmt()
	case lexer.TOKEN_FOR:
		return p.parseForEachStmt()
	case lexer.TOKEN_RETURN, lexer.TOKEN_ANSWER:
		return p.parseReturnStmt()
	case lexer.TOKEN_BREAK:
		p.advance()
		p.expectNewlineOrSemi()
		return &ast.BreakStmt{Pos: p.tokPos(tok)}
	case lexer.TOKEN_CONTINUE:
		p.advance()
		p.expectNewlineOrSemi()
		return &ast.ContinueStmt{Pos: p.tokPos(tok)}
	case lexer.TOKEN_CALL:
		return p.parseCallStmt()
	case lexer.TOKEN_NEWLINE:
		p.advance()
		return nil
	default:
		return p.parseExprOrAssignStmt()
	}
}

// parseBody reads statements until one of the stop tokens is encountered.
func (p *Parser) parseBody(stopTokens ...lexer.TokenType) []ast.Stmt {
	var stmts []ast.Stmt
	for {
		p.skipNewlines()
		if p.isAtEnd() {
			break
		}
		if p.match(stopTokens...) {
			break
		}
		s := p.parseStmt()
		if s != nil {
			stmts = append(stmts, s)
		}
	}
	return stmts
}

// ─── declarations ─────────────────────────────────────────────────────────────

func (p *Parser) parseClassDecl() *ast.ClassDecl {
	kw := p.advance() // consume 'class'
	pos := p.tokPos(kw)
	nameTok := p.consume(lexer.TOKEN_IDENT, "expected class name after 'class'")
	name := nameTok.Literal

	superClass := ""
	if p.match(lexer.TOKEN_EXTENDS) || p.match(lexer.TOKEN_INHERITS) {
		p.advance()
		super := p.consume(lexer.TOKEN_IDENT, "expected superclass name")
		superClass = super.Literal
	}
	p.expectNewlineOrSemi()

	var props []*ast.PropertyDecl
	var methods []*ast.MethodDecl
	for {
		p.skipNewlines()
		if p.isAtEnd() || p.check(lexer.TOKEN_END) {
			break
		}
		switch p.peekType() {
		case lexer.TOKEN_PROPERTY:
			props = append(props, p.parsePropertyDecl())
		case lexer.TOKEN_METHOD:
			methods = append(methods, p.parseMethodDecl())
		default:
			p.addError(p.peek(), fmt.Sprintf("unexpected token %q inside class body", p.peek().Literal))
			p.advance()
		}
	}
	p.consume(lexer.TOKEN_END, "expected 'end' to close class '"+name+"'")
	p.expectNewlineOrSemi()

	return &ast.ClassDecl{Pos: pos, Name: name, SuperClass: superClass, Properties: props, Methods: methods}
}

func (p *Parser) parsePropertyDecl() *ast.PropertyDecl {
	kw := p.advance() // 'property'
	pos := p.tokPos(kw)
	nameTok := p.consume(lexer.TOKEN_IDENT, "expected property name")
	var def ast.Expr
	if p.check(lexer.TOKEN_IS) || p.check(lexer.TOKEN_ASSIGN) {
		p.advance()
		def = p.parseExpr()
	}
	p.expectNewlineOrSemi()
	return &ast.PropertyDecl{Pos: pos, Name: nameTok.Literal, Default: def}
}

func (p *Parser) parseMethodDecl() *ast.MethodDecl {
	kw := p.advance() // 'method'
	pos := p.tokPos(kw)
	nameTok := p.consume(lexer.TOKEN_IDENT, "expected method name")

	var params []string
	if p.check(lexer.TOKEN_WITH) || p.check(lexer.TOKEN_LPAREN) {
		if p.check(lexer.TOKEN_WITH) {
			p.advance()
		} else {
			p.advance() // '('
		}
		for p.check(lexer.TOKEN_IDENT) {
			params = append(params, p.advance().Literal)
			if p.check(lexer.TOKEN_COMMA) {
				p.advance()
			}
		}
		if p.check(lexer.TOKEN_RPAREN) {
			p.advance()
		}
	}

	// optional 'returns' keyword (documentation hint)
	if p.check(lexer.TOKEN_RETURNS) {
		p.advance()
	}

	p.expectNewlineOrSemi()
	body := p.parseBody(lexer.TOKEN_END)
	p.consume(lexer.TOKEN_END, "expected 'end' to close method '"+nameTok.Literal+"'")
	p.expectNewlineOrSemi()

	return &ast.MethodDecl{Pos: pos, Name: nameTok.Literal, Params: params, Body: body}
}

func (p *Parser) parseObjectDecl() *ast.ObjectDecl {
	kw := p.advance() // 'object'
	pos := p.tokPos(kw)
	nameTok := p.consume(lexer.TOKEN_IDENT, "expected object name after 'object'")
	p.consume(lexer.TOKEN_OF, "expected 'of' after object name")
	classTok := p.consume(lexer.TOKEN_IDENT, "expected class name after 'of'")

	initFields := map[string]ast.Expr{}
	if p.check(lexer.TOKEN_WITH) {
		p.advance()
		for {
			fieldTok := p.consume(lexer.TOKEN_IDENT, "expected field name")
			if p.check(lexer.TOKEN_IS) || p.check(lexer.TOKEN_ASSIGN) {
				p.advance()
			}
			val := p.parseExpr()
			initFields[fieldTok.Literal] = val
			if !p.check(lexer.TOKEN_COMMA) {
				break
			}
			p.advance()
		}
	}
	p.expectNewlineOrSemi()
	return &ast.ObjectDecl{Pos: pos, Name: nameTok.Literal, ClassName: classTok.Literal, InitFields: initFields}
}

func (p *Parser) parseImport() *ast.ImportStmt {
	kw := p.advance() // 'import'
	pos := p.tokPos(kw)
	nameTok := p.consume(lexer.TOKEN_IDENT, "expected module name after 'import'")
	path := ""
	if p.check(lexer.TOKEN_FROM) {
		p.advance()
		pathTok := p.consume(lexer.TOKEN_STRING, "expected path string after 'from'")
		path = pathTok.Literal
	}
	p.expectNewlineOrSemi()
	return &ast.ImportStmt{Pos: pos, Name: nameTok.Literal, Path: path}
}

// ─── statements ───────────────────────────────────────────────────────────────

func (p *Parser) parseSetStmt() *ast.SetStmt {
	kw := p.advance() // 'set'
	pos := p.tokPos(kw)
	target := p.parsePostfix()
	if p.check(lexer.TOKEN_TO) || p.check(lexer.TOKEN_IS) || p.check(lexer.TOKEN_ASSIGN) {
		p.advance()
	}
	val := p.parseExpr()
	p.expectNewlineOrSemi()
	return &ast.SetStmt{Pos: pos, Target: target, Value: val}
}

func (p *Parser) parsePrintStmt() *ast.PrintStmt {
	kw := p.advance() // 'print'
	expr := p.parseExpr()
	p.expectNewlineOrSemi()
	return &ast.PrintStmt{Pos: p.tokPos(kw), Expr: expr}
}

func (p *Parser) parseInputStmt() *ast.InputStmt {
	kw := p.advance() // 'input'
	pos := p.tokPos(kw)
	targetTok := p.consume(lexer.TOKEN_IDENT, "expected variable name after 'input'")
	var prompt ast.Expr
	if p.check(lexer.TOKEN_AS) {
		p.advance()
		prompt = p.parseExpr()
	}
	p.expectNewlineOrSemi()
	return &ast.InputStmt{Pos: pos, Target: targetTok.Literal, Prompt: prompt}
}

func (p *Parser) parseIfStmt() *ast.IfStmt {
	kw := p.advance() // 'if'
	pos := p.tokPos(kw)
	cond := p.parseExpr()
	if p.check(lexer.TOKEN_THEN) {
		p.advance()
	}
	p.expectNewlineOrSemi()

	consequent := p.parseBody(lexer.TOKEN_ELSE, lexer.TOKEN_END)

	var alternate []ast.Stmt
	if p.check(lexer.TOKEN_ELSE) {
		p.advance()
		// else if → nest
		if p.check(lexer.TOKEN_IF) {
			alternate = []ast.Stmt{p.parseIfStmt()}
		} else {
			p.expectNewlineOrSemi()
			alternate = p.parseBody(lexer.TOKEN_END)
			p.consume(lexer.TOKEN_END, "expected 'end' to close if")
			p.expectNewlineOrSemi()
		}
	} else {
		p.consume(lexer.TOKEN_END, "expected 'end' or 'else' to close if")
		p.expectNewlineOrSemi()
	}

	return &ast.IfStmt{Pos: pos, Condition: cond, Consequent: consequent, Alternate: alternate}
}

func (p *Parser) parseWhileStmt() *ast.WhileStmt {
	kw := p.advance() // 'while'
	pos := p.tokPos(kw)
	cond := p.parseExpr()
	if p.check(lexer.TOKEN_DO) {
		p.advance()
	}
	p.expectNewlineOrSemi()
	body := p.parseBody(lexer.TOKEN_END)
	p.consume(lexer.TOKEN_END, "expected 'end' to close while")
	p.expectNewlineOrSemi()
	return &ast.WhileStmt{Pos: pos, Condition: cond, Body: body}
}

func (p *Parser) parseForEachStmt() *ast.ForEachStmt {
	kw := p.advance() // 'for'
	pos := p.tokPos(kw)
	if p.check(lexer.TOKEN_EACH) {
		p.advance()
	}
	varTok := p.consume(lexer.TOKEN_IDENT, "expected loop variable after 'for each'")
	p.consume(lexer.TOKEN_IN, "expected 'in' after loop variable")
	iterable := p.parseExpr()
	if p.check(lexer.TOKEN_DO) {
		p.advance()
	}
	p.expectNewlineOrSemi()
	body := p.parseBody(lexer.TOKEN_END)
	p.consume(lexer.TOKEN_END, "expected 'end' to close for-each")
	p.expectNewlineOrSemi()
	return &ast.ForEachStmt{Pos: pos, VarName: varTok.Literal, Iterable: iterable, Body: body}
}

func (p *Parser) parseReturnStmt() *ast.ReturnStmt {
	kw := p.advance() // 'return' or 'answer'
	pos := p.tokPos(kw)
	var val ast.Expr
	if !p.check(lexer.TOKEN_NEWLINE) && !p.check(lexer.TOKEN_SEMICOLON) && !p.isAtEnd() {
		val = p.parseExpr()
	}
	p.expectNewlineOrSemi()
	return &ast.ReturnStmt{Pos: pos, Value: val}
}

func (p *Parser) parseCallStmt() *ast.CallStmt {
	kw := p.advance() // 'call'
	pos := p.tokPos(kw)
	methodTok := p.consume(lexer.TOKEN_IDENT, "expected method name after 'call'")
	p.consume(lexer.TOKEN_ON, "expected 'on' after method name")
	obj := p.parsePostfix()

	var args []ast.Expr
	if p.check(lexer.TOKEN_WITH) || p.check(lexer.TOKEN_LPAREN) {
		if p.check(lexer.TOKEN_WITH) {
			p.advance()
		} else {
			p.advance()
		}
		args = p.parseArgList(lexer.TOKEN_RPAREN)
	}
	p.expectNewlineOrSemi()
	return &ast.CallStmt{Pos: pos, Method: methodTok.Literal, Object: obj, Args: args}
}

// parseExprOrAssignStmt handles   <ident> is <expr>
//
//	<ident> = <expr>
//	<expr> += <expr>  etc.
//	or a plain expression statement.
func (p *Parser) parseExprOrAssignStmt() ast.Stmt {
	expr := p.parsePostfix()
	tok := p.peek()

	switch tok.Type {
	case lexer.TOKEN_IS, lexer.TOKEN_ASSIGN:
		p.advance()
		val := p.parseExpr()
		p.expectNewlineOrSemi()
		return &ast.SetStmt{Pos: expr.nodePos(), Target: expr, Value: val}

	case lexer.TOKEN_PLUS_EQ, lexer.TOKEN_MINUS_EQ, lexer.TOKEN_STAR_EQ, lexer.TOKEN_SLASH_EQ:
		op := p.advance().Literal
		val := p.parseExpr()
		p.expectNewlineOrSemi()
		return &ast.CompoundAssignStmt{Pos: expr.nodePos(), Operator: op, Target: expr, Value: val}
	}

	// finish the expression (it might be a binary), then wrap as ExprStmt
	fullExpr := p.continueExpr(expr)
	p.expectNewlineOrSemi()
	return &ast.ExprStmt{Pos: fullExpr.nodePos(), Expr: fullExpr}
}

// ─── expressions ──────────────────────────────────────────────────────────────

func (p *Parser) parseExpr() ast.Expr       { return p.parseLogicOr() }
func (p *Parser) continueExpr(left ast.Expr) ast.Expr { return p.finishBinary(left) }

func (p *Parser) finishBinary(left ast.Expr) ast.Expr {
	return p.finishLogicOr(left)
}

func (p *Parser) parseLogicOr() ast.Expr {
	left := p.parseLogicAnd()
	return p.finishLogicOr(left)
}

func (p *Parser) finishLogicOr(left ast.Expr) ast.Expr {
	for p.check(lexer.TOKEN_OR) {
		op := p.advance()
		right := p.parseLogicAnd()
		left = &ast.BinaryExpr{Pos: p.tokPos(op), Operator: "or", Left: left, Right: right}
	}
	return left
}

func (p *Parser) parseLogicAnd() ast.Expr {
	left := p.parseEquality()
	for p.check(lexer.TOKEN_AND) {
		op := p.advance()
		right := p.parseEquality()
		left = &ast.BinaryExpr{Pos: p.tokPos(op), Operator: "and", Left: left, Right: right}
	}
	return left
}

func (p *Parser) parseEquality() ast.Expr {
	left := p.parseComparison()
	for p.match(lexer.TOKEN_EQ, lexer.TOKEN_NEQ) {
		op := p.advance()
		right := p.parseComparison()
		left = &ast.BinaryExpr{Pos: p.tokPos(op), Operator: op.Literal, Left: left, Right: right}
	}
	return left
}

func (p *Parser) parseComparison() ast.Expr {
	left := p.parseTerm()
	for p.match(lexer.TOKEN_LT, lexer.TOKEN_GT, lexer.TOKEN_LTE, lexer.TOKEN_GTE) {
		op := p.advance()
		right := p.parseTerm()
		left = &ast.BinaryExpr{Pos: p.tokPos(op), Operator: op.Literal, Left: left, Right: right}
	}
	return left
}

func (p *Parser) parseTerm() ast.Expr {
	left := p.parseFactor()
	for p.match(lexer.TOKEN_PLUS, lexer.TOKEN_MINUS) {
		op := p.advance()
		right := p.parseFactor()
		left = &ast.BinaryExpr{Pos: p.tokPos(op), Operator: op.Literal, Left: left, Right: right}
	}
	return left
}

func (p *Parser) parseFactor() ast.Expr {
	left := p.parseUnary()
	for p.match(lexer.TOKEN_STAR, lexer.TOKEN_SLASH, lexer.TOKEN_PERCENT) {
		op := p.advance()
		right := p.parseUnary()
		left = &ast.BinaryExpr{Pos: p.tokPos(op), Operator: op.Literal, Left: left, Right: right}
	}
	return left
}

func (p *Parser) parseUnary() ast.Expr {
	if p.check(lexer.TOKEN_NOT) {
		op := p.advance()
		operand := p.parseUnary()
		return &ast.UnaryExpr{Pos: p.tokPos(op), Operator: "not", Operand: operand}
	}
	if p.check(lexer.TOKEN_MINUS) {
		op := p.advance()
		operand := p.parseUnary()
		return &ast.UnaryExpr{Pos: p.tokPos(op), Operator: "-", Operand: operand}
	}
	return p.parsePower()
}

func (p *Parser) parsePower() ast.Expr {
	base := p.parsePostfix()
	if p.check(lexer.TOKEN_CARET) {
		op := p.advance()
		exp := p.parseUnary()
		return &ast.BinaryExpr{Pos: p.tokPos(op), Operator: "^", Left: base, Right: exp}
	}
	return base
}

func (p *Parser) parsePostfix() ast.Expr {
	expr := p.parsePrimary()
	for {
		switch p.peekType() {
		case lexer.TOKEN_DOT:
			dot := p.advance()
			field := p.consume(lexer.TOKEN_IDENT, "expected field name after '.'")
			// method call?
			if p.check(lexer.TOKEN_LPAREN) {
				p.advance()
				args := p.parseArgList(lexer.TOKEN_RPAREN)
				p.consume(lexer.TOKEN_RPAREN, "expected ')' after argument list")
				expr = &ast.CallExpr{Pos: p.tokPos(dot), Object: expr, Method: field.Literal, Args: args}
			} else {
				expr = &ast.DotExpr{Pos: p.tokPos(dot), Object: expr, Field: field.Literal}
			}
		case lexer.TOKEN_LBRACKET:
			lbr := p.advance()
			idx := p.parseExpr()
			p.consume(lexer.TOKEN_RBRACKET, "expected ']' after index")
			expr = &ast.IndexExpr{Pos: p.tokPos(lbr), Object: expr, Index: idx}
		default:
			return expr
		}
	}
}

func (p *Parser) parsePrimary() ast.Expr {
	tok := p.peek()

	switch tok.Type {
	case lexer.TOKEN_NUMBER:
		p.advance()
		val, _ := strconv.ParseFloat(tok.Literal, 64)
		return &ast.NumberLit{Pos: p.tokPos(tok), Value: val, Raw: tok.Literal}

	case lexer.TOKEN_STRING:
		p.advance()
		return &ast.StringLit{Pos: p.tokPos(tok), Value: tok.Literal}

	case lexer.TOKEN_TRUE:
		p.advance()
		return &ast.BoolLit{Pos: p.tokPos(tok), Value: true}

	case lexer.TOKEN_FALSE:
		p.advance()
		return &ast.BoolLit{Pos: p.tokPos(tok), Value: false}

	case lexer.TOKEN_NIL, lexer.TOKEN_NOTHING:
		p.advance()
		return &ast.NilLit{Pos: p.tokPos(tok)}

	case lexer.TOKEN_THIS:
		p.advance()
		return &ast.Identifier{Pos: p.tokPos(tok), Name: "this"}

	case lexer.TOKEN_SUPER:
		p.advance()
		return &ast.Identifier{Pos: p.tokPos(tok), Name: "super"}

	case lexer.TOKEN_NEW:
		return p.parseNewExpr()

	case lexer.TOKEN_IDENT:
		p.advance()
		// function call shorthand: ident(args)
		if p.check(lexer.TOKEN_LPAREN) {
			p.advance()
			args := p.parseArgList(lexer.TOKEN_RPAREN)
			p.consume(lexer.TOKEN_RPAREN, "expected ')' after arguments")
			return &ast.FuncCallExpr{Pos: p.tokPos(tok), Name: tok.Literal, Args: args}
		}
		return &ast.Identifier{Pos: p.tokPos(tok), Name: tok.Literal}

	case lexer.TOKEN_LPAREN:
		p.advance()
		expr := p.parseExpr()
		p.consume(lexer.TOKEN_RPAREN, "expected ')' after expression")
		return expr

	case lexer.TOKEN_LBRACKET:
		return p.parseListLiteral()
	}

	p.addError(tok, fmt.Sprintf("unexpected token %q in expression", tok.Literal))
	p.advance() // skip it to avoid infinite loop
	return &ast.NilLit{Pos: p.tokPos(tok)}
}

func (p *Parser) parseNewExpr() *ast.NewExpr {
	kw := p.advance() // 'new'
	pos := p.tokPos(kw)
	classTok := p.consume(lexer.TOKEN_IDENT, "expected class name after 'new'")
	initFields := map[string]ast.Expr{}
	if p.check(lexer.TOKEN_WITH) {
		p.advance()
		for p.check(lexer.TOKEN_IDENT) {
			fieldTok := p.advance()
			if p.check(lexer.TOKEN_IS) || p.check(lexer.TOKEN_ASSIGN) {
				p.advance()
			}
			val := p.parseExpr()
			initFields[fieldTok.Literal] = val
			if !p.check(lexer.TOKEN_COMMA) {
				break
			}
			p.advance()
		}
	}
	return &ast.NewExpr{Pos: pos, ClassName: classTok.Literal, InitFields: initFields}
}

func (p *Parser) parseListLiteral() *ast.ListLiteral {
	lbr := p.advance() // '['
	elems := p.parseArgList(lexer.TOKEN_RBRACKET)
	p.consume(lexer.TOKEN_RBRACKET, "expected ']' after list elements")
	return &ast.ListLiteral{Pos: p.tokPos(lbr), Elements: elems}
}

func (p *Parser) parseArgList(terminator lexer.TokenType) []ast.Expr {
	var args []ast.Expr
	for !p.check(terminator) && !p.isAtEnd() {
		args = append(args, p.parseExpr())
		if !p.check(lexer.TOKEN_COMMA) {
			break
		}
		p.advance()
	}
	return args
}

package parser

import (
	"fmt"

	"hunterlang/lang/ast"
	"hunterlang/lang/lexer"
	"hunterlang/lang/token"
)

type Parser struct {
	lx  *lexer.Lexer
	cur token.Token
}

func New(lx *lexer.Lexer) *Parser {
	p := &Parser{lx: lx}
	t, err := p.lx.NextToken()
	if err != nil {
		p.cur = token.Token{Type: token.ILLEGAL, Lexeme: "", Line: 1, Column: 1}
	} else {
		p.cur = t
	}
	return p
}

func (p *Parser) advance() error {
	t, err := p.lx.NextToken()
	if err != nil {
		return err
	}
	p.cur = t
	return nil
}

func (p *Parser) expect(tt token.TokenType) (token.Token, error) {
	if p.cur.Type != tt {
		return token.Token{}, fmt.Errorf("expected %s at %d:%d, got %s (%q)", tt, p.cur.Line, p.cur.Column, p.cur.Type, p.cur.Lexeme)
	}
	t := p.cur
	if err := p.advance(); err != nil {
		return token.Token{}, err
	}
	return t, nil
}

func (p *Parser) match(tt token.TokenType) bool {
	if p.cur.Type != tt {
		return false
	}
	_ = p.advance()
	return true
}

func (p *Parser) ParseProgram() (*ast.Program, error) {
	var lines []ast.Line
	for p.cur.Type != token.EOF {
		line, err := p.parseLine()
		if err != nil {
			return nil, err
		}
		if line != nil {
			lines = append(lines, line)
		}
	}
	return &ast.Program{Lines: lines}, nil
}

func (p *Parser) parseLine() (ast.Line, error) {
	switch p.cur.Type {
	case token.K_SUMMON:
		return p.parseSummon()
	case token.K_CREATE, token.K_MAKE:
		return p.parseObjectLine()
	case token.K_PRINT, token.K_SAY, token.K_SET, token.K_WAIT, token.K_CAN_TELL:
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		return &ast.GlobalStatementLine{Statement: stmt}, nil
	default:
		return nil, fmt.Errorf("unexpected token %s (%q) at %d:%d", p.cur.Type, p.cur.Lexeme, p.cur.Line, p.cur.Column)
	}
}

func (p *Parser) parseSummon() (ast.Line, error) {
	if _, err := p.expect(token.K_SUMMON); err != nil {
		return nil, err
	}
	modTok, err := p.expect(token.STRING)
	if err != nil {
		return nil, err
	}

	if p.cur.Type != token.K_NAMED && p.cur.Type != token.K_AS {
		return nil, fmt.Errorf("expected 'named' or 'as' at %d:%d", p.cur.Line, p.cur.Column)
	}
	_ = p.advance()

	aliasTok, err := p.expect(token.IDENT)
	if err != nil {
		return nil, err
	}
	alias := &ast.Identifier{Name: aliasTok.Lexeme, Tok: aliasTok}

	return ast.SummonLine{ModuleName: modTok, Alias: alias}, nil
}

func (p *Parser) parseObjectLine() (*ast.ObjectLine, error) {
	kindTok := p.cur
	if p.cur.Type != token.K_CREATE && p.cur.Type != token.K_MAKE {
		return nil, fmt.Errorf("expected 'create' or 'make' at %d:%d", p.cur.Line, p.cur.Column)
	}
	_ = p.advance()

	nameTok, err := p.expect(token.IDENT)
	if err != nil {
		return nil, err
	}
	obj := &ast.ObjectLine{
		Kind: kindTok.Lexeme,
		Name: ast.Identifier{Name: nameTok.Lexeme, Tok: nameTok},
	}

	if p.match(token.K_LIKE) {
		likeTok, err := p.expect(token.IDENT)
		if err != nil {
			return nil, err
		}
		obj.Like = &ast.Identifier{Name: likeTok.Lexeme, Tok: likeTok}
	}

	var end token.TokenType
	switch p.cur.Type {
	case token.LBRACE:
		end = token.RBRACE
	case token.LBRACK:
		end = token.RBRACK
	default:
		return nil, fmt.Errorf("expected '{' or '[' at %d:%d, got %s (%q)", p.cur.Line, p.cur.Column, p.cur.Type, p.cur.Lexeme)
	}
	_ = p.advance()

	body, err := p.parseObjectBody(end)
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(end); err != nil {
		return nil, err
	}
	obj.Body = body
	return obj, nil
}

func (p *Parser) parseObjectBody(end token.TokenType) (ast.ObjectBody, error) {
	var body ast.ObjectBody

	for p.cur.Type != end && p.cur.Type != token.EOF {
		if p.cur.Type == token.K_HAS || p.cur.Type == token.K_OWNS || p.cur.Type == token.K_CARRIES {
			sectionKind := p.cur.Lexeme
			_ = p.advance()
			_ = p.match(token.COLON)

			if _, err := p.expect(token.LBRACK); err != nil {
				return ast.ObjectBody{}, err
			}

			props, err := p.parsePropertiesSection(sectionKind)
			if err != nil {
				return ast.ObjectBody{}, err
			}
			body.Sections = append(body.Sections, props)
			continue
		}

		if p.cur.Type == token.K_CAN {
			_ = p.advance()
			_ = p.match(token.COLON)

			if _, err := p.expect(token.LBRACK); err != nil {
				return ast.ObjectBody{}, err
			}
			abilities, err := p.parseAbilitiesList()
			if err != nil {
				return ast.ObjectBody{}, err
			}
			if _, err := p.expect(token.RBRACK); err != nil {
				return ast.ObjectBody{}, err
			}
			body.Sections = append(body.Sections, abilities)
			continue
		}

		if p.cur.Type == token.K_ON || p.cur.Type == token.K_WHEN || p.cur.Type == token.K_WHENEVER {
			react, err := p.parseReactionsSection()
			if err != nil {
				return ast.ObjectBody{}, err
			}
			body.Sections = append(body.Sections, react)
			continue
		}

		if p.cur.Type == token.K_REMEMBERS {
			return ast.ObjectBody{}, fmt.Errorf("memories not supported yet at %d:%d", p.cur.Line, p.cur.Column)
		}

		if p.cur.Type == token.K_CREATE || p.cur.Type == token.K_MAKE {
			obj, err := p.parseObjectLine()
			if err != nil {
				return ast.ObjectBody{}, err
			}
			body.NestedObjects = append(body.NestedObjects, *obj)
			continue
		}

		return ast.ObjectBody{}, fmt.Errorf("unexpected token in object body: %s (%q) at %d:%d", p.cur.Type, p.cur.Lexeme, p.cur.Line, p.cur.Column)
	}

	return body, nil
}

func (p *Parser) parsePropertiesSection(sectionKind string) (ast.PropertiesSection, error) {
	var props []ast.Property
	for p.cur.Type != token.RBRACK && p.cur.Type != token.EOF {
		nameTok, err := p.expect(token.IDENT)
		if err != nil {
			return ast.PropertiesSection{}, err
		}

		kind := ""
		if p.match(token.K_IS) {
			kind = "is"
		} else if p.cur.Type == token.IDENT && p.cur.Lexeme == "starts" {
			_ = p.advance()
			if _, err := p.expect(token.K_AS); err != nil {
				return ast.PropertiesSection{}, err
			}
			kind = "starts as"
		} else {
			return ast.PropertiesSection{}, fmt.Errorf("expected 'is' or 'starts as' at %d:%d", p.cur.Line, p.cur.Column)
		}

		ex, err := p.parseExpression()
		if err != nil {
			return ast.PropertiesSection{}, err
		}

		props = append(props, ast.Property{
			Name:  nameTok.Lexeme,
			Kind:  sectionKind + " " + kind,
			Value: ex,
		})
	}

	if _, err := p.expect(token.RBRACK); err != nil {
		return ast.PropertiesSection{}, err
	}
	return ast.PropertiesSection{Props: props}, nil
}

func (p *Parser) parseAbilitiesList() (ast.AbilitiesSection, error) {
	var actions []ast.ActionDef
	for p.cur.Type != token.RBRACK && p.cur.Type != token.EOF {
		nameTok, err := p.expect(token.IDENT)
		if err != nil {
			return ast.AbilitiesSection{}, err
		}
		action := ast.ActionDef{Name: nameTok.Lexeme}

		if p.match(token.LPAREN) {
			if p.cur.Type != token.RPAREN {
				params, err := p.parseParamList()
				if err != nil {
					return ast.AbilitiesSection{}, err
				}
				action.Params = params
				action.HasParams = true
			}
			if _, err := p.expect(token.RPAREN); err != nil {
				return ast.AbilitiesSection{}, err
			}
		}

		if p.cur.Type == token.IDENT && p.cur.Lexeme == "returns" {
			_ = p.advance()
			typ, err := p.parseType()
			if err != nil {
				return ast.AbilitiesSection{}, err
			}
			action.Returns = &typ
		}

		if _, err := p.expect(token.LBRACK); err != nil {
			return ast.AbilitiesSection{}, err
		}
		bodyStmts, err := p.parseStatementListUntil(token.RBRACK)
		if err != nil {
			return ast.AbilitiesSection{}, err
		}
		if _, err := p.expect(token.RBRACK); err != nil {
			return ast.AbilitiesSection{}, err
		}
		action.Body = bodyStmts
		actions = append(actions, action)
	}
	return ast.AbilitiesSection{Abilities: actions}, nil
}

func (p *Parser) parseParamList() ([]ast.Param, error) {
	var params []ast.Param
	for {
		nameTok, err := p.expect(token.IDENT)
		if err != nil {
			return nil, err
		}

		if p.cur.Type == token.COLON {
			_ = p.advance()
			typ, err := p.parseType()
			if err != nil {
				return nil, err
			}
			params = append(params, ast.Param{Name: nameTok.Lexeme, Type: typ})
		} else {
			params = append(params, ast.Param{Name: nameTok.Lexeme, Type: ast.TypeNode{Name: "unknown"}})
		}

		if p.cur.Type != token.COMMA {
			break
		}
		_ = p.advance()
	}
	return params, nil
}

func (p *Parser) parseType() (ast.TypeNode, error) {
	if p.cur.Type != token.IDENT {
		return ast.TypeNode{}, fmt.Errorf("expected type at %d:%d", p.cur.Line, p.cur.Column)
	}
	name := p.cur.Lexeme
	_ = p.advance()

	if p.match(token.LBRACK) {
		if _, err := p.expect(token.RBRACK); err != nil {
			return ast.TypeNode{}, err
		}
		name = name + "[]"
	}
	return ast.TypeNode{Name: name}, nil
}

func (p *Parser) parseReactionsSection() (ast.ReactionsSection, error) {
	_ = p.advance()

	trigger, err := p.parseTrigger()
	if err != nil {
		return ast.ReactionsSection{}, err
	}
	if _, err := p.expect(token.LBRACK); err != nil {
		return ast.ReactionsSection{}, err
	}
	body, err := p.parseStatementListUntil(token.RBRACK)
	if err != nil {
		return ast.ReactionsSection{}, err
	}
	if _, err := p.expect(token.RBRACK); err != nil {
		return ast.ReactionsSection{}, err
	}
	return ast.ReactionsSection{Reactions: []ast.ReactionDef{{Trigger: trigger, Body: body}}}, nil
}

func (p *Parser) parseTrigger() (ast.Trigger, error) {
	// Identifier or keyword-tokenized triggers like "collision"
	if p.cur.Type != token.IDENT && p.cur.Type != token.K_COLLISION {
		return ast.Trigger{}, fmt.Errorf("expected trigger identifier at %d:%d, got %s (%q)", p.cur.Line, p.cur.Column, p.cur.Type, p.cur.Lexeme)
	}

	leftTok := p.cur
	_ = p.advance()

	tr := ast.Trigger{Left: leftTok.Lexeme}

	if p.match(token.K_WITH) {
		rightTok, err := p.expect(token.IDENT)
		if err != nil {
			return ast.Trigger{}, err
		}
		tr.Right = &rightTok.Lexeme
	}
	return tr, nil
}

func (p *Parser) parseStatementListUntil(end token.TokenType) ([]ast.Statement, error) {
	var stmts []ast.Statement
	for p.cur.Type != end && p.cur.Type != token.EOF {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
	}
	return stmts, nil
}

func (p *Parser) parseStatement() (ast.Statement, error) {
	switch p.cur.Type {
	case token.K_PRINT, token.K_SAY:
		_ = p.advance()
		ex, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		return ast.SayStmt{Expr: ex}, nil

	case token.K_SET, token.K_MAKE:
		kindTok := p.cur
		_ = p.advance()

		nameTok, err := p.expect(token.IDENT)
		if err != nil {
			return nil, err
		}
		if p.cur.Type != token.K_TO {
			return nil, fmt.Errorf("expected 'to' after %s at %d:%d", kindTok.Lexeme, p.cur.Line, p.cur.Column)
		}
		_ = p.advance()

		ex, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		return ast.ChangeStmt{Kind: kindTok.Lexeme, Name: nameTok.Lexeme, Expr: ex}, nil
	}

	if p.cur.Type == token.IDENT && p.cur.Lexeme == "if" {
		_ = p.advance()
		cond, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(token.K_THEN); err != nil {
			return nil, err
		}
		if _, err := p.expect(token.LBRACK); err != nil {
			return nil, err
		}
		thenStmts, err := p.parseStatementListUntil(token.RBRACK)
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(token.RBRACK); err != nil {
			return nil, err
		}

		var elseStmts []ast.Statement
		if p.match(token.K_ELSE) {
			if _, err := p.expect(token.LBRACK); err != nil {
				return nil, err
			}
			elseStmts, err = p.parseStatementListUntil(token.RBRACK)
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(token.RBRACK); err != nil {
				return nil, err
			}
		}
		return ast.DecideStmt{Cond: cond, Then: thenStmts, Else: elseStmts}, nil
	}

	if p.cur.Type == token.IDENT && (p.cur.Lexeme == "repeat" || p.cur.Lexeme == "while") {
		kind := p.cur.Lexeme
		_ = p.advance()
		cond, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(token.LBRACK); err != nil {
			return nil, err
		}
		body, err := p.parseStatementListUntil(token.RBRACK)
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(token.RBRACK); err != nil {
			return nil, err
		}
		return ast.RepeatStmt{Kind: kind, Cond: cond, Body: body}, nil
	}

	if p.cur.Type == token.K_CAN_TELL {
		_ = p.advance()
		objTok, err := p.expect(token.IDENT)
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(token.K_TO); err != nil {
			return nil, err
		}
		actionTok, err := p.expect(token.IDENT)
		if err != nil {
			return nil, err
		}

		var args []ast.Expression
		if p.match(token.K_WITH) {
			args, err = p.parseArgList()
			if err != nil {
				return nil, err
			}
		}
		return ast.CallStmt{Object: objTok.Lexeme, Action: actionTok.Lexeme, Args: args}, nil
	}

	if p.cur.Type == token.K_WAIT {
		_ = p.advance()
		ex, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		unit := ""
		if p.cur.Type == token.K_SECONDS {
			unit = "seconds"
			_ = p.advance()
		} else if p.cur.Type == token.K_FRAMES {
			unit = "frames"
			_ = p.advance()
		} else {
			return nil, fmt.Errorf("expected seconds|frames after wait at %d:%d", p.cur.Line, p.cur.Column)
		}
		return ast.WaitStmt{Expr: ex, Unit: unit}, nil
	}

	if p.cur.Type == token.IDENT && p.cur.Lexeme == "return" {
		_ = p.advance()
		ex, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		return ast.ReturnStmt{Expr: ex}, nil
	}

	// Expression statement
	if p.cur.Type == token.IDENT {
		ex, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		return ast.ExprStmt{Expr: ex}, nil
	}

	return nil, fmt.Errorf("unexpected statement token %s (%q) at %d:%d", p.cur.Type, p.cur.Lexeme, p.cur.Line, p.cur.Column)
}

func (p *Parser) parseArgList() ([]ast.Expression, error) {
	var args []ast.Expression
	ex, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	args = append(args, ex)

	for p.match(token.COMMA) {
		ex2, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		args = append(args, ex2)
	}
	return args, nil
}

func (p *Parser) parseExpression() (ast.Expression, error) {
	return p.parseLogicOr()
}

func (p *Parser) parseLogicOr() (ast.Expression, error) {
	left, err := p.parseLogicAnd()
	if err != nil {
		return nil, err
	}
	for p.match(token.K_OR) {
		right, err := p.parseLogicAnd()
		if err != nil {
			return nil, err
		}
		left = ast.LogicOrExpr{Left: left, Right: right}
	}
	return left, nil
}

func (p *Parser) parseLogicAnd() (ast.Expression, error) {
	left, err := p.parseComparison()
	if err != nil {
		return nil, err
	}
	for p.match(token.K_AND) {
		right, err := p.parseComparison()
		if err != nil {
			return nil, err
		}
		left = ast.LogicAndExpr{Left: left, Right: right}
	}
	return left, nil
}

func (p *Parser) parseComparison() (ast.Expression, error) {
	left, err := p.parseArithmetic()
	if err != nil {
		return nil, err
	}
	switch p.cur.Type {
	case token.EQEQ, token.NEQ, token.GT, token.LT, token.GTE, token.LTE:
		opTok := p.cur
		_ = p.advance()
		right, err := p.parseArithmetic()
		if err != nil {
			return nil, err
		}
		return ast.CompareExpr{Left: left, Op: opTok.Lexeme, Right: right}, nil
	}
	return left, nil
}

func (p *Parser) parseArithmetic() (ast.Expression, error) {
	left, err := p.parseTerm()
	if err != nil {
		return nil, err
	}
	for p.cur.Type == token.PLUS || p.cur.Type == token.MINUS {
		op := p.cur.Lexeme
		_ = p.advance()
		right, err := p.parseTerm()
		if err != nil {
			return nil, err
		}
		left = ast.BinaryExpr{Left: left, Op: op, Right: right}
	}
	return left, nil
}

func (p *Parser) parseTerm() (ast.Expression, error) {
	left, err := p.parseFactor()
	if err != nil {
		return nil, err
	}
	for p.cur.Type == token.ASTERISK || p.cur.Type == token.SLASH {
		op := p.cur.Lexeme
		_ = p.advance()
		right, err := p.parseFactor()
		if err != nil {
			return nil, err
		}
		left = ast.BinaryExpr{Left: left, Op: op, Right: right}
	}
	return left, nil
}

func (p *Parser) parseFactor() (ast.Expression, error) {
	left, err := p.parseUnary()
	if err != nil {
		return nil, err
	}
	if p.cur.Type == token.POW {
		_ = p.advance()
		right, err := p.parseFactor()
		if err != nil {
			return nil, err
		}
		return ast.PowerExpr{Left: left, Right: right}, nil
	}
	return left, nil
}

func (p *Parser) parseUnary() (ast.Expression, error) {
	if p.cur.Type == token.MINUS {
		_ = p.advance()
		x, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		return ast.UnaryExpr{Op: "-", X: x}, nil
	}
	if p.cur.Type == token.BANG {
		_ = p.advance()
		x, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		return ast.UnaryExpr{Op: "!", X: x}, nil
	}
	return p.parsePrimary()
}

func (p *Parser) parsePrimary() (ast.Expression, error) {
	switch p.cur.Type {
	case token.NUMBER:
		tok := p.cur
		_ = p.advance()

		var v float64
		if _, err := fmt.Sscanf(tok.Lexeme, "%f", &v); err != nil {
			return nil, fmt.Errorf("invalid number %q at %d:%d", tok.Lexeme, tok.Line, tok.Column)
		}
		return ast.NumberLit{Value: v}, nil

	case token.STRING:
		tok := p.cur
		_ = p.advance()
		return ast.StringLit{Value: tok.Lexeme}, nil

	case token.BOOLEAN:
		tok := p.cur
		_ = p.advance()
		return ast.BooleanLit{Value: tok.Lexeme == "true"}, nil

	case token.IDENT:
		identTok := p.cur
		_ = p.advance()
		recv := ast.IdentifierExpr{Name: identTok.Lexeme}

		if p.cur.Type == token.DOT {
			_ = p.advance()
			methodTok, err := p.expect(token.IDENT)
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(token.LPAREN); err != nil {
				return nil, err
			}

			var args []ast.Expression
			if p.cur.Type != token.RPAREN {
				argsParsed, err := p.parseArgList()
				if err != nil {
					return nil, err
				}
				args = argsParsed
			}
			if _, err := p.expect(token.RPAREN); err != nil {
				return nil, err
			}
			return ast.MethodCallExpr{Receiver: recv, Method: methodTok.Lexeme, Args: args}, nil
		}

		if p.match(token.LPAREN) {
			var args []ast.Expression
			if p.cur.Type != token.RPAREN {
				argsParsed, err := p.parseArgList()
				if err != nil {
					return nil, err
				}
				args = argsParsed
			}
			if _, err := p.expect(token.RPAREN); err != nil {
				return nil, err
			}
			return ast.CallExpr{Callee: identTok.Lexeme, Args: args}, nil
		}

		return ast.IdentifierExpr{Name: identTok.Lexeme}, nil

	case token.LPAREN:
		_ = p.advance()
		ex, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(token.RPAREN); err != nil {
			return nil, err
		}
		return ast.ParenExpr{X: ex}, nil

	case token.LBRACK:
		_ = p.advance()
		var elems []ast.Expression
		if p.cur.Type != token.RBRACK {
			ex, err := p.parseExpression()
			if err != nil {
				return nil, err
			}
			elems = append(elems, ex)
			for p.match(token.COMMA) {
				ex2, err := p.parseExpression()
				if err != nil {
					return nil, err
				}
				elems = append(elems, ex2)
			}
		}
		if _, err := p.expect(token.RBRACK); err != nil {
			return nil, err
		}
		return ast.ArrayLit{Elems: elems}, nil
	}

	return nil, fmt.Errorf("unexpected primary token %s (%q) at %d:%d", p.cur.Type, p.cur.Lexeme, p.cur.Line, p.cur.Column)
}

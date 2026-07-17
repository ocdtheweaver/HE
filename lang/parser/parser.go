package parser

import (
	"fmt"

	"strings"

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
		p.cur = token.Token{Type: token.ILLEGAL, Line: 1, Column: 1}
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
		return token.Token{}, fmt.Errorf(
			"line %d:%d — expected %q but got %q (%s)",
			p.cur.Line, p.cur.Column, tt, p.cur.Lexeme, p.cur.Type,
		)
	}
	t := p.cur
	_ = p.advance()
	return t, nil
}

func (p *Parser) match(tt token.TokenType) bool {
	if p.cur.Type != tt {
		return false
	}
	_ = p.advance()
	return true
}

// isStatementStart returns true if the current token can begin a statement.
func (p *Parser) isStatementStart() bool {
	switch p.cur.Type {
	case token.K_PRINT, token.K_SAY, token.K_SHOW,
		token.K_SET, token.K_LET, token.K_CHANGE, token.K_MAKE,
		token.K_GROW, token.K_SHRINK,
		token.K_IF, token.K_CHECK,
		token.K_REPEAT, token.K_WHILE,
		token.K_RETURN,
		token.K_WAIT,
		token.K_CAN_TELL, token.K_FOR, token.K_ASK, token.K_GIVE, token.K_TRY,
		token.K_REMEMBER, token.K_FORGET, token.K_SUMMON, token.K_EACH,
		token.K_WITH, token.K_FIELDS, token.STRING, token.INTERP_STRING, token.NUMBER,
		token.K_BECOMES,
		token.IDENT:
		return true
	}
	return false
}

// ── Top-level ─────────────────────────────────────────────────────────────────

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
		// "make" could be an object creator OR a set statement at top level.
		// If followed by IDENT then '{' or '[', it's an object. Otherwise statement.
		if p.cur.Type == token.K_MAKE {
			next, err := p.lx.Peek()
			if err != nil {
				return nil, err
			}
			if next.Type == token.IDENT {
				// Peek one more ahead
				// We do a speculative parse: save position isn't available, so
				// we check if next-next looks like a block opener. Since we only
				// have one lookahead, we treat MAKE at top-level as an object
				// declaration. Inside statements, K_MAKE is handled separately.
				return p.parseObjectLine()
			}
		}
		return p.parseObjectLine()
	}

	// Everything else is a statement at global scope.
	if p.isStatementStart() {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		return &ast.GlobalStatementLine{Statement: stmt}, nil
	}

	return nil, fmt.Errorf(
		"line %d:%d — unexpected token %q (%s) at top level",
		p.cur.Line, p.cur.Column, p.cur.Lexeme, p.cur.Type,
	)
}

// ── Summon ────────────────────────────────────────────────────────────────────

func (p *Parser) parseSummon() (ast.Line, error) {
	if _, err := p.expect(token.K_SUMMON); err != nil {
		return nil, err
	}
	modTok, err := p.expect(token.STRING)
	if err != nil {
		return nil, err
	}

	alias := ""
	hasAlias := false
	// "as"/"named" is now optional — flat-merge summons skip it entirely.
	if p.cur.Type == token.K_NAMED || p.cur.Type == token.K_AS {
		_ = p.advance()
		aliasTok, err := p.expect(token.IDENT)
		if err != nil {
			return nil, err
		}
		alias = aliasTok.Lexeme
		hasAlias = true
	}

	// If module name ends in .he — it's a file module (always flat-merges
	// its contents into the importing scope; alias just adds a qualified
	// path on top, it never gates access).
	if strings.HasSuffix(modTok.Lexeme, ".he") {
		return &ast.GlobalStatementLine{
			Statement: ast.LoadModuleStmt{FilePath: modTok.Lexeme, Alias: alias, HasAlias: hasAlias},
		}, nil
	}

	// Built-in stdlib modules (math, text, list, ...) still require an
	// alias today — they aren't backed by HE source to flat-merge from.
	if !hasAlias {
		return nil, fmt.Errorf(
			"line %d:%d — built-in module %q needs 'as <name>' (e.g. summon %q as m) — only .he file summons can omit it",
			p.cur.Line, p.cur.Column, modTok.Lexeme, modTok.Lexeme,
		)
	}
	return ast.SummonLine{
		ModuleName: modTok,
		Alias:      &ast.Identifier{Name: alias},
	}, nil
}

// ── Object ────────────────────────────────────────────────────────────────────

func (p *Parser) parseObjectLine() (*ast.ObjectLine, error) {
	kindTok := p.cur
	_ = p.advance()

	nameTok, err := p.expect(token.IDENT)
	if err != nil {
		return nil, err
	}
	obj := &ast.ObjectLine{
		Kind: kindTok.Lexeme,
		Name: ast.Identifier{Name: nameTok.Lexeme, Tok: nameTok},
	}

	// Optional protection tag: make Car3 #protected [...]
	if p.cur.Type == token.PROTECT_TAG {
		obj.ProtectTag = p.cur.Lexeme
		_ = p.advance()
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
		return nil, fmt.Errorf(
			"line %d:%d — expected '{' or '[' to open object body, got %q",
			p.cur.Line, p.cur.Column, p.cur.Lexeme,
		)
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
		switch p.cur.Type {
		case token.K_HAS, token.K_OWNS, token.K_CARRIES:
			sectionKind := p.cur.Lexeme
			_ = p.advance()
			_ = p.match(token.COLON)
			if _, err := p.expect(token.LBRACK); err != nil {
				return body, err
			}
			props, err := p.parsePropertiesSection(sectionKind)
			if err != nil {
				return body, err
			}
			body.Sections = append(body.Sections, props)

		case token.K_CAN, token.K_KNOW:
			_ = p.advance()
			// "know how to" is an alias for "can"
			if p.cur.Type == token.K_HOW {
				_ = p.advance() // consume "how"
				_ = p.match(token.K_TO) // consume optional "to"
			}
			_ = p.match(token.COLON)
			if _, err := p.expect(token.LBRACK); err != nil {
				return body, err
			}
			abilities, err := p.parseAbilitiesList()
			if err != nil {
				return body, err
			}
			if _, err := p.expect(token.RBRACK); err != nil {
				return body, err
			}
			body.Sections = append(body.Sections, abilities)

		case token.K_ON, token.K_WHEN, token.K_WHENEVER:
			react, err := p.parseReactionsSection()
			if err != nil {
				return body, err
			}
			body.Sections = append(body.Sections, react)

		case token.K_REMEMBERS:
			_ = p.advance()
			_ = p.match(token.COLON)
			// Skip memories section body for now — consume until matching ]
			if p.cur.Type == token.LBRACK {
				_ = p.advance()
				depth := 1
				for depth > 0 && p.cur.Type != token.EOF {
					if p.cur.Type == token.LBRACK {
						depth++
					} else if p.cur.Type == token.RBRACK {
						depth--
					}
					_ = p.advance()
				}
			}
			body.Sections = append(body.Sections, ast.MemoriesSection{})

		case token.K_CREATE, token.K_MAKE:
			obj, err := p.parseObjectLine()
			if err != nil {
				return body, err
			}
			body.NestedObjects = append(body.NestedObjects, *obj)

		default:
			return body, fmt.Errorf(
				"line %d:%d — unexpected %q inside object body",
				p.cur.Line, p.cur.Column, p.cur.Lexeme,
			)
		}
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

		// Optional type annotation: name: typehint is value
		if p.cur.Type == token.COLON {
			_ = p.advance() // consume ":"
			_ = p.advance() // consume type name (ignore for now, stored as hint)
		}

		kind := ""
		switch {
		case p.match(token.K_IS):
			kind = "is"
		case p.match(token.K_BECOMES):
			kind = "is"
		case p.cur.Type == token.IDENT && p.cur.Lexeme == "starts":
			_ = p.advance()
			if p.cur.Type == token.K_AS || p.cur.Type == token.K_BE {
				_ = p.advance()
			}
			kind = "starts as"
		default:
			return ast.PropertiesSection{}, fmt.Errorf(
				"line %d:%d — expected 'is', 'becomes', or 'starts as' after property name %q",
				p.cur.Line, p.cur.Column, nameTok.Lexeme,
			)
		}

		ex, err := p.parseExpression()
		if err != nil {
			return ast.PropertiesSection{}, err
		}
		_ = p.match(token.COMMA)

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

		// Optional protection tag right after the name: past20Seconds #protected1 [...]
		if p.cur.Type == token.PROTECT_TAG {
			action.ProtectTag = p.cur.Lexeme
			_ = p.advance()
		}

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

		// Also allow the tag after params, for "past20Seconds(arg) #protected1 [...]"
		if p.cur.Type == token.PROTECT_TAG {
			action.ProtectTag = p.cur.Lexeme
			_ = p.advance()
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
		if p.match(token.COLON) {
			typ, err := p.parseType()
			if err != nil {
				return nil, err
			}
			params = append(params, ast.Param{Name: nameTok.Lexeme, Type: typ})
		} else {
			params = append(params, ast.Param{Name: nameTok.Lexeme, Type: ast.TypeNode{Name: "any"}})
		}
		if !p.match(token.COMMA) {
			break
		}
	}
	return params, nil
}

func (p *Parser) parseType() (ast.TypeNode, error) {
	if p.cur.Type != token.IDENT {
		return ast.TypeNode{}, fmt.Errorf(
			"line %d:%d — expected type name, got %q",
			p.cur.Line, p.cur.Column, p.cur.Lexeme,
		)
	}
	name := p.cur.Lexeme
	_ = p.advance()

	// Only treat as array type if we see [] (open then immediately close bracket).
	// Peek ahead: if next is LBRACK and the one after is RBRACK, it's a type modifier.
	if p.cur.Type == token.LBRACK {
		next, err := p.lx.Peek()
		if err == nil && next.Type == token.RBRACK {
			_ = p.advance() // consume [
			_ = p.advance() // consume ]
			name = name + "[]"
		}
		// Otherwise leave [ alone — it's the action body opener
	}
	return ast.TypeNode{Name: name}, nil
}

// ── Reactions ─────────────────────────────────────────────────────────────────

func (p *Parser) parseReactionsSection() (ast.ReactionsSection, error) {
	_ = p.advance() // consume on/when/whenever
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
	return ast.ReactionsSection{
		Reactions: []ast.ReactionDef{{Trigger: trigger, Body: body}},
	}, nil
}

func (p *Parser) parseTrigger() (ast.Trigger, error) {
	// Accept any identifier or collision keyword as the trigger name
	if p.cur.Type != token.IDENT && p.cur.Type != token.K_COLLISION {
		return ast.Trigger{}, fmt.Errorf(
			"line %d:%d — expected event name, got %q",
			p.cur.Line, p.cur.Column, p.cur.Lexeme,
		)
	}
	leftTok := p.cur
	_ = p.advance()
	tr := ast.Trigger{Left: leftTok.Lexeme}

	if p.match(token.K_WITH) {
		switch p.cur.Type {
		case token.IDENT:
			s := p.cur.Lexeme
			_ = p.advance()
			tr.Right = &s
			tr.RightIsStr = false
		case token.STRING:
			s := p.cur.Lexeme
			_ = p.advance()
			tr.Right = &s
			tr.RightIsStr = true
		default:
			return ast.Trigger{}, fmt.Errorf(
				"line %d:%d — expected identifier or string after 'with' in trigger",
				p.cur.Line, p.cur.Column,
			)
		}
	}
	return tr, nil
}

// ── Statements ────────────────────────────────────────────────────────────────

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

	// ── Output ──
	case token.K_PRINT, token.K_SAY, token.K_SHOW:
		_ = p.advance()
		ex, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		return ast.SayStmt{Expr: ex}, nil

	// ── Assignment: set X to Y | set X, Y to ... | set X.field to Y ──
	case token.K_SET, token.K_LET:
		_ = p.advance()
		nameTok, err := p.expect(token.IDENT)
		if err != nil {
			return nil, err
		}
		// dot-assign: set obj.field to expr
		if p.cur.Type == token.DOT {
			return p.parseDotAssign(nameTok.Lexeme)
		}
		// multi-assign: set a, b to expr1, expr2
		if p.cur.Type == token.COMMA {
			names := []string{nameTok.Lexeme}
			for p.match(token.COMMA) {
				ntok, err := p.expect(token.IDENT)
				if err != nil {
					return nil, err
				}
				names = append(names, ntok.Lexeme)
			}
			if p.cur.Type != token.K_TO && p.cur.Type != token.K_BE && p.cur.Type != token.K_IS {
				return nil, fmt.Errorf("line %d:%d — expected 'to' after names in multi-assign", p.cur.Line, p.cur.Column)
			}
			_ = p.advance()
			var exprs []ast.Expression
			first, err := p.parseExpression()
			if err != nil {
				return nil, err
			}
			exprs = append(exprs, first)
			for p.match(token.COMMA) {
				ex2, err := p.parseExpression()
				if err != nil {
					return nil, err
				}
				exprs = append(exprs, ex2)
			}
			return ast.MultiAssignStmt{Names: names, Exprs: exprs}, nil
		}
		// accept "to", "be", "is", "becomes"
		if p.cur.Type != token.K_TO && p.cur.Type != token.K_BE &&
			p.cur.Type != token.K_IS && p.cur.Type != token.K_BECOMES {
			return nil, fmt.Errorf(
				"line %d:%d — expected 'to', 'be', or 'is' after variable name",
				p.cur.Line, p.cur.Column,
			)
		}
		_ = p.advance()
		ex, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		return ast.ChangeStmt{Kind: "set", Name: nameTok.Lexeme, Expr: ex}, nil

	// ── Change: change X to Y ──
	case token.K_CHANGE:
		_ = p.advance()
		nameTok, err := p.expect(token.IDENT)
		if err != nil {
			return nil, err
		}
		if p.cur.Type != token.K_TO && p.cur.Type != token.K_IS && p.cur.Type != token.K_BECOMES {
			return nil, fmt.Errorf(
				"line %d:%d — expected 'to' after 'change %s'",
				p.cur.Line, p.cur.Column, nameTok.Lexeme,
			)
		}
		_ = p.advance()
		ex, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		return ast.ChangeStmt{Kind: "change", Name: nameTok.Lexeme, Expr: ex}, nil

	// ── Make (as statement): make X to Y | make X be Y ──
	case token.K_MAKE:
		_ = p.advance()
		nameTok, err := p.expect(token.IDENT)
		if err != nil {
			return nil, err
		}
		if p.cur.Type == token.K_TO || p.cur.Type == token.K_BE ||
			p.cur.Type == token.K_IS || p.cur.Type == token.K_BECOMES {
			_ = p.advance()
			ex, err := p.parseExpression()
			if err != nil {
				return nil, err
			}
			return ast.ChangeStmt{Kind: "set", Name: nameTok.Lexeme, Expr: ex}, nil
		}
		// Otherwise it might be "make X like Y [...]" — object creation as statement
		// Put cur back isn't possible with our simple lexer, so error gracefully.
		return nil, fmt.Errorf(
			"line %d:%d — 'make' as a statement expects 'to', 'be', or 'is'",
			p.cur.Line, p.cur.Column,
		)

	// ── grow X by N | grow obj.field by N ──
	case token.K_GROW:
		_ = p.advance()
		nameTok, err := p.expect(token.IDENT)
		if err != nil {
			return nil, err
		}
		// dot form: grow obj.field by N  →  DotAssign with binary add
		if p.cur.Type == token.DOT {
			_ = p.advance()
			fieldTok, err := p.expect(token.IDENT)
			if err != nil {
				return nil, err
			}
			if p.cur.Type == token.IDENT && p.cur.Lexeme == "by" {
				_ = p.advance()
			}
			delta, err := p.parseExpression()
			if err != nil {
				return nil, err
			}
			// build: obj.field + delta
			addExpr := ast.BinaryExpr{
				Left: ast.FieldAccessExpr{Receiver: ast.IdentifierExpr{Name: nameTok.Lexeme}, Field: fieldTok.Lexeme},
				Op: "+",
				Right: delta,
			}
			return ast.DotAssignStmt{Object: nameTok.Lexeme, Field: fieldTok.Lexeme, Expr: addExpr}, nil
		}
		if p.cur.Type == token.IDENT && p.cur.Lexeme == "by" {
			_ = p.advance()
		}
		ex, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		return ast.GrowStmt{Name: nameTok.Lexeme, Expr: ex}, nil

	// ── shrink X by N ──
	case token.K_SHRINK:
		_ = p.advance()
		nameTok, err := p.expect(token.IDENT)
		if err != nil {
			return nil, err
		}
		if p.cur.Type == token.IDENT && p.cur.Lexeme == "by" {
			_ = p.advance()
		}
		ex, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		return ast.ShrinkStmt{Name: nameTok.Lexeme, Expr: ex}, nil

	// ── if / check if ──
	case token.K_IF, token.K_CHECK:
		_ = p.advance()
		// "check if" — consume optional "if"
		if p.cur.Type == token.K_IF {
			_ = p.advance()
		}
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

	// ── repeat while cond | repeat N times | repeat until cond ──
	case token.K_REPEAT:
		_ = p.advance()
		if p.match(token.K_UNTIL) {
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
			return ast.RepeatUntilStmt{Cond: cond, Body: body}, nil
		}
		if p.match(token.K_WHILE) {
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
			return ast.RepeatStmt{Kind: "while", Cond: cond, Body: body}, nil
		}
		// repeat N times [as i] [...]
		countExpr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		_ = p.match(token.K_TIMES) // consume optional "times"
		counterVar := ""
		if p.cur.Type == token.K_AS {
			_ = p.advance()
			ctrTok, err := p.expect(token.IDENT)
			if err != nil {
				return nil, err
			}
			counterVar = ctrTok.Lexeme
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
		if counterVar != "" {
			return ast.CountedRepeatStmt{Count: countExpr, CountVar: counterVar, Body: body}, nil
		}
		return ast.RepeatStmt{Kind: "times", Cond: countExpr, Body: body}, nil

	// ── while cond ──
	case token.K_WHILE:
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
		return ast.RepeatStmt{Kind: "while", Cond: cond, Body: body}, nil

	// ── return [a, b, c] ──
	case token.K_RETURN:
		_ = p.advance()
		first, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		// multiple return: return a, b, c
		if p.cur.Type == token.COMMA {
			exprs := []ast.Expression{first}
			for p.match(token.COMMA) {
				ex2, err := p.parseExpression()
				if err != nil {
					return nil, err
				}
				exprs = append(exprs, ex2)
			}
			return ast.MultiReturnStmt{Exprs: exprs}, nil
		}
		return ast.ReturnStmt{Expr: first}, nil

	// ── wait N seconds | wait N frames ──
	case token.K_WAIT:
		_ = p.advance()
		ex, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		unit := ""
		switch p.cur.Type {
		case token.K_SECONDS:
			unit = "seconds"
			_ = p.advance()
		case token.K_FRAMES:
			unit = "frames"
			_ = p.advance()
		default:
			return nil, fmt.Errorf(
				"line %d:%d — expected 'seconds' or 'frames' after wait",
				p.cur.Line, p.cur.Column,
			)
		}
		return ast.WaitStmt{Expr: ex, Unit: unit}, nil

	// ── tell X to action [with args] ──
	case token.K_CAN_TELL:
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





	// ── summon inside a method body ──
	case token.K_SUMMON:
		line, err := p.parseSummon()
		if err != nil {
			return nil, err
		}
		switch sl := line.(type) {
		case ast.SummonLine:
			return ast.ExprStmt{Expr: ast.StringLit{Value: "__summon__:" + sl.ModuleName.Lexeme + ":" + sl.Alias.Name}}, nil
		case *ast.GlobalStatementLine:
			if ld, ok := sl.Statement.(ast.LoadModuleStmt); ok {
				return ld, nil
			}
		}
		return nil, fmt.Errorf("line %d:%d — invalid summon", p.cur.Line, p.cur.Column)

	// ── remember name / forget name / recall name ──
	case token.K_REMEMBER:
		_ = p.advance()
		nameTok, err := p.expect(token.IDENT)
		if err != nil {
			return nil, err
		}
		key := nameTok.Lexeme
		if p.cur.Type == token.K_AS {
			_ = p.advance()
			keyTok, err := p.expect(token.STRING)
			if err != nil {
				return nil, err
			}
			key = keyTok.Lexeme
		}
		return ast.RememberStmt{Name: nameTok.Lexeme, Key: key}, nil

	case token.K_FORGET:
		_ = p.advance()
		nameTok, err := p.expect(token.IDENT)
		if err != nil {
			return nil, err
		}
		return ast.ForgetStmt{Name: nameTok.Lexeme}, nil

	// ── try [...] or [...] / try [...] or (errVar) [...] ──
	case token.K_TRY:
		_ = p.advance()
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
		if p.cur.Type != token.K_OR {
			return ast.TryStmt{Body: body, Handler: nil}, nil
		}
		_ = p.advance() // consume "or"
		// Named error capture: or (errVar) [...]
		if p.cur.Type == token.LPAREN {
			_ = p.advance()
			errTok, err := p.expect(token.IDENT)
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(token.RPAREN); err != nil {
				return nil, err
			}
			if _, err := p.expect(token.LBRACK); err != nil {
				return nil, err
			}
			handler, err := p.parseStatementListUntil(token.RBRACK)
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(token.RBRACK); err != nil {
				return nil, err
			}
			return ast.TryWithVarStmt{Body: body, ErrVar: errTok.Lexeme, Handler: handler}, nil
		}
		// Anonymous: or [...] / or if it fails [...]
		for p.cur.Type == token.K_IF || (p.cur.Type == token.IDENT && (p.cur.Lexeme == "it" || p.cur.Lexeme == "something")) || p.cur.Type == token.K_FAILS {
			_ = p.advance()
		}
		if _, err := p.expect(token.LBRACK); err != nil {
			return nil, err
		}
		handler, err := p.parseStatementListUntil(token.RBRACK)
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(token.RBRACK); err != nil {
			return nil, err
		}
		return ast.TryStmt{Body: body, Handler: handler}, nil

	// ── repeat until cond ──
	case token.K_UNTIL:
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
		return ast.RepeatUntilStmt{Cond: cond, Body: body}, nil



	// ── with expr as alias [ body ] ──
	case token.K_WITH:
		_ = p.advance()
		// Could be: with obj as alias [...] or named-arg block context
		// Detect: with IDENT/expr as IDENT [...]
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if p.cur.Type != token.K_AS {
			// Not a with-scope — treat as tell-with (shouldn't reach here normally)
			return nil, fmt.Errorf("line %d:%d — expected 'as' after 'with expr'", p.cur.Line, p.cur.Column)
		}
		_ = p.advance()
		aliasTok, err := p.expect(token.IDENT)
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
		return ast.WithScopeStmt{Alias: aliasTok.Lexeme, Expr: expr, Body: body}, nil

	// ── each N from A to B [...] — short range loop (no "for" prefix needed) ──
	case token.K_EACH:
		_ = p.advance()
		varTok, err := p.expect(token.IDENT)
		if err != nil {
			return nil, err
		}
		if p.cur.Type != token.K_FROM {
			return nil, fmt.Errorf("line %d:%d — expected 'from' after variable in 'each'", p.cur.Line, p.cur.Column)
		}
		_ = p.advance()
		fromExpr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if p.cur.Type != token.K_TO && p.cur.Type != token.K_UPTO {
			return nil, fmt.Errorf("line %d:%d — expected 'to' in 'each ... from A to B'", p.cur.Line, p.cur.Column)
		}
		_ = p.advance()
		toExpr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		var stepExpr ast.Expression
		if p.cur.Type == token.K_STEP {
			_ = p.advance()
			stepExpr, err = p.parseExpression()
			if err != nil {
				return nil, err
			}
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
		return ast.RangeLoopStmt{VarName: varTok.Lexeme, From: fromExpr, To: toExpr, Step: stepExpr, Body: body}, nil

	// ── for each item in list or range ──
	case token.K_FOR:
		return p.parseForEach()

	// ── ask "prompt" [as varName] ──
	case token.K_ASK:
		return p.parseAsk()

	// ── give expr to obj.field  ──
	case token.K_GIVE:
		_ = p.advance()
		valExpr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(token.K_TO); err != nil {
			return nil, err
		}
		// expect obj.field
		objTok, err := p.expect(token.IDENT)
		if err != nil {
			return nil, err
		}
		if p.cur.Type != token.DOT {
			return nil, fmt.Errorf("line %d:%d — expected obj.field after 'give ... to'", p.cur.Line, p.cur.Column)
		}
		_ = p.advance()
		fieldTok, err := p.expect(token.IDENT)
		if err != nil {
			return nil, err
		}
		return ast.DotAssignStmt{Object: objTok.Lexeme, Field: fieldTok.Lexeme, Expr: valExpr}, nil

	// ── STRING/NUMBER/INTERP — expression statement (e.g. "hello".upper()) ──
	case token.STRING, token.INTERP_STRING, token.NUMBER:
		ex, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		return ast.ExprStmt{Expr: ex}, nil

	// ── IDENT — expression statement, dot-assign, or dot-call ──
	case token.IDENT:
		// "name becomes value" — mutation shorthand (same as set name to value)
		if p.cur.Type == token.IDENT {
			next, _ := p.lx.Peek()
			if next.Type == token.K_BECOMES {
				nameTok := p.cur
				_ = p.advance() // consume name
				_ = p.advance() // consume "becomes"
				ex, err := p.parseExpression()
				if err != nil {
					return nil, err
				}
				return ast.ChangeStmt{Kind: "set", Name: nameTok.Lexeme, Expr: ex}, nil
			}
		}
		// "recall name [as alias]" — load from persistent store
		if p.cur.Lexeme == "recall" {
			_ = p.advance()
			nameTok, err := p.expect(token.IDENT)
			if err != nil {
				return nil, err
			}
			key := nameTok.Lexeme
			alias := nameTok.Lexeme
			if p.cur.Type == token.K_AS {
				_ = p.advance()
				askeyTok, err := p.expect(token.IDENT)
				if err != nil {
					return nil, err
				}
				alias = askeyTok.Lexeme
			}
			return ast.RecallStmt{Name: alias, Key: key}, nil
		}
		ex, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		return ast.ExprStmt{Expr: ex}, nil
	}

	return nil, fmt.Errorf(
		"line %d:%d — unexpected %q (%s) — not a valid statement",
		p.cur.Line, p.cur.Column, p.cur.Lexeme, p.cur.Type,
	)
}

// ── Arguments ─────────────────────────────────────────────────────────────────

func (p *Parser) parseArgList() ([]ast.Expression, error) {
	ex, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	args := []ast.Expression{ex}
	for p.match(token.COMMA) {
		ex2, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		args = append(args, ex2)
	}
	return args, nil
}

// ── Expressions ───────────────────────────────────────────────────────────────

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
	case token.K_IS:
		// "X is one of [...]" | "X is between A and B" | "X is not Y" | "X is Y"
		_ = p.advance()
		// "is one of [list]"
		if p.cur.Type == token.K_ONE {
			_ = p.advance() // consume "one"
			if p.cur.Type == token.K_OF {
				_ = p.advance() // consume "of"
			}
			listExpr, err := p.parseArithmetic()
			if err != nil {
				return nil, err
			}
			return ast.MembershipExpr{Value: left, List: listExpr}, nil
		}
		if p.cur.Type == token.K_BETWEEN {
			_ = p.advance()
			low, err := p.parseArithmetic()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(token.K_AND); err != nil {
				return nil, fmt.Errorf("line %d:%d — expected 'and' in 'is between A and B'", p.cur.Line, p.cur.Column)
			}
			high, err := p.parseArithmetic()
			if err != nil {
				return nil, err
			}
			return ast.BetweenExpr{Value: left, Low: low, High: high}, nil
		}
		op := "=="
		if p.cur.Type == token.K_NOT {
			_ = p.advance()
			op = "!="
		}
		right, err := p.parseArithmetic()
		if err != nil {
			return nil, err
		}
		return ast.CompareExpr{Left: left, Op: op, Right: right}, nil
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
	if p.cur.Type == token.POW || p.cur.Type == token.CARET {
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
	if p.cur.Type == token.BANG || p.cur.Type == token.K_NOT {
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
			return nil, fmt.Errorf("line %d:%d — invalid number %q", tok.Line, tok.Column, tok.Lexeme)
		}
		return p.chainExpr(ast.NumberLit{Value: v})

	case token.STRING:
		tok := p.cur
		_ = p.advance()
		lit := ast.Expression(ast.StringLit{Value: tok.Lexeme})
		return p.chainExpr(lit)

	case token.INTERP_STRING:
		tok := p.cur
		_ = p.advance()
		return p.parseInterpString(tok.Lexeme), nil

	case token.BOOLEAN:
		tok := p.cur
		_ = p.advance()
		return ast.BooleanLit{Value: tok.Lexeme == "true"}, nil

	case token.IDENT:
		identTok := p.cur
		_ = p.advance()
		recv := ast.IdentifierExpr{Name: identTok.Lexeme}

		// dot access: obj.method(args) or obj.field
		if p.cur.Type == token.DOT {
			_ = p.advance()
			memberTok, err := p.expect(token.IDENT)
			if err != nil {
				return nil, err
			}
			if p.match(token.LPAREN) {
				var args []ast.Expression
				if p.cur.Type != token.RPAREN {
					args, err = p.parseArgList()
					if err != nil {
						return nil, err
					}
				}
				if _, err := p.expect(token.RPAREN); err != nil {
					return nil, err
				}
				callExpr := ast.Expression(ast.MethodCallExpr{Receiver: recv, Method: memberTok.Lexeme, Args: args})
				return p.chainExpr(callExpr)
			}
			// field access (no parens)
			return ast.FieldAccessExpr{Receiver: recv, Field: memberTok.Lexeme}, nil
		}

		// function call: name(args)
		if p.match(token.LPAREN) {
			var args []ast.Expression
			var err error
			if p.cur.Type != token.RPAREN {
				args, err = p.parseArgList()
				if err != nil {
					return nil, err
				}
			}
			if _, err := p.expect(token.RPAREN); err != nil {
				return nil, err
			}
			return ast.CallExpr{Callee: identTok.Lexeme, Args: args}, nil
		}

		return recv, nil

	case token.LPAREN:
		_ = p.advance()
		ex, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(token.RPAREN); err != nil {
			return nil, err
		}
		return p.chainExpr(ast.ParenExpr{X: ex})

	case token.LBRACK:
		_ = p.advance()
		// Check if this is a named-arg block: [key is val, key: val, ...]
		// We detect by peeking: if next two tokens are IDENT (IS|COLON), it's named.
		if p.cur.Type == token.IDENT {
			next, err := p.lx.Peek()
			if err != nil {
				return nil, err
			}
			if next.Type == token.K_IS || next.Type == token.COLON {
				return p.parseNamedArgLit()
			}
		}
		// Regular array literal
		var elems []ast.Expression
		if p.cur.Type != token.RBRACK {
			ex, err := p.parseExpression()
			if err != nil {
				return nil, err
			}
			elems = append(elems, ex)
			for p.match(token.COMMA) {
				if p.cur.Type == token.RBRACK {
					break // trailing comma
				}
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
		return p.chainExpr(ast.ArrayLit{Elems: elems})

	case token.K_ABILITY:
		_ = p.advance()
		var params []string
		if p.match(token.LPAREN) {
			for p.cur.Type != token.RPAREN && p.cur.Type != token.EOF {
				ptok, err := p.expect(token.IDENT)
				if err != nil {
					return nil, err
				}
				params = append(params, ptok.Lexeme)
				_ = p.match(token.COMMA)
			}
			if _, err := p.expect(token.RPAREN); err != nil {
				return nil, err
			}
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
		// Emit ClosureExpr — runtime will capture current scope
		return ast.ClosureExpr{Params: params, Body: body, Captured: nil}, nil
	}

	return nil, fmt.Errorf(
		"line %d:%d — unexpected %q — expected a value or expression",
		p.cur.Line, p.cur.Column, p.cur.Lexeme,
	)
}


// parseInterpExpr parses a single expression from inside a {..} interpolation.
// Supports: name, obj.field, fn(args), obj.method(args), arithmetic, etc.
func (p *Parser) parseInterpExpr(raw string) ast.Expression {
	lx2 := lexer.New(raw)
	p2 := &Parser{lx: lx2}
	if t, err := p2.lx.NextToken(); err == nil {
		p2.cur = t
	}
	ex, err := p2.parseExpression()
	if err != nil || ex == nil {
		// Fallback: plain identifier
		return ast.IdentifierExpr{Name: strings.TrimSpace(raw)}
	}
	return ex
}


// chainExpr wraps an expression with any following .method() or .field suffixes
func (p *Parser) chainExpr(expr ast.Expression) (ast.Expression, error) {
	for p.cur.Type == token.DOT {
		_ = p.advance()
		methodTok, err := p.expect(token.IDENT)
		if err != nil {
			return nil, err
		}
		if p.match(token.LPAREN) {
			var args []ast.Expression
			if p.cur.Type != token.RPAREN {
				args, err = p.parseArgList()
				if err != nil {
					return nil, err
				}
			}
			if _, err := p.expect(token.RPAREN); err != nil {
				return nil, err
			}
			expr = ast.MethodChainExpr{Recv: expr, Method: methodTok.Lexeme, Args: args}
		} else {
			expr = ast.MethodChainExpr{Recv: expr, Method: methodTok.Lexeme, Args: nil}
		}
	}
	return expr, nil
}

// parseNamedArgLit parses [key is val, key: val, ...] — called after '[' consumed
func (p *Parser) parseNamedArgLit() (ast.NamedArgLit, error) {
	var pairs []ast.NamedPair
	for p.cur.Type != token.RBRACK && p.cur.Type != token.EOF {
		keyTok, err := p.expect(token.IDENT)
		if err != nil {
			return ast.NamedArgLit{}, err
		}
		// accept "is" or ":"
		if p.cur.Type != token.K_IS && p.cur.Type != token.COLON {
			return ast.NamedArgLit{}, fmt.Errorf(
				"line %d:%d — expected 'is' or ':' after key %q in named block",
				p.cur.Line, p.cur.Column, keyTok.Lexeme,
			)
		}
		_ = p.advance()
		val, err := p.parseExpression()
		if err != nil {
			return ast.NamedArgLit{}, err
		}
		_ = p.match(token.COMMA)
		pairs = append(pairs, ast.NamedPair{Key: keyTok.Lexeme, Value: val})
	}
	if _, err := p.expect(token.RBRACK); err != nil {
		return ast.NamedArgLit{}, err
	}
	return ast.NamedArgLit{Pairs: pairs}, nil
}

// ── New statements added in evolution pass ────────────────────────────────────

// parseForEach: for each varName in listExpr [ body ]
//               for each varName from N to M [step S] [ body ]
func (p *Parser) parseForEach() (ast.Statement, error) {
	// consume "for"
	_ = p.advance()
	// consume "each"
	if _, err := p.expect(token.K_EACH); err != nil {
		return nil, fmt.Errorf("line %d:%d — expected 'each' after 'for'", p.cur.Line, p.cur.Column)
	}
	varTok, err := p.expect(token.IDENT)
	if err != nil {
		return nil, err
	}

	// Range loop: for each i from N to M [step S]
	if p.cur.Type == token.K_FROM {
		_ = p.advance()
		fromExpr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		// accept "to" or "up to" or "upto"
		if p.cur.Type != token.K_TO && p.cur.Type != token.K_UPTO {
			return nil, fmt.Errorf("line %d:%d — expected 'to' after from-value", p.cur.Line, p.cur.Column)
		}
		_ = p.advance()
		toExpr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		var stepExpr ast.Expression
		if p.cur.Type == token.K_STEP {
			_ = p.advance()
			stepExpr, err = p.parseExpression()
			if err != nil {
				return nil, err
			}
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
		return ast.RangeLoopStmt{VarName: varTok.Lexeme, From: fromExpr, To: toExpr, Step: stepExpr, Body: body}, nil
	}

	// Key-value iteration: for each key, val in obj.fields
	if p.cur.Type == token.COMMA {
		_ = p.advance()
		valTok, err := p.expect(token.IDENT)
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(token.K_IN); err != nil {
			return nil, fmt.Errorf("line %d:%d — expected 'in' after key, val", p.cur.Line, p.cur.Column)
		}
		objExpr, err := p.parseExpression()
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
		return ast.ForEachFieldStmt{KeyVar: varTok.Lexeme, ValVar: valTok.Lexeme, Object: objExpr, Body: body}, nil
	}

	// List loop: for each item in list
	if _, err := p.expect(token.K_IN); err != nil {
		return nil, fmt.Errorf("line %d:%d — expected 'in' or 'from' after variable name in 'for each'", p.cur.Line, p.cur.Column)
	}
	listExpr, err := p.parseExpression()
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
	return ast.ForEachStmt{VarName: varTok.Lexeme, List: listExpr, Body: body}, nil
}

// parseAsk: ask "prompt" then set varName to answer
//           ask "prompt" storing result in varName
func (p *Parser) parseAsk() (ast.Statement, error) {
	_ = p.advance() // consume "ask"
	prompt, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	// accept "then set X to answer" or "storing result in X" or just "as X"
	varName := "__answer__"
	switch {
	case p.match(token.K_THEN):
		// "then set varName to answer"
		_ = p.match(token.K_SET) // optional "set"
		nameTok, err := p.expect(token.IDENT)
		if err != nil {
			return nil, err
		}
		_ = p.match(token.K_TO) // optional "to"
		varName = nameTok.Lexeme
	case p.cur.Type == token.K_AS:
		_ = p.advance()
		nameTok, err := p.expect(token.IDENT)
		if err != nil {
			return nil, err
		}
		varName = nameTok.Lexeme
	case p.cur.Type == token.IDENT && p.cur.Lexeme == "storing":
		_ = p.advance() // "storing"
		_ = p.match(token.IDENT) // skip "result" if present
		_ = p.match(token.K_IN)  // "in"
		nameTok, err := p.expect(token.IDENT)
		if err != nil {
			return nil, err
		}
		varName = nameTok.Lexeme
	}
	return ast.AskStmt{Prompt: prompt, VarName: varName}, nil
}

// parseDotAssign: set obj.field to expr
func (p *Parser) parseDotAssign(objName string) (ast.Statement, error) {
	// cur is DOT
	_ = p.advance() // consume "."
	fieldTok, err := p.expect(token.IDENT)
	if err != nil {
		return nil, err
	}
	if p.cur.Type != token.K_TO && p.cur.Type != token.K_IS &&
		p.cur.Type != token.K_BE && p.cur.Type != token.K_BECOMES {
		return nil, fmt.Errorf(
			"line %d:%d — expected 'to', 'is', or 'be' after '%s.%s'",
			p.cur.Line, p.cur.Column, objName, fieldTok.Lexeme,
		)
	}
	_ = p.advance()
	val, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	return ast.DotAssignStmt{Object: objName, Field: fieldTok.Lexeme, Expr: val}, nil
}

// ── Pass 3 parser helpers ─────────────────────────────────────────────────────

// parseInterpString splits "Hello, {name}! Score: {score}." into segments.
// Handles simple identifiers {name} and dot-access {obj.field}.
func (p *Parser) parseInterpString(raw string) ast.InterpStringExpr {
	var segs []ast.Expression
	i := 0
	for i < len(raw) {
		// Find next {
		j := i
		for j < len(raw) && raw[j] != '{' {
			j++
		}
		// Literal segment (always emit, even if empty, for correct alternation)
		segs = append(segs, ast.StringLit{Value: raw[i:j]})
		if j >= len(raw) {
			break
		}
		// Inside { ... }
		j++ // skip {
		k := j
		for k < len(raw) && raw[k] != '}' {
			k++
		}
		expr := strings.TrimSpace(raw[j:k])
		if expr != "" {
			// Parse the expression inside {} properly
			parsed := p.parseInterpExpr(expr)
			segs = append(segs, parsed)
		}
		i = k + 1 // skip }
	}
	return ast.InterpStringExpr{Segments: segs}
}

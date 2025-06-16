package ipx

import (
        "encoding/json"
        "errors"
        "fmt"
        "strings"
        "unicode"
)

// ─── Token Definitions ────────────────────────────────────────────────

type TokenType int

const (
        TOKEN_EOF TokenType = iota
        TOKEN_IDENT
        TOKEN_STRING
        TOKEN_AND
        TOKEN_OR
        TOKEN_NOT
        TOKEN_COLON
        TOKEN_ASSIGN
        TOKEN_LPAREN
        TOKEN_RPAREN
        TOKEN_LBRACKET
        TOKEN_RBRACKET
        TOKEN_COMMA
)

type Token struct {
        Type  TokenType
        Value string
        Start int
        End   int
}

// ─── Lexer ─────────────────────────────────────────────────────────────

type Lexer struct {
        input string
        pos   int
}

func NewLexer(input string) *Lexer {
        return &Lexer{input: input}
}

func (l *Lexer) next() rune {
        if l.pos >= len(l.input) {
                return 0
        }
        ch := rune(l.input[l.pos])
        l.pos++
        return ch
}

func (l *Lexer) peek() rune {
        if l.pos >= len(l.input) {
                return 0
        }
        return rune(l.input[l.pos])
}

func (l *Lexer) skipWhitespace() {
        for unicode.IsSpace(l.peek()) {
                l.next()
        }
}

func (l *Lexer) NextToken() (Token, error) {
        l.skipWhitespace()
        start := l.pos
        ch := l.peek()
        switch {
        case ch == 0:
                return Token{Type: TOKEN_EOF, Start: l.pos, End: l.pos}, nil
        case isLetter(ch):
                return l.lexIdent()
        case ch == '"':
                return l.lexString()
        case ch == '&':
                l.next()
                if l.peek() == '&' {
                        l.next()
                        return Token{Type: TOKEN_AND, Value: "&&", Start: start, End: l.pos}, nil
                }
                return Token{}, fmt.Errorf("unexpected character '&' at %d", start)
        case ch == '|':
                l.next()
                if l.peek() == '|' {
                        l.next()
                        return Token{Type: TOKEN_OR, Value: "||", Start: start, End: l.pos}, nil
                }
                return Token{}, fmt.Errorf("unexpected character '|' at %d", start)
        case ch == '!':
                l.next()
                return Token{Type: TOKEN_NOT, Value: "!", Start: start, End: l.pos}, nil
        case ch == ':':
                l.next()
                if l.peek() == '=' {
                        l.next()
                        return Token{Type: TOKEN_ASSIGN, Value: ":=", Start: start, End: l.pos}, nil
                }
                return Token{Type: TOKEN_COLON, Value: ":", Start: start, End: l.pos}, nil
        case ch == '(':
                l.next()
                return Token{Type: TOKEN_LPAREN, Value: "(", Start: start, End: l.pos}, nil
        case ch == ')':
                l.next()
                return Token{Type: TOKEN_RPAREN, Value: ")", Start: start, End: l.pos}, nil
        case ch == '[':
                l.next()
                return Token{Type: TOKEN_LBRACKET, Value: "[", Start: start, End: l.pos}, nil
        case ch == ']':
                l.next()
                return Token{Type: TOKEN_RBRACKET, Value: "]", Start: start, End: l.pos}, nil
        case ch == ',':
                l.next()
                return Token{Type: TOKEN_COMMA, Value: ",", Start: start, End: l.pos}, nil
        default:
                return Token{}, fmt.Errorf("unexpected character %q at position %d", ch, l.pos)
        }
}

func (l *Lexer) lexIdent() (Token, error) {
        start := l.pos
        for {
                ch := l.peek()
                if !isLetter(ch) && !unicode.IsDigit(ch) && ch != '.' && ch != '_' {
                        break
                }
                l.next()
        }
        return Token{Type: TOKEN_IDENT, Value: l.input[start:l.pos], Start: start, End: l.pos}, nil
}

func (l *Lexer) lexString() (Token, error) {
        start := l.pos
        l.next()
        strStart := l.pos
        for {
                ch := l.next()
                if ch == 0 {
                        return Token{}, errors.New("unterminated string")
                }
                if ch == '"' {
                        break
                }
                if ch == '\\' {
                        l.next()
                }
        }
        strEnd := l.pos - 1
        return Token{Type: TOKEN_STRING, Value: l.input[strStart:strEnd], Start: start, End: l.pos}, nil
}

func isLetter(ch rune) bool {
        return unicode.IsLetter(ch) || ch == '_'
}

// ─── AST ───────────────────────────────────────────────────────────────

type ExprType int

const (
        EXPR_MATCH ExprType = iota
        EXPR_AND
        EXPR_OR
        EXPR_NOT
)

type Expr struct {
        Type   ExprType `json:"type"`
        Left   *Expr    `json:"left,omitempty"`
        Right  *Expr    `json:"right,omitempty"`
        Key    string   `json:"key,omitempty"`
        Op     string   `json:"op,omitempty"`
        Value  string   `json:"value,omitempty"`
        Values []string `json:"values,omitempty"`
}

func (t ExprType) MarshalJSON() ([]byte, error) {
        switch t {
        case EXPR_MATCH:
                return json.Marshal("MATCH")
        case EXPR_AND:
                return json.Marshal("AND")
        case EXPR_OR:
                return json.Marshal("OR")
        case EXPR_NOT:
                return json.Marshal("NOT")
        default:
                return json.Marshal("UNKNOWN")
        }
}

// ─── Parser ─────────────────────────────────────────────────────────────

type Parser struct {
        input   string
        lexer   *Lexer
        curTok  Token
        prevTok Token
        err     error
}

func NewParser(input string) *Parser {
        l := NewLexer(input)
        p := &Parser{lexer: l, input: input}
        p.nextToken()
        return p
}

func (p *Parser) nextToken() {
        p.prevTok = p.curTok
        tok, err := p.lexer.NextToken()
        if err != nil {
                p.errorf("%v", err)
        }
        p.curTok = tok
}

func (p *Parser) errorf(format string, args ...any) {
        if p.err == nil {
                msg := fmt.Sprintf(format, args...)
                p.err = fmt.Errorf("%s\n%s\n%s^", msg, p.input, strings.Repeat(" ", p.curTok.Start))
        }
}

func (p *Parser) noSpaceBetweenTokens() bool {
        return p.prevTok.End == p.curTok.Start
}

func (p *Parser) parseExpression() *Expr {
        return p.parseOr()
}

func (p *Parser) parseOr() *Expr {
        left := p.parseAnd()
        for p.curTok.Type == TOKEN_OR && p.err == nil {
                p.nextToken()
                right := p.parseAnd()
                left = &Expr{Type: EXPR_OR, Left: left, Right: right}
        }
        return left
}

func (p *Parser) parseAnd() *Expr {
        left := p.parseNot()
        for p.curTok.Type == TOKEN_AND && p.err == nil {
                p.nextToken()
                right := p.parseNot()
                left = &Expr{Type: EXPR_AND, Left: left, Right: right}
        }
        return left
}

func (p *Parser) parseNot() *Expr {
        if p.curTok.Type == TOKEN_NOT {
                p.nextToken()
                expr := p.parseAtom()
                return &Expr{Type: EXPR_NOT, Right: expr}
        }
        return p.parseAtom()
}

func (p *Parser) parseAtom() *Expr {
        switch p.curTok.Type {
        case TOKEN_LPAREN:
                p.nextToken()
                expr := p.parseExpression()
                if p.curTok.Type != TOKEN_RPAREN {
                        p.errorf("expected ')' at position %d", p.curTok.Start)
                        return nil
                }
                p.nextToken()
                return expr
        case TOKEN_IDENT:
                key := p.curTok.Value
                p.nextToken()
                if p.curTok.Type != TOKEN_COLON && p.curTok.Type != TOKEN_ASSIGN {
                        p.errorf("expected ':' or ':=' after key %q at position %d", key, p.curTok.Start)
                        return nil
                }
                op := p.curTok.Value
                p.nextToken()
                if !p.noSpaceBetweenTokens() {
                        p.errorf("unexpected space between %q and value", op)
                        return nil
                }
                switch p.curTok.Type {
                case TOKEN_STRING:
                        val := p.curTok.Value
                        p.nextToken()
                        return &Expr{Type: EXPR_MATCH, Key: key, Op: op, Value: val}
                case TOKEN_LBRACKET:
                        p.nextToken()
                        var vals []string
                        for {
                                if p.curTok.Type != TOKEN_STRING {
                                        p.errorf("expected string inside list at position %d", p.curTok.Start)
                                        return nil
                                }
                                vals = append(vals, p.curTok.Value)
                                p.nextToken()
                                if p.curTok.Type == TOKEN_COMMA {
                                        p.nextToken()
                                } else if p.curTok.Type == TOKEN_RBRACKET {
                                        p.nextToken()
                                        break
                                } else {
                                        p.errorf("expected ',' or ']' at position %d", p.curTok.Start)
                                        return nil
                                }
                        }
                        return &Expr{Type: EXPR_MATCH, Key: key, Op: op, Values: vals}
                default:
                        p.errorf("expected string or '[' after operator %q", op)
                        return nil
                }
        default:
                p.errorf("unexpected token %q at position %d", p.curTok.Value, p.curTok.Start)
                return nil
        }
}

func IsValid(input string) (bool, error) {
        _, err := Parse(input)
        if err != nil {
                return false, err
        }
        return true, nil
}

func (p *Parser) Parse() (*Expr, error) {
        expr := p.parseExpression()
        if p.err != nil {
                return nil, p.err
        }
        if p.curTok.Type != TOKEN_EOF {
                p.errorf("unexpected trailing token %q at position %d", p.curTok.Value, p.curTok.Start)
                return nil, p.err
        }
        return expr, nil
}

// ─── API ───────────────────────────────────────────────────────────────

func Parse(input string) (*Expr, error) {
        return NewParser(input).Parse()
}

func MarshalExprToJSON(expr *Expr) (string, error) {
        data, err := json.MarshalIndent(expr, "", "  ")
        if err != nil {
                return "", err
        }
        return string(data), nil
}

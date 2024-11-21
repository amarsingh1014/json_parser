package parser

import (
	"json_parser/internal/tokenizer"
	"fmt"
	"errors"
)

type Parser struct {
	tokens []tokenizer.Token
	current int
}

func NewParser(tokens []tokenizer.Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) currentToken() tokenizer.Token {
	if p.current >= len(p.tokens) {
		return tokenizer.Token{Type: tokenizer.TOKEN_EOF}
	}
	return p.tokens[p.current]
}

func (p *Parser) advance() {
	if p.current < len(p.tokens) - 1 {
		p.current++
	}
}

func (p *Parser) Parse() (interface{}, error) {
	switch p.currentToken().Type {
	case tokenizer.TOKEN_LEFT_BRACE:
		return p.parseObject()
	case tokenizer.TOKEN_LEFT_BRACKET:
		return p.parseArray()
	default:
		return nil, errors.New("invalid JSON")
	}
}

func (p *Parser) parseObject() (map[string]interface{}, error) {
	obj := make(map[string]interface{})

	if p.currentToken().Type != tokenizer.TOKEN_LEFT_BRACE {
		return nil, errors.New("expected {")
	}

	p.advance()

	for {
		if p.currentToken().Type == tokenizer.TOKEN_RIGHT_BRACE{
			p.advance()
			break
		}

		keyToken := p.currentToken()
		if keyToken.Type != tokenizer.TOKEN_STRING {
			return nil, errors.New("expected string")
		}

		key := keyToken.Value
		p.advance()

		if p.currentToken().Type != tokenizer.TOKEN_COLON {
			return nil, errors.New("expected ':' ")
		}
		p.advance()

		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}
		obj[key] = value

		if p.currentToken().Type == tokenizer.TOKEN_COMMA {
			p.advance()
		} else if p.currentToken().Type != tokenizer.TOKEN_RIGHT_BRACE {
			return nil, errors.New("expected ',' or '}'")
		}
	}

	return obj, nil
}

// parseValue handles individual JSON values
func (p *Parser) parseValue() (interface{}, error) {
    switch p.currentToken().Type {
    case tokenizer.TOKEN_STRING:
        value := p.currentToken().Value
        p.advance()
        return value, nil
    case tokenizer.TOKEN_NUMBER:
        value := p.currentToken().Value
        p.advance()
        return value, nil
    case tokenizer.TOKEN_TRUE:
        p.advance()
        return true, nil
    case tokenizer.TOKEN_FALSE:
        p.advance()
        return false, nil
    case tokenizer.TOKEN_NULL:
        p.advance()
        return nil, nil
    case tokenizer.TOKEN_LEFT_BRACE:
        return p.parseObject()
    case tokenizer.TOKEN_LEFT_BRACKET:
        return p.parseArray()
    default:
        return nil, fmt.Errorf("unexpected token: %v", p.currentToken())
    }
}

// parseArray parses a JSON array and returns it as a Go slice
func (p *Parser) parseArray() ([]interface{}, error) {
    array := []interface{}{}

    // Ensure the first token is '['
    if p.currentToken().Type != tokenizer.TOKEN_LEFT_BRACKET {
        return nil, errors.New("expected '[' at the beginning of an array")
    }
    p.advance() // Move past '['

    for {
        // Check for closing bracket ']' which ends the array
        if p.currentToken().Type == tokenizer.TOKEN_RIGHT_BRACKET {
            p.advance()
            break
        }

        // Parse the next value in the array
        value, err := p.parseValue()
        if err != nil {
            return nil, err
        }
        array = append(array, value)

        // Handle comma or closing bracket
        if p.currentToken().Type == tokenizer.TOKEN_COMMA {
            p.advance() // Move past comma and continue loop
        } else if p.currentToken().Type != tokenizer.TOKEN_RIGHT_BRACKET {
            return nil, errors.New("expected ',' or ']' in array")
        }
    }

    return array, nil
}

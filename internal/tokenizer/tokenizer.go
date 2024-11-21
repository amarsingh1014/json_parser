package tokenizer

import (
    "unicode"
    "strings"
    "errors"
)

// TokenType represents the type of JSON token
type TokenType int

const (
    TOKEN_LEFT_BRACE TokenType = iota
    TOKEN_RIGHT_BRACE
    TOKEN_LEFT_BRACKET
    TOKEN_RIGHT_BRACKET
    TOKEN_COLON
    TOKEN_COMMA
    TOKEN_STRING
    TOKEN_NUMBER
    TOKEN_TRUE
    TOKEN_FALSE
    TOKEN_NULL
    TOKEN_EOF
    TOKEN_INVALID
)

// Token represents a JSON token with type and value
type Token struct {
    Type  TokenType
    Value string
}

// Tokenizer holds the input string and current position
type Tokenizer struct {
    input   string
    current int
}

// NewTokenizer creates a new tokenizer for the given input
func NewTokenizer(input string) *Tokenizer {
    return &Tokenizer{input: strings.TrimSpace(input)}
}

// NextToken returns the next token in the input
func (t *Tokenizer) NextToken() (Token, error) {
    for t.current < len(t.input) {
        char := t.input[t.current]

        switch char {
        case '{':
            t.current++
            return Token{Type: TOKEN_LEFT_BRACE, Value: "{"}, nil
        case '}':
            t.current++
            return Token{Type: TOKEN_RIGHT_BRACE, Value: "}"}, nil
        case '[':
            t.current++
            return Token{Type: TOKEN_LEFT_BRACKET, Value: "["}, nil
        case ']':
            t.current++
            return Token{Type: TOKEN_RIGHT_BRACKET, Value: "]"}, nil
        case ':':
            t.current++
            return Token{Type: TOKEN_COLON, Value: ":"}, nil
        case ',':
            t.current++
            return Token{Type: TOKEN_COMMA, Value: ","}, nil
        case '"':
            return t.readString()
        default:
            if unicode.IsDigit(rune(char)) || char == '-' {
                return t.readNumber()
            } else if strings.HasPrefix(t.input[t.current:], "true") {
                t.current += 4
                return Token{Type: TOKEN_TRUE, Value: "true"}, nil
            } else if strings.HasPrefix(t.input[t.current:], "false") {
                t.current += 5
                return Token{Type: TOKEN_FALSE, Value: "false"}, nil
            } else if strings.HasPrefix(t.input[t.current:], "null") {
                t.current += 4
                return Token{Type: TOKEN_NULL, Value: "null"}, nil
            } else if unicode.IsSpace(rune(char)) {
                t.current++
            } else {
                return Token{Type: TOKEN_INVALID, Value: string(char)}, errors.New("invalid character")
            }
        }
    }
    return Token{Type: TOKEN_EOF, Value: ""}, nil
}

// readString reads a JSON string token
func (t *Tokenizer) readString() (Token, error) {
    t.current++ // skip the opening quote
    start := t.current

    for t.current < len(t.input) {
        if t.input[t.current] == '"' {
            value := t.input[start:t.current]
            t.current++ // skip the closing quote
            return Token{Type: TOKEN_STRING, Value: value}, nil
        }
        t.current++
    }
    return Token{Type: TOKEN_INVALID, Value: ""}, errors.New("unterminated string")
}

// readNumber reads a JSON number token
func (t *Tokenizer) readNumber() (Token, error) {
    start := t.current

    for t.current < len(t.input) && (unicode.IsDigit(rune(t.input[t.current])) || t.input[t.current] == '.') {
        t.current++
    }

    value := t.input[start:t.current]
    return Token{Type: TOKEN_NUMBER, Value: value}, nil
}

// Main function to test the tokenizer
// func main() {
//     input := `{"name": "ChatGPT", "age": 3, "isAI": true, "meta": null}`
//     tokenizer := NewTokenizer(input)

//     for {
//         token, err := tokenizer.NextToken()
//         if err != nil {
//             fmt.Println("Error:", err)
//             break
//         }
//         if token.Type == TOKEN_EOF {
//             break
//         }
//         fmt.Printf("Token: %+v\n", token)
//     }
// }

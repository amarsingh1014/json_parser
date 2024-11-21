package main

import (
	"fmt"
	"json_parser/internal/tokenizer"
	"json_parser/internal/parser"
	"encoding/json"
)

func main() {
	input := `{"name": "ChatGPT", "age": 3, "isAI": true, "meta": null}`
	tok := tokenizer.NewTokenizer(input)

	var tokens []tokenizer.Token
	for {
		token, err := tok.NextToken()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		if token.Type == tokenizer.TOKEN_EOF {
			break
		}
		tokens = append(tokens, token)
	}

	p := parser.NewParser(tokens)
	result, err := p.Parse()
	if err != nil {
		fmt.Println("Parse Error:", err)
		return
	}

	// Pretty print the JSON result
	prettyJSON, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		fmt.Println("Error formatting JSON:", err)
		return
	}

	fmt.Println("Parsed JSON:")
	fmt.Println(string(prettyJSON))
}
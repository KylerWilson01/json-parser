package internal_test

import (
	"testing"

	"github.com/KylerWilson01/json-parser/internal"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		name, input string
		expected    []internal.Token
	}{
		{
			"Empty string", "", []internal.Token{
				{Literal: "", Type: internal.EOF},
			},
		},
		{
			"String", `"value"`, []internal.Token{
				{Literal: "value", Type: internal.STRING},
			},
		},
		{
			"Number", `42`, []internal.Token{
				{Literal: "42", Type: internal.NUMBER},
			},
		},
		{
			"Illegal Number", `42f`, []internal.Token{
				{Literal: "42", Type: internal.NUMBER},
				{Literal: "", Type: internal.ILLEGAL},
			},
		},
		{
			"Simple object", "{}", []internal.Token{
				{Literal: "{", Type: internal.OPENING_CURLY},
				{Literal: "}", Type: internal.CLOSING_CURLY},
				{Literal: "", Type: internal.EOF},
			},
		},
		{
			"Simple key value object", "{\"key\": \"value\"}", []internal.Token{
				{Literal: "{", Type: internal.OPENING_CURLY},
				{Literal: "key", Type: internal.STRING},
				{Literal: ":", Type: internal.COLON},
				{Literal: "value", Type: internal.STRING},
				{Literal: "}", Type: internal.CLOSING_CURLY},
				{Literal: "", Type: internal.EOF},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := internal.NewLexer(tt.input)
			for _, expected := range tt.expected {
				actual := l.NextToken()
				if actual != expected {
					t.Errorf("expected %v, got %v", expected, actual)
				}
			}
		})
	}
}

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
			"Empty string", "", []internal.Token{},
		},
		{
			"String", `"value"`, []internal.Token{
				{Literal: "value", Type: internal.String},
			},
		},
		{
			"Number", `42`, []internal.Token{
				{Literal: "42", Type: internal.Number},
			},
		},
		{
			"Simple object", "{}", []internal.Token{
				{Literal: "{", Type: internal.OpeningCurly},
				{Literal: "}", Type: internal.ClosingCurly},
			},
		},
		{
			"Simple key value object", "{\"key\": \"value\"}", []internal.Token{
				{Literal: "{", Type: internal.OpeningCurly},
				{Literal: "key", Type: internal.String},
				{Literal: ":", Type: internal.Colon},
				{Literal: "value", Type: internal.String},
				{Literal: "}", Type: internal.ClosingCurly},
			},
		},
		{
			"Nested Object and array", `{"key": "value","key-n": 101,"key-o": {},"key-l": []}`,
			[]internal.Token{
				{Literal: "{", Type: internal.OpeningCurly},
				{Literal: "key", Type: internal.String},
				{Literal: ":", Type: internal.Colon},
				{Literal: "value", Type: internal.String},
				{Literal: ",", Type: internal.Comma},
				{Literal: "key-n", Type: internal.String},
				{Literal: ":", Type: internal.Colon},
				{Literal: "101", Type: internal.Number},
				{Literal: ",", Type: internal.Comma},
				{Literal: "key-o", Type: internal.String},
				{Literal: ":", Type: internal.Colon},
				{Literal: "{", Type: internal.OpeningCurly},
				{Literal: "}", Type: internal.ClosingCurly},
				{Literal: ",", Type: internal.Comma},
				{Literal: "key-l", Type: internal.String},
				{Literal: ":", Type: internal.Colon},
				{Literal: "[", Type: internal.OpeningBracket},
				{Literal: "]", Type: internal.CloseingBracket},
				{Literal: "}", Type: internal.ClosingCurly},
			},
		},
		{
			"Multiline Object", `{\n"key1": true,\n"key2": false,\n"key3": null,\n"key4": "value",\n"key5": 101\n}`,
			[]internal.Token{
				{Literal: "{", Type: internal.OpeningCurly},
				{Literal: "key1", Type: internal.String},
				{Literal: ":", Type: internal.Colon},
				{Literal: "true", Type: internal.True},
				{Literal: ",", Type: internal.Comma},
				{Literal: "key2", Type: internal.String},
				{Literal: ":", Type: internal.Colon},
				{Literal: "false", Type: internal.False},
				{Literal: ",", Type: internal.Comma},
				{Literal: "key3", Type: internal.String},
				{Literal: ":", Type: internal.Colon},
				{Literal: "null", Type: internal.Null},
				{Literal: ",", Type: internal.Comma},
				{Literal: "key4", Type: internal.String},
				{Literal: ":", Type: internal.Colon},
				{Literal: "value", Type: internal.String},
				{Literal: ",", Type: internal.Comma},
				{Literal: "key5", Type: internal.String},
				{Literal: ":", Type: internal.Colon},
				{Literal: "101", Type: internal.Number},
				{Literal: "}", Type: internal.ClosingCurly},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := internal.NewLexer(tt.input)
			err := l.ValidateTokens()
			if err != nil {
				t.Fail()
			}

			for i, expected := range tt.expected {
				actual := l.Tokens[i]
				if actual.Type != expected.Type || actual.Literal != expected.Literal {
					t.Errorf("expected %v, got %v", expected, actual)
				}
			}
		})
	}
}

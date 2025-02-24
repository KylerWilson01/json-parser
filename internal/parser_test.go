package internal_test

import (
	"reflect"
	"testing"

	"github.com/KylerWilson01/json-parser/internal"
)

func TestParseTokens_Valid(t *testing.T) {
	testCases := []struct {
		name         string
		tokenList    []internal.Token
		expectedBool bool
	}{
		{
			name: "Valid string json",
			tokenList: func() []internal.Token {
				l := internal.NewLexer(`{"name":"ayo"}`)
				err := l.ValidateTokens()
				if err != nil {
					return []internal.Token{}
				}
				return l.Tokens
			}(),
			expectedBool: true,
		},
		{
			name: "Valid num json",
			tokenList: func() []internal.Token {
				l := internal.NewLexer(`{
						  "num1": 123,
						  "num2": 987
						}`)
				err := l.ValidateTokens()
				if err != nil {
					return []internal.Token{}
				}
				return l.Tokens
			}(),
			expectedBool: true,
		},
		{
			name: "Valid array json with string at idx 0",
			tokenList: func() []internal.Token {
				l := internal.NewLexer(`{
						  "array": ["testing", "hello", 123, {"name":"hello"}]
						}`)
				err := l.ValidateTokens()
				if err != nil {
					return []internal.Token{}
				}
				return l.Tokens
			}(),
			expectedBool: true,
		},
		{
			name: "simple empty object",
			tokenList: func() []internal.Token {
				l := internal.NewLexer(`{}`)
				err := l.ValidateTokens()
				if err != nil {
					return []internal.Token{}
				}
				return l.Tokens
			}(),
			expectedBool: true,
		},
		{
			name: "big complicated array",
			tokenList: func() []internal.Token {
				l := internal.NewLexer(`
				[
					"JSON Test Pattern pass1",
					{"object with 1 member":["array with 1 element"]},
					{},
					[],
					-42,
					true,
					false,
					null,
					{
						"integer": 1234567890,
						"real": -9876.543210,
						"e": 0.123456789e-12,
						"E": 1.234567890E+34,
						"":  23456789012E66,
						"zero": 0,
						"one": 1,
						"space": " ",
						"quote": "\"",
						"backslash": "\\",
						"controls": "\b\f\n\r\t",
						"slash": "/ & \/",
						"alpha": "abcdefghijklmnopqrstuvwyz",
						"ALPHA": "ABCDEFGHIJKLMNOPQRSTUVWYZ",
						"digit": "0123456789",
						"0123456789": "digit",
						"special": "1~!@#$%^&*()_+-={':[,]}|;.</>?",
									"hex": "\u0123\u4567\u89AB\uCDEF\uabcd\uef4A",
										"true": true,
										"false": false,
										"null": null,
										"array":[  ],
							"object":{  },
							"address": "50 St. James Street",
							"url": "http://www.JSON.org/",
							"comment": "// /* <!-- --",
							"# -- --> */": " ",
							" s p a c e d " :[1,2 , 3
				
							,
				
							4 , 5        ,          6           ,7        ],"compact":[1,2,3,4,5,6,7],
							"jsontext": "{\"object with 1 member\":[\"array with 1 element\"]}",
							"quotes": "&#34; \u0022 %22 0x22 034 &#x22;",
							"\/\\\"\uCAFE\uBABE\uAB98\uFCDE\ubcda\uef4A\b\f\n\r\t1~!@#$%^&*()_+-=[]{}|;:',./<>?"
							: "A key can be any string"
						},
						0.5 ,98.6
						,
						99.44
						,
				
						1066,
						1e1,
						0.1e1,
						1e-1,
						1e00,2e+00,2e-00
						,"rosebud"]
`)
				err := l.ValidateTokens()
				if err != nil {
					return []internal.Token{}
				}
				return l.Tokens
			}(),
			expectedBool: true,
		},
		{
			name: "deeply nested array",
			tokenList: func() []internal.Token {
				l := internal.NewLexer(`[[[[[[[[[[[[[[[[[[["Not too deep"]]]]]]]]]]]]]]]]]]]`)
				err := l.ValidateTokens()
				if err != nil {
					return []internal.Token{}
				}
				return l.Tokens
			}(),
			expectedBool: true,
		},
		{
			name: "nested json",
			tokenList: func() []internal.Token {
				l := internal.NewLexer(`
					{
						"JSON Test Pattern pass3": {
							"The outermost value": "must be an object or array.",
							"In this test": "It is an object."
						}
					}`)
				err := l.ValidateTokens()
				if err != nil {
					return []internal.Token{}
				}
				return l.Tokens
			}(),
			expectedBool: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := internal.NewParser(tc.tokenList)
			actualBool, err := p.ParseTokens()
			if actualBool != tc.expectedBool {
				t.Errorf("expected: %v, got: %v. details: %v", tc.expectedBool, actualBool, err)
			}
		})
	}
}

func TestParseTokens_Invalid(t *testing.T) {
	testCases := []struct {
		name         string
		tokenList    []internal.Token
		expectedBool bool
		errStr       string
	}{
		{
			name: "Unclosed Array",
			tokenList: func() []internal.Token {
				l := internal.NewLexer(`["Unclosed array"`)
				err := l.ValidateTokens()
				if err != nil {
					return []internal.Token{}
				}
				return l.Tokens
			}(),
			expectedBool: false,
		},
		{
			name: "Extra comma",
			tokenList: func() []internal.Token {
				l := internal.NewLexer(`["extra comma",]`)
				err := l.ValidateTokens()
				if err != nil {
					return []internal.Token{}
				}
				return l.Tokens
			}(),
			expectedBool: false,
		},
		{
			name: "Double comma",
			tokenList: func() []internal.Token {
				l := internal.NewLexer(`["double comma",,]`)
				err := l.ValidateTokens()
				if err != nil {
					return []internal.Token{}
				}
				return l.Tokens
			}(),
			expectedBool: false,
		},
		{
			name: "Missing array value",
			tokenList: func() []internal.Token {
				l := internal.NewLexer(`[   , "<-- missing value"]`)
				err := l.ValidateTokens()
				if err != nil {
					return []internal.Token{}
				}
				return l.Tokens
			}(),
			expectedBool: false,
		},
		{
			name: "comma after close",
			tokenList: func() []internal.Token {
				l := internal.NewLexer(`["Comma after the close"],`)
				err := l.ValidateTokens()
				if err != nil {
					return []internal.Token{}
				}
				return l.Tokens
			}(),
			expectedBool: false,
		},
		{
			name: "Extra comma object",
			tokenList: func() []internal.Token {
				l := internal.NewLexer(`{"Extra comma": true,}`)
				err := l.ValidateTokens()
				if err != nil {
					return []internal.Token{}
				}
				return l.Tokens
			}(),
			expectedBool: false,
		},
		{
			name: "Missing Colon",
			tokenList: func() []internal.Token {
				l := internal.NewLexer(`{"Missing colon" null}`)
				err := l.ValidateTokens()
				if err != nil {
					return []internal.Token{}
				}
				return l.Tokens
			}(),
			expectedBool: false,
		},
		{
			name: "double colon",
			tokenList: func() []internal.Token {
				l := internal.NewLexer(`{"Double colon":: null}`)
				err := l.ValidateTokens()
				if err != nil {
					return []internal.Token{}
				}
				return l.Tokens
			}(),
			expectedBool: false,
		},
		{
			name: "comma instead of colon",
			tokenList: func() []internal.Token {
				l := internal.NewLexer(`{"Comma instead of colon", null}`)
				err := l.ValidateTokens()
				if err != nil {
					return []internal.Token{}
				}
				return l.Tokens
			}(),
			expectedBool: false,
		},
		{
			name: "colon instead of comma",
			tokenList: func() []internal.Token {
				l := internal.NewLexer(`["Colon instead of comma": false]`)
				err := l.ValidateTokens()
				if err != nil {
					return []internal.Token{}
				}
				return l.Tokens
			}(),
			expectedBool: false,
		},
		{
			name: "comma instead of closing brace",
			tokenList: func() []internal.Token {
				l := internal.NewLexer(`{"Comma instead if closing brace": true,`)
				err := l.ValidateTokens()
				if err != nil {
					return []internal.Token{}
				}
				return l.Tokens
			}(),
			expectedBool: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := internal.NewParser(tc.tokenList)
			actualBool, err := p.ParseTokens()
			if !reflect.DeepEqual(actualBool, tc.expectedBool) {
				t.Errorf("expected: %v, got: %v, details: %v", tc.expectedBool, actualBool, err)
			}
		})
	}
}

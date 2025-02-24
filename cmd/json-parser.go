package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/KylerWilson01/json-parser/internal"
)

func main() {
	flag.Parse()

	fp := flag.Args()
	if len(fp) == 0 {
		fp = append(fp, "-")
	}

	for _, file := range fp {
		var f *os.File
		var err error

		if file == "-" {
			f = os.Stdin
		} else {
			f, err = os.Open(file)
			if err != nil {
				os.Exit(1)
				return
			}
			defer f.Close()
		}

		data := ""
		scanner := bufio.NewReader(f)
		for {
			line, err := scanner.ReadString('\n')
			data += line
			if err != nil {
				if err == io.EOF {
					break
				}
				fmt.Println("Something went wrong with reading the file")
				os.Exit(0)
			}
		}

		lexer := internal.NewLexer(data)
		lexer.ValidateTokens()

		parser := internal.NewParser(lexer.Tokens)
		v, err := parser.ParseTokens()
		if err != nil {
			fmt.Printf("Not valid json, details: %v\n", err)
			os.Exit(1)
			return
		}

		if v != true {
			fmt.Println("Not valid json")
			os.Exit(1)
			return
		}

		fmt.Println("Valid json")
		os.Exit(0)
	}
}

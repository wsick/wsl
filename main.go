package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/wsick/wsl/lex"
)

func main() {
	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	if fi.Mode()&os.ModeNamedPipe == 0 {
		fmt.Println("no input")
	} else {
		raw, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}

		l := lex.FromString(fi.Name(), string(raw))
		for {
			token := l.NextToken()
			if token.Type() != lex.TokenWS {
				fmt.Printf("%d:%d:%s\n", token.Line(), token.Col()+1, token)
			}
			if token.Type() == lex.TokenEOF || token.Type() == lex.TokenError {
				break
			}
		}
	}
}

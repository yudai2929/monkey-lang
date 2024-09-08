package repl

import (
	"bufio"
	"fmt"
	"gtihub.com/yudai2929/monkey-lang/lexer"
	"gtihub.com/yudai2929/monkey-lang/token"
	"io"
)

const PROMPT = ">> "

func Start(in io.Reader, _ io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)

		for tok := l.NextToken(); !tok.IsType(token.EOF); tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}

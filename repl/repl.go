package repl

import (
	"bufio"
	"fmt"
	"gtihub.com/yudai2929/monkey-lang/lexer"
	"gtihub.com/yudai2929/monkey-lang/parser"
	"io"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) error {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return nil
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			if err := printParserErrors(out, p.Errors()); err != nil {
				return err
			}
			continue
		}

		if _, err := io.WriteString(out, program.String()); err != nil {
			return err
		}
		if _, err := io.WriteString(out, "\n"); err != nil {
			return err
		}

	}
}

func printParserErrors(out io.Writer, errors []string) error {
	for _, msg := range errors {
		if _, err := io.WriteString(out, "\t"+msg+"\n"); err != nil {
			return err
		}
	}

	return nil
}

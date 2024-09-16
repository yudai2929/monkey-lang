package repl

import (
	"bufio"
	"fmt"
	"gtihub.com/yudai2929/monkey-lang/evalutor"
	"gtihub.com/yudai2929/monkey-lang/lexer"
	"gtihub.com/yudai2929/monkey-lang/object"
	"gtihub.com/yudai2929/monkey-lang/parser"
	"io"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) error {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

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

		evaluated := evalutor.Eval(program, env)
		if evaluated != nil {
			if _, err := io.WriteString(out, evaluated.Inspect()); err != nil {
				return err
			}
			if _, err := io.WriteString(out, "\n"); err != nil {
				return err
			}
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

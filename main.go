package main

import (
	"fmt"
	"gtihub.com/yudai2929/monkey-lang/repl"
	"os"
	"os/user"
)

func main() {

	u, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the Monkey programming language!\n", u.Username)
	fmt.Printf("Feel free to type in commands\n")
	err = repl.Start(os.Stdin, os.Stdout)
	if err != nil {
		panic(err)
	}
}

package main

import (
	"flag"
	"fmt"
	"github.com/alessio/shellescape"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	var env map[string]string
	flag.Parse()

	env, err := godotenv.Read()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	args := flag.Args()

	if len(args) == 1 {
		fmt.Printf("%s", shellescape.Quote(arg))
		fmt.Printf("%s", arg)
	} else if len(args) > 1 {
		for _, k := range args {
			fmt.Printf("%s=%s\n", k, shellescape.Quote(env[k]))
		}
	} else {
		for k, v := range env {
			fmt.Printf("%s=%s\n", k, shellescape.Quote(v))
		}
	}
}

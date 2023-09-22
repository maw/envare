package main

import (
	"flag"
	"fmt"
	"github.com/alessio/shellescape"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func read(path string) map[string]string {
	var env map[string]string

	if path == "" {
		env, err := godotenv.Read()
		if err != nil {
			log.Fatal("Error loading .env file")
		}

		return env
	}

	realPath := path
	env, err := godotenv.Read(realPath)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error loading %s", realPath))
	}

	return env
}

func main() {
	var env map[string]string

	var path string
	var captureMode bool

	flag.StringVar(&path, "f", "", "path of alternate .env file")
	flag.BoolVar(&captureMode, "c", false, "run in capture-for-eval mode")
	flag.Parse()

	env = read(path)

	args := flag.Args()

	if (len(args)) > 1 {
		captureMode = true
	}
	if len(args) == 0 {
		for k, v := range env {
			fmt.Printf("%s=%s\n", k, shellescape.Quote(v))
		}
		os.Exit(0)
	}

	if len(args) == 1 && captureMode == false {
		fmt.Printf("%s", shellescape.Quote(env[args[0]]))
		os.Exit(0)
	}

	for _, k := range args {
		fmt.Printf("%s=%s\n", k, shellescape.Quote(env[k]))
	}
}

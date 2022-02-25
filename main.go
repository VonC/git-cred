package main

import (
	"fmt"
	"log"
	"minimalcli/internal/syscall"
	"minimalcli/version"
	"os"
	"path/filepath"

	"github.com/jpillora/opts"
	"github.com/ryboe/q"
	"github.com/spewerspew/spew"
)

// Config stores arguments and subcommands
type Config struct {
	Arg     string `help:"a string argument"`
	Version bool   `help:"if true, print Version and exit."`
}

var c = &Config{}

func fatal(msg string, err error) {
	if err != nil {
		log.Fatalf("%s: error '%+v'", msg, err)
	}
}

// myproject main entry
func main() {

	var err error

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	fatal("Unable to find current program execution directory", err)
	log.Println(dir)
	opts.New(c).
		ConfigPath(filepath.Join(dir, "conf.json")).
		Parse()
	if c.Version {
		fmt.Println(version.String())
		os.Exit(0)
	}

	spew.Dump(c)
	q.Q(c)

	stderr, stdout, err := syscall.ExecCmd("hostname")
	fatal("Unable to call hostname", err)

	fmt.Printf("Stdout: '%s'\n", stdout.String())
	fmt.Printf("Stderr: '%s'\n", stderr.String())

	fmt.Println(os.Args[0])
}

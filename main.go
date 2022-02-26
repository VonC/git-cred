package main

import (
	"fmt"
	"gitcred/internal/credhelper"
	"gitcred/version"
	"log"
	"os"
	"path/filepath"

	"github.com/jpillora/opts"
	"github.com/ryboe/q"
	"github.com/spewerspew/spew"
)

// Config stores arguments and subcommands
type Config struct {
	Host    string `opts:"help=repository hosting service name, mode=arg"`
	Version bool   `help:"if true, print Version and exit."`
	Debug   bool   `help:"if true, print Debug information."`
}

var c = &Config{}

func fatal(msg string, err error) {
	if err != nil {
		log.Fatalf("%s: error '%+v'", msg, err)
	}
}

type CredHelper interface {
	Get() (string, error)
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

	if c.Debug {
		spew.Dump(c)
		q.Q(c)
	}

	fmt.Println(os.Args[0])

	var ch CredHelper
	ch, err = credhelper.NewCredHelper(c.Host)
	fatal("Unable to get Credential Helper", err)
	get, err := ch.Get()
	fatal("Get error", err)
	fmt.Println(get)
}

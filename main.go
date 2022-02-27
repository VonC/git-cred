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
	Host     string `opts:"help=repository hosting service name, mode=arg"`
	Debug    bool   `help:"if true, print Debug information."`
	Username string `help:"Get: username, by default current user if not set"`
	*Set     `opts:"mode=cmd,help=set user password for a given host"`
	*Erase   `opts:"mode=cmd,help=erase password for a given host and username"`
	ch       CredHelper
}

type Set struct {
	User     string `opts:"help=username, mode=arg"`
	Password string `opts:"help=user password, mode=arg"`
}

func (s *Set) isEmpty() bool {
	return !(s != nil && s.User != "" && s.Password != "")
}

type Erase struct {
	User string `opts:"help=username, mode=arg"`
}

func (e *Erase) isEmpty() bool {
	return !(e != nil && e.User != "")
}

var c = &Config{}

func fatal(msg string, err error) {
	if err != nil {
		log.Fatalf("%s: error '%+v'", msg, err)
	}
}

type CredHelper interface {
	Get(username string) (string, error)
	Set(username, password string) error
	Erase(username string) error
}

// myproject main entry
func main() {

	var err error

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	fatal("Unable to find current program execution directory", err)
	//log.Println(dir)
	opts.New(c).
		ConfigPath(filepath.Join(dir, "conf.json")).
		Version(version.String()).
		Parse()

	if c.Debug {
		spew.Dump(c)
		q.Q(c)
	}

	//fmt.Println(os.Args[0])

	var ch CredHelper
	ch, err = credhelper.NewCredHelper(c.Host)
	fatal("Unable to get Credential Helper", err)
	c.ch = ch
	err = c.Run()
	fatal("gitcred ERROR", err)
}

func (c *Config) Run() error {
	if !c.Set.isEmpty() {
		return c.Set.Run()
	}
	if !c.Erase.isEmpty() {
		return c.Erase.Run()
	}
	get, err := c.ch.Get(c.Username)
	if err != nil {
		return err
	}
	fmt.Println(get)
	return nil
}

func (s *Set) Run() error {
	err := c.ch.Set(s.User, s.Password)
	return err
}

func (e *Erase) Run() error {
	err := c.ch.Erase(e.User)
	return err
}

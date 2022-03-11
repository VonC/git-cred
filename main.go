package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/VonC/gitcred/version"

	"github.com/VonC/gitcred/internal/credhelper"

	"github.com/jpillora/opts"
	"github.com/ryboe/q"
	"github.com/spewerspew/spew"
)

// Config stores arguments and subcommands
type Config struct {
	Servername string `help:"help=repository hosting server name (hostname). If not set, use the one from current repository folder, if present in pwd"`
	Debug      bool   `help:"if true, print Debug information."`
	Username   string `help:"Get: username. If not set, use the one from from current repository remote URL, if present in pwd"`
	*Set       `opts:"mode=cmd,help=[password] set user password for a given host: -u/--username mandatory"`
	*Erase     `opts:"mode=cmd,help=erase password for a given host and username: -u/--username and -s/--servername mandatory"`
	ch         CredHelper
}

type Set struct {
	Password string `opts:"help=user password, mode=arg"`
}

func (s *Set) isEmpty() bool {
	return !(s != nil && s.Password != "")
}

type Erase struct {
}

func (e *Erase) isEmpty() bool {
	return !(e != nil)
}

var c = &Config{}

func fatal(msg string, err error) {
	if err != nil {
		log.Fatalf("%s: error '%+v'", msg, err)
	}
}

type CredHelper interface {
	Get(username string) (string, error)
	Set(username, password, host string) error
	Erase(username, host string) error
	Host() string
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
	ch, err = credhelper.NewCredHelper(c.Servername)
	fatal("Unable to get Credential Helper", err)
	c.ch = ch
	if c.Servername == "" {
		c.Servername = c.ch.Host()
	}
	err = c.Run()
	fatal("gitcred ERROR", err)
}

func (c *Config) Run() error {
	if !c.Set.isEmpty() && c.Username != "" {
		return c.Set.Run()
	}
	if !c.Erase.isEmpty() && c.Username != "" && c.Servername != "" {
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
	err := c.ch.Set(c.Username, s.Password, c.ch.Host())
	return err
}

func (e *Erase) Run() error {
	err := c.ch.Erase(c.Username, c.ch.Host())
	return err
}

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/VonC/gitcred/version"
	"github.com/ryboe/q"
	"github.com/spewerspew/spew"

	"github.com/VonC/gitcred/internal/credhelper"

	"github.com/jpillora/opts"
)

// Config stores arguments and subcommands
type Config struct {
	Servername string `help:"help=repository hosting server name (hostname). If not set, use the one from current repository folder, if present in pwd"`
	Debug      bool   `help:"if true, print Debug information."`
	Username   string `help:"Get: username. If not set, use the one from from current repository remote URL, if present in pwd"`
	*Get       `opts:"mode=cmd,help=get password for a given host and username: can read those from current folder repository"`
	*Set       `opts:"mode=cmd,help=[password] set user password for a given host: -u/--username mandatory"`
	*Erase     `opts:"mode=cmd,help=erase password for a given host and username: -u/--username and -s/--servername mandatory"`
	ch         CredHelper
}

type Set struct {
	Password string `opts:"help=user password, mode=arg"`
}

type Erase struct {
}

type Get struct {
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

	//fmt.Println(os.Args[0])

	var ch CredHelper
	ch, err = credhelper.NewCredHelper(c.Servername)
	fatal("Unable to get Credential Helper", err)
	c.ch = ch
	if c.Servername == "" {
		c.Servername = ch.Host()
	}
	if !hasCmd("erase") {
		c.Erase = nil
	}
	if !hasCmd("set") {
		c.Set = nil
	}

	if c.Debug {
		spew.Dump(c)
		q.Q(c)
	}

	err = c.Run()
	fatal("gitcred ERROR", err)
}

func (c *Config) Run() error {
	if c.Erase != nil {
		return c.Erase.Run()
	}
	if c.Set != nil {
		return c.Set.Run()
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

func hasCmd(cmd string) bool {
	previousIsOption := false
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-") {
			previousIsOption = true
			continue
		}
		if cmd == arg && !previousIsOption {
			return true
		}
		previousIsOption = false
	}
	return false
}

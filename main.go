package main

import (
	"fmt"
	"log"

	"github.com/ryboe/q"
	"github.com/spewerspew/spew"

	"github.com/VonC/gitcred/internal/credhelper"

	"github.com/alecthomas/kong"
)

// CLI stores arguments and subcommands
type CLI struct {
	Servername string   `help:"repository hosting server name (hostname). If not set, use the one from current repository folder, if present in pwd" short:"s" type:"string"`
	Debug      bool     `help:"if true, print Debug information." type:"bool" short:"d"`
	Username   string   `help:"Get: username. If not set, use the one from from current repository remote URL, if present in pwd" short:"u"`
	Get        GetCmd   `cmd:"" help:"get password for a given host and username: can read those from current folder repository" name:"get"`
	Set        SetCmd   `cmd:"" help:"[password] set user password for a given host: -u/--username mandatory" name:"set" aliases:"store"`
	Erase      EraseCmd `cmd:"" help:"erase password for a given host and username: -u/--username and -s/--servername mandatory" name:"erase" aliases:"rm,del,delete,remove"`
	ch         CredHelper
}

type SetCmd struct {
	Password string `arg:"" help:"user password or token" short:"p"`
}

type EraseCmd struct{}

type GetCmd struct{}

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
	User() string
}

type Context struct {
	*CLI
}

// myproject main entry
func main() {

	var err error

	//dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	//fatal("Unable to find current program execution directory", err)
	//log.Println(dir)
	var cli CLI
	ctx := kong.Parse(&cli)
	//ctx.BindTo(os.Stdout, (*io.Writer)(nil))

	//fmt.Println(os.Args[0])

	var ch CredHelper
	ch, err = credhelper.NewCredHelper(cli.Servername, cli.Username)
	fatal("Unable to get Credential Helper", err)
	cli.ch = ch
	if cli.Servername == "" {
		cli.Servername = ch.Host()
	}
	if cli.Username == "" {
		cli.Username = ch.User()
	}

	spew.Dump(cli)

	if cli.Debug {
		spew.Dump(cli)
		q.Q(cli)
	}

	fmt.Printf("ctx command '%s'\n", ctx.Command())

	err = ctx.Run(&Context{CLI: &cli})
	fatal("gitcred Unable to run:", err)
}

func (s *GetCmd) Run(c *Context) error {
	get, err := c.ch.Get(c.Username)
	if err != nil {
		return err
	}
	fmt.Println(get)
	return nil
}

func (s *SetCmd) Run(c *Context) error {
	err := c.ch.Set(c.Username, s.Password, c.ch.Host())
	return err
}

func (e *EraseCmd) Run(c *Context) error {
	err := c.ch.Erase(c.Username, c.ch.Host())
	return err
}

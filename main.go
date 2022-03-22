package main

import (
	"embed"
	"fmt"
	"log"

	"github.com/ryboe/q"
	"github.com/spewerspew/spew"

	"github.com/VonC/gitcred/internal/credhelper"
	"github.com/VonC/gitcred/version"

	"github.com/alecthomas/kong"
)

// CLI stores arguments and subcommands
type CLI struct {
	Servername string   `help:"repository hosting server name (hostname). If not set, use the one from current repository folder, if present in pwd" short:"s" type:"string"`
	Debug      bool     `help:"if true, print Debug information." type:"bool" short:"d" env:"DEBUG"`
	Username   string   `help:"Get: username. If not set, use the one from from current repository remote URL, if present in pwd" short:"u"`
	Get        GetCmd   `cmd:"" help:"get password for a given host and username: can read those from current folder repository" name:"get" default:""`
	Set        SetCmd   `cmd:"" help:"[password] set user password for a given host: -u/--username mandatory" name:"set" aliases:"store"`
	Erase      EraseCmd `cmd:"" help:"erase password for a given host and username: -u/--username and -s/--servername mandatory" name:"erase" aliases:"rm,del,delete,remove"`
	ch         CredHelper
	Version    VersionFlag `name:"version" help:"Print version information and quit" short:"v" type:"counter"`
	VersionC   VersionCmd  `cmd:"" help:"Show the version information" name:"version" aliases:"ver"`
	versionFs  embed.FS
}

type VersionFlag int

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

type VersionCmd struct{}

type CredHelper interface {
	Get(username string, servername string) (string, error)
	Set(username, password, host string) error
	Erase(username, host string) error
}

type Context struct {
	*CLI
}

// https://github.com/golang/go/issues/41191
// https://stackoverflow.com/a/67357103/6309
//go:embed version/*
var versionFs embed.FS

// myproject main entry
func main() {

	var err error

	var cli CLI
	ctx := kong.Parse(&cli)

	if cli.Debug {
		spew.Dump(cli)
		q.Q(cli)
		fmt.Printf("ctx command '%s'\n", ctx.Command())
	}

	cli.versionFs = versionFs

	if ctx.Command() != "version" && cli.Version > 0 {
		fmt.Printf(version.String(int(cli.Version), cli.versionFs))
		ctx.Exit(0)
	}

	if ctx.Command() != "version" {
		var ch CredHelper
		ch, err = credhelper.NewCredHelper(cli.Servername, cli.Username)
		fatal("Unable to get Credential Helper", err)
		cli.ch = ch
	}

	err = ctx.Run(&Context{CLI: &cli})
	fatal("gitcred Unable to run:", err)
}

func (s *GetCmd) Run(c *Context) error {
	get, err := c.ch.Get(c.Username, c.Servername)
	if err != nil {
		return err
	}
	fmt.Println(get)
	return nil
}

func (s *SetCmd) Run(c *Context) error {
	err := c.ch.Set(c.Username, s.Password, c.Servername)
	return err
}

func (e *EraseCmd) Run(c *Context) error {
	err := c.ch.Erase(c.Username, c.Servername)
	return err
}

func (v *VersionCmd) Run(c *Context) error {
	//spew.Dump(c)
	fmt.Printf(version.String(int(c.Version)+1, c.versionFs))
	return nil
}

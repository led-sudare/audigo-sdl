package main

import (
	"os"

	"github.com/code560/audigo-sdl/app"
	"github.com/code560/audigo-sdl/util"
	"github.com/urfave/cli"
)

var log = util.GetLogger()

func main() {
	// reset cd
	execDir, _ := os.Executable()
	log.Debugf("executable: %s", execDir)

	cl := cli.NewApp()
	cl.Name = "audigo"
	cl.Usage = "Audio service by LED CUBU"
	cl.Version = "1.0.1"

	cl.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "cd, c",
			Usage: "change current directory by repl",
			Value: "",
		},
	}
	clientFlags := append(cl.Flags,
		cli.StringFlag{
			Name:  "domain, d",
			Usage: "set request domain url by client",
			Value: "http://audigo.local",
		})
	cl.Commands = []cli.Command{
		{
			Name:    "server",
			Aliases: []string{"s"},
			Usage:   "Instant server mode. (default)",
			Action:  doServe,
			Flags:   cl.Flags,
		},
		{
			Name:    "client",
			Aliases: []string{"c"},
			Usage:   "Instant client REPL mode.",
			Action:  doClientRepl,
			Flags:   clientFlags,
		},
		{
			Name:    "repl",
			Aliases: []string{"r"},
			Usage:   "Instant local REPL mode.",
			Action:  doRepl,
			Flags:   cl.Flags,
		},
	}
	cl.Action = doServe

	cl.Run(os.Args)
}

func doServe(ctx *cli.Context) error {
	asCd(ctx)
	app.Serve(ctx.Args().Get(0))
	return nil
}

func doClientRepl(ctx *cli.Context) error {
	asCd(ctx)
	url := ctx.String("d")
	app.ClientRepl(url, "abc")
	return nil
}

func doRepl(ctx *cli.Context) error {
	asCd(ctx)
	app.Repl()
	return nil
}

func asCd(ctx *cli.Context) {
	if ctx.IsSet("cd") {
		dir := ctx.String("cd")
		cd(dir)
	}
}

func cd(dir string) {
	if dir != "" {
		stat, _ := os.Stat(dir)
		if stat.IsDir() {
			os.Chdir(dir)
			log.Debugf("change current directory: %s", dir)
		}
	}
}

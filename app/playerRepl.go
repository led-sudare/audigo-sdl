package app

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/code560/audigo/player"
	"github.com/code560/audigo/util"
)

var log = util.GetLogger()

type playerRepl struct {
	player player.Proxy
}

func Repl() {
	// init
	c := &playerRepl{
		player: player.NewProxy(),
	}
	c.printInit()

	// pipeline
	if terminal.IsTerminal(0) {
		src, _ := ioutil.ReadAll(os.Stdin)
		c.play(string(src))
		return
	}

	// wait input
	c.cli()
}

func (c *playerRepl) sendChan(a *player.Action) bool {
	select {
	case c.player.GetChannel() <- a:
		return true
	default:
		log.Errorf("dont send chan: %s", a.Act)
		return false
	}
}

func (c *playerRepl) play(src string) {
	c.sendChan(
		&player.Action{
			Act:  player.Play,
			Args: &player.PlayArgs{Src: src},
		})
}

func (c *playerRepl) stop() {
	c.sendChan(&player.Action{Act: player.Stop})
}

func (c *playerRepl) pause() {
	c.sendChan(&player.Action{Act: player.Pause})
}

func (c *playerRepl) resume() {
	c.sendChan(&player.Action{Act: player.Resume})
}

func (c *playerRepl) volume(a string) {
	v, err := strconv.ParseFloat(a, 64)
	if err != nil {
		log.Warn(fmt.Sprintf("invalid format volume string: %s", a))
		return
	}
	c.sendChan(&player.Action{
		Act:  player.Volume,
		Args: &player.VolumeArgs{Vol: v},
	})
}

func (c *playerRepl) cli() {
	s := bufio.NewScanner(os.Stdin)

	for {
		c.printHead()
		s.Scan() // <- wait input
		inputs := strings.Split(s.Text(), " ")
		switch strings.ToLower(inputs[0]) {
		case "play":
			c.play(inputs[1])
		case "stop":
			c.stop()
		case "pause":
			c.pause()
		case "resume":
			c.resume()
		case "volume":
			c.volume(inputs[1])
		case "exit":
			return
		case "help":
			c.help()
		default:
			fmt.Println()
		}
	}
}

func (c *playerRepl) printInit() {
	fmt.Println("welcome to audigo REPL 1.0")
}

func (c *playerRepl) help() {
	fmt.Print(`
        play <file name>    play sound file
        stop                stop music
        pause               pause music
        resume              resume music
        volume <new vol>    change volume (float)

		help                this is it
        exit                finish audigo REPL
    `)
	fmt.Println()
}

func (c *playerRepl) printHead() {
	fmt.Print("audigo >>> ")
}

func (c *playerRepl) printf(f string, a ...interface{}) {
	c.printHead()
	fmt.Printf(f, a...)
}

func (c *playerRepl) println(s ...interface{}) {
	c.printHead()
	fmt.Println(s...)
}
